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
	"os"

	"gopkg.in/yaml.v3"
)

type ColumnMapping struct {
	Target string `yaml:"target"`
	Source string `yaml:"source"`
}

type OutputConfig struct {
	Path       string          `yaml:"path"`
	SheetName  string          `yaml:"sheet_name"`
	DateFormat string          `yaml:"date_format"`
	Columns    []ColumnMapping `yaml:"columns"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Output   OutputConfig   `yaml:"output"`
}

// DefaultConfig returns a sane default configuration.
func DefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Path: "residents.csv",
		},
		Output: OutputConfig{
			Path:       "studentenliste.xlsx",
			SheetName:  "Sheet1",
			DateFormat: "02.01.2006",
			Columns: []ColumnMapping{
				{Target: "Vorname", Source: "first_name"},
				{Target: "Nachname", Source: "last_name"},
				{Target: "Email", Source: "email"},
				{Target: "Tel", Source: "phone"},
				{Target: "Zimmernummer", Source: "room"},
				{Target: "Geburtstag", Source: "birthday"},
				{Target: "Eintragsdatum", Source: "entry_date"},
			},
		},
	}
}

// Load loads the config from a YAML file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves the config to a YAML file.
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
