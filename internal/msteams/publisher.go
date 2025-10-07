package msteams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Publisher handles publishing AdaptiveCards to Microsoft Teams webhooks
type Publisher struct {
	webhookURL string
}

// NewPublisher creates a new Teams publisher
func NewPublisher(webhookURL string) *Publisher {
	return &Publisher{webhookURL: webhookURL}
}

// PublishAdaptiveCard posts an AdaptiveCard to the Teams webhook
func (p *Publisher) PublishAdaptiveCard(adaptiveCard AdaptiveCard) error {
	message := FormatTeamsMessage(adaptiveCard)

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal Teams message: %w", err)
	}

	resp, err := http.Post(p.webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to post to Teams webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Teams webhook request failed with status: %d", resp.StatusCode)
	}

	return nil
}
