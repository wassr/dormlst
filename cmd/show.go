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
	"github.com/wassr/dormlst/internal/model"
	"github.com/wassr/dormlst/internal/ui"
)

var (
	showActive   bool
	showInactive bool
)

var showCmd = &cobra.Command{
	Use:     "show [query]",
	Aliases: []string{"info"},
	Short:   "See all available data on one person",
	Long:    `Search for a resident and display their full profile. If multiple matches are found, you will be prompted to select one.`,
	Args:    cobra.MaximumNArgs(1),
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
		if showActive {
			filter = model.FilterActive
		} else if showInactive {
			filter = model.FilterInactive
		}

		res, _, auto, err := ui.SelectFromMatches(residents, query, filter, "Select Resident to Show")
		if err != nil {
			if errors.Is(err, ui.ErrNotFound) || errors.Is(err, ui.ErrAborted) {
				if errors.Is(err, ui.ErrNotFound) {
					fmt.Printf("\033[31m\033[1m[ERROR]\033[0m %s\n", err.Error())
				}
				return nil
			}
			return err
		}

		// ANSI color codes
		bold := "\033[1m"
		reset := "\033[0m"
		green := "\033[32m"
		red := "\033[31m"
		faint := "\033[2m"

		statusStr := green + "Active \u2713" + reset
		if !res.Active {
			statusStr = red + "Inactive \u2717" + reset
		}

		// Only show leading newline if we had to pick from multiple matches
		if !auto {
			fmt.Println()
		}

		fmt.Printf("%s%s %s%s %s(%d)%s\n", bold, res.FirstName, res.LastName, reset, faint, res.RoomNumber, reset)
		fmt.Printf("%s\n\n", faint+strings.Repeat("\u2500", 45)+reset)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		if _, err := fmt.Fprintf(w, "Status\t%s\n", statusStr); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "\t"); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(w, "%s[ Contact ]%s\t\n", bold, reset); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "Email\t%s\n", res.Email); err != nil {
			return err
		}
		if res.PhoneNumber != "" {
			if _, err := fmt.Fprintf(w, "Phone\t%s\n", res.PhoneNumber); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, "Phone\t%s--%s\n", faint, reset); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w, "\t"); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(w, "%s[ Details ]%s\t\n", bold, reset); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "Birthday\t%s (%d years)\n", res.Birthday.Format("2006-01-02"), calculateAge(res.Birthday)); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "\t"); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(w, "%s[ Timestamps ]%s\t\n", bold, reset); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "Signed Up\t%s\n", res.DateSignedUp.Format("2006-01-02")); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "Added\t%s\n", res.DateAdded.Format("2006-01-02")); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "Modified\t%s\n", res.DateModified.Format("2006-01-02")); err != nil {
			return err
		}

		if err := w.Flush(); err != nil {
			return err
		}

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
	showCmd.Flags().BoolVar(&showActive, "active", false, "Filter for active residents only")
	showCmd.Flags().BoolVar(&showInactive, "inactive", false, "Filter for inactive residents only")
	rootCmd.AddCommand(showCmd)
}
