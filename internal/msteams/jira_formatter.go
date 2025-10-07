package msteams

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// getStatusEmoji returns an emoji for the given status
func getStatusEmoji(status string) string {
	switch strings.ToLower(status) {
	case "to do", "open", "new", "created":
		return "📋"
	case "in progress", "in review", "in development":
		return "🔄"
	case "done", "closed", "resolved", "completed":
		return "✅"
	case "blocked", "on hold":
		return "🚫"
	case "testing", "qa", "review":
		return "🧪"
	default:
		return "📝"
	}
}

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
	SubTasks    []IssueUpdate // Sub-tasks belonging to this issue
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
		"ExtraLarge",
		"Bolder",
		true,
	)

	// Add subtitle
	card.AddTextBlock(
		"> From last updates in the last 24 hours",
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

	// Add "Everything else" section
	if len(noEpicIssues) > 0 {
		addEverythingElseSection(&card, noEpicIssues, loc)
	}

	return card
}

// addEpicSection adds an epic section to the AdaptiveCard
func addEpicSection(card *AdaptiveCard, group *EpicGroup, loc *time.Location) {
	// Epic header with clickable epic key using Markdown-style link and status emoji
	statusEmoji := getStatusEmoji(group.EpicStatus)
	var epicText string
	if group.EpicURL != "" {
		epicText = fmt.Sprintf("[%s](%s) | %s %s | %s", group.EpicKey, group.EpicURL, statusEmoji, group.EpicStatus, group.EpicSummary)
	} else {
		epicText = fmt.Sprintf("%s | %s %s | %s", group.EpicKey, statusEmoji, group.EpicStatus, group.EpicSummary)
	}

	card.AddTextBlock(epicText, "Large", "Bolder", true)

	// Sort issues by last updated time
	sort.Slice(group.Issues, func(i, j int) bool {
		return group.Issues[i].LastUpdated.After(group.Issues[j].LastUpdated)
	})

	for _, issue := range group.Issues {
		addIssueSection(card, issue, loc)
	}
}

// addEverythingElseSection adds the "Everything else" section
func addEverythingElseSection(card *AdaptiveCard, noEpicIssues []IssueUpdate, loc *time.Location) {
	// "Everything else" header
	card.AddTextBlock(
		"Everything else",
		"Large",
		"Bolder",
		true,
	)

	// Sort issues by last updated time
	sort.Slice(noEpicIssues, func(i, j int) bool {
		return noEpicIssues[i].LastUpdated.After(noEpicIssues[j].LastUpdated)
	})

	for _, issue := range noEpicIssues {
		addIssueSection(card, issue, loc)
	}
}

// addIssueSection adds an issue section to the AdaptiveCard
func addIssueSection(card *AdaptiveCard, issue IssueUpdate, loc *time.Location) {
	// Issue header with clickable issue key using Markdown-style link, hierarchical numbering, and status emoji
	statusEmoji := getStatusEmoji(issue.Status)
	var issueText string
	if issue.URL != "" {
		issueText = fmt.Sprintf("%s | [%s](%s) | %s %s | %s", issue.IssueType, issue.Key, issue.URL, statusEmoji, issue.Status, issue.Summary)
	} else {
		issueText = fmt.Sprintf("%s | %s | %s %s | %s", issue.IssueType, issue.Key, statusEmoji, issue.Status, issue.Summary)
	}

	card.AddTextBlock(issueText, "Medium", "Bolder", true)

	// Add task updates and sub-tasks as individual TextBlocks
	itemNumber := 1

	// Add task updates first as individual TextBlocks
	for _, update := range issue.Updates {
		updateTime := update.Time.In(loc)
		var updateText string

		switch update.Type {
		case "comment":
			updateText = fmt.Sprintf("%d. %s → %s commented: %s", itemNumber, updateTime.Format("15:04"), update.AuthorName, truncateText(update.Content, 200))
		case "worklog":
			updateText = fmt.Sprintf("%d. %s → %s log work %s: %s", itemNumber, updateTime.Format("15:04"), update.AuthorName, update.TimeSpent, truncateText(update.Content, 200))
		}

		card.AddTextBlock(updateText, "Default", "", true)
		itemNumber++
	}

	// Add sub-tasks continuing the numbering sequence as individual TextBlocks
	if len(issue.SubTasks) > 0 {
		// Sort sub-tasks by last updated time
		sort.Slice(issue.SubTasks, func(i, j int) bool {
			return issue.SubTasks[i].LastUpdated.After(issue.SubTasks[j].LastUpdated)
		})

		for _, subTask := range issue.SubTasks {
			addSubTaskSection(card, subTask, loc, itemNumber)
			itemNumber++
		}
	}
}

// addSubTaskSection adds a sub-task section as individual TextBlocks
func addSubTaskSection(card *AdaptiveCard, subTask IssueUpdate, loc *time.Location, itemNumber int) {
	// Sub-task header with clickable issue key using Markdown-style link, numbered format, and status emoji
	statusEmoji := getStatusEmoji(subTask.Status)
	var subTaskText string
	if subTask.URL != "" {
		subTaskText = fmt.Sprintf("%d. %s | [%s](%s) | %s %s | %s", itemNumber, subTask.IssueType, subTask.Key, subTask.URL, statusEmoji, subTask.Status, subTask.Summary)
	} else {
		subTaskText = fmt.Sprintf("%d. %s | %s | %s %s | %s", itemNumber, subTask.IssueType, subTask.Key, statusEmoji, subTask.Status, subTask.Summary)
	}

	// Add sub-task header as individual TextBlock
	card.AddTextBlock(subTaskText, "Default", "Bolder", true)

	// Add sub-task updates as individual TextBlocks with indentation
	for _, update := range subTask.Updates {
		updateTime := update.Time.In(loc)
		var updateText string

		switch update.Type {
		case "comment":
			updateText = fmt.Sprintf("- %s → %s commented: %s", updateTime.Format("15:04"), update.AuthorName, update.Content)
		case "worklog":
			updateText = fmt.Sprintf("- %s → %s log work %s: %s", updateTime.Format("15:04"), update.AuthorName, update.TimeSpent, update.Content)
		}

		card.AddTextBlock(updateText, "Default", "", true)
	}
}

// truncateText truncates text to a maximum length
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}
