package responsetypes

// ProjectListResponse represents a paginated list of projects
type ProjectListResponse struct {
	// The REST API URL for this resource
	Self string `json:"self,omitempty"`

	// The URL for the next page of results
	NextPage string `json:"nextPage,omitempty"`

	// The maximum number of results per page
	MaxResults int `json:"maxResults,omitempty"`

	// The index of the first item returned in the page
	StartAt int `json:"startAt,omitempty"`

	// The total number of items available
	Total int `json:"total,omitempty"`

	// Whether this is the last page of results
	IsLast bool `json:"isLast,omitempty"`

	// The list of projects in this page
	Values []Project `json:"values,omitempty"`
}
