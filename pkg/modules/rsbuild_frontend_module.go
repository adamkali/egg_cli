// This module is used to build the frontend if the user chooses to do so
// and it will be used in the m.Run() function. EggCli does not give a hoot
// about which frontend framework the user is using, so this module will assume
// that the user will implement the frontend framework of their choice, and integrate
//
// config.Server.Frontend.Dir to hold the {config.Name}/{config.Server.Frontend.Dir} directory
//
//	as their distribution directory ie the result of running npm build
//	e.g. egg_app/web/dist
//
// config.Server.Frontend.API to hold the {config.Name}/{config.Server.Frontend.API} directory
//
//	as their api directory generated by ./tmp/{config.Name} swag
//	e.g. egg_app/web/api
package modules

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
)

const (
	PnpmInstall = "pnpm create rsbuild@latest"
	NpmInstall  = "npm create rsbuild@latest"
	YarnInstall = "yarn create rsbuild"
	BunInstall  = "bun create rsbuild"
)

// RsbuildFrontendModule
//
// implements:
//   - IModule
//
// description:
//
//	This struct is used to build the frontend if the user chooses to do so
type RsbuildFrontendModule struct {
	configuration *configuration.Configuration
	error         error
	progress      int
	eggl          *models.EggLog
	InputFunc     func(prompt string) string                       // For testing - can be injected to mock user input
	ExecFunc      func(cmd string, args ...string) ([]byte, error) // For testing - can be injected to mock command execution
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
func (m *RsbuildFrontendModule) Name() string {
	return "egg::rsbuild_frontend"
}

// IncrProg
//
// description:
//
//	increments the progress
func (m *RsbuildFrontendModule) IncrProg() {
	m.progress += 1
}

// GetProgress
//
// returns:
//
//	float64: the progress
func (m *RsbuildFrontendModule) GetProgress() float64 {
	return float64(m.progress)
}

// Run
//
// description:
//
//		This function is used to build the frontend and uses the following command as a reference
//	  pnpm create rsbuild@latest || npm create rsbuild@latest || yarn create rsbuild || bun create rsbuild
//	  we firrst ask the user which frontend framework they are using and then we build the frontend and
//	  passing off the control flow to rsbuild for the user to interact with the cli interface
func (m *RsbuildFrontendModule) Run() {
	// ask the user if they want to use js for the frontend
	// if they do not want to use js for the frontend, exit
	var useJs string
	if m.InputFunc != nil {
		useJs = m.InputFunc("Do you want to use a JavaScript Framework for the frontend? (y/n)")
	} else {
		fmt.Println("Do you want to use a JavaScript Framework for the frontend? (y/n)")
		fmt.Scanln(&useJs)
	}
	if useJs == "n" || useJs == "N" || useJs == "no" || useJs == "No" || useJs == "NO" {
		return
	}

	// ask which package manager do they want to use?
	var packageManager string
	if m.InputFunc != nil {
		packageManager = m.InputFunc("Which package manager do you want to use? (pnpm, npm, yarn, bun)")
	} else {
		fmt.Println("Which package manager do you want to use? (pnpm, npm, yarn, bun)")
		fmt.Scanln(&packageManager)
	}
	switch packageManager {
	case "pnpm":
		m.installAndWaitForRsBuild(PnpmInstall)
	case "p":
		m.installAndWaitForRsBuild(PnpmInstall)
	case "npm":
		m.installAndWaitForRsBuild(NpmInstall)
	case "n":
		m.installAndWaitForRsBuild(NpmInstall)
	case "yarn":
		m.installAndWaitForRsBuild(YarnInstall)
	case "y":
		m.installAndWaitForRsBuild(YarnInstall)
	case "bun":
		m.installAndWaitForRsBuild(BunInstall)
	case "b":
		m.installAndWaitForRsBuild(BunInstall)
	default:
		m.error = fmt.Errorf("invalid package manager: %s", packageManager)
		m.eggl.Error("error: %s", m.error.Error())
		return
	}
	return
}

// installAndWaitForRsBuild
//
// params:
//
//	string: the package manager
//
// returns:
//
//	error: any errors that occur to install the frontend
//
// description:
//
//			This function is used to install the frontend framework and wait for rsbuild
//	     to finish building the frontend. It will wait for rsbuild to take control of
//	     stdin and stdout and then pass off the control flow to rsbuild. then when it
//	     is done, it will return the error
func (m *RsbuildFrontendModule) installAndWaitForRsBuild(packageManager string) {
	p := strings.Split(packageManager, " ")[0]
	rest := strings.Split(packageManager, " ")[1:]
	var err error
	if m.ExecFunc != nil {
		_, err = m.ExecFunc(p, rest...)
	} else {
		_, err = exec.LookPath(p)
		if err != nil {
			m.error = err
			m.eggl.Error("error: %s", m.error.Error())
			return
		}
		var output []byte
		output, m.error = exec.Command(packageManager, rest...).Output()
		if m.error != nil {
			m.eggl.Error("error: %s", m.error.Error())
			return
		}
		fmt.Println(string(output))
	}
	if err != nil {
		m.error = err
		m.eggl.Error("error: %s", m.error.Error())
		return
	}
	return
}

// IsError
//
// returns:
//
//	error: the error
func (m *RsbuildFrontendModule) IsError() error {
	return m.error
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
//	This function is used to load the module from the configuration and sets up the empty structure so that it can be used
//	by m.Run().
func (m *RsbuildFrontendModule) LoadFromConfig(configuration *configuration.Configuration, eggl *models.EggLog) {
	m.configuration = configuration
	m.eggl = eggl
	return
}
