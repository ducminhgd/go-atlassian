# go-atlassian

Atlassian Products client, written in Go

## Authentication

The library supports two authentication methods for Jira:

### Personal Access Token (Recommended)

```go
auth := auth.NewTokenAuth("your-personal-access-token")
```

### Basic Authentication

```go
auth := auth.NewBasicAuth("your-username", "your-password")
```

Both authentication methods implement the `Authenticator` interface which can be used to add authentication headers to your requests.
