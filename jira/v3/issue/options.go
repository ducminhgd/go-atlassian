package issue

// IssueGetOpts represents options for getting an issue
type IssueGetOpts struct {
	// Expand options that include additional issue details in the response
	Expand []string
	
	// List of fields to return for the issue
	Fields []string
	
	// List of properties to return for the issue
	Properties []string
}

// IssueSearchOpts represents options for searching issues
type IssueSearchOpts struct {
	// Expand options that include additional issue details in the response
	Expand string
	
	// List of fields to return for each issue
	Fields []string
	
	// Whether to use field keys instead of field names in the response
	FieldsByKeys bool
	
	// Maximum number of results to return (default: 50, max: 100)
	MaxResults int
	
	// Token for pagination to get the next page of results
	NextPageToken string
	
	// List of properties to return for each issue
	Properties []string
	
	// List of issue IDs to reconcile
	ReconcileIssues []int
}

// IssueCreateOpts represents options for creating an issue
type IssueCreateOpts struct {
	// Whether to update the issue history
	UpdateHistory bool
}

// IssueUpdateOpts represents options for updating an issue
type IssueUpdateOpts struct {
	// Whether to notify users about the update
	NotifyUsers bool
	
	// Whether to override screen security
	OverrideScreenSecurity bool
	
	// Whether to override editable fields
	OverrideEditableFields bool
}

// IssueDeleteOpts represents options for deleting an issue
type IssueDeleteOpts struct {
	// Whether to delete subtasks
	DeleteSubtasks bool
}

// IssueTransitionOpts represents options for transitioning an issue
type IssueTransitionOpts struct {
	// Expand options for the transition response
	Expand []string
}

// IssueCommentOpts represents options for issue comments
type IssueCommentOpts struct {
	// Expand options for comments
	Expand []string
	
	// Maximum number of results to return
	MaxResults int
	
	// Starting index for pagination
	StartAt int
	
	// Order by field (created, updated)
	OrderBy string
}

// IssueWorklogOpts represents options for issue worklogs
type IssueWorklogOpts struct {
	// Expand options for worklogs
	Expand []string
	
	// Maximum number of results to return
	MaxResults int
	
	// Starting index for pagination
	StartAt int
}
