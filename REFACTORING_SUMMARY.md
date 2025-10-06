# Refactoring Summary: Jira Daily Report

## Overview

Refactored the Jira Daily Report implementation to separate business logic from the example application. The business logic is now in a reusable package `pkg/jira-report`, while `cmd/jira-daily-report/main.go` serves as a simple example of how to use the package.

## Changes Made

### 1. Created `pkg/jira-report` Package

The package contains all the business logic for generating and publishing Jira reports.

#### Package Structure

```
pkg/jira-report/
├── config.go              # Configuration management
├── errors.go              # Error definitions
├── types.go               # Data structures
├── generator.go           # Report generation logic
├── formatter_markdown.go  # Markdown formatting
├── formatter_html.go      # HTML formatting
├── publisher.go           # Webhook publishing
├── runner.go              # Orchestration
└── README.md              # Package documentation
```

#### Files Created

**`config.go`**
- `Config` struct with all configuration fields
- `NewConfig()` constructor with defaults
- `Validate()` method for configuration validation

**`errors.go`**
- Predefined error variables for common errors
- `ErrMissingJiraHost`, `ErrMissingJiraPassword`, etc.
- `ErrSearchIssues`, `ErrGenerateReport`, `ErrPostToWebhook`

**`types.go`**
- `EpicGroup` - Group of issues under an epic
- `IssueUpdate` - Issue with its updates
- `Update` - Single update (comment or worklog)
- `Report` - Contains both Markdown and HTML

**`generator.go`**
- `Generator` struct - Main report generator
- `NewGenerator()` - Creates generator with Jira client
- `Generate()` - Generates the report
- `processIssue()` - Processes individual issues
- `extractTextFromBody()` - Extracts text from ADF
- `extractTextFromADF()` - Recursively extracts text
- `truncateText()` - Truncates long text

**`formatter_markdown.go`**
- `formatMarkdownReport()` - Formats report in Markdown
- `writeMarkdownIssueSection()` - Writes issue sections

**`formatter_html.go`**
- `formatHTMLReport()` - Formats report in HTML
- `writeHTMLIssueSection()` - Writes issue sections
- `escapeHTML()` - Escapes HTML special characters

**`publisher.go`**
- `Publisher` struct - Handles webhook publishing
- `NewPublisher()` - Creates publisher
- `Publish()` - Posts HTML to webhook

**`runner.go`**
- `Runner` struct - Orchestrates generation and publishing
- `NewRunner()` - Creates runner
- `Run()` - Generates and publishes in one call

**`README.md`**
- Complete package documentation
- Usage examples
- API reference
- Configuration guide

### 2. Simplified `cmd/jira-daily-report/main.go`

The main.go file is now a simple example (51 lines vs 496 lines):

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    jirareport "github.com/ducminhgd/go-atlassian/pkg/jira-report"
)

func main() {
    // Load configuration from environment variables
    config := loadConfigFromEnv()

    // Create and run the report generator
    runner, err := jirareport.NewRunner(config)
    if err != nil {
        log.Fatalf("Failed to create runner: %v", err)
    }

    // Generate and publish report
    report, err := runner.Run(context.Background())
    if err != nil {
        log.Fatalf("Failed to run report: %v", err)
    }

    // Print markdown report to console
    fmt.Println(report.Markdown)

    log.Println("Daily report posted successfully!")
}

func loadConfigFromEnv() *jirareport.Config {
    config := jirareport.NewConfig()
    
    config.JiraHost = os.Getenv("JIRA_HOST")
    config.JiraUsername = os.Getenv("JIRA_USERNAME")
    config.JiraPassword = os.Getenv("JIRA_PASSWORD")
    config.JiraProject = os.Getenv("JIRA_PROJECT")
    config.WebhookURL = os.Getenv("TEAMS_WEBHOOK_URL")
    
    if timezone := os.Getenv("REPORT_TIMEZONE"); timezone != "" {
        config.Timezone = timezone
    }

    return config
}
```

## Benefits

### 1. Separation of Concerns
- **Business Logic**: In `pkg/jira-report` package
- **Example Usage**: In `cmd/jira-daily-report/main.go`
- Clear separation makes code easier to understand and maintain

### 2. Reusability
- The package can be imported and used in other projects
- Different applications can use the same report generation logic
- Easy to create custom implementations

### 3. Testability
- Business logic can be unit tested independently
- Mock configurations and dependencies easily
- Test different scenarios without running the full application

### 4. Flexibility
- Users can choose to use:
  - `Runner` - Full orchestration (generate + publish)
  - `Generator` - Just generate reports
  - `Publisher` - Just publish to webhook
- Load configuration from any source (env vars, files, databases, etc.)
- Customize behavior by implementing interfaces

### 5. Maintainability
- Smaller, focused files are easier to maintain
- Clear responsibilities for each component
- Easy to add new features or modify existing ones

### 6. Documentation
- Package-level documentation in README.md
- Clear API reference
- Usage examples for common scenarios

## Usage Examples

### Basic Usage (as in main.go)

```go
config := jirareport.NewConfig()
config.JiraHost = "https://your-domain.atlassian.net"
config.JiraPassword = "your-api-token"
config.JiraProject = "PROJ"
config.WebhookURL = "https://webhook-url"

runner, _ := jirareport.NewRunner(config)
report, _ := runner.Run(context.Background())
fmt.Println(report.Markdown)
```

### Generate Only (No Publishing)

```go
generator, _ := jirareport.NewGenerator(config)
report, _ := generator.Generate(context.Background())

// Use the report however you want
saveToFile(report.HTML)
sendEmail(report.Markdown)
```

### Custom Publishing

```go
generator, _ := jirareport.NewGenerator(config)
report, _ := generator.Generate(context.Background())

// Publish to multiple destinations
publisher1 := jirareport.NewPublisher("https://teams-webhook")
publisher2 := jirareport.NewPublisher("https://slack-webhook")

publisher1.Publish(report.HTML)
publisher2.Publish(report.HTML)
```

### Custom Configuration Source

```go
// Load from JSON file
config := loadConfigFromJSON("config.json")

// Load from database
config := loadConfigFromDB(db, "report-config")

// Load from command-line flags
config := loadConfigFromFlags()

runner, _ := jirareport.NewRunner(config)
report, _ := runner.Run(context.Background())
```

## Migration Guide

### For Existing Users

No changes required! The `cmd/jira-daily-report` still works the same way:

```bash
export JIRA_HOST="https://your-domain.atlassian.net"
export JIRA_PASSWORD="your-api-token"
export JIRA_PROJECT="PROJ"
export TEAMS_WEBHOOK_URL="https://webhook-url"

go run ./cmd/jira-daily-report/main.go
```

### For Developers

If you want to use the package in your own code:

```go
import jirareport "github.com/ducminhgd/go-atlassian/pkg/jira-report"

// Create config
config := jirareport.NewConfig()
// ... set config fields ...

// Use the package
runner, err := jirareport.NewRunner(config)
report, err := runner.Run(context.Background())
```

## Testing

All existing tests still pass:

```bash
go test ./...
```

Output:
```
?   	github.com/ducminhgd/go-atlassian/cmd/jira-daily-report	[no test files]
ok  	github.com/ducminhgd/go-atlassian/jira/v3/auth	(cached)
ok  	github.com/ducminhgd/go-atlassian/jira/v3/issue	(cached)
ok  	github.com/ducminhgd/go-atlassian/jira/v3/project	(cached)
?   	github.com/ducminhgd/go-atlassian/jira/v3/responsetypes	[no test files]
?   	github.com/ducminhgd/go-atlassian/jira/v3/utils	[no test files]
?   	github.com/ducminhgd/go-atlassian/pkg/jira-report	[no test files]
```

Build succeeds:
```bash
go build ./cmd/jira-daily-report
```

## Future Enhancements

With this refactoring, it's now easier to add:

1. **Unit Tests** - Test each component independently
2. **Multiple Output Formats** - Add PDF, JSON, etc.
3. **Multiple Webhooks** - Publish to multiple destinations
4. **Custom Filters** - Filter issues by assignee, labels, etc.
5. **Scheduled Reports** - Add scheduling logic
6. **Report Templates** - Customizable report templates
7. **Metrics** - Add statistics and metrics to reports

## Conclusion

The refactoring successfully separates business logic from the example application, making the code more maintainable, testable, and reusable. The `pkg/jira-report` package can now be used in any Go application that needs to generate Jira reports.

✅ **Status**: Complete and ready for use

