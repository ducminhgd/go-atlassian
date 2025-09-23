package issue

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ducminhgd/go-atlassian/jira/v3/auth"
)

func TestService_SearchJQL(t *testing.T) {
	// Mock response
	mockResponse := JQLSearchResponse{
		IsLast: true,
		Issues: []Issue{
			{
				ID:  "10001",
				Key: "TEST-1",
				Fields: IssueFields{
					Summary: "Test issue",
					Status: StatusDetails{
						Name: "Open",
						ID:   "1",
					},
				},
			},
		},
		MaxResults: 50,
		StartAt:    0,
		Total:      1,
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/rest/api/3/search/jql" {
			t.Errorf("Expected path /rest/api/3/search/jql, got %s", r.URL.Path)
		}

		// Verify request body
		var request JQLSearchRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if request.JQL != "project = TEST" {
			t.Errorf("Expected JQL 'project = TEST', got '%s'", request.JQL)
		}

		// Send mock response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create service
	auth := auth.NewBasicAuth("test", "test")
	service := NewService(nil, server.URL, auth)

	// Test SearchJQL
	request := JQLSearchRequest{
		JQL:        "project = TEST",
		MaxResults: 50,
	}

	response, err := service.SearchJQL(context.Background(), request)
	if err != nil {
		t.Fatalf("SearchJQL failed: %v", err)
	}

	// Verify response
	if response.IsLast != true {
		t.Errorf("Expected IsLast to be true, got %v", response.IsLast)
	}
	if len(response.Issues) != 1 {
		t.Errorf("Expected 1 issue, got %d", len(response.Issues))
	}
	if response.Issues[0].Key != "TEST-1" {
		t.Errorf("Expected issue key 'TEST-1', got '%s'", response.Issues[0].Key)
	}
}

func TestService_SearchJQL_EmptyJQL(t *testing.T) {
	auth := auth.NewBasicAuth("test", "test")
	service := NewService(nil, "http://example.com", auth)

	request := JQLSearchRequest{
		JQL: "",
	}

	_, err := service.SearchJQL(context.Background(), request)
	if err == nil {
		t.Error("Expected error for empty JQL, got nil")
	}
	if err.Error() != "JQL query is required" {
		t.Errorf("Expected 'JQL query is required' error, got '%s'", err.Error())
	}
}

func TestService_SearchJQL_MaxResultsLimit(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request JQLSearchRequest
		json.NewDecoder(r.Body).Decode(&request)
		
		// Verify max results is capped at 100
		if request.MaxResults != 100 {
			t.Errorf("Expected MaxResults to be capped at 100, got %d", request.MaxResults)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JQLSearchResponse{})
	}))
	defer server.Close()

	auth := auth.NewBasicAuth("test", "test")
	service := NewService(nil, server.URL, auth)

	request := JQLSearchRequest{
		JQL:        "project = TEST",
		MaxResults: 200, // Should be capped at 100
	}

	_, err := service.SearchJQL(context.Background(), request)
	if err != nil {
		t.Fatalf("SearchJQL failed: %v", err)
	}
}

func TestService_Get(t *testing.T) {
	// Mock response
	mockIssue := Issue{
		ID:  "10001",
		Key: "TEST-1",
		Fields: IssueFields{
			Summary: "Test issue",
			Status: StatusDetails{
				Name: "Open",
				ID:   "1",
			},
		},
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1" {
			t.Errorf("Expected path /rest/api/3/issue/TEST-1, got %s", r.URL.Path)
		}

		// Send mock response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockIssue)
	}))
	defer server.Close()

	// Create service
	auth := auth.NewBasicAuth("test", "test")
	service := NewService(nil, server.URL, auth)

	// Test Get
	issue, err := service.Get(context.Background(), "TEST-1", nil, nil, nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Verify response
	if issue.Key != "TEST-1" {
		t.Errorf("Expected issue key 'TEST-1', got '%s'", issue.Key)
	}
	if issue.Fields.Summary != "Test issue" {
		t.Errorf("Expected summary 'Test issue', got '%s'", issue.Fields.Summary)
	}
}

func TestService_Get_EmptyIssueKey(t *testing.T) {
	auth := auth.NewBasicAuth("test", "test")
	service := NewService(nil, "http://example.com", auth)

	_, err := service.Get(context.Background(), "", nil, nil, nil)
	if err == nil {
		t.Error("Expected error for empty issue key, got nil")
	}
	if err.Error() != "issue ID or key is required" {
		t.Errorf("Expected 'issue ID or key is required' error, got '%s'", err.Error())
	}
}
