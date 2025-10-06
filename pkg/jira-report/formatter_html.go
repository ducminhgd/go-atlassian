package jirareport

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// formatHTMLReport formats the report in HTML with ordered lists
func formatHTMLReport(epicGroups map[string]*EpicGroup, noEpicIssues []IssueUpdate, reportDate time.Time, timezone string) string {
	var report strings.Builder

	// Load timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	reportDate = reportDate.In(loc)

	// Title
	report.WriteString(fmt.Sprintf("<h1>Daily Report %s</h1>\n", reportDate.Format("02-Jan-2006")))
	report.WriteString("<p>From last updates in the last 24 hours</p>\n\n")
	report.WriteString("<ol>\n")

	// Sort epic groups by epic key
	var epicKeys []string
	for key := range epicGroups {
		epicKeys = append(epicKeys, key)
	}
	sort.Strings(epicKeys)

	// Write epic groups
	for _, epicKey := range epicKeys {
		group := epicGroups[epicKey]
		report.WriteString("  <li>\n")
		report.WriteString(fmt.Sprintf("    <h2><a href=\"%s\">%s %s: %s</a></h2>\n",
			escapeHTML(group.EpicURL), escapeHTML(group.EpicKey), escapeHTML(group.EpicStatus), escapeHTML(group.EpicSummary)))
		report.WriteString("    <ol>\n")

		// Sort issues by last updated time
		sort.Slice(group.Issues, func(i, j int) bool {
			return group.Issues[i].LastUpdated.After(group.Issues[j].LastUpdated)
		})

		for _, iss := range group.Issues {
			writeHTMLIssueSection(&report, iss, loc)
		}

		report.WriteString("    </ol>\n")
		report.WriteString("  </li>\n\n")
	}

	// Write "Anything else" section
	if len(noEpicIssues) > 0 {
		report.WriteString("  <li>\n")
		report.WriteString("    <h2>Anything else</h2>\n")
		report.WriteString("    <ol>\n")

		// Sort issues by last updated time
		sort.Slice(noEpicIssues, func(i, j int) bool {
			return noEpicIssues[i].LastUpdated.After(noEpicIssues[j].LastUpdated)
		})

		for _, iss := range noEpicIssues {
			writeHTMLIssueSection(&report, iss, loc)
		}

		report.WriteString("    </ol>\n")
		report.WriteString("  </li>\n")
	}

	report.WriteString("</ol>\n")
	return report.String()
}

// writeHTMLIssueSection writes a single issue section to the HTML report
func writeHTMLIssueSection(report *strings.Builder, iss IssueUpdate, loc *time.Location) {
	report.WriteString("      <li>\n")
	report.WriteString(fmt.Sprintf("        <h3><a href=\"%s\">[%s | %s %s: %s]</a></h3>\n",
		escapeHTML(iss.URL), escapeHTML(iss.IssueType), escapeHTML(iss.Key), escapeHTML(iss.Status), escapeHTML(iss.Summary)))
	report.WriteString("        <ol>\n")

	for _, update := range iss.Updates {
		updateTime := update.Time.In(loc)
		if update.Type == "comment" {
			report.WriteString(fmt.Sprintf("          <li>%s %s commented: %s</li>\n",
				updateTime.Format("15:04"), escapeHTML(update.AuthorName), escapeHTML(truncateText(update.Content, 200))))
		} else if update.Type == "worklog" {
			report.WriteString(fmt.Sprintf("          <li>%s %s log work %s: %s</li>\n",
				updateTime.Format("15:04"), escapeHTML(update.AuthorName), escapeHTML(update.TimeSpent), escapeHTML(truncateText(update.Content, 200))))
		}
	}

	report.WriteString("        </ol>\n")
	report.WriteString("      </li>\n")
}

// escapeHTML escapes special HTML characters
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

