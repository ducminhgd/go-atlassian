package responsetypes

// ProjectComponent represents a component within a project, used to organize and categorize issues
type ProjectComponent struct {
	// The unique identifier of the component
	ID string `json:"id"`

	// The name of the component
	Name string `json:"name"`

	// A description of the component
	Description string `json:"description,omitempty"`

	// The user assigned as the component's lead
	Lead *User `json:"lead,omitempty"`

	// The key of the project this component belongs to
	Project string `json:"project,omitempty"`

	// The numeric ID of the project this component belongs to
	ProjectID int64 `json:"projectId,omitempty"`

	// The REST API URL for this component resource
	Self string `json:"self,omitempty"`
}
