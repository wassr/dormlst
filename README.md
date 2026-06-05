# dormlst

`dormlst` is a Go CLI tool for managing student dormitory resident lists without the headaches of Excel. It uses a Git-friendly CSV as the source of truth to ensure data integrity and easy synchronization, generating .xlsx files only when needed.

## Installation

If you have `go` installed simply run:

```bash
go install github.com/wassr/dormlst@latest
```

Else download the prebuild binary from the releasse page [here](https://github.com/wassr/dormlst/releases).

Or build it from source yourself.

```bash
git clone https://github.com/wassr/dormlst.git
cd dormlst
go install .
```

## Why not just use Excel?

Excel is notoriously bad for version control, often corrupts formatting (like phone numbers and dates), and makes collaboration difficult. `dormlst` solves this by:
- Storing everything in a sorted, deterministic CSV.
- Normalizing all data (E.164 phone numbers, standardized dates) on entry.
- Generating the required Excel artifact only as a temporary output for external services.

## Features

- **Git-Friendly Storage**: CSV entries are automatically sorted (Room > Name) for clean, meaningful Git diffs.
- **Interactive Management**: Built-in prompts for adding and updating residents with real-time validation.
- **Explicit Modes**: Support for both guided interactive modes (`-i`) and surgical flag-based updates.
- **Active Status Management**: Easily `disable` residents (e.g., for exchange students) and `enable` them later without deleting data.
- **Strict Validation**: E.164 phone number normalization and verified email formats.
- **Modern CLI UX**: A clean, unified design using status symbols (✓/✗), bold sections, and flush-left formatting.
- **Excel Generation**: Custom YAML mapping to match the exact format required by external vendors.

## Quick Start

1. **Initialize**: Create a default configuration.
   ```bash
   dormlst init
   ```
2. **Add Residents**: Use flags or interactive mode (`-i`).
   ```bash
   dormlst add -i
   ```
3. **Update**: Use flags of interactive mode (`-i`)
   ```bash
   dormlst update "Noel" --room 555
   dormlst update -i "Noel"
   ```
4. **Search & Show**: Find residents and view detailed profiles.
   ```bash
   dormlst se "Noel"
   dormlst show 666
   ```
5. **Enable/Disable**: Toggle resident status.
   ```bash
   dormlst dis "Noel" # Set to inactive
   dormlst en "Noel"  # Set to active
   ```
6. **Generate Upload File**: Produce the .xlsx for external use.
   ```bash
   dormlst gen
   ```
7. **Vacuum**: Permanently remove all disabled entries.
   ```bash
   dormlst vac --yes
   ```
8. **Check Health**: Audit the database for any validation errors.
   ```bash
   dormlst status
   ```

## Commands & Aliases

| Command    | Alias  | Description                                      |
|------------|--------|--------------------------------------------------|
| `add`      | -      | Create a new resident entry                      |
| `update`   | `up`   | Modify an existing resident                      |
| `remove`   | `rm`   | Delete a resident from the CSV                   |
| `enable`   | `en`   | Set a resident's status to active                |
| `disable`  | `dis`  | Set a resident's status to inactive              |
| `vacuum`   | `vac`  | Permanently remove all disabled entries          |
| `search`   | `se`   | Find residents based on specific criteria        |
| `show`     | -      | See all available data on one person             |
| `generate` | `gen`  | Generate the .xlsx file required for uploads     |
| `upload`   | `push` | Upload the resident list (Stubbed)               |
| `status`   | -      | Get stats and health check on the database       |
| `init`     | -      | Initialize a default configuration file          |


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
