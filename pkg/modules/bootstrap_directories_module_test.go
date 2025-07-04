package modules

import (
	"errors"
	"testing"
	"time"
)

func TestBootstrapDirectoriesModule_Name(t *testing.T) {
	m := &BootstrapDirectoriesModule{}
	if got, want := m.Name(), "egg::bootstrap_directories"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestBootstrapDirectoriesModule_GetProgressAndIncrProg(t *testing.T) {
	m := &BootstrapDirectoriesModule{Directories: []string{"a", "b"}, Progress: 0}
	if m.GetProgress() != 0.0 {
		t.Errorf("GetProgress() = %v, want 0.0", m.GetProgress())
	}
	m.IncrProg()
	if m.Progress != 1 {
		t.Errorf("IncrProg() did not increment Progress, got %d", m.Progress)
	}
}

func TestBootstrapDirectoriesModule_LoadFromConfig(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	cfg := createTestConfiguration()
	m := &BootstrapDirectoriesModule{}
	m.LoadFromConfig(cfg, logger)
	if m.eggl != logger {
		t.Error("LoadFromConfig did not set logger")
	}
	if m.Progress != 0 {
		t.Errorf("LoadFromConfig did not reset Progress, got %d", m.Progress)
	}
	if len(m.Directories) == 0 {
		t.Error("LoadFromConfig did not set Directories")
	}
}

func TestBootstrapDirectoriesModule_IsError(t *testing.T) {
	m := &BootstrapDirectoriesModule{Error: errors.New("fail")}
	if m.IsError() == nil {
		t.Error("IsError() = nil, want error")
	}
	m.Error = nil
	if m.IsError() != nil {
		t.Error("IsError() != nil, want nil")
	}
}

func TestBootstrapDirectoriesModule_Run_Success(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	m := &BootstrapDirectoriesModule{eggl: logger, Directories: []string{"dir1", "dir2"}}

	created := make(map[string]bool)
	m.MkdirFunc = func(dir string) error {
		created[dir] = true
		return nil
	}

	m.Run()
	time.Sleep(50 * time.Millisecond)

	if m.IsError() != nil {
		t.Errorf("Run() set unexpected error: %v", m.IsError())
	}
	for _, dir := range m.Directories {
		if !created[dir] {
			t.Errorf("Directory %q was not created", dir)
		}
	}
}

func TestBootstrapDirectoriesModule_Run_Error(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	m := &BootstrapDirectoriesModule{eggl: logger, Directories: []string{"dir1", "dir2"}}

	errorDir := "dir2"
	m.MkdirFunc = func(dir string) error {
		if dir == errorDir {
			return errors.New("simulated error")
		}
		return nil
	}

	m.Run()
	time.Sleep(50 * time.Millisecond)

	if m.IsError() == nil {
		t.Error("Run() did not set error on directory creation failure")
	}
}
