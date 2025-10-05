package responsetypes

// User represents a Jira user with their account details and metadata
type User struct {
	// The account ID of the user, which uniquely identifies the user across all Atlassian products. For example, 5b10ac8d82e05b22cc7d4ef5. Required in requests.
	AccountID string `json:"accountId"`

	// The user account type. Can take the following values:
	// - `atlassian` regular Atlassian user account
	// - `app` system account used for Connect applications and OAuth to represent external systems
	// - `customer` Jira Service Desk account representing an external service desk
	// - `unknown` account type
	AccountType string `json:"accountType,omitempty"`

	// Whether the user account is active
	Active bool `json:"active"`

	// The application roles the user is assigned to.
	ApplicationRoles SimpleListWrapperApplicationRole `json:"applicationRoles,omitempty"`

	// A map of avatar URLs for different sizes (16x16, 24x24, 32x32, 48x48)
	AvatarURLs AvatarUrls `json:"avatarUrls,omitempty"`

	// The display name of the user. Depending on the user’s privacy setting, this may return an alternative value.
	DisplayName string `json:"displayName"`

	// The email address of the user. Depending on the user’s privacy setting, this may be returned as null.
	EmailAddress string `json:"emailAddress,omitempty"`

	// Expand options that include additional user details in the response.
	Expand string `json:"expand,omitempty"`

	// The groups that the user belongs to.
	Groups SimpleListWrapperGroup `json:"groups,omitempty"`

	// The locale of the user. Depending on the user’s privacy setting, this may be returned as null.
	Locale string `json:"locale,omitempty"`

	// The URL of the user.
	Self string `json:"self"`

	// The time zone specified in the user's profile. If the user's time zone is not visible to the current user (due to user's profile setting), or if a time zone has not been set, the instance's default time zone will be returned.
	TimeZone string `json:"timeZone,omitempty"`
}

type SimpleUser struct {
	// The account ID of the user
	AccountID string `json:"accountId,omitempty"`

	// The user account type
	AccountType string `json:"accountType,omitempty"`

	// Whether the user account is active
	Active bool `json:"active,omitempty"`

	// A map of avatar URLs
	AvatarUrls AvatarUrls `json:"avatarUrls,omitempty"`

	// The display name of the user
	DisplayName string `json:"displayName,omitempty"`

	// The email address of the user
	EmailAddress string `json:"emailAddress,omitempty"`

	// The key of the user
	Key string `json:"key,omitempty"`

	// The name of the user
	Name string `json:"name,omitempty"`

	// The URL of the user
	Self string `json:"self,omitempty"`

	// The time zone of the user
	TimeZone string `json:"timeZone,omitempty"`
}
