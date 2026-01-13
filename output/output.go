package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/brendan-myers/temporal-cost-report/report"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// PrintTable outputs the report as a formatted ASCII table.
func PrintTable(r *report.Report) {
	fmt.Println()
	fmt.Println("Temporal Cloud Usage Report")
	fmt.Printf("Period: %s to %s\n", r.Period.Start, r.Period.End)
	fmt.Printf("Pricing: $%.2f/M actions, $%.4f/GBh active, $%.5f/GBh retained\n",
		r.Pricing.ActionPricePerMillion,
		r.Pricing.ActiveStoragePricePerGBh,
		r.Pricing.RetainedStoragePricePerGBh)
	fmt.Println()

	headers := []string{
		"Namespace",
		"Count", "Cost", "%",
		"GBh", "Cost", "%",
		"GBh", "Cost", "%",
		"Cost", "%",
	}

	// First, render to buffer to get column widths
	var buf bytes.Buffer
	table := tablewriter.NewTable(&buf,
		tablewriter.WithHeader(headers),
		tablewriter.WithHeaderAlignmentConfig(tw.CellAlignment{
			PerColumn: []tw.Align{
				tw.AlignLeft,
				tw.AlignCenter, tw.AlignCenter, tw.AlignCenter,
				tw.AlignCenter, tw.AlignCenter, tw.AlignCenter,
				tw.AlignCenter, tw.AlignCenter, tw.AlignCenter,
				tw.AlignCenter, tw.AlignCenter,
			},
		}),
		tablewriter.WithRowAlignmentConfig(tw.CellAlignment{
			PerColumn: []tw.Align{
				tw.AlignLeft,
				tw.AlignRight, tw.AlignRight, tw.AlignRight,
				tw.AlignRight, tw.AlignRight, tw.AlignRight,
				tw.AlignRight, tw.AlignRight, tw.AlignRight,
				tw.AlignRight, tw.AlignRight,
			},
		}),
		tablewriter.WithFooterAlignmentConfig(tw.CellAlignment{
			PerColumn: []tw.Align{
				tw.AlignLeft,
				tw.AlignRight, tw.AlignRight, tw.AlignRight,
				tw.AlignRight, tw.AlignRight, tw.AlignRight,
				tw.AlignRight, tw.AlignRight, tw.AlignRight,
				tw.AlignRight, tw.AlignRight,
			},
		}),
	)

	for _, ns := range r.Namespaces {
		table.Append([]string{
			ns.Name,
			formatNumber(ns.Actions),
			formatCurrency(ns.ActionCost),
			formatPercent(ns.ActionsPercent),
			fmt.Sprintf("%.2f", ns.ActiveStorageGBh),
			formatCurrency(ns.ActiveStorageCost),
			formatPercent(ns.ActiveStoragePercent),
			fmt.Sprintf("%.2f", ns.RetainedStorageGBh),
			formatCurrency(ns.RetainedStorageCost),
			formatPercent(ns.RetainedStoragePercent),
			formatCurrency(ns.TotalCost),
			formatPercent(ns.TotalCostPercent),
		})
	}

	table.Footer(
		"TOTAL",
		formatNumber(r.Totals.Actions),
		formatCurrency(r.Totals.ActionCost),
		"100.00%",
		fmt.Sprintf("%.2f", r.Totals.ActiveStorageGBh),
		formatCurrency(r.Totals.ActiveStorageCost),
		"100.00%",
		fmt.Sprintf("%.2f", r.Totals.RetainedStorageGBh),
		formatCurrency(r.Totals.RetainedStorageCost),
		"100.00%",
		formatCurrency(r.Totals.TotalCost),
		"100.00%",
	)

	table.Render()

	// Parse the rendered table to get column positions from header row
	lines := strings.Split(buf.String(), "\n")
	if len(lines) < 3 {
		fmt.Print(buf.String())
		return
	}

	// Find the header row (second line, after the top border)
	headerLine := lines[1]

	// Build group header components
	topBorder, groupHeader, separator := buildGroupHeader(headerLine)

	// Print: top border, group header, separator, then rest of table (skipping original top border)
	fmt.Println(topBorder)
	fmt.Println(groupHeader)
	fmt.Println(separator)
	for i, line := range lines {
		if i == 0 { // Skip original top border
			continue
		}
		fmt.Println(line)
	}
	fmt.Println("* Costs are estimates based on the provided pricing and may differ from actual invoiced amounts.")
	fmt.Println()
}

// findColumnWidths parses the header row to find the display width of each column
func findColumnWidths(headerLine string) []int {
	var widths []int
	currentWidth := 0
	runes := []rune(headerLine)

	for i, r := range runes {
		if r == '│' {
			if i > 0 { // Skip the first separator
				widths = append(widths, currentWidth)
			}
			currentWidth = 0
		} else {
			currentWidth++
		}
	}
	return widths
}

// buildGroupHeader creates group header components: top border, header row, and separator
func buildGroupHeader(headerLine string) (topBorder, groupHeader, separator string) {
	widths := findColumnWidths(headerLine)
	if len(widths) < 12 {
		return "", "", ""
	}

	// Calculate group widths (including separators between columns in the group)
	// Group spans: Namespace (col 0), Actions (cols 1-3), Active Storage (cols 4-6),
	// Retained Storage (cols 7-9), Total (cols 10-11)
	namespaceWidth := widths[0]
	actionsWidth := widths[1] + 1 + widths[2] + 1 + widths[3]     // 3 columns + 2 separators
	activeWidth := widths[4] + 1 + widths[5] + 1 + widths[6]      // 3 columns + 2 separators
	retainedWidth := widths[7] + 1 + widths[8] + 1 + widths[9]    // 3 columns + 2 separators
	totalWidth := widths[10] + 1 + widths[11]                      // 2 columns + 1 separator

	groups := []struct {
		name  string
		width int
	}{
		{"", namespaceWidth},
		{"ACTIONS", actionsWidth},
		{"ACTIVE STORAGE", activeWidth},
		{"RETAINED STORAGE", retainedWidth},
		{"TOTAL", totalWidth},
	}

	// Build top border
	var topB strings.Builder
	topB.WriteString("┌")
	topB.WriteString(strings.Repeat("─", namespaceWidth))
	topB.WriteString("┬")
	topB.WriteString(strings.Repeat("─", actionsWidth))
	topB.WriteString("┬")
	topB.WriteString(strings.Repeat("─", activeWidth))
	topB.WriteString("┬")
	topB.WriteString(strings.Repeat("─", retainedWidth))
	topB.WriteString("┬")
	topB.WriteString(strings.Repeat("─", totalWidth))
	topB.WriteString("┐")

	// Build group header row (no leading │ for namespace column)
	var header strings.Builder
	header.WriteString("│")
	header.WriteString(strings.Repeat(" ", namespaceWidth))
	for i, g := range groups {
		if i == 0 { // Skip namespace, already handled
			continue
		}
		header.WriteString("│")
		// Center the group name
		padding := g.width - len(g.name)
		leftPad := padding / 2
		rightPad := padding - leftPad
		if leftPad < 0 {
			leftPad = 0
		}
		if rightPad < 0 {
			rightPad = 0
		}
		header.WriteString(strings.Repeat(" ", leftPad))
		header.WriteString(g.name)
		header.WriteString(strings.Repeat(" ", rightPad))
	}
	header.WriteString("│")

	// Build separator line between group header and column headers
	var sep strings.Builder
	sep.WriteString("├")
	sep.WriteString(strings.Repeat("─", namespaceWidth))
	sep.WriteString("┼")
	// Actions group: 3 columns
	sep.WriteString(strings.Repeat("─", widths[1]))
	sep.WriteString("┬")
	sep.WriteString(strings.Repeat("─", widths[2]))
	sep.WriteString("┬")
	sep.WriteString(strings.Repeat("─", widths[3]))
	sep.WriteString("┼")
	// Active Storage group: 3 columns
	sep.WriteString(strings.Repeat("─", widths[4]))
	sep.WriteString("┬")
	sep.WriteString(strings.Repeat("─", widths[5]))
	sep.WriteString("┬")
	sep.WriteString(strings.Repeat("─", widths[6]))
	sep.WriteString("┼")
	// Retained Storage group: 3 columns
	sep.WriteString(strings.Repeat("─", widths[7]))
	sep.WriteString("┬")
	sep.WriteString(strings.Repeat("─", widths[8]))
	sep.WriteString("┬")
	sep.WriteString(strings.Repeat("─", widths[9]))
	sep.WriteString("┼")
	// Total group: 2 columns
	sep.WriteString(strings.Repeat("─", widths[10]))
	sep.WriteString("┬")
	sep.WriteString(strings.Repeat("─", widths[11]))
	sep.WriteString("┤")

	return topB.String(), header.String(), sep.String()
}

// PrintJSON outputs the report as formatted JSON.
func PrintJSON(r *report.Report) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r)
}

func formatNumber(n float64) string {
	if n >= 1_000_000_000 {
		return fmt.Sprintf("%.2fB", n/1_000_000_000)
	}
	if n >= 1_000_000 {
		return fmt.Sprintf("%.2fM", n/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.2fK", n/1_000)
	}
	return fmt.Sprintf("%.0f", n)
}

func formatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

func formatPercent(pct float64) string {
	return fmt.Sprintf("%.2f%%", pct)
}
