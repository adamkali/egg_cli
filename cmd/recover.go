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
	"os"

	"github.com/adamkali/egg_cli/pkg"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/spf13/cobra"
)

var recoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "Recover a project from the .scrambled file",
	Long:  `Attempt to recover a project from the .scrambled file, if it exists. This will resume execution from where the previous run failed.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if .scrambled file exists before attempting recovery
		if !pkg.CheckScrambled() {
			fmt.Println("❌ No .scrambled file found. Nothing to recover.")
			fmt.Println("💡 Run 'egg init' first to create a project.")
			os.Exit(1)
		}

		// Initialize logger
		logger, err := models.NewLogger("egg-log")
		if err != nil {
			fmt.Printf("❌ Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Close()

		fmt.Println("🔄 Attempting to recover project from .scrambled file...")

		// Attempt recovery
		if err := pkg.RecoverFromScrambled(logger); err != nil {
			fmt.Printf("❌ Recovery failed: %v\n", err)
			fmt.Println("💡 Check the .scrambled file for details about the failure.")
			logger.Error("Recovery failed: %v", err)
			os.Exit(2)
		}

		fmt.Println("✅ Project recovered successfully!")
		logger.Info("Project recovery completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(recoverCmd)
}
