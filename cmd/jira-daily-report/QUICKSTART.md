# Quick Start Guide - Jira Daily Report

Get your daily Jira reports up and running in 5 minutes!

## Prerequisites

- Go 1.24 or later
- Jira Cloud account with API access
- Microsoft Teams channel with webhook access

## Step 1: Get Jira API Token

1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Click **"Create API token"**
3. Give it a name (e.g., "Daily Report Bot")
4. **Copy the token** - you won't see it again!

## Step 2: Create Teams Webhook

1. Open Microsoft Teams
2. Go to the channel where you want reports
3. Click **"..."** → **"Connectors"**
4. Find **"Incoming Webhook"** → **"Configure"**
5. Give it a name (e.g., "Jira Daily Report")
6. Optionally upload an icon
7. **Copy the webhook URL**

## Step 3: Configure Environment

Create a `.env` file in the project root:

```bash
# Copy the example file
cp cmd/jira-daily-report/.env.example .env

# Edit with your values
nano .env
```

Fill in your values:
```bash
JIRA_HOST=https://your-company.atlassian.net
JIRA_USERNAME=your-email@company.com
JIRA_PASSWORD=your-api-token-from-step-1
JIRA_PROJECT=PROJ
TEAMS_WEBHOOK_URL=your-webhook-url-from-step-2
REPORT_TIMEZONE=America/New_York
```

## Step 4: Test Run

```bash
# Load environment variables
source .env

# Run the report
go run ./cmd/jira-daily-report/main.go
```

You should see:
1. Markdown report printed to console (for easy reading)
2. Message "Daily report posted successfully!"
3. HTML-formatted report appears in your Teams channel with proper formatting

## Step 5: Automate (Choose One)

### Option A: GitHub Actions (Recommended)

1. Add secrets to your GitHub repository:
   - Go to **Settings** → **Secrets and variables** → **Actions**
   - Add these secrets:
     - `JIRA_HOST`
     - `JIRA_USERNAME`
     - `JIRA_PASSWORD`
     - `JIRA_PROJECT`
     - `TEAMS_WEBHOOK_URL`

2. The workflow is already configured in `.github/workflows/daily-report.yml`

3. It will run automatically at 9 AM UTC on weekdays

4. Test it manually:
   - Go to **Actions** tab
   - Select **"Daily Jira Report"**
   - Click **"Run workflow"**

### Option B: Cron Job

```bash
# Build the binary
go build -o jira-daily-report ./cmd/jira-daily-report

# Move to a permanent location
sudo mv jira-daily-report /usr/local/bin/

# Create a script with environment variables
cat > /usr/local/bin/run-jira-report.sh << 'EOF'
#!/bin/bash
export JIRA_HOST="https://your-company.atlassian.net"
export JIRA_USERNAME="your-email@company.com"
export JIRA_PASSWORD="your-api-token"
export JIRA_PROJECT="PROJ"
export TEAMS_WEBHOOK_URL="your-webhook-url"
export REPORT_TIMEZONE="America/New_York"

/usr/local/bin/jira-daily-report
EOF

# Make it executable
sudo chmod +x /usr/local/bin/run-jira-report.sh

# Add to crontab (runs at 9 AM on weekdays)
crontab -e
# Add this line:
0 9 * * 1-5 /usr/local/bin/run-jira-report.sh >> /var/log/jira-daily-report.log 2>&1
```

### Option C: Docker

```bash
# Build the image
docker build -t jira-daily-report -f cmd/jira-daily-report/Dockerfile .

# Run with environment file
docker run --rm --env-file .env jira-daily-report

# Or run with environment variables
docker run --rm \
  -e JIRA_HOST="https://your-company.atlassian.net" \
  -e JIRA_USERNAME="your-email@company.com" \
  -e JIRA_PASSWORD="your-api-token" \
  -e JIRA_PROJECT="PROJ" \
  -e TEAMS_WEBHOOK_URL="your-webhook-url" \
  -e REPORT_TIMEZONE="America/New_York" \
  jira-daily-report
```

## Troubleshooting

### "JIRA_PASSWORD is required"
- Make sure you've set all required environment variables
- Check that `.env` file exists and is properly formatted
- Try running `source .env` before running the command

### "Failed to search issues"
- Verify your Jira host URL is correct (include `https://`)
- Check your API token is valid
- Ensure your username is your email address
- Verify you have permission to view the project

### "Failed to post to Teams"
- Check the webhook URL is correct
- Verify the webhook hasn't been deleted in Teams
- Try posting a test message to the webhook:
  ```bash
  curl -H "Content-Type: application/json" \
       -d '{"text":"Test message"}' \
       $TEAMS_WEBHOOK_URL
  ```

### "No issues found"
- Check that the project key is correct
- Verify there are issues updated in the last 24 hours
- Try running with a different project that has recent activity

### Report shows wrong timezone
- Check the `REPORT_TIMEZONE` value
- Use standard timezone names like:
  - `America/New_York`
  - `Europe/London`
  - `Asia/Tokyo`
  - `UTC`

## Next Steps

- Customize the report format (see `cmd/jira-daily-report/main.go`)
- Adjust the lookback period (default: 24 hours)
- Add multiple projects (requires code modification)
- Set up monitoring and alerts for failed reports

## Support

For detailed documentation, see:
- [cmd/jira-daily-report/README.md](README.md) - Full documentation
- [DAILY_REPORT_IMPLEMENTATION.md](../../DAILY_REPORT_IMPLEMENTATION.md) - Implementation details

## Example Report

### Console Output (Markdown)

```markdown
# Daily Report 05-Oct-2025

From last updates in the last 24 hours

## [EPIC-123 In Progress: User Authentication System](https://jira.example.com/browse/EPIC-123)

### [Task | AUTH-456 In Progress: Implement OAuth2 flow](https://jira.example.com/browse/AUTH-456)

1. 14:30 John Doe commented: Started working on the OAuth2 implementation
2. 15:45 John Doe log work 2h: Implemented authorization endpoint
3. 16:20 Jane Smith commented: Looks good, please add unit tests

### [Bug | AUTH-457 Done: Fix token expiration](https://jira.example.com/browse/AUTH-457)

1. 10:15 Alice Johnson commented: Fixed the token expiration issue
2. 11:30 Alice Johnson log work 1h 30m: Testing and verification complete

## Anything else

### [Task | MISC-789 In Progress: Update documentation](https://jira.example.com/browse/MISC-789)

1. 09:00 Bob Wilson commented: Updated API documentation
2. 09:45 Bob Wilson log work 45m: Documentation updates
```

### Microsoft Teams (HTML)

The report in Teams will be rendered as a nicely formatted HTML list with:
- Clickable links to all issues and epics
- Proper heading hierarchy (H1, H2, H3)
- Ordered lists for easy navigation
- All special characters properly escaped

Enjoy your automated daily reports! 🎉

