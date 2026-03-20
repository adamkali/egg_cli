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
	"github.com/adamkali/egg_cli/pkg/targets"
	"github.com/adamkali/egg_cli/styles"
)

type InstallLibrariesModule struct {
	eggl      *models.EggLog
	Progress  int
	Error     error
	GoGetFunc func(pac string) error // For testing - can be injected to mock go get
}

func (*InstallLibrariesModule) Name() string {
	return "egg::install_libraries"

}

var maxprog_modules = float64(len(targets.GolangPackages))

func (m *InstallLibrariesModule) GetProgress() float64 {
	return float64(m.Progress) / maxprog_modules
}

// incrprog increments the progress by 1
func (m *InstallLibrariesModule) IncrProg() {
	m.Progress += 1
	return
}

func (m *InstallLibrariesModule) Run() {
	installLibrariesStart := styles.EggProgressInfo.Render("🥚 " + m.Name() + " start")
	fmt.Println(installLibrariesStart)
	for _, pac := range targets.GolangPackages {
		installLibrariesMessage := fmt.Sprintf(
			"🥚 %s installing %s",
			m.Name(),
			pac,
		)
		m.eggl.Info(installLibrariesMessage)
		installLibrariesMessage = styles.EggProgressInfo.Render(installLibrariesMessage)
		fmt.Println(installLibrariesMessage)
		err := m.GoGet(pac)
		if err != nil {
			return
		}
	}
}

func (m *InstallLibrariesModule) GoGet(pac string) error {
	// Use injected function if available (for testing), otherwise use real implementation
	if m.GoGetFunc != nil {
		return m.GoGetFunc(pac)
	}

	_, err := exec.Command(
		"go", "get", pac).Output()
	if err != nil {
		m.Error = err
	}
	return m.Error
}

func (m *InstallLibrariesModule) IsError() error {
	return m.Error
}

func (m *InstallLibrariesModule) LoadFromConfig(_ *configuration.Configuration, eggl *models.EggLog) {
	m.eggl = eggl
	m.Progress = 0
	m.eggl.Info("Installing libraries")
	return
}
