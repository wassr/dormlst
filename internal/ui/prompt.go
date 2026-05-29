package ui

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
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/wassr/dormlst/internal/model"
)

// Custom error types for graceful handling in the CLI
var (
	ErrNotFound = errors.New("no resident found")
	ErrAborted  = errors.New("aborted by user")
)

// PromptResident guides the user through entering or updating resident data.
func PromptResident(res model.Resident, isUpdate bool) (model.Resident, error) {
	if isUpdate {
		fmt.Printf("Updating Resident in Room %d: %s %s\n", res.RoomNumber, res.FirstName, res.LastName)
		fmt.Println("Press Enter to keep current value.")
	}

	validateEmpty := func(input string) error {
		if !isUpdate && input == "" {
			return fmt.Errorf("this field is required")
		}
		return nil
	}

	pFirstName := promptui.Prompt{Label: "First Name", Default: res.FirstName, Validate: validateEmpty}
	if val, err := pFirstName.Run(); err == nil && (val != "" || !isUpdate) {
		res.FirstName = val
	} else if err != nil {
		if err == promptui.ErrInterrupt {
			return res, ErrAborted
		}
		return res, err
	}

	pLastName := promptui.Prompt{Label: "Last Name", Default: res.LastName, Validate: validateEmpty}
	if val, err := pLastName.Run(); err == nil && (val != "" || !isUpdate) {
		res.LastName = val
	} else if err != nil {
		if err == promptui.ErrInterrupt {
			return res, ErrAborted
		}
		return res, err
	}

	pRoom := promptui.Prompt{
		Label: "Room Number",
		Default: func() string {
			if isUpdate {
				return strconv.Itoa(res.RoomNumber)
			}
			return ""
		}(),
		Validate: func(input string) error {
			if input == "" && isUpdate {
				return nil
			}
			_, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid room number")
			}
			return nil
		},
	}
	roomStr, err := pRoom.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			return res, ErrAborted
		}
		return res, err
	}
	if roomStr != "" {
		res.RoomNumber, _ = strconv.Atoi(roomStr)
	}

	pEmail := promptui.Prompt{
		Label:   "Email",
		Default: res.Email,
		Validate: func(input string) error {
			if err := validateEmpty(input); err != nil {
				return err
			}
			// Use the model's logic or a simple check for the UI
			_, err := mail.ParseAddress(input)
			if err != nil {
				return fmt.Errorf("invalid email format")
			}
			return nil
		},
	}
	if val, err := pEmail.Run(); err == nil && (val != "" || !isUpdate) {
		res.Email = val
	} else if err != nil {
		if err == promptui.ErrInterrupt {
			return res, ErrAborted
		}
		return res, err
	}

	pPhone := promptui.Prompt{
		Label:   "Phone",
		Default: res.PhoneNumber,
		Validate: func(input string) error {
			if input == "" {
				return nil
			}
			_, err := model.NormalizePhoneNumber(input)
			return err
		},
	}
	if val, err := pPhone.Run(); err == nil {
		if val != "" {
			res.PhoneNumber, _ = model.NormalizePhoneNumber(val)
		} else if !isUpdate {
			res.PhoneNumber = ""
		}
	} else {
		if err == promptui.ErrInterrupt {
			return res, ErrAborted
		}
		return res, err
	}

	pBirthday := promptui.Prompt{
		Label: "Birthday (YYYY-MM-DD)",
		Default: func() string {
			if isUpdate {
				return res.Birthday.Format("2006-01-02")
			}
			return ""
		}(),
		Validate: func(input string) error {
			if input == "" && isUpdate {
				return nil
			}
			_, err := time.Parse("2006-01-02", input)
			return err
		},
	}
	if val, err := pBirthday.Run(); err == nil && (val != "" || !isUpdate) {
		t, _ := time.Parse("2006-01-02", val)
		res.Birthday = t
	} else if err != nil {
		if err == promptui.ErrInterrupt {
			return res, ErrAborted
		}
		return res, err
	}

	pSignedUp := promptui.Prompt{
		Label: "Date Signed Up (YYYY-MM-DD)",
		Default: func() string {
			if isUpdate {
				return res.DateSignedUp.Format("2006-01-02")
			}
			return time.Now().Format("2006-01-02")
		}(),
		Validate: func(input string) error {
			if input == "" && isUpdate {
				return nil
			}
			_, err := time.Parse("2006-01-02", input)
			return err
		},
	}
	if val, err := pSignedUp.Run(); err == nil && (val != "" || !isUpdate) {
		t, _ := time.Parse("2006-01-02", val)
		res.DateSignedUp = t
	} else if err != nil {
		if err == promptui.ErrInterrupt {
			return res, ErrAborted
		}
		return res, err
	}

	if isUpdate {
		pActive := promptui.Select{
			Label: "Active Status",
			Items: []string{"Active", "Inactive"},
			CursorPos: func() int {
				if res.Active {
					return 0
				} else {
					return 1
				}
			}(),
		}
		if i, _, err := pActive.Run(); err == nil {
			res.Active = (i == 0)
		} else {
			if err == promptui.ErrInterrupt {
				return res, ErrAborted
			}
			return res, err
		}
	} else {
		res.Active = true
		res.DateAdded = time.Now()
		res.DateModified = res.DateAdded
	}

	return res, nil
}

// SelectResident prompts the user to select a single resident from a list.
// Returns the resident, the index in the input slice, and whether it was the ONLY option (suggesting autoselect could happen).
func SelectResident(residents []model.Resident, label string) (model.Resident, int, bool, error) {
	if len(residents) == 0 {
		return model.Resident{}, -1, false, ErrNotFound
	}
	if len(residents) == 1 {
		return residents[0], 0, true, nil
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ .RoomNumber | cyan }} ({{ .FirstName | red }} {{ .LastName | red }})",
		Inactive: "  {{ .RoomNumber | cyan }} ({{ .FirstName }} {{ .LastName }})",
		Selected: "\U0001F449 Selected: {{ .RoomNumber | cyan }} ({{ .FirstName }} {{ .LastName }})",
	}

	searcher := func(input string, index int) bool {
		r := residents[index]
		name := strings.ToLower(r.FirstName + " " + r.LastName)
		input = strings.ToLower(input)
		return strings.Contains(name, input) || strings.Contains(strconv.Itoa(r.RoomNumber), input)
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     residents,
		Templates: templates,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			return model.Resident{}, -1, false, ErrAborted
		}
		return model.Resident{}, -1, false, err
	}

	return residents[i], i, false, nil
}

// SelectFromMatches filters residents by query and prompts for selection if multiple matches exist.
// Returns the resident, the index in the original database, and a boolean indicating if it was autoselected (single match).
func SelectFromMatches(residents []model.Resident, query string, label string) (model.Resident, int, bool, error) {
	matched := model.FilterResidents(residents, query)
	if len(matched) == 0 {
		return model.Resident{}, -1, false, fmt.Errorf("%w matching: %s", ErrNotFound, query)
	}

	res, _, auto, err := SelectResident(matched, label)
	if err != nil {
		return model.Resident{}, -1, false, err
	}

	for i, r := range residents {
		if r.Email == res.Email && r.FirstName == res.FirstName && r.LastName == res.LastName && r.RoomNumber == res.RoomNumber {
			return r, i, auto, nil
		}
	}

	return model.Resident{}, -1, false, errors.New("failed to correlate selection with database")
}

// Confirm prompts the user for a yes/no confirmation.
func Confirm(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label + " [y/N]",
		IsConfirm: true,
	}
	res, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort || err == promptui.ErrInterrupt {
			return false, nil
		}
		return false, err
	}
	return strings.ToLower(res) == "y", nil
}
