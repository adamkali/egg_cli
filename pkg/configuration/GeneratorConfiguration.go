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
	"fmt"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"github.com/adamkali/egg_cli/styles"
)

type GeneratorConfiguration struct {
	EggVersion string `yaml:"egg_version"`
	Namespace  string `yaml:"namespace"`
	Name       string `yaml:"name"`
	Semver     string `yaml:"semver"`
	License    string `yaml:"license"`
	Copyright  struct {
		Year   int    `yaml:"year"`
		Author string `yaml:"author"`
	} `yaml:"copyright"`
	Auth struct {
		Provider          string `yaml:"provider"`
		JWT               string `yaml:"jwt"`
		Auth0Domain       string `yaml:"auth0_domain"`
		Auth0Audience     string `yaml:"auth0_audience"`
		Auth0ClientID     string `yaml:"auth0_client_id"`
		Auth0ClientSecret string `yaml:"auth0_client_secret"`
	} `yaml:"auth"`
	Server struct {
		Port     int    `yaml:"port"`
		Frontend struct {
			Dir string `yaml:"dir"`
			Api string `yaml:"api"`
		} `yaml:"frontend"`
	} `yaml:"server"`
	Database struct {
		Host           string `yaml:"host"`
		Port           int    `yaml:"port"`
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		DBName         string `yaml:"db"`
		TursoAuthToken string `yaml:"turso_auth_token,omitempty"`
		Sqlc           struct {
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
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db"`
	} `yaml:"cache"`
	S3 struct {
		URL    string `yaml:"url"`
		Access string `yaml:"access"`
		Secret string `yaml:"secret"`
	} `yaml:"s3"`
}

func gitTag() (string, error) {
	cmd := exec.Command("egg_cli", "version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func calcVersion() string {
	// First get the current git hash
	// then get the current git tag from the repository
	// then "v{git tag}+{git hash}"
	tag := "unknown"
	var err error
	tag, err = gitTag()
	if err != nil {
		tag = "unknown"
	}
	return tag

}

func parsePostgresURL(rawURL string) (string, string, string, int, string, error) {

	portint := -1
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", portint, "", err
	}

	// Extract username and password
	username := u.User.Username()
	password, _ := u.User.Password()

	// Split host and port
	hostPort := strings.Split(u.Host, ":")
	host := hostPort[0]
	port := ""
	if len(hostPort) > 1 {
		port = hostPort[1]
	}

	// Extract database name (remove leading '/')
	db := strings.TrimPrefix(u.Path, "/")
	portint, err = strconv.Atoi(port)
	if err != nil {
		return "", "", "", -1, "", err
	}

	return username, password, host, portint, db, nil
}

func parseCacheUrl(rawURL string) (string, string, string, int, string, error) {
	portint := -1
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", portint, "", err
	}
	// username can be empty
	username := u.User.Username()
	password, _ := u.User.Password()
	hostPort := strings.Split(u.Host, ":")
	host := hostPort[0]
	port := ""
	if len(hostPort) > 1 {
		port = hostPort[1]
	}
	portint, err = strconv.Atoi(port)
	if err != nil {
		return "", "", "", -1, "", err
	}
	db := strings.TrimPrefix(u.Path, "/")
	return username, password, host, portint, db, nil
}

func parseS3Url(url string) (string, int, error) {
	// this will inputted with host:port NO PROTOCOL
	portint := -1
	var err error
	var splitted = strings.Split(url, ":")
	host := splitted[0]
	port := splitted[1]
	portint, err = strconv.Atoi(port)
	if err != nil {
		return "", -1, err
	}
	return host, portint, nil
}

func (conf *GeneratorConfiguration) PostgresUrl(url string) {
	var err error
	conf.Database.Host,
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Port,
		conf.Database.DBName,
		err = parsePostgresURL(url)
	if err != nil {
		conf.Database.Host = "localhost"
		conf.Database.Port = 5432
		conf.Database.DBName = "postgres"
		conf.Database.Username = "guy"
		conf.Database.Password = "secret"
		fmt.Println(styles.EggProgressError.Render("Failed to parse postgres url: " + err.Error()))
		fmt.Println(styles.EggProgressInfo.Render("Using default postgres url"))
	}
}

func (conf *GeneratorConfiguration) CacheUrl(url string) {
	var err error
	conf.Cache.Host,
		conf.Cache.Username,
		conf.Cache.Password,
		conf.Cache.Port,
		conf.Cache.DBName,
		err = parseCacheUrl(url)
	if err != nil {
		conf.Cache.Host = "localhost"
		conf.Cache.Port = 6379
		conf.Cache.DBName = "postgres"
		conf.Cache.Username = "guy"
		conf.Cache.Password = "secret"
		fmt.Println(styles.EggProgressError.Render("Failed to parse cache url: " + err.Error()))
		fmt.Println(styles.EggProgressInfo.Render("Using default cache url (redis://localhost:6379)"))
	}
}

func (conf *GeneratorConfiguration) Init(version string) *GeneratorConfiguration {
	conf.EggVersion = version
	return conf
}
