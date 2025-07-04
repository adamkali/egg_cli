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
	case "egg::install-tools":
		return &InstallToolsModule{}
	case "egg::install-libraries":
		return &InstallLibrariesModule{}
	case "egg::bootstrap-directories":
		return &BootstrapDirectoriesModule{}
	case "egg::generate-configuration":
		return &GenerateConfigurationModule{}
	case "egg::bootstrap-framework-files-from-templates":
		return &BootstrapFrameworkFilesFromTemplatesModule{}
	case "egg::rsbuild-frontend":
		return &RsbuildFrontendModule{}
	default:
		return nil
	}
}
