// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/oliverandrich/eol-date/internal/api"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct { //nolint:govet // field alignment irrelevant for test struct
		name     string
		duration time.Duration
		want     string
	}{
		{duration: 12 * time.Hour, name: "less than a day", want: "<1d"},
		{duration: 24 * time.Hour, name: "one day", want: "1d"},
		{duration: 15 * 24 * time.Hour, name: "multiple days", want: "15d"},
		{duration: 30 * 24 * time.Hour, name: "one month", want: "1m"},
		{duration: 90 * 24 * time.Hour, name: "multiple months", want: "3m"},
		{duration: 365 * 24 * time.Hour, name: "one year", want: "1y"},
		{duration: 450 * 24 * time.Hour, name: "one year and months", want: "1y 3m"},
		{duration: 800 * 24 * time.Hour, name: "multiple years", want: "2y 2m"},
		{duration: 730 * 24 * time.Hour, name: "exact years", want: "2y"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}

func TestFormatRelease(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		releaseTime  time.Time
		wantRelative string
		wantDate     string
	}{
		{
			name:         "zero time",
			releaseTime:  time.Time{},
			wantRelative: "",
			wantDate:     "",
		},
		{
			name:         "recent release",
			releaseTime:  now.AddDate(0, -3, 0),
			wantRelative: "3m ago",
			wantDate:     now.AddDate(0, -3, 0).Format("2006-01-02"),
		},
		{
			name:         "old release",
			releaseTime:  now.AddDate(-2, -6, 0),
			wantRelative: "2y 6m ago",
			wantDate:     now.AddDate(-2, -6, 0).Format("2006-01-02"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatRelease(tt.releaseTime)
			if got.relative != tt.wantRelative {
				t.Errorf("formatRelease().relative = %q, want %q", got.relative, tt.wantRelative)
			}
			if got.date != tt.wantDate {
				t.Errorf("formatRelease().date = %q, want %q", got.date, tt.wantDate)
			}
		})
	}
}

func TestFormatSupport(t *testing.T) {
	now := time.Now()
	futureDate := now.AddDate(1, 6, 0)
	pastDate := now.AddDate(-1, -3, 0)

	tests := []struct {
		name         string
		support      api.EOLValue
		wantRelative string
		wantHasDate  bool
	}{
		{
			name:         "boolean true (active)",
			support:      api.EOLValue{IsBoolean: true, BoolValue: true},
			wantRelative: "Active",
			wantHasDate:  false,
		},
		{
			name:         "boolean false (no info)",
			support:      api.EOLValue{IsBoolean: true, BoolValue: false},
			wantRelative: "-",
			wantHasDate:  false,
		},
		{
			name:         "future date",
			support:      api.EOLValue{IsBoolean: false, DateValue: futureDate},
			wantRelative: "in 1y 6m",
			wantHasDate:  true,
		},
		{
			name:         "past date",
			support:      api.EOLValue{IsBoolean: false, DateValue: pastDate},
			wantRelative: "1y 3m ago",
			wantHasDate:  true,
		},
		{
			name:         "zero date",
			support:      api.EOLValue{IsBoolean: false, DateValue: time.Time{}},
			wantRelative: "",
			wantHasDate:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatSupport(tt.support)
			if got.relative != tt.wantRelative {
				t.Errorf("formatSupport().relative = %q, want %q", got.relative, tt.wantRelative)
			}
			if tt.wantHasDate && got.date == "" {
				t.Error("formatSupport().date is empty, expected a date")
			}
			if !tt.wantHasDate && got.date != "" {
				t.Errorf("formatSupport().date = %q, expected empty", got.date)
			}
		})
	}
}

func TestFormatEOL(t *testing.T) {
	now := time.Now()
	futureDate := now.AddDate(2, 0, 0)
	pastDate := now.AddDate(0, -6, 0)

	tests := []struct {
		name         string
		eol          api.EOLValue
		wantRelative string
		wantHasDate  bool
	}{
		{
			name:         "boolean true (ended)",
			eol:          api.EOLValue{IsBoolean: true, BoolValue: true},
			wantRelative: "Ended",
			wantHasDate:  false,
		},
		{
			name:         "boolean false (active)",
			eol:          api.EOLValue{IsBoolean: true, BoolValue: false},
			wantRelative: "Active",
			wantHasDate:  false,
		},
		{
			name:         "future date",
			eol:          api.EOLValue{IsBoolean: false, DateValue: futureDate},
			wantRelative: "in 2y",
			wantHasDate:  true,
		},
		{
			name:         "past date",
			eol:          api.EOLValue{IsBoolean: false, DateValue: pastDate},
			wantRelative: "6m ago",
			wantHasDate:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatEOL(tt.eol)
			if got.relative != tt.wantRelative {
				t.Errorf("formatEOL().relative = %q, want %q", got.relative, tt.wantRelative)
			}
			if tt.wantHasDate && got.date == "" {
				t.Error("formatEOL().date is empty, expected a date")
			}
			if !tt.wantHasDate && got.date != "" {
				t.Errorf("formatEOL().date = %q, expected empty", got.date)
			}
		})
	}
}

func TestCombinedCell(t *testing.T) {
	green := lipgloss.Color("42")
	dim := lipgloss.Color("240")

	tests := []struct {
		name  string
		rel   relativeDate
		width int
	}{
		{
			name:  "empty",
			rel:   relativeDate{"", ""},
			width: 20,
		},
		{
			name:  "only relative",
			rel:   relativeDate{"Active", ""},
			width: 20,
		},
		{
			name:  "both values",
			rel:   relativeDate{"3m ago", "2025-10-07"},
			width: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := combinedCell(tt.rel, green, dim, tt.width)
			if tt.rel.relative == "" && tt.rel.date == "" {
				if result != "" {
					t.Errorf("combinedCell() = %q, want empty string", result)
				}
				return
			}
			// For non-empty cases, just verify we get some output
			if result == "" {
				t.Error("combinedCell() returned empty string for non-empty input")
			}
		})
	}
}

func TestPrepareDisplayRows(t *testing.T) {
	futureDate := time.Now().AddDate(2, 0, 0)
	pastDate := time.Now().AddDate(-1, 0, 0)

	cycles := []api.Cycle{
		{
			Cycle:       "1.0",
			Latest:      "1.0.5",
			ReleaseDate: api.Date{Time: pastDate},
			EOL:         api.EOLValue{IsBoolean: false, DateValue: futureDate},
			Support:     api.EOLValue{IsBoolean: false, DateValue: futureDate},
			LTS:         api.LTSValue{IsBoolean: true, BoolValue: true},
		},
		{
			Cycle:       "0.9",
			Latest:      "0.9.10",
			ReleaseDate: api.Date{Time: pastDate.AddDate(-1, 0, 0)},
			EOL:         api.EOLValue{IsBoolean: false, DateValue: pastDate},
			Support:     api.EOLValue{IsBoolean: false, DateValue: pastDate},
			LTS:         api.LTSValue{IsBoolean: true, BoolValue: false},
		},
	}

	t.Run("showAll=false filters EOL", func(t *testing.T) {
		rows := prepareDisplayRows(cycles, false)
		if len(rows) != 1 {
			t.Errorf("expected 1 row, got %d", len(rows))
		}
		if rows[0].Cycle != "1.0" {
			t.Errorf("expected cycle 1.0, got %s", rows[0].Cycle)
		}
	})

	t.Run("showAll=true includes all", func(t *testing.T) {
		rows := prepareDisplayRows(cycles, true)
		if len(rows) != 2 {
			t.Errorf("expected 2 rows, got %d", len(rows))
		}
	})

	t.Run("LTS flag is set correctly", func(t *testing.T) {
		rows := prepareDisplayRows(cycles, true)
		if !rows[0].LTS {
			t.Error("expected LTS=true for cycle 1.0")
		}
		if rows[1].LTS {
			t.Error("expected LTS=false for cycle 0.9")
		}
	})
}

func TestFormatRawValue(t *testing.T) {
	tests := []struct {
		name string
		val  api.EOLValue
		want string
	}{
		{
			name: "boolean true",
			val:  api.EOLValue{IsBoolean: true, BoolValue: true},
			want: "true",
		},
		{
			name: "boolean false",
			val:  api.EOLValue{IsBoolean: true, BoolValue: false},
			want: "false",
		},
		{
			name: "date value",
			val:  api.EOLValue{IsBoolean: false, DateValue: time.Date(2025, 10, 7, 0, 0, 0, 0, time.UTC)},
			want: "2025-10-07",
		},
		{
			name: "zero date",
			val:  api.EOLValue{IsBoolean: false, DateValue: time.Time{}},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatRawValue(tt.val)
			if got != tt.want {
				t.Errorf("formatRawValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatMarkdownDate(t *testing.T) {
	tests := []struct {
		name string
		rel  string
		raw  string
		want string
	}{
		{name: "both empty", rel: "", raw: "", want: ""},
		{name: "only relative", rel: "Active", raw: "", want: "Active"},
		{name: "boolean true", rel: "Active", raw: "true", want: "Active"},
		{name: "boolean false", rel: "-", raw: "false", want: "-"},
		{name: "with date", rel: "3m ago", raw: "2025-10-07", want: "3m ago (2025-10-07)"},
		{name: "only date", rel: "", raw: "2025-10-07", want: "2025-10-07"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatMarkdownDate(tt.rel, tt.raw)
			if got != tt.want {
				t.Errorf("formatMarkdownDate(%q, %q) = %q, want %q", tt.rel, tt.raw, got, tt.want)
			}
		})
	}
}

// captureStdout captures stdout output from a function
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestFormatAsCSV(t *testing.T) {
	rows := []displayRow{
		{
			Cycle:       "3.14",
			Latest:      "3.14.2",
			ReleasedRel: "3m ago",
			ReleasedRaw: "2025-10-07",
			SupportRel:  "in 1y 8m",
			SupportRaw:  "2027-10-01",
			EOLRel:      "in 4y 10m",
			EOLRaw:      "2030-10-31",
			LTS:         false,
			IsEOL:       false,
		},
		{
			Cycle:       "3.13",
			Latest:      "3.13.11",
			ReleasedRel: "1y ago",
			ReleasedRaw: "2024-10-07",
			SupportRel:  "in 8m",
			SupportRaw:  "2026-10-01",
			EOLRel:      "in 3y 10m",
			EOLRaw:      "2029-10-31",
			LTS:         true,
			IsEOL:       false,
		},
	}

	output := captureStdout(func() {
		formatAsCSV(rows)
	})

	// Check header
	if !strings.Contains(output, "CYCLE,LATEST,RELEASED,SUPPORT,EOL,LTS") {
		t.Error("CSV output missing header")
	}

	// Check first row
	if !strings.Contains(output, "3.14,3.14.2,2025-10-07,2027-10-01,2030-10-31,false") {
		t.Error("CSV output missing first row data")
	}

	// Check second row with LTS
	if !strings.Contains(output, "3.13,3.13.11,2024-10-07,2026-10-01,2029-10-31,true") {
		t.Error("CSV output missing second row data")
	}
}

func TestFormatAsMarkdown(t *testing.T) {
	rows := []displayRow{
		{
			Cycle:       "3.14",
			Latest:      "3.14.2",
			ReleasedRel: "3m ago",
			ReleasedRaw: "2025-10-07",
			SupportRel:  "in 1y 8m",
			SupportRaw:  "2027-10-01",
			EOLRel:      "in 4y 10m",
			EOLRaw:      "2030-10-31",
			LTS:         false,
			IsEOL:       false,
		},
	}

	output := captureStdout(func() {
		formatAsMarkdown("python", rows)
	})

	// Check header
	if !strings.Contains(output, "# Release cycles for python") {
		t.Error("Markdown output missing title")
	}

	// Check table header
	if !strings.Contains(output, "| CYCLE | LATEST | RELEASED | SUPPORT | EOL | LTS |") {
		t.Error("Markdown output missing table header")
	}

	// Check separator
	if !strings.Contains(output, "|-------|--------|----------|---------|-----|-----|") {
		t.Error("Markdown output missing separator")
	}

	// Check data row
	if !strings.Contains(output, "| 3.14 | 3.14.2 |") {
		t.Error("Markdown output missing data row")
	}
}

func TestFormatAsHTML(t *testing.T) {
	rows := []displayRow{
		{
			Cycle:       "3.14",
			Latest:      "3.14.2",
			ReleasedRel: "3m ago",
			ReleasedRaw: "2025-10-07",
			SupportRel:  "in 1y 8m",
			SupportRaw:  "2027-10-01",
			EOLRel:      "in 4y 10m",
			EOLRaw:      "2030-10-31",
			LTS:         true,
			IsEOL:       false,
		},
		{
			Cycle:       "2.7",
			Latest:      "2.7.18",
			ReleasedRel: "5y ago",
			ReleasedRaw: "2020-04-20",
			SupportRel:  "Ended",
			SupportRaw:  "true",
			EOLRel:      "4y ago",
			EOLRaw:      "2020-01-01",
			LTS:         false,
			IsEOL:       true,
		},
	}

	output := captureStdout(func() {
		formatAsHTML("python", rows)
	})

	// Check structure
	if !strings.Contains(output, "<h1>Release cycles for python</h1>") {
		t.Error("HTML output missing title")
	}
	if !strings.Contains(output, "<table>") {
		t.Error("HTML output missing table tag")
	}
	if !strings.Contains(output, "<thead>") {
		t.Error("HTML output missing thead")
	}
	if !strings.Contains(output, "<tbody>") {
		t.Error("HTML output missing tbody")
	}

	// Check green row (active)
	if !strings.Contains(output, `style="color: green;"`) {
		t.Error("HTML output missing green style for active version")
	}

	// Check red row (EOL)
	if !strings.Contains(output, `style="color: red;"`) {
		t.Error("HTML output missing red style for EOL version")
	}

	// Check LTS checkmark
	if !strings.Contains(output, "✔") {
		t.Error("HTML output missing LTS checkmark")
	}
}
