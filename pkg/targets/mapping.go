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

package targets

import (
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"gopkg.in/yaml.v3"
)

type templates struct {
	Version   float64 `yaml:"version"`
	Templates []struct {
		Name     string `yaml:"name"`
		Filepath string `yaml:"filepath"`
		Template string `yaml:"template"`
	} `yaml:"templates"`
}

//go:embed mapping.yaml
var mapping string

// Mapping
//
// params:
//
//	config: *configuration.Configuration
//
// returns:
//
//		type: map[string]*template.Template
//	 description:
//		  the map of templates with the key as the same as the file name to be rendered
func Mapping(config *configuration.Configuration, templatesFS embed.FS) (map[string]*template.Template, error) {
	templates := &templates{}
	err := yaml.Unmarshal([]byte(mapping), templates)
	if err != nil {
		return nil, err
	}
	templateMap := make(map[string]*template.Template, len(templates.Templates))

	for _, t := range templates.Templates {
		templateContent, err := templatesFS.ReadFile(t.Template)
		if err != nil {
			return nil, err
		}
		FilePath := t.Filepath
		if strings.Contains(FilePath, "database/migrations") {
			FilePath = strings.Replace(FilePath, "database/migrations", config.Database.Migration.Destination, 1)
		} else if strings.Contains(FilePath, "database/queries") {
			FilePath = strings.Replace(FilePath, "database/queries", config.Database.QueriesLocation, 1)
		}
		templateMap[FilePath] = template.Must(template.New(t.Name).Parse(string(templateContent)))
	}
	fmt.Printf("Loaded %d templates\n", len(templateMap))
	return templateMap, nil
}
