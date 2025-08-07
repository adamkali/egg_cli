package modules

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/styles"
)

type PostInstallPushMigrationsToServer struct {
	eggl                *models.EggLog
	Progress            int
	Error               error
	RunMigration        func() error
	GOOSE_DRIVER        string
	GOOSE_DBSTRING      string
	GOOSE_MIGRATION_DIR string
}

func (m *PostInstallPushMigrationsToServer) Name() string {
	return "post_install_run_migration"
}

func (m *PostInstallPushMigrationsToServer) LoadFromConfig(config *configuration.Configuration, eggl *models.EggLog) {
	m.GOOSE_DRIVER = config.Database.Migration.Protocol
	m.GOOSE_DBSTRING = config.Database.URL
	m.GOOSE_MIGRATION_DIR = config.Database.Migration.Destination

	m.eggl = eggl
	m.Progress = 0
	m.eggl.Info("Pushing migrations %s", m.GOOSE_DBSTRING)
}

func (m *PostInstallPushMigrationsToServer) Run() {

	installRunMigrationMessage := fmt.Sprintf(
		"🥚 %s %s",
		m.Name(),
		"Running migrations",
	)
	m.eggl.Info(installRunMigrationMessage)
	installRunMigrationMessage = styles.EggProgressInfo.Render(installRunMigrationMessage)
	fmt.Println(installRunMigrationMessage)

	if m.RunMigration != nil {
		err := m.RunMigration()
		if err != nil {
			m.Error = err
			m.eggl.Error("error: %s", m.Error.Error())
			return
		}
	} else {
		os.Setenv("GOOSE_DRIVER", m.GOOSE_DRIVER)
		os.Setenv("GOOSE_DBSTRING", m.GOOSE_DBSTRING)
		os.Setenv("GOOSE_MIGRATION_DIR", m.GOOSE_MIGRATION_DIR)
		cmd := exec.Command("goose", "up")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			m.Error = err
			m.eggl.Error("error: %s", m.Error.Error())
			return
		}
	}

	m.eggl.Info("🥚 " + m.Name() + " complete")
}

func (m *PostInstallPushMigrationsToServer) ProgressUpdate(progress int) {
	m.Progress = progress
}

func (m *PostInstallPushMigrationsToServer) GetProgress() float64 {
	return float64(m.Progress)
}

func (m *PostInstallPushMigrationsToServer) IsError() error {
	return m.Error
}

func (m *PostInstallPushMigrationsToServer) SetProgress(progress int) {
	m.Progress = progress
}
