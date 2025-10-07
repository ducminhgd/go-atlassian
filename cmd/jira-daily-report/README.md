# Jira Daily Report

A tool to generate daily reports from Jira and post them to Microsoft Teams.

## Features

- Collects information from work items updated in the last 24 hours
- Groups issues by Epic
- Includes comments and worklogs from the last 24 hours
- Generates reports in Markdown, HTML, and AdaptiveCard formats
- Posts AdaptiveCard-formatted reports to Microsoft Teams via webhook (default)
- AdaptiveCards support full-width display in Teams for better readability
- Backward compatibility with HTML format

## Configuration

The tool is configured via environment variables:

| Variable | Required | Description | Default |
|----------|----------|-------------|---------|
| `JIRA_HOST` | Yes | Jira instance URL (e.g., `https://your-domain.atlassian.net`) | - |
| `JIRA_USERNAME` | No | Jira username (for Basic Auth) | - |
| `JIRA_PASSWORD` | Yes | Jira password or API token | - |
| `JIRA_PROJECT` | Yes | Jira project key to generate report for | - |
| `TEAMS_WEBHOOK_URL` | Yes | Microsoft Teams incoming webhook URL | - |
| `REPORT_TIMEZONE` | No | Timezone for report timestamps | `UTC` |

## Setup

### 1. Create a Jira API Token

1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Click "Create API token"
3. Give it a name and copy the token
4. Use your email as `JIRA_USERNAME` and the token as `JIRA_PASSWORD`

### 2. Create a Microsoft Teams Webhook

1. In Microsoft Teams, go to the channel where you want to post reports
2. Click the "..." menu and select "Connectors"
3. Find "Incoming Webhook" and click "Configure"
4. Give it a name and optionally upload an image
5. Copy the webhook URL

### 3. Set Environment Variables

Create a `.env` file or set environment variables:

```bash
export JIRA_HOST="https://your-domain.atlassian.net"
export JIRA_USERNAME="your-email@example.com"
export JIRA_PASSWORD="your-api-token"
export JIRA_PROJECT="PROJ"
export TEAMS_WEBHOOK_URL="https://your-org.webhook.office.com/webhookb2/..."
export REPORT_TIMEZONE="America/New_York"
```

## Usage

### Build and Run

```bash
# Build the binary
go build -o jira-daily-report ./cmd/jira-daily-report

# Run with environment variables
./jira-daily-report
```

### Run Directly

```bash
go run ./cmd/jira-daily-report/main.go
```

### Docker

```bash
# Build Docker image
docker build -t jira-daily-report -f cmd/jira-daily-report/Dockerfile .

# Run with environment variables
docker run --rm \
  -e JIRA_HOST="https://your-domain.atlassian.net" \
  -e JIRA_USERNAME="your-email@example.com" \
  -e JIRA_PASSWORD="your-api-token" \
  -e JIRA_PROJECT="PROJ" \
  -e TEAMS_WEBHOOK_URL="https://your-webhook-url" \
  jira-daily-report
```

### Schedule with Cron

Add to your crontab to run daily at 9 AM:

```bash
0 9 * * * cd /path/to/project && /path/to/jira-daily-report >> /var/log/jira-daily-report.log 2>&1
```

### Schedule with GitHub Actions

Create `.github/workflows/daily-report.yml`:

```yaml
name: Daily Jira Report

on:
  schedule:
    - cron: '0 9 * * *'  # Run at 9 AM UTC daily
  workflow_dispatch:  # Allow manual trigger

jobs:
  report:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Run Daily Report
        env:
          JIRA_HOST: ${{ secrets.JIRA_HOST }}
          JIRA_USERNAME: ${{ secrets.JIRA_USERNAME }}
          JIRA_PASSWORD: ${{ secrets.JIRA_PASSWORD }}
          JIRA_PROJECT: ${{ secrets.JIRA_PROJECT }}
          TEAMS_WEBHOOK_URL: ${{ secrets.TEAMS_WEBHOOK_URL }}
          REPORT_TIMEZONE: 'America/New_York'
        run: go run ./cmd/jira-daily-report/main.go
```

## Report Format

The tool generates reports in two formats:

### Console Output (Markdown)

The console displays a markdown-formatted report for easy reading:

```markdown
# Daily Report DD-MMM-YYYY

From last updates in the last 24 hours

## [EPIC-123 In Progress: Epic Summary](https://jira.example.com/browse/EPIC-123)

### [Task | TASK-456 In Progress: Task Summary](https://jira.example.com/browse/TASK-456)

1. 14:30 John Doe commented: This is a comment
2. 15:45 Jane Smith log work 2h: Worked on implementation
3. 16:20 John Doe commented: Another update

## Anything else

### [Bug | BUG-789 Done: Bug Summary](https://jira.example.com/browse/BUG-789)

1. 10:15 Alice Johnson commented: Fixed the issue
2. 11:30 Alice Johnson log work 1h 30m: Testing and verification
```

### Microsoft Teams (HTML)

The report posted to Teams uses HTML format for better rendering:

```html
<h1>Daily Report DD-MMM-YYYY</h1>
<p>From last updates in the last 24 hours</p>

<ol>
  <li>
    <h2><a href="https://jira.example.com/browse/EPIC-123">EPIC-123 In Progress: Epic Summary</a></h2>
    <ol>
      <li>
        <h3><a href="https://jira.example.com/browse/TASK-456">[Task | TASK-456 In Progress: Task Summary]</a></h3>
        <ol>
          <li>14:30 John Doe commented: This is a comment</li>
          <li>15:45 Jane Smith log work 2h: Worked on implementation</li>
        </ol>
      </li>
    </ol>
  </li>
  <li>
    <h2>Anything else</h2>
    <ol>
      <li>
        <h3><a href="https://jira.example.com/browse/BUG-789">[Bug | BUG-789 Done: Bug Summary]</a></h3>
        <ol>
          <li>10:15 Alice Johnson commented: Fixed the issue</li>
        </ol>
      </li>
    </ol>
  </li>
</ol>
```

## Troubleshooting

### Authentication Errors

- Make sure you're using an API token, not your password
- Verify your email and token are correct
- Check that your Jira instance URL is correct (include `https://`)

### No Issues Found

- Verify the project key is correct
- Check that there are actually issues updated in the last 24 hours
- Ensure your user has permission to view the project

### Teams Webhook Errors

- Verify the webhook URL is correct and active
- Check that the webhook hasn't been deleted in Teams
- Ensure the payload format is correct

## Development

### Run Tests

```bash
go test ./cmd/jira-daily-report/...
```

### Build for Multiple Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o jira-daily-report-linux ./cmd/jira-daily-report

# macOS
GOOS=darwin GOARCH=amd64 go build -o jira-daily-report-macos ./cmd/jira-daily-report

# Windows
GOOS=windows GOARCH=amd64 go build -o jira-daily-report.exe ./cmd/jira-daily-report
```

## License

MIT License

