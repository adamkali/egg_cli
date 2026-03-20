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

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/styles"
)

type CompleteModule struct {
	eggl               *models.EggLog
	Progress           int
	Error              error
	ProjectName        string
	ProjectFrontendDir string
	ProjectPort        int 
}

func (m *CompleteModule) Name() string {
	return "egg::complete"
}

// TODO: using lipgloss here color the egg 
//       logo to use egg based colors
const eggComplete = `
	           ████████████████        
	         ██                ██      
	     ████    ░░░░░░░░        ██    
	   ██      ░░      ░░░░        ██  
	 ██      ░░          ░░░░        ██
	 ██      ░░          ░░░░        ██
	██        ░░▒▒░░  ░░░░░░░░        ██
	██░░        ░░░░░░░░░░░░        ░░██
	  ██░░        ░░░░░░░░        ░░██  
	  ██░░░░                    ░░██    
	    ████░░░░            ░░░░██      
 	       ████░░░░░░░░░░░░████        
	           ████████████            
`

func (m *CompleteModule) Run() {
	installCompleteMessage := fmt.Sprintf(
		"🥚 %s %s",
		m.Name(),
		"Complete",
	)
	m.eggl.Info(installCompleteMessage)
	installCompleteMessage = styles.EggProgressInfo.Render(installCompleteMessage)
	fmt.Println(installCompleteMessage)

	// Now that sqlc has generated db/repository/, we can run go mod tidy
	// and verify the project builds cleanly.
	fmt.Println(styles.EggProgressInfo.Render("Running go mod tidy..."))
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Stdout = nil
	tidy.Stderr = nil
	if err := tidy.Run(); err != nil {
		m.Error = fmt.Errorf("go mod tidy failed: %w", err)
		m.eggl.Error("error: %s", m.Error.Error())
		return
	}

	fmt.Println(styles.EggProgressInfo.Render("Verifying project builds..."))
	build := exec.Command("go", "build", "./...")
	build.Stdout = nil
	build.Stderr = nil
	if err := build.Run(); err != nil {
		m.Error = fmt.Errorf("go build ./... failed: %w", err)
		m.eggl.Error("error: %s", m.Error.Error())
		return
	}
	fmt.Println(styles.EggProgressComplete.Render("Project builds successfully"))

	cmd := exec.Command("clear")
	err := cmd.Run()
	if err != nil {
		m.Error = err
		m.eggl.Error("error: %s", m.Error.Error())
		return
	}

	installCompleteMessageWhatToDoNext := fmt.Sprintf(
		`%s

		Now that the generator is complete, 
		$ cd ./ %s 
	
		and here is what you can do next:

		1.  Run "go run main.go" to start the server
		2.  Run "make build" to build the project
		3.  Run "air" to also start the server but with hot reloading
		4.  Create a frontend, and put it in the "web" directory &&
		    when you are done with the frontend, build it so that the dist
		    directory matches 
			config.Server.Frontend.Dir: %s
			config variable
		5.  Run the server and make sure that the database is workning by 
		    going to http://localhost:%d/swagger

		Have fun hacking!
		`,
		eggComplete,
		m.ProjectName,
		m.ProjectFrontendDir,
		m.ProjectPort,
	)
	m.eggl.Info(installCompleteMessageWhatToDoNext)
	installCompleteMessageWhatToDoNext = styles.EggProgressComplete.Render(installCompleteMessageWhatToDoNext)
	fmt.Println(installCompleteMessageWhatToDoNext)
}

func (m *CompleteModule) IsError() error {
	return m.Error
}

func (m *CompleteModule) GetProgress() int {
	return m.Progress
}

func (m *CompleteModule) LoadFromConfig(config *configuration.Configuration, eggl *models.EggLog) {
	m.eggl = eggl
	m.Progress = 0
	m.ProjectName = config.Name
	m.ProjectFrontendDir = config.Server.Frontend.Dir
	m.ProjectPort = config.Server.Port
	m.Error = nil
}

func CompleteModuleInstance() *CompleteModule { return &CompleteModule{} }
