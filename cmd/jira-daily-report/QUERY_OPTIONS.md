# Query Options Guide

The Jira Daily Report tool supports three different query methods. Choose the one that best fits your needs.

## Option 1: Project + Hours (Default)

Search for issues in a specific project updated within the last N hours.

### Environment Variables

```bash
JIRA_PROJECT=PROJ
JIRA_LOOKBACK_HOURS=24  # Optional, defaults to 24
```

### Example

```bash
export JIRA_HOST="https://your-domain.atlassian.net"
export JIRA_PASSWORD="your-api-token"
export JIRA_PROJECT="MYPROJECT"
export JIRA_LOOKBACK_HOURS="48"
export TEAMS_WEBHOOK_URL="https://webhook-url"

go run ./cmd/jira-daily-report/main.go
```

### Generated JQL

```
project = MYPROJECT AND updated >= -48h ORDER BY updated DESC
```

### Use Cases

- Daily team reports
- Project status updates
- Simple, consistent queries

---

## Option 2: Custom JQL

Use any custom JQL query for maximum flexibility.

### Environment Variables

```bash
JIRA_CUSTOM_JQL="your custom JQL query"
```

### Example

```bash
export JIRA_HOST="https://your-domain.atlassian.net"
export JIRA_PASSWORD="your-api-token"
export JIRA_CUSTOM_JQL="assignee = currentUser() AND status = 'In Progress' AND updated >= -24h"
export TEAMS_WEBHOOK_URL="https://webhook-url"

go run ./cmd/jira-daily-report/main.go
```

### Use Cases

- Personal reports (assignee = currentUser())
- Specific status tracking
- Complex filters with multiple conditions
- Cross-project queries

### JQL Examples

**My issues updated today:**
```
assignee = currentUser() AND updated >= -24h
```

**All bugs in multiple projects:**
```
project in (PROJ1, PROJ2) AND type = Bug AND status != Done
```

**Issues by specific team:**
```
assignee in membersOf("team-developers") AND updated >= -48h
```

**High priority issues:**
```
project = PROJ AND priority in (Highest, High) AND status != Done
```

---

## Option 3: Saved Filter

Use a saved Jira filter by its ID. This is useful for:
- Reusing existing filters
- Sharing queries across team
- Complex filters managed in Jira UI

### Environment Variables

```bash
JIRA_FILTER_ID=12345
```

### Finding Filter ID

1. Go to Jira and open your filter
2. Look at the URL: `https://your-domain.atlassian.net/issues/?filter=12345`
3. The number after `filter=` is your Filter ID

### Example

```bash
export JIRA_HOST="https://your-domain.atlassian.net"
export JIRA_PASSWORD="your-api-token"
export JIRA_FILTER_ID="12345"
export TEAMS_WEBHOOK_URL="https://webhook-url"

go run ./cmd/jira-daily-report/main.go
```

### Use Cases

- Reuse existing team filters
- Complex filters managed in Jira UI
- Filters shared across multiple tools
- Filters with permissions/visibility rules

---

## Priority Order

When multiple options are set, the tool uses this priority:

1. **JIRA_FILTER_ID** (highest priority)
2. **JIRA_CUSTOM_JQL**
3. **JIRA_PROJECT** (default)

### Example

```bash
# This will use the filter (highest priority)
export JIRA_PROJECT="PROJ"
export JIRA_CUSTOM_JQL="assignee = currentUser()"
export JIRA_FILTER_ID="12345"
```

To switch methods, simply unset the higher priority variables:

```bash
# Switch from filter to custom JQL
unset JIRA_FILTER_ID

# Switch from custom JQL to project
unset JIRA_CUSTOM_JQL
```

---

## Complete Examples

### Example 1: Daily Team Report

```bash
#!/bin/bash
export JIRA_HOST="https://mycompany.atlassian.net"
export JIRA_PASSWORD="my-api-token"
export JIRA_PROJECT="BACKEND"
export JIRA_LOOKBACK_HOURS="24"
export TEAMS_WEBHOOK_URL="https://mycompany.webhook.office.com/..."
export REPORT_TIMEZONE="America/New_York"

./jira-daily-report
```

### Example 2: Personal Weekly Report

```bash
#!/bin/bash
export JIRA_HOST="https://mycompany.atlassian.net"
export JIRA_PASSWORD="my-api-token"
export JIRA_CUSTOM_JQL="assignee = currentUser() AND updated >= -7d ORDER BY updated DESC"
export TEAMS_WEBHOOK_URL="https://mycompany.webhook.office.com/..."
export REPORT_TIMEZONE="America/New_York"

./jira-daily-report
```

### Example 3: Sprint Report Using Filter

```bash
#!/bin/bash
export JIRA_HOST="https://mycompany.atlassian.net"
export JIRA_PASSWORD="my-api-token"
export JIRA_FILTER_ID="10234"  # "Current Sprint" filter
export TEAMS_WEBHOOK_URL="https://mycompany.webhook.office.com/..."
export REPORT_TIMEZONE="America/New_York"

./jira-daily-report
```

---

## Troubleshooting

### "JIRA_PROJECT is required" error

You're using the default query type but haven't set `JIRA_PROJECT`. Either:
- Set `JIRA_PROJECT` environment variable, or
- Use `JIRA_CUSTOM_JQL` or `JIRA_FILTER_ID` instead

### "CUSTOM_JQL is required" error

You set `JIRA_CUSTOM_JQL` to an empty string. Either:
- Provide a valid JQL query, or
- Unset `JIRA_CUSTOM_JQL` to use project + hours

### "FILTER_ID is required" error

You set `JIRA_FILTER_ID` to an empty string. Either:
- Provide a valid filter ID, or
- Unset `JIRA_FILTER_ID` to use another method

### "failed to get filter" error

The filter ID doesn't exist or you don't have permission. Check:
- Filter ID is correct
- You have permission to view the filter
- The filter exists in your Jira instance

---

## Tips

### Testing JQL Queries

Before using a custom JQL in the tool, test it in Jira:
1. Go to Issues → Search for issues
2. Switch to Advanced search
3. Enter your JQL
4. Verify the results

### Creating Reusable Filters

1. Create and test your JQL in Jira
2. Save it as a filter
3. Note the filter ID
4. Use `JIRA_FILTER_ID` in the tool

### Scheduling Reports

Use cron or GitHub Actions to schedule reports:

```cron
# Daily report at 9 AM
0 9 * * * /path/to/jira-daily-report
```

### Multiple Reports

Create different scripts for different reports:

```bash
# daily-team-report.sh
export JIRA_PROJECT="BACKEND"
./jira-daily-report

# weekly-personal-report.sh
export JIRA_CUSTOM_JQL="assignee = currentUser() AND updated >= -7d"
./jira-daily-report

# sprint-report.sh
export JIRA_FILTER_ID="10234"
./jira-daily-report
```

---

## See Also

- [Jira JQL Documentation](https://support.atlassian.com/jira-software-cloud/docs/use-advanced-search-with-jira-query-language-jql/)
- [JQL Functions Reference](https://support.atlassian.com/jira-software-cloud/docs/jql-functions/)
- [Creating and Sharing Filters](https://support.atlassian.com/jira-software-cloud/docs/save-your-search-as-a-filter/)

