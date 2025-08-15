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
package cmd

import (
	"fmt"
	"os"

	"github.com/adamkali/egg_cli/pkg"
	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/styles"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd represents the init command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a new project without the wizard",
	Long:  `Initialize a new project without the wizard by using a
	configuration file in the yaml context.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := models.NewLogger("egg-log")
		if err != nil {
			panic(err)
		}

		config := new(configuration.Configuration)
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}
		configFileContent, err := os.ReadFile(configFile)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(configFileContent, config)
		if err != nil {
			panic(err)
		}
		
		configPretty, err := yaml.Marshal(config)
		if err != nil {
			panic(err)
		}
		fmt.Println(styles.EggProgressInfo.Render(string(configPretty)))
		fmt.Println(styles.EggProgressTitle.Render("🥚 Creating Project: " + config.Name))
		err = pkg.ProjectFactory(config, logger, embeddedTemplates, embeddedMapping)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("config", "c", "egg.yaml", "Path to the configuration file to generate from")
}
