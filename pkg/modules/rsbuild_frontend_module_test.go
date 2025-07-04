package modules

import (
	"errors"
	"testing"
)

func TestRsbuildFrontendModule_Name(t *testing.T) {
	m := &RsbuildFrontendModule{}
	if got, want := m.Name(), "egg::rsbuild_frontend"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestRsbuildFrontendModule_LoadFromConfig(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	cfg := createTestConfiguration()
	m := &RsbuildFrontendModule{}
	m.LoadFromConfig(cfg, logger)
	if m.configuration != cfg {
		t.Error("LoadFromConfig did not set configuration")
	}
	if m.eggl != logger {
		t.Error("LoadFromConfig did not set logger")
	}
}

func TestRsbuildFrontendModule_IsError(t *testing.T) {
	m := &RsbuildFrontendModule{error: errors.New("fail")}
	if m.IsError() == nil {
		t.Error("IsError() = nil, want error")
	}
	m.error = nil
	if m.IsError() != nil {
		t.Error("IsError() != nil, want nil")
	}
}

func TestRsbuildFrontendModule_Run_UserSaysNo(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	cfg := createTestConfiguration()
	m := &RsbuildFrontendModule{configuration: cfg, eggl: logger}

	m.InputFunc = func(prompt string) string {
		if prompt == "Do you want to use a JavaScript Framework for the frontend? (y/n)" {
			return "n"
		}
		return ""
	}

	m.Run()
	if m.IsError() != nil {
		t.Errorf("Run() set unexpected error: %v", m.IsError())
	}
}

func TestRsbuildFrontendModule_Run_Success(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	cfg := createTestConfiguration()
	m := &RsbuildFrontendModule{configuration: cfg, eggl: logger}

	inputCount := 0
	m.InputFunc = func(prompt string) string {
		inputCount++
		if inputCount == 1 {
			return "y" // Use JS
		}
		return "pnpm" // Use pnpm
	}

	execCalled := false
	m.ExecFunc = func(cmd string, args ...string) ([]byte, error) {
		execCalled = true
		return []byte("rsbuild success"), nil
	}

	m.Run()
	if m.IsError() != nil {
		t.Errorf("Run() set unexpected error: %v", m.IsError())
	}
	if !execCalled {
		t.Error("ExecFunc was not called for rsbuild install")
	}
}

func TestRsbuildFrontendModule_Run_ExecError(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	cfg := createTestConfiguration()
	m := &RsbuildFrontendModule{configuration: cfg, eggl: logger}

	inputCount := 0
	m.InputFunc = func(prompt string) string {
		inputCount++
		if inputCount == 1 {
			return "y" // Use JS
		}
		return "pnpm" // Use pnpm
	}

	m.ExecFunc = func(cmd string, args ...string) ([]byte, error) {
		return nil, errors.New("simulated exec error")
	}

	m.Run()
	if m.IsError() == nil {
		t.Error("Run() did not set error on exec failure")
	}
}

func TestRsbuildFrontendModule_Run_InvalidPackageManager(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	cfg := createTestConfiguration()
	m := &RsbuildFrontendModule{configuration: cfg, eggl: logger}

	inputCount := 0
	m.InputFunc = func(prompt string) string {
		inputCount++
		if inputCount == 1 {
			return "y" // Use JS
		}
		return "invalid" // Invalid package manager
	}

	m.Run()
	if m.IsError() == nil {
		t.Error("Run() did not set error on invalid package manager")
	}
}
