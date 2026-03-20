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
	"fmt"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/styles"
)

type GenerateConfigurationModule struct {
	Configuration      *configuration.Configuration
	Error              error
	Progress           int
	eggl               *models.EggLog
	GenerateConfigFunc func(environment string) error // For testing - can be injected to mock config generation
}

func (m *GenerateConfigurationModule) Name() string   { return "egg::generate_configuration" }
func (m *GenerateConfigurationModule) IsError() error { return m.Error }

func (m *GenerateConfigurationModule) Run() {
	generateConfigurationStart := "🥚 " + m.Name() + " start"
	m.eggl.Info(generateConfigurationStart)
	generateConfigurationStart = styles.EggProgressInfo.Render(generateConfigurationStart)
	fmt.Println(generateConfigurationStart)
	var err error
	if m.GenerateConfigFunc != nil {
		err = m.GenerateConfigFunc("development")
	} else {
		err = m.Configuration.GenerateConfigurationFile("development")
	}
	if err != nil {
		m.Error = err
		m.eggl.Error("error: %s", m.Error.Error())
		return
	}
	m.eggl.Info("🥚 " + m.Name() + " complete")
	generateConfigurationComplete := styles.EggProgressInfo.Render("🥚 " + m.Name() + " complete")
	fmt.Println(generateConfigurationComplete)
}

func (m *GenerateConfigurationModule) LoadFromConfig(
	configuration *configuration.Configuration,
	eggl *models.EggLog,
) {
	m.Configuration = configuration
	m.eggl = eggl
	m.Progress = 0
	return
}
