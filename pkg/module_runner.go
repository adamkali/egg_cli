package pkg

import (
	"errors"
	"fmt"
	"os"

	"github.com/adamkali/egg_cli/models"
	"github.com/adamkali/egg_cli/pkg/configuration"
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
		panic(err)
	}
	for _, module := range failedModules {
		module.LoadFromConfig(configuration, nil)
		module.Run()
		err = module.IsError()
		if PrintError(module, nil) {
			// if there is an error, then we have to write the .scrambled file
			bigErr := WriteScrambled(configuration, succeededModules, module, err)
			if bigErr != nil {
				return bigErr
			}
			// print out to check the Scrambled file
			fmt.Println("check the .scrambled for the breaking error and what module failed")
			return err
		}
		succeededModules = append(succeededModules, module)
	}
	return nil
}

func CheckScrambled() bool {
	// check if there is a .scrambled file in the current directory
	_, err := os.Stat(ScrambledFileName)
	if os.IsNotExist(err) {
		// if the .scrambled file does not exist, return false
		return false
	}
	// if the .scrambled file exists, return true
	return true
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
		return nil, nil, nil, err
	}
	// read the .scrambled file
	var scramble scrambleFile
	err = yaml.NewDecoder(file).Decode(&scramble)
	if err != nil {
		return nil, nil, nil, err
	}
	succeededModules := []modules.IModule{}
	for _, moduleName := range scramble.Succeeded {
		if modules.ModuleFactory(moduleName) == nil {
			return nil, nil, nil, errors.New("module not found")
		}
		succeededModules = append(succeededModules, modules.ModuleFactory(moduleName))
	}

	// cut out all the succeeded modules from Modules list
	// and take the remaining modules as failed modules
	var failedModules []modules.IModule
	for _, module := range Modules {
		found := false
		for _, succeededModule := range succeededModules {
			if module.Name() == succeededModule.Name() {
				found = true
				break
			}
		}
		if !found {
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
	var f *os.File
	// create a .scrambled file in the current directory
	if _, err := os.Stat(ScrambledFileName); errors.Is(err, os.ErrNotExist) {
		if f, err = os.Create(ScrambledFileName); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(ScrambledFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	// init the scramble file
	succeededNames := make([]string, len(succeeded))
	for _, module := range succeeded {
		succeededNames = append(succeededNames, module.Name())
	}
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
	// marshal the scramble file
	scrambled, err := yaml.Marshal(scr)
	if err != nil {
		return err
	}
	// write the scramble file
	if _, err := f.Write(scrambled); err != nil {
		return err
	}
	return nil
	// if the .scrambled file already exists then
}
