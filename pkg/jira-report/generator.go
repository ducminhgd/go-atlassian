package jirareport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

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

	for _, iss := range response.Issues {
		issueUpdate := g.processIssue(iss, lookbackTime)
		if len(issueUpdate.Updates) == 0 {
			continue // Skip issues with no relevant updates
		}

		// Determine epic
		epicKey := ""
		if iss.Fields.Parent.Key != "" {
			// Get parent issue to check if it's an epic
			parentIssue, err := g.issueService.Get(ctx, iss.Fields.Parent.Key, nil, []string{"summary", "status", "issuetype"}, nil)
			if err == nil && parentIssue.Fields.IssueType.Name == "Epic" {
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
			}
		}

		if epicKey != "" {
			epicGroups[epicKey].Issues = append(epicGroups[epicKey].Issues, issueUpdate)
		} else {
			noEpicIssues = append(noEpicIssues, issueUpdate)
		}
	}

	// Generate markdown and HTML reports
	markdownReport := formatMarkdownReport(epicGroups, noEpicIssues, now, g.config.Timezone)
	htmlReport := formatHTMLReport(epicGroups, noEpicIssues, now, g.config.Timezone)

	return &Report{
		Markdown: markdownReport,
		HTML:     htmlReport,
	}, nil
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
		JQL string `json:"jql"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&filterResponse); err != nil {
		return "", fmt.Errorf("failed to decode filter response: %w", err)
	}

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
