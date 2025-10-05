package responsetypes

type VoteInfo struct {
	// Whether the current user has voted
	HasVoted bool `json:"hasVoted"`
	// Number of votes
	Votes int `json:"votes,omitempty"`
	// URL of the votes resource
	Self string `json:"self,omitempty"`
}

type WatchInfo struct {
	// Whether the current user is watching
	IsWatching bool `json:"isWatching"`
	// Number of watchers
	WatchCount int `json:"watchCount,omitempty"`
	// URL of the watches resource
	Self string `json:"self,omitempty"`
}

// ParentIssue represents a minimal version of an issue used when an issue is referenced as a parent
type ParentIssue struct {
	// The ID of the parent issue
	ID string `json:"id,omitempty"`

	// The key of the parent issue
	Key string `json:"key,omitempty"`

	// The self URL of the parent issue
	Self string `json:"self,omitempty"`

	// The fields of the parent issue
	Fields struct {
		// The summary of the parent issue
		Summary string `json:"summary,omitempty"`

		// The status of the parent issue
		Status StatusDetails `json:"status,omitempty"`

		// The priority of the parent issue
		Priority Priority `json:"priority,omitempty"`

		// The issue type of the parent issue
		IssueType IssueType `json:"issuetype,omitempty"`
	} `json:"fields,omitempty"`
}

type IssueFields struct {
	// Person to whom the issue is assigned
	Assignee SimpleUser `json:"assignee,omitempty"`

	// Person who created the issue
	Creator SimpleUser `json:"creator,omitempty"`

	// Parent issue if this issue is a subtask
	Parent ParentIssue `json:"parent,omitempty"`

	// Type of the issue
	IssueType IssueType `json:"issuetype,omitempty"`

	// Labels attached to the issue
	Labels []string `json:"labels,omitempty"`

	// Priority of the issue
	Priority Priority `json:"priority,omitempty"`

	// Project the issue belongs to
	Project SimpleProject `json:"project,omitempty"`

	// Person who reported the issue
	Reporter SimpleUser `json:"reporter,omitempty"`

	// Resolution details of the issue
	Resolution *StatusDetails `json:"resolution,omitempty"`

	// Current status of the issue
	Status StatusDetails `json:"status,omitempty"`

	// Status category of the issue
	StatusCategory StatusCategory `json:"statusCategory,omitempty"`

	// Summary of the issue
	Summary string `json:"summary,omitempty"`

	// Voting information
	Votes VoteInfo `json:"votes,omitempty"`

	// Watch information
	Watches WatchInfo `json:"watches,omitempty"`

	// Worklog information
	Worklog PagedWorklog `json:"worklog,omitempty"`
}

type Issue struct {
	// Expand options that include additional issue details in the response
	Expand string `json:"expand,omitempty"`

	// All fields of the issue
	Fields IssueFields `json:"fields,omitempty"`

	// The ID of the issue
	ID string `json:"id,omitempty"`

	// The key of the issue
	Key string `json:"key,omitempty"`

	// The self URL of the issue
	Self string `json:"self,omitempty"`
}
