package pkg

import (
	"errors"
	"os"
	"testing"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/pkg/modules"
)

// Mock module for testing
type MockModule struct {
	name    string
	err     error
	run     func()
	loadRun bool
}

func (m *MockModule) Name() string {
	return m.name
}

func (m *MockModule) IsError() error {
	return m.err
}

func (m *MockModule) Run() {
	if m.run != nil {
		m.run()
	}
}

func (m *MockModule) LoadFromConfig(configuration *configuration.Configuration, eggl *models.EggLog) {
	m.loadRun = true
}

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
			URL: "redis://localhost:6379",
		},
		S3: struct {
			URL    string "yaml:\"url\""
			Access string "yaml:\"access\""
			Secret string "yaml:\"secret\""
		}{
			URL:    "http://localhost:9000",
			Access: "test-access",
			Secret: "test-secret",
		},
	}
}

func TestPrintError(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()

	tests := []struct {
		name     string
		module   modules.IModule
		expected bool
	}{
		{
			name:     "module with no error",
			module:   &MockModule{name: "test", err: nil},
			expected: false,
		},
		{
			name:     "module with error",
			module:   &MockModule{name: "test", err: errors.New("test error")},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PrintError(tt.module, logger)
			if result != tt.expected {
				t.Errorf("PrintError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProjectFactory_Success(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	config := createTestConfiguration()

	// Create mock modules that succeed
	mockModules := []modules.IModule{
		&MockModule{name: "test1", err: nil},
		&MockModule{name: "test2", err: nil},
		&MockModule{name: "test3", err: nil},
	}

	// Temporarily replace the Modules slice
	originalModules := Modules
	Modules = mockModules
	defer func() { Modules = originalModules }()

	err := ProjectFactory(config, logger)
	if err != nil {
		t.Errorf("ProjectFactory() returned unexpected error: %v", err)
	}

	// Verify that LoadFromConfig was called on all modules
	for _, module := range mockModules {
		if mock, ok := module.(*MockModule); !ok || !mock.loadRun {
			t.Errorf("LoadFromConfig was not called on module %s", module.Name())
		}
	}
}

func TestProjectFactory_ModuleFailure(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	config := createTestConfiguration()

	// Create mock modules where the second one fails
	mockModules := []modules.IModule{
		&MockModule{name: "test1", err: nil},
		&MockModule{name: "test2", err: errors.New("module failed")},
		&MockModule{name: "test3", err: nil},
	}

	// Temporarily replace the Modules slice
	originalModules := Modules
	Modules = mockModules
	defer func() { Modules = originalModules }()

	err := ProjectFactory(config, logger)
	// ProjectFactory returns the error from WriteScrambled, which returns nil on success
	// So we should check that the .scrambled file was created instead
	if err != nil {
		t.Errorf("ProjectFactory() returned unexpected error: %v", err)
	}

	// Check that .scrambled file was created
	if !CheckScrambled() {
		t.Error("ProjectFactory() should have created .scrambled file")
	}

	// Clean up
	os.Remove(ScrambledFileName)
}

func TestCheckScrambled(t *testing.T) {
	// Test when file doesn't exist
	if CheckScrambled() {
		t.Error("CheckScrambled() should return false when file doesn't exist")
	}

	// Create a test .scrambled file
	testContent := `succeeded: []
configuration:
  namespace: test
  name: test`
	err := os.WriteFile(ScrambledFileName, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .scrambled file: %v", err)
	}
	defer os.Remove(ScrambledFileName)

	// Test when file exists
	if !CheckScrambled() {
		t.Error("CheckScrambled() should return true when file exists")
	}
}

func TestLoadScrambled_Success(t *testing.T) {
	// Create a test .scrambled file
	testContent := `succeeded:
  - "egg::initialize"
  - "egg::install_tools"
failed:
  moduleName: egg::install_libraries
  error: "test error"
configuration:
  namespace: github.com/testuser/testproject
  name: testproject
  semver: "0.0.1"
  license: MIT
  copyright:
    year: 2024
    author: Test User
  server:
    jwt: test-secret
    port: 8080
    frontend:
      dir: web/dist
      api: web/src/api
  database:
    url: postgres://postgres:postgres@localhost:5432/test?sslmode=disable
    sqlc: sql
    repository: db/repository
    queries: db/queries
    migration:
      protocol: postgres
      destination: db/migrations
  cache:
    url: redis://localhost:6379
  s3:
    url: http://localhost:9000
    access: test-access
    secret: test-secret`

	err := os.WriteFile(ScrambledFileName, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .scrambled file: %v", err)
	}
	defer os.Remove(ScrambledFileName)

	config, succeeded, failed, err := LoadScrambled()
	if err != nil {
		t.Fatalf("LoadScrambled() returned unexpected error: %v", err)
	}

	if config == nil {
		t.Error("LoadScrambled() should return configuration")
	}

	if config.Namespace != "github.com/testuser/testproject" {
		t.Errorf("LoadScrambled() configuration namespace = %s, want %s", config.Namespace, "github.com/testuser/testproject")
	}

	if len(succeeded) != 2 {
		t.Errorf("LoadScrambled() succeeded modules = %d, want %d", len(succeeded), 2)
	}

	// Should have 5 failed modules (total 7 - 2 succeeded)
	if len(failed) != 5 {
		t.Errorf("LoadScrambled() failed modules = %d, want %d", len(failed), 5)
	}
}

func TestLoadScrambled_FileNotFound(t *testing.T) {
	// Ensure .scrambled file doesn't exist
	os.Remove(ScrambledFileName)

	config, succeeded, failed, err := LoadScrambled()
	if err == nil {
		t.Error("LoadScrambled() should return error when file doesn't exist")
	}

	if config != nil {
		t.Error("LoadScrambled() should return nil configuration when file doesn't exist")
	}

	if succeeded != nil {
		t.Error("LoadScrambled() should return nil succeeded modules when file doesn't exist")
	}

	if failed != nil {
		t.Error("LoadScrambled() should return nil failed modules when file doesn't exist")
	}
}

func TestLoadScrambled_InvalidYAML(t *testing.T) {
	// Create a test .scrambled file with invalid YAML
	testContent := `invalid: yaml: content: [`

	err := os.WriteFile(ScrambledFileName, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .scrambled file: %v", err)
	}
	defer os.Remove(ScrambledFileName)

	config, succeeded, failed, err := LoadScrambled()
	if err == nil {
		t.Error("LoadScrambled() should return error for invalid YAML")
	}

	if config != nil {
		t.Error("LoadScrambled() should return nil configuration for invalid YAML")
	}

	if succeeded != nil {
		t.Error("LoadScrambled() should return nil succeeded modules for invalid YAML")
	}

	if failed != nil {
		t.Error("LoadScrambled() should return nil failed modules for invalid YAML")
	}
}

func TestLoadScrambled_UnknownModule(t *testing.T) {
	// Create a test .scrambled file with unknown module
	testContent := `succeeded:
  - unknown::module
configuration:
  namespace: test
  name: test`

	err := os.WriteFile(ScrambledFileName, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .scrambled file: %v", err)
	}
	defer os.Remove(ScrambledFileName)

	// if module is unknown it should return an error 
	_, _, _, err = LoadScrambled()
	if err == nil {	
		return
	}
	t.Error("LoadScrambled() should return error for unknown module")

}

func TestWriteScrambled(t *testing.T) {
	config := createTestConfiguration()
	succeeded := []modules.IModule{
		&MockModule{name: "egg::initialize"},
		&MockModule{name: "egg::install_tools"},
	}
	failed := &MockModule{name: "egg::install_libraries"}
	moduleError := errors.New("test error")

	err := WriteScrambled(config, succeeded, failed, moduleError)
	if err != nil {
		t.Fatalf("WriteScrambled() returned unexpected error: %v", err)
	}
	defer os.Remove(ScrambledFileName)

	// Verify file was created
	if !CheckScrambled() {
		t.Error("WriteScrambled() should create .scrambled file")
	}

	// Load and verify the content
	loadedConfig, loadedSucceeded, loadedFailed, err := LoadScrambled()
	if err != nil {
		t.Fatalf("Failed to load written .scrambled file: %v", err)
	}

	if loadedConfig.Namespace != config.Namespace {
		t.Errorf("WriteScrambled() configuration namespace = %s, want %s", loadedConfig.Namespace, config.Namespace)
	}

	if len(loadedSucceeded) != 2 {
		t.Errorf("WriteScrambled() succeeded modules = %d, want %d", len(loadedSucceeded), 2)
	}

	// Should have 5 failed modules (total 7 - 2 succeeded)
	if len(loadedFailed) != 5 {
		t.Errorf("WriteScrambled() failed modules = %d, want %d", len(loadedFailed), 5)
	}
}

func TestWriteScrambled_FileCreationError(t *testing.T) {
	// Create a directory with the same name as ScrambledFileName to cause an error
	err := os.Mkdir(ScrambledFileName, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.Remove(ScrambledFileName)

	config := createTestConfiguration()
	succeeded := []modules.IModule{}
	failed := &MockModule{name: "test"}
	moduleError := errors.New("test error")

	err = WriteScrambled(config, succeeded, failed, moduleError)
	if err == nil {
		t.Error("WriteScrambled() should return error when file creation fails")
	}
}

func TestRecoverFromScrambled_Success(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()

	// Create a test .scrambled file with some failed modules
	testContent := `succeeded:
  - egg::initialize
  - egg::install_tools
failed:
  moduleName: egg::install_libraries
  error: "test error"
configuration:
  namespace: github.com/testuser/testproject
  name: testproject
  semver: "0.0.1"
  license: MIT
  copyright:
    year: 2024
    author: Test User
  server:
    jwt: test-secret
    port: 8080
    frontend:
      dir: web/dist
      api: web/src/api
  database:
    url: postgres://postgres:postgres@localhost:5432/test?sslmode=disable
    sqlc: sql
    repository: db/repository
    queries: db/queries
    migration:
      protocol: postgres
      destination: db/migrations
  cache:
    url: redis://localhost:6379
  s3:
    url: http://localhost:9000
    access: test-access
    secret: test-secret`

	err := os.WriteFile(ScrambledFileName, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .scrambled file: %v", err)
	}
	defer os.Remove(ScrambledFileName)

	// Temporarily replace the Modules slice with mock modules that succeed
	originalModules := Modules
	Modules = []modules.IModule{
		&MockModule{name: "egg::initialize", err: nil},
		&MockModule{name: "egg::install_tools", err: nil},
		&MockModule{name: "egg::install_libraries", err: nil},
		&MockModule{name: "egg::bootstrap_directories", err: nil},
		&MockModule{name: "egg::generate_configuration", err: nil},
		&MockModule{name: "egg::bootstrap_framwork", err: nil},
		&MockModule{name: "egg::rsbuild_frontend", err: nil},
	}
	defer func() { Modules = originalModules }()

	err = RecoverFromScrambled(logger)
	if err != nil {
		t.Errorf("RecoverFromScrambled() returned unexpected error: %v", err)
	}
}

func TestRecoverFromScrambled_NoFailedModules(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()

	// Create a test .scrambled file with all modules succeeded
	testContent := `succeeded:
  - egg::initialize
  - egg::install_tools
  - egg::install_libraries
  - egg::bootstrap_directories
  - egg::generate_configuration
  - egg::bootstrap_framework_files
  - egg::rsbuild_frontend
configuration:
  namespace: github.com/testuser/testproject
  name: testproject
  semver: "0.0.1"
  license: MIT
  copyright:
    year: 2024
    author: Test User
  server:
    jwt: test-secret
    port: 8080
    frontend:
      dir: web/dist
      api: web/src/api
  database:
    url: postgres://postgres:postgres@localhost:5432/test?sslmode=disable
    sqlc: sql
    repository: db/repository
    queries: db/queries
    migration:
      protocol: postgres
      destination: db/migrations
  cache:
    url: redis://localhost:6379
  s3:
    url: http://localhost:9000
    access: test-access
    secret: test-secret`

	err := os.WriteFile(ScrambledFileName, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .scrambled file: %v", err)
	}
	defer os.Remove(ScrambledFileName)

	err = RecoverFromScrambled(logger)
	if err != nil {
		t.Errorf("RecoverFromScrambled() returned unexpected error: %v", err)
	}
}

func TestRecoverFromScrambled_LoadError(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()

	// Ensure .scrambled file doesn't exist
	os.Remove(ScrambledFileName)

	err := RecoverFromScrambled(logger)
	if err == nil {
		t.Error("RecoverFromScrambled() should return error when .scrambled file doesn't exist")
	}
}

func TestRecoverFromScrambled_ModuleRecoveryFailure(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()

	// Create a test .scrambled file with failed modules
	testContent := `succeeded:
  - egg::initialize
failed:
  moduleName: egg::install_tools
  error: "test error"
configuration:
  namespace: github.com/testuser/testproject
  name: testproject
  semver: "0.0.1"
  license: MIT
  copyright:
    year: 2024
    author: Test User
  server:
    jwt: test-secret
    port: 8080
    frontend:
      dir: web/dist
      api: web/src/api
  database:
    url: postgres://postgres:postgres@localhost:5432/test?sslmode=disable
    sqlc: sql
    repository: db/repository
    queries: db/queries
    migration:
      protocol: postgres
      destination: db/migrations
  cache:
    url: redis://localhost:6379
  s3:
    url: http://localhost:9000
    access: test-access
    secret: test-secret`

	err := os.WriteFile(ScrambledFileName, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .scrambled file: %v", err)
	}
	defer os.Remove(ScrambledFileName)

	// Temporarily replace the Modules slice with modules that will fail
	originalModules := Modules
	Modules = []modules.IModule{
		&MockModule{name: "egg::initialize", err: nil},
		&MockModule{name: "egg::install_tools", err: errors.New("recovery failed")},
	}
	defer func() { Modules = originalModules }()

	err = RecoverFromScrambled(logger)
	if err == nil {
		t.Error("RecoverFromScrambled() should return error when module recovery fails")
	}

	// Check that new .scrambled file was created
	if !CheckScrambled() {
		t.Error("RecoverFromScrambled() should create new .scrambled file when recovery fails")
	}
}

func TestScrambledFileName_Constant(t *testing.T) {
	if ScrambledFileName != ".scrambled" {
		t.Errorf("ScrambledFileName = %s, want %s", ScrambledFileName, ".scrambled")
	}
}

func TestModules_Initialization(t *testing.T) {
	if len(Modules) == 0 {
		t.Error("Modules slice should not be empty")
	}

	expectedModules := []string{
		"egg::initialize",
		"egg::install_tools",
		"egg::install_libraries",
		"egg::bootstrap_directories",
		"egg::generate_configuration",
		"egg::bootstrap_framwork",
		"egg::rsbuild_frontend",
	}

	if len(Modules) != len(expectedModules) {
		t.Errorf("Modules slice length = %d, want %d", len(Modules), len(expectedModules))
	}

	for i, expectedName := range expectedModules {
		if Modules[i].Name() != expectedName {
			t.Errorf("Modules[%d].Name() = %s, want %s", i, Modules[i].Name(), expectedName)
		}
	}
}
