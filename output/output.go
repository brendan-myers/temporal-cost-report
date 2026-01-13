package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/brendan-myers/temporal-cost-report/report"
	"github.com/olekukonko/tablewriter"
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
		"Actions %",
		"Active (GBh)",
		"Active %",
		"Retained (GBh)",
		"Retained %",
		"Action Cost",
		"Storage Cost",
		"Total Cost",
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeader(headers),
	)

	for _, ns := range r.Namespaces {
		storageCost := ns.ActiveStorageCost + ns.RetainedStorageCost
		table.Append([]string{
			ns.Name,
			formatNumber(ns.Actions),
			formatPercent(ns.ActionsPercent),
			fmt.Sprintf("%.2f", ns.ActiveStorageGBh),
			formatPercent(ns.ActiveStoragePercent),
			fmt.Sprintf("%.2f", ns.RetainedStorageGBh),
			formatPercent(ns.RetainedStoragePercent),
			formatCurrency(ns.ActionCost),
			formatCurrency(storageCost),
			formatCurrency(ns.TotalCost),
		})
	}

	// Add totals as footer
	table.Footer(
		"TOTAL",
		formatNumber(r.Totals.Actions),
		"100.00%",
		fmt.Sprintf("%.2f", r.Totals.ActiveStorageGBh),
		"100.00%",
		fmt.Sprintf("%.2f", r.Totals.RetainedStorageGBh),
		"100.00%",
		formatCurrency(r.Totals.ActionCost),
		formatCurrency(r.Totals.StorageCost),
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
