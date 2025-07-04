package pkg

import (
	"fmt"
	"os"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/pkg/modules"
	"github.com/adamkali/egg_cli/styles"
	"gopkg.in/yaml.v3"
)

const (
	ScrambledFileName = ".scrambled"
)

var (
	Modules = []modules.IModule{
		&modules.InitializeModule{},
		&modules.InstallToolsModule{},
		&modules.InstallLibrariesModule{},
		&modules.BootstrapDirectoriesModule{},
		&modules.GenerateConfigurationModule{},
		&modules.BootstrapFrameworkFilesFromTemplatesModule{},
		&modules.RsbuildFrontendModule{},
	}
)

type scrambleFile struct {
	Failed struct {
		ModuleName string `yaml:"moduleName"`
		Error      string `yaml:"error"`
	}
	Succeeded     []string                    `yaml:"succeeded"`
	Configuration configuration.Configuration `yaml:"configuration"`
}

func PrintError(m modules.IModule, eggl *models.EggLog) bool {
	if m.IsError() != nil {
		styles.EggProgressError.Render(fmt.Sprintf("🥚 %s encountered error: %v", m.Name(), m.IsError().Error()))
		eggl.Error("error: %s", m.IsError().Error())
		return true
	}
	return false
}

func ProjectFactory(configuration *configuration.Configuration, eggl *models.EggLog) error {
	var err error
	var succeededModules []modules.IModule
	for _, module := range Modules {
		module.LoadFromConfig(configuration, eggl)
		module.Run()
		err = module.IsError()
		if PrintError(module, eggl) {
			// if there is an error, then we to still write the .scrambled file
			err = WriteScrambled(configuration, succeededModules, module, err)
			// print out to check the Scrambled file
			fmt.Println("check the .scrambled for the breaking error and what module failed")
			return err
		}
		succeededModules = append(succeededModules, module)
	}
	return nil
}

func RecoverFromScrambled(eggl *models.EggLog) error {
	configuration, succeededModules, failedModules, err := LoadScrambled()
	if err != nil {
		return fmt.Errorf("failed to load .scrambled file: %w", err)
	}

	if len(failedModules) == 0 {
		eggl.Info("No failed modules to recover - all modules completed successfully")
		return nil
	}

	eggl.Info("Recovering %d failed modules", len(failedModules))

	for _, module := range failedModules {
		eggl.Info("Attempting to recover module: %s", module.Name())

		module.LoadFromConfig(configuration, eggl)
		module.Run()
		err = module.IsError()

		if PrintError(module, eggl) {
			// if there is an error, then we have to write the .scrambled file
			bigErr := WriteScrambled(configuration, succeededModules, module, err)
			if bigErr != nil {
				return fmt.Errorf("failed to write .scrambled file: %w", bigErr)
			}
			// print out to check the Scrambled file
			fmt.Println("check the .scrambled for the breaking error and what module failed")
			return err
		}

		eggl.Info("Successfully recovered module: %s", module.Name())
		succeededModules = append(succeededModules, module)
	}

	eggl.Info("All modules recovered successfully")
	return nil
}

func CheckScrambled() bool {
	_, err := os.Stat(ScrambledFileName)
	return !os.IsNotExist(err)
}

func LoadScrambled() (
	*configuration.Configuration,
	[]modules.IModule,
	[]modules.IModule,
	error,
) {
	// open the .scrambled file
	file, err := os.Open(ScrambledFileName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open .scrambled file: %w", err)
	}
	defer file.Close()

	// read the .scrambled file
	var scramble scrambleFile
	err = yaml.NewDecoder(file).Decode(&scramble)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode .scrambled file: %w", err)
	}

	// Create a map for O(1) lookup of succeeded modules
	succeededMap := make(map[string]bool)
	succeededModules := make([]modules.IModule, 0, len(scramble.Succeeded))

	for _, moduleName := range scramble.Succeeded {
		module := modules.ModuleFactory(moduleName)
		if module == nil {
			return nil, nil, nil, fmt.Errorf("module not found: %s", moduleName)
		}
		succeededModules = append(succeededModules, module)
		succeededMap[moduleName] = true
	}

	// Find failed modules using efficient map lookup
	var failedModules []modules.IModule
	for _, module := range Modules {
		if !succeededMap[module.Name()] {
			failedModules = append(failedModules, module)
		}
	}

	return &scramble.Configuration, succeededModules, failedModules, nil
}

// WriteScrambled
//
// params:
//
//	configuration:
//	  type: *configuration.Configuration
//	  description:
//	    the configuration of the project to be written to the .scrambled file
//	    if so that if there was any errors during the project creation we can
//	    recreate the state of the project from the .scrambled file
//	succeeded:
//	  type: []modules.IModule
//	  description:
//	    the modules that were successfully run and can be written to the
//	    .scrambled file so that we do not need to repeat them
//	failed:
//	  type: modules.IModule
//	  description:
//	    the module that failed and can be written to the .scrambled file
//	err:
//	  type: error
//	  description:
//	    the error that occurred during the project creation
//
// returns:
//
//	error:
//	  - Any error that occurs while writing the .scrambled file
//
// description:
//
//		This function is used to write the .scrambled file to the current directory
//	 if there were any errors during the project creation. The .scrambled file
//	 contains the modules that were successfully run and the module that failed
//	 this allows us to recreate the state of the project from the .scrambled file
//	 and re-run the failed module by trusting that the rest of the modules were
//	 correct
func WriteScrambled(
	configuration *configuration.Configuration,
	succeeded []modules.IModule,
	failed modules.IModule,
	ModuleError error,
) error {
	// Create or truncate the .scrambled file
	f, err := os.Create(ScrambledFileName)
	if err != nil {
		return fmt.Errorf("failed to create .scrambled file: %w", err)
	}
	defer f.Close()

	// Build succeeded module names slice
	succeededNames := make([]string, len(succeeded))
	for i, module := range succeeded {
		succeededNames[i] = module.Name()
	}

	// Create scramble file structure
	scr := scrambleFile{
		Failed: struct {
			ModuleName string `yaml:"moduleName"`
			Error      string `yaml:"error"`
		}{
			ModuleName: failed.Name(),
			Error:      ModuleError.Error(),
		},
		Succeeded:     succeededNames,
		Configuration: *configuration,
	}

	// Marshal the scramble file
	scrambled, err := yaml.Marshal(scr)
	if err != nil {
		return fmt.Errorf("failed to marshal scramble file: %w", err)
	}

	// Write the scramble file
	if _, err := f.Write(scrambled); err != nil {
		return fmt.Errorf("failed to write .scrambled file: %w", err)
	}

	return nil
}
