package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/brendan-myers/temporal-cost-report/workflow"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// PrintWorkflowTable outputs the workflow cost report as formatted ASCII tables.
func PrintWorkflowTable(r *workflow.WorkflowCostReport) {
	fmt.Println()
	fmt.Println("Workflow Cost Analysis")
	fmt.Printf("Type: %s\n", r.WorkflowType)
	fmt.Printf("Namespace: %s\n", r.Namespace)
	if r.SampleSize > 0 {
		fmt.Printf("Sample: %d executions (%s to %s)\n", r.SampleSize, r.Period.Start, r.Period.End)
	}
	fmt.Printf("Pricing: $%.2f/M actions\n", r.ActionPricePerMillion)
	fmt.Println()

	if r.SampleSize == 0 {
		fmt.Println("No completed workflows found for this type.")
		fmt.Println()
		return
	}

	// Summary table
	summaryHeaders := []string{"Metric", "Value"}
	summaryTable := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeader(summaryHeaders),
		tablewriter.WithHeaderAlignmentConfig(tw.CellAlignment{
			PerColumn: []tw.Align{tw.AlignLeft, tw.AlignRight},
		}),
		tablewriter.WithRowAlignmentConfig(tw.CellAlignment{
			PerColumn: []tw.Align{tw.AlignLeft, tw.AlignRight},
		}),
	)

	summaryTable.Append([]string{"Min Actions/Exec", fmt.Sprintf("%d", r.MinActionsPerExec)})
	summaryTable.Append([]string{"Max Actions/Exec", fmt.Sprintf("%d", r.MaxActionsPerExec)})
	summaryTable.Append([]string{"Avg Actions/Exec", fmt.Sprintf("%.1f", r.AverageActionsPerExec)})
	summaryTable.Append([]string{"Avg Cost/Exec", fmt.Sprintf("$%.6f", r.AverageCostPerExec)})
	summaryTable.Append([]string{"Executions Sampled", fmt.Sprintf("%d", r.SampleSize)})
	summaryTable.Append([]string{"Sample Period (days)", fmt.Sprintf("%.1f", r.PeriodDays)})
	summaryTable.Append([]string{"Est. Monthly Execs", formatNumber(float64(r.EstimatedMonthlyExecs))})
	summaryTable.Append([]string{"Est. Monthly Cost", fmt.Sprintf("$%.2f", r.EstimatedMonthlyCost)})

	summaryTable.Render()
	fmt.Println()

	// Action breakdown table
	fmt.Println("Action Breakdown (avg per execution):")
	breakdownHeaders := []string{"Event Type", "Count", "Actions"}
	breakdownTable := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeader(breakdownHeaders),
		tablewriter.WithHeaderAlignmentConfig(tw.CellAlignment{
			PerColumn: []tw.Align{tw.AlignLeft, tw.AlignRight, tw.AlignRight},
		}),
		tablewriter.WithRowAlignmentConfig(tw.CellAlignment{
			PerColumn: []tw.Align{tw.AlignLeft, tw.AlignRight, tw.AlignRight},
		}),
		tablewriter.WithFooterAlignmentConfig(tw.CellAlignment{
			PerColumn: []tw.Align{tw.AlignLeft, tw.AlignRight, tw.AlignRight},
		}),
	)

	b := r.AverageActionBreakdown

	if b.WorkflowStarts > 0 {
		breakdownTable.Append([]string{"Workflow Starts", fmt.Sprintf("%.1f", b.WorkflowStarts), fmt.Sprintf("%.1f", b.WorkflowStarts)})
	}
	if b.Activities > 0 {
		breakdownTable.Append([]string{"Activities", fmt.Sprintf("%.1f", b.Activities), fmt.Sprintf("%.1f", b.Activities)})
	}
	if b.Timers > 0 {
		breakdownTable.Append([]string{"Timers", fmt.Sprintf("%.1f", b.Timers), fmt.Sprintf("%.1f", b.Timers)})
	}
	if b.Signals > 0 {
		breakdownTable.Append([]string{"Signals", fmt.Sprintf("%.1f", b.Signals), fmt.Sprintf("%.1f", b.Signals)})
	}
	if b.ChildWorkflows > 0 {
		breakdownTable.Append([]string{"Child Workflows", fmt.Sprintf("%.1f", b.ChildWorkflows), fmt.Sprintf("%.1f", b.ChildWorkflows*2)})
	}
	if b.Updates > 0 {
		breakdownTable.Append([]string{"Updates", fmt.Sprintf("%.1f", b.Updates), fmt.Sprintf("%.1f", b.Updates)})
	}
	if b.SearchAttrUpserts > 0 {
		breakdownTable.Append([]string{"Search Attr Upserts", fmt.Sprintf("%.1f", b.SearchAttrUpserts), fmt.Sprintf("%.1f", b.SearchAttrUpserts)})
	}
	if b.SideEffects > 0 {
		breakdownTable.Append([]string{"Side Effects", fmt.Sprintf("%.1f", b.SideEffects), fmt.Sprintf("%.1f", b.SideEffects)})
	}

	breakdownTable.Footer("TOTAL", "", fmt.Sprintf("%.1f", b.TotalActions))
	breakdownTable.Render()

	fmt.Println()
	fmt.Println("* Costs are estimates based on sampled data and may differ from actual invoiced amounts.")
	fmt.Println()
}

// PrintWorkflowJSON outputs the workflow cost report as formatted JSON.
func PrintWorkflowJSON(r *workflow.WorkflowCostReport) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r)
}
