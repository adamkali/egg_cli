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

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/pkg/modules"
	"github.com/adamkali/egg_cli/styles"
	"github.com/spf13/cobra"
)

func loadConfig(configFile string) (*configuration.Configuration, error) {
	return configuration.LoadConfigurationFile(configFile)
}

// dotenvCmd represents the dotenv command
var dotenvCmd = &cobra.Command{
	Use:   "dotenv",
	Args:  cobra.MaximumNArgs(1),
	Short: "Generate a .env file",
	Long: `
Generate a .env file from a configuration file.

This will generate a .env file from a configuration file.

Example:
	egg_cli generate dotenv <path to config>/<some>.yaml
	or
	egg_cli generate dotenv

If no configuration file is passed in, egg will look for a config/developent.yaml.
If a configuration file is passed in, it will be used as the path to the configuration file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// get the path to the config file
		eggl, err := models.NewLogger("egg-log")
		configFile := "./config/developent.yaml"
		config := new(configuration.Configuration)
		if len(args) == 1 {
			configFile = args[0]
		}
		if len(args) == 0 {
			// check if there is anything in stdin and use that as the config file instead
			stdin, err := os.Stdin.Stat()
			if err != nil {
				// if not then use the default config fileA
				configFile = "./config/developent.yaml"
				config, err = loadConfig(configFile)
				if err != nil {
					fmt.Println(
						styles.EggProgressError.Render("Failed to load config file: " + err.Error()),
					)
					os.Exit(9)
				}
			} else {
				if stdin.Mode()&os.ModeNamedPipe == 0 {
					// if not then use the default config file
					configFile = "./config/developent.yaml"
					config, err = loadConfig(configFile)
					if err != nil {
						fmt.Println(
							styles.EggProgressError.Render("Failed to load config file: " + err.Error()),
						)
						os.Exit(9)
					}
				} else {
					// if there is something in stdin then use that as the config file
					// and we can load as bytes
					config, err = configuration.LoadConfigurationReader(os.Stdin)
					if err != nil {
						fmt.Println(
							styles.EggProgressError.Render("Failed to load config file: " + err.Error()),
						)
						os.Exit(9)
					}
				}
			}
		}

		// generate the .env file
		if err != nil {
			fmt.Println(
				styles.EggProgressError.Render("Failed to generate .env file: " + err.Error()),
			)
			os.Exit(10)
		}

		module := new(modules.GenerateDotenvModule)
		module.SetEggVersion(appVersion)
		module.SetOutput(cmd.Flag("output").Value.String())
		module.LoadFromConfig(config, eggl)
		module.Run()
		err = module.IsError()
		if err != nil {
			fmt.Println(
				styles.EggProgressError.Render("Writing .env file failed: " + err.Error()),
			)
			os.Exit(11)
		}
		fmt.Println(styles.EggProgressComplete.Render("Generated .env file"))
		fmt.Println(styles.EggProgressComplete.Render(fmt.Sprintf("egg_cli %s - .env generation complete", appVersion)))
		os.Exit(0)
	},
}

func init() {
	generateCmd.AddCommand(dotenvCmd)
	dotenvCmd.Flags().StringP("output", "o", "egg.yaml", "Path to the output .env file")
}
