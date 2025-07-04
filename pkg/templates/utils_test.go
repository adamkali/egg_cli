package templates_test

import (
	"strings"

	"github.com/adamkali/egg_cli/pkg/configuration"
)

func createConfiguration() *configuration.Configuration {
	return &configuration.Configuration{
		Namespace: "github.com/adamkali/egg",
		Name:      "egg",
		Semver:    "0.0.1",
		License:   "Apache-2.0",
		Copyright: struct {
			Year   int    "yaml:\"year\""
			Author string "yaml:\"author\""
		}{
			Year:   2022,
			Author: "Adam Kalinowski",
		},
		Server: struct {
			JWT      string "yaml:\"jwt\""
			Port     int    "yaml:\"port\""
			Frontend struct {
				Dir string "yaml:\"dir\""
				Api string "yaml:\"api\""
			} "yaml:\"frontend\""
		}{
			JWT:  "secret",
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
			URL:                    "postgres://postgres:postgres@localhost:5432/egg?sslmode=disable",
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
			Access: "AKIAIOSFODNN7EXAMPLE",
			Secret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		},
	}
}

type DiffMap struct {
	Expected string
	Actual   string
}

func Diff(expected, actual string) map[int]DiffMap {
	// Take the two strings and lexically compare them
	expectedLines := []string{}
	actualLines := []string{}
	strings.Split(expected, "\n")
	strings.Split(actual, "\n")
	changes := make(map[int]DiffMap, len(expectedLines))
	for i, _ := range expectedLines {
		if expectedLines[i] != actualLines[i] {
			changes[i] = DiffMap{
				Expected: expectedLines[i],
				Actual:   actualLines[i],
			}
		}

	}
	return changes
}
