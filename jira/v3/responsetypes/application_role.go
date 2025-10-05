package responsetypes

// ApplicationRole represents an application role in Jira.
type ApplicationRole struct {
	// The groups that are granted default access for this application role. As a group's name can change, use of defaultGroupsDetails is recommended to identify a groups.
	DefaultGroups []string `json:"defaultGroups,omitempty"`

	// The groups that are granted default access for this application role.
	DefaultGroupsDetails []Group `json:"defaultGroupsDetails,omitempty"`

	// The groups associated with the application role.
	GroupDetails []Group `json:"groupDetails,omitempty"`

	// The groups associated with the application role. As a group's name can change, use of groupDetails is recommended to identify a groups.
	Groups []string `json:"groups,omitempty"`

	// Has unlimited seats
	HasUnlimitedSeats bool `json:"hasUnlimitedSeats,omitempty"`

	// The key of the application role.
	Key string `json:"key,omitempty"`

	// The display name of the application role.
	Name string `json:"name,omitempty"`

	// The maximum count of users on your license.
	NumberOfSeats int `json:"numberOfSeats,omitempty"`

	// Indicates if the application role belongs to Jira platform (jira-core).
	Platform bool `json:"platform,omitempty"`

	// The count of users remaining on your license.
	RemainingSeats int `json:"remainingSeats,omitempty"`

	// Determines whether this application role should be selected by default on user creation.
	SelectedByDefault bool `json:"selectedByDefault,omitempty"`

	// The number of users counting against your license.
	UserCount int `json:"userCount,omitempty"`

	// The type of users being counted against your license.
	UserCountDescription string `json:"userCountDescription,omitempty"`
}

type SimpleListWrapperApplicationRole struct {
	Items []ApplicationRole `json:"items,omitempty"`
	Size  int               `json:"size,omitempty"`
}
