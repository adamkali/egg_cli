package modules

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/models"
	"github.com/adamkali/egg_cli/styles"
)

//lipgloss "github.com/charmbracelet/lipgloss"

type InitializeModule struct {
	eggl        *models.EggLog
	Namespace   string
	ProjectName string
	License     string
	Progress    int
	Error       error
}

func (m *InitializeModule) Name() string {
	return "egg::initialize"
}

var maxprog_init = float64(4)

func (m *InitializeModule) GetProgress() float64 {
	return float64(m.Progress) / maxprog_init
}

// incrprog increments the progress by 1
func (m *InitializeModule) IncrProg() {
	m.Progress += 1
	return
}

func (m *InitializeModule) LoadFromConfig(configuration *configuration.Configuration, eggl *models.EggLog) {
	m.Namespace = configuration.Namespace
	m.ProjectName = configuration.Name
	m.License = configuration.License
	m.Progress = 0
	m.eggl = eggl
	return
}

// IsError checks if there is an error in the module.
// if there is none it returns nil
// if there is an error it returns the stored error
func (m InitializeModule) IsError() error {
	return m.Error
}

func (m *InitializeModule) Run() {
	// create a directory for the project
	initModuleStart := styles.EggProgressInfo.Render("🥚 " + m.Name() + " start\n")
	initModuleMkdirMessage := styles.EggProgressInfo.Render("🥚 " + m.Name() + " creating project root directory\n")
	initModuleGoVersionMessage := styles.EggProgressInfo.Render("🥚 " + m.Name() + " checking go version\n")
	initModuleGoModInitMessage := styles.EggProgressInfo.Render("🥚 " + m.Name() + " initializing go module\n")
	initModuleChangingDirectoryMessage := styles.EggProgressInfo.Render("🥚 " + m.Name() + " changing directory to " + m.ProjectName + "\n")
	initModuleCompletSuccessMessage := styles.EggProgressInfo.Render("🥚 " + m.Name() + " initialization complete\n")
	fmt.Println(initModuleStart)

	fmt.Println(initModuleMkdirMessage)
	output, err := exec.Command("mkdir", m.ProjectName).Output()
	if err != nil {
		if err.Error() == "exit status 1" {
			// ask the user if they want to overwrite the directory
			// if they do not want to overwrite the directory, return the error
			fmt.Println("project directory already exists, do you want to overwrite it? (y/n)")
			var overwrite string
			fmt.Scanln(&overwrite)
			if overwrite == "y" || overwrite == "Y" || overwrite == "yes" || overwrite == "Yes" || overwrite == "YES" {
				deleteErr := os.RemoveAll(m.ProjectName)
				if deleteErr != nil {
					m.Error = errors.New("error deleting project directory: " + m.ProjectName + " " + deleteErr.Error())
					m.eggl.Error("error: %s", m.Error.Error())
					return
				}
				// create the directory again
				exec.Command("mkdir", m.ProjectName).Output()
				m.Error = nil
			} else {
				m.Error = errors.New("project directory already exists")
				m.eggl.Error("error: %s", m.Error.Error())
				return
			}
		}
	}

	fmt.Println(initModuleGoVersionMessage)
	output, err = exec.Command("go", "version").Output()
	if err != nil {
		m.eggl.Error("error: %s", err.Error())
		m.Error = err
		return
	}
	m.Error = nil
	m.IncrProg()

	// go version returns:
	//go version go1.24.4 linux/amd64
	goVersion := strings.Split(string(output), " ")[2]
	goVersion = strings.Split(goVersion, ".")[1]

	goVersionF, err := strconv.ParseFloat(goVersion, 64)
	if err != nil {
		m.eggl.Error("error: %s", err.Error())
		m.Error = err
		return
	}
	if goVersionF < 23.0 {
		m.Error = errors.New("go version must be at least 1.23 current version: " + goVersion)
		m.eggl.Error("error: %s", m.Error.Error())
		return
	}
	m.Error = nil
	m.IncrProg()

	fmt.Println(initModuleChangingDirectoryMessage)
	err = os.Chdir(m.ProjectName)
	if err != nil {
		m.eggl.Error("error: %s", err.Error())
		m.Error = err
		return
	}
	m.Error = nil
	m.IncrProg()

	fmt.Println(initModuleGoModInitMessage)
	// go mod init
	// put it in the root project
	output, err = exec.Command(
		"go",
		"mod",
		"init",
		m.Namespace,
	).Output()
	if err != nil {
		m.eggl.Error("error: %s", err.Error())
		m.Error = err
		return
	}
	m.Error = nil
	m.IncrProg()

	fmt.Println(initModuleCompletSuccessMessage)
	return
}
