package msteams

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestFormatJiraReportAsAdaptiveCard(t *testing.T) {
	// Create test data
	reportDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	epicGroups := map[string]*EpicGroup{
		"EPIC-123": {
			EpicKey:     "EPIC-123",
			EpicSummary: "Test Epic",
			EpicStatus:  "In Progress",
			EpicURL:     "https://jira.example.com/browse/EPIC-123",
			Issues: []IssueUpdate{
				{
					Key:         "TASK-456",
					Summary:     "Test Task",
					Status:      "In Progress",
					IssueType:   "Task",
					URL:         "https://jira.example.com/browse/TASK-456",
					LastUpdated: reportDate.Add(-2 * time.Hour),
					Updates: []Update{
						{
							Time:       reportDate.Add(-1 * time.Hour),
							AuthorName: "John Doe",
							Type:       "comment",
							Content:    "This is a test comment",
						},
						{
							Time:       reportDate.Add(-30 * time.Minute),
							AuthorName: "Jane Smith",
							Type:       "worklog",
							Content:    "Worked on implementation",
							TimeSpent:  "2h",
						},
					},
				},
			},
		},
	}

	noEpicIssues := []IssueUpdate{
		{
			Key:         "BUG-789",
			Summary:     "Test Bug",
			Status:      "Open",
			IssueType:   "Bug",
			URL:         "https://jira.example.com/browse/BUG-789",
			LastUpdated: reportDate.Add(-3 * time.Hour),
			Updates: []Update{
				{
					Time:       reportDate.Add(-2 * time.Hour),
					AuthorName: "Bob Wilson",
					Type:       "comment",
					Content:    "Found a bug in the system",
				},
			},
		},
	}

	// Generate AdaptiveCard
	card := FormatJiraReportAsAdaptiveCard(epicGroups, noEpicIssues, reportDate, "UTC")

	// Verify basic structure
	if card.Type != "AdaptiveCard" {
		t.Errorf("Expected Type to be 'AdaptiveCard', got %s", card.Type)
	}

	if card.Version != "1.5" {
		t.Errorf("Expected Version to be '1.5', got %s", card.Version)
	}

	if card.MSTeams == nil || card.MSTeams.Width != "Full" {
		t.Error("Expected MSTeams.Width to be 'Full'")
	}

	// Verify body has content
	if len(card.Body) == 0 {
		t.Error("Expected Body to have content")
	}

	// Verify title is present
	if len(card.Body) > 0 {
		titleBlock := card.Body[0]
		if titleBlock.Type != "TextBlock" {
			t.Errorf("Expected first element to be TextBlock, got %s", titleBlock.Type)
		}
		if titleBlock.Text != "Daily Report 15-Jan-2024" {
			t.Errorf("Expected title to be 'Daily Report 15-Jan-2024', got %s", titleBlock.Text)
		}
	}

	// Test JSON marshaling
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
}

func TestTruncateText(t *testing.T) {
	tests := []struct {
		input     string
		maxLength int
		expected  string
	}{
		{"Short text", 20, "Short text"},
		{"This is a very long text that should be truncated", 20, "This is a very lo..."},
		{"Exact length", 12, "Exact length"},
		{"", 10, ""},
		{"ABC", 3, "ABC"},
		{"ABCD", 3, "..."},
	}

	for _, test := range tests {
		result := truncateText(test.input, test.maxLength)
		if result != test.expected {
			t.Errorf("truncateText(%q, %d) = %q, expected %q", test.input, test.maxLength, result, test.expected)
		}
	}
}

func TestAddEpicSection(t *testing.T) {
	card := NewAdaptiveCard()

	group := &EpicGroup{
		EpicKey:     "EPIC-123",
		EpicSummary: "Test Epic",
		EpicStatus:  "In Progress",
		EpicURL:     "https://example.com/epic",
		Issues:      []IssueUpdate{},
	}

	addEpicSection(&card, group, time.UTC)

	// Should have at least the title, subtitle, and epic header
	if len(card.Body) < 1 {
		t.Error("Expected at least 1 element in card body after adding epic section")
	}

	// Find the epic header (should be a TextBlock with Markdown link)
	var epicHeader *AdaptiveCardElement
	for i := range card.Body {
		if card.Body[i].Type == "TextBlock" && card.Body[i].Weight == "Bolder" {
			epicHeader = &card.Body[i]
			break
		}
	}

	if epicHeader == nil {
		t.Error("Expected to find TextBlock for epic header")
		return
	}

	// Check that the epic text contains Markdown-style link with new format and emoji
	expectedText := fmt.Sprintf("[%s](%s) | 🔄 %s | %s", group.EpicKey, group.EpicURL, group.EpicStatus, group.EpicSummary)
	if epicHeader.Text != expectedText {
		t.Errorf("Expected epic text to be %s, got %s", expectedText, epicHeader.Text)
	}
}

func TestAddIssueSection(t *testing.T) {
	card := NewAdaptiveCard()

	issue := IssueUpdate{
		Key:       "TASK-456",
		Summary:   "Test Task",
		Status:    "In Progress",
		IssueType: "Task",
		URL:       "https://example.com/task",
		Updates: []Update{
			{
				Time:       time.Now(),
				AuthorName: "John Doe",
				Type:       "comment",
				Content:    "Test comment",
			},
		},
	}

	addIssueSection(&card, issue, time.UTC, 1)

	// Should have issue header and updates container
	if len(card.Body) < 2 {
		t.Error("Expected at least 2 elements in card body after adding issue section")
	}

	// Find the issue header (should be a TextBlock with Markdown link)
	var issueHeader *AdaptiveCardElement
	for i := range card.Body {
		if card.Body[i].Type == "TextBlock" && card.Body[i].Weight == "Bolder" {
			issueHeader = &card.Body[i]
			break
		}
	}

	if issueHeader == nil {
		t.Error("Expected to find TextBlock for issue header")
		return
	}

	// Check that the issue text contains Markdown-style link with numbering and emoji
	expectedText := fmt.Sprintf("1. %s | [%s](%s) | 🔄 %s | %s", issue.IssueType, issue.Key, issue.URL, issue.Status, issue.Summary)
	if issueHeader.Text != expectedText {
		t.Errorf("Expected issue text to be %s, got %s", expectedText, issueHeader.Text)
	}
}
