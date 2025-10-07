package msteams

import (
	"encoding/json"
	"testing"
)

func TestNewAdaptiveCard(t *testing.T) {
	card := NewAdaptiveCard()

	if card.Type != "AdaptiveCard" {
		t.Errorf("Expected Type to be 'AdaptiveCard', got %s", card.Type)
	}

	if card.Version != "1.5" {
		t.Errorf("Expected Version to be '1.5', got %s", card.Version)
	}

	if card.MSTeams == nil || card.MSTeams.Width != "Full" {
		t.Error("Expected MSTeams.Width to be 'Full'")
	}

	if len(card.Body) != 0 {
		t.Errorf("Expected empty Body, got %d elements", len(card.Body))
	}
}

func TestAdaptiveCard_AddTextBlock(t *testing.T) {
	card := NewAdaptiveCard()
	
	card.AddTextBlock("Test Text", "Large", "Bolder", true)

	if len(card.Body) != 1 {
		t.Errorf("Expected 1 element in Body, got %d", len(card.Body))
	}

	element := card.Body[0]
	if element.Type != "TextBlock" {
		t.Errorf("Expected Type to be 'TextBlock', got %s", element.Type)
	}

	if element.Text != "Test Text" {
		t.Errorf("Expected Text to be 'Test Text', got %s", element.Text)
	}

	if element.Size != "Large" {
		t.Errorf("Expected Size to be 'Large', got %s", element.Size)
	}

	if element.Weight != "Bolder" {
		t.Errorf("Expected Weight to be 'Bolder', got %s", element.Weight)
	}

	if !element.Wrap {
		t.Error("Expected Wrap to be true")
	}
}

func TestAdaptiveCard_AddRichTextBlock(t *testing.T) {
	card := NewAdaptiveCard()
	
	inlines := []AdaptiveCardInline{
		CreateTextRun("Bold Text", "Bolder", "Accent"),
	}
	
	card.AddRichTextBlock(inlines, "Medium", true)

	if len(card.Body) != 1 {
		t.Errorf("Expected 1 element in Body, got %d", len(card.Body))
	}

	element := card.Body[0]
	if element.Type != "RichTextBlock" {
		t.Errorf("Expected Type to be 'RichTextBlock', got %s", element.Type)
	}

	if element.Spacing != "Medium" {
		t.Errorf("Expected Spacing to be 'Medium', got %s", element.Spacing)
	}

	if !element.Separator {
		t.Error("Expected Separator to be true")
	}

	if len(element.Inlines) != 1 {
		t.Errorf("Expected 1 inline element, got %d", len(element.Inlines))
	}
}

func TestCreateTextRun(t *testing.T) {
	textRun := CreateTextRun("Test Text", "Bolder", "Accent")

	if textRun.Type != "TextRun" {
		t.Errorf("Expected Type to be 'TextRun', got %s", textRun.Type)
	}

	if textRun.Text != "Test Text" {
		t.Errorf("Expected Text to be 'Test Text', got %s", textRun.Text)
	}

	if textRun.Weight != "Bolder" {
		t.Errorf("Expected Weight to be 'Bolder', got %s", textRun.Weight)
	}

	if textRun.Color != "Accent" {
		t.Errorf("Expected Color to be 'Accent', got %s", textRun.Color)
	}
}

func TestCreateOpenUrlAction(t *testing.T) {
	action := CreateOpenUrlAction("Open Link", "https://example.com")

	if action.Type != "Action.OpenUrl" {
		t.Errorf("Expected Type to be 'Action.OpenUrl', got %s", action.Type)
	}

	if action.Title != "Open Link" {
		t.Errorf("Expected Title to be 'Open Link', got %s", action.Title)
	}

	if action.URL != "https://example.com" {
		t.Errorf("Expected URL to be 'https://example.com', got %s", action.URL)
	}
}

func TestAdaptiveCard_JSONMarshaling(t *testing.T) {
	card := NewAdaptiveCard()
	card.AddTextBlock("Test", "Large", "Bolder", true)

	jsonData, err := json.Marshal(card)
	if err != nil {
		t.Errorf("Failed to marshal AdaptiveCard to JSON: %v", err)
	}

	// Verify JSON is valid by unmarshaling
	var unmarshaledCard AdaptiveCard
	err = json.Unmarshal(jsonData, &unmarshaledCard)
	if err != nil {
		t.Errorf("Failed to unmarshal AdaptiveCard JSON: %v", err)
	}

	// Verify unmarshaled data matches original
	if unmarshaledCard.Type != card.Type {
		t.Errorf("Unmarshaled Type doesn't match: expected %s, got %s", card.Type, unmarshaledCard.Type)
	}

	if len(unmarshaledCard.Body) != len(card.Body) {
		t.Errorf("Unmarshaled Body length doesn't match: expected %d, got %d", len(card.Body), len(unmarshaledCard.Body))
	}
}
