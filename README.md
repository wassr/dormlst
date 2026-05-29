# dormlst

`dormlst` is a simple Go CLI tool for managing student dormitory resident lists. It stores data in a Git-friendly CSV format and exports to professional Excel (.xlsx) files.

## Features

- **Git-Friendly Storage**: CSV entries are automatically sorted (Room > Name) to ensure clean, deterministic Git diffs.
- **Interactive Management**: Built-in prompts for adding and updating residents with real-time validation.
- **Strict Validation**: E.164 phone number normalization and strict email format verification.
- **Excel Export**: Custom YAML mapping to define Excel column headers and data sources.
- **System Audit**: Comprehensive status command to check data integrity and system health.

## Quick Start

1. **Initialize**: Create a default configuration.
   ```bash
   dormlst init
   ```
2. **Add Residents**: Use flags or interactive mode (`-i`).
   ```bash
   dormlst add -i
   ```
3. **Search & Show**: Find residents and view detailed profiles.
   ```bash
   dormlst search "Noel"
   dormlst show 104
   ```
4. **Export**: Generate the Excel file based on your config.
   ```bash
   dormlst convert
   ```
5. **Check Health**: Audit the database for validation errors.
   ```bash
   dormlst status
   ```

## Configuration (`dormlst.yaml`)

The configuration defines the paths for the CSV database and Excel output, as well as the column mapping for the Excel export.

```yaml
database:
  path: "residents.csv"
output:
  path: "studentenliste.xlsx"
  date_format: "02.01.2006" # German format
  columns:
    - target: "Vorname"
      source: "first_name"
    - target: "Nachname"
      source: "last_name"
```

## Development

- **Language**: Go 1.24+
- **Core Libs**: Cobra (CLI), PromptUI (UX), Excelize (XLSX), Phonenumbers (Validation).
- **Project Structure**:
  - `cmd/`: CLI command definitions.
  - `internal/model/`: Resident struct and validation.
  - `internal/csvdb/`: CSV persistence logic.
  - `internal/ui/`: Shared interactive prompt logic.
  - `internal/excel/`: Excel generation.
