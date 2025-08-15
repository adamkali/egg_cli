package modules

import (
	"fmt"
	"os/exec"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/styles"
)

type PostInstallSqlcModule struct {
	eggl         *models.EggLog
	Progress     int
	Error        error
	SwagGenerate func() error // For testing - can be injected to mock go get
}

func (m *PostInstallSqlcModule) Name() string {
	return "egg::post_install_sqlc"
}

func (m *PostInstallSqlcModule) GetProgress() float64 {
	return float64(m.Progress) / maxprog_modules
}

// incrprog increments the progress by 1
func (m *PostInstallSqlcModule) IncrProg() {
	m.Progress += 1
	return
}

func (m *PostInstallSqlcModule) Run() {
	installSqlcMessage := fmt.Sprintf(
		"🥚 %s %s",
		m.Name(),
		"Running sqlc generate",
	)
	m.eggl.Info(installSqlcMessage)
	installSqlcMessage = styles.EggProgressInfo.Render(installSqlcMessage)
	fmt.Println(installSqlcMessage)
	if m.SwagGenerate != nil {
		err := m.SwagGenerate()
		if err != nil {
			m.Error = err
			m.eggl.Error("error: %s", m.Error.Error())
			return
		}
	} else {
		// run sqlc generate
		cmd := exec.Command("sqlc", "generate")
		err := cmd.Run()
		if err != nil {
			m.Error = err
			m.eggl.Error("error: %s", m.Error.Error())
			return
		}
	}

	m.eggl.Info("🥚 " + m.Name() + " complete")
	installSqlcComplete := styles.EggProgressInfo.Render("🥚 " + m.Name() + " complete\n")
	fmt.Println(installSqlcComplete)
}

func (m *PostInstallSqlcModule) IsError() error {
	return m.Error
}

func (m *PostInstallSqlcModule) LoadFromConfig(config *configuration.Configuration, eggl *models.EggLog) {
	m.eggl = eggl
	m.Progress = 0
	m.eggl.Info("🥚 generating %s using sqlc",  config.Database.Sqlc.RepositoryLocation)
	return
}

