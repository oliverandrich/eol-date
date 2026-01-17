// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

package search

import (
	"sort"
	"strings"
)

// FindExact performs a case-insensitive exact match search
func FindExact(products []string, query string) (string, bool) {
	queryLower := strings.ToLower(query)
	for _, p := range products {
		if strings.ToLower(p) == queryLower {
			return p, true
		}
	}
	return "", false
}

// FindSimilar returns products that contain the query as substring (case-insensitive)
func FindSimilar(products []string, query string, limit int) []string {
	queryLower := strings.ToLower(query)

	var results []string
	for _, p := range products {
		if strings.Contains(strings.ToLower(p), queryLower) {
			results = append(results, p)
		}
	}

	sort.Strings(results)

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}
