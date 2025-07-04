package templates_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/adamkali/egg_cli/pkg/templates"
)

const ResultDBCmdTemplate = `
/* Generated by egg v0.0.1
Copyright © 2022 Adam Kalinowski 

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
package cmd

import (
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database interactions",
	Long: "Database interactions such as migrations",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
`

func TestDBCmdTemplate(t *testing.T) {
	// load the template
	temp := templates.DBCmdTemplate
	templateTest := template.Must(template.New("db.go").Parse(temp))

	// execute the template
	stringWriter := new(bytes.Buffer)
	err := templateTest.ExecuteTemplate(stringWriter, "db.go", createConfiguration())
	if err != nil {
		t.Error(err)
	}

}
