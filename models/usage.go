package models

// GetUsageRequest represents a request to the Temporal Cloud usage API.
type GetUsageRequest struct {
	StartTimeInclusive string `json:"startTimeInclusive,omitempty"`
	EndTimeExclusive   string `json:"endTimeExclusive,omitempty"`
	PageSize           int    `json:"pageSize,omitempty"`
	PageToken          string `json:"pageToken,omitempty"`
}

// GetUsageResponse represents the response from the Temporal Cloud usage API.
type GetUsageResponse struct {
	Summaries     []Summary `json:"summaries"`
	NextPageToken string    `json:"nextPageToken"`
}

// Summary contains usage data for a specific time period.
type Summary struct {
	StartTime    string        `json:"startTime"`
	EndTime      string        `json:"endTime"`
	RecordGroups []RecordGroup `json:"recordGroups"`
	Incomplete   bool          `json:"incomplete"`
}

// RecordGroup organizes records by grouping criteria.
type RecordGroup struct {
	GroupBys []GroupBy `json:"groupBys"`
	Records  []Record  `json:"records"`
}

// GroupBy represents a key-value pair for record grouping.
type GroupBy struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Record represents an individual usage measurement.
type Record struct {
	Type  string  `json:"type"`
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

// RecordType constants for usage record types.
const (
	RecordTypeActions         = "RECORD_TYPE_ACTIONS"
	RecordTypeActiveStorage   = "RECORD_TYPE_ACTIVE_STORAGE"
	RecordTypeRetainedStorage = "RECORD_TYPE_RETAINED_STORAGE"
)

// GroupByKey constants for grouping dimensions.
const (
	GroupByKeyNamespace = "GROUP_BY_KEY_NAMESPACE"
)
