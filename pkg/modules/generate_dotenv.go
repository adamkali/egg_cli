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
	"text/template"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/pkg/templates"
	"github.com/adamkali/egg_cli/styles"
)

type GenerateDotenvModule struct {
	EggVersion      string
	Output         string
	Config         *configuration.GeneratorConfiguration
	Error          error
	eggl           *models.EggLog
	GenerateDotenv func() error
}

func (g *GenerateDotenvModule) Name() string {
	return "egg::generate_dotenv"
}

func (g *GenerateDotenvModule) IsError() error {
	return g.Error
}

const fpath = ".env"

func (g *GenerateDotenvModule) Run() {
	log := fmt.Sprintf("🥚 %s creating %s", g.Name(), g.GetOutput())
	fmt.Println(styles.EggProgressInfo.Render(log))
	g.eggl.Info(log)
	if g.GenerateDotenv != nil {
		g.Error = g.GenerateDotenv()
	} else {
		templatesFS := templates.GetTemplates()
		// load the template from ../templates/.env.tmpl
		templateContent, err := templatesFS.ReadFile(".env.tmpl")
		if err != nil {
			g.Error = err
		}
		templateObject, err := template.New(g.GetOutput()).Parse(string(templateContent))
		if err != nil {
			g.Error = err
		}

		fileWriter, err := os.Create(g.GetOutput())
		if err != nil {
			g.Error = err
		}
		defer fileWriter.Close()

		err = templateObject.Execute(fileWriter, g.Config)
		if err != nil {
			g.Error = err
		}
	}
}

func (g *GenerateDotenvModule) GetOutput() string {
	return g.Output
}

func (g *GenerateDotenvModule) SetOutput(output string) {
	if output == "" {
		g.Output = fpath
	}
	g.Output = output
}

func (g *GenerateDotenvModule) GetEggVersion() string {
	return g.EggVersion
}

func (g *GenerateDotenvModule) SetEggVersion(eggVersion string) {
	g.EggVersion = eggVersion
}

func (g *GenerateDotenvModule) LoadFromConfig(
	config *configuration.Configuration,
	eggl *models.EggLog,
) {
	g.Config = config.IntoGeneratorConfig(g.GetOutput())
	g.eggl = eggl
}
