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

var disableCmd = &cobra.Command{
	Use:     "disable [query]",
	Aliases: []string{"dis"},
	Short:   "Set a resident's status to inactive without deleting",
	Long:  `Search for a resident by name, room number, or email. If multiple matches are found, you will be prompted to select one to disable.`,
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

		res, index, _, err := ui.SelectFromMatches(residents, query, model.FilterActive, "Select Resident to Disable")
		if err != nil {
			if errors.Is(err, ui.ErrNotFound) || errors.Is(err, ui.ErrAborted) {
				if errors.Is(err, ui.ErrNotFound) {
					fmt.Printf("\033[31m\033[1m[ERROR]\033[0m %s\n", err.Error())
				}
				return nil
			}
			return err
		}

		residents[index].Active = false
		residents[index].DateModified = time.Now()

		if err := csvdb.Save(cfg.Database.Path, residents); err != nil {
			return err
		}

		fmt.Printf("\033[32m\033[1m[OK]\033[0m Successfully disabled resident: %s %s\n", res.FirstName, res.LastName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)
}
