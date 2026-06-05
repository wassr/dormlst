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
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/wassr/dormlst/internal/csvdb"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get stats and health check on the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ANSI color codes
		bold := "\033[1m"
		reset := "\033[0m"
		green := "\033[32m"
		red := "\033[31m"
		faint := "\033[2m"

		fmt.Printf("%sSystem Status%s\n", bold, reset)
		fmt.Printf("%s\n", faint+strings.Repeat("\u2500", 45)+reset)

		// 1. Config Check
		fmt.Printf("Config file:   %s [%sOK%s]\n", cfgFile, green, reset)

		// 2. Database Check
		residents, err := csvdb.Load(cfg.Database.Path)
		if err != nil {
			fmt.Printf("Database:      %s [%sERROR%s: %v]\n", cfg.Database.Path, red, reset, err)
			return nil
		}
		fmt.Printf("Database:      %s [%sOK%s]\n", cfg.Database.Path, green, reset)
		fmt.Println()

		// 3. Stats Calculation
		total := len(residents)
		active := 0
		inactive := 0
		var invalidMsgs []string
		uniqueRooms := make(map[int]bool)
		var totalAge int
		recentMods := 0
		now := time.Now()
		sevenDaysAgo := now.AddDate(0, 0, -7)

		for _, r := range residents {
			if r.Active {
				active++
			} else {
				inactive++
			}

			if err := r.Validate(); err != nil {
				invalidMsgs = append(invalidMsgs, fmt.Sprintf("%s %s (Room %d): %v", r.FirstName, r.LastName, r.RoomNumber, err))
			}

			uniqueRooms[r.RoomNumber] = true
			totalAge += calculateAge(r.Birthday)

			if r.DateModified.After(sevenDaysAgo) {
				recentMods++
			}
		}

		avgAge := 0.0
		if total > 0 {
			avgAge = float64(totalAge) / float64(total)
		}

		// 4. Output Stats
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "%s[ Residents ]%s\t\n", bold, reset)
		fmt.Fprintf(w, "Total\t%d\n", total)
		fmt.Fprintf(w, "Active\t%s%d%s\n", green, active, reset)
		fmt.Fprintf(w, "Inactive\t%s%d%s\n", red, inactive, reset)
		fmt.Fprintf(w, "Avg. Age\t%.1f years\n", avgAge)
		fmt.Fprintln(w, "\t")

		fmt.Fprintf(w, "%s[ Infrastructure ]%s\t\n", bold, reset)
		fmt.Fprintf(w, "Unique Rooms\t%d\n", len(uniqueRooms))
		fmt.Fprintln(w, "\t")

		fmt.Fprintf(w, "%s[ Health ]%s\t\n", bold, reset)
		healthColor := green
		if len(invalidMsgs) > 0 {
			healthColor = red
		}
		fmt.Fprintf(w, "Valid Entries\t%s%d/%d%s\n", healthColor, total-len(invalidMsgs), total, reset)
		fmt.Fprintf(w, "Recent Changes (7d)\t%d\n", recentMods)
		if err := w.Flush(); err != nil {
			return err
		}

		// 5. Report Invalid Entries
		if len(invalidMsgs) > 0 {
			fmt.Printf("\n%sInvalid Entries Found:%s\n", red+bold, reset)
			for _, msg := range invalidMsgs {
				fmt.Printf("%s- %s%s\n", red, msg, reset)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
