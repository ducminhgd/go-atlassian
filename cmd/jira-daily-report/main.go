package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ducminhgd/go-atlassian/jira/v3/auth"
	"github.com/ducminhgd/go-atlassian/jira/v3/issue"
	"github.com/ducminhgd/go-atlassian/jira/v3/utils"
)

// Config holds the configuration for the daily report
type Config struct {
	JiraHost        string
	JiraUsername    string
	JiraPassword    string
	JiraProject     string
	TeamsWebhookURL string
	ReportTimezone  string
	LookbackHours   int
}

// loadConfig loads configuration from environment variables
func loadConfig() (*Config, error) {
	config := &Config{
		JiraHost:        os.Getenv("JIRA_HOST"),
		JiraUsername:    os.Getenv("JIRA_USERNAME"),
		JiraPassword:    os.Getenv("JIRA_PASSWORD"),
		JiraProject:     os.Getenv("JIRA_PROJECT"),
		TeamsWebhookURL: os.Getenv("TEAMS_WEBHOOK_URL"),
		ReportTimezone:  os.Getenv("REPORT_TIMEZONE"),
		LookbackHours:   24,
	}

	if config.JiraHost == "" {
		return nil, fmt.Errorf("JIRA_HOST is required")
	}
	if config.JiraPassword == "" {
		return nil, fmt.Errorf("JIRA_PASSWORD is required")
	}
	if config.JiraProject == "" {
		return nil, fmt.Errorf("JIRA_PROJECT is required")
	}
	if config.TeamsWebhookURL == "" {
		return nil, fmt.Errorf("TEAMS_WEBHOOK_URL is required")
	}
	if config.ReportTimezone == "" {
		config.ReportTimezone = "UTC"
	}

	return config, nil
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
}

// Update represents a single update (comment or worklog)
type Update struct {
	Time       time.Time
	AuthorName string
	Type       string // "comment" or "worklog"
	Content    string
	TimeSpent  string // for worklogs
}

func main() {
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Setup Jira client
	authenticator := auth.NewBasicAuth(config.JiraUsername, config.JiraPassword)
	client := &http.Client{}
	issueService := issue.NewService(client, config.JiraHost, authenticator)

	// Generate report
	markdownReport, htmlReport, err := generateDailyReport(context.Background(), issueService, config)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	// Print markdown report to console
	fmt.Println(markdownReport)

	// Post HTML report to Microsoft Teams Workflow
	if err := postToTeams(config.TeamsWebhookURL, htmlReport); err != nil {
		log.Fatalf("Failed to post to Teams: %v", err)
	}

	log.Println("Daily report posted successfully!")
}

// generateDailyReport generates the daily report in both markdown and HTML formats
func generateDailyReport(ctx context.Context, issueService *issue.Service, config *Config) (string, string, error) {
	// Calculate time range
	now := time.Now()
	lookbackTime := now.Add(-time.Duration(config.LookbackHours) * time.Hour)

	// Build JQL query to find updated issues
	jql := fmt.Sprintf(
		"project = %s AND updated >= -%dh ORDER BY updated DESC",
		config.JiraProject,
		config.LookbackHours,
	)

	// Search for issues
	searchRequest := issue.JQLSearchRequest{
		JQL:        jql,
		MaxResults: 1000,
		Fields: []string{
			"summary", "status", "issuetype", "parent", "updated", "created",
			"comment", "worklog",
		},
		Expand: "changelog",
	}

	response, err := issueService.SearchJQL(ctx, searchRequest)
	if err != nil {
		return "", "", fmt.Errorf("failed to search issues: %w", err)
	}

	// Process issues and group by epic
	epicGroups := make(map[string]*EpicGroup)
	var noEpicIssues []IssueUpdate

	for _, iss := range response.Issues {
		issueUpdate := processIssue(iss, lookbackTime, config.JiraHost)
		if len(issueUpdate.Updates) == 0 {
			continue // Skip issues with no relevant updates
		}

		// Determine epic
		epicKey := ""
		if iss.Fields.Parent.Key != "" {
			// Get parent issue to check if it's an epic
			parentIssue, err := issueService.Get(ctx, iss.Fields.Parent.Key, nil, []string{"summary", "status", "issuetype"}, nil)
			if err == nil && parentIssue.Fields.IssueType.Name == "Epic" {
				epicKey = parentIssue.Key
				if _, exists := epicGroups[epicKey]; !exists {
					epicGroups[epicKey] = &EpicGroup{
						EpicKey:     parentIssue.Key,
						EpicSummary: parentIssue.Fields.Summary,
						EpicStatus:  parentIssue.Fields.Status.Name,
						EpicURL:     fmt.Sprintf("%s/browse/%s", config.JiraHost, parentIssue.Key),
						Issues:      []IssueUpdate{},
					}
				}
			}
		}

		if epicKey != "" {
			epicGroups[epicKey].Issues = append(epicGroups[epicKey].Issues, issueUpdate)
		} else {
			noEpicIssues = append(noEpicIssues, issueUpdate)
		}
	}

	// Generate markdown report for console and HTML for Teams
	markdownReport := formatMarkdownReport(epicGroups, noEpicIssues, now, config.ReportTimezone)
	htmlReport := formatHTMLReport(epicGroups, noEpicIssues, now, config.ReportTimezone)

	return markdownReport, htmlReport, nil
}

// processIssue processes a single issue and extracts relevant updates
func processIssue(iss issue.Issue, lookbackTime time.Time, jiraHost string) IssueUpdate {
	issueUpdate := IssueUpdate{
		Key:       iss.Key,
		Summary:   iss.Fields.Summary,
		Status:    iss.Fields.Status.Name,
		IssueType: iss.Fields.IssueType.Name,
		URL:       fmt.Sprintf("%s/browse/%s", jiraHost, iss.Key),
		Updates:   []Update{},
	}

	// Process comments
	for _, comment := range iss.Fields.Comment.Comments {
		commentTime, err := time.Parse(utils.JIRATIMEFORMAT, comment.Created)
		if err == nil && commentTime.After(lookbackTime) {
			content := extractTextFromBody(comment.Body)
			issueUpdate.Updates = append(issueUpdate.Updates, Update{
				Time:       commentTime,
				AuthorName: comment.Author.DisplayName,
				Type:       "comment",
				Content:    content,
			})
			if commentTime.After(issueUpdate.LastUpdated) {
				issueUpdate.LastUpdated = commentTime
			}
		}
	}

	// Process worklogs
	for _, worklog := range iss.Fields.Worklog.Worklogs {
		worklogTime, err := time.Parse(utils.JIRATIMEFORMAT, worklog.Created)
		if err == nil && worklogTime.After(lookbackTime) {
			content := extractTextFromBody(worklog.Comment)
			issueUpdate.Updates = append(issueUpdate.Updates, Update{
				Time:       worklogTime,
				AuthorName: worklog.Author.DisplayName,
				Type:       "worklog",
				Content:    content,
				TimeSpent:  worklog.TimeSpent,
			})
			if worklogTime.After(issueUpdate.LastUpdated) {
				issueUpdate.LastUpdated = worklogTime
			}
		}
	}

	// Sort updates by time
	sort.Slice(issueUpdate.Updates, func(i, j int) bool {
		return issueUpdate.Updates[i].Time.Before(issueUpdate.Updates[j].Time)
	})

	return issueUpdate
}

// extractTextFromBody extracts plain text from ADF or string body
func extractTextFromBody(body interface{}) string {
	if body == nil {
		return ""
	}

	// Try to handle as string first
	if str, ok := body.(string); ok {
		return str
	}

	// Try to handle as ADF (Atlassian Document Format)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Sprintf("%v", body)
	}

	var adf map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &adf); err != nil {
		return string(bodyBytes)
	}

	// Extract text from ADF content
	return extractTextFromADF(adf)
}

// extractTextFromADF recursively extracts text from ADF structure
func extractTextFromADF(node map[string]interface{}) string {
	var text strings.Builder

	if nodeType, ok := node["type"].(string); ok && nodeType == "text" {
		if textContent, ok := node["text"].(string); ok {
			return textContent
		}
	}

	if content, ok := node["content"].([]interface{}); ok {
		for _, item := range content {
			if itemMap, ok := item.(map[string]interface{}); ok {
				text.WriteString(extractTextFromADF(itemMap))
				text.WriteString(" ")
			}
		}
	}

	return strings.TrimSpace(text.String())
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

// truncateText truncates text to a maximum length
func truncateText(text string, maxLen int) string {
	text = strings.TrimSpace(text)
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// postToTeams posts the HTML report to Microsoft Teams Workflow
func postToTeams(webhookURL, htmlReport string) error {
	// Microsoft Teams Workflow can accept HTML content in the body field
	// Send as a simple message with HTML body
	payload := map[string]interface{}{
		"body": htmlReport,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to post to Teams: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Teams webhook returned status %d", resp.StatusCode)
	}

	return nil
}
