# Microsoft Teams Package

This package provides Microsoft Teams integration functionality, specifically for creating and publishing AdaptiveCards to Teams webhooks.

## Features

- **AdaptiveCard Creation**: Build rich, interactive AdaptiveCards with proper Teams formatting
- **Full-Width Support**: Cards automatically stretch to use the full width of Teams windows
- **Jira Report Formatting**: Convert Jira report data into beautifully formatted AdaptiveCards
- **Teams Message Publishing**: Send AdaptiveCards to Teams via webhooks
- **Type Safety**: Strongly typed structures for all AdaptiveCard elements

## Core Types

### AdaptiveCard
The main AdaptiveCard structure with Teams-specific properties:

```go
type AdaptiveCard struct {
    Type    string                `json:"type"`
    Schema  string                `json:"$schema"`
    Version string                `json:"version"`
    Body    []AdaptiveCardElement `json:"body"`
    MSTeams *MSTeamsProperties    `json:"msteams,omitempty"`
}
```

### AdaptiveCardElement
Represents any element that can be added to an AdaptiveCard:

```go
type AdaptiveCardElement struct {
    Type      string                `json:"type"`
    Text      string                `json:"text,omitempty"`
    Size      string                `json:"size,omitempty"`
    Weight    string                `json:"weight,omitempty"`
    Color     string                `json:"color,omitempty"`
    Wrap      bool                  `json:"wrap,omitempty"`
    Spacing   string                `json:"spacing,omitempty"`
    Separator bool                  `json:"separator,omitempty"`
    Items     []AdaptiveCardElement `json:"items,omitempty"`
    Columns   []AdaptiveCardElement `json:"columns,omitempty"`
    Width     any                   `json:"width,omitempty"`
    Actions   []AdaptiveCardAction  `json:"actions,omitempty"`
    URL       string                `json:"url,omitempty"`
    Title     string                `json:"title,omitempty"`
    Style     string                `json:"style,omitempty"`
    Inlines   []AdaptiveCardInline  `json:"inlines,omitempty"`
}
```

## Usage Examples

### Creating a Basic AdaptiveCard

```go
package main

import (
    "encoding/json"
    "fmt"
    
    "github.com/ducminhgd/go-atlassian/internal/msteams"
)

func main() {
    // Create a new AdaptiveCard with Teams defaults
    card := msteams.NewAdaptiveCard()
    
    // Add a title
    card.AddTextBlock("Daily Report", "Large", "Bolder", true)
    
    // Add a subtitle
    card.AddTextBlock("From last updates in the last 24 hours", "Medium", "", true)
    
    // Convert to JSON
    jsonData, _ := json.MarshalIndent(card, "", "  ")
    fmt.Println(string(jsonData))
}
```

### Creating a Jira Report AdaptiveCard

```go
package main

import (
    "time"
    
    "github.com/ducminhgd/go-atlassian/internal/msteams"
)

func main() {
    // Create sample Jira data
    epicGroups := map[string]*msteams.EpicGroup{
        "EPIC-123": {
            EpicKey:     "EPIC-123",
            EpicSummary: "Improve User Experience",
            EpicStatus:  "In Progress",
            EpicURL:     "https://jira.example.com/browse/EPIC-123",
            Issues: []msteams.IssueUpdate{
                {
                    Key:       "TASK-456",
                    Summary:   "Update login page design",
                    Status:    "In Progress",
                    IssueType: "Task",
                    URL:       "https://jira.example.com/browse/TASK-456",
                    Updates: []msteams.Update{
                        {
                            Time:       time.Now().Add(-1 * time.Hour),
                            AuthorName: "John Doe",
                            Type:       "comment",
                            Content:    "Updated the mockups",
                        },
                    },
                },
            },
        },
    }
    
    // Generate AdaptiveCard
    card := msteams.FormatJiraReportAsAdaptiveCard(epicGroups, nil, time.Now(), "UTC")
    
    // Publish to Teams
    publisher := msteams.NewPublisher("https://your-teams-webhook-url")
    err := publisher.PublishAdaptiveCard(card)
    if err != nil {
        panic(err)
    }
}
```

### Publishing to Teams

```go
package main

import (
    "github.com/ducminhgd/go-atlassian/internal/msteams"
)

func main() {
    // Create an AdaptiveCard
    card := msteams.NewAdaptiveCard()
    card.AddTextBlock("Hello Teams!", "Large", "Bolder", true)
    
    // Create publisher and send to Teams
    publisher := msteams.NewPublisher("https://your-teams-webhook-url")
    err := publisher.PublishAdaptiveCard(card)
    if err != nil {
        panic(err)
    }
}
```

## API Reference

### Functions

#### NewAdaptiveCard() AdaptiveCard
Creates a new AdaptiveCard with default Teams settings (full width, version 1.5).

#### FormatJiraReportAsAdaptiveCard(epicGroups, noEpicIssues, reportDate, timezone) AdaptiveCard
Formats Jira report data into a rich AdaptiveCard with proper styling and interactive elements.

#### FormatTeamsMessage(adaptiveCard AdaptiveCard) TeamsMessage
Wraps an AdaptiveCard in a Teams message structure ready for webhook posting.

#### NewPublisher(webhookURL string) *Publisher
Creates a new Teams webhook publisher.

### Methods

#### AdaptiveCard.AddTextBlock(text, size, weight string, wrap bool)
Adds a text block to the AdaptiveCard.

#### AdaptiveCard.AddRichTextBlock(inlines []AdaptiveCardInline, spacing string, separator bool)
Adds a rich text block with inline formatting.

#### AdaptiveCard.AddContainer(items []AdaptiveCardElement, spacing, style string)
Adds a container with child elements.

#### Publisher.PublishAdaptiveCard(adaptiveCard AdaptiveCard) error
Publishes an AdaptiveCard to the Teams webhook.



### Helper Functions

#### CreateTextRun(text, weight, color string) AdaptiveCardInline
Creates a text run for use in RichTextBlock elements.

#### CreateOpenUrlAction(title, url string) AdaptiveCardAction
Creates an action that opens a URL when clicked.

## Teams Integration

This package is specifically designed for Microsoft Teams and includes:

- **Full-width cards**: All cards use `msteams.width: "Full"` for optimal display
- **Teams-compatible JSON**: Proper message structure for Teams webhooks
- **Interactive elements**: Clickable links and actions that work in Teams
- **Rich formatting**: Proper use of Teams-supported AdaptiveCard features

## Error Handling

All functions return appropriate errors for:
- Invalid webhook URLs
- Network failures
- Teams API errors
- JSON marshaling issues

Always check for errors when publishing to Teams:

```go
err := publisher.PublishAdaptiveCard(card)
if err != nil {
    log.Printf("Failed to publish to Teams: %v", err)
}
```
