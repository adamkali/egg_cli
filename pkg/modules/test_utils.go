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

package modules

import (
	"testing"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
)

// createTestLogger creates a logger for testing
func createTestLogger(t *testing.T) *models.EggLog {
	logger, err := models.NewLogger("test.log")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return logger
}

// createTestConfiguration creates a mock configuration for testing
func createTestConfiguration() *configuration.Configuration {
	return &configuration.Configuration{
		Namespace: "github.com/testuser/testproject",
		Name:      "testproject",
		Semver:    "0.0.1",
		License:   "MIT",
		Copyright: struct {
			Year   int    "yaml:\"year\""
			Author string "yaml:\"author\""
		}{
			Year:   2024,
			Author: "Test User",
		},
		Auth: struct {
			Provider          string "yaml:\"provider\""
			JWT               string "yaml:\"jwt\""
			Auth0Domain       string "yaml:\"auth0_domain\""
			Auth0Audience     string "yaml:\"auth0_audience\""
			Auth0ClientID     string "yaml:\"auth0_client_id\""
			Auth0ClientSecret string "yaml:\"auth0_client_secret\""
		}{
			Provider: "jwt",
			JWT:      "test-secret",
		},
		Server: struct {
			Port     int    "yaml:\"port\""
			Frontend struct {
				Dir string "yaml:\"dir\""
				Api string "yaml:\"api\""
			} "yaml:\"frontend\""
		}{
			Port: 8080,
			Frontend: struct {
				Dir string "yaml:\"dir\""
				Api string "yaml:\"api\""
			}{Dir: "web/dist", Api: "web/src/api"},
		},
		Database: struct {
			URL                    string "yaml:\"url\""
			Sqlc struct {
				RepositoryLocation string `yaml:"repository"`
				Schema             string `yaml:"schema"`
				SqlOrGo            string `yaml:"sql_or_go"`
			} `yaml:"sqlc"`
			QueriesLocation        string "yaml:\"queries\""
			Migration              struct {
				Protocol    string "yaml:\"protocol\""
				Destination string "yaml:\"destination\""
			} "yaml:\"migration\""
		}{
			URL:                    "postgres://postgres:postgres@localhost:5432/test?sslmode=disable",
			Sqlc: struct {
				RepositoryLocation string "yaml:\"repository\""
				Schema             string "yaml:\"schema\""
				SqlOrGo            string "yaml:\"sql_or_go\""
			}{
				RepositoryLocation: "db/repository",
				Schema:             "postgres",
				SqlOrGo:            "sql",
			},
			QueriesLocation:        "db/queries",
			Migration: struct {
				Protocol    string "yaml:\"protocol\""
				Destination string "yaml:\"destination\""
			}{
				Protocol:    "postgres",
				Destination: "db/migrations",
			},
		},
		Cache: struct {
			URL string "yaml:\"url\""
		}{
			URL: "redis://localhost:6379/0",
		},
		S3: struct {
			URL    string "yaml:\"url\""
			Access string "yaml:\"access\""
			Secret string "yaml:\"secret\""
		}{
			URL:    "https://localhost:9000",
			Access: "test-access",
			Secret: "test-secret",
		},
	}
}
