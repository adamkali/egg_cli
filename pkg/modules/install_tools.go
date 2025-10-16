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
	"strings"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/pkg/targets"
	"github.com/adamkali/egg_cli/styles"
)

type InstallToolsModule struct {
	eggl            *models.EggLog
	Progress        int
	Error           error
	LookPathFunc    func(file string) (string, error) // For testing - can be injected to mock tool lookup
	InstallToolFunc func(tool string) error           // For testing - can be injected to mock tool installation
}

func (m *InstallToolsModule) Name() string {
	return "egg::install_tools"
}

var maxprog_tools = float64(len(targets.RequiredTools))

func (m *InstallToolsModule) IncrProg() {
	m.Progress += 1
	return
}

func (m *InstallToolsModule) GetProgress() float64 {
	return float64(m.Progress) / maxprog_tools
}

func (m *InstallToolsModule) Run() {
	installToolsStart := styles.EggProgressInfo.Render("🥚 " + m.Name() + " start")
	fmt.Println(installToolsStart)

	// install go tools
	for _, tool := range targets.RequiredTools {
		toolStr := tool[strings.LastIndex(tool, "/")+1:]
		toolStr = toolStr[:strings.Index(toolStr, "@")]

		// Use injected function if available, otherwise use real implementation
		var err error
		if m.LookPathFunc != nil {
			_, err = m.LookPathFunc(tool)
		} else {
			_, err = exec.LookPath(tool)
		}

		if err == nil {
			installToolInstalledMessage := fmt.Sprintf(
				"🥚 %s %s is already installed",
				m.Name(),
				toolStr,
			)
			m.eggl.Info(installToolInstalledMessage)
			styles.EggProgressInfo.Render(installToolInstalledMessage)
			fmt.Println(installToolInstalledMessage)
			continue
		}
		// err should be handled at this point
		err = nil

		installToolsMessage := fmt.Sprintf(
			"🥚 %s installing %s",
			m.Name(),
			tool,
		)
		m.eggl.Info(installToolsMessage)
		installToolsMessage = styles.EggProgressInfo.Render(installToolsMessage)
		fmt.Println(installToolsMessage)
		// 

		// Use injected function if available, otherwise use real implementation
		if m.InstallToolFunc != nil {
			err = m.InstallToolFunc(tool)
		} else {
			output, err := exec.Command("go", "install", tool).Output()
			if err != nil {
				m.Error = err
				return
			}
			fmt.Println(string(output))
		}

	
		if err != nil {
			m.Error = err
			return
		}
		m.IncrProg()
		m.Error = nil
	}
}

func (m *InstallToolsModule) IsError() error {
	return m.Error
}

func (m *InstallToolsModule) LoadFromConfig(_ *configuration.Configuration, eggl *models.EggLog) {
	m.Progress = 0
	m.eggl = eggl
	return
}
