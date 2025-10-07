package main

import (
	"context"
	"log"
	"os"
	"strconv"

	jirareport "github.com/ducminhgd/go-atlassian/pkg/jira-report"
)

// main is an example of how to use the jira-report package
func main() {
	// Load configuration from environment variables
	config := loadConfigFromEnv()

	// Create generator
	generator, err := jirareport.NewGenerator(config)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	// Generate report
	report, err := generator.Generate(context.Background())
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	// Publish AdaptiveCard to webhook (default format for Teams)
	publisher := jirareport.NewPublisher(config.WebhookURL)
	if err := publisher.PublishAdaptiveCard(report.AdaptiveCard); err != nil {
		log.Fatalf("Failed to publish report: %v", err)
	}

	log.Println("Daily report posted successfully!")
}

// loadConfigFromEnv loads configuration from environment variables
// This is just an example - you can load config from any source
func loadConfigFromEnv() *jirareport.Config {
	config := jirareport.NewConfig()

	config.JiraHost = os.Getenv("JIRA_HOST")
	config.JiraUsername = os.Getenv("JIRA_USERNAME")
	config.JiraPassword = os.Getenv("JIRA_PASSWORD")
	config.WebhookURL = os.Getenv("TEAMS_WEBHOOK_URL")

	if timezone := os.Getenv("REPORT_TIMEZONE"); timezone != "" {
		config.Timezone = timezone
	}

	// Determine query type based on environment variables
	if filterID := os.Getenv("JIRA_FILTER_ID"); filterID != "" {
		// Option 3: Use saved filter
		config.QueryType = jirareport.QueryTypeFilter
		config.FilterID = filterID
	} else if customJQL := os.Getenv("JIRA_CUSTOM_JQL"); customJQL != "" {
		// Option 2: Use custom JQL
		config.QueryType = jirareport.QueryTypeCustomJQL
		config.CustomJQL = customJQL
	} else {
		// Option 1: Use project + hours (default)
		config.QueryType = jirareport.QueryTypeProjectAndHours
		config.JiraProject = os.Getenv("JIRA_PROJECT")

		if hoursStr := os.Getenv("JIRA_LOOKBACK_HOURS"); hoursStr != "" {
			if hours, err := strconv.Atoi(hoursStr); err == nil {
				config.LookbackHours = hours
			}
		}
	}

	return config
}
