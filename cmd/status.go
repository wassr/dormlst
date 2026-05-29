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
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/wassr/dormlst/internal/csvdb"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get stats and health check on the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("\n  System Status")
		fmt.Println("  -------------")

		// 1. Config Check
		fmt.Printf("  Config file:   %s [OK]\n", cfgFile)

		// 2. Database Check
		residents, err := csvdb.Load(cfg.Database.Path)
		if err != nil {
			fmt.Printf("  Database:      %s [ERROR: %v]\n", cfg.Database.Path, err)
			return nil
		}
		fmt.Printf("  Database:      %s [OK]\n", cfg.Database.Path)
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
		w := tabwriter.NewWriter(os.Stdout, 4, 0, 2, ' ', 0)
		fmt.Fprintf(w, "  RESIDENTS\t\n")
		fmt.Fprintf(w, "  Total\t%d\n", total)
		fmt.Fprintf(w, "  Active\t%d\n", active)
		fmt.Fprintf(w, "  Inactive\t%d\n", inactive)
		fmt.Fprintf(w, "  Avg. Age\t%.1f\n", avgAge)
		fmt.Fprintf(w, "\t\n")
		fmt.Fprintf(w, "  INFRASTRUCTURE\t\n")
		fmt.Fprintf(w, "  Unique Rooms\t%d\n", len(uniqueRooms))
		fmt.Fprintf(w, "\t\n")
		fmt.Fprintf(w, "  HEALTH\t\n")
		fmt.Fprintf(w, "  Valid Entries\t%d/%d\n", total-len(invalidMsgs), total)
		fmt.Fprintf(w, "  Recent Changes (7d)\t%d\n", recentMods)
		w.Flush()

		// 5. Report Invalid Entries
		if len(invalidMsgs) > 0 {
			fmt.Println("\n  Invalid Entries Found:")
			for _, msg := range invalidMsgs {
				fmt.Printf("  - %s\n", msg)
			}
		}

		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
