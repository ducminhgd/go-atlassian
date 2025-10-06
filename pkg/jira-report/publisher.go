package jirareport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Publisher handles publishing reports to webhooks
type Publisher struct {
	webhookURL string
	client     *http.Client
}

// NewPublisher creates a new publisher
func NewPublisher(webhookURL string) *Publisher {
	return &Publisher{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Publish posts the HTML report to the webhook
func (p *Publisher) Publish(htmlReport string) error {
	// Microsoft Teams Workflow can accept HTML content in the body field
	// Send as a simple message with HTML body
	payload := map[string]interface{}{
		"body": htmlReport,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("%w: failed to marshal payload: %v", ErrPostToWebhook, err)
	}

	resp, err := p.client.Post(p.webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPostToWebhook, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("%w: webhook returned status %d", ErrPostToWebhook, resp.StatusCode)
	}

	return nil
}

