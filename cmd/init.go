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

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a default configuration file",
	Long:  `Creates a default dormlst.yaml file in the current directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(cfgFile); err == nil {
			return fmt.Errorf("config file %s already exists", cfgFile)
		}

		defaultCfg := config.DefaultConfig()
		if err := defaultCfg.Save(cfgFile); err != nil {
			return fmt.Errorf("failed to save default config: %w", err)
		}

		fmt.Printf("\033[32m\033[1m[OK]\033[0m Initialized default configuration in %s\n", cfgFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
