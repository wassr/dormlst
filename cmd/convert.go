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
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Transform the CSV file into a Microsoft Excel (.xlsx) format",
	RunE: func(cmd *cobra.Command, args []string) error {
		residents, err := csvdb.Load(cfg.Database.Path)
		if err != nil {
			return fmt.Errorf("failed to load database: %w", err)
		}

		if len(residents) == 0 {
			fmt.Println("No residents found in database. Excel file will be empty.")
		}

		if err := excel.Generate(residents, cfg.Output); err != nil {
			return fmt.Errorf("failed to generate excel file: %w", err)
		}

		fmt.Printf("Successfully converted %s to %s\n", cfg.Database.Path, cfg.Output.Path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
