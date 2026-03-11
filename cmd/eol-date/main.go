// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/oliverandrich/eol-date/internal/api"
	"github.com/oliverandrich/eol-date/internal/search"
	"github.com/oliverandrich/eol-date/internal/ui"
	"github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	cmd := &cli.Command{
		Name:      "eol-date",
		Usage:     "Check end-of-life dates for software products",
		ArgsUsage: "<product>",
		Version:   version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "show all cycles including end-of-life versions",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output format: table, markdown, csv, html",
				Value:   "table",
			},
		},
		Action: run,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() < 1 {
		return fmt.Errorf("product name required\n\nUsage: eol-date <product>\n\nExample: eol-date python")
	}

	query := cmd.Args().First()
	showAll := cmd.Bool("all")
	format := cmd.String("format")

	products, err := api.FetchProducts(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch product list: %w", err)
	}

	product, found := search.FindExact(products, query)
	if !found {
		matches := search.FindSimilar(products, query, 10)
		if len(matches) == 0 {
			return fmt.Errorf("no products found matching '%s'", query)
		}

		selected, selectErr := ui.SelectProduct(matches)
		if selectErr != nil {
			return selectErr
		}
		product = selected
	}

	cycles, err := api.FetchProduct(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to fetch product details: %w", err)
	}

	ui.DisplayCycles(product, cycles, showAll, format)

	return nil
}
