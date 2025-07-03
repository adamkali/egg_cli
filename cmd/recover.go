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
	"os"

	"github.com/adamkali/egg_cli/models"
	"github.com/adamkali/egg_cli/pkg"
	"github.com/spf13/cobra"
)

var recoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "Initialize a project with the .scrambled file",
	Long:  `Try to recover a project from the .scrambled file, if it exists`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := models.NewLogger("egg-log")
		if err != nil {
			os.Exit(1)
		}
		if err := pkg.RecoverFromScrambled(logger); err != nil {
			os.Exit(2)

		}
	},
}

func init() {
	rootCmd.AddCommand(recoverCmd)
}
