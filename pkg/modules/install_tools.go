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

		installToolsMessage := fmt.Sprintf(
			"🥚 %s installing %s",
			m.Name(),
			tool,
		)
		m.eggl.Info(installToolsMessage)
		installToolsMessage = styles.EggProgressInfo.Render(installToolsMessage)
		fmt.Println(installToolsMessage)

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
