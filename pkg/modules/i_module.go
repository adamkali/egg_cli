/*
Copyright © 2025 Adam Kalinowski <adam.kalilarosa@proton.me>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	case "egg::clone_template":
		return &CloneTemplateModule{}
	case "egg::bootstrap_framework":
		return &BootstrapFrameworkFilesFromTemplatesModule{}
	case "egg::post_install_sqlc":
		return &PostInstallSqlcModule{}
	case "egg::post_install_push_migrations":
		return &PostInstallPushMigrationsToServer{}
	case "egg::post_install_swag":
		return &PostInstallSwagModule{}
	case "egg::complete":
		return CompleteModuleInstance()
	default:
		return nil
	}
}
