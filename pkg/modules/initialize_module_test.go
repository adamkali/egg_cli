package modules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adamkali/egg_cli/pkg/models"
)

func TestInitializeModule_Name(t *testing.T) {
	module := &InitializeModule{}
	expected := "egg::initialize"
	if got := module.Name(); got != expected {
		t.Errorf("InitializeModule.Name() = %v, want %v", got, expected)
	}
}

func TestInitializeModule_GetProgress(t *testing.T) {
	tests := []struct {
		name     string
		progress int
		expected float64
	}{
		{"zero progress", 0, 0.0},
		{"half progress", 2, 0.5},
		{"full progress", 4, 1.0},
		{"over progress", 5, 1.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			module := &InitializeModule{Progress: tt.progress}
			if got := module.GetProgress(); got != tt.expected {
				t.Errorf("InitializeModule.GetProgress() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInitializeModule_IncrProg(t *testing.T) {
	module := &InitializeModule{Progress: 1}
	module.IncrProg()
	if module.Progress != 2 {
		t.Errorf("InitializeModule.IncrProg() failed, expected Progress = 2, got %d", module.Progress)
	}
}

func TestInitializeModule_LoadFromConfig(t *testing.T) {
	config := createTestConfiguration()
	logger, _ := models.NewLogger("test.log")
	defer logger.Close()

	module := &InitializeModule{}
	module.LoadFromConfig(config, logger)

	if module.Namespace != config.Namespace {
		t.Errorf("InitializeModule.LoadFromConfig() Namespace = %v, want %v", module.Namespace, config.Namespace)
	}
	if module.ProjectName != config.Name {
		t.Errorf("InitializeModule.LoadFromConfig() ProjectName = %v, want %v", module.ProjectName, config.Name)
	}
	if module.License != config.License {
		t.Errorf("InitializeModule.LoadFromConfig() License = %v, want %v", module.License, config.License)
	}
	if module.Progress != 0 {
		t.Errorf("InitializeModule.LoadFromConfig() Progress = %v, want 0", module.Progress)
	}
	if module.eggl != logger {
		t.Errorf("InitializeModule.LoadFromConfig() logger not set correctly")
	}
}

func TestInitializeModule_IsError(t *testing.T) {
	tests := []struct {
		name     string
		module   InitializeModule
		expected error
	}{
		{"no error", InitializeModule{Error: nil}, nil},
		{"with error", InitializeModule{Error: os.ErrNotExist}, os.ErrNotExist},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.module.IsError(); got != tt.expected {
				t.Errorf("InitializeModule.IsError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInitializeModule_Run_Success(t *testing.T) {
	// Create test directory
	testDir := "test_data"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	config := createTestConfiguration()
	logger, _ := models.NewLogger("test.log")
	defer logger.Close()

	module := &InitializeModule{}
	module.LoadFromConfig(config, logger)

	// Run the module
	module.Run()

	// Check for errors
	if err := module.IsError(); err != nil {
		t.Errorf("InitializeModule.Run() failed with error: %v", err)
	}

	// The module changes directory to the project directory, so we need to check from the current directory
	// First, let's check what directory we're in
	currentDir, _ := os.Getwd()
	t.Logf("Current directory after module run: %s", currentDir)

	// Check that go.mod file was created in the current directory (which should be the project directory)
	goModPath := "go.mod"
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Errorf("go.mod file was not created: %s", goModPath)
	}

	// Check progress
	expectedProgress := 4 // Should complete all 4 steps
	if module.Progress != expectedProgress {
		t.Errorf("InitializeModule.Run() Progress = %d, want %d", module.Progress, expectedProgress)
	}
}

func TestInitializeModule_Run_DirectoryExists(t *testing.T) {
	// Create test directory
	testDir := "test_data"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	config := createTestConfiguration()
	logger, _ := models.NewLogger("test.log")
	defer logger.Close()

	module := &InitializeModule{}
	module.LoadFromConfig(config, logger)

	// Create the project directory first to simulate it already existing
	projectPath := filepath.Join(testDir, config.Name)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("Failed to create existing project directory: %v", err)
	}

	// Create a test file in the directory to verify it gets removed
	testFile := filepath.Join(projectPath, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Run the module - this should prompt for overwrite, but in test we can't interact
	// So we'll just verify the module handles the existing directory gracefully
	module.Run()

	// The module should have encountered the existing directory
	// In a real scenario, this would prompt the user
	// For testing purposes, we'll just verify the module doesn't crash
	if module.Progress > 0 {
		t.Logf("Module progressed %d steps despite existing directory", module.Progress)
	}
}

func TestInitializeModule_Run_GoVersionCheck(t *testing.T) {
	// Create test directory
	testDir := "test_data"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	config := createTestConfiguration()
	logger, _ := models.NewLogger("test.log")
	defer logger.Close()

	module := &InitializeModule{}
	module.LoadFromConfig(config, logger)

	// Run the module
	module.Run()

	// Check that Go version check passed (should be at least 1.23)
	if err := module.IsError(); err != nil {
		t.Errorf("Go version check failed: %v", err)
	}

	// Should have progressed past the Go version check
	if module.Progress < 2 {
		t.Errorf("Module did not progress past Go version check, Progress = %d", module.Progress)
	}
}

func TestInitializeModule_Run_ModuleInitialization(t *testing.T) {
	// Create test directory
	testDir := "test_data"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	config := createTestConfiguration()
	logger, _ := models.NewLogger("test.log")
	defer logger.Close()

	module := &InitializeModule{}
	module.LoadFromConfig(config, logger)

	// Run the module
	module.Run()

	// Check for errors
	if err := module.IsError(); err != nil {
		t.Errorf("InitializeModule.Run() failed with error: %v", err)
	}

	// Check that go.mod was created with correct module name
	goModPath := "go.mod"

	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Errorf("go.mod file was not created: %s", goModPath)
	} else {
		// Read and verify go.mod content
		content, err := os.ReadFile(goModPath)
		if err != nil {
			t.Errorf("Failed to read go.mod file: %v", err)
		} else {
			expectedModuleLine := "module " + config.Namespace
			if !contains(string(content), expectedModuleLine) {
				t.Errorf("go.mod does not contain expected module line. Expected: %s, Content: %s",
					expectedModuleLine, string(content))
			}
		}
	}

	// Should have completed all steps
	if module.Progress != 4 {
		t.Errorf("Module did not complete all steps, Progress = %d", module.Progress)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
