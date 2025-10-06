package jirareport

import "errors"

var (
	// Configuration errors
	ErrMissingJiraHost     = errors.New("JIRA_HOST is required")
	ErrMissingJiraPassword = errors.New("JIRA_PASSWORD is required")
	ErrMissingJiraProject  = errors.New("JIRA_PROJECT is required")
	ErrMissingWebhookURL   = errors.New("WEBHOOK_URL is required")
	ErrMissingCustomJQL    = errors.New("CUSTOM_JQL is required when using custom JQL query type")
	ErrMissingFilterID     = errors.New("FILTER_ID is required when using filter query type")
	ErrInvalidQueryType    = errors.New("invalid query type")

	// Report generation errors
	ErrSearchIssues   = errors.New("failed to search issues")
	ErrGenerateReport = errors.New("failed to generate report")
	ErrPostToWebhook  = errors.New("failed to post to webhook")
)
