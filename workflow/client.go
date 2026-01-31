package workflow

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/brendan-myers/temporal-cost-report/config"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

// nopLogger is a no-op logger that discards all log messages.
type nopLogger struct{}

func (nopLogger) Debug(string, ...any) {}
func (nopLogger) Info(string, ...any)  {}
func (nopLogger) Warn(string, ...any)  {}
func (nopLogger) Error(string, ...any) {}

// NewTemporalClient creates a new Temporal client using the provided configuration.
// It supports both mTLS certificate authentication and API key authentication.
// If both are configured, mTLS takes precedence.
func NewTemporalClient(cfg config.ConnectionConfig) (client.Client, error) {
	opts := client.Options{
		HostPort:  cfg.Address,
		Namespace: cfg.Namespace,
		Logger:    nopLogger{},
	}

	if cfg.UsesMTLS() {
		// Load mTLS certificates
		cert, err := tls.LoadX509KeyPair(cfg.TLSCertPath, cfg.TLSKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS certificates: %w", err)
		}

		opts.ConnectionOptions = client.ConnectionOptions{
			TLS: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}
	} else {
		// Use API key authentication
		apiKey := cfg.APIKey
		if apiKey == "" {
			apiKey = os.Getenv("TEMPORAL_API_KEY")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("API key required: set TEMPORAL_API_KEY environment variable or use --api-key flag")
		}
		opts.Credentials = client.NewAPIKeyStaticCredentials(apiKey)
	}

	c, err := client.Dial(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create Temporal client: %w", err)
	}

	return c, nil
}

// WorkflowExecution represents a workflow execution to analyze.
type WorkflowExecution struct {
	WorkflowID string
	RunID      string
	StartTime  int64 // Unix timestamp in nanoseconds
	CloseTime  int64 // Unix timestamp in nanoseconds
}

// ListWorkflowsByType queries completed workflows of the given type.
func ListWorkflowsByType(ctx context.Context, c client.Client, namespace, workflowType string, limit int) ([]WorkflowExecution, error) {
	query := fmt.Sprintf("WorkflowType = '%s' AND CloseTime IS NOT NULL", workflowType)

	var executions []WorkflowExecution
	var nextPageToken []byte

	for {
		resp, err := c.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
			Namespace:     namespace,
			Query:         query,
			PageSize:      int32(min(limit-len(executions), 100)),
			NextPageToken: nextPageToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list workflows: %w", err)
		}

		for _, exec := range resp.Executions {
			var startNano, closeNano int64
			if exec.StartTime != nil {
				startNano = exec.StartTime.AsTime().UnixNano()
			}
			if exec.CloseTime != nil {
				closeNano = exec.CloseTime.AsTime().UnixNano()
			}

			executions = append(executions, WorkflowExecution{
				WorkflowID: exec.Execution.WorkflowId,
				RunID:      exec.Execution.RunId,
				StartTime:  startNano,
				CloseTime:  closeNano,
			})

			if len(executions) >= limit {
				return executions, nil
			}
		}

		nextPageToken = resp.NextPageToken
		if len(nextPageToken) == 0 {
			break
		}
	}

	return executions, nil
}
