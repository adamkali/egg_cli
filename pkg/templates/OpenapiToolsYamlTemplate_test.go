package templates_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/adamkali/egg_cli/pkg/templates"
)

const ResultOpenapiToolsJSONTemplate = `
{
  "$schema": "./node_modules/@openapitools/openapi-generator-cli/config.schema.json",
  "spaces": 2,
  "generator-cli": {
    "version": "7.12.0"
  }
}`

func TestOpenapiToolsJSONTemplate(t *testing.T) {
	// load the template
	temp := templates.OpenapitoolsJSONTemplate
	templateTest := template.Must(template.New("openapitools.json").Parse(temp))

	// execute the template
	stringWriter := new(bytes.Buffer)
	err := templateTest.ExecuteTemplate(stringWriter, "openapitools.json", createConfiguration())
	if err != nil {
		t.Error(err)
	}

	// check the result
	if stringWriter.String() != ResultOpenapiToolsJSONTemplate {
		diff := Diff(stringWriter.String(), ResultOpenapiToolsJSONTemplate)
		for i, v := range diff {
			t.Errorf("line %d: expected %s, got %s", i, v.Expected, v.Actual)
		}
	}
}
