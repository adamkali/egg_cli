package modules

import (
	"fmt"
	"os/exec"

	"github.com/adamkali/egg_cli/configuration"
	"github.com/adamkali/egg_cli/models"
	"github.com/adamkali/egg_cli/pkg/targets"
	"github.com/adamkali/egg_cli/styles"
)

type InstallLibrariesModule struct {
	eggl *models.EggLog
	Progress int
	Error error
}

func (*InstallLibrariesModule) Name() string {
	return "egg::install_libraries"
	
}

var maxprog_modules = float64(len(targets.GolangPackages))

func (m *InstallLibrariesModule) GetProgress() float64{
	return float64(m.Progress ) / maxprog_modules
}

// incrprog increments the progress by 1
func (m *InstallLibrariesModule) IncrProg() {
	m.Progress += 1
	return
}

func (m *InstallLibrariesModule) Run() {
	installLibrariesStart := styles.EggProgressInfo.Render("🥚 " +m.Name() + " start")
	fmt.Println(installLibrariesStart)
	for _, pac := range targets.GolangPackages{
		installLibrariesMessage := fmt.Sprintf(
			"🥚 %s installing %s",
			m.Name(),
			pac,
		)
		m.eggl.Info(installLibrariesMessage)
		installLibrariesMessage = styles.EggProgressInfo.Render(installLibrariesMessage)
		fmt.Println(installLibrariesMessage)
		err := m.GoGet(pac)
		if err != nil {
			return
		}
	}
}

func (m *InstallLibrariesModule) GoGet(pac string) error {
	_, err := exec.Command(
		"go", "get", pac).Output()
	if err != nil {
		m.Error = err
	}
	return m.Error
}

func (m *InstallLibrariesModule) IsError() error {
	return m.Error
}

func (m *InstallLibrariesModule) LoadFromConfig(_ *configuration.Configuration, eggl *models.EggLog) {
	m.eggl = eggl 
	m.Progress = 0
	m.eggl.Info("Installing libraries")
	return
}
