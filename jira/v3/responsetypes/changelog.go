package responsetypes

type PageOfChangelogs struct {
	Histories  []Changelog `json:"histories,omitempty"`
	MaxResults int         `json:"maxResults,omitempty"`
	StartAt    int         `json:"startAt,omitempty"`
	Total      int         `json:"total,omitempty"`
}

type Changelog struct {
	ID              string             `json:"id,omitempty"`
	Author          SimpleUser         `json:"author,omitempty"`
	Created         string             `json:"created,omitempty"`
	Items           []ChangelogDetails `json:"items,omitempty"`
	HistoryMetadata []HistoryMetadata  `json:"historyMetadata,omitempty"`
}

type ChangelogDetails struct {
	// The name of the field changed
	Field string `json:"field,omitempty"`

	// The ID of the field changed
	FieldID string `json:"fieldId,omitempty"`

	// The type of the field changed
	FieldType string `json:"fieldtype,omitempty"`

	// The details of the original value
	From string `json:"from,omitempty"`

	// The details of the original value as a string
	FromString string `json:"fromString,omitempty"`

	// The details of the new value
	To string `json:"to,omitempty"`

	// The details of the new value as a string
	ToString string `json:"toString,omitempty"`
}

// HistoryMetadata represents metadata about a history record
type HistoryMetadata struct {
	// The activity description
	ActivityDescription string `json:"activityDescription,omitempty"`

	// The key for the activity description
	ActivityDescriptionKey string `json:"activityDescriptionKey,omitempty"`

	// The actor who performed the activity
	Actor HistoryMetadataParticipant `json:"actor,omitempty"`

	// The cause of the activity
	Cause HistoryMetadataParticipant `json:"cause,omitempty"`

	// The description of the activity
	Description string `json:"description,omitempty"`

	// The key for the description
	DescriptionKey string `json:"descriptionKey,omitempty"`

	// The email description of the activity
	EmailDescription string `json:"emailDescription,omitempty"`

	// The key for the email description
	EmailDescriptionKey string `json:"emailDescriptionKey,omitempty"`

	// Additional data associated with the activity
	ExtraData map[string]interface{} `json:"extraData,omitempty"`

	// The generator of the activity
	Generator HistoryMetadataParticipant `json:"generator,omitempty"`

	// The type of the activity
	Type string `json:"type,omitempty"`
}

// HistoryMetadataParticipant represents a participant in a history metadata entry
type HistoryMetadataParticipant struct {
	// The avatar URL of the participant
	AvatarURL string `json:"avatarUrl,omitempty"`

	// The display name of the participant
	DisplayName string `json:"displayName,omitempty"`

	// The display name key of the participant
	DisplayNameKey string `json:"displayNameKey,omitempty"`

	// The ID of the participant
	ID string `json:"id,omitempty"`

	// The type of the participant
	Type string `json:"type,omitempty"`

	// The URL of the participant
	URL string `json:"url,omitempty"`
}
