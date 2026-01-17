// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

package api //nolint:revive // package name is intentional

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://endoflife.date/api"

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// FetchProducts retrieves the list of all product names from endoflife.date
func FetchProducts(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/all.json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var products []string
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

// FetchProduct retrieves the release cycles for a specific product
func FetchProduct(ctx context.Context, name string) ([]Cycle, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s.json", baseURL, name), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product %s: %w", name, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product %s not found", name)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d for product %s", resp.StatusCode, name)
	}

	var cycles []Cycle
	if err := json.NewDecoder(resp.Body).Decode(&cycles); err != nil {
		return nil, fmt.Errorf("failed to decode cycles for %s: %w", name, err)
	}

	return cycles, nil
}
