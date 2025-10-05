package responsetypes

type ProjectVersion struct {
	// If the expand option approvers is used, returns a list containing the approvers for this version.
	Approvers []VersionApprover `json:"approvers,omitempty"`

	// Indicates that the version is archived. Optional when creating or updating a version.
	Archived bool `json:"archived,omitempty"`

	// The description of the version. Optional when creating or updating a version. The maximum size is 16,384 bytes.
	Description string `json:"description,omitempty"`

	// If the expand option driver is used, returns the Atlassian account ID of the driver.
	Driver string `json:"driver,omitempty"`

	// Use expand to include additional information about version in the response. This parameter accepts a comma-separated list. Expand options include:
	// - `operationss` Returns the list of operations available for this version.
	// - `issuesstatus` Returns the count of issues in this version for each of the status categories to do, in progress, done, and unmapped. The unmapped property contains a count of issues with a status other than to do, in progress, and done.
	// - `driver` Returns the Atlassian account ID of the version driver.
	// - `approvers` Returns a list containing approvers for this version.
	// Optional for create and update.
	Expand string `json:"expand,omitempty"`

	// The ID of the version.
	ID string `json:"id,omitempty"`

	// If the expand option issuesstatus is used, returns the count of issues in this version for each of the status categories to do, in progress, done, and unmapped. The unmapped property contains a count of issues with a status other than to do, in progress, and done.
	IssuesStatusForFixVersion map[string]int `json:"issuesStatusForFixVersion,omitempty"`

	// The URL of the self link to the version to which all unfixed issues are moved when a version is released. Not applicable when creating a version. Optional when updating a version.
	MoveUnfixedIssuesTo string `json:"moveUnfixedIssuesTo,omitempty"`

	// The unique name of the version. Required when creating a version. Optional when updating a version. The maximum length is 255 characters.
	Name string `json:"name,omitempty"`

	// If the expand option operations is used, returns the list of operations available for this version.
	Operations []SimpleLink `json:"operations,omitempty"`

	// Indicates that the version is overdue.
	Overdue bool `json:"overdue,omitempty"`

	// The ID of the project to which this version is attached. Required when creating a version. Not applicable when updating a version.
	ProjectID string `json:"projectId,omitempty"`

	// The release date of the version. Expressed in ISO 8601 format (yyyy-mm-dd). Optional when creating or updating a version.
	ReleaseDate string `json:"releaseDate,omitempty"`

	// Indicates that the version is released. If the version is released a request to release again is ignored. Not applicable when creating a version. Optional when updating a version.
	Released bool `json:"released,omitempty"`

	// The URL of the version.
	Self string `json:"self,omitempty"`

	// The start date of the version. Expressed in ISO 8601 format (yyyy-mm-dd). Optional when creating or updating a version.
	StartDate string `json:"startDate,omitempty"`

	// The date on which work on this version is expected to finish, expressed in the instance's Day/Month/Year Format date format.
	UserReleaseDate string `json:"userReleaseDate,omitempty"`

	// The date on which work on this version is expected to start, expressed in the instance's Day/Month/Year Format date format.
	UserStartDate string `json:"userStartDate,omitempty"`
}

type VersionApprover struct {
	// The Atlassian account ID of the approver.
	AccountID string `json:"accountId,omitempty"`

	// A description of why the user is declining the approval.
	DeclineReason string `json:"declineReason,omitempty"`

	// A description of what the user is approving within the specified version.
	Description string `json:"description,omitempty"`

	// The status of the approval, which can be PENDING, APPROVED, or DECLINED
	Status string `json:"status,omitempty"`
}

type SimpleLink struct {
	Href       string `json:"href,omitempty"`
	IconClass  string `json:"iconClass,omitempty"`
	ID         string `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	StyleClass string `json:"styleClass,omitempty"`
	Title      string `json:"title,omitempty"`
	Weight     int    `json:"weight,omitempty"`
}
