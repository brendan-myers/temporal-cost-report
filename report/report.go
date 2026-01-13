package report

import (
	"sort"

	"github.com/brendan-myers/temporal-cost-report/models"
)

// Pricing holds the configurable prices for cost calculation.
type Pricing struct {
	ActionPricePerMillion      float64 `json:"actionPricePerMillion"`
	ActiveStoragePricePerGBh   float64 `json:"activeStoragePricePerGBh"`
	RetainedStoragePricePerGBh float64 `json:"retainedStoragePricePerGBh"`
}

// NamespaceUsage holds aggregated usage data for a single namespace.
type NamespaceUsage struct {
	Name                   string  `json:"name"`
	Actions                float64 `json:"actions"`
	ActionsPercent         float64 `json:"actionsPercent"`
	ActiveStorageGBh       float64 `json:"activeStorageGBh"`
	ActiveStoragePercent   float64 `json:"activeStoragePercent"`
	RetainedStorageGBh     float64 `json:"retainedStorageGBh"`
	RetainedStoragePercent float64 `json:"retainedStoragePercent"`
	ActionCost             float64 `json:"actionCost"`
	ActiveStorageCost      float64 `json:"activeStorageCost"`
	RetainedStorageCost    float64 `json:"retainedStorageCost"`
	TotalCost              float64 `json:"totalCost"`
}

// Totals holds aggregated totals across all namespaces.
type Totals struct {
	Actions            float64 `json:"actions"`
	ActiveStorageGBh   float64 `json:"activeStorageGBh"`
	RetainedStorageGBh float64 `json:"retainedStorageGBh"`
	ActionCost         float64 `json:"actionCost"`
	StorageCost        float64 `json:"storageCost"`
	TotalCost          float64 `json:"totalCost"`
}

// Report contains the complete cost report data.
type Report struct {
	Period     Period           `json:"period"`
	Pricing    Pricing          `json:"pricing"`
	Namespaces []NamespaceUsage `json:"namespaces"`
	Totals     Totals           `json:"totals"`
}

// Period represents the date range for the report.
type Period struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// Generate creates a cost report from usage summaries.
func Generate(summaries []models.Summary, pricing Pricing, startDate, endDate string) *Report {
	// Aggregate usage by namespace
	namespaceData := make(map[string]*namespaceAggregator)

	for _, summary := range summaries {
		for _, group := range summary.RecordGroups {
			namespace := extractNamespace(group.GroupBys)
			if namespace == "" {
				continue
			}

			if _, exists := namespaceData[namespace]; !exists {
				namespaceData[namespace] = &namespaceAggregator{}
			}

			agg := namespaceData[namespace]
			for _, record := range group.Records {
				switch record.Type {
				case models.RecordTypeActions:
					agg.actions += record.Value
				case models.RecordTypeActiveStorage:
					agg.activeStorageByteSeconds += record.Value
				case models.RecordTypeRetainedStorage:
					agg.retainedStorageByteSeconds += record.Value
				}
			}
		}
	}

	// Convert to NamespaceUsage with cost calculations
	var namespaces []NamespaceUsage
	var totals Totals

	for name, agg := range namespaceData {
		usage := calculateNamespaceUsage(name, agg, pricing)
		namespaces = append(namespaces, usage)

		totals.Actions += usage.Actions
		totals.ActiveStorageGBh += usage.ActiveStorageGBh
		totals.RetainedStorageGBh += usage.RetainedStorageGBh
		totals.ActionCost += usage.ActionCost
		totals.StorageCost += usage.ActiveStorageCost + usage.RetainedStorageCost
		totals.TotalCost += usage.TotalCost
	}

	// Calculate percentages
	for i := range namespaces {
		if totals.Actions > 0 {
			namespaces[i].ActionsPercent = (namespaces[i].Actions / totals.Actions) * 100
		}
		if totals.ActiveStorageGBh > 0 {
			namespaces[i].ActiveStoragePercent = (namespaces[i].ActiveStorageGBh / totals.ActiveStorageGBh) * 100
		}
		if totals.RetainedStorageGBh > 0 {
			namespaces[i].RetainedStoragePercent = (namespaces[i].RetainedStorageGBh / totals.RetainedStorageGBh) * 100
		}
	}

	// Sort namespaces by name for consistent output
	sort.Slice(namespaces, func(i, j int) bool {
		return namespaces[i].Name < namespaces[j].Name
	})

	return &Report{
		Period: Period{
			Start: startDate,
			End:   endDate,
		},
		Pricing:    pricing,
		Namespaces: namespaces,
		Totals:     totals,
	}
}

type namespaceAggregator struct {
	actions                    float64
	activeStorageByteSeconds   float64
	retainedStorageByteSeconds float64
}

func extractNamespace(groupBys []models.GroupBy) string {
	for _, gb := range groupBys {
		if gb.Key == models.GroupByKeyNamespace {
			return gb.Value
		}
	}
	return ""
}

func calculateNamespaceUsage(name string, agg *namespaceAggregator, pricing Pricing) NamespaceUsage {
	// Convert byte-seconds to GBh:
	// GBh = byte_seconds / (3600 seconds/hour) / (1024^3 bytes/GB)
	const bytesPerGB = 1024.0 * 1024.0 * 1024.0
	const secondsPerHour = 3600.0

	activeStorageGBh := agg.activeStorageByteSeconds / secondsPerHour / bytesPerGB
	retainedStorageGBh := agg.retainedStorageByteSeconds / secondsPerHour / bytesPerGB

	// Calculate costs
	actionCost := (agg.actions / 1_000_000.0) * pricing.ActionPricePerMillion
	activeStorageCost := activeStorageGBh * pricing.ActiveStoragePricePerGBh
	retainedStorageCost := retainedStorageGBh * pricing.RetainedStoragePricePerGBh

	return NamespaceUsage{
		Name:                name,
		Actions:             agg.actions,
		ActiveStorageGBh:    activeStorageGBh,
		RetainedStorageGBh:  retainedStorageGBh,
		ActionCost:          actionCost,
		ActiveStorageCost:   activeStorageCost,
		RetainedStorageCost: retainedStorageCost,
		TotalCost:           actionCost + activeStorageCost + retainedStorageCost,
	}
}
