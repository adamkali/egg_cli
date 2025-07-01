package pkg

import (
	"fmt"

	"github.com/adamkali/egg_cli/configuration"
	"github.com/adamkali/egg_cli/models"
	"github.com/adamkali/egg_cli/pkg/modules"
	"github.com/adamkali/egg_cli/styles"
)

var (
	Modules = []modules.IModule{
		&modules.InitializeModule{},
		&modules.InstallToolsModule{},
		&modules.InstallLibrariesModule{},
		&modules.BootstrapDirectoriesModule{},
	}
)

func PrintError(m modules.IModule, eggl *models.EggLog) bool{
	if m.IsError() != nil {
		styles.EggProgressError.Render(fmt.Sprintf("🥚 %s encountered error: %v", m.Name(), m.IsError().Error()))
		eggl.Error("error: %s", m.IsError().Error())
		return true
	}
	return false
}


func ProjectFactory(configuration *configuration.Configuration, eggl *models.EggLog) error {
	var err error
	for _, module := range Modules {
		module.LoadFromConfig(configuration, eggl)
		module.Run()
		err = module.IsError()
		if PrintError(module, eggl) {
			return err
		}
	}
	return nil
}

