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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/wassr/dormlst/internal/config"
	"github.com/wassr/dormlst/internal/model"
)

func TestGenerate(t *testing.T) {
	tmpDir := t.TempDir()
	xlsxPath := filepath.Join(tmpDir, "test.xlsx")

	residents := []model.Resident{
		{
			FirstName:    "John",
			LastName:     "Doe",
			RoomNumber:   101,
			Email:        "john@example.com",
			Birthday:     time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
			DateSignedUp: time.Now(),
			DateAdded:    time.Now(),
			DateModified: time.Now(),
			Active:       true,
		},
	}

	cfg := config.DefaultConfig().Output
	cfg.Path = xlsxPath

	err := Generate(residents, cfg)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if _, err := os.Stat(xlsxPath); os.IsNotExist(err) {
		t.Fatal("Excel file was not created")
	}
}
