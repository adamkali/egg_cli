package targets

import (
	"embed"
	"text/template"


	"gopkg.in/yaml.v3"
	"github.com/adamkali/egg_cli/pkg/configuration"
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
		templateMap[t.Filepath] = template.Must(template.New(t.Name).Parse(string(templateContent)))
	}
	return templateMap, nil
}
