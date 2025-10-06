package jirareport

// QueryType defines the type of query to use
type QueryType string

const (
	// QueryTypeProjectAndHours uses project key and lookback hours
	QueryTypeProjectAndHours QueryType = "project_hours"
	// QueryTypeCustomJQL uses a custom JQL query
	QueryTypeCustomJQL QueryType = "custom_jql"
	// QueryTypeFilter uses a saved filter ID
	QueryTypeFilter QueryType = "filter"
)

// Config holds the configuration for the daily report
type Config struct {
	JiraHost     string
	JiraUsername string
	JiraPassword string
	WebhookURL   string
	Timezone     string

	// Query configuration - use one of these options
	QueryType QueryType

	// Option 1: Project + Hours
	JiraProject   string
	LookbackHours int

	// Option 2: Custom JQL
	CustomJQL string

	// Option 3: Filter ID
	FilterID string
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		Timezone:      "UTC",
		LookbackHours: 24,
		QueryType:     QueryTypeProjectAndHours,
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.JiraHost == "" {
		return ErrMissingJiraHost
	}
	if c.JiraPassword == "" {
		return ErrMissingJiraPassword
	}
	if c.WebhookURL == "" {
		return ErrMissingWebhookURL
	}

	// Validate query configuration based on type
	switch c.QueryType {
	case QueryTypeProjectAndHours:
		if c.JiraProject == "" {
			return ErrMissingJiraProject
		}
	case QueryTypeCustomJQL:
		if c.CustomJQL == "" {
			return ErrMissingCustomJQL
		}
	case QueryTypeFilter:
		if c.FilterID == "" {
			return ErrMissingFilterID
		}
	default:
		return ErrInvalidQueryType
	}

	return nil
}
