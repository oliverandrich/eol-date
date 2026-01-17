// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/oliverandrich/eol-date/internal/api"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("252"))
)

// formatDuration formats a duration as "Xy Xm" or "Xm" or "Xd"
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	months := days / 30
	years := months / 12

	if years > 0 {
		remainingMonths := months % 12
		if remainingMonths > 0 {
			return fmt.Sprintf("%dy %dm", years, remainingMonths)
		}
		return fmt.Sprintf("%dy", years)
	}
	if months > 0 {
		return fmt.Sprintf("%dm", months)
	}
	if days > 0 {
		return fmt.Sprintf("%dd", days)
	}
	return "<1d"
}

// relativeDate holds both the relative string and the date
type relativeDate struct {
	relative string
	date     string
}

// combinedCell creates a single string with relative left-aligned and date right-aligned
func combinedCell(rel relativeDate, relColor, dateColor lipgloss.Color, width int) string {
	if rel.relative == "" && rel.date == "" {
		return ""
	}

	relStyle := lipgloss.NewStyle().Foreground(relColor)
	dateStyle := lipgloss.NewStyle().Foreground(dateColor)

	relStr := relStyle.Render(rel.relative)
	dateStr := ""
	if rel.date != "" {
		dateStr = dateStyle.Render(rel.date)
	}

	// Calculate visible lengths (without ANSI codes)
	relLen := len(rel.relative)
	dateLen := len(rel.date)

	// Calculate padding needed between relative and date
	padding := max(1, width-relLen-dateLen)

	if dateStr == "" {
		return relStr
	}

	return relStr + strings.Repeat(" ", padding) + dateStr
}

// formatRelease formats a release date
func formatRelease(t time.Time) relativeDate {
	if t.IsZero() {
		return relativeDate{"", ""}
	}
	diff := time.Since(t)
	return relativeDate{
		relative: fmt.Sprintf("%s ago", formatDuration(diff)),
		date:     t.Format("2006-01-02"),
	}
}

// formatSupport formats a support end date
func formatSupport(support api.EOLValue) relativeDate {
	if support.IsBoolean {
		if support.BoolValue {
			// true = support is active
			return relativeDate{"Active", ""}
		}
		// false = no support info available
		return relativeDate{"-", ""}
	}
	if support.DateValue.IsZero() {
		return relativeDate{"", ""}
	}

	now := time.Now()
	diff := support.DateValue.Sub(now)

	if diff > 0 {
		return relativeDate{
			relative: fmt.Sprintf("in %s", formatDuration(diff)),
			date:     support.DateValue.Format("2006-01-02"),
		}
	}
	return relativeDate{
		relative: fmt.Sprintf("%s ago", formatDuration(-diff)),
		date:     support.DateValue.Format("2006-01-02"),
	}
}

// formatEOL formats an EOL date
func formatEOL(eol api.EOLValue) relativeDate {
	if eol.IsBoolean {
		if eol.BoolValue {
			return relativeDate{"Ended", ""}
		}
		return relativeDate{"Active", ""}
	}
	if eol.DateValue.IsZero() {
		return relativeDate{"", ""}
	}

	now := time.Now()
	diff := eol.DateValue.Sub(now)

	if diff > 0 {
		return relativeDate{
			relative: fmt.Sprintf("in %s", formatDuration(diff)),
			date:     eol.DateValue.Format("2006-01-02"),
		}
	}
	return relativeDate{
		relative: fmt.Sprintf("%s ago", formatDuration(-diff)),
		date:     eol.DateValue.Format("2006-01-02"),
	}
}

// DisplayCycles prints the release cycles in a formatted table
func DisplayCycles(product string, cycles []api.Cycle, showAll bool) {
	var displayCycles []api.Cycle
	if showAll {
		displayCycles = cycles
	} else {
		for _, c := range cycles {
			if !c.EOL.IsEOL() {
				displayCycles = append(displayCycles, c)
			}
		}
	}

	if len(displayCycles) == 0 {
		if showAll {
			fmt.Println("No release cycles found for", product)
		} else {
			fmt.Println("No active release cycles found for", product)
			fmt.Println(dimStyle.Render("Use --all to show end-of-life versions"))
		}
		return
	}

	fmt.Println()
	fmt.Println(headerStyle.Render(fmt.Sprintf("Release cycles for %s", product)))
	fmt.Println()

	// Prepare data and calculate column widths
	type rowData struct {
		cycle    api.Cycle
		release  relativeDate
		support  relativeDate
		eol      relativeDate
		ltsStr   string
		rowColor lipgloss.Color
	}

	data := make([]rowData, 0, len(displayCycles))
	releasedWidth, supportWidth, eolWidth := 0, 0, 0

	for _, c := range displayCycles {
		release := formatRelease(c.ReleaseDate.Time)
		support := formatSupport(c.Support)
		eol := formatEOL(c.EOL)

		ltsStr := ""
		if c.LTS.IsLTS() {
			ltsStr = "✔"
		}

		rowColor := lipgloss.Color("42") // green
		if c.EOL.IsEOL() {
			rowColor = lipgloss.Color("203") // red
		}

		data = append(data, rowData{c, release, support, eol, ltsStr, rowColor})

		// Calculate widths (relative + 1 space + date)
		if w := len(release.relative) + 1 + len(release.date); w > releasedWidth {
			releasedWidth = w
		}
		if w := len(support.relative) + 1 + len(support.date); w > supportWidth {
			supportWidth = w
		}
		if w := len(eol.relative) + 1 + len(eol.date); w > eolWidth {
			eolWidth = w
		}
	}

	dimColor := lipgloss.Color("240")
	rows := make([][]string, 0, len(data))
	for _, d := range data {
		rows = append(rows, []string{
			d.cycle.Cycle,
			d.cycle.Latest,
			combinedCell(d.release, d.rowColor, dimColor, releasedWidth),
			combinedCell(d.support, d.rowColor, dimColor, supportWidth),
			combinedCell(d.eol, d.rowColor, dimColor, eolWidth),
			d.ltsStr,
		})
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		Headers("CYCLE", "LATEST", "RELEASED", "SUPPORT", "EOL", "LTS").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			baseStyle := lipgloss.NewStyle().Padding(0, 1)

			// LTS column is centered
			if col == 5 {
				baseStyle = baseStyle.Align(lipgloss.Center)
			}

			if row == table.HeaderRow {
				if col == 5 {
					return tableHeaderStyle.Padding(0, 1).Align(lipgloss.Center)
				}
				return tableHeaderStyle.Padding(0, 1)
			}

			cycle := displayCycles[row]
			if cycle.EOL.IsEOL() {
				baseStyle = baseStyle.Foreground(lipgloss.Color("203"))
			} else {
				baseStyle = baseStyle.Foreground(lipgloss.Color("42"))
			}

			// LTS column gets special styling
			if col == 5 && cycle.LTS.IsLTS() {
				return baseStyle.Foreground(lipgloss.Color("220"))
			}

			return baseStyle
		})

	fmt.Println(t.Render())
	fmt.Println()

	activeCount := 0
	eolCount := 0
	for _, c := range cycles {
		if c.EOL.IsEOL() {
			eolCount++
		} else {
			activeCount++
		}
	}

	summary := fmt.Sprintf("%d active", activeCount)
	if eolCount > 0 && !showAll {
		summary += dimStyle.Render(fmt.Sprintf(", %d EOL (use --all to show)", eolCount))
	} else if eolCount > 0 {
		summary += fmt.Sprintf(", %d EOL", eolCount)
	}
	fmt.Println(dimStyle.Render(summary))
}
