# Query Options Update - Multiple Query Methods

## Summary

Enhanced the Jira Report generator to support three different query methods, allowing users to easily switch between:
1. **Project + Hours** - Search by project key and lookback hours (default)
2. **Custom JQL** - Use any custom JQL query
3. **Saved Filter** - Use a saved Jira filter by ID

Also removed the `runner.go` file and moved its logic into `main.go` for simplicity.

## Changes Made

### 1. Updated `pkg/jira-report/config.go`

Added query type configuration:

```go
type QueryType string

const (
    QueryTypeProjectAndHours QueryType = "project_hours"
    QueryTypeCustomJQL       QueryType = "custom_jql"
    QueryTypeFilter          QueryType = "filter"
)

type Config struct {
    // ... existing fields ...
    
    // Query configuration - use one of these options
    QueryType     QueryType
    
    // Option 1: Project + Hours
    JiraProject   string
    LookbackHours int
    
    // Option 2: Custom JQL
    CustomJQL     string
    
    // Option 3: Filter ID
    FilterID      string
}
```

The `Validate()` method now checks the appropriate fields based on `QueryType`.

### 2. Updated `pkg/jira-report/errors.go`

Added new error types:

```go
ErrMissingCustomJQL = errors.New("CUSTOM_JQL is required when using custom JQL query type")
ErrMissingFilterID  = errors.New("FILTER_ID is required when using filter query type")
ErrInvalidQueryType = errors.New("invalid query type")
```

### 3. Enhanced `pkg/jira-report/generator.go`

#### Added Query Building Logic

**`buildSearchRequest()`** - Builds the search request based on query type:
```go
func (g *Generator) buildSearchRequest(ctx context.Context) (*issue.JQLSearchRequest, error)
```

**`buildProjectAndHoursJQL()`** - Builds JQL for project + hours:
```go
func (g *Generator) buildProjectAndHoursJQL() string {
    return fmt.Sprintf(
        "project = %s AND updated >= -%dh ORDER BY updated DESC",
        g.config.JiraProject,
        g.config.LookbackHours,
    )
}
```

**`getFilterJQL()`** - Retrieves JQL from a saved filter:
```go
func (g *Generator) getFilterJQL(ctx context.Context) (string, error) {
    // Calls Jira API: GET /rest/api/3/filter/{filterID}
    // Returns the JQL from the filter
}
```

#### Added Fluent API Methods

Users can easily switch between query types:

```go
// Switch to project + hours
generator.WithProjectAndHours("PROJ", 24)

// Switch to custom JQL
generator.WithCustomJQL("assignee = currentUser() AND updated >= -48h")

// Switch to filter
generator.WithFilter("12345")
```

### 4. Removed `pkg/jira-report/runner.go`

The runner logic was simple orchestration that's better suited for the example application. Removed the file to simplify the package.

### 5. Updated `cmd/jira-daily-report/main.go`

Moved runner logic into main.go and added support for all three query types:

```go
func main() {
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
    
    // Print markdown report
    fmt.Println(report.Markdown)
    
    // Publish to webhook
    publisher := jirareport.NewPublisher(config.WebhookURL)
    if err := publisher.Publish(report.HTML); err != nil {
        log.Fatalf("Failed to publish report: %v", err)
    }
    
    log.Println("Daily report posted successfully!")
}
```

The `loadConfigFromEnv()` function now determines query type based on environment variables:

```go
// Priority order:
// 1. If JIRA_FILTER_ID is set -> use filter
// 2. Else if JIRA_CUSTOM_JQL is set -> use custom JQL
// 3. Else -> use project + hours (default)
```

### 6. Updated `cmd/jira-daily-report/.env.example`

Added documentation for all three query options:

```bash
# Query Configuration - Choose ONE of the following options:

# Option 1: Project + Hours (Default)
JIRA_PROJECT=PROJ
JIRA_LOOKBACK_HOURS=24

# Option 2: Custom JQL (uncomment to use)
# JIRA_CUSTOM_JQL=project = PROJ AND assignee = currentUser()

# Option 3: Saved Filter (uncomment to use)
# JIRA_FILTER_ID=12345
```

### 7. Updated `pkg/jira-report/README.md`

Added comprehensive documentation with examples for all three query types.

## Usage Examples

### Option 1: Project + Hours (Default)

```bash
export JIRA_HOST="https://your-domain.atlassian.net"
export JIRA_PASSWORD="your-api-token"
export JIRA_PROJECT="PROJ"
export JIRA_LOOKBACK_HOURS="24"
export TEAMS_WEBHOOK_URL="https://webhook-url"

go run ./cmd/jira-daily-report/main.go
```

Generates JQL: `project = PROJ AND updated >= -24h ORDER BY updated DESC`

### Option 2: Custom JQL

```bash
export JIRA_HOST="https://your-domain.atlassian.net"
export JIRA_PASSWORD="your-api-token"
export JIRA_CUSTOM_JQL="assignee = currentUser() AND status = 'In Progress'"
export TEAMS_WEBHOOK_URL="https://webhook-url"

go run ./cmd/jira-daily-report/main.go
```

Uses the exact JQL provided.

### Option 3: Saved Filter

```bash
export JIRA_HOST="https://your-domain.atlassian.net"
export JIRA_PASSWORD="your-api-token"
export JIRA_FILTER_ID="12345"
export TEAMS_WEBHOOK_URL="https://webhook-url"

go run ./cmd/jira-daily-report/main.go
```

Fetches JQL from the saved filter and uses it.

## Programmatic Usage

### Using Config

```go
// Option 1: Project + Hours
config := jirareport.NewConfig()
config.QueryType = jirareport.QueryTypeProjectAndHours
config.JiraProject = "PROJ"
config.LookbackHours = 24

// Option 2: Custom JQL
config := jirareport.NewConfig()
config.QueryType = jirareport.QueryTypeCustomJQL
config.CustomJQL = "assignee = currentUser()"

// Option 3: Filter
config := jirareport.NewConfig()
config.QueryType = jirareport.QueryTypeFilter
config.FilterID = "12345"
```

### Using Fluent API

```go
generator, _ := jirareport.NewGenerator(config)

// Generate report with project + hours
generator.WithProjectAndHours("PROJ", 24)
report1, _ := generator.Generate(ctx)

// Switch to custom JQL
generator.WithCustomJQL("assignee = currentUser()")
report2, _ := generator.Generate(ctx)

// Switch to filter
generator.WithFilter("12345")
report3, _ := generator.Generate(ctx)
```

## Benefits

### 1. Flexibility
Users can choose the query method that best fits their needs:
- **Project + Hours**: Simple, good for daily reports
- **Custom JQL**: Maximum flexibility for complex queries
- **Saved Filter**: Reuse existing filters, share queries across team

### 2. Easy Switching
The fluent API makes it easy to switch between query types without recreating the generator.

### 3. Backward Compatibility
The default behavior (project + hours) remains unchanged, so existing users don't need to update their code.

### 4. Simplified Package
Removing `runner.go` simplifies the package structure. The orchestration logic is better suited for the example application.

## Testing

All tests pass:

```bash
go test ./...
```

Build succeeds:

```bash
go build ./cmd/jira-daily-report
```

## Migration Guide

### For Existing Users

No changes required! The default behavior is unchanged:

```bash
export JIRA_PROJECT="PROJ"
go run ./cmd/jira-daily-report/main.go
```

### To Use Custom JQL

Just set `JIRA_CUSTOM_JQL` instead of `JIRA_PROJECT`:

```bash
export JIRA_CUSTOM_JQL="your custom JQL"
go run ./cmd/jira-daily-report/main.go
```

### To Use Saved Filter

Just set `JIRA_FILTER_ID`:

```bash
export JIRA_FILTER_ID="12345"
go run ./cmd/jira-daily-report/main.go
```

## Implementation Details

### Filter API Call

When using `QueryTypeFilter`, the generator makes an API call to fetch the filter:

```
GET /rest/api/3/filter/{filterID}
```

Response:
```json
{
  "id": "12345",
  "name": "My Filter",
  "jql": "project = PROJ AND assignee = currentUser()",
  ...
}
```

The JQL is extracted and used for the search.

### Priority Order

When loading from environment variables, the priority is:
1. `JIRA_FILTER_ID` (highest priority)
2. `JIRA_CUSTOM_JQL`
3. `JIRA_PROJECT` + `JIRA_LOOKBACK_HOURS` (default)

This allows users to easily switch between methods by setting/unsetting environment variables.

## Conclusion

The generator now supports three flexible query methods, making it suitable for a wide range of use cases. The fluent API makes it easy to switch between methods, and the package structure is simplified by removing the runner.

✅ **Status**: Complete and ready for use

