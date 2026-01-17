// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

package search

import (
	"reflect"
	"testing"
)

func TestFindExact(t *testing.T) {
	products := []string{"python", "Python3", "nodejs", "go", "rust"}

	tests := []struct {
		name      string
		query     string
		wantMatch string
		wantFound bool
	}{
		{"exact match", "python", "python", true},
		{"case insensitive", "PYTHON", "python", true},
		{"mixed case", "Python3", "Python3", true},
		{"no match", "java", "", false},
		{"partial no match", "pyth", "", false},
		{"empty query", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := FindExact(products, tt.query)
			if found != tt.wantFound {
				t.Errorf("FindExact() found = %v, want %v", found, tt.wantFound)
			}
			if got != tt.wantMatch {
				t.Errorf("FindExact() = %q, want %q", got, tt.wantMatch)
			}
		})
	}
}

func TestFindSimilar(t *testing.T) {
	products := []string{"python", "python2", "python3", "postgres", "postgresql", "go", "golang"}

	tests := []struct { //nolint:govet // field alignment irrelevant for test struct
		name  string
		query string
		limit int
		want  []string
	}{
		{
			name:  "prefix match",
			query: "python",
			limit: 10,
			want:  []string{"python", "python2", "python3"},
		},
		{
			name:  "substring match",
			query: "post",
			limit: 10,
			want:  []string{"postgres", "postgresql"},
		},
		{
			name:  "case insensitive",
			query: "GO",
			limit: 10,
			want:  []string{"go", "golang"},
		},
		{
			name:  "limit results",
			query: "python",
			limit: 2,
			want:  []string{"python", "python2"},
		},
		{
			name:  "no match",
			query: "java",
			limit: 10,
			want:  nil,
		},
		{
			name:  "empty query matches all",
			query: "",
			limit: 3,
			want:  []string{"go", "golang", "postgres"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindSimilar(products, tt.query, tt.limit)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindSimilar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindSimilar_Sorted(t *testing.T) {
	products := []string{"zython", "python", "aython"}

	got := FindSimilar(products, "ython", 10)
	want := []string{"aython", "python", "zython"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindSimilar() results not sorted: got %v, want %v", got, want)
	}
}
