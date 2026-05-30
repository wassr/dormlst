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
	"os"
	"testing"
	"time"

	"github.com/wassr/dormlst/internal/model"
)

func TestSaveAndLoad(t *testing.T) {
	tmpFile := "test_residents.csv"
	defer func() { _ = os.Remove(tmpFile) }()

	validDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	residents := []model.Resident{
		{
			FirstName:    "Zelda",
			LastName:     "Wild",
			RoomNumber:   202,
			Email:        "zelda@example.com",
			Birthday:     validDate,
			DateSignedUp: validDate,
			DateAdded:    validDate,
			DateModified: validDate,
			Active:       true,
		},
		{
			FirstName:    "Link",
			LastName:     "Past",
			RoomNumber:   101,
			Email:        "link@example.com",
			Birthday:     validDate,
			DateSignedUp: validDate,
			DateAdded:    validDate,
			DateModified: validDate,
			Active:       true,
		},
	}

	err := Save(tmpFile, residents)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Fatalf("Expected 2 residents, got %d", len(loaded))
	}

	// Check sorting: Room 101 should be first
	if loaded[0].RoomNumber != 101 {
		t.Errorf("Expected first resident to be room 101, got %d", loaded[0].RoomNumber)
	}
	if loaded[1].RoomNumber != 202 {
		t.Errorf("Expected second resident to be room 202, got %d", loaded[1].RoomNumber)
	}

	if loaded[0].FirstName != "Link" {
		t.Errorf("Expected first resident to be Link, got %s", loaded[0].FirstName)
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	loaded, err := Load("does_not_exist_xyz123.csv")
	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got %v", err)
	}
	if len(loaded) != 0 {
		t.Fatalf("Expected 0 residents, got %d", len(loaded))
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	tmpFile := "empty_test.csv"
	if err := os.WriteFile(tmpFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write empty file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	loaded, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error for empty file, got %v", err)
	}
	if len(loaded) != 0 {
		t.Fatalf("Expected 0 residents, got %d", len(loaded))
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	tmpFile := "corrupt_test.csv"
	// Write a file with a header but a corrupted row (too few columns)
	corruptData := "first_name,last_name,room_number,phone_number,email,birthday,date_signed_up,date_added,date_modified,active\nMax,Mustermann"
	if err := os.WriteFile(tmpFile, []byte(corruptData), 0644); err != nil {
		t.Fatalf("Failed to write corrupt file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	_, err := Load(tmpFile)
	if err == nil {
		t.Fatal("Expected an error for corrupt file, got nil")
	}
}
