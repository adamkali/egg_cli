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
	return "egg_cli::post_install_run_migration"
}

func (m *PostInstallPushMigrationsToServer) LoadFromConfig(config *configuration.Configuration, eggl *models.EggLog) {
	m.GOOSE_DRIVER = config.Database.Migration.Protocol
	m.GOOSE_DBSTRING = config.Database.URL
	if m.GOOSE_DRIVER == "turso" && config.Database.TursoAuthToken != "" {
		m.GOOSE_DBSTRING = m.GOOSE_DBSTRING + "?authToken=" + config.Database.TursoAuthToken
	}
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
		cmd, err := exec.Command("goose", "up", "-v").Output()
		if err != nil {
			m.Error = err
			m.eggl.Error("error: %s", m.Error.Error())
			return
		}
		log := fmt.Sprintf("🥚 %s %s", m.Name(), string(cmd))
		m.eggl.Info(log)
		fmt.Println(styles.EggProgressInfo.Render(log))
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
