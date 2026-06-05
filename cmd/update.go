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
	"time"

	"github.com/spf13/cobra"
	"github.com/wassr/dormlst/internal/csvdb"
	"github.com/wassr/dormlst/internal/model"
	"github.com/wassr/dormlst/internal/ui"
)

var (
	updateRoom     int
	updateEmail    string
	updateSignedUp string
	updateActive   bool
	updateInactive bool
)

var updateCmd = &cobra.Command{
	Use:     "update [query]",
	Aliases: []string{"up"},
	Short:   "Modify an existing resident's information",
	Long:  `Search for a resident by name, room number, or email. If multiple matches are found, you will be prompted to select one.`,
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
		if updateActive {
			filter = model.FilterActive
		} else if updateInactive {
			filter = model.FilterInactive
		}

		res, index, _, err := ui.SelectFromMatches(residents, query, filter, "Select Resident to Update")
		if err != nil {
			if errors.Is(err, ui.ErrNotFound) || errors.Is(err, ui.ErrAborted) {
				if errors.Is(err, ui.ErrNotFound) {
					fmt.Println(err.Error())
				}
				return nil
			}
			return err
		}

		if updateRoom != 0 || updateEmail != "" || updateSignedUp != "" {
			if updateRoom != 0 {
				res.RoomNumber = updateRoom
			}
			if updateEmail != "" {
				res.Email = updateEmail
			}
			if updateSignedUp != "" {
				t, err := time.Parse("2006-01-02", updateSignedUp)
				if err != nil {
					return fmt.Errorf("invalid date format for signed-up: %w", err)
				}
				res.DateSignedUp = t
			}
		} else {
			res, err = ui.PromptResident(res, true)
			if err != nil {
				if errors.Is(err, ui.ErrAborted) {
					return nil
				}
				return err
			}
		}

		res.DateModified = time.Now()
		residents[index] = res
		if err := residents[index].Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		if err := csvdb.Save(cfg.Database.Path, residents); err != nil {
			return err
		}

		fmt.Printf("Successfully updated resident %s %s\n", res.FirstName, res.LastName)
		return nil
	},
}

func init() {
	updateCmd.Flags().IntVar(&updateRoom, "room", 0, "New room number")
	updateCmd.Flags().StringVar(&updateEmail, "email", "", "New email address")
	updateCmd.Flags().StringVar(&updateSignedUp, "signed-up", "", "New date signed up (YYYY-MM-DD)")
	updateCmd.Flags().BoolVar(&updateActive, "active", false, "Filter for active residents only")
	updateCmd.Flags().BoolVar(&updateInactive, "inactive", false, "Filter for inactive residents only")

	rootCmd.AddCommand(updateCmd)
}
