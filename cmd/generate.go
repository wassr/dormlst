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

	"github.com/spf13/cobra"
	"github.com/wassr/dormlst/internal/csvdb"
	"github.com/wassr/dormlst/internal/excel"
	"github.com/wassr/dormlst/internal/model"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate the .xlsx file required for uploads",
	Long:  `Reads the CSV database and produces an Excel file based on the YAML configuration. This is intended to be used as a compatibility artifact for external systems.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		residents, err := csvdb.Load(cfg.Database.Path)
		if err != nil {
			return fmt.Errorf("failed to load database: %w", err)
		}

		residents = model.FilterResidents(residents, "", model.FilterActive)

		if len(residents) == 0 {
			fmt.Printf("\033[2m\033[1m[INFO]\033[0m No residents found in database. Output file will be empty.\n")
		}

		if err := excel.Generate(residents, cfg.Output); err != nil {
			return fmt.Errorf("failed to generate output file: %w", err)
		}

		fmt.Printf("\033[32m\033[1m[OK]\033[0m Successfully generated %s from %s\n", cfg.Output.Path, cfg.Database.Path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
