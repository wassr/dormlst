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
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/wassr/dormlst/internal/csvdb"
	"github.com/wassr/dormlst/internal/model"
)

var (
	searchSort     string
	searchActive   bool
	searchInactive bool
)

var searchCmd = &cobra.Command{
	Use:     "search [query]",
	Aliases: []string{"se"},
	Short:   "Find residents based on specific criteria",
	Long:  `Displays a tabular list of residents matching the query in name, room, or email. If no query is provided, all residents are shown.`,
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
		if searchActive {
			filter = model.FilterActive
		} else if searchInactive {
			filter = model.FilterInactive
		}

		matched := model.FilterResidents(residents, query, filter)

		if len(matched) == 0 {
			if query != "" {
				fmt.Printf("No residents found matching: %s\n", query)
			} else {
				fmt.Println("No residents found.")
			}
			return nil
		}

		// Sorting logic
		sort.Slice(matched, func(i, j int) bool {
			switch searchSort {
			case "name":
				if matched[i].FirstName != matched[j].FirstName {
					return matched[i].FirstName < matched[j].FirstName
				}
				return matched[i].LastName < matched[j].LastName
			case "added":
				return matched[i].DateAdded.Before(matched[j].DateAdded)
			case "signup":
				return matched[i].DateSignedUp.Before(matched[j].DateSignedUp)
			case "modified":
				return matched[i].DateModified.Before(matched[j].DateModified)
			case "room":
				fallthrough
			default:
				if matched[i].RoomNumber != matched[j].RoomNumber {
					return matched[i].RoomNumber < matched[j].RoomNumber
				}
				return matched[i].FirstName < matched[j].FirstName
			}
		})

		bold := "\033[1m"
		reset := "\033[0m"
		green := "\033[32m"
		red := "\033[31m"

		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ROOM\tNAME\tEMAIL\tPHONE\tACTIVE")

		for _, res := range matched {
			activeStr := "\u2713"
			if !res.Active {
				activeStr = "\u2717"
			}
			fmt.Fprintf(w, "%d\t%s %s\t%s\t%s\t%s\n",
				res.RoomNumber, res.FirstName, res.LastName, res.Email, res.PhoneNumber, activeStr)
		}
		w.Flush()

		// Apply colors line by line
		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		for i, line := range lines {
			if i == 0 {
				// Header
				fmt.Printf("%s%s%s\n", bold, line, reset)
				continue
			}

			// Data rows
			if strings.HasSuffix(line, "\u2713") {
				// Active
				fmt.Printf("%s%s%s\n", strings.TrimSuffix(line, "\u2713"), green, "\u2713"+reset)
			} else if strings.HasSuffix(line, "\u2717") {
				// Inactive
				fmt.Printf("%s%s%s\n", strings.TrimSuffix(line, "\u2717"), red, "\u2717"+reset)
			} else {
				fmt.Println(line)
			}
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().StringVarP(&searchSort, "sort", "s", "room", "Sort by: name, room, added, signup, modified")
	searchCmd.Flags().BoolVar(&searchActive, "active", false, "Filter for active residents only")
	searchCmd.Flags().BoolVar(&searchInactive, "inactive", false, "Filter for inactive residents only")
	rootCmd.AddCommand(searchCmd)
}
