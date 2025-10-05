package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ducminhgd/go-atlassian/jira/v3/auth"
	"github.com/ducminhgd/go-atlassian/jira/v3/issue"
)

func main() {
	// Load environment variables
	jiraHost := os.Getenv("JIRA_HOST")
	jiraUsername := os.Getenv("JIRA_USERNAME")
	jiraPassword := os.Getenv("JIRA_PASSWORD")
	if jiraPassword == "" {
		panic("JIRA_PASSWORD is not set in the environment variables")
	}

	// Create authenticator and HTTP client
	authenticator := auth.NewBasicAuth(jiraUsername, jiraPassword)
	client := &http.Client{}
	issueService := issue.NewService(client, jiraHost, authenticator)

	// Example 1: Search for issues using JQL
	fmt.Println("=== Searching for issues using JQL ===")
	searchRequest := issue.JQLSearchRequest{
		JQL:        "project = TEST ORDER BY created DESC",
		MaxResults: 10,
		Fields:     []string{"summary", "status", "assignee", "created"},
	}

	searchResponse, err := issueService.SearchJQL(context.Background(), searchRequest)
	if err != nil {
		fmt.Printf("Error searching issues: %v\n", err)
	} else {
		fmt.Printf("Found %d issues (total: %d)\n", len(searchResponse.Issues), searchResponse.Total)
		for _, issue := range searchResponse.Issues {
			fmt.Printf("- %s: %s (Status: %s)\n", 
				issue.Key, 
				issue.Fields.Summary, 
				issue.Fields.Status.Name)
		}
	}

	// Example 2: Get a specific issue by key
	fmt.Println("\n=== Getting specific issue ===")
	if len(searchResponse.Issues) > 0 {
		issueKey := searchResponse.Issues[0].Key
		specificIssue, err := issueService.Get(
			context.Background(), 
			issueKey, 
			[]string{"changelog"}, // expand changelog
			[]string{"summary", "description", "status", "assignee", "reporter"}, // specific fields
			nil, // no properties
		)
		if err != nil {
			fmt.Printf("Error getting issue %s: %v\n", issueKey, err)
		} else {
			fmt.Printf("Issue: %s\n", specificIssue.Key)
			fmt.Printf("Summary: %s\n", specificIssue.Fields.Summary)
			fmt.Printf("Status: %s\n", specificIssue.Fields.Status.Name)
			if specificIssue.Fields.Assignee.DisplayName != "" {
				fmt.Printf("Assignee: %s\n", specificIssue.Fields.Assignee.DisplayName)
			}
			if specificIssue.Fields.Reporter.DisplayName != "" {
				fmt.Printf("Reporter: %s\n", specificIssue.Fields.Reporter.DisplayName)
			}
		}
	}

	// Example 3: Search with different JQL queries
	fmt.Println("\n=== Different JQL search examples ===")
	
	jqlQueries := []string{
		"assignee = currentUser() AND status != Done",
		"project = TEST AND created >= -7d",
		"priority = High AND status = 'In Progress'",
	}

	for i, jql := range jqlQueries {
		fmt.Printf("\nQuery %d: %s\n", i+1, jql)
		searchReq := issue.JQLSearchRequest{
			JQL:        jql,
			MaxResults: 5,
			Fields:     []string{"summary", "status", "priority"},
		}
		
		resp, err := issueService.SearchJQL(context.Background(), searchReq)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Found %d issues\n", len(resp.Issues))
			for _, issue := range resp.Issues {
				priority := "None"
				if issue.Fields.Priority.Name != "" {
					priority = issue.Fields.Priority.Name
				}
				fmt.Printf("  - %s: %s (Priority: %s, Status: %s)\n", 
					issue.Key, 
					issue.Fields.Summary, 
					priority,
					issue.Fields.Status.Name)
			}
		}
	}
}
