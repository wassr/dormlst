package model

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
	"testing"
	"time"
)

func TestResident_Validate(t *testing.T) {
	validDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		resident      Resident
		wantErr       bool
		expectedPhone string
	}{
		{
			name: "Valid Resident with AT local phone",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "max@example.com",
				PhoneNumber:  "0677 62331692",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
				Active:       true,
			},
			wantErr:       false,
			expectedPhone: "+4367762331692",
		},
		{
			name: "Valid Resident with international phone",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "max@example.com",
				PhoneNumber:  "+43 677 62331692",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
				Active:       true,
			},
			wantErr:       false,
			expectedPhone: "+4367762331692",
		},
		{
			name: "Valid Resident with complex email",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "max.m-ustermann+tag@student.tuwien.ac.at",
				PhoneNumber:  "+43 677 62331692",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
				Active:       true,
			},
			wantErr:       false,
			expectedPhone: "+4367762331692",
		},
		{
			name: "Missing FirstName",
			resident: Resident{
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "max@example.com",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
			},
			wantErr: true,
		},
		{
			name: "Missing LastName",
			resident: Resident{
				FirstName:    "Max",
				RoomNumber:   101,
				Email:        "max@example.com",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
			},
			wantErr: true,
		},
		{
			name: "Negative Room",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   -5,
				Email:        "max@example.com",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
			},
			wantErr: true,
		},
		{
			name: "Zero Room",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   0,
				Email:        "max@example.com",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
			},
			wantErr: true,
		},
		{
			name: "Invalid Email",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "invalid-email",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
			},
			wantErr: true,
		},
		{
			name: "Missing Email",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
			},
			wantErr: true,
		},
		{
			name: "Invalid Phone",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "max@example.com",
				PhoneNumber:  "12345", // too short / invalid
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
			},
			wantErr: true,
		},
		{
			name: "Missing Birthday",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "max@example.com",
				DateSignedUp: validDate,
				DateAdded:    validDate,
				DateModified: validDate,
			},
			wantErr: true,
		},
		{
			name: "Missing DateSignedUp",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "max@example.com",
				Birthday:     validDate,
				DateAdded:    validDate,
				DateModified: validDate,
			},
			wantErr: true,
		},
		{
			name: "Missing DateAdded",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "max@example.com",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateModified: validDate,
			},
			wantErr: true,
		},
		{
			name: "Missing DateModified",
			resident: Resident{
				FirstName:    "Max",
				LastName:     "Mustermann",
				RoomNumber:   101,
				Email:        "max@example.com",
				Birthday:     validDate,
				DateSignedUp: validDate,
				DateAdded:    validDate,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.resident
			err := r.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Resident.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && r.PhoneNumber != tt.expectedPhone {
				t.Errorf("Expected phone %s, got %s", tt.expectedPhone, r.PhoneNumber)
			}
		})
	}
}

func TestFilterResidents(t *testing.T) {
	residents := []Resident{
		{FirstName: "Alice", LastName: "Alpha", RoomNumber: 101, Email: "alice@example.com", Active: true},
		{FirstName: "Bob", LastName: "Beta", RoomNumber: 202, Email: "bob@example.com", Active: false},
	}

	tests := []struct {
		name     string
		query    string
		filter   ActiveFilter
		expected int
	}{
		{"Empty query, all", "", FilterAll, 2},
		{"Empty query, active", "", FilterActive, 1},
		{"Empty query, inactive", "", FilterInactive, 1},
		{"By first name, all", "alice", FilterAll, 1},
		{"By first name, inactive", "alice", FilterInactive, 0},
		{"By room, active", "202", FilterActive, 0},
		{"By room, inactive", "202", FilterInactive, 1},
		{"By email, all", "example.com", FilterAll, 2},
		{"No match", "gamma", FilterAll, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterResidents(residents, tt.query, tt.filter)
			if len(got) != tt.expected {
				t.Errorf("FilterResidents() returned %d results, expected %d", len(got), tt.expected)
			}
		})
	}
}
