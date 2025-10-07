package jirareport

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

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
		report.WriteString(fmt.Sprintf("## [%s %s: %s](%s)\n\n", group.EpicKey, group.EpicStatus, group.EpicSummary, group.EpicURL))

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
	report.WriteString(fmt.Sprintf("### [%s | %s %s: %s](%s)\n\n", iss.IssueType, iss.Key, iss.Status, iss.Summary, iss.URL))

	for i, update := range iss.Updates {
		updateTime := update.Time.In(loc)
		if update.Type == "comment" {
			report.WriteString(fmt.Sprintf("%d. %s **%s** commented: %s\n", i+1, updateTime.Format("15:04"), update.AuthorName, truncateText(update.Content, 200)))
		} else if update.Type == "worklog" {
			report.WriteString(fmt.Sprintf("%d. %s **%s** log work %s: %s\n", i+1, updateTime.Format("15:04"), update.AuthorName, update.TimeSpent, truncateText(update.Content, 200)))
		}
	}
	report.WriteString("\n")
}

