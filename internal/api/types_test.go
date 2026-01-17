// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

package api //nolint:revive // package name is intentional

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		want    time.Time
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "valid date",
			json: `"2025-10-07"`,
			want: time.Date(2025, 10, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "empty string",
			json: `""`,
			want: time.Time{},
		},
		{
			name: "null",
			json: `null`,
			want: time.Time{},
		},
		{
			name: "invalid format",
			json: `"2025/10/07"`,
			want: time.Time{},
		},
		{
			name: "invalid value",
			json: `123`,
			want: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Date
			err := json.Unmarshal([]byte(tt.json), &d)
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !d.Equal(tt.want) {
				t.Errorf("Date.UnmarshalJSON() = %v, want %v", d.Time, tt.want)
			}
		})
	}
}

func TestEOLValue_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		wantDate      time.Time
		name          string
		json          string
		wantIsBoolean bool
		wantBoolValue bool
	}{
		{
			name:          "boolean true",
			json:          `true`,
			wantIsBoolean: true,
			wantBoolValue: true,
		},
		{
			name:          "boolean false",
			json:          `false`,
			wantIsBoolean: true,
			wantBoolValue: false,
		},
		{
			name:          "date string",
			json:          `"2025-10-31"`,
			wantIsBoolean: false,
			wantDate:      time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "invalid date string",
			json:          `"invalid"`,
			wantIsBoolean: true,
			wantBoolValue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e EOLValue
			if err := json.Unmarshal([]byte(tt.json), &e); err != nil {
				t.Errorf("EOLValue.UnmarshalJSON() error = %v", err)
				return
			}
			if e.IsBoolean != tt.wantIsBoolean {
				t.Errorf("EOLValue.IsBoolean = %v, want %v", e.IsBoolean, tt.wantIsBoolean)
			}
			if e.BoolValue != tt.wantBoolValue {
				t.Errorf("EOLValue.BoolValue = %v, want %v", e.BoolValue, tt.wantBoolValue)
			}
			if !e.DateValue.Equal(tt.wantDate) {
				t.Errorf("EOLValue.DateValue = %v, want %v", e.DateValue, tt.wantDate)
			}
		})
	}
}

func TestEOLValue_IsEOL(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		eol  EOLValue
		want bool
	}{
		{
			name: "boolean true",
			eol:  EOLValue{IsBoolean: true, BoolValue: true},
			want: true,
		},
		{
			name: "boolean false",
			eol:  EOLValue{IsBoolean: true, BoolValue: false},
			want: false,
		},
		{
			name: "future date",
			eol:  EOLValue{IsBoolean: false, DateValue: now.AddDate(1, 0, 0)},
			want: false,
		},
		{
			name: "past date",
			eol:  EOLValue{IsBoolean: false, DateValue: now.AddDate(-1, 0, 0)},
			want: true,
		},
		{
			name: "zero date",
			eol:  EOLValue{IsBoolean: false, DateValue: time.Time{}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.eol.IsEOL(); got != tt.want {
				t.Errorf("EOLValue.IsEOL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEOLValue_String(t *testing.T) {
	tests := []struct {
		name string
		eol  EOLValue
		want string
	}{
		{
			name: "boolean true",
			eol:  EOLValue{IsBoolean: true, BoolValue: true},
			want: "Yes",
		},
		{
			name: "boolean false",
			eol:  EOLValue{IsBoolean: true, BoolValue: false},
			want: "No",
		},
		{
			name: "date",
			eol:  EOLValue{IsBoolean: false, DateValue: time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)},
			want: "2025-10-31",
		},
		{
			name: "zero date",
			eol:  EOLValue{IsBoolean: false, DateValue: time.Time{}},
			want: "N/A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.eol.String(); got != tt.want {
				t.Errorf("EOLValue.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLTSValue_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		wantDate      time.Time
		name          string
		json          string
		wantIsBoolean bool
		wantBoolValue bool
	}{
		{
			name:          "boolean true",
			json:          `true`,
			wantIsBoolean: true,
			wantBoolValue: true,
		},
		{
			name:          "boolean false",
			json:          `false`,
			wantIsBoolean: true,
			wantBoolValue: false,
		},
		{
			name:          "date string",
			json:          `"2032-04-30"`,
			wantIsBoolean: false,
			wantDate:      time.Date(2032, 4, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l LTSValue
			if err := json.Unmarshal([]byte(tt.json), &l); err != nil {
				t.Errorf("LTSValue.UnmarshalJSON() error = %v", err)
				return
			}
			if l.IsBoolean != tt.wantIsBoolean {
				t.Errorf("LTSValue.IsBoolean = %v, want %v", l.IsBoolean, tt.wantIsBoolean)
			}
			if l.BoolValue != tt.wantBoolValue {
				t.Errorf("LTSValue.BoolValue = %v, want %v", l.BoolValue, tt.wantBoolValue)
			}
			if !l.DateValue.Equal(tt.wantDate) {
				t.Errorf("LTSValue.DateValue = %v, want %v", l.DateValue, tt.wantDate)
			}
		})
	}
}

func TestLTSValue_IsLTS(t *testing.T) {
	tests := []struct {
		name string
		lts  LTSValue
		want bool
	}{
		{
			name: "boolean true",
			lts:  LTSValue{IsBoolean: true, BoolValue: true},
			want: true,
		},
		{
			name: "boolean false",
			lts:  LTSValue{IsBoolean: true, BoolValue: false},
			want: false,
		},
		{
			name: "has date",
			lts:  LTSValue{IsBoolean: false, DateValue: time.Date(2032, 4, 30, 0, 0, 0, 0, time.UTC)},
			want: true,
		},
		{
			name: "zero date",
			lts:  LTSValue{IsBoolean: false, DateValue: time.Time{}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lts.IsLTS(); got != tt.want {
				t.Errorf("LTSValue.IsLTS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCycle_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"cycle": "3.13",
		"releaseDate": "2024-10-07",
		"eol": "2029-10-31",
		"support": "2026-10-01",
		"lts": false,
		"latest": "3.13.11",
		"latestReleaseDate": "2025-01-14"
	}`

	var cycle Cycle
	if err := json.Unmarshal([]byte(jsonData), &cycle); err != nil {
		t.Fatalf("Cycle.UnmarshalJSON() error = %v", err)
	}

	if cycle.Cycle != "3.13" {
		t.Errorf("Cycle.Cycle = %q, want %q", cycle.Cycle, "3.13")
	}
	if cycle.Latest != "3.13.11" {
		t.Errorf("Cycle.Latest = %q, want %q", cycle.Latest, "3.13.11")
	}
	if cycle.ReleaseDate.Format("2006-01-02") != "2024-10-07" {
		t.Errorf("Cycle.ReleaseDate = %v, want 2024-10-07", cycle.ReleaseDate)
	}
	if cycle.EOL.IsBoolean || cycle.EOL.DateValue.Format("2006-01-02") != "2029-10-31" {
		t.Errorf("Cycle.EOL = %v, want date 2029-10-31", cycle.EOL)
	}
	if cycle.LTS.IsLTS() {
		t.Error("Cycle.LTS.IsLTS() = true, want false")
	}
}
