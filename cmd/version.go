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

	"github.com/spf13/cobra"
)

// Version information variables set by main package
var (
	appVersion   = "unknown"
	appBuildTime = "unknown"
	appGitCommit = "unknown"
)

// SetVersionInfo sets the version information from main package
func SetVersionInfo(version, buildTime, gitCommit string) {
	appVersion = version
	appBuildTime = buildTime
	appGitCommit = gitCommit
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of egg_cli",
	Long: `Print the version number of egg_cli as described by the ldflags
during compilation.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		hash, _ := cmd.Flags().GetBool("hash")
		buildTime, _ := cmd.Flags().GetBool("build-time")
		oneline, _ := cmd.Flags().GetBool("oneline")
		compact, _ := cmd.Flags().GetBool("compact")

		if verbose {
			hash = true
			buildTime = true
		}

		if compact {
			fmt.Printf("%s+%s+%s\n", appVersion, appGitCommit, appBuildTime)
			return
		}

		if oneline {
			fmt.Printf("egg_cli %s (%s) built %s\n", appVersion, appGitCommit, appBuildTime)
			return
		}

		fmt.Printf("egg_cli %s\n", appVersion)

		if hash {
			fmt.Printf("Commit: %s\n", appGitCommit)
		}
		if buildTime {
			fmt.Printf("Built: %s\n", appBuildTime)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	versionCmd.Flags().BoolP("hash", "#", false, "Print the commit hash from ldflags")
	versionCmd.Flags().BoolP("build-time", "t", false, "Print the build time from ldflags")
	versionCmd.Flags().BoolP("verbose", "v", false, "Print the version number, commit hash and build time from ldflags. same as -t and -# together ")
	versionCmd.Flags().BoolP("oneline", "1", false, "Print all version information on a single line")
	versionCmd.Flags().BoolP("compact", "c", false, "Print version+hash+buildtime in compact format")
}
