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
	"time"

	"github.com/spf13/cobra"
	"github.com/wassr/dormlst/internal/csvdb"
	"github.com/wassr/dormlst/internal/model"
	"github.com/wassr/dormlst/internal/ui"
)

var (
	addInteractive bool
	addFirstName   string
	addLastName    string
	addRoom        int
	addEmail       string
	addPhone       string
	addBirthday    string
	addSignedUp    string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new resident entry in the CSV list",
	RunE: func(cmd *cobra.Command, args []string) error {
		var res model.Resident
		var err error

		if addInteractive {
			res, err = ui.PromptResident(model.Resident{}, false)
			if err != nil {
				if err == ui.ErrAborted {
					return nil
				}
				return err
			}
		} else {
			birthday, _ := time.Parse("2006-01-02", addBirthday)
			signedUp, _ := time.Parse("2006-01-02", addSignedUp)
			if addSignedUp == "" {
				signedUp = time.Now()
			}

			now := time.Now()
			res = model.Resident{
				FirstName:    addFirstName,
				LastName:     addLastName,
				RoomNumber:   addRoom,
				Email:        addEmail,
				PhoneNumber:  addPhone,
				Birthday:     birthday,
				DateSignedUp: signedUp,
				DateAdded:    now,
				DateModified: now,
				Active:       true,
			}
		}

		if err := res.Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		residents, err := csvdb.Load(cfg.Database.Path)
		if err != nil {
			return err
		}

		residents = append(residents, res)
		if err := csvdb.Save(cfg.Database.Path, residents); err != nil {
			return err
		}

		fmt.Printf("Successfully added resident: %s %s (Room %d)\n", res.FirstName, res.LastName, res.RoomNumber)
		return nil
	},
}

func init() {
	addCmd.Flags().BoolVarP(&addInteractive, "interactive", "i", false, "Start interactive add mode")
	addCmd.Flags().StringVar(&addFirstName, "first-name", "", "First name of the resident")
	addCmd.Flags().StringVar(&addLastName, "last-name", "", "Last name of the resident")
	addCmd.Flags().IntVar(&addRoom, "room", 0, "Room number")
	addCmd.Flags().StringVar(&addEmail, "email", "", "Email address")
	addCmd.Flags().StringVar(&addPhone, "phone", "", "Phone number")
	addCmd.Flags().StringVar(&addBirthday, "birthday", "", "Birthday (YYYY-MM-DD)")
	addCmd.Flags().StringVar(&addSignedUp, "signed-up", "", "Date signed up (YYYY-MM-DD, default: today)")

	rootCmd.AddCommand(addCmd)
}
