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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestE2ECLI acts as a comprehensive end-to-end integration test.
// It compiles the application into a temporary binary and runs it through its lifecycle.
func TestE2ECLI(t *testing.T) {
	// 1. Compile the binary
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "dormlst")

	// We need to point to where main.go is. In a standard project it's in the root.
	// Since we are in cmd/, the root is ../
	buildCmd := exec.Command("go", "build", "-o", binPath, "../main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI binary: %v", err)
	}

	// Prepare environment context
	configPath := filepath.Join(tmpDir, "dormlst.yaml")
	xlsxPath := filepath.Join(tmpDir, "studentenliste.xlsx")

	runCmd := func(args ...string) (string, string, error) {
		// Use the specific config file to isolate from the dev environment
		cmdArgs := append([]string{"--config", configPath}, args...)
		cmd := exec.Command(binPath, cmdArgs...)
		cmd.Dir = tmpDir // Run in tmpDir so defaults point there
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		return out.String(), stderr.String(), err
	}

	// Step 1: Initialize
	t.Run("init", func(t *testing.T) {
		out, _, err := runCmd("init")
		if err != nil {
			t.Fatalf("init failed: %v", err)
		}
		if !strings.Contains(out, "Initialized default configuration") {
			t.Errorf("Unexpected init output: %s", out)
		}
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Fatal("Config file not created")
		}
	})

	// Step 2: Add resident (Flag based)
	t.Run("add_resident", func(t *testing.T) {
		out, _, err := runCmd(
			"add",
			"--first-name", "John",
			"--last-name", "Doe",
			"--room", "101",
			"--email", "john.doe@example.com",
			"--birthday", "1995-05-15",
		)
		if err != nil {
			t.Fatalf("add failed: %v", err)
		}
		if !strings.Contains(out, "Successfully added resident") {
			t.Errorf("Unexpected add output: %s", out)
		}
	})

	// Step 3: Check status
	t.Run("status", func(t *testing.T) {
		out, _, err := runCmd("status")
		if err != nil {
			t.Fatalf("status failed: %v", err)
		}
		if !strings.Contains(out, "Total                1") || !strings.Contains(out, "Active               1") {
			t.Errorf("Status output missing expected metrics: %s", out)
		}
	})

	// Step 4: Show resident (finds automatically since 1 match)
	t.Run("show_resident", func(t *testing.T) {
		out, _, err := runCmd("show", "101")
		if err != nil {
			t.Fatalf("show failed: %v", err)
		}
		if !strings.Contains(out, "John Doe") {
			t.Errorf("Show did not find the resident: %s", out)
		}
	})

	// Step 5: Show nonexistent (Graceful exit)
	t.Run("show_nonexistent", func(t *testing.T) {
		out, _, err := runCmd("show", "999")
		if err != nil {
			t.Fatalf("show failed for nonexistent (should be graceful): %v", err)
		}
		if !strings.Contains(out, "no resident found") {
			t.Errorf("Expected not found message, got: %s", out)
		}
	})

	// Step 6: Update resident (Flag based)
	t.Run("update_resident", func(t *testing.T) {
		out, _, err := runCmd("update", "john", "--room", "102")
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}
		if !strings.Contains(out, "Successfully updated") {
			t.Errorf("Unexpected update output: %s", out)
		}

		// Verify update via search
		out, _, _ = runCmd("search", "john")
		if !strings.Contains(out, "102") {
			t.Errorf("Update did not change room number: %s", out)
		}
	})

	// Step 7: Disable
	t.Run("disable", func(t *testing.T) {
		out, _, err := runCmd("disable", "john")
		if err != nil {
			t.Fatalf("disable failed: %v", err)
		}
		if !strings.Contains(out, "Successfully disabled") {
			t.Errorf("Unexpected disable output: %s", out)
		}
	})

	// Step 8: Convert to Excel
	t.Run("convert", func(t *testing.T) {
		out, _, err := runCmd("convert")
		if err != nil {
			t.Fatalf("convert failed: %v", err)
		}
		if !strings.Contains(out, "Successfully converted") {
			t.Errorf("Unexpected convert output: %s", out)
		}
		if _, err := os.Stat(xlsxPath); os.IsNotExist(err) {
			t.Fatal("Excel file was not created")
		}
	})
}
