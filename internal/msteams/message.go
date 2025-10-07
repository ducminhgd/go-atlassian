package msteams

// TeamsMessage represents a message to be sent to Microsoft Teams
type TeamsMessage struct {
	Type        string       `json:"type"`
	Attachments []Attachment `json:"attachments"`
}

// Attachment represents an attachment in a Teams message
type Attachment struct {
	ContentType string        `json:"contentType"`
	Content     *AdaptiveCard `json:"content"`
}

// FormatTeamsMessage creates a Teams message with an AdaptiveCard
func FormatTeamsMessage(adaptiveCard AdaptiveCard) TeamsMessage {
	return TeamsMessage{
		Type: "message",
		Attachments: []Attachment{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content:     &adaptiveCard,
			},
		},
	}
}
