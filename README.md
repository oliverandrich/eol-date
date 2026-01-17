# eol-date

CLI tool to check end-of-life dates for software products using the [endoflife.date](https://endoflife.date) API.

## Features

- Query EOL information for 300+ software products
- Shows active and end-of-life release cycles
- Displays release dates, support end dates, EOL dates, and LTS status
- Fuzzy search with interactive product selection
- Color-coded output (green = active, red = EOL)

## Installation

### Homebrew

```bash
brew install oliverandrich/tap/eol-date
```

### Go

```bash
go install github.com/oliverandrich/eol-date/cmd/eol-date@latest
```

### From Source

```bash
git clone https://github.com/oliverandrich/eol-date.git
cd eol-date
just build
```

## Usage

```bash
# Check EOL dates for a product
eol-date python

# Show all versions including EOL
eol-date python --all

# Fuzzy search (shows interactive selection)
eol-date post  # matches postgres, postgresql, etc.
```

### Example Output

```
Release cycles for python

╭───────┬─────────┬──────────────────────┬──────────────────────┬──────────────────────┬─────╮
│ CYCLE │ LATEST  │ RELEASED             │ SUPPORT              │ EOL                  │ LTS │
├───────┼─────────┼──────────────────────┼──────────────────────┼──────────────────────┼─────┤
│ 3.14  │ 3.14.2  │ 3m ago    2025-10-07 │ in 1y 8m  2027-10-01 │ in 4y 10m 2030-10-31 │     │
│ 3.13  │ 3.13.11 │ 1y 3m ago 2024-10-07 │ in 8m     2026-10-01 │ in 3y 10m 2029-10-31 │     │
│ 3.12  │ 3.12.12 │ 2y 3m ago 2023-10-02 │ 9m ago    2025-04-02 │ in 2y 9m  2028-10-31 │     │
╰───────┴─────────┴──────────────────────┴──────────────────────┴──────────────────────┴─────╯

5 active, 12 EOL (use --all to show)
```

## Column Description

| Column   | Description |
|----------|-------------|
| CYCLE    | Version/release cycle identifier |
| LATEST   | Latest patch version |
| RELEASED | Release date (relative + absolute) |
| SUPPORT  | Active support end date |
| EOL      | End-of-life date |
| LTS      | Long-term support indicator |

## Development

```bash
just build   # Build binary to build/eol-date
just test    # Run tests
just fmt     # Format code
just lint    # Run linter
just check   # Run fmt, lint, and test
```

## License

EUPL-1.2 - see [LICENSE](LICENSE) for details.
