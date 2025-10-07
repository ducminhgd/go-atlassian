package jirareport

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

// formatMarkdownReport formats the report in markdown
func formatMarkdownReport(epicGroups map[string]*EpicGroup, noEpicIssues []IssueUpdate, reportDate time.Time, timezone string) string {
	var report strings.Builder

	// Load timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	reportDate = reportDate.In(loc)

	// Title
	report.WriteString(fmt.Sprintf("# Daily Report %s\n\n", reportDate.Format("02-Jan-2006")))
	report.WriteString("From last updates in the last 24 hours\n\n")

	// Sort epic groups by epic key
	var epicKeys []string
	for key := range epicGroups {
		epicKeys = append(epicKeys, key)
	}
	sort.Strings(epicKeys)

	// Write epic groups
	for _, epicKey := range epicKeys {
		group := epicGroups[epicKey]
		statusEmoji := getStatusEmoji(group.EpicStatus)
		if group.EpicURL != "" {
			report.WriteString(fmt.Sprintf("## [%s](%s) | %s %s | %s\n\n", group.EpicKey, group.EpicURL, statusEmoji, group.EpicStatus, group.EpicSummary))
		} else {
			report.WriteString(fmt.Sprintf("## %s | %s %s | %s\n\n", group.EpicKey, statusEmoji, group.EpicStatus, group.EpicSummary))
		}

		// Sort issues by last updated time
		sort.Slice(group.Issues, func(i, j int) bool {
			return group.Issues[i].LastUpdated.After(group.Issues[j].LastUpdated)
		})

		for _, iss := range group.Issues {
			writeMarkdownIssueSection(&report, iss, loc)
		}
	}

	// Write "Anything else" section
	if len(noEpicIssues) > 0 {
		report.WriteString("## Anything else\n\n")

		// Sort issues by last updated time
		sort.Slice(noEpicIssues, func(i, j int) bool {
			return noEpicIssues[i].LastUpdated.After(noEpicIssues[j].LastUpdated)
		})

		for _, iss := range noEpicIssues {
			writeMarkdownIssueSection(&report, iss, loc)
		}
	}

	return report.String()
}

// writeMarkdownIssueSection writes a single issue section to the markdown report
func writeMarkdownIssueSection(report *strings.Builder, iss IssueUpdate, loc *time.Location) {
	statusEmoji := getStatusEmoji(iss.Status)
	if iss.URL != "" {
		report.WriteString(fmt.Sprintf("### %s | [%s](%s) | %s %s | %s\n\n", iss.IssueType, iss.Key, iss.URL, statusEmoji, iss.Status, iss.Summary))
	} else {
		report.WriteString(fmt.Sprintf("### %s | %s | %s %s | %s\n\n", iss.IssueType, iss.Key, statusEmoji, iss.Status, iss.Summary))
	}

	for i, update := range iss.Updates {
		updateTime := update.Time.In(loc)
		if update.Type == "comment" {
			report.WriteString(fmt.Sprintf("%d. %s → **%s** commented: %s\n", i+1, updateTime.Format("15:04"), update.AuthorName, truncateText(update.Content, 200)))
		} else if update.Type == "worklog" {
			report.WriteString(fmt.Sprintf("%d. %s → **%s** log work %s: %s\n", i+1, updateTime.Format("15:04"), update.AuthorName, update.TimeSpent, truncateText(update.Content, 200)))
		}
	}

	// Add sub-tasks
	if len(iss.SubTasks) > 0 {
		// Sort sub-tasks by last updated time
		sort.Slice(iss.SubTasks, func(i, j int) bool {
			return iss.SubTasks[i].LastUpdated.After(iss.SubTasks[j].LastUpdated)
		})

		for _, subTask := range iss.SubTasks {
			writeMarkdownSubTaskSection(report, subTask, loc)
		}
	}

	report.WriteString("\n")
}

// writeMarkdownSubTaskSection writes a single sub-task section to the markdown report
func writeMarkdownSubTaskSection(report *strings.Builder, subTask IssueUpdate, loc *time.Location) {
	statusEmoji := getStatusEmoji(subTask.Status)
	if subTask.URL != "" {
		report.WriteString(fmt.Sprintf("#### %s | [%s](%s) | %s %s | %s\n\n", subTask.IssueType, subTask.Key, subTask.URL, statusEmoji, subTask.Status, subTask.Summary))
	} else {
		report.WriteString(fmt.Sprintf("#### %s | %s | %s %s | %s\n\n", subTask.IssueType, subTask.Key, statusEmoji, subTask.Status, subTask.Summary))
	}

	for i, update := range subTask.Updates {
		updateTime := update.Time.In(loc)
		if update.Type == "comment" {
			report.WriteString(fmt.Sprintf("%d. %s → **%s** commented: %s\n", i+1, updateTime.Format("15:04"), update.AuthorName, truncateText(update.Content, 200)))
		} else if update.Type == "worklog" {
			report.WriteString(fmt.Sprintf("%d. %s → **%s** log work %s: %s\n", i+1, updateTime.Format("15:04"), update.AuthorName, update.TimeSpent, truncateText(update.Content, 200)))
		}
	}
	report.WriteString("\n")
}
