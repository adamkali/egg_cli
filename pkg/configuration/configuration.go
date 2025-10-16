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

package configuration

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
	Semver    string `yaml:"semver"`
	License   string `yaml:"license"`
	Copyright struct {
		Year   int    `yaml:"year"`
		Author string `yaml:"author"`
	} `yaml:"copyright"`
	Server struct {
		JWT      string `yaml:"jwt"`
		Port     int    `yaml:"port"`
		Frontend struct {
			Dir string `yaml:"dir"`
			Api string `yaml:"api"`
		} `yaml:"frontend"`
	} `yaml:"server"`
	Database struct {
		URL  string `yaml:"url"`
		Sqlc struct {
			RepositoryLocation string `yaml:"repository"`
			Schema             string `yaml:"schema"`
			SqlOrGo            string `yaml:"sql_or_go"`
		} `yaml:"sqlc"`
		QueriesLocation string `yaml:"queries"`
		Migration       struct {
			Protocol    string `yaml:"protocol"`
			Destination string `yaml:"destination"`
		} `yaml:"migration"`
	} `yaml:"database"`
	Cache struct {
		URL string `yaml:"url"`
	} `yaml:"cache"`
	S3 struct {
		URL    string `yaml:"url"`
		Access string `yaml:"access"`
		Secret string `yaml:"secret"`
	} `yaml:"s3"`
}

func (configuration *Configuration) SaveConfiguration(configFile string) error {
	configurationFile := configFile
	configBytes, err := yaml.Marshal(configuration)
	if err != nil {
		return err
	}

	// create the config file if not exists
	if _, err := os.Stat(configurationFile); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(configurationFile); err != nil {
			return err
		}
	}

	// save the configuration
	if err := os.WriteFile(configurationFile, configBytes, 0777); err != nil {
		return err
	}
	return nil
}

const ConfigurationDir = "config/"

func LoadConfigurationFile(configurationFile string) (*Configuration, error) {
	configuration := new(Configuration)
	file, err := os.ReadFile(configurationFile)
	if err != nil {
		return configuration, err
	}
	if err = yaml.Unmarshal(file, configuration); err != nil {
		return configuration, err
	}
	return configuration, nil
}

func LoadConfigurationReader(configBytes io.Reader) (*Configuration, error) {
	configuration := new(Configuration)
	bytes, err := io.ReadAll(configBytes)
	if err != nil {
		return configuration, err
	}
	if err := yaml.Unmarshal(bytes, configuration); err != nil {
		return configuration, err
	}
	return configuration, nil
}

func LoadConfiguration(environment string) (*Configuration, error) {
	configuration := new(Configuration)
	configurationFile := ConfigurationDir + environment + ".yaml"
	file, err := os.ReadFile(configurationFile)
	if err != nil {
		return configuration, err
	}
	if err = yaml.Unmarshal(file, configuration); err != nil {
		return configuration, err
	}
	return configuration, nil
}

func SaveConfiguration(configBytes []byte, environment string) error {
	configurationFile := ConfigurationDir + environment + ".yaml"
	if _, err := os.Stat(configurationFile); errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.WriteFile(configurationFile, configBytes, 0777); err != nil {
		return err
	}
	return nil
}

func SaveConfigurationWriter(configBytes []byte, configWriter io.Writer) error {
	if _, err := configWriter.Write(configBytes); err != nil {
		return err
	}
	return nil
}

func (configuration *Configuration) GenerateConfigurationFile(environment string) error {
	// create the config directory if not exists config/
	if _, err := os.Stat(ConfigurationDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(ConfigurationDir, 0777); err != nil {
			return err
		}
	}
	if _, err := os.Stat(ConfigurationDir + environment + ".yaml"); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(ConfigurationDir + environment + ".yaml"); err != nil {
			return err
		}
	}
	// write to the file with the yaml content as the configuration
	configBytes, err := yaml.Marshal(configuration)
	if err != nil {
		return err
	}
	if err := SaveConfiguration(configBytes, environment); err != nil {
		return err
	}
	return nil
}

func (configuration *Configuration) IntoGeneratorConfig(version string) *GeneratorConfiguration {
	genConfig := new(GeneratorConfiguration).Init(version)
	genConfig.Namespace = configuration.Namespace
	genConfig.Name = configuration.Name
	genConfig.Semver = configuration.Semver
	genConfig.License = configuration.License
	genConfig.Copyright.Author = configuration.Copyright.Author
	genConfig.Copyright.Year = configuration.Copyright.Year
	genConfig.Server.JWT = configuration.Server.JWT
	genConfig.Server.Port = configuration.Server.Port
	genConfig.Server.Frontend.Dir = configuration.Server.Frontend.Dir
	genConfig.Server.Frontend.Api = configuration.Server.Frontend.Api
	genConfig.PostgresUrl(configuration.Database.URL)
	genConfig.Database.Sqlc.RepositoryLocation = configuration.Database.Sqlc.RepositoryLocation
	genConfig.Database.Sqlc.Schema = configuration.Database.Sqlc.Schema
	genConfig.Database.Sqlc.SqlOrGo = configuration.Database.Sqlc.SqlOrGo
	genConfig.Database.QueriesLocation = configuration.Database.QueriesLocation
	genConfig.Database.Migration.Destination = configuration.Database.Migration.Destination
	genConfig.Database.Migration.Protocol = configuration.Database.Migration.Protocol
	genConfig.CacheUrl(configuration.Cache.URL)
	genConfig.S3.URL = configuration.S3.URL
	genConfig.S3.Access = configuration.S3.Access
	genConfig.S3.Secret = configuration.S3.Secret
	return genConfig
}
