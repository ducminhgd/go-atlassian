package responsetypes

// Project represents a Jira project. It contains all the details and metadata for a project.
type Project struct {
	// Whether the project has been archived
	Archived bool `json:"archived,omitempty"`

	// User archived the project
	ArchivedBy User `json:"archivedBy,omitempty"`

	// The date the project was archived
	ArchivedDate string `json:"archivedDate,omitempty"`

	// The default assignee when creating issues for this project. Valid values: PROJECT_LEAD, UNASSIGNED
	AssigneeType string `json:"assigneeType,omitempty"`

	// The URLs of the project's avatars.
	// The keys are the size of the avatar: "48x48", "32x32", "24x24", "16x16"
	AvatarURLs AvatarUrls `json:"avatarUrls,omitempty"`

	// List of the components contained in the project.
	Components []ProjectComponent `json:"components,omitempty"`

	// Whether the project is marked as deleted.
	Deleted bool `json:"deleted,omitempty"`

	// The user who marked the project as deleted.
	DeletedBy User `json:"deletedBy,omitempty"`

	// The date when the project was marked as deleted.
	DeletedDate string `json:"deletedDate,omitempty"`

	// The description of the project.
	Description string `json:"description,omitempty"`

	// An email address associated with the project.
	Email string `json:"email,omitempty"`

	// Expand options that include additional project details in the response.
	Expand string `json:"expand,omitempty"`

	// Whether the project is selected as a favorite.
	Favourite bool `json:"favourite,omitempty"`

	// The ID of the project
	ID string `json:"id,omitempty"`

	// Insights about the project.
	Insight ProjectInsight `json:"insight,omitempty"`

	// Whether the project is private from the user's perspective. This means the user can't see the project or any associated issues.
	IsPrivate bool `json:"isPrivate,omitempty"`

	// The issue type hierarchy for the project.
	IssueTypeHierarchy []Hierarchy `json:"issueTypeHierarchy,omitempty"`

	// List of the issue types available in the project.
	IssueTypes []IssueType `json:"issueTypes,omitempty"`

	// The key of the project. Must be unique and match Jira project key requirements
	Key string `json:"key"`

	// The project landing page info.
	LandingPageInfo ProjectLandingPageInfo `json:"landingPageInfo,omitempty"`

	// The username of the project lead.
	Lead User `json:"lead,omitempty"`

	// The name of the project.
	Name string `json:"name,omitempty"`

	// User permissions on the project
	Permissions ProjectPermissions `json:"permissions,omitempty"`

	// The category the project belongs to.
	ProjectCategory ProjectCategory `json:"projectCategory,omitempty"`

	// The project type of the project. Valid values: software, service_desk, business
	ProjectTypeKey string `json:"projectTypeKey,omitempty"`

	// Map of project properties
	Properties map[string]interface{} `json:"properties,omitempty"`

	// The date when the project is deleted permanently.
	RetentionTillDate string `json:"retentionTillDate,omitempty"`

	// Roles of the project.
	Roles map[string]string `json:"roles,omitempty"`

	// The URL of the project.
	Self string `json:"self,omitempty"`

	// Whether the project is simplified.
	Simplified bool `json:"simplified,omitempty"`

	// The type of the project.
	Style string `json:"style,omitempty"`

	// The URL of the project.
	URL string `json:"url,omitempty"`

	// Unique ID for next-gen projects.
	UUID string `json:"uuid,omitempty"`

	// The versions defined in the project
	Versions []ProjectVersion `json:"versions,omitempty"`
}

type SimpleProject struct {
	// The ID of the project
	ID string `json:"id,omitempty"`

	// The key of the project
	Key string `json:"key,omitempty"`

	// The name of the project
	Name string `json:"name,omitempty"`

	// The URL of the project
	Self string `json:"self,omitempty"`

	// The type of the project
	ProjectTypeKey string `json:"projectTypeKey,omitempty"`

	// Simplified project
	Simplified bool `json:"simplified,omitempty"`

	// Avatar URLs for the project
	AvatarURLs AvatarUrls `json:"avatarUrls,omitempty"`
}
