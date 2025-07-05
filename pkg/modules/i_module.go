package modules

import (
	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
)

type IModule interface {
	Run()
	Name() string
	IsError() error
	LoadFromConfig(configuration *configuration.Configuration, eggl *models.EggLog)
}

func ModuleFactory(moduleName string) IModule {
	switch moduleName {
	case "egg::initialize":
		return &InitializeModule{}
	case "egg::install_tools":
		return &InstallToolsModule{}
	case "egg::install_libraries":
		return &InstallLibrariesModule{}
	case "egg::bootstrap_directories":
		return &BootstrapDirectoriesModule{}
	case "egg::generate_configuration":
		return &GenerateConfigurationModule{}
	case "egg::bootstrap_framework":
		return &BootstrapFrameworkFilesFromTemplatesModule{}
	case "egg::rsbuild_frontend":
		return &RsbuildFrontendModule{}
	default:
		return nil
	}
}
