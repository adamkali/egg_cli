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
		Server: struct {
			JWT      string "yaml:\"jwt\""
			Port     int    "yaml:\"port\""
			Frontend struct {
				Dir string "yaml:\"dir\""
				Api string "yaml:\"api\""
			} "yaml:\"frontend\""
		}{
			JWT:  "test-secret",
			Port: 8080,
			Frontend: struct {
				Dir string "yaml:\"dir\""
				Api string "yaml:\"api\""
			}{Dir: "web/dist", Api: "web/src/api"},
		},
		Database: struct {
			URL                    string "yaml:\"url\""
			Sqlc                   string "yaml:\"sqlc\""
			SqlcRepositoryLocation string "yaml:\"repository\""
			QueriesLocation        string "yaml:\"queries\""
			Migration              struct {
				Protocol    string "yaml:\"protocol\""
				Destination string "yaml:\"destination\""
			} "yaml:\"migration\""
		}{
			URL:                    "postgres://postgres:postgres@localhost:5432/test?sslmode=disable",
			Sqlc:                   "sql",
			SqlcRepositoryLocation: "db/repository",
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
