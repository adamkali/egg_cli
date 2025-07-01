package modules

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/adamkali/egg_cli/configuration"
	"github.com/adamkali/egg_cli/models"
	"github.com/adamkali/egg_cli/styles"
)

type BootstrapDirectoriesModule struct {
	Directories []string
	Error       error
	Progress    int
	eggl        *models.EggLog
}

func (m *BootstrapDirectoriesModule) Name() string {
	return "egg::bootstrap_directories"
}

func (m *BootstrapDirectoriesModule) IncrProg() {
	m.Progress += 1
	return
}

func (m *BootstrapDirectoriesModule) GetProgress() float64 {
	return float64(m.Progress) / float64(len(m.Directories))
}

func (m *BootstrapDirectoriesModule) Run() {
	bootstrapDirectoriesStart := styles.EggProgressInfo.Render("🥚 " + m.Name() + " start")
	fmt.Println(bootstrapDirectoriesStart)

	for _, dir := range m.Directories {
		bootstrapDirectoriesMessage := styles.EggProgressInfo.Render(fmt.Sprintf(
			"🥚 %s creating %s",
			m.Name(),
			dir,
		))
		fmt.Println(bootstrapDirectoriesMessage)
		m.mkdir(dir)
		if m.IsError() != nil {
			return
		}
	}

	return
}

func (m *BootstrapDirectoriesModule) mkdir(dir string) {
	_, err := exec.Command("mkdir", "-p", dir).Output()
	// if err is not File Exists, return the error
	if err != nil {
		if err.Error() != "exit status 1" {
			err = errors.New("error creating directory: " + dir + " " + err.Error())
			m.Error = err
		}
	}
}

func (m *BootstrapDirectoriesModule) IsError() error {
	return m.Error
}

func (m *BootstrapDirectoriesModule) LoadFromConfig(configuration *configuration.Configuration, eggl *models.EggLog) {
	m.Directories = []string{
		"cmd/configuration",
		"cmd/database_cli",
		"controllers",
		"docs",
		"middlwares/configs",
		"models",
		"public",
		"services",
		"tmp",
	}
	m.Directories = append(m.Directories, configuration.Database.Migration.Destination)
	m.Directories = append(m.Directories, configuration.Database.QueriesLocation)
	m.Directories = append(m.Directories, configuration.Database.SqlcRepositoryLocation)
	m.eggl = eggl
	m.Progress = 0
	return
}
