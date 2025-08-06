package targets

import (
	"embed"
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

// Mapping
//
// params:
// 	config: *configuration.Configuration
// returns:
// 	type: map[string]*template.Template
//  description:
// 	  the map of templates with the key as the same as the file name to be rendered
func Mapping(config *configuration.Configuration, templatesFS embed.FS, mappingYaml string) (map[string]*template.Template, error) {
	templates := &templates{}
	err := yaml.Unmarshal([]byte(mappingYaml), templates)
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
	return templateMap, nil
}
