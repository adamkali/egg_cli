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
	"strconv"
	"time"

	"github.com/adamkali/egg_cli/configuration"
	"github.com/adamkali/egg_cli/models"
	"github.com/adamkali/egg_cli/pkg"
	"github.com/adamkali/egg_cli/state"
	"github.com/adamkali/egg_cli/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	defaultPort         = 8080
	defaultDatabaseRoot = "db"
	defaultDatabaseURL  = "postgres://postgres@localhost:5432/egg?sslmode=disable"
)

func GenerateJWTSecret(nBytes int) (string, error) {
	// Generate nBytes of random data as the secret key
	bytes := make([]byte, nBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("Failed to generate JWT secret: %v", err)
	}

	// Encode the byte array to base64 string
	encoded := base64.StdEncoding.EncodeToString(bytes)

	return encoded, nil
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Long:  `Initialize a new project with interactive setup`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := models.NewLogger("egg-log")
		if err != nil {
			panic(err)
		}

		pageModel := models.CreatePageModel(logger)
		p := tea.NewProgram(pageModel)
		if _, err := p.Run(); err != nil {
			logger.Error("Error running program: %v", err)
		}
		config := new(configuration.Configuration)
		config.Namespace = state.ProjectNamespace
		config.Name = state.ProjectName
		config.License = state.License
		config.Semver = "0.0.1"
		config.Copyright.Author = state.ProjectUsername
		config.Copyright.Year = time.Now().Year()

		if state.CopyrightYear == "" {
			// get current year
			generatingCopyrightYearMessage := fmt.Sprintf("Generating copyright year")
			generatingCopyrightYearMessage = styles.EggProgressInfo.Render(generatingCopyrightYearMessage)
			fmt.Println(generatingCopyrightYearMessage)
			config.Copyright.Year = time.Now().Year()
		}

		var port int
		if state.ServerPort == "" {
			port = defaultPort
		} else {
			port, _ = strconv.Atoi(state.ServerPort)
		}
		if state.ServerJWT == "" {
			generatingJWTSecretMessage := fmt.Sprintf("Generating JWT secret")
			generatingJWTSecretMessage = styles.EggProgressInfo.Render(generatingJWTSecretMessage)
			fmt.Println(generatingJWTSecretMessage)
			secret, err := GenerateJWTSecret(32)
			if err != nil {
				panic(err)
			}
			state.ServerJWT = secret
		}

		config.Server = struct {
			JWT      string "yaml:\"jwt\""
			Port     int    "yaml:\"port\""
			Frontend struct {
				Dir string "yaml:\"dir\""
				Api string "yaml:\"api\""
			} "yaml:\"frontend\""
		}{
			Port: port,
			JWT:  state.ServerJWT,
			Frontend: struct {
				Dir string "yaml:\"dir\""
				Api string "yaml:\"api\""
			}{
				Dir: "web/dist",
				Api: "web/src/api",
			},
		}

		if state.DatabaseURL == "" {
			defaultingDatabaseURLMessage := fmt.Sprintf("Defaulting to %s as database url", defaultDatabaseURL)
			logger.Info(defaultingDatabaseURLMessage)
			defaultingDatabaseURLMessage = styles.EggProgressInfo.Render(defaultingDatabaseURLMessage)
			fmt.Println(defaultingDatabaseURLMessage)
			state.DatabaseURL = defaultDatabaseURL 
		}

		if state.DatabaseSqlcOrGo == "" {
			defaultingDatabaseSqlcOrGoMessage := fmt.Sprintf("Defaulting to %s as database sqlc or go", "sql")
			logger.Info(defaultingDatabaseSqlcOrGoMessage)
			defaultingDatabaseSqlcOrGoMessage = styles.EggProgressInfo.Render(defaultingDatabaseSqlcOrGoMessage)
			fmt.Println(defaultingDatabaseSqlcOrGoMessage)
			state.DatabaseSqlcOrGo = "sql"
		}

		if state.DatabaseRoot == "" {
			defaultingDatabaseRootMessage := fmt.Sprintf("Defaulting to %s as database root directory", defaultDatabaseRoot)
			logger.Info(defaultingDatabaseRootMessage)
			defaultingDatabaseRootMessage = styles.EggProgressInfo.Render(defaultingDatabaseRootMessage)
			fmt.Println(defaultingDatabaseRootMessage)
			state.DatabaseRoot = "db"
		}
		config.Database = struct {
			URL                    string "yaml:\"url\""
			Sqlc                   string "yaml:\"sqlc\""
			SqlcRepositoryLocation string "yaml:\"repository\""
			QueriesLocation        string "yaml:\"queries\""
			Migration              struct {
				Protocol    string "yaml:\"protocol\""
				Destination string "yaml:\"destination\""
			} "yaml:\"migration\""
		}{
			URL:                    state.DatabaseURL,
			Sqlc:                   state.DatabaseSqlcOrGo,
			SqlcRepositoryLocation: state.DatabaseRoot + "/repository",
			QueriesLocation:        state.DatabaseRoot + "/queries",
			Migration: struct {
				Protocol    string "yaml:\"protocol\""
				Destination string "yaml:\"destination\""
			}{
				// for now we only support postgres
				Protocol:    "postgresql",
				Destination: state.DatabaseRoot + "/migrations",
			},
		}

		config.Cache = struct {
			URL string "yaml:\"url\""
		}{
			URL: "redis://localhost:6379",
		}
		if state.MinioAccessKey == "" {
			secret, err := GenerateJWTSecret(32)
			if err != nil {
				panic(err)
			}
			defaultingMinioAccessKeyMessage := fmt.Sprintf("Defaulting to %s as minio access key", secret)
			logger.Info(defaultingMinioAccessKeyMessage)
			defaultingMinioAccessKeyMessage = styles.EggProgressInfo.Render(defaultingMinioAccessKeyMessage)
			fmt.Println(defaultingMinioAccessKeyMessage)
			fmt.Println("You should change this in the future to what AWS S3 / Minio generates when registering a new user")
			state.MinioAccessKey = secret
		}
		if state.MinioSecretKey == "" {
			secret, err := GenerateJWTSecret(32)
			if err != nil {
				panic(err)
			}
			defaultingMinioSecretKeyMessage := fmt.Sprintf("Defaulting to %s as minio secret key", secret)
			logger.Info(defaultingMinioSecretKeyMessage)
			defaultingMinioSecretKeyMessage = styles.EggProgressInfo.Render(defaultingMinioSecretKeyMessage)
			fmt.Println(defaultingMinioSecretKeyMessage)
			fmt.Println("You should change this in the future to what AWS S3 / Minio generates when registering a new user")
			state.MinioSecretKey = secret
		}

		config.S3 = struct {
			URL    string "yaml:\"url\""
			Access string "yaml:\"access\""
			Secret string "yaml:\"secret\""
		}{
			URL:    state.MinioURL,
			Access: state.MinioAccessKey,
			Secret: state.MinioSecretKey,
		}

		fmt.Printf("\n")
		configPretty, err := yaml.Marshal(config)
		if err != nil {
			panic(err)
		}
		fmt.Println(styles.EggProgressInfo.Render(string(configPretty)))
		fmt.Println(styles.EggProgressTitle.Render("🥚 Creating Project: " + config.Name))
		err = pkg.ProjectFactory(config, logger)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
