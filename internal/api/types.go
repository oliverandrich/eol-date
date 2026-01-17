// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

package api //nolint:revive // package name is intentional

import (
	"encoding/json"
	"time"
)

// Cycle represents a release cycle from endoflife.date
type Cycle struct {
	ReleaseDate       Date     `json:"releaseDate"`
	LatestReleaseDate Date     `json:"latestReleaseDate"`
	EOL               EOLValue `json:"eol"`
	Support           EOLValue `json:"support"`
	LTS               LTSValue `json:"lts"`
	Cycle             string   `json:"cycle"`
	Latest            string   `json:"latest"`
}

// Date handles date parsing from the API (YYYY-MM-DD format)
type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if json.Unmarshal(data, &s) != nil {
		d.Time = time.Time{}
		return nil //nolint:nilerr // lenient parsing: accept invalid data
	}
	if s == "" {
		d.Time = time.Time{}
		return nil
	}
	t, parseErr := time.Parse("2006-01-02", s)
	if parseErr != nil {
		d.Time = time.Time{}
		return nil //nolint:nilerr // lenient parsing: accept invalid date format
	}
	d.Time = t
	return nil
}

// EOLValue can be a boolean (false = still supported, true = EOL) or a date string
type EOLValue struct {
	DateValue time.Time
	IsBoolean bool
	BoolValue bool
}

func (e *EOLValue) UnmarshalJSON(data []byte) error {
	var b bool
	if json.Unmarshal(data, &b) == nil {
		e.IsBoolean = true
		e.BoolValue = b
		return nil
	}

	var s string
	if json.Unmarshal(data, &s) == nil {
		e.IsBoolean = false
		t, parseErr := time.Parse("2006-01-02", s)
		if parseErr != nil {
			e.BoolValue = true
			e.IsBoolean = true
			return nil //nolint:nilerr // lenient parsing: treat invalid date as EOL
		}
		e.DateValue = t
		return nil
	}

	return nil
}

// IsEOL returns true if the product has reached end of life
func (e *EOLValue) IsEOL() bool {
	if e.IsBoolean {
		return e.BoolValue
	}
	return !e.DateValue.IsZero() && time.Now().After(e.DateValue)
}

// String returns a string representation of the EOL value
func (e *EOLValue) String() string {
	if e.IsBoolean {
		if e.BoolValue {
			return "Yes"
		}
		return "No"
	}
	if e.DateValue.IsZero() {
		return "N/A"
	}
	return e.DateValue.Format("2006-01-02")
}

// LTSValue can be a boolean or a date string for LTS releases
type LTSValue struct {
	DateValue time.Time
	IsBoolean bool
	BoolValue bool
}

func (l *LTSValue) UnmarshalJSON(data []byte) error {
	var b bool
	if json.Unmarshal(data, &b) == nil {
		l.IsBoolean = true
		l.BoolValue = b
		return nil
	}

	var s string
	if json.Unmarshal(data, &s) == nil {
		l.IsBoolean = false
		t, parseErr := time.Parse("2006-01-02", s)
		if parseErr != nil {
			l.BoolValue = false
			l.IsBoolean = true
			return nil //nolint:nilerr // lenient parsing: treat invalid date as non-LTS
		}
		l.DateValue = t
		return nil
	}

	return nil
}

// IsLTS returns true if this is an LTS release
func (l *LTSValue) IsLTS() bool {
	if l.IsBoolean {
		return l.BoolValue
	}
	return !l.DateValue.IsZero()
}

// String returns a string representation of the LTS value
func (l *LTSValue) String() string {
	if l.IsBoolean {
		if l.BoolValue {
			return "Yes"
		}
		return "No"
	}
	if l.DateValue.IsZero() {
		return "No"
	}
	return l.DateValue.Format("2006-01-02")
}
