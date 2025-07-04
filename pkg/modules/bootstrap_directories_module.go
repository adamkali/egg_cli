package modules

import (
	"errors"
	"fmt"
	"os/exec"
	"sync"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/styles"
)

type BootstrapDirectoriesModule struct {
	Directories []string
	Error       error
	Progress    int
	eggl        *models.EggLog
}

// Name
//
// returns:
//
//	string: the name
//
// description:
//
//	This function returns the name of the module
func (m *BootstrapDirectoriesModule) Name() string {
	return "egg::bootstrap_directories"
}

// IncrProg
//
// description:
//
//	increments the progress
func (m *BootstrapDirectoriesModule) IncrProg() {
	m.Progress += 1
	return
}

// GetProgress
//
// returns:
//
//	float64: the progress
func (m *BootstrapDirectoriesModule) GetProgress() float64 {
	return float64(m.Progress) / float64(len(m.Directories))
}

// Run
//
// description:
//
//		This function is used to bootstrap the directories for the project.
//	  this is done by iterating over the directories and creating a goroutine for each directory
//	  and then waiting for all the goroutines to finish
//	  and collecting the errors
//	  and logging the errors
//	  and canceling the context if there is an error
func (m *BootstrapDirectoriesModule) Run() {
	bootstrapDirectoriesStart := styles.EggProgressInfo.Render("🥚 " + m.Name() + " start")
	fmt.Println(bootstrapDirectoriesStart)
	errChan := make(chan error)
	logChan := make(chan string)
	wg := new(sync.WaitGroup)

	for _, dir := range m.Directories {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			err := m.mkdir(dir)
			logChan <- fmt.Sprintf("🥚 %s creating %s", m.Name(), dir)
			if err != nil {
				errChan <- err
			}
		}(dir)
	}

	go func() {
		for log := range logChan {
			m.eggl.Info(log)
			log = styles.EggProgressInfo.Render(log)
			fmt.Println(log)
		}
	}()
	go func() {
		for err := range errChan {
			m.Error = err
			return
		}
	}()
	go func() {
		wg.Wait()
		bootstrapDirectoriesComplete := styles.EggProgressInfo.Render("🥚 " + m.Name() + " complete")
		fmt.Println(bootstrapDirectoriesComplete)
		m.eggl.Info("🥚 " + m.Name() + " complete")
	}()
	return
}

// mkdir
//
// params:
//
//	dir: string
//
// returns:
//
//	error: if there is an error creating the directory
//
// description:
//
//	This function is used to create a directory using the builtin mkdir command
//	and passing in the -p flag to create the entire parent directory if it does not exist
func (m *BootstrapDirectoriesModule) mkdir(dir string) error {
	_, err := exec.Command("mkdir", "-p", dir).Output()
	// if err is not File Exists, return the error
	if err != nil {
		if err.Error() != "exit status 1" {
			return errors.New("error creating directory: " + dir + " " + err.Error())
		}
	}
	return nil
}

// IsError
//
// returns:
//
//	error: The stored error in the monad.
func (m *BootstrapDirectoriesModule) IsError() error {
	return m.Error
}

// LoadFromConfig
//
// params:
//
//	configuration: *configuration.Configuration
//	eggl: *models.EggLog
//
// description:
//
//	This function is used to load the module from the configuration and initialize the module
//	with the needed data
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
