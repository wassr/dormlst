package csvdb

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
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/wassr/dormlst/internal/model"
)

const dateFormat = "2006-01-02"

// Load reads residents from a CSV file.
func Load(path string) ([]model.Resident, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Resident{}, nil
		}
		return nil, err
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	// Read header
	header, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return []model.Resident{}, nil
		}
		return nil, err
	}

	var residents []model.Resident
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		res, err := parseRecord(record, header)
		if err != nil {
			return nil, fmt.Errorf("error parsing record: %w", err)
		}
		residents = append(residents, res)
	}

	return residents, nil
}

// Save writes residents to a CSV file in a git-friendly (sorted) way.
func Save(path string, residents []model.Resident) error {
	// Sort by RoomNumber, then LastName for deterministic output
	sort.Slice(residents, func(i, j int) bool {
		if residents[i].RoomNumber != residents[j].RoomNumber {
			return residents[i].RoomNumber < residents[j].RoomNumber
		}
		return residents[i].LastName < residents[j].LastName
	})

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"first_name", "last_name", "room_number", "phone_number",
		"email", "birthday", "date_signed_up", "date_added", "date_modified", "active",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, res := range residents {
		record := []string{
			res.FirstName,
			res.LastName,
			strconv.Itoa(res.RoomNumber),
			res.PhoneNumber,
			res.Email,
			res.Birthday.Format(dateFormat),
			res.DateSignedUp.Format(dateFormat),
			res.DateAdded.Format(dateFormat),
			res.DateModified.Format(dateFormat),
			strconv.FormatBool(res.Active),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func parseRecord(record []string, header []string) (model.Resident, error) {
	if len(record) < 10 {
		return model.Resident{}, fmt.Errorf("too few columns: expected 10, got %d", len(record))
	}

	room, _ := strconv.Atoi(record[2])
	birthday, _ := time.Parse(dateFormat, record[5])
	signedUp, _ := time.Parse(dateFormat, record[6])
	added, _ := time.Parse(dateFormat, record[7])
	modified, _ := time.Parse(dateFormat, record[8])
	active, _ := strconv.ParseBool(record[9])

	return model.Resident{
		FirstName:    record[0],
		LastName:     record[1],
		RoomNumber:   room,
		PhoneNumber:  record[3],
		Email:        record[4],
		Birthday:     birthday,
		DateSignedUp: signedUp,
		DateAdded:    added,
		DateModified: modified,
		Active:       active,
	}, nil
}
