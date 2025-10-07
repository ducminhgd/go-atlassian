package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ducminhgd/go-atlassian/internal/msteams"
)

func main() {
	// Create sample data
	reportDate := time.Now()

	epicGroups := map[string]*msteams.EpicGroup{
		"EPIC-123": {
			EpicKey:     "EPIC-123",
			EpicSummary: "Improve User Experience",
			EpicStatus:  "In Progress",
			EpicURL:     "https://jira.example.com/browse/EPIC-123",
			Issues: []msteams.IssueUpdate{
				{
					Key:         "TASK-456",
					Summary:     "Update login page design",
					Status:      "In Progress",
					IssueType:   "Task",
					URL:         "https://jira.example.com/browse/TASK-456",
					LastUpdated: reportDate.Add(-2 * time.Hour),
					Updates: []msteams.Update{
						{
							Time:       reportDate.Add(-1 * time.Hour),
							AuthorName: "John Doe",
							Type:       "comment",
							Content:    "Updated the mockups based on feedback from the design team",
						},
						{
							Time:       reportDate.Add(-30 * time.Minute),
							AuthorName: "Jane Smith",
							Type:       "worklog",
							Content:    "Implemented responsive design for mobile devices",
							TimeSpent:  "3h",
						},
					},
				},
			},
		},
	}

	noEpicIssues := []msteams.IssueUpdate{
		{
			Key:         "BUG-789",
			Summary:     "Fix memory leak in data processing",
			Status:      "Open",
			IssueType:   "Bug",
			URL:         "https://jira.example.com/browse/BUG-789",
			LastUpdated: reportDate.Add(-3 * time.Hour),
			Updates: []msteams.Update{
				{
					Time:       reportDate.Add(-2 * time.Hour),
					AuthorName: "Bob Wilson",
					Type:       "comment",
					Content:    "Identified the root cause in the batch processing module",
				},
			},
		},
	}

	// Generate AdaptiveCard using the msteams package
	card := msteams.FormatJiraReportAsAdaptiveCard(epicGroups, noEpicIssues, reportDate, "UTC")

	// Convert to JSON for display
	jsonData, err := json.MarshalIndent(card, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling AdaptiveCard: %v\n", err)
		return
	}

	fmt.Println("Generated AdaptiveCard JSON:")
	fmt.Println(string(jsonData))

	// Also demonstrate Teams message format
	teamsMessage := msteams.FormatTeamsMessage(card)
	teamsJSON, err := json.MarshalIndent(teamsMessage, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling Teams message: %v\n", err)
		return
	}

	fmt.Println("\n\nTeams Message JSON (ready to send to webhook):")
	fmt.Println(string(teamsJSON))
}
