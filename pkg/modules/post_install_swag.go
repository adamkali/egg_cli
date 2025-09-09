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
	"os/exec"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/styles"
)

type PostInstallSwagModule struct {
	eggl         *models.EggLog
	Progress     int
	Error        error
	SwagGenerate func() error
}

func (m *PostInstallSwagModule) Name() string {
	return "post_install_swag"
}

func (m *PostInstallSwagModule) LoadFromConfig(_ *configuration.Configuration, eggl *models.EggLog) {
	m.eggl = eggl
	m.Progress = 0
	m.eggl.Info("Installing swag")
	return
}

func (m *PostInstallSwagModule) Run() {
	installSwagMessage := fmt.Sprintf(
		"🥚 %s %s",
		m.Name(),
		"Running swag generate",
	)
	m.eggl.Info(installSwagMessage)
	installSwagMessage = styles.EggProgressInfo.Render(installSwagMessage)
	fmt.Println(installSwagMessage)
	if m.SwagGenerate == nil {
		// run swag generate
		cmd := exec.Command("swag", "init")
		err := cmd.Run()
		if err != nil {
			m.Error = err
			m.eggl.Error("error: %s", m.Error.Error())
			return
		}
	} else {
		err := m.SwagGenerate()
		if err != nil {
			m.Error = err
			m.eggl.Error("error: %s", m.Error.Error())
			return
		}

	}

	m.eggl.Info("🥚 " + m.Name() + " complete")
}

func (m *PostInstallSwagModule) GetProgress() float64 {
	return float64(m.Progress)
}

func (m *PostInstallSwagModule) IncrProg() {
	m.Progress += 1
	return
}

func (m *PostInstallSwagModule) IsError() error {
	return m.Error
}
