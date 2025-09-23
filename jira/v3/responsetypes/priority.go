package responsetypes

// Priority represents a priority level that can be assigned to an issue
type Priority struct {
	// The URL of the priority.
	Self string `json:"self,omitempty"`

	// A URL to an icon representing the priority.
	IconURL string `json:"iconUrl,omitempty"`

	// The name of the priority.
	Name string `json:"name,omitempty"`

	// The ID of the priority.
	ID string `json:"id,omitempty"`
}
