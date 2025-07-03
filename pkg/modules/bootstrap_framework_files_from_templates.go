// This module is used to bootstrap the framework files from the templates found in the templates directory
// it will try by using the configuration file to determine what values to substitute into the templates
// and use the templates/mapping function to iterate over the templates and output the files to the correct location
// and provide the logging names when printing to stdout

package modules

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"text/template"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/models"
	"github.com/adamkali/egg_cli/pkg/templates"
	"github.com/adamkali/egg_cli/styles"
)

// BootstrapFrameworkFilesFromTemplatesModule
//
// implements:
//   - IModule
//
// description:
//
//	This struct is used to bootstrap the framework files from the templates found in the templates directory
//	this is done by iterating over the mapping and creating a goroutine for each template
//	and then waiting for all the goroutines to finish
//	and collecting the errors
//	and logging the errors
//	canceling the context if there is an error
type BootstrapFrameworkFilesFromTemplatesModule struct {
	mapping       map[string]*template.Template
	configuration *configuration.Configuration
	error         error
	progress      int
	eggl          *models.EggLog
}

// Name
//
// returns:
//
//	string: the name
//
// description:
//
//	This function returns the name of the module
func (m *BootstrapFrameworkFilesFromTemplatesModule) Name() string {
	return "egg::bootstrap_framwork"
}

// IncrProg
//
// description:
//
//	increments the progress
func (m *BootstrapFrameworkFilesFromTemplatesModule) IncrProg() {
	m.progress += 1
	return
}

// GetProgress
//
// returns:
//
//	float64: the progress
func (m *BootstrapFrameworkFilesFromTemplatesModule) GetProgress() float64 {
	return float64(m.progress) / float64(len(m.mapping))
}

// Run
//
// description:
//   This function is used to bootstrap the framework files from the templates found in the templates directory
//   this is done by iterating over the mapping and creating a goroutine for each template
//   and then waiting for all the goroutines to finish
//   and collecting the errors
//   and logging the errors
func (m *BootstrapFrameworkFilesFromTemplatesModule) Run() {
	m.progress = 0
	// iterate over the mapping and create a goroutine for each template
	errChan := make(chan error)
	logChan := make(chan string)
	wg := new(sync.WaitGroup)
	for name, templ := range m.mapping {
		wg.Add(1)
		go func(name string, t *template.Template) {
			defer wg.Done()
			err := m.populateTemplate(name, t)
			logChan <- fmt.Sprintf("🥚 %s creating %s", m.Name(), name)
			if err != nil {
				errChan <- err
			}
		}(name, templ)
	}

	// handle logging
	go func() {
		for log := range logChan {
			m.eggl.Info(log)
			log = styles.EggProgressInfo.Render(log)
			fmt.Println(log)
		}
	}()

	handleErrors := func() {
		for err := range errChan {
			m.error = err
			// return if there is an error
			wg.Wait()
			close(errChan)
			close(logChan)
			return
		}
	}

	// handle errors
	go func() {
		handleErrors()
	}()

	go func() {
		wg.Wait()
		close(errChan)
		close(logChan)
	}()
}

// IsError
//
// returns:
//
//	error:
func (m *BootstrapFrameworkFilesFromTemplatesModule) IsError() error {
	return m.error
}

// populateTemplate
//
// params:
//
//	name: string
//	template: *template.Template
//
// returns:
//
//	  error:
//	    - if there is an error creating the file
//		   - if there is an error executing the template
//	    - if the generater failed to use a correct directory
//
// description:
//
//	This function is used to populate the template with the values from the configuration file
//	and output the file to the correct location. Because template.Mapping is created with keys
//	that are the same as the file names, this function is used as a single instance the mapping
//	so that the loop can be split into a goroutine and be ran concurrently. We also return an error
//	so that we can collect all the errors into a channel and stop the downstream goroutines
func (m *BootstrapFrameworkFilesFromTemplatesModule) populateTemplate(name string, template *template.Template) error {
	// create an io writer to write the file to
	var err error
	var f *os.File
	f, err = os.Create(name)
	if err != nil {
		// if m.error is PathError then we need to make sure that the directory exists
		// we seperate the name as the relative directory in which the file is located
		// Example:
		//   name: "./cmd/configuration/configuration.go"
		//   relativeDir: "<configuration.Name>/cmd/configuration"
		//   fileName: "configuration.go"
		//   dir: "./cmd/configuration"
		relativeDir := path.Dir(name)
		// if relativeDir is not empty then we need to create the directory
		if relativeDir != "" {
			err := os.MkdirAll(relativeDir, os.ModePerm)
			if err != nil {
				err = errors.New(m.Name() + "error creating directory: " + relativeDir + " " + err.Error())
				return err
			}
			f, err = os.Create(name)
		}
	}
	defer f.Close()

	// execute the template
	err = template.Execute(f, m.configuration)
	if err != nil {
		err = errors.New(m.Name() + "error executing template: " + name + " " + err.Error())
		return err
	}
	return nil
}

// LoadFromConfig
//
// params:
//
//	configuration: *configuration.Configuration
//	eggl: *models.EggLog
//
// description:
//
//	This function is used to load the module from the configuration and sets up the empty structure so that it can be used
//	by m.Run().
func (m *BootstrapFrameworkFilesFromTemplatesModule) LoadFromConfig(configuration *configuration.Configuration, eggl *models.EggLog) {
	m.configuration = configuration
	m.mapping = templates.Mapping(configuration)
	m.eggl = eggl
	return
}
