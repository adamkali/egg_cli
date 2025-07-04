package modules

import (
	"errors"
	"testing"

	"github.com/adamkali/egg_cli/pkg/configuration"
)

func TestInstallLibrariesModule_Name(t *testing.T) {
	m := &InstallLibrariesModule{}
	if got, want := m.Name(), "egg::install_libraries"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestInstallLibrariesModule_GetProgressAndIncrProg(t *testing.T) {
	m := &InstallLibrariesModule{Progress: 0}
	if m.GetProgress() != 0.0 {
		t.Errorf("GetProgress() = %v, want 0.0", m.GetProgress())
	}
	m.IncrProg()
	if m.Progress != 1 {
		t.Errorf("IncrProg() did not increment Progress, got %d", m.Progress)
	}
}

func TestInstallLibrariesModule_LoadFromConfig(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	m := &InstallLibrariesModule{}
	m.LoadFromConfig(&configuration.Configuration{}, logger)
	if m.eggl != logger {
		t.Error("LoadFromConfig did not set logger")
	}
	if m.Progress != 0 {
		t.Errorf("LoadFromConfig did not reset Progress, got %d", m.Progress)
	}
}

func TestInstallLibrariesModule_IsError(t *testing.T) {
	m := &InstallLibrariesModule{Error: errors.New("fail")}
	if m.IsError() == nil {
		t.Error("IsError() = nil, want error")
	}
	m.Error = nil
	if m.IsError() != nil {
		t.Error("IsError() != nil, want nil")
	}
}

func TestInstallLibrariesModule_Run_Simulated(t *testing.T) {
	// Use a test logger and inject GoGetFunc to avoid real installs
	logger := createTestLogger(t)
	defer logger.Close()
	m := &InstallLibrariesModule{eggl: logger}

	// Simulate GoGet always succeeds
	m.GoGetFunc = func(pac string) error {
		m.IncrProg()
		return nil
	}

	m.Run()
	if m.IsError() != nil {
		t.Errorf("Run() set unexpected error: %v", m.IsError())
	}
	if m.Progress == 0 {
		t.Error("Run() did not increment Progress")
	}
}

func TestInstallLibrariesModule_Run_SimulatedError(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	m := &InstallLibrariesModule{eggl: logger}

	// Simulate GoGet fails on first package
	m.GoGetFunc = func(pac string) error {
		m.Error = errors.New("simulated error")
		return m.Error
	}

	m.Run()
	if m.IsError() == nil {
		t.Error("Run() did not set error on GoGet failure")
	}

}
