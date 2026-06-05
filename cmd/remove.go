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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wassr/dormlst/internal/csvdb"
	"github.com/wassr/dormlst/internal/model"
	"github.com/wassr/dormlst/internal/ui"
)

var (
	removeActive   bool
	removeInactive bool
)

var removeCmd = &cobra.Command{
	Use:   "remove [query]",
	Short: "Delete a resident from the CSV list",
	Long:  `Search for a resident by name, room number, or email. If multiple matches are found, you will be prompted to select one for removal.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		residents, err := csvdb.Load(cfg.Database.Path)
		if err != nil {
			return err
		}

		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		filter := model.FilterAll
		if removeActive {
			filter = model.FilterActive
		} else if removeInactive {
			filter = model.FilterInactive
		}

		res, index, auto, err := ui.SelectFromMatches(residents, query, filter, "Select Resident to Remove")
		if err != nil {
			if errors.Is(err, ui.ErrNotFound) || errors.Is(err, ui.ErrAborted) {
				if errors.Is(err, ui.ErrNotFound) {
					fmt.Println(err.Error())
				}
				return nil
			}
			return err
		}

		if auto {
			ok, err := ui.Confirm(fmt.Sprintf("Are you sure you want to remove %s %s (Room %d)?", res.FirstName, res.LastName, res.RoomNumber))
			if err != nil {
				return err
			}
			if !ok {
				fmt.Println("Aborted.")
				return nil
			}
		}

		// Remove the resident at the found index
		residents = append(residents[:index], residents[index+1:]...)

		if err := csvdb.Save(cfg.Database.Path, residents); err != nil {
			return err
		}

		fmt.Printf("Successfully removed resident: %s %s\n", res.FirstName, res.LastName)
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVar(&removeActive, "active", false, "Filter for active residents only")
	removeCmd.Flags().BoolVar(&removeInactive, "inactive", false, "Filter for inactive residents only")
	rootCmd.AddCommand(removeCmd)
}
