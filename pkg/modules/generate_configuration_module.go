package modules

import (
	"fmt"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/styles"
)

type GenerateConfigurationModule struct {
	Configuration *configuration.Configuration
	Error         error
	Progress      int
	eggl          *models.EggLog
}

func (m *GenerateConfigurationModule) Name() string   { return "egg::generate_configuration" }
func (m *GenerateConfigurationModule) IsError() error { return m.Error }

func (m *GenerateConfigurationModule) Run() {
	generateConfigurationStart := "🥚 " + m.Name() + " start\n"
	m.eggl.Info(generateConfigurationStart)
	generateConfigurationStart = styles.EggProgressInfo.Render(generateConfigurationStart)
	fmt.Println(generateConfigurationStart)
	err := m.Configuration.GenerateConfigurationFile("development")
	if err != nil {
		m.Error = err
		m.eggl.Error("error: %s", m.Error.Error())
		return
	}
	m.eggl.Info("🥚 " + m.Name() + " complete")
	generateConfigurationComplete := styles.EggProgressInfo.Render("🥚 " + m.Name() + " complete\n")
	fmt.Println(generateConfigurationComplete)
}

func (m *GenerateConfigurationModule) LoadFromConfig(
	configuration *configuration.Configuration,
	eggl *models.EggLog,
) {
	m.Configuration = configuration
	m.eggl = eggl
	m.Progress = 0
	return
}
