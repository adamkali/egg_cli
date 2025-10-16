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
	"errors"
	"testing"
)

func TestGenerateConfigurationModule_Name(t *testing.T) {
	m := &GenerateConfigurationModule{}
	if got, want := m.Name(), "egg::generate_configuration"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestGenerateConfigurationModule_LoadFromConfig(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	cfg := createTestConfiguration()
	m := &GenerateConfigurationModule{}
	m.LoadFromConfig(cfg, logger)
	if m.Configuration != cfg {
		t.Error("LoadFromConfig did not set configuration")
	}
	if m.eggl != logger {
		t.Error("LoadFromConfig did not set logger")
	}
	if m.Progress != 0 {
		t.Errorf("LoadFromConfig did not reset Progress, got %d", m.Progress)
	}
}

func TestGenerateConfigurationModule_IsError(t *testing.T) {
	m := &GenerateConfigurationModule{Error: errors.New("fail")}
	if m.IsError() == nil {
		t.Error("IsError() = nil, want error")
	}
	m.Error = nil
	if m.IsError() != nil {
		t.Error("IsError() != nil, want nil")
	}
}

func TestGenerateConfigurationModule_Run_Success(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	cfg := createTestConfiguration()
	m := &GenerateConfigurationModule{Configuration: cfg, eggl: logger}

	// Inject mock config generation function
	m.GenerateConfigFunc = func(environment string) error {
		return nil
	}

	m.Run()
	if m.IsError() != nil {
		t.Errorf("Run() set unexpected error: %v", m.IsError())
	}
}

func TestGenerateConfigurationModule_Run_Error(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	cfg := createTestConfiguration()
	m := &GenerateConfigurationModule{Configuration: cfg, eggl: logger}

	// Inject mock config generation function that returns error
	m.GenerateConfigFunc = func(environment string) error {
		return errors.New("simulated error")
	}

	m.Run()
	if m.IsError() == nil {
		t.Error("Run() did not set error on failure")
	}
}
