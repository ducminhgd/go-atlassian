package responsetypes

// ProjectInsight represents project statistics and metrics
type ProjectInsight struct {
	// The timestamp when the last issue was updated in the project (Unix timestamp in milliseconds)
	LastIssueUpdateTime int64 `json:"lastIssueUpdateTime,omitempty"`

	// The total number of issues in the project
	TotalIssueCount int `json:"totalIssueCount,omitempty"`
}
