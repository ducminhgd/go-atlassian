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

	// Process issues and build hierarchical tree structure
	epicGroups := make(map[string]*EpicGroup)
	var noEpicIssues []IssueUpdate
	processedIssues := make(map[string]*IssueUpdate) // Track all processed issues

	for _, iss := range response.Issues {
		// Skip if this issue has already been processed and added to report
		if existingIssue, exists := processedIssues[iss.Key]; exists && existingIssue.AddedToReport {
			continue
		}

		issueUpdate := g.processIssue(iss, lookbackTime)
		if len(issueUpdate.Updates) == 0 {
			continue // Skip issues with no relevant updates
		}

		// Store the processed issue
		processedIssues[iss.Key] = &issueUpdate

		// Build the hierarchical tree based on issue level
		g.addIssueToHierarchy(ctx, iss, &issueUpdate, epicGroups, &noEpicIssues, processedIssues, lookbackTime)
	}

	// Generate markdown and AdaptiveCard reports
	markdownReport := formatMarkdownReport(epicGroups, noEpicIssues, now, g.config.Timezone)
	adaptiveCardReport := g.formatAdaptiveCardReport(epicGroups, noEpicIssues, now, g.config.Timezone)

	return &Report{
		Markdown:     markdownReport,
		AdaptiveCard: adaptiveCardReport,
	}, nil
}

// addIssueToHierarchy adds an issue to the appropriate place in the hierarchical tree
func (g *Generator) addIssueToHierarchy(ctx context.Context, originalIssue issue.Issue, issueUpdate *IssueUpdate, epicGroups map[string]*EpicGroup, noEpicIssues *[]IssueUpdate, processedIssues map[string]*IssueUpdate, lookbackTime time.Time) {
	// Skip if already added to report
	if issueUpdate.AddedToReport {
		return
	}

	// Determine the hierarchy level and parent relationships using original issue data
	parentKey, epicKey := g.findParentAndEpicFromIssue(originalIssue)

	if parentKey == "" {
		// Level 1 (Epic) or Level 2 (Story/Task without Epic)
		g.addTopLevelIssue(issueUpdate, epicKey, epicGroups, noEpicIssues)
	} else {
		// Level 3 (Subtask) - find or create parent (Level 2)
		g.addSubtaskToParent(ctx, issueUpdate, parentKey, epicKey, epicGroups, noEpicIssues, processedIssues, lookbackTime)
	}
}

// findParentAndEpicFromIssue determines the parent and epic for an issue using existing data
func (g *Generator) findParentAndEpicFromIssue(iss issue.Issue) (parentKey, epicKey string) {
	if iss.Fields.Parent.Key == "" {
		// No parent - this is either Epic or top-level Story/Task
		return "", ""
	}

	parentKey = iss.Fields.Parent.Key

	// We need to fetch parent details to determine hierarchy
	// This is unavoidable as we need to know if parent is Epic or Story/Task
	parentDetails, err := g.issueService.Get(context.Background(), parentKey, nil, []string{"parent", "issuetype"}, nil)
	if err != nil {
		return parentKey, ""
	}

	if parentDetails.Fields.IssueType.Name == "Epic" {
		// Parent is Epic (Level 1) - current issue is Level 2
		return "", parentKey
	} else {
		// Parent is Story/Task (Level 2) - current issue is Level 3
		// Check if parent has an Epic (grandparent)
		if parentDetails.Fields.Parent.Key != "" {
			grandParentDetails, err := g.issueService.Get(context.Background(), parentDetails.Fields.Parent.Key, nil, []string{"issuetype"}, nil)
			if err == nil && grandParentDetails.Fields.IssueType.Name == "Epic" {
				epicKey = grandParentDetails.Key
			}
		}
		return parentKey, epicKey
	}
}

// addTopLevelIssue adds a top-level issue (Epic or Story/Task) to the appropriate group
func (g *Generator) addTopLevelIssue(issue *IssueUpdate, epicKey string, epicGroups map[string]*EpicGroup, noEpicIssues *[]IssueUpdate) {
	if epicKey != "" {
		// This issue belongs to an Epic - ensure Epic group exists
		g.ensureEpicGroupExists(epicKey, epicGroups)
		epicGroups[epicKey].Issues = append(epicGroups[epicKey].Issues, *issue)
	} else {
		// No Epic - add to standalone issues
		*noEpicIssues = append(*noEpicIssues, *issue)
	}
	issue.AddedToReport = true
}

// addSubtaskToParent adds a subtask to its parent task
func (g *Generator) addSubtaskToParent(ctx context.Context, subtask *IssueUpdate, parentKey, epicKey string, epicGroups map[string]*EpicGroup, noEpicIssues *[]IssueUpdate, processedIssues map[string]*IssueUpdate, lookbackTime time.Time) {
	// Check if parent is already processed
	if parentIssue, exists := processedIssues[parentKey]; exists {
		// Parent exists - add subtask to it
		parentIssue.SubTasks = append(parentIssue.SubTasks, *subtask)

		// Ensure parent is added to report if not already
		if !parentIssue.AddedToReport {
			g.addTopLevelIssue(parentIssue, epicKey, epicGroups, noEpicIssues)
		}
	} else {
		// Parent not processed yet - fetch and process it
		parentTask, err := g.fetchAndProcessParentTask(ctx, parentKey, lookbackTime)
		if err == nil {
			parentTask.SubTasks = append(parentTask.SubTasks, *subtask)
			processedIssues[parentKey] = parentTask

			// Add parent to report
			g.addTopLevelIssue(parentTask, epicKey, epicGroups, noEpicIssues)
		}
	}

	// Mark subtask as processed (it's included in its parent)
	subtask.AddedToReport = true
}

// ensureEpicGroupExists creates an Epic group if it doesn't exist
func (g *Generator) ensureEpicGroupExists(epicKey string, epicGroups map[string]*EpicGroup) {
	if _, exists := epicGroups[epicKey]; !exists {
		// Fetch Epic details
		epicDetails, err := g.issueService.Get(context.Background(), epicKey, nil, []string{"summary", "status", "issuetype"}, nil)
		if err == nil {
			epicGroups[epicKey] = &EpicGroup{
				EpicKey:     epicKey,
				EpicSummary: epicDetails.Fields.Summary,
				EpicStatus:  epicDetails.Fields.Status.Name,
				EpicURL:     fmt.Sprintf("%s/browse/%s", g.config.JiraHost, epicKey),
				Issues:      []IssueUpdate{},
			}
		}
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
			Key:           issue.Key,
			Summary:       issue.Summary,
			Status:        issue.Status,
			IssueType:     issue.IssueType,
			URL:           issue.URL,
			Updates:       convertUpdates(issue.Updates),
			LastUpdated:   issue.LastUpdated,
			SubTasks:      convertIssueUpdates(issue.SubTasks), // Recursive conversion for sub-tasks
			AddedToReport: issue.AddedToReport,
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
