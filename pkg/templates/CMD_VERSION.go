package templates
const VersionCmdTemplate = `
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
package cmd

import (
	"fmt"

	"{{.Namespace}}/cmd/configuration"
	"github.com/labstack/echo"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Gets the semantic version of the server ",
	Long: ` + "`" + `
*** Help Text
Gets the semantic version of the server.
Based on the environment, passed you could have a different version 
of what is a production or a nightly based on ci/cd.
this can also be used in ci/cd piplines to if you want to 
use the versioning to determine deployments or using docker to tag
to version docker images.

*** Command 
**** Default 
--- bash
go build main.go -o {{.Name }}
./{{.Name }} version
---

**** with -e passed
If one had some configuration file really-sick-config.yaml
--- bash
go build main.go -o egg_app
./egg_app version -e really-sick-config
---

**** using to make a docker version
--- bash 
go build -o egg_app
EGG_APP_VER=echo(./egg_app version -e really-sick-config)
docker tag repository/user/egg_app:EGG_APP_VER
---
    ` + "`" + `,
	Run: func(cmd *cobra.Command, args []string) {
        e := echo.New()
        config, err := configuration.LoadConfiguration(Environment)
        if err != nil {
            e.Logger.Fatal(err.Error())
            panic(err.Error())
        }
        fmt.Println(fmt.Sprintf("%s", config.Semver))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
`
