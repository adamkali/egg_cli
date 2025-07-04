package modules

import (
	"os"
	"strings"
	"testing"
	"text/template"
	"time"
)

func TestBootstrapFrameworkFilesFromTemplatesModule_Name(t *testing.T) {
	m := &BootstrapFrameworkFilesFromTemplatesModule{}
	if got, want := m.Name(), "egg::bootstrap_framwork"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestBootstrapFrameworkFilesFromTemplatesModule_ProgressAndIncrProg(t *testing.T) {
	m := &BootstrapFrameworkFilesFromTemplatesModule{mapping: map[string]*template.Template{"a": nil, "b": nil}, progress: 0}
	if m.GetProgress() != 0.0 {
		t.Errorf("GetProgress() = %v, want 0.0", m.GetProgress())
	}
	m.IncrProg()
	if m.progress != 1 {
		t.Errorf("IncrProg() did not increment progress, got %d", m.progress)
	}
}

func TestBootstrapFrameworkFilesFromTemplatesModule_LoadFromConfig(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	m := &BootstrapFrameworkFilesFromTemplatesModule{}
	cfg := createTestConfiguration()
	m.LoadFromConfig(cfg, logger)
	if m.configuration != cfg {
		t.Error("LoadFromConfig did not set configuration")
	}
	if m.eggl != logger {
		t.Error("LoadFromConfig did not set logger")
	}
}

func TestBootstrapFrameworkFilesFromTemplatesModule_IsError(t *testing.T) {
	m := &BootstrapFrameworkFilesFromTemplatesModule{error: os.ErrNotExist}
	if m.IsError() == nil {
		t.Error("IsError() = nil, want error")
	}
	m.error = nil
	if m.IsError() != nil {
		t.Error("IsError() != nil, want nil")
	}
}

func TestBootstrapFrameworkFilesFromTemplatesModule_Run(t *testing.T) {
	logger := createTestLogger(t)
	defer logger.Close()
	cfg := createTestConfiguration()

	// Create a simple template
	tmpl := template.Must(template.New("testfile.txt").Parse("Hello, {{.Name}}!"))
	mapping := map[string]*template.Template{
		"testfile.txt": tmpl,
	}

	// Track generated content for verification
	generatedContent := make(map[string]string)

	m := &BootstrapFrameworkFilesFromTemplatesModule{
		mapping:       mapping,
		configuration: cfg,
		eggl:          logger,
		PopulateTemplatesFunc: func(name string, template *template.Template) error {
			// Instead of writing to file, capture the generated content
			var buf strings.Builder
			err := template.Execute(&buf, cfg)
			if err != nil {
				return err
			}
			generatedContent[name] = buf.String()
			return nil
		},
	}

	m.Run()

	// Wait a bit for goroutines to complete (since Run() is asynchronous)
	time.Sleep(100 * time.Millisecond)

	// Check that the content was generated correctly
	expectedContent := "Hello, testproject!"
	if content, exists := generatedContent["testfile.txt"]; !exists {
		t.Error("Template was not processed")
	} else if content != expectedContent {
		t.Errorf("Generated content = %q, want %q", content, expectedContent)
	}

	if m.IsError() != nil {
		t.Errorf("Run() set unexpected error: %v", m.IsError())
	}
}
