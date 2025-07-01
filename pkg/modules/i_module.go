package modules

import (
	"github.com/adamkali/egg_cli/configuration"
	"github.com/adamkali/egg_cli/models"
)

type IModule interface {
	Run()
	Name() string
	IsError() error
	LoadFromConfig(configuration *configuration.Configuration, eggl *models.EggLog)
}
