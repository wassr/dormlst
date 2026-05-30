# dormlst

`dormlst` is a Go CLI tool for managing student dormitory resident lists without the headaches of Excel. It uses a Git-friendly CSV as the source of truth to ensure data integrity and easy synchronization, generating .xlsx files only when needed for uploads (e.g., to workitout.at).

## Why not just use Excel?

Excel is notoriously bad for version control, often corrupts formatting (like phone numbers and dates), and makes collaboration difficult. `dormlst` solves this by:
- Storing everything in a sorted, deterministic CSV.
- Normalizing all data (E.164 phone numbers, standardized dates) on entry.
- Generating the required Excel artifact only as a temporary output for external services.

## Features

- **Git-Friendly Storage**: CSV entries are automatically sorted (Room > Name) for clean, meaningful Git diffs.
- **Interactive Management**: Built-in prompts for adding and updating residents with real-time validation.
- **Strict Validation**: E.164 phone number normalization and verified email formats.
- **Excel Generation**: Custom YAML mapping to match the exact format required by external vendors.
- **System Audit**: Quick status command to check database health and statistics.

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
4. **Generate Upload File**: Produce the .xlsx for external use.
   ```bash
   dormlst generate
   ```
5. **Check Health**: Audit the database for any validation errors.
   ```bash
   dormlst status
   ```

## Configuration (`dormlst.yaml`)

Define your database path and the exact column mappings required for your Excel output.

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

- **Language**: Go 1.25.10
- **Core Libs**: Cobra (CLI), PromptUI (UX), Excelize (XLSX), Phonenumbers (Validation).
