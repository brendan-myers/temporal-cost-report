package workflow

import (
	"time"
)

// WorkflowCostReport contains the cost analysis for a workflow type.
type WorkflowCostReport struct {
	WorkflowType             string          `json:"workflowType"`
	Namespace                string          `json:"namespace"`
	SampleSize               int             `json:"sampleSize"`
	Period                   Period          `json:"period"`
	PeriodDays               float64         `json:"periodDays"`
	MinActionsPerExec        int             `json:"minActionsPerExecution"`
	MaxActionsPerExec        int             `json:"maxActionsPerExecution"`
	AverageActionsPerExec    float64         `json:"averageActionsPerExecution"`
	AverageCostPerExec       float64         `json:"averageCostPerExecution"`
	EstimatedMonthlyExecs    int             `json:"estimatedMonthlyExecutions"`
	EstimatedMonthlyCost     float64         `json:"estimatedMonthlyCost"`
	ActionPricePerMillion    float64         `json:"actionPricePerMillion"`
	AverageActionBreakdown   ActionBreakdown `json:"actionBreakdown"`
}

// Period represents a date range.
type Period struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// ActionBreakdown shows average actions per execution by type.
type ActionBreakdown struct {
	WorkflowStarts    float64 `json:"workflowStarts"`
	Timers            float64 `json:"timers"`
	Signals           float64 `json:"signals"`
	SearchAttrUpserts float64 `json:"searchAttrUpserts"`
	Updates           float64 `json:"updates"`
	Activities        float64 `json:"activities"`
	ChildWorkflows    float64 `json:"childWorkflows"`
	SideEffects       float64 `json:"sideEffects"`
	TotalActions      float64 `json:"totalActions"`
}

// GenerateReport creates a cost report from analyzed workflow executions.
func GenerateReport(workflowType, namespace string, executions []AnalyzedExecution, actionPricePerMillion float64) *WorkflowCostReport {
	if len(executions) == 0 {
		return &WorkflowCostReport{
			WorkflowType:          workflowType,
			Namespace:             namespace,
			SampleSize:            0,
			ActionPricePerMillion: actionPricePerMillion,
		}
	}

	// Find date range from executions
	var minStart, maxClose int64
	minStart = executions[0].Execution.StartTime
	maxClose = executions[0].Execution.CloseTime

	for _, exec := range executions {
		if exec.Execution.StartTime < minStart {
			minStart = exec.Execution.StartTime
		}
		if exec.Execution.CloseTime > maxClose {
			maxClose = exec.Execution.CloseTime
		}
	}

	startTime := time.Unix(0, minStart)
	endTime := time.Unix(0, maxClose)
	periodDays := endTime.Sub(startTime).Hours() / 24
	if periodDays < 1 {
		periodDays = 1 // Minimum 1 day to avoid division issues
	}

	// Sum up all actions and track min/max
	var totalActions ActionCount
	minActions := executions[0].Actions.Total
	maxActions := executions[0].Actions.Total

	for _, exec := range executions {
		totalActions.WorkflowStarts += exec.Actions.WorkflowStarts
		totalActions.Timers += exec.Actions.Timers
		totalActions.Signals += exec.Actions.Signals
		totalActions.SearchAttrUpserts += exec.Actions.SearchAttrUpserts
		totalActions.Updates += exec.Actions.Updates
		totalActions.Activities += exec.Actions.Activities
		totalActions.ChildWorkflows += exec.Actions.ChildWorkflows
		totalActions.SideEffects += exec.Actions.SideEffects
		totalActions.Total += exec.Actions.Total

		if exec.Actions.Total < minActions {
			minActions = exec.Actions.Total
		}
		if exec.Actions.Total > maxActions {
			maxActions = exec.Actions.Total
		}
	}

	sampleSize := float64(len(executions))

	// Calculate averages
	avgActions := float64(totalActions.Total) / sampleSize
	avgCost := (avgActions / 1_000_000) * actionPricePerMillion

	// Estimate monthly executions (scale to 30 days)
	execsPerDay := sampleSize / periodDays
	monthlyExecs := int(execsPerDay * 30)
	monthlyCost := float64(monthlyExecs) * avgCost

	return &WorkflowCostReport{
		WorkflowType:          workflowType,
		Namespace:             namespace,
		SampleSize:            len(executions),
		Period: Period{
			Start: startTime.Format("2006-01-02"),
			End:   endTime.Format("2006-01-02"),
		},
		PeriodDays:            periodDays,
		MinActionsPerExec:     minActions,
		MaxActionsPerExec:     maxActions,
		AverageActionsPerExec: avgActions,
		AverageCostPerExec:    avgCost,
		EstimatedMonthlyExecs: monthlyExecs,
		EstimatedMonthlyCost:  monthlyCost,
		ActionPricePerMillion: actionPricePerMillion,
		AverageActionBreakdown: ActionBreakdown{
			WorkflowStarts:    float64(totalActions.WorkflowStarts) / sampleSize,
			Timers:            float64(totalActions.Timers) / sampleSize,
			Signals:           float64(totalActions.Signals) / sampleSize,
			SearchAttrUpserts: float64(totalActions.SearchAttrUpserts) / sampleSize,
			Updates:           float64(totalActions.Updates) / sampleSize,
			Activities:        float64(totalActions.Activities) / sampleSize,
			ChildWorkflows:    float64(totalActions.ChildWorkflows) / sampleSize,
			SideEffects:       float64(totalActions.SideEffects) / sampleSize,
			TotalActions:      avgActions,
		},
	}
}
