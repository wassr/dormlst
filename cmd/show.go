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
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/wassr/dormlst/internal/csvdb"
	"github.com/wassr/dormlst/internal/ui"
)

var showCmd = &cobra.Command{
	Use:   "show [query]",
	Short: "See all available data on one person",
	Long:  `Search for a resident and display their full profile. If multiple matches are found, you will be prompted to select one.`,
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

		res, _, auto, err := ui.SelectFromMatches(residents, query, "Select Resident to Show")
		if err != nil {
			if errors.Is(err, ui.ErrNotFound) || errors.Is(err, ui.ErrAborted) {
				if errors.Is(err, ui.ErrNotFound) {
					fmt.Println(err.Error())
				}
				return nil
			}
			return err
		}

		status := "Active"
		if !res.Active {
			status = "Inactive"
		}

		// Only show leading newlines if we had to pick from multiple matches
		if !auto {
			fmt.Println()
		}

		fmt.Printf("  %s %s (Room %d)\n", res.FirstName, res.LastName, res.RoomNumber)
		fmt.Printf("  %s\n\n", strings.Repeat("-", 40))

		w := tabwriter.NewWriter(os.Stdout, 4, 0, 2, ' ', 0)

		fmt.Fprintf(w, "  STATUS\t%s\n", status)
		fmt.Fprintln(w, "\t")

		fmt.Fprintln(w, "  CONTACT\t")
		fmt.Fprintf(w, "  Email\t%s\n", res.Email)
		fmt.Fprintf(w, "  Phone\t%s\n", res.PhoneNumber)
		fmt.Fprintln(w, "\t")

		fmt.Fprintln(w, "  DETAILS\t")
		fmt.Fprintf(w, "  Room\t%d\n", res.RoomNumber)
		fmt.Fprintf(w, "  Birthday\t%s (Age: %d)\n", res.Birthday.Format("2006-01-02"), calculateAge(res.Birthday))
		fmt.Fprintln(w, "\t")

		fmt.Fprintln(w, "  TIMESTAMPS\t")
		fmt.Fprintf(w, "  Signed Up\t%s\n", res.DateSignedUp.Format("2006-01-02"))
		fmt.Fprintf(w, "  Added\t%s\n", res.DateAdded.Format("2006-01-02"))
		fmt.Fprintf(w, "  Modified\t%s\n", res.DateModified.Format("2006-01-02"))

		w.Flush()
		fmt.Println()

		return nil
	},
}

func calculateAge(birthdate time.Time) int {
	now := time.Now()
	years := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		years--
	}
	return years
}

func init() {
	rootCmd.AddCommand(showCmd)
}
