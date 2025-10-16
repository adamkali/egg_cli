// This module is used to bootstrap the framework files from the templates found in the templates directory
// it will try by using the configuration file to determine what values to substitute into the templates
// and use the templates/mapping function to iterate over the templates and output the files to the correct location
// and provide the logging names when printing to stdout

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
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	tgts "github.com/adamkali/egg_cli/pkg/targets"
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
	mapping               map[string]*template.Template
	configuration         *configuration.Configuration
	error                 error
	progress              int
	eggl                  *models.EggLog
	PopulateTemplatesFunc func(name string, template *template.Template) error // For testing - can be injected to mock template population
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
//
//	This function is used to bootstrap the framework files from the templates found in the templates directory
//	this is done by iterating over the mapping and creating a goroutine for each template
//	and then waiting for all the goroutines to finish
//	and collecting the errors
//	and logging the errors
func (m *BootstrapFrameworkFilesFromTemplatesModule) Run() {
	m.progress = 0
	if m.error != nil {
		m.eggl.Error(m.error.Error())
		return
	}
	for name, templ := range m.mapping {
		if m.PopulateTemplatesFunc != nil {
			m.error = m.PopulateTemplatesFunc(name, templ)
		} else {
			m.error = m.populateTemplate(name, templ)
		}
		if m.error != nil {
			m.eggl.Error(m.error.Error())
			return
		}
	}

	return
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
	log := fmt.Sprintf("🥚 %s creating %s", m.Name(), name)
	m.eggl.Info(log)
	fmt.Println(styles.EggProgressInfo.Render(log))
	var err error
	var f *os.File
	f, err = os.Create(name)
	if err != nil {
		relativeDir := path.Dir(name)
		// if relativeDir is not empty then we need to create the directory
		if relativeDir != "" {
			err := os.MkdirAll(relativeDir, os.ModePerm)
			log:= fmt.Sprintf("🥚 %s creating %s", m.Name(), relativeDir)
			m.eggl.Info(log)
			fmt.Println(styles.EggProgressInfo.Render(log))
			if err != nil {
				err = errors.New(m.Name() + "error creating directory: " + relativeDir + " " + err.Error())
				return err
			}
			f, err = os.Create(name)
			if err != nil {
				err = errors.New(m.Name() + "error creating file: " + name + " " + err.Error())
				return err
			}
		}
		if err != nil {
			err = errors.New(m.Name() + "error creating file: " + name + " " + err.Error())
			return err
		}
	}
	defer f.Close()
	err = template.Execute(f, m.configuration)
	if err != nil {
		print("%s error executing template: %s %s\n", m.Name(), name, err.Error())
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
	m.mapping, m.error = tgts.Mapping(configuration, templates.GetTemplates())
	m.eggl = eggl
	return
}
