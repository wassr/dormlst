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
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"
)

var (
	// Basic email regex for "verified" format
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
)

// Resident represents a person living in the dormitory.
type Resident struct {
	FirstName    string    `csv:"first_name"`
	LastName     string    `csv:"last_name"`
	RoomNumber   int       `csv:"room_number"`
	PhoneNumber  string    `csv:"phone_number"`
	Email        string    `csv:"email"`
	Birthday     time.Time `csv:"birthday"`
	DateSignedUp time.Time `csv:"date_signed_up"`
	DateAdded    time.Time `csv:"date_added"`
	DateModified time.Time `csv:"date_modified"`
	Active       bool      `csv:"active"`
}

// ActiveFilter defines how to filter residents by their active status.
type ActiveFilter int

const (
	// FilterAll returns both active and inactive residents.
	FilterAll ActiveFilter = iota
	// FilterActive returns only active residents.
	FilterActive
	// FilterInactive returns only inactive residents.
	FilterInactive
)

// FilterResidents returns a slice of residents that match the query and active filter.
func FilterResidents(residents []Resident, query string, activeFilter ActiveFilter) []Resident {
	query = strings.ToLower(query)
	var matched []Resident
	for _, r := range residents {
		// Apply active filter
		if activeFilter == FilterActive && !r.Active {
			continue
		}
		if activeFilter == FilterInactive && r.Active {
			continue
		}

		// Apply search query
		if query == "" {
			matched = append(matched, r)
			continue
		}

		name := strings.ToLower(r.FirstName + " " + r.LastName)
		room := strconv.Itoa(r.RoomNumber)
		email := strings.ToLower(r.Email)
		if strings.Contains(name, query) || strings.Contains(room, query) || strings.Contains(email, query) {
			matched = append(matched, r)
		}
	}
	return matched
}

// NormalizePhoneNumber parses and formats the phone number to E.164.
func NormalizePhoneNumber(phone string) (string, error) {
	if phone == "" {
		return "", nil
	}
	num, err := phonenumbers.Parse(phone, "AT")
	if err != nil {
		return "", fmt.Errorf("invalid phone number: %w", err)
	}
	if !phonenumbers.IsPossibleNumber(num) || len(fmt.Sprintf("%d", *num.NationalNumber)) < 6 {
		return "", errors.New("invalid phone number")
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}

// Validate checks if the resident data is correct and follows the required formats.
func (r *Resident) Validate() error {
	if r.FirstName == "" {
		return errors.New("first name is required")
	}
	if r.LastName == "" {
		return errors.New("last name is required")
	}
	if r.RoomNumber <= 0 {
		return errors.New("room number must be a positive integer")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}
	if !emailRegex.MatchString(strings.ToLower(r.Email)) {
		return errors.New("email does not follow a verified format")
	}

	if r.PhoneNumber != "" {
		normalized, err := NormalizePhoneNumber(r.PhoneNumber)
		if err != nil {
			return err
		}
		r.PhoneNumber = normalized
	}

	if r.Birthday.IsZero() {
		return errors.New("birthday is required")
	}
	if r.DateSignedUp.IsZero() {
		return errors.New("date signed up is required")
	}
	if r.DateAdded.IsZero() {
		return errors.New("date added is required")
	}
	if r.DateModified.IsZero() {
		return errors.New("date modified is required")
	}

	return nil
}
