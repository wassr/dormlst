package excel

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
	"strconv"

	"github.com/wassr/dormlst/internal/config"
	"github.com/wassr/dormlst/internal/model"
	"github.com/xuri/excelize/v2"
)

// Generate creates an Excel file based on the provided residents and configuration.
func Generate(residents []model.Resident, cfg config.OutputConfig) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := cfg.SheetName
	if sheetName == "" {
		sheetName = "Sheet1"
	}

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	if sheetName != "Sheet1" {
		if err := f.DeleteSheet("Sheet1"); err != nil {
			return err
		}
	}

	// Write Header
	for i, col := range cfg.Columns {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return err
		}
		if err := f.SetCellValue(sheetName, cell, col.Target); err != nil {
			return err
		}
	}

	// Write Data
	dateFormat := cfg.DateFormat
	if dateFormat == "" {
		dateFormat = "02.01.2006" // Default to German format
	}

	for rowIndex, res := range residents {
		for colIndex, col := range cfg.Columns {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			if err != nil {
				return err
			}
			val := getFieldValue(res, col.Source, dateFormat)
			if err := f.SetCellValue(sheetName, cell, val); err != nil {
				return err
			}
		}
	}

	if err := f.SaveAs(cfg.Path); err != nil {
		return err
	}

	return nil
}

func getFieldValue(res model.Resident, source string, dateFormat string) interface{} {
	switch source {
	case "first_name":
		return res.FirstName
	case "last_name":
		return res.LastName
	case "email":
		return res.Email
	case "phone":
		return res.PhoneNumber
	case "room":
		return res.RoomNumber
	case "birthday":
		return res.Birthday.Format(dateFormat)
	case "entry_date":
		return res.DateSignedUp.Format(dateFormat)
	case "date_added":
		return res.DateAdded.Format(dateFormat)
	case "date_modified":
		return res.DateModified.Format(dateFormat)
	case "active":
		return strconv.FormatBool(res.Active)
	default:
		return ""
	}
}
