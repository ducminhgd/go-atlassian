package issue

// JQLSearchRequest represents the request body for JQL search
type JQLSearchRequest struct {
	// Expand options that include additional issue details in the response
	Expand string `json:"expand,omitempty"`

	// List of fields to return for each issue
	Fields []string `json:"fields,omitempty"`

	// Whether to use field keys instead of field names in the response
	FieldsByKeys bool `json:"fieldsByKeys,omitempty"`

	// JQL query string
	JQL string `json:"jql"`

	// Maximum number of results to return (default: 50, max: 100)
	MaxResults int `json:"maxResults,omitempty"`

	// Token for pagination to get the next page of results
	NextPageToken string `json:"nextPageToken,omitempty"`

	// List of properties to return for each issue
	Properties []string `json:"properties,omitempty"`

	// List of issue IDs to reconcile
	ReconcileIssues []int `json:"reconcileIssues,omitempty"`
}

// JQLSearchResponse represents the response from JQL search
type JQLSearchResponse struct {
	// Whether this is the last page of results
	IsLast bool `json:"isLast"`

	// List of issues returned by the search
	Issues []Issue `json:"issues"`

	// Maximum number of results that can be returned
	MaxResults int `json:"maxResults,omitempty"`

	// Token for getting the next page of results
	NextPageToken string `json:"nextPageToken,omitempty"`

	// Starting index of the results
	StartAt int `json:"startAt,omitempty"`

	// Total number of issues matching the query
	Total int `json:"total,omitempty"`
}

// Issue represents a Jira issue (re-exported from responsetypes for convenience)
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

	// Changelog information (when expanded)
	Changelog PageOfChangelogs `json:"changelog,omitempty"`
}

// IssueFields represents the fields of an issue
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

	// Description of the issue
	Description interface{} `json:"description,omitempty"`

	// Voting information
	Votes VoteInfo `json:"votes,omitempty"`

	// Watch information
	Watches WatchInfo `json:"watches,omitempty"`

	// Worklog information
	Worklog PagedWorklog `json:"worklog,omitempty"`

	// Comment information
	Comment PagedComment `json:"comment,omitempty"`

	// Created timestamp
	Created string `json:"created,omitempty"`

	// Updated timestamp
	Updated string `json:"updated,omitempty"`
}

// SimpleUser represents a basic user structure
type SimpleUser struct {
	AccountID    string            `json:"accountId,omitempty"`
	AccountType  string            `json:"accountType,omitempty"`
	Active       bool              `json:"active,omitempty"`
	AvatarUrls   map[string]string `json:"avatarUrls,omitempty"`
	DisplayName  string            `json:"displayName,omitempty"`
	EmailAddress string            `json:"emailAddress,omitempty"`
	Key          string            `json:"key,omitempty"`
	Name         string            `json:"name,omitempty"`
	Self         string            `json:"self,omitempty"`
	TimeZone     string            `json:"timeZone,omitempty"`
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

// IssueType represents a type of issue that can be created in Jira
type IssueType struct {
	// The ID of the issue type
	ID string `json:"id,omitempty"`

	// The name of the issue type (e.g., "Bug", "Story", "Task")
	Name string `json:"name,omitempty"`

	// A description of the issue type
	Description string `json:"description,omitempty"`

	// The URL of the issue type's icon
	IconURL string `json:"iconUrl,omitempty"`

	// The ID of the issue type's avatar
	AvatarID int64 `json:"avatarId,omitempty"`

	// The hierarchy level of the issue type (e.g., 0 for regular issues, -1 for subtasks)
	HierarchyLevel int `json:"hierarchyLevel,omitempty"`

	// Whether this issue type is a subtask type
	Subtask bool `json:"subtask,omitempty"`
}

// Priority represents the priority of an issue
type Priority struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"iconUrl,omitempty"`
	StatusColor string `json:"statusColor,omitempty"`
	Self        string `json:"self,omitempty"`
}

// SimpleProject represents a basic project structure
type SimpleProject struct {
	ID             string            `json:"id,omitempty"`
	Key            string            `json:"key,omitempty"`
	Name           string            `json:"name,omitempty"`
	ProjectTypeKey string            `json:"projectTypeKey,omitempty"`
	Self           string            `json:"self,omitempty"`
	AvatarUrls     map[string]string `json:"avatarUrls,omitempty"`
}

// StatusDetails represents the status of an issue
type StatusDetails struct {
	Self           string         `json:"self"`
	Description    string         `json:"description"`
	IconURL        string         `json:"iconUrl"`
	Name           string         `json:"name"`
	ID             string         `json:"id"`
	StatusCategory StatusCategory `json:"statusCategory"`
}

// StatusCategory represents the category of a status
type StatusCategory struct {
	Self      string `json:"self"`
	ID        int    `json:"id"`
	Key       string `json:"key"`
	ColorName string `json:"colorName"`
	Name      string `json:"name"`
}

// VoteInfo represents voting information for an issue
type VoteInfo struct {
	// Whether the current user has voted
	HasVoted bool `json:"hasVoted"`
	// Number of votes
	Votes int `json:"votes,omitempty"`
	// URL of the votes resource
	Self string `json:"self,omitempty"`
}

// WatchInfo represents watch information for an issue
type WatchInfo struct {
	// Whether the current user is watching
	IsWatching bool `json:"isWatching"`
	// Number of watchers
	WatchCount int `json:"watchCount,omitempty"`
	// URL of the watches resource
	Self string `json:"self,omitempty"`
}

// PagedWorklog represents worklog information with pagination
type PagedWorklog struct {
	StartAt    int       `json:"startAt,omitempty"`
	MaxResults int       `json:"maxResults,omitempty"`
	Total      int       `json:"total,omitempty"`
	Worklogs   []Worklog `json:"worklogs,omitempty"`
}

// Worklog represents a worklog entry
type Worklog struct {
	Self             string      `json:"self,omitempty"`
	Author           SimpleUser  `json:"author,omitempty"`
	UpdateAuthor     SimpleUser  `json:"updateAuthor,omitempty"`
	Comment          interface{} `json:"comment,omitempty"`
	Created          string      `json:"created,omitempty"`
	Updated          string      `json:"updated,omitempty"`
	Started          string      `json:"started,omitempty"`
	TimeSpent        string      `json:"timeSpent,omitempty"`
	TimeSpentSeconds int         `json:"timeSpentSeconds,omitempty"`
	ID               string      `json:"id,omitempty"`
	IssueID          string      `json:"issueId,omitempty"`
}

// IssueComment represents a single comment on a Jira issue
type IssueComment struct {
	Self         string      `json:"self"`
	ID           string      `json:"id"`
	Author       SimpleUser  `json:"author"`
	Body         interface{} `json:"body"`
	UpdateAuthor SimpleUser  `json:"updateAuthor"`
	Created      string      `json:"created"`
	Updated      string      `json:"updated"`
	JsdPublic    bool        `json:"jsdPublic"`
}

// PagedComment represents a paged list of comments with pagination information
type PagedComment struct {
	StartAt    int            `json:"startAt"`
	MaxResults int            `json:"maxResults"`
	Total      int            `json:"total"`
	Comments   []IssueComment `json:"comments"`
}

// PageOfChangelogs represents a paged list of changelogs
type PageOfChangelogs struct {
	Histories  []Changelog `json:"histories,omitempty"`
	MaxResults int         `json:"maxResults,omitempty"`
	StartAt    int         `json:"startAt,omitempty"`
	Total      int         `json:"total,omitempty"`
}

// Changelog represents a single changelog entry
type Changelog struct {
	ID      string             `json:"id,omitempty"`
	Author  SimpleUser         `json:"author,omitempty"`
	Created string             `json:"created,omitempty"`
	Items   []ChangelogDetails `json:"items,omitempty"`
}

// ChangelogDetails represents the details of a changelog item
type ChangelogDetails struct {
	// The name of the field changed
	Field string `json:"field,omitempty"`

	// The ID of the field changed
	FieldID string `json:"fieldId,omitempty"`

	// The type of the field changed
	FieldType string `json:"fieldtype,omitempty"`

	// The details of the original value
	From string `json:"from,omitempty"`

	// The details of the original value as a string
	FromString string `json:"fromString,omitempty"`

	// The details of the new value
	To string `json:"to,omitempty"`

	// The details of the new value as a string
	ToString string `json:"toString,omitempty"`
}
