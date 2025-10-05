package responsetypes

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

	// The scope of the issue type (project or global)
	Scope map[string]interface{} `json:"scope,omitempty"`

	// The REST API URL for this issue type resource
	Self string `json:"self,omitempty"`
}
