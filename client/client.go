package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/brendan-myers/temporal-cost-report/models"
)

const (
	baseURL    = "https://saas-api.tmprl.cloud/cloud/usage"
	apiVersion = "2024-10-01-00"
)

// Client handles communication with the Temporal Cloud API.
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// New creates a new Temporal Cloud API client.
// If apiKey is empty, it reads from the TEMPORAL_API_KEY environment variable.
func New(apiKey string) (*Client, error) {
	if apiKey == "" {
		apiKey = os.Getenv("TEMPORAL_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("API key not provided: use --api-key flag or set TEMPORAL_API_KEY environment variable")
	}

	return &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}, nil
}

// FetchUsage retrieves usage data for the specified date range.
// Dates should be in RFC3339 format (e.g., "2026-01-01T00:00:00Z").
func (c *Client) FetchUsage(startTime, endTime string) ([]models.Summary, error) {
	var allSummaries []models.Summary
	pageToken := ""

	for {
		resp, err := c.doRequest(startTime, endTime, 1000, pageToken)
		if err != nil {
			return nil, err
		}

		allSummaries = append(allSummaries, resp.Summaries...)

		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	return allSummaries, nil
}

func (c *Client) doRequest(startTime, endTime string, pageSize int, pageToken string) (*models.GetUsageResponse, error) {
	// Build URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("start_time_inclusive", startTime)
	q.Set("end_time_exclusive", endTime)
	q.Set("page_size", strconv.Itoa(pageSize))
	if pageToken != "" {
		q.Set("page_token", pageToken)
	}
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("temporal-cloud-api-version", apiVersion)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var usageResp models.GetUsageResponse
	if err := json.Unmarshal(respBody, &usageResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &usageResp, nil
}
