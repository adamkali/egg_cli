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
	"strconv"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/styles"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const (
	initMessage = "Creating config from environment variables"
)


var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Create Config from environment variables",
	Long: `Create Config from environment variables.
	
	Use this command to create a config from environment variables, 
	which can be useful for a variety of purposes. To make sure that 
	important secrets don't get committed to source control, you can
	use this command to scaffold a config file from environment variables.

	The structure of the environment variables is similar to the
	configuration file. with a few exceptions.

	1. "EGG_" is required at the beginning of the variable names.
	2. The structure of the configuration file is the structure of the
	   string name of the environment variable. i.e. Server.Frontend.Api 
	   is the same as EGG_SERVER_FRONTEND_API.
	3. Variables that like EGG_DATABASE_URL are not supported. This is so
	   that you can also use EGG_DATABASE_PASSWORD and EGG_DATABASE_USERNAME
	   in a docker-compose file, that you could use if you do not wish to 
	   use managed systems for the cache, database, etc.
	4. As A rule of thumb,  

	# The following environment variables are supported, 
	# and will be used to create the config file, 
	# any with a X after them are REQUIRED during runtime of an
	# egg project.
	#-----------------------------------------------------#
	EGG_NAMESPACE="github.com/yourname/example"           # 
	EGG_NAME="example"                                    # X 
	EGG_LICENSE="MIT"                                     # 
	EGG_SEMVER="0.0.1"                                    # X
	#-----------------------------------------------------#
	EGG_COPYRIGHT_AUTHOR="Your Name"                      #
	EGG_COPYRIGHT_YEAR="2023"                             #
	#-----------------------------------------------------#
	EGG_SERVER_PORT=8080                                  # X
	EGG_SERVER_JWT="your_base64_encoded_jwt"              # X
	EGG_SERVER_FRONTEND_DIR="web/dist"                    # X
	EGG_SERVER_FRONTEND_API="web/src/api"                 #
	#-----------------------------------------------------#
	EGG_DATABASE_USERNAME="postgres"                      # X
	EGG_DATABASE_PASSWORD="your_base64_encoded_password"  # X
	EGG_DATABASE_DBNAME="example"                         # X
	EGG_DATABASE_PORT=5432                                # X
	EGG_DATABASE_HOST="localhost"                         # X
	EGG_DATABASE_QUERIES_LOCATION="db/queries"            #
	EGG_DATABASE_MIGRATION_DESTINATION="queries"          #
	EGG_DATABASE_MIGRATION_PROTOCOL="postgres"            #
	EGG_DATABASE_SQLC_REPOSITORY_LOCATION="db/repository" #
	EGG_DATABASE_SQLC_SCHEMA="postgresql"                 #
	EGG_DATABASE_SQLC_SQL_OR_GO="sql"                     #
	#-----------------------------------------------------#
	EGG_CACHE_USERNAME=""                                 # X
	EGG_CACHE_PASSWORD="something"                        # X
	EGG_CACHE_HOST="localhost"	                          # X
	EGG_CACHE_PORT=6379                                   # X
	EGG_CACHE_DB=0                                        # X
	#-----------------------------------------------------#
	EGG_S3_HOST="localhost"                               # X
	EGG_S3_PORT=9000                                      # X
	EGG_S3_ACCESS_KEY_ID=""                               # X
	EGG_S3_SECRET_ACCESS_KEY=""                           # X
	#-----------------------------------------------------#
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config := new(configuration.Configuration)
		envFile := cmd.Flag("input").Value.String()
		if envFile == "" {
			envFile = ".env"
		}
		godotenv.Load(envFile)

		fmt.Println(styles.EggProgressInfo.Render(initMessage))

		// need to check for required variables
		// anything with a X after it is required
		config.Namespace = os.Getenv("EGG_NAMESPACE")
		config.Name = os.Getenv("EGG_NAME")
		config.License = os.Getenv("EGG_LICENSE")
		config.Semver = os.Getenv("EGG_SEMVER")
		config.Copyright.Author = os.Getenv("EGG_COPYRIGHT_AUTHOR")
		config.Copyright.Year, _ = strconv.Atoi(
			os.Getenv("EGG_COPYRIGHT_YEAR"),
		)
		if config.Name == "" {
			fmt.Println(styles.EggProgressError.Render("EGG_NAME is required"))
			os.Exit(1)
		}
		if config.Semver == "" {
			fmt.Println(styles.EggProgressError.Render("EGG_SEMVER is required"))
			os.Exit(1)
		}

		var port int
		port, err := strconv.Atoi(os.Getenv("EGG_SERVER_PORT"))
		if err != nil {
			fmt.Println(styles.EggProgressError.Render("EGG_SERVER_PORT is required"))
			os.Exit(1)
		}
		if port <= 0 {
			fmt.Println(styles.EggProgressError.Render("EGG_SERVER_PORT Must be between 0 and 65535"))
			os.Exit(6)
		}
		if port > 65535 {
			fmt.Println(styles.EggProgressError.Render("EGG_SERVER_PORT Must be between 0 and 65535"))
			os.Exit(6)
		}
		config.Server.Port = port
		config.Server.JWT = os.Getenv("EGG_SERVER_JWT")
		config.Server.Frontend.Dir = os.Getenv("EGG_SERVER_FRONTEND_DIR")
		config.Server.Frontend.Api = os.Getenv("EGG_SERVER_FRONTEND_API")
		if config.Server.JWT == "" {
			fmt.Println(styles.EggProgressError.Render("EGG_SERVER_JWT is required"))
			os.Exit(1)
		}
		if config.Server.Frontend.Dir == "" {
			fmt.Println(styles.EggProgressError.Render("EGG_SERVER_FRONTEND_DIR is required"))
			os.Exit(1)
		}

		dbport, err := strconv.Atoi(os.Getenv("EGG_DATABASE_PORT"))
		if err != nil {
			fmt.Println(styles.EggProgressError.Render("EGG_DATABASE_PORT is required"))
			os.Exit(1)
		}
		username := os.Getenv("EGG_DATABASE_USERNAME")
		password := os.Getenv("EGG_DATABASE_PASSWORD")
		dbname := os.Getenv("EGG_DATABASE_DBNAME")
		host := os.Getenv("EGG_DATABASE_HOST")

		if username == "" {
			fmt.Println(styles.EggProgressError.Render("EGG_DATABASE_USERNAME is required"))
			os.Exit(1)
		}
		if password == "" {
			fmt.Println(styles.EggProgressError.Render("EGG_DATABASE_PASSWORD is required"))
			os.Exit(1)
		}
		if dbname == "" {
			fmt.Println(styles.EggProgressError.Render("EGG_DATABASE_DBNAME is required"))
			os.Exit(1)
		}
		if dbport <= 0 {
			fmt.Println(styles.EggProgressError.Render("EGG_DATABASE_PORT Must be between 0 and 65535"))
			os.Exit(6)
		}
		if dbport > 65535 {
			fmt.Println(styles.EggProgressError.Render("EGG_DATABASE_PORT Must be between 0 and 65535"))
			os.Exit(6)
		}

		config.Database.URL = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			username,
			password,
			host,
			dbport,
			dbname,
		)
		config.Database.QueriesLocation = os.Getenv("EGG_DATABASE_QUERIES_LOCATION")
		config.Database.Migration.Destination = os.Getenv("EGG_DATABASE_MIGRATION_DESTINATION")
		config.Database.Migration.Protocol = os.Getenv("EGG_DATABASE_MIGRATION_PROTOCOL")
		config.Database.Sqlc.RepositoryLocation = os.Getenv("EGG_DATABASE_SQLC_REPOSITORY_LOCATION")
		config.Database.Sqlc.Schema = os.Getenv("EGG_DATABASE_SQLC_SCHEMA")
		config.Database.Sqlc.SqlOrGo = os.Getenv("EGG_DATABASE_SQLC_SQL_OR_GO")

		port, _ = strconv.Atoi(os.Getenv("EGG_CACHE_PORT"))
		if port <= 0 {
			fmt.Println(styles.EggProgressError.Render("EGG_CACHE_PORT Must be between 0 and 65535"))
			os.Exit(6)
		}
		if port > 65535 {
			fmt.Println(styles.EggProgressError.Render("EGG_CACHE_PORT Must be between 0 and 65535"))
			os.Exit(6)
		}
		username = os.Getenv("EGG_CACHE_USERNAME")
		password = os.Getenv("EGG_CACHE_PASSWORD")
		host = os.Getenv("EGG_CACHE_HOST")
		dbname = os.Getenv("EGG_CACHE_DB")
		config.Cache.URL = fmt.Sprintf(
			"redis://%s:%s@%s:%d/%s",
			username,
			password,
			host,
			port,
			dbname,
		)
		port, err = strconv.Atoi(os.Getenv("EGG_S3_PORT"))
		config.S3.Access = os.Getenv("EGG_S3_ACCESS_KEY_ID")
		config.S3.Secret = os.Getenv("EGG_S3_SECRET_ACCESS_KEY")
		host = os.Getenv("EGG_S3_HOST")
		if err != nil {
			fmt.Println(styles.EggProgressError.Render("EGG_S3_PORT is required"))
			os.Exit(1)
		}
		if port <= 0 {
			fmt.Println(styles.EggProgressError.Render("EGG_S3_PORT Must be between 0 and 65535"))
			os.Exit(6)
		}
		if port > 65535 {
			fmt.Println(styles.EggProgressError.Render("EGG_S3_PORT Must be between 0 and 65535"))
			os.Exit(6)
		}
		if host == "" {
			fmt.Println(styles.EggProgressError.Render("EGG_S3_HOST is required"))
			os.Exit(1)
		}
		if config.S3.Access == "" {
			fmt.Println(styles.EggProgressError.Render("EGG_S3_ACCESS_KEY_ID is required"))
			os.Exit(1)
		}
		if config.S3.Secret == "" {
			fmt.Println(styles.EggProgressError.Render("EGG_S3_SECRET_ACCESS_KEY is required"))
			os.Exit(1)
		}
		config.S3.URL = fmt.Sprintf(
			"%s:%d",
			host,
			port,
		)

		environment := cmd.Flag("env").Value.String() 
		config.GenerateConfigurationFile(environment)
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.Flags().StringP("env", "e", "dev", "Environment to generate")
	envCmd.Flags().StringP("input", "i", ".env", "Path to the environment file to generate from")
}
