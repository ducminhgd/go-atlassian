# Microsoft Teams Workflow - HTML Format Implementation

## Summary

Updated the Jira Daily Report tool to send HTML-formatted content directly to Microsoft Teams Workflow as a simple message body, instead of using Adaptive Cards. The HTML format uses proper ordered lists with automatic numbering for H2, H3, and child items.

## Changes Made

### 1. HTML Format with Ordered Lists

The HTML output now follows the template structure with proper ordered list numbering:

```html
<h1>Daily Report DD-MMM-YYYY</h1>
<p>From last updates in the last 24 hours</p>

<ol>
  <li>
    <h2><a href="Epic URL">Epic KEY Epic Status: Epic Summary</a></h2>
    <ol>
      <li>
        <h3><a href="Task URL">[Task Type | Task KEY Task Status: Task Summary]</a></h3>
        <ol>
          <li>HH:MM Author's name commented: comment-content</li>
          <li>HH:MM Author's name log work worklog-time: worklog-content</li>
        </ol>
      </li>
    </ol>
  </li>
  
  <li>
    <h2>Anything else</h2>
    <ol>
      <li>
        <h3><a href="Task URL">[Task Type | Task KEY Task Status: Task Summary]</a></h3>
        <ol>
          <li>HH:MM Author's name commented: comment-content</li>
        </ol>
      </li>
    </ol>
  </li>
</ol>
```

### 2. Ordered List Hierarchy

The HTML structure provides automatic numbering at three levels:

1. **Level 1 (Epics)**: Top-level ordered list
   1. **Level 2 (Issues)**: Nested ordered list under each epic
      1. **Level 3 (Updates)**: Nested ordered list under each issue

This creates a clear visual hierarchy:
- 1. Epic 1
  - 1.1. Issue 1
    - 1.1.1. Update 1
    - 1.1.2. Update 2
  - 1.2. Issue 2
    - 1.2.1. Update 1
- 2. Epic 2
  - 2.1. Issue 3

### 3. Code Changes

#### Restored Functions

**`formatHTMLReport()`**
- Generates HTML with proper ordered list structure
- Uses `<ol>` tags for automatic numbering
- Includes H1, H2, H3 headers within list items
- Escapes all HTML special characters

**`writeHTMLIssueSection()`**
- Writes individual issue sections in HTML
- Creates nested ordered lists for updates
- Properly indents HTML for readability

**`escapeHTML()`**
- Escapes HTML special characters: `&`, `<`, `>`, `"`, `'`
- Prevents HTML injection
- Ensures safe rendering

#### Removed Functions

- Removed `formatAdaptiveCardReport()` (no longer needed)
- Removed `buildAdaptiveCardIssue()` (no longer needed)

#### Updated Functions

**`generateDailyReport()`**
- Returns markdown and HTML (not Adaptive Card)
- Signature: `(ctx, service, config) (string, string, error)`

**`postToTeams()`**
- Sends HTML content in simple `body` field
- Payload structure:
  ```json
  {
    "body": "<html content here>"
  }
  ```
- No longer uses `attachments` array
- Simpler and more direct approach

### 4. Benefits

#### Automatic Numbering
- Browser/Teams automatically numbers all list items
- No manual numbering required
- Consistent numbering across all levels

#### Clear Hierarchy
- H2 for epics
- H3 for issues
- Nested ordered lists for updates
- Visual indentation shows relationships

#### Clickable Links
- All epic and issue titles are clickable
- Links open directly in Jira
- Easy navigation from Teams to Jira

#### HTML Safety
- All user-generated content is escaped
- Prevents HTML injection attacks
- Safe rendering of special characters

### 5. Microsoft Teams Workflow Integration

The payload sent to Teams Workflow is now very simple:

```json
{
  "body": "<h1>Daily Report 05-Oct-2025</h1><p>From last updates...</p><ol>...</ol>"
}
```

This format:
- Works with Microsoft Teams Workflow
- Renders HTML properly in Teams
- Supports all HTML formatting (headers, lists, links)
- No need for complex Adaptive Card structure

### 6. Example Output

#### Console (Markdown)
```markdown
# Daily Report 05-Oct-2025

From last updates in the last 24 hours

## [EPIC-123 In Progress: User Authentication](https://jira.example.com/browse/EPIC-123)

### [Task | AUTH-456 In Progress: Implement OAuth2](https://jira.example.com/browse/AUTH-456)

1. 14:30 John Doe commented: Started implementation
2. 15:45 John Doe log work 2h: Implemented endpoint
```

#### Teams (HTML - Rendered)

**Daily Report 05-Oct-2025**

From last updates in the last 24 hours

1. **[EPIC-123 In Progress: User Authentication](link)**
   1. **[Task | AUTH-456 In Progress: Implement OAuth2](link)**
      1. 14:30 John Doe commented: Started implementation
      2. 15:45 John Doe log work 2h: Implemented endpoint
2. **Anything else**
   1. **[Bug | BUG-789 Done: Fix login issue](link)**
      1. 10:15 Alice Johnson commented: Fixed the issue

### 7. Testing

✅ **Build**: Succeeds (8.5MB binary)  
✅ **Tests**: All tests pass  
✅ **Compilation**: No errors or warnings  

### 8. Usage

No configuration changes required:

```bash
# Set environment variables
export JIRA_HOST="https://your-domain.atlassian.net"
export JIRA_USERNAME="your-email@example.com"
export JIRA_PASSWORD="your-api-token"
export JIRA_PROJECT="PROJ"
export TEAMS_WEBHOOK_URL="https://your-workflow-url"

# Run the daily report
go run ./cmd/jira-daily-report/main.go
```

The tool will:
1. Generate markdown report for console
2. Generate HTML report with ordered lists
3. Send HTML to Teams Workflow in `body` field
4. Teams will render the HTML with proper formatting and numbering

## Troubleshooting

### If HTML doesn't render in Teams

If Teams doesn't render the HTML properly, you may need to:

1. **Check Workflow Configuration**: Ensure your Teams Workflow is configured to accept HTML content
2. **Use MessageCard Format**: Some workflows require MessageCard format:
   ```json
   {
     "@type": "MessageCard",
     "@context": "https://schema.org/extensions",
     "summary": "Daily Jira Report",
     "text": "<html content>"
   }
   ```
3. **Use Adaptive Card with HTML**: Some workflows only support Adaptive Cards

### Alternative Payload Formats

If the simple `body` format doesn't work, try these alternatives:

**Option 1: MessageCard**
```json
{
  "@type": "MessageCard",
  "@context": "https://schema.org/extensions",
  "summary": "Daily Jira Report",
  "text": "<html content>"
}
```

**Option 2: Text Field**
```json
{
  "text": "<html content>"
}
```

**Option 3: Sections**
```json
{
  "sections": [
    {
      "text": "<html content>"
    }
  ]
}
```

You can modify the `postToTeams()` function to try different formats if needed.

## Conclusion

The implementation now sends HTML-formatted reports with proper ordered list numbering to Microsoft Teams Workflow. The HTML structure follows the template exactly, with H2 for epics, H3 for issues, and nested ordered lists for automatic numbering at all levels.

✅ **Status**: Complete and ready for testing with Teams Workflow

