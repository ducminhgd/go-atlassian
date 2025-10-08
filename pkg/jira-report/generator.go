package jirareport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/ducminhgd/go-atlassian/internal/msteams"
	"github.com/ducminhgd/go-atlassian/jira/v3/auth"
	"github.com/ducminhgd/go-atlassian/jira/v3/issue"
	"github.com/ducminhgd/go-atlassian/jira/v3/utils"
)

// Generator handles the report generation
type Generator struct {
	config        *Config
	issueService  *issue.Service
	authenticator auth.Authenticator
}

// NewGenerator creates a new report generator
func NewGenerator(config *Config) (*Generator, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Setup Jira client
	var authenticator auth.Authenticator
	if config.JiraUsername != "" {
		authenticator = auth.NewBasicAuth(config.JiraUsername, config.JiraPassword)
	} else {
		// If no username provided, use empty username with token as password
		authenticator = auth.NewBasicAuth("", config.JiraPassword)
	}

	client := &http.Client{}
	issueService := issue.NewService(client, config.JiraHost, authenticator)

	return &Generator{
		config:        config,
		issueService:  issueService,
		authenticator: authenticator,
	}, nil
}

// Generate generates the daily report
func (g *Generator) Generate(ctx context.Context) (*Report, error) {
	// Calculate time range
	now := time.Now()
	lookbackTime := now.Add(-time.Duration(g.config.LookbackHours) * time.Hour)

	// Build search request based on query type
	searchRequest, err := g.buildSearchRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSearchIssues, err)
	}

	// Search for issues
	response, err := g.issueService.SearchJQL(ctx, *searchRequest)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSearchIssues, err)
	}

	// Process issues and group by epic
	epicGroups := make(map[string]*EpicGroup)
	var noEpicIssues []IssueUpdate
	processedIssues := make(map[string]*IssueUpdate) // Track processed issues to avoid duplicates
	addedToReport := make(map[string]bool)           // Track which issues have been added to the final report

	for _, iss := range response.Issues {
		// Skip if this issue has already been processed and added to the report
		if addedToReport[iss.Key] {
			continue
		}

		issueUpdate := g.processIssue(iss, lookbackTime)
		if len(issueUpdate.Updates) == 0 {
			continue // Skip issues with no relevant updates
		}

		processedIssues[iss.Key] = &issueUpdate

		// Handle parent-child relationships
		epicKey, parentTaskKey := g.determineParentRelationships(ctx, iss, epicGroups)

		if parentTaskKey != "" {
			// This is a sub-task, add it to its parent task
			if parentIssue, exists := processedIssues[parentTaskKey]; exists {
				parentIssue.SubTasks = append(parentIssue.SubTasks, issueUpdate)

				// If parent exists but hasn't been added to report yet, add it now
				if !addedToReport[parentTaskKey] {
					if epicKey != "" {
						epicGroups[epicKey].Issues = append(epicGroups[epicKey].Issues, *parentIssue)
					} else {
						noEpicIssues = append(noEpicIssues, *parentIssue)
					}
					addedToReport[parentTaskKey] = true
				}
			} else {
				// Parent task not yet processed, fetch it and treat as updated
				parentTask, err := g.fetchAndProcessParentTask(ctx, parentTaskKey, lookbackTime)
				if err == nil {
					parentTask.SubTasks = append(parentTask.SubTasks, issueUpdate)
					processedIssues[parentTaskKey] = parentTask

					// Only add parent task to report if it hasn't been added yet
					if !addedToReport[parentTaskKey] {
						if epicKey != "" {
							epicGroups[epicKey].Issues = append(epicGroups[epicKey].Issues, *parentTask)
						} else {
							noEpicIssues = append(noEpicIssues, *parentTask)
						}
						addedToReport[parentTaskKey] = true
					}
				}
			}
			// Mark this sub-task as processed (it's included in its parent)
			addedToReport[iss.Key] = true
		} else {
			// This is a top-level task or epic-level task
			// Only add to report if it hasn't been added yet
			if !addedToReport[iss.Key] {
				if epicKey != "" {
					epicGroups[epicKey].Issues = append(epicGroups[epicKey].Issues, issueUpdate)
				} else {
					noEpicIssues = append(noEpicIssues, issueUpdate)
				}
				addedToReport[iss.Key] = true
			}
		}
	}

	// Generate markdown and AdaptiveCard reports
	markdownReport := formatMarkdownReport(epicGroups, noEpicIssues, now, g.config.Timezone)
	adaptiveCardReport := g.formatAdaptiveCardReport(epicGroups, noEpicIssues, now, g.config.Timezone)

	return &Report{
		Markdown:     markdownReport,
		AdaptiveCard: adaptiveCardReport,
	}, nil
}

// determineParentRelationships determines the epic and parent task relationships for an issue
func (g *Generator) determineParentRelationships(ctx context.Context, iss issue.Issue, epicGroups map[string]*EpicGroup) (epicKey, parentTaskKey string) {
	if iss.Fields.Parent.Key == "" {
		return "", ""
	}

	// Get parent issue
	parentIssue, err := g.issueService.Get(ctx, iss.Fields.Parent.Key, nil, []string{"summary", "status", "issuetype", "parent"}, nil)
	if err != nil {
		return "", ""
	}

	if parentIssue.Fields.IssueType.Name == "Epic" {
		// Parent is an epic
		epicKey = parentIssue.Key
		if _, exists := epicGroups[epicKey]; !exists {
			epicGroups[epicKey] = &EpicGroup{
				EpicKey:     parentIssue.Key,
				EpicSummary: parentIssue.Fields.Summary,
				EpicStatus:  parentIssue.Fields.Status.Name,
				EpicURL:     fmt.Sprintf("%s/browse/%s", g.config.JiraHost, parentIssue.Key),
				Issues:      []IssueUpdate{},
			}
		}
		return epicKey, ""
	} else {
		// Parent is a task/story, this is a sub-task
		parentTaskKey = parentIssue.Key

		// Check if the parent task has an epic
		if parentIssue.Fields.Parent.Key != "" {
			grandParent, err := g.issueService.Get(ctx, parentIssue.Fields.Parent.Key, nil, []string{"summary", "status", "issuetype"}, nil)
			if err == nil && grandParent.Fields.IssueType.Name == "Epic" {
				epicKey = grandParent.Key
				if _, exists := epicGroups[epicKey]; !exists {
					epicGroups[epicKey] = &EpicGroup{
						EpicKey:     grandParent.Key,
						EpicSummary: grandParent.Fields.Summary,
						EpicStatus:  grandParent.Fields.Status.Name,
						EpicURL:     fmt.Sprintf("%s/browse/%s", g.config.JiraHost, grandParent.Key),
						Issues:      []IssueUpdate{},
					}
				}
			}
		}

		return epicKey, parentTaskKey
	}
}

// fetchAndProcessParentTask fetches and processes a parent task that wasn't in the original search results
func (g *Generator) fetchAndProcessParentTask(ctx context.Context, parentKey string, lookbackTime time.Time) (*IssueUpdate, error) {
	parentIssue, err := g.issueService.Get(ctx, parentKey, nil, []string{"summary", "status", "issuetype", "comment", "worklog"}, []string{"changelog"})
	if err != nil {
		return nil, err
	}

	// Process the parent task as if it was updated (since its child was updated)
	parentUpdate := g.processIssue(*parentIssue, lookbackTime)

	// If the parent has no direct updates, we still want to include it because its child was updated
	if len(parentUpdate.Updates) == 0 {
		// Add a synthetic update to indicate the parent should be included
		parentUpdate.LastUpdated = time.Now()
	}

	return &parentUpdate, nil
}

// processIssue processes a single issue and extracts relevant updates
func (g *Generator) processIssue(iss issue.Issue, lookbackTime time.Time) IssueUpdate {
	issueUpdate := IssueUpdate{
		Key:       iss.Key,
		Summary:   iss.Fields.Summary,
		Status:    iss.Fields.Status.Name,
		IssueType: iss.Fields.IssueType.Name,
		URL:       fmt.Sprintf("%s/browse/%s", g.config.JiraHost, iss.Key),
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

// buildSearchRequest builds the search request based on query type
func (g *Generator) buildSearchRequest(ctx context.Context) (*issue.JQLSearchRequest, error) {
	var jql string
	var err error

	switch g.config.QueryType {
	case QueryTypeProjectAndHours:
		jql = g.buildProjectAndHoursJQL()
	case QueryTypeCustomJQL:
		jql = g.config.CustomJQL
	case QueryTypeFilter:
		jql, err = g.getFilterJQL(ctx)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidQueryType
	}

	return &issue.JQLSearchRequest{
		JQL:        jql,
		MaxResults: 1000,
		Fields: []string{
			"summary", "status", "issuetype", "parent", "updated", "created",
			"comment", "worklog",
		},
		Expand: "changelog",
	}, nil
}

// buildProjectAndHoursJQL builds JQL for project + hours query
func (g *Generator) buildProjectAndHoursJQL() string {
	return fmt.Sprintf(
		"project = %s AND updated >= -%dh ORDER BY updated DESC",
		g.config.JiraProject,
		g.config.LookbackHours,
	)
}

// getFilterJQL retrieves JQL from a saved filter
func (g *Generator) getFilterJQL(ctx context.Context) (string, error) {
	// Build the URL for getting filter details
	url := fmt.Sprintf("%s/rest/api/3/filter/%s", g.config.JiraHost, g.config.FilterID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication
	if err := g.authenticator.AddAuthentication(req); err != nil {
		return "", fmt.Errorf("failed to add authentication: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get filter: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get filter: status %d", resp.StatusCode)
	}

	var filterResponse struct {
		JQL  string `json:"jql"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&filterResponse); err != nil {
		return "", fmt.Errorf("failed to decode filter response: %w", err)
	}

	// Store filter name in config for subtitle generation
	g.config.FilterName = filterResponse.Name

	return filterResponse.JQL, nil
}

// WithProjectAndHours configures the generator to use project + hours query
func (g *Generator) WithProjectAndHours(project string, hours int) *Generator {
	g.config.QueryType = QueryTypeProjectAndHours
	g.config.JiraProject = project
	g.config.LookbackHours = hours
	return g
}

// WithCustomJQL configures the generator to use custom JQL query
func (g *Generator) WithCustomJQL(jql string) *Generator {
	g.config.QueryType = QueryTypeCustomJQL
	g.config.CustomJQL = jql
	return g
}

// WithFilter configures the generator to use a saved filter
func (g *Generator) WithFilter(filterID string) *Generator {
	g.config.QueryType = QueryTypeFilter
	g.config.FilterID = filterID
	return g
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

// truncateText truncates text to a maximum length
func truncateText(text string, maxLen int) string {
	text = strings.TrimSpace(text)
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// formatAdaptiveCardReport converts internal types to msteams types and formats as AdaptiveCard
func (g *Generator) formatAdaptiveCardReport(epicGroups map[string]*EpicGroup, noEpicIssues []IssueUpdate, reportDate time.Time, timezone string) msteams.AdaptiveCard {
	// Convert internal types to msteams types
	msteamsEpicGroups := make(map[string]*msteams.EpicGroup)
	for key, group := range epicGroups {
		msteamsEpicGroups[key] = &msteams.EpicGroup{
			EpicKey:     group.EpicKey,
			EpicSummary: group.EpicSummary,
			EpicStatus:  group.EpicStatus,
			EpicURL:     group.EpicURL,
			Issues:      convertIssueUpdates(group.Issues),
		}
	}

	msteamsNoEpicIssues := convertIssueUpdates(noEpicIssues)

	// Create subtitle config
	subtitleConfig := msteams.SubtitleConfig{
		QueryType:     string(g.config.QueryType),
		FilterName:    g.config.FilterName,
		FilterID:      g.config.FilterID,
		CustomJQL:     g.config.CustomJQL,
		JiraProject:   g.config.JiraProject,
		LookbackHours: g.config.LookbackHours,
		JiraHost:      g.config.JiraHost,
	}

	return msteams.FormatJiraReportAsAdaptiveCard(msteamsEpicGroups, msteamsNoEpicIssues, reportDate, timezone, subtitleConfig)
}

// convertIssueUpdates converts internal IssueUpdate slice to msteams IssueUpdate slice
func convertIssueUpdates(issues []IssueUpdate) []msteams.IssueUpdate {
	var msteamsIssues []msteams.IssueUpdate
	for _, issue := range issues {
		msteamsIssues = append(msteamsIssues, msteams.IssueUpdate{
			Key:         issue.Key,
			Summary:     issue.Summary,
			Status:      issue.Status,
			IssueType:   issue.IssueType,
			URL:         issue.URL,
			Updates:     convertUpdates(issue.Updates),
			LastUpdated: issue.LastUpdated,
			SubTasks:    convertIssueUpdates(issue.SubTasks), // Recursive conversion for sub-tasks
		})
	}
	return msteamsIssues
}

// convertUpdates converts internal Update slice to msteams Update slice
func convertUpdates(updates []Update) []msteams.Update {
	var msteamsUpdates []msteams.Update
	for _, update := range updates {
		msteamsUpdates = append(msteamsUpdates, msteams.Update{
			Time:       update.Time,
			AuthorName: update.AuthorName,
			Type:       update.Type,
			Content:    update.Content,
			TimeSpent:  update.TimeSpent,
		})
	}
	return msteamsUpdates
}

// FormatAsAdaptiveCard creates an AdaptiveCard from the provided data
// This is a public function that allows external users to create AdaptiveCards
// without needing access to the internal msteams package
func FormatAsAdaptiveCard(epicGroups map[string]*EpicGroup, noEpicIssues []IssueUpdate, reportDate time.Time, timezone string) msteams.AdaptiveCard {
	// Convert public types to internal msteams types
	msteamsEpicGroups := make(map[string]*msteams.EpicGroup)
	for key, group := range epicGroups {
		msteamsEpicGroups[key] = &msteams.EpicGroup{
			EpicKey:     group.EpicKey,
			EpicSummary: group.EpicSummary,
			EpicStatus:  group.EpicStatus,
			EpicURL:     group.EpicURL,
			Issues:      convertIssueUpdates(group.Issues),
		}
	}

	msteamsNoEpicIssues := convertIssueUpdates(noEpicIssues)

	// Create a basic subtitle config for standalone usage
	subtitleConfig := msteams.SubtitleConfig{
		QueryType:     "manual",
		LookbackHours: 24,
	}

	return msteams.FormatJiraReportAsAdaptiveCard(msteamsEpicGroups, msteamsNoEpicIssues, reportDate, timezone, subtitleConfig)
}

// FormatAsTeamsMessage creates a Teams message from an AdaptiveCard
// This wraps the AdaptiveCard in the proper Teams message structure
func FormatAsTeamsMessage(adaptiveCard msteams.AdaptiveCard) msteams.TeamsMessage {
	return msteams.FormatTeamsMessage(adaptiveCard)
}
