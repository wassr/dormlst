package cmd

/*
Copyright © 2026 Noel Atzwanger (@wassr) <me@wassr.cc>

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

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wassr/dormlst/internal/config"
)

var (
	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "dormlst",
	Version: "1.0.0",
	Short:   "Manage resident lists in Git-friendly CSV and generate .xlsx for uploads",
	Long: `dormlst is a CLI tool designed to bypass the sync and formatting issues of Excel. 
It maintains your resident database in a simple, sorted CSV format for version control, 
while providing tools to generate the .xlsx files required by external services 
(like workitout.at).`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "dormlst.yaml", "config file (default is dormlst.yaml)")
}

func initConfig() error {
	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		if os.IsNotExist(err) {
			// If config doesn't exist, use default and optionally warn
			cfg = config.DefaultConfig()
			return nil
		}
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}
