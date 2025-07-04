package modules

import (
	"errors"
	"testing"

	"github.com/adamkali/egg_cli/pkg/configuration"
)

func TestInstallToolsModule_Name(t *testing.T) {
	m := &InstallToolsModule{}
	if got, want := m.Name(), "egg::install_tools"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestInstallToolsModule_GetProgressAndIncrProg(t *testing.T) {
	m := &InstallToolsModule{Progress: 0}
	if m.GetProgress() != 0.0 {
		t.Errorf("GetProgress() = %v, want 0.0", m.GetProgress())
	}
	m.IncrProg()
	if m.Progress != 1 {
		t.Errorf("IncrProg() did not increment Progress, got %d", m.Progress)
	}
}

func TestInstallToolsModule_LoadFromConfig(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	m := &InstallToolsModule{}
	m.LoadFromConfig(&configuration.Configuration{}, logger)
	if m.eggl != logger {
		t.Error("LoadFromConfig did not set logger")
	}
	if m.Progress != 0 {
		t.Errorf("LoadFromConfig did not reset Progress, got %d", m.Progress)
	}
}

func TestInstallToolsModule_IsError(t *testing.T) {
	m := &InstallToolsModule{Error: errors.New("fail")}
	if m.IsError() == nil {
		t.Error("IsError() = nil, want error")
	}
	m.Error = nil
	if m.IsError() != nil {
		t.Error("IsError() != nil, want nil")
	}
}

func TestInstallToolsModule_Run_AllToolsAlreadyInstalled(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	m := &InstallToolsModule{eggl: logger}

	// Mock LookPathFunc to simulate all tools already installed
	m.LookPathFunc = func(file string) (string, error) {
		return "/usr/local/bin/" + file, nil // Simulate tool found
	}

	m.Run()
	if m.IsError() != nil {
		t.Errorf("Run() set unexpected error: %v", m.IsError())
	}
	// Should not have incremented progress since no tools were installed
	if m.Progress != 0 {
		t.Errorf("Progress should be 0 when all tools are already installed, got %d", m.Progress)
	}
}

func TestInstallToolsModule_Run_InstallSomeTools(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	m := &InstallToolsModule{eggl: logger}

	// Mock LookPathFunc to simulate some tools not found
	toolsNotFound := 0
	m.LookPathFunc = func(file string) (string, error) {
		if toolsNotFound < 2 {
			toolsNotFound++
			return "", errors.New("tool not found")
		}
		return "/usr/local/bin/" + file, nil // Simulate tool found
	}

	// Mock InstallToolFunc to simulate successful installation
	m.InstallToolFunc = func(tool string) error {
		return nil
	}

	m.Run()
	if m.IsError() != nil {
		t.Errorf("Run() set unexpected error: %v", m.IsError())
	}
	// Should have incremented progress for each tool installed
	if m.Progress < 1 {
		t.Errorf("Progress should be incremented for installed tools, got %d", m.Progress)
	}
}

func TestInstallToolsModule_Run_InstallError(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	m := &InstallToolsModule{eggl: logger}

	// Mock LookPathFunc to simulate tool not found
	m.LookPathFunc = func(file string) (string, error) {
		return "", errors.New("tool not found")
	}

	// Mock InstallToolFunc to simulate installation error
	m.InstallToolFunc = func(tool string) error {
		return errors.New("installation failed")
	}

	m.Run()
	if m.IsError() == nil {
		t.Error("Run() did not set error on installation failure")
	}
}

func TestInstallToolsModule_Run_MixedScenario(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	m := &InstallToolsModule{eggl: logger}

	// Mock LookPathFunc to simulate mixed scenario
	toolCount := 0
	m.LookPathFunc = func(file string) (string, error) {
		toolCount++
		if toolCount <= 2 {
			return "/usr/local/bin/" + file, nil // First 2 tools already installed
		}
		return "", errors.New("tool not found") // Rest need installation
	}

	// Mock InstallToolFunc to simulate successful installation
	installCount := 0
	m.InstallToolFunc = func(tool string) error {
		installCount++
		return nil
	}

	m.Run()
	if m.IsError() != nil {
		t.Errorf("Run() set unexpected error: %v", m.IsError())
	}
	// Should have incremented progress for each tool that needed installation
	if m.Progress < 1 {
		t.Errorf("Progress should be incremented for installed tools, got %d", m.Progress)
	}
}
