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
		cmd := exec.Command("swag", "generate")
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
