package jirareport

import (
	"github.com/ducminhgd/go-atlassian/internal/msteams"
)

// Publisher handles publishing reports to webhooks
type Publisher struct {
	teamsPublisher *msteams.Publisher
}

// NewPublisher creates a new publisher
func NewPublisher(webhookURL string) *Publisher {
	return &Publisher{
		teamsPublisher: msteams.NewPublisher(webhookURL),
	}
}

// PublishAdaptiveCard posts an AdaptiveCard report to the webhook
func (p *Publisher) PublishAdaptiveCard(adaptiveCard msteams.AdaptiveCard) error {
	return p.teamsPublisher.PublishAdaptiveCard(adaptiveCard)
}
