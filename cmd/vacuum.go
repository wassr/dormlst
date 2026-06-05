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
	"github.com/wassr/dormlst/internal/model"
	"github.com/wassr/dormlst/internal/ui"
)

var (
	vacuumYes bool
)

var vacuumCmd = &cobra.Command{
	Use:     "vacuum",
	Aliases: []string{"vac"},
	Short:   "Permanently remove all disabled (inactive) resident entries",
	Long:    `Identifies all residents marked as inactive and removes them from the CSV database. This is a destructive operation and will prompt for confirmation unless --yes is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		residents, err := csvdb.Load(cfg.Database.Path)
		if err != nil {
			return err
		}

		var activeOnly []model.Resident
		inactiveCount := 0
		for _, r := range residents {
			if r.Active {
				activeOnly = append(activeOnly, r)
			} else {
				inactiveCount++
			}
		}

		if inactiveCount == 0 {
			fmt.Printf("\033[2m\033[1m[INFO]\033[0m No inactive residents found. Nothing to vacuum.\n")
			return nil
		}

		if !vacuumYes {
			ok, err := ui.Confirm(fmt.Sprintf("Are you sure you want to permanently remove %d inactive resident(s)?", inactiveCount))
			if err != nil {
				return err
			}
			if !ok {
				fmt.Println("\033[31m\033[1m[ABORT]\033[0m Aborted.")
				return nil
			}
		}

		if err := csvdb.Save(cfg.Database.Path, activeOnly); err != nil {
			return err
		}

		fmt.Printf("\033[32m\033[1m[OK]\033[0m Successfully removed %d inactive resident(s).\n", inactiveCount)
		return nil
	},
}

func init() {
	vacuumCmd.Flags().BoolVarP(&vacuumYes, "yes", "y", false, "Skip confirmation prompt")
	rootCmd.AddCommand(vacuumCmd)
}
