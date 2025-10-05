package project

// ProjectUpdateOpts contains the options for updating a project
type ProjectUpdateOpts struct {
	// The project category ID associated with the project
	CategoryID int64 `json:"categoryId,omitempty"`

	// The description of the project
	Description string `json:"description,omitempty"`

	// The issue security scheme ID associated with the project
	IssueSecurityScheme int64 `json:"issueSecurityScheme,omitempty"`

	// The project key. Must be unique and match Jira project key requirements
	Key string `json:"key,omitempty"`

	// The account ID of the project lead
	LeadAccountID string `json:"leadAccountId,omitempty"`

	// The name of the project
	Name string `json:"name,omitempty"`

	// The notification scheme ID associated with the project
	NotificationScheme int64 `json:"notificationScheme,omitempty"`

	// The permission scheme ID associated with the project
	PermissionScheme int64 `json:"permissionScheme,omitempty"`

	// The URL of the project
	URL string `json:"url,omitempty"`

	// Use expand to include additional information in the response. This parameter accepts a comma-separated list.
	// Expanded options include: "description", "issueTypes", "lead", "projectKeys"
	Expand string `url:"expand,omitempty"`
}

// ProjectQueryOpts contains the options for the Get method
type ProjectQueryOpts struct {
	// Use expand to include additional information in the response. This parameter accepts a comma-separated list.
	// Expanded options include:
	// "description" - Returns the project description
	// "issueTypes" - Returns all issue types associated with the project
	// "lead" - Returns information about the project lead
	// "projectKeys" - Returns all project keys associated with the project
	Expand string `url:"expand,omitempty"`

	// A list of project properties to return for the project. This parameter accepts a comma-separated list.
	Properties []string `url:"properties,omitempty"`
}

// ProjectGetAllOpts contains the options for the GetAll method
type ProjectGetAllOpts struct {
	// Use expand to include additional information in the response. This parameter accepts a comma-separated list.
	// Expanded options include:
	// "description" - Returns the project description
	// "issueTypes" - Returns all issue types associated with the project
	// "lead" - Returns information about the project lead
	// "projectKeys" - Returns all project keys associated with the project
	Expand string `url:"expand,omitempty"`

	// A list of project properties to return for the project. This parameter accepts a comma-separated list.
	Properties []string `url:"properties,omitempty"`

	// Returns the user's most recently accessed projects. You may specify the number of results to return up to a maximum of 20.
	// If access is anonymous, then the recently accessed projects are based on the current HTTP session.
	Recent int `url:"recent,omitempty"`
}
