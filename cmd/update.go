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
	updateInteractive bool
	updateFirstName   string
	updateLastName    string
	updateRoom        int
	updateEmail       string
	updatePhone       string
	updateBirthday    string
	updateSignedUp    string
	updateSetEnable   bool
	updateSetDisable  bool
	updateActive      bool
	updateInactive    bool
)

var updateCmd = &cobra.Command{
	Use:     "update [query]",
	Aliases: []string{"up"},
	Short:   "Modify an existing resident's information",
	Long:  `Search for a resident by name, room number, or email. If multiple matches are found, you will be prompted to select one. Use --interactive (-i) for a guided update, or use flags to update specific fields non-interactively.`,
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
					fmt.Printf("\033[31m\033[1m[ERROR]\033[0m %s\n", err.Error())
				}
				return nil
			}
			return err
		}

		if updateInteractive {
			res, err = ui.PromptResident(res, true)
			if err != nil {
				if errors.Is(err, ui.ErrAborted) {
					return nil
				}
				return err
			}
		} else {
			changed := false
			if cmd.Flags().Changed("first-name") {
				res.FirstName = updateFirstName
				changed = true
			}
			if cmd.Flags().Changed("last-name") {
				res.LastName = updateLastName
				changed = true
			}
			if cmd.Flags().Changed("room") {
				res.RoomNumber = updateRoom
				changed = true
			}
			if cmd.Flags().Changed("email") {
				res.Email = updateEmail
				changed = true
			}
			if cmd.Flags().Changed("phone") {
				res.PhoneNumber = updatePhone
				changed = true
			}
			if cmd.Flags().Changed("birthday") {
				t, err := time.Parse("2006-01-02", updateBirthday)
				if err != nil {
					return fmt.Errorf("invalid date format for birthday: %w", err)
				}
				res.Birthday = t
				changed = true
			}
			if cmd.Flags().Changed("signed-up") {
				t, err := time.Parse("2006-01-02", updateSignedUp)
				if err != nil {
					return fmt.Errorf("invalid date format for signed-up: %w", err)
				}
				res.DateSignedUp = t
				changed = true
			}
			if updateSetEnable {
				res.Active = true
				changed = true
			}
			if updateSetDisable {
				res.Active = false
				changed = true
			}

			if !changed {
				fmt.Printf("\033[31m\033[1m[ERROR]\033[0m No update flags provided. Use --interactive (-i) or specify fields to update.\n")
				return nil
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

		fmt.Printf("\033[32m\033[1m[OK]\033[0m Successfully updated resident: %s %s\n", res.FirstName, res.LastName)
		return nil
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&updateInteractive, "interactive", "i", false, "Start interactive update mode")
	updateCmd.Flags().StringVar(&updateFirstName, "first-name", "", "New first name")
	updateCmd.Flags().StringVar(&updateLastName, "last-name", "", "New last name")
	updateCmd.Flags().IntVar(&updateRoom, "room", 0, "New room number")
	updateCmd.Flags().StringVar(&updateEmail, "email", "", "New email address")
	updateCmd.Flags().StringVar(&updatePhone, "phone", "", "New phone number")
	updateCmd.Flags().StringVar(&updateBirthday, "birthday", "", "New birthday (YYYY-MM-DD)")
	updateCmd.Flags().StringVar(&updateSignedUp, "signed-up", "", "New date signed up (YYYY-MM-DD)")
	updateCmd.Flags().BoolVar(&updateSetEnable, "enable", false, "Set resident to active")
	updateCmd.Flags().BoolVar(&updateSetDisable, "disable", false, "Set resident to inactive")

	updateCmd.Flags().BoolVar(&updateActive, "filter-active", false, "Only search for active residents")
	updateCmd.Flags().BoolVar(&updateInactive, "filter-inactive", false, "Only search for inactive residents")

	rootCmd.AddCommand(updateCmd)
}
