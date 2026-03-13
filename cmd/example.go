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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/styles"
	"github.com/spf13/cobra"
)

func generate_random_base64(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// exampleCmd represents the example command
var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Generate a sample egg.yaml file for project configuration.",
	Long: `
Generate a sample egg.yaml file for project configuration.

This will generate a sample egg.yaml file for project configuration.
with a sample configuration for a simple project.

For jwt, Database password, redis password, and s3 keys: 
the program will generate random base64 strings.
for all other fields, the program will use generic values like
default, yourname, etc.

The file will be saved to the current working directory.
but you can still use the --config flag to specify a different path.

be advised, when using the egg_cli generate example --config <file> command
, the following egg_cli generate command should be called with the same --config flag

example:
	egg_cli generate example --config <file>
	egg_cli generate --config <file>
`,
	Run: func(cmd *cobra.Command, args []string) {
		config := new(configuration.Configuration)
		config.Namespace = "github.com/yourname/example"
		config.Name = "example"
		config.License = "MIT"
		config.Semver = "0.0.1"
		config.Copyright.Author = "Your Name"
		config.Copyright.Year = time.Now().Year()

		var port int
		port = 8080
		config.Auth.Provider = "jwt"
		config.Auth.JWT = base64.StdEncoding.EncodeToString([]byte(generate_random_base64(32)))
		config.Server.Port = port
		config.Server.Frontend.Dir = "web/dist"
		config.Server.Frontend.Api = "web/src/api"

		config.Database.URL = fmt.Sprintf("postgres://postgres:%s@localhost:5432/example?sslmode=disable", generate_random_base64(64))
		config.Database.QueriesLocation = "db/queries"
		config.Database.Migration.Destination = "queries"
		config.Database.Migration.Protocol = "postgres"
		config.Database.Sqlc.RepositoryLocation = "db/repository"
		config.Database.Sqlc.Schema = "postgresql"
		config.Database.Sqlc.SqlOrGo = "sql"

		config.Cache.URL = fmt.Sprintf("redis://%s:%s@localhost:6379/0",
			generate_random_base64(8),
			generate_random_base64(32),
		)

		config.S3.URL = fmt.Sprintf("localhost:9000")
		config.S3.Access = generate_random_base64(32)
		config.S3.Secret = generate_random_base64(32)

		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}
		if configFile == "" {
			configFile = "egg.yaml"
		}

		if err := config.SaveConfiguration(configFile); err != nil {
			fmt.Println(
				styles.EggProgressError.Render("Failed to generate sample egg.yaml file: " + err.Error()),
			)
			os.Exit(4)
		}

		fmt.Println(styles.EggProgressInfo.Render("Sample egg.yaml file generated and saved to " + configFile))
	},
}

func init() {
	generateCmd.AddCommand(exampleCmd)
	exampleCmd.Flags().StringP("config", "c", "egg.yaml", "Path to the configuration file to generate from")
}
