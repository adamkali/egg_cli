package models

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
)



func ServerJWTInput() textinput.Model {
	serverJWTInput := textinput.New()
	serverJWTInput.Placeholder = "Default: auto-generated"
	serverJWTInput.CharLimit = 30
	serverJWTInput.Width = 30
	serverJWTInput.Prompt = ""
	return serverJWTInput
}

func ServerPortInput() textinput.Model {
	serverPortInput := textinput.New()
	serverPortInput.Placeholder = "Default: 8080"
	serverPortInput.CharLimit = 5
	serverPortInput.Width = 30
	serverPortInput.Prompt = ""
	return serverPortInput
}

func ServerFrontendDirInput() textinput.Model {
	serverFrontendDirInput := textinput.New()
	serverFrontendDirInput.Placeholder = "frontend/dist"
	serverFrontendDirInput.CharLimit = 30
	serverFrontendDirInput.Width = 30
	serverFrontendDirInput.Prompt = ""
	return serverFrontendDirInput
}

func DatabaseURLInput() textinput.Model {
	databaseURLInput := textinput.New()
	databaseURLInput.Placeholder = "Default: postgres://postgres@localhost:5432/egg?sslmode=disable"
	databaseURLInput.CharLimit = 70
	databaseURLInput.Width = 70
	databaseURLInput.Prompt = ""
	return databaseURLInput
}

func DatabaseSqlOrGoSQLCInput() textinput.Model {
	databaseURLInput := textinput.New()
	databaseURLInput.Placeholder = "sqlc"
	databaseURLInput.CharLimit = 30
	databaseURLInput.Width = 30
	databaseURLInput.Prompt = ""
	databaseURLInput.Validate = func(s string) error {
		if s == "" {
			return nil // Allow empty input
		}
		if s != "go" && s != "sql" {
			return fmt.Errorf("must be 'go' or 'sql'")
		}
		return nil
	}
	return databaseURLInput
}

func DatabaseRootLocation() textinput.Model {
	databaseMigrationLocationInput := textinput.New()
	databaseMigrationLocationInput.Placeholder = "Default: db"
	databaseMigrationLocationInput.CharLimit = 30
	databaseMigrationLocationInput.Width = 30
	databaseMigrationLocationInput.Prompt = ""
	databaseMigrationLocationInput.Validate = func(s string) error {
		if s == "" {
			return nil // Allow empty input
		}

		// Check if the path is valid for the current OS
		if !filepath.IsAbs(s) {
			// For relative paths, check if they're valid
			if filepath.Clean(s) == "." || filepath.Clean(s) == ".." {
				return fmt.Errorf("invalid path: %s", s)
			}
		}

		// Check for invalid characters in the path
		if filepath.Clean(s) != s {
			return fmt.Errorf("invalid path characters")
		}

		return nil
	}
	return databaseMigrationLocationInput
}

func ProjectHostInput() textinput.Model {
	projectHostInput := textinput.New()
	projectHostInput.Placeholder = "github.com"
	projectHostInput.CharLimit = 20
	projectHostInput.Width = 10
	projectHostInput.Prompt = ""
	return projectHostInput
}
func ProjectUsernameInput() textinput.Model {
	projectUsernameInput := textinput.New()
	projectUsernameInput.Placeholder = "adamkali"
	projectUsernameInput.CharLimit = 20
	projectUsernameInput.Width = 10
	projectUsernameInput.Prompt = ""
	return projectUsernameInput
}
func ProjectNameInput() textinput.Model {
	projectUsernameInput := textinput.New()
	projectUsernameInput.Placeholder = "egg-app"
	projectUsernameInput.CharLimit = 20
	projectUsernameInput.Width = 10
	projectUsernameInput.Prompt = ""
	return projectUsernameInput
}

func LicenseInput() textinput.Model {
	licenseInput := textinput.New()
	licenseInput.Placeholder = "MIT"
	licenseInput.CharLimit = 20
	licenseInput.Width = 20
	licenseInput.Prompt = ""
	licenseInput.Validate = func(s string) error {
		if s == "" {
			return nil // Allow empty input
		}
		validLicenses := []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "ISC", "Unlicense"}
		for _, license := range validLicenses {
			if s == license {
				return nil
			}
		}
		return fmt.Errorf("must be one of: MIT, Apache-2.0, GPL-3.0, BSD-3-Clause, ISC, Unlicense")
	}
	return licenseInput
}

func CopyrightYearInput() textinput.Model {
	copyrightYearInput := textinput.New()
	copyrightYearInput.Placeholder = "2024"
	copyrightYearInput.CharLimit = 4
	copyrightYearInput.Width = 10
	copyrightYearInput.Prompt = ""
	copyrightYearInput.Validate = func(s string) error {
		if s == "" {
			return nil // Allow empty input
		}
		// Basic year validation (1900-2100)
		if len(s) != 4 {
			return fmt.Errorf("year must be 4 digits")
		}
		year := 0
		if _, err := fmt.Sscanf(s, "%d", &year); err != nil {
			return fmt.Errorf("invalid year format")
		}
		if year < 1900 || year > 2100 {
			return fmt.Errorf("year must be between 1900 and 2100")
		}
		return nil
	}
	return copyrightYearInput
}

func CopyrightAuthorInput() textinput.Model {
	copyrightAuthorInput := textinput.New()
	copyrightAuthorInput.Placeholder = "Your Name"
	copyrightAuthorInput.CharLimit = 50
	copyrightAuthorInput.Width = 30
	copyrightAuthorInput.Prompt = ""
	return copyrightAuthorInput
}

func ServerFrontendApiInput() textinput.Model {
	serverFrontendApiInput := textinput.New()
	serverFrontendApiInput.Placeholder = "frontend/src/api"
	serverFrontendApiInput.CharLimit = 30
	serverFrontendApiInput.Width = 30
	serverFrontendApiInput.Prompt = ""
	return serverFrontendApiInput
}


func MinioUrlInput() textinput.Model {
	minioUrlInput := textinput.New()
	minioUrlInput.Placeholder = "localhost:9000"
	minioUrlInput.CharLimit = 100
	minioUrlInput.Width = 100
	minioUrlInput.Prompt = ""
	minioUrlInput.Blur()
	return minioUrlInput
}

func MinioAccessKeyInput() textinput.Model {
	minioAccessKeyInput := textinput.New()
	minioAccessKeyInput.Placeholder = "Default: auto-generated" 
	minioAccessKeyInput.CharLimit = 100
	minioAccessKeyInput.Width = 100
	minioAccessKeyInput.Prompt = ""
	minioAccessKeyInput.Blur()
	return minioAccessKeyInput
}

func MinioSecretKeyInput() textinput.Model {
	minioSecretKeyInput := textinput.New()
	minioSecretKeyInput.Placeholder = "Default: auto-generated"
	minioSecretKeyInput.CharLimit = 100
	minioSecretKeyInput.Width = 100
	minioSecretKeyInput.Prompt = ""
	minioSecretKeyInput.Blur()
	return minioSecretKeyInput
}
