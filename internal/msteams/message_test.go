package msteams

import (
	"encoding/json"
	"testing"
)

func TestFormatTeamsMessage(t *testing.T) {
	// Create a simple AdaptiveCard
	card := NewAdaptiveCard()
	card.AddTextBlock("Test message", "", "", true)

	// Format as Teams message
	message := FormatTeamsMessage(card)

	// Verify structure
	if message.Type != "message" {
		t.Errorf("Expected Type to be 'message', got %s", message.Type)
	}

	if len(message.Attachments) != 1 {
		t.Errorf("Expected 1 attachment, got %d", len(message.Attachments))
	}

	attachment := message.Attachments[0]
	if attachment.ContentType != "application/vnd.microsoft.card.adaptive" {
		t.Errorf("Expected ContentType to be 'application/vnd.microsoft.card.adaptive', got %s", attachment.ContentType)
	}

	if attachment.Content.Type != "AdaptiveCard" {
		t.Errorf("Expected Content.Type to be 'AdaptiveCard', got %s", attachment.Content.Type)
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Errorf("Failed to marshal TeamsMessage to JSON: %v", err)
	}

	// Verify JSON is valid
	var unmarshaledMessage TeamsMessage
	err = json.Unmarshal(jsonData, &unmarshaledMessage)
	if err != nil {
		t.Errorf("Failed to unmarshal TeamsMessage JSON: %v", err)
	}
}

func TestTeamsMessage_JSONStructure(t *testing.T) {
	card := NewAdaptiveCard()
	card.AddTextBlock("Test", "Large", "Bolder", true)

	message := FormatTeamsMessage(card)

	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Errorf("Failed to marshal TeamsMessage: %v", err)
	}

	// Parse JSON to verify structure
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	// Verify top-level structure
	if jsonMap["type"] != "message" {
		t.Errorf("Expected type to be 'message', got %v", jsonMap["type"])
	}

	attachments, ok := jsonMap["attachments"].([]interface{})
	if !ok || len(attachments) != 1 {
		t.Error("Expected attachments to be an array with 1 element")
	}

	attachment, ok := attachments[0].(map[string]interface{})
	if !ok {
		t.Error("Expected attachment to be an object")
	}

	if attachment["contentType"] != "application/vnd.microsoft.card.adaptive" {
		t.Errorf("Expected contentType to be 'application/vnd.microsoft.card.adaptive', got %v", attachment["contentType"])
	}

	content, ok := attachment["content"].(map[string]interface{})
	if !ok {
		t.Error("Expected content to be an object")
	}

	if content["type"] != "AdaptiveCard" {
		t.Errorf("Expected content type to be 'AdaptiveCard', got %v", content["type"])
	}
}
