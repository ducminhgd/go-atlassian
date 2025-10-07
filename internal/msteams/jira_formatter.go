package msteams

import (
	"fmt"
	"sort"
	"time"
)

// EpicGroup represents a group of issues under an epic
type EpicGroup struct {
	EpicKey     string
	EpicSummary string
	EpicStatus  string
	EpicURL     string
	Issues      []IssueUpdate
}

// IssueUpdate represents an issue with its updates
type IssueUpdate struct {
	Key         string
	Summary     string
	Status      string
	IssueType   string
	URL         string
	Updates     []Update
	LastUpdated time.Time
}

// Update represents a single update (comment or worklog)
type Update struct {
	Time       time.Time
	AuthorName string
	Type       string // "comment" or "worklog"
	Content    string
	TimeSpent  string // for worklogs
}

// FormatJiraReportAsAdaptiveCard formats a Jira report as an AdaptiveCard
func FormatJiraReportAsAdaptiveCard(epicGroups map[string]*EpicGroup, noEpicIssues []IssueUpdate, reportDate time.Time, timezone string) AdaptiveCard {
	// Load timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	reportDate = reportDate.In(loc)

	// Create new AdaptiveCard
	card := NewAdaptiveCard()

	// Add title
	card.AddTextBlock(
		fmt.Sprintf("Daily Report %s", reportDate.Format("02-Jan-2006")),
		"Large",
		"Bolder",
		true,
	)

	// Add subtitle
	card.AddTextBlock(
		"From last updates in the last 24 hours",
		"Medium",
		"",
		true,
	)

	// Sort epic groups by epic key
	var epicKeys []string
	for key := range epicGroups {
		epicKeys = append(epicKeys, key)
	}
	sort.Strings(epicKeys)

	// Add epic groups
	for _, epicKey := range epicKeys {
		group := epicGroups[epicKey]
		addEpicSection(&card, group, loc)
	}

	// Add "Anything else" section
	if len(noEpicIssues) > 0 {
		addAnythingElseSection(&card, noEpicIssues, loc)
	}

	return card
}

// addEpicSection adds an epic section to the AdaptiveCard
func addEpicSection(card *AdaptiveCard, group *EpicGroup, loc *time.Location) {
	// Epic header with clickable epic key using Markdown-style link
	var epicText string
	if group.EpicURL != "" {
		epicText = fmt.Sprintf("[%s](%s) %s: %s", group.EpicKey, group.EpicURL, group.EpicStatus, group.EpicSummary)
	} else {
		epicText = fmt.Sprintf("%s %s: %s", group.EpicKey, group.EpicStatus, group.EpicSummary)
	}

	card.AddTextBlock(epicText, "Medium", "Bolder", true)

	// Sort issues by last updated time
	sort.Slice(group.Issues, func(i, j int) bool {
		return group.Issues[i].LastUpdated.After(group.Issues[j].LastUpdated)
	})

	// Add issues
	for _, issue := range group.Issues {
		addIssueSection(card, issue, loc)
	}
}

// addAnythingElseSection adds the "Anything else" section
func addAnythingElseSection(card *AdaptiveCard, noEpicIssues []IssueUpdate, loc *time.Location) {
	// "Anything else" header
	card.AddTextBlock(
		"Anything else",
		"Medium",
		"Bolder",
		true,
	)

	// Sort issues by last updated time
	sort.Slice(noEpicIssues, func(i, j int) bool {
		return noEpicIssues[i].LastUpdated.After(noEpicIssues[j].LastUpdated)
	})

	// Add issues
	for _, issue := range noEpicIssues {
		addIssueSection(card, issue, loc)
	}
}

// addIssueSection adds an issue section to the AdaptiveCard
func addIssueSection(card *AdaptiveCard, issue IssueUpdate, loc *time.Location) {
	// Issue header with clickable issue key using Markdown-style link
	var issueText string
	if issue.URL != "" {
		issueText = fmt.Sprintf("%s | [%s](%s) %s: %s", issue.IssueType, issue.Key, issue.URL, issue.Status, issue.Summary)
	} else {
		issueText = fmt.Sprintf("%s | %s %s: %s", issue.IssueType, issue.Key, issue.Status, issue.Summary)
	}

	card.AddTextBlock(issueText, "Small", "Bolder", true)

	// Add updates
	if len(issue.Updates) > 0 {
		var updateItems []AdaptiveCardElement

		for _, update := range issue.Updates {
			updateTime := update.Time.In(loc)
			var updateText string

			switch update.Type {
			case "comment":
				updateText = fmt.Sprintf("%s %s commented: %s",
					updateTime.Format("15:04"), update.AuthorName, truncateText(update.Content, 200))
			case "worklog":
				updateText = fmt.Sprintf("%s %s log work %s: %s",
					updateTime.Format("15:04"), update.AuthorName, update.TimeSpent, truncateText(update.Content, 200))
			}

			updateItems = append(updateItems, AdaptiveCardElement{
				Type: "TextBlock",
				Text: updateText,
				Wrap: true,
				Size: "Small",
			})
		}

		card.AddContainer(updateItems, "Small", "emphasis")
	}
}

// truncateText truncates text to a maximum length
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}
