package output

import (
	"encoding/json"
	"fmt"
	"os"

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
		"Actions",
		"Action Cost",
		"Actions %",
		"Active (GBh)",
		"Active Cost",
		"Active %",
		"Retained (GBh)",
		"Retained Cost",
		"Retained %",
		"Total Cost",
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeader(headers),
		tablewriter.WithRowAlignmentConfig(tw.CellAlignment{
			PerColumn: []tw.Align{
				tw.AlignLeft,  // Namespace
				tw.AlignRight, // Actions
				tw.AlignRight, // Action Cost
				tw.AlignRight, // Actions %
				tw.AlignRight, // Active (GBh)
				tw.AlignRight, // Active Cost
				tw.AlignRight, // Active %
				tw.AlignRight, // Retained (GBh)
				tw.AlignRight, // Retained Cost
				tw.AlignRight, // Retained %
				tw.AlignRight, // Total Cost
			},
		}),
		tablewriter.WithFooterAlignmentConfig(tw.CellAlignment{
			PerColumn: []tw.Align{
				tw.AlignLeft,  // TOTAL label
				tw.AlignRight, // Actions
				tw.AlignRight, // Action Cost
				tw.AlignRight, // Actions %
				tw.AlignRight, // Active (GBh)
				tw.AlignRight, // Active Cost
				tw.AlignRight, // Active %
				tw.AlignRight, // Retained (GBh)
				tw.AlignRight, // Retained Cost
				tw.AlignRight, // Retained %
				tw.AlignRight, // Total Cost
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
		})
	}

	// Add totals as footer
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
	)

	table.Render()
	fmt.Println()
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
