package main

import (
	"fmt"
	"os"
	"time"

	"github.com/brendan-myers/temporal-cost-report/client"
	"github.com/brendan-myers/temporal-cost-report/output"
	"github.com/brendan-myers/temporal-cost-report/report"
	"github.com/spf13/cobra"
)

const (
	defaultActionPrice          = 50.0
	defaultActiveStoragePrice   = 0.042
	defaultRetainedStoragePrice = 0.00105
)

var (
	startDate            string
	endDate              string
	actionPrice          float64
	activeStoragePrice   float64
	retainedStoragePrice float64
	outputFormat string
	apiKey       string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "temporal-cost-report",
		Short: "Generate cost reports from Temporal Cloud usage data",
		Long: `A CLI tool that fetches usage data from Temporal Cloud and generates
cost reports per namespace for platform team chargebacks.

The tool reads the TEMPORAL_API_KEY environment variable for authentication.`,
		RunE: run,
	}

	// Disable alphabetical sorting of flags
	rootCmd.Flags().SortFlags = false

	// Date range flags
	rootCmd.Flags().StringVar(&startDate, "start-date", "", "Start date in YYYY-MM-DD format (default: first day of current month)")
	rootCmd.Flags().StringVar(&endDate, "end-date", "", "End date in YYYY-MM-DD format (default: today)")

	// Pricing flags
	rootCmd.Flags().Float64Var(&actionPrice, "action-price", defaultActionPrice, "Price per million actions (USD)")
	rootCmd.Flags().Float64Var(&activeStoragePrice, "active-storage-price", defaultActiveStoragePrice, "Price per GBh of active storage (USD)")
	rootCmd.Flags().Float64Var(&retainedStoragePrice, "retained-storage-price", defaultRetainedStoragePrice, "Price per GBh of retained storage (USD)")

	// Output format flag
	rootCmd.Flags().StringVar(&outputFormat, "format", "table", "Output format: table or json")

	// API key flag
	rootCmd.Flags().StringVar(&apiKey, "api-key", "", "Temporal Cloud API key (defaults to TEMPORAL_API_KEY env var)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Parse and validate dates
	start, end, err := parseDates(startDate, endDate)
	if err != nil {
		return err
	}

	// Validate output format
	if outputFormat != "table" && outputFormat != "json" {
		return fmt.Errorf("invalid format '%s': must be 'table' or 'json'", outputFormat)
	}

	// Create API client
	apiClient, err := client.New(apiKey)
	if err != nil {
		return err
	}

	// Fetch usage data
	summaries, err := apiClient.FetchUsage(start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to fetch usage data: %w", err)
	}

	// Generate report
	pricing := report.Pricing{
		ActionPricePerMillion:      actionPrice,
		ActiveStoragePricePerGBh:   activeStoragePrice,
		RetainedStoragePricePerGBh: retainedStoragePrice,
	}

	r := report.Generate(summaries, pricing, start.Format("2006-01-02"), end.Format("2006-01-02"))

	// Output report
	switch outputFormat {
	case "json":
		if err := output.PrintJSON(r); err != nil {
			return fmt.Errorf("failed to output JSON: %w", err)
		}
	default:
		output.PrintTable(r)
	}

	return nil
}

func parseDates(startStr, endStr string) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	var start, end time.Time
	var err error

	// Parse start date or default to first day of current month
	if startStr == "" {
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start date '%s': use YYYY-MM-DD format", startStr)
		}
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	}

	// Parse end date or default to today
	if endStr == "" {
		end = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end date '%s': use YYYY-MM-DD format", endStr)
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
	}

	// Validate date range
	if end.Before(start) {
		return time.Time{}, time.Time{}, fmt.Errorf("end date cannot be before start date")
	}

	// Add one day to end date for exclusive end time
	end = end.AddDate(0, 0, 1)

	return start, end, nil
}
