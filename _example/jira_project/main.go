package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ducminhgd/go-atlassian/jira/v3/auth"
	"github.com/ducminhgd/go-atlassian/jira/v3/project"
)

func main() {
	// Load environment variables from .env file
	jiraHost := os.Getenv("JIRA_HOST")
	jiraUsername := os.Getenv("JIRA_USERNAME")
	jiraPassword := os.Getenv("JIRA_PASSWORD")
	if jiraPassword == "" {
		panic("JIRA_PASSWORD is not set in the environment variables")
	}

	authenticator := auth.NewBasicAuth(jiraUsername, jiraPassword)
	client := &http.Client{}
	projectService := project.NewService(client, jiraHost, authenticator)
	// Example: Get all projects
	projects, err := projectService.GetAll(context.Background(), project.ProjectGetAllOpts{})
	if err != nil {
		panic(err)
	}
	// Print project keys
	for _, p := range projects {
		if p.Key == "NXS" {
			fmt.Printf("Project Key: %s\n\t%v", p.Key, p)
		}
	}
}
