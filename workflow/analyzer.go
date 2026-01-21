package workflow

import (
	"context"

	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

// ActionCount holds the breakdown of billable actions for a workflow execution.
type ActionCount struct {
	WorkflowStarts    int `json:"workflowStarts"`
	Timers            int `json:"timers"`
	Signals           int `json:"signals"`
	SearchAttrUpserts int `json:"searchAttrUpserts"`
	Updates           int `json:"updates"`
	Activities        int `json:"activities"`
	ChildWorkflows    int `json:"childWorkflows"` // Each counts as 2 actions
	SideEffects       int `json:"sideEffects"`
	Total             int `json:"total"`
}

// AnalyzedExecution contains a workflow execution with its action count.
type AnalyzedExecution struct {
	Execution WorkflowExecution
	Actions   ActionCount
}

// CountActions analyzes a workflow's history and counts billable actions.
// Based on Temporal Cloud billing: https://docs.temporal.io/cloud/actions
func CountActions(ctx context.Context, c client.Client, workflowID, runID string) (ActionCount, error) {
	var count ActionCount

	iter := c.GetWorkflowHistory(ctx, workflowID, runID, false, enumspb.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)

	for iter.HasNext() {
		event, err := iter.Next()
		if err != nil {
			return count, err
		}

		switch event.EventType {
		// Workflow starts (1 action)
		case enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_STARTED:
			count.WorkflowStarts++

		// Timer starts (1 action each)
		case enumspb.EVENT_TYPE_TIMER_STARTED:
			count.Timers++

		// Signals (1 action each)
		case enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_SIGNALED:
			count.Signals++

		// Search attribute upserts (1 action each)
		case enumspb.EVENT_TYPE_UPSERT_WORKFLOW_SEARCH_ATTRIBUTES:
			count.SearchAttrUpserts++

		// Updates - accepted counts as action
		case enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_UPDATE_ACCEPTED:
			count.Updates++

		// Activity starts (1 action each, including retries)
		case enumspb.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED:
			count.Activities++

		// Child workflow starts (2 actions each)
		case enumspb.EVENT_TYPE_START_CHILD_WORKFLOW_EXECUTION_INITIATED:
			count.ChildWorkflows++

		// Side effects via markers (simplified: count all markers as potential side effects)
		case enumspb.EVENT_TYPE_MARKER_RECORDED:
			// Check if this is a SideEffect marker
			if attrs := event.GetMarkerRecordedEventAttributes(); attrs != nil {
				if attrs.MarkerName == "SideEffect" {
					count.SideEffects++
				}
			}
		}
	}

	// Calculate total (child workflows count as 2 actions each)
	count.Total = count.WorkflowStarts +
		count.Timers +
		count.Signals +
		count.SearchAttrUpserts +
		count.Updates +
		count.Activities +
		(count.ChildWorkflows * 2) +
		count.SideEffects

	return count, nil
}

// AnalyzeWorkflows fetches history for each execution and counts actions.
func AnalyzeWorkflows(ctx context.Context, c client.Client, executions []WorkflowExecution) ([]AnalyzedExecution, error) {
	results := make([]AnalyzedExecution, 0, len(executions))

	for _, exec := range executions {
		actions, err := CountActions(ctx, c, exec.WorkflowID, exec.RunID)
		if err != nil {
			return nil, err
		}

		results = append(results, AnalyzedExecution{
			Execution: exec,
			Actions:   actions,
		})
	}

	return results, nil
}
