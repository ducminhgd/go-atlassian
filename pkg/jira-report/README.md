# Jira Report Package

A Go package for generating daily reports from Jira and publishing them to webhooks (e.g., Microsoft Teams).

## Features

- Generate daily reports from Jira issues updated in the last N hours
- Group issues by Epic
- Include comments and worklogs
- Generate reports in both Markdown and HTML formats
- Publish HTML reports to webhooks (Microsoft Teams, Slack, etc.)
- Configurable timezone for timestamps
- Proper HTML escaping for security

## Installation

```bash
go get github.com/ducminhgd/go-atlassian/pkg/jira-report
```

## Usage

### Basic Example - Project + Hours

```go
package main

import (
    "context"
    "fmt"
    "log"

    jirareport "github.com/ducminhgd/go-atlassian/pkg/jira-report"
)

func main() {
    // Create configuration
    config := jirareport.NewConfig()
    config.JiraHost = "https://your-domain.atlassian.net"
    config.JiraUsername = "your-email@example.com"
    config.JiraPassword = "your-api-token"
    config.WebhookURL = "https://your-webhook-url"
    config.Timezone = "America/New_York"

    // Option 1: Use project + hours (default)
    config.QueryType = jirareport.QueryTypeProjectAndHours
    config.JiraProject = "PROJ"
    config.LookbackHours = 24

    // Create generator
    generator, err := jirareport.NewGenerator(config)
    if err != nil {
        log.Fatal(err)
    }

    // Generate report
    report, err := generator.Generate(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    // Print markdown report
    fmt.Println(report.Markdown)

    // Publish to webhook
    publisher := jirareport.NewPublisher(config.WebhookURL)
    publisher.Publish(report.HTML)
}
```

### Using Custom JQL

```go
config := jirareport.NewConfig()
config.JiraHost = "https://your-domain.atlassian.net"
config.JiraPassword = "your-api-token"
config.WebhookURL = "https://your-webhook-url"

// Option 2: Use custom JQL
config.QueryType = jirareport.QueryTypeCustomJQL
config.CustomJQL = "project = PROJ AND assignee = currentUser() AND updated >= -24h"

generator, _ := jirareport.NewGenerator(config)
report, _ := generator.Generate(context.Background())
```

### Using Saved Filter

```go
config := jirareport.NewConfig()
config.JiraHost = "https://your-domain.atlassian.net"
config.JiraPassword = "your-api-token"
config.WebhookURL = "https://your-webhook-url"

// Option 3: Use saved filter ID
config.QueryType = jirareport.QueryTypeFilter
config.FilterID = "12345"

generator, _ := jirareport.NewGenerator(config)
report, _ := generator.Generate(context.Background())
```

### Fluent API - Switching Query Types

```go
generator, _ := jirareport.NewGenerator(config)

// Switch to project + hours
generator.WithProjectAndHours("PROJ", 24)
report1, _ := generator.Generate(ctx)

// Switch to custom JQL
generator.WithCustomJQL("assignee = currentUser() AND updated >= -48h")
report2, _ := generator.Generate(ctx)

// Switch to filter
generator.WithFilter("12345")
report3, _ := generator.Generate(ctx)
```

### Advanced Usage

#### Generate Report Without Publishing

```go
// Create generator
generator, err := jirareport.NewGenerator(config)
if err != nil {
    log.Fatal(err)
}

// Generate report
report, err := generator.Generate(context.Background())
if err != nil {
    log.Fatal(err)
}

// Use the report however you want
fmt.Println("Markdown:", report.Markdown)
fmt.Println("HTML:", report.HTML)
saveToFile(report.HTML)
sendEmail(report.Markdown)
```

#### Publish to Multiple Webhooks

```go
report, _ := generator.Generate(ctx)

// Publish to Teams
teamsPublisher := jirareport.NewPublisher("https://teams-webhook")
teamsPublisher.Publish(report.HTML)

// Publish to Slack
slackPublisher := jirareport.NewPublisher("https://slack-webhook")
slackPublisher.Publish(report.HTML)
```

#### Load Configuration from Environment Variables

```go
func loadConfig() *jirareport.Config {
    config := jirareport.NewConfig()
    
    config.JiraHost = os.Getenv("JIRA_HOST")
    config.JiraUsername = os.Getenv("JIRA_USERNAME")
    config.JiraPassword = os.Getenv("JIRA_PASSWORD")
    config.JiraProject = os.Getenv("JIRA_PROJECT")
    config.WebhookURL = os.Getenv("WEBHOOK_URL")
    
    if tz := os.Getenv("TIMEZONE"); tz != "" {
        config.Timezone = tz
    }
    
    return config
}
```

## Configuration

### Config Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `JiraHost` | string | Yes | - | Jira instance URL (e.g., `https://your-domain.atlassian.net`) |
| `JiraUsername` | string | No | - | Jira username (email). If empty, uses token-only auth |
| `JiraPassword` | string | Yes | - | Jira API token or password |
| `WebhookURL` | string | Yes | - | Webhook URL to post reports to |
| `Timezone` | string | No | `"UTC"` | Timezone for timestamps (e.g., `"America/New_York"`) |
| `QueryType` | QueryType | No | `QueryTypeProjectAndHours` | Type of query to use |
| `JiraProject` | string | Conditional | - | Required if `QueryType` is `QueryTypeProjectAndHours` |
| `LookbackHours` | int | No | `24` | Number of hours to look back (for project + hours query) |
| `CustomJQL` | string | Conditional | - | Required if `QueryType` is `QueryTypeCustomJQL` |
| `FilterID` | string | Conditional | - | Required if `QueryType` is `QueryTypeFilter` |

### Query Types

The package supports three query types:

#### 1. QueryTypeProjectAndHours (Default)
Search for issues in a specific project updated within the last N hours.

```go
config.QueryType = jirareport.QueryTypeProjectAndHours
config.JiraProject = "PROJ"
config.LookbackHours = 24
```

Generates JQL: `project = PROJ AND updated >= -24h ORDER BY updated DESC`

#### 2. QueryTypeCustomJQL
Use a custom JQL query for maximum flexibility.

```go
config.QueryType = jirareport.QueryTypeCustomJQL
config.CustomJQL = "assignee = currentUser() AND status = 'In Progress'"
```

#### 3. QueryTypeFilter
Use a saved Jira filter by its ID.

```go
config.QueryType = jirareport.QueryTypeFilter
config.FilterID = "12345"
```

The generator will fetch the JQL from the saved filter and use it.

### Validation

The `Config.Validate()` method checks that all required fields are set:

```go
if err := config.Validate(); err != nil {
    log.Fatal(err)
}
```

## Report Format

### Markdown Format

```markdown
# Daily Report DD-MMM-YYYY

From last updates in the last 24 hours

## [EPIC-123 In Progress: Epic Summary](https://jira.example.com/browse/EPIC-123)

### [Task | TASK-456 In Progress: Task Summary](https://jira.example.com/browse/TASK-456)

1. 14:30 **John Doe** commented: This is a comment
2. 15:45 **Jane Smith** log work 2h: Worked on implementation

## Anything else

### [Bug | BUG-789 Done: Bug Summary](https://jira.example.com/browse/BUG-789)

1. 10:15 **Alice Johnson** commented: Fixed the issue
```

### HTML Format

The HTML format uses ordered lists with automatic numbering:

```html
<h1>Daily Report DD-MMM-YYYY</h1>
<p>From last updates in the last 24 hours</p>

<ol>
  <li>
    <h2><a href="...">EPIC-123 In Progress: Epic Summary</a></h2>
    <ol>
      <li>
        <h3><a href="...">[Task | TASK-456 In Progress: Task Summary]</a></h3>
        <ol>
          <li>14:30 John Doe commented: This is a comment</li>
          <li>15:45 Jane Smith log work 2h: Worked on implementation</li>
        </ol>
      </li>
    </ol>
  </li>
</ol>
```

## API Reference

### Types

#### Config
Configuration for the report generator.

#### Report
Contains both Markdown and HTML versions of the report.

```go
type Report struct {
    Markdown string
    HTML     string
}
```

#### EpicGroup
Represents a group of issues under an epic.

#### IssueUpdate
Represents an issue with its updates.

#### Update
Represents a single update (comment or worklog).

### Functions

#### NewConfig() *Config
Creates a new Config with default values.

#### NewGenerator(config *Config) (*Generator, error)
Creates a new report generator.

#### NewPublisher(webhookURL string) *Publisher
Creates a new publisher for posting reports to webhooks.

#### NewRunner(config *Config) (*Runner, error)
Creates a new runner that combines generation and publishing.

### Methods

#### Generator.Generate(ctx context.Context) (*Report, error)
Generates the daily report.

#### Publisher.Publish(htmlReport string) error
Publishes the HTML report to the webhook.

#### Runner.Run(ctx context.Context) (*Report, error)
Generates and publishes the report in one call.

## Error Handling

The package defines several error types:

- `ErrMissingJiraHost` - JIRA_HOST is required
- `ErrMissingJiraPassword` - JIRA_PASSWORD is required
- `ErrMissingJiraProject` - JIRA_PROJECT is required
- `ErrMissingWebhookURL` - WEBHOOK_URL is required
- `ErrSearchIssues` - Failed to search issues
- `ErrGenerateReport` - Failed to generate report
- `ErrPostToWebhook` - Failed to post to webhook

Example error handling:

```go
report, err := runner.Run(ctx)
if err != nil {
    if errors.Is(err, jirareport.ErrPostToWebhook) {
        // Handle webhook error
        log.Printf("Failed to post to webhook: %v", err)
        // Report was still generated, can use it
        fmt.Println(report.Markdown)
    } else {
        log.Fatal(err)
    }
}
```

## Examples

See `cmd/jira-daily-report/main.go` for a complete example of how to use this package.

## License

MIT License

