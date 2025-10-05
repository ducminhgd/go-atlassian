# go-atlassian

Atlassian Products client, written in Go

## Features

- **Jira Cloud API v3** support
  - Project management (create, read, update, delete, search)
  - Issue management (search with JQL, get issue details)
  - Authentication (Basic Auth, Token Auth)
- Type-safe API clients with comprehensive error handling
- Full test coverage with unit tests
- Examples for common use cases

## Installation

```bash
go get github.com/ducminhgd/go-atlassian
```

## Authentication

The library supports two authentication methods for Jira:

### Personal Access Token (Recommended)

```go
import "github.com/ducminhgd/go-atlassian/jira/v3/auth"

authenticator := auth.NewTokenAuth("your-personal-access-token")
```

### Basic Authentication

```go
import "github.com/ducminhgd/go-atlassian/jira/v3/auth"

authenticator := auth.NewBasicAuth("your-username", "your-password")
```

Both authentication methods implement the `Authenticator` interface which can be used to add authentication headers to your requests.

## Usage

### Working with Projects

```go
package main

import (
    "context"
    "fmt"
    "net/http"

    "github.com/ducminhgd/go-atlassian/jira/v3/auth"
    "github.com/ducminhgd/go-atlassian/jira/v3/project"
)

func main() {
    // Setup authentication
    authenticator := auth.NewBasicAuth("your-username", "your-api-token")
    client := &http.Client{}

    // Create project service
    projectService := project.NewService(client, "https://your-domain.atlassian.net", authenticator)

    // Get all projects
    projects, err := projectService.GetAll(context.Background(), project.ProjectGetAllOpts{})
    if err != nil {
        panic(err)
    }

    for _, p := range projects {
        fmt.Printf("Project: %s (%s)\n", p.Name, p.Key)
    }
}
```

### Searching Issues with JQL

```go
package main

import (
    "context"
    "fmt"
    "net/http"

    "github.com/ducminhgd/go-atlassian/jira/v3/auth"
    "github.com/ducminhgd/go-atlassian/jira/v3/issue"
)

func main() {
    // Setup authentication
    authenticator := auth.NewBasicAuth("your-username", "your-api-token")
    client := &http.Client{}

    // Create issue service
    issueService := issue.NewService(client, "https://your-domain.atlassian.net", authenticator)

    // Search for issues using JQL
    searchRequest := issue.JQLSearchRequest{
        JQL:        "project = TEST AND status = 'In Progress' ORDER BY created DESC",
        MaxResults: 50,
        Fields:     []string{"summary", "status", "assignee", "created"},
    }

    response, err := issueService.SearchJQL(context.Background(), searchRequest)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d issues (total: %d)\n", len(response.Issues), response.Total)
    for _, issue := range response.Issues {
        fmt.Printf("- %s: %s (Status: %s)\n",
            issue.Key,
            issue.Fields.Summary,
            issue.Fields.Status.Name)
    }
}
```

### Getting a Specific Issue

```go
// Get issue by key
issue, err := issueService.Get(
    context.Background(),
    "TEST-123",
    []string{"changelog"}, // expand options
    []string{"summary", "description", "status", "assignee"}, // fields to retrieve
    nil, // no properties
)
if err != nil {
    panic(err)
}

fmt.Printf("Issue: %s - %s\n", issue.Key, issue.Fields.Summary)
```

## Project Structure

```
jira/v3/
├── auth/           # Authentication implementations
├── project/        # Project API client
├── issue/          # Issue API client
├── responsetypes/  # Common response type definitions
└── utils/          # Utility functions and constants
```

## Testing

Run all tests:

```bash
go test ./jira/v3/...
```

Run tests with coverage:

```bash
go test -cover ./jira/v3/...
```

## Examples

See the `_example/` directory for complete working examples:

- `_example/jira_project/` - Project management examples
- `_example/jira_issue/` - Issue search and retrieval examples

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
