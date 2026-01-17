// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

package ui

import (
	"encoding/csv"
	"fmt"
	"os"
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

// displayRow holds processed row data for output formatting
type displayRow struct {
	Cycle       string
	Latest      string
	ReleasedRel string // relative format (e.g., "3m ago")
	ReleasedRaw string // raw date (e.g., "2025-10-07")
	SupportRel  string // relative format
	SupportRaw  string // raw date or boolean as string
	EOLRel      string // relative format
	EOLRaw      string // raw date or boolean as string
	LTS         bool
	IsEOL       bool
}

// prepareDisplayRows converts cycles to displayRow slice
func prepareDisplayRows(cycles []api.Cycle, showAll bool) []displayRow {
	var rows []displayRow
	for _, c := range cycles {
		if !showAll && c.EOL.IsEOL() {
			continue
		}

		release := formatRelease(c.ReleaseDate.Time)
		support := formatSupport(c.Support)
		eol := formatEOL(c.EOL)

		row := displayRow{
			Cycle:       c.Cycle,
			Latest:      c.Latest,
			ReleasedRel: release.relative,
			ReleasedRaw: release.date,
			SupportRel:  support.relative,
			SupportRaw:  formatRawValue(c.Support),
			EOLRel:      eol.relative,
			EOLRaw:      formatRawValue(c.EOL),
			LTS:         c.LTS.IsLTS(),
			IsEOL:       c.EOL.IsEOL(),
		}
		rows = append(rows, row)
	}
	return rows
}

// formatRawValue returns the raw value for CSV/machine-readable output
func formatRawValue(v api.EOLValue) string {
	if v.IsBoolean {
		if v.BoolValue {
			return "true"
		}
		return "false"
	}
	if v.DateValue.IsZero() {
		return ""
	}
	return v.DateValue.Format("2006-01-02")
}

// DisplayCycles prints the release cycles in the specified format
func DisplayCycles(product string, cycles []api.Cycle, showAll bool, format string) {
	rows := prepareDisplayRows(cycles, showAll)

	if len(rows) == 0 {
		if showAll {
			fmt.Println("No release cycles found for", product)
		} else {
			fmt.Println("No active release cycles found for", product)
			fmt.Println(dimStyle.Render("Use --all to show end-of-life versions"))
		}
		return
	}

	switch format {
	case "markdown":
		formatAsMarkdown(product, rows)
	case "csv":
		formatAsCSV(rows)
	case "html":
		formatAsHTML(product, rows)
	default:
		formatAsTable(product, cycles, rows, showAll)
	}
}

// formatAsTable renders the lipgloss table (original format)
func formatAsTable(product string, cycles []api.Cycle, rows []displayRow, showAll bool) {
	fmt.Println()
	fmt.Println(headerStyle.Render(fmt.Sprintf("Release cycles for %s", product)))
	fmt.Println()

	// Calculate column widths for combined cells
	releasedWidth, supportWidth, eolWidth := 0, 0, 0
	for _, r := range rows {
		if w := len(r.ReleasedRel) + 1 + len(r.ReleasedRaw); w > releasedWidth {
			releasedWidth = w
		}
		supportDisplay := r.SupportRel
		if supportDisplay != "" && r.SupportRaw != "" && r.SupportRaw != "true" && r.SupportRaw != "false" {
			if w := len(r.SupportRel) + 1 + len(r.SupportRaw); w > supportWidth {
				supportWidth = w
			}
		} else if len(supportDisplay) > supportWidth {
			supportWidth = len(supportDisplay)
		}
		eolDisplay := r.EOLRel
		if eolDisplay != "" && r.EOLRaw != "" && r.EOLRaw != "true" && r.EOLRaw != "false" {
			if w := len(r.EOLRel) + 1 + len(r.EOLRaw); w > eolWidth {
				eolWidth = w
			}
		} else if len(eolDisplay) > eolWidth {
			eolWidth = len(eolDisplay)
		}
	}

	dimColor := lipgloss.Color("240")
	tableRows := make([][]string, 0, len(rows))
	for _, r := range rows {
		rowColor := lipgloss.Color("42") // green
		if r.IsEOL {
			rowColor = lipgloss.Color("203") // red
		}

		releaseRel := relativeDate{r.ReleasedRel, r.ReleasedRaw}
		supportDate := r.SupportRaw
		if supportDate == "true" || supportDate == "false" {
			supportDate = ""
		}
		supportRel := relativeDate{r.SupportRel, supportDate}
		eolDate := r.EOLRaw
		if eolDate == "true" || eolDate == "false" {
			eolDate = ""
		}
		eolRel := relativeDate{r.EOLRel, eolDate}

		ltsStr := ""
		if r.LTS {
			ltsStr = "✔"
		}

		tableRows = append(tableRows, []string{
			r.Cycle,
			r.Latest,
			combinedCell(releaseRel, rowColor, dimColor, releasedWidth),
			combinedCell(supportRel, rowColor, dimColor, supportWidth),
			combinedCell(eolRel, rowColor, dimColor, eolWidth),
			ltsStr,
		})
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		Headers("CYCLE", "LATEST", "RELEASED", "SUPPORT", "EOL", "LTS").
		Rows(tableRows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			baseStyle := lipgloss.NewStyle().Padding(0, 1)

			if col == 5 {
				baseStyle = baseStyle.Align(lipgloss.Center)
			}

			if row == table.HeaderRow {
				if col == 5 {
					return tableHeaderStyle.Padding(0, 1).Align(lipgloss.Center)
				}
				return tableHeaderStyle.Padding(0, 1)
			}

			if rows[row].IsEOL {
				baseStyle = baseStyle.Foreground(lipgloss.Color("203"))
			} else {
				baseStyle = baseStyle.Foreground(lipgloss.Color("42"))
			}

			if col == 5 && rows[row].LTS {
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

// formatAsMarkdown renders a Markdown table
func formatAsMarkdown(product string, rows []displayRow) {
	fmt.Printf("# Release cycles for %s\n\n", product)
	fmt.Println("| CYCLE | LATEST | RELEASED | SUPPORT | EOL | LTS |")
	fmt.Println("|-------|--------|----------|---------|-----|-----|")

	for _, r := range rows {
		released := formatMarkdownDate(r.ReleasedRel, r.ReleasedRaw)
		support := formatMarkdownDate(r.SupportRel, r.SupportRaw)
		eol := formatMarkdownDate(r.EOLRel, r.EOLRaw)
		lts := ""
		if r.LTS {
			lts = "✔"
		}

		fmt.Printf("| %s | %s | %s | %s | %s | %s |\n",
			r.Cycle, r.Latest, released, support, eol, lts)
	}
}

// formatMarkdownDate combines relative and raw date for Markdown output
func formatMarkdownDate(rel, raw string) string {
	if rel == "" && raw == "" {
		return ""
	}
	if raw == "" || raw == "true" || raw == "false" {
		return rel
	}
	if rel == "" {
		return raw
	}
	return fmt.Sprintf("%s (%s)", rel, raw)
}

// formatAsCSV renders CSV output
func formatAsCSV(rows []displayRow) {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	_ = w.Write([]string{"CYCLE", "LATEST", "RELEASED", "SUPPORT", "EOL", "LTS"})

	for _, r := range rows {
		lts := "false"
		if r.LTS {
			lts = "true"
		}
		_ = w.Write([]string{
			r.Cycle,
			r.Latest,
			r.ReleasedRaw,
			r.SupportRaw,
			r.EOLRaw,
			lts,
		})
	}
}

// formatAsHTML renders an HTML table
func formatAsHTML(product string, rows []displayRow) {
	fmt.Printf("<h1>Release cycles for %s</h1>\n", product)
	fmt.Println("<table>")
	fmt.Println("  <thead>")
	fmt.Println("    <tr><th>CYCLE</th><th>LATEST</th><th>RELEASED</th><th>SUPPORT</th><th>EOL</th><th>LTS</th></tr>")
	fmt.Println("  </thead>")
	fmt.Println("  <tbody>")

	for _, r := range rows {
		color := "green"
		if r.IsEOL {
			color = "red"
		}

		released := formatHTMLDate(r.ReleasedRel, r.ReleasedRaw)
		support := formatHTMLDate(r.SupportRel, r.SupportRaw)
		eol := formatHTMLDate(r.EOLRel, r.EOLRaw)
		lts := ""
		if r.LTS {
			lts = "✔"
		}

		fmt.Printf("    <tr style=\"color: %s;\"><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
			color, r.Cycle, r.Latest, released, support, eol, lts)
	}

	fmt.Println("  </tbody>")
	fmt.Println("</table>")
}

// formatHTMLDate combines relative and raw date for HTML output
func formatHTMLDate(rel, raw string) string {
	if rel == "" && raw == "" {
		return ""
	}
	if raw == "" || raw == "true" || raw == "false" {
		return rel
	}
	if rel == "" {
		return raw
	}
	return fmt.Sprintf("%s (%s)", rel, raw)
}
