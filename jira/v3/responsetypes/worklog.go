package responsetypes

// Worklog represents a single worklog entry for an issue
type Worklog struct {
	// The URL of the worklog entry
	Self string `json:"self,omitempty"`
	// The author of the worklog
	Author SimpleUser `json:"author,omitempty"`
	// The user who last updated the worklog
	UpdateAuthor SimpleUser `json:"updateAuthor,omitempty"`
	// The comment about the work in Atlassian Document Format
	Comment AtlassianDocumentFormat `json:"comment,omitempty"`
	// When the worklog was created
	Created string `json:"created,omitempty"`
	// When the worklog was last updated
	Updated string `json:"updated,omitempty"`
	// When the work was started
	Started string `json:"started,omitempty"`
	// The time spent working (e.g., "15m", "2h 30m")
	TimeSpent string `json:"timeSpent,omitempty"`
	// The time spent in seconds
	TimeSpentSeconds int `json:"timeSpentSeconds,omitempty"`
	// The ID of the worklog entry
	ID string `json:"id,omitempty"`
	// The ID of the issue this worklog belongs to
	IssueID string `json:"issueId,omitempty"`
}

// PagedWorklog represents the collection of worklog entries for an issue
type PagedWorklog struct {
	// List of worklog entries
	Worklogs []Worklog `json:"worklogs,omitempty"`
	// Maximum results returned
	MaxResults int `json:"maxResults,omitempty"`
	// Start index of the page
	StartAt int `json:"startAt,omitempty"`
	// Total number of worklogs
	Total int `json:"total,omitempty"`
}
