// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

package ui

import (
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
