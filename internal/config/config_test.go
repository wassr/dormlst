package config

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
	"path/filepath"
	"testing"
)

func TestConfig_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "test_config.yaml")

	cfg := DefaultConfig()
	cfg.Database.Path = "test_db.csv"

	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	loaded, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loaded.Database.Path != "test_db.csv" {
		t.Errorf("Expected database path test_db.csv, got %s", loaded.Database.Path)
	}
}

func TestLoad_NonExistent(t *testing.T) {
	_, err := Load("this_config_does_not_exist.yaml")
	if err == nil {
		t.Fatal("Expected error when loading non-existent config, got nil")
	}
}
