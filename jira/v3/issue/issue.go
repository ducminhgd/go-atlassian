package issue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ducminhgd/go-atlassian/jira/v3/auth"
	"github.com/ducminhgd/go-atlassian/jira/v3/utils"
)

// Service handles communication with the issue related methods
type Service struct {
	client  *http.Client
	baseURL string
	auth    auth.Authenticator
}

// NewService creates a new service instance
func NewService(client *http.Client, baseURL string, auth auth.Authenticator) *Service {
	if client == nil {
		client = http.DefaultClient
	}
	return &Service{
		client:  client,
		baseURL: baseURL,
		auth:    auth,
	}
}

// newRequest creates a new HTTP request
func (s *Service) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(s.baseURL + path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	err = s.auth.AddAuthentication(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// do sends an HTTP request and returns the response
func (s *Service) do(req *http.Request, v interface{}) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}

// SearchJQL searches for issues using JQL (Jira Query Language)
// See: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-search/#api-rest-api-3-search-jql-post
func (s *Service) SearchJQL(ctx context.Context, request JQLSearchRequest) (*JQLSearchResponse, error) {
	if request.JQL == "" {
		return nil, fmt.Errorf("JQL query is required")
	}

	// Set default max results if not specified
	if request.MaxResults <= 0 || request.MaxResults > utils.MAX_RESULTS {
		request.MaxResults = utils.MAX_RESULTS_DEFAULT
	}

	req, err := s.newRequest(ctx, http.MethodPost, ISSUE_SEARCH_JQL_ENDPOINT, request)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	response := new(JQLSearchResponse)
	if err := s.do(req, response); err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return response, nil
}

// Get retrieves an issue by its ID or key
// See: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-get
func (s *Service) Get(ctx context.Context, issueIDOrKey string, expand []string, fields []string, properties []string) (*Issue, error) {
	if issueIDOrKey == "" {
		return nil, fmt.Errorf("issue ID or key is required")
	}

	path := fmt.Sprintf(ISSUE_GET_ENDPOINT, issueIDOrKey)
	params := url.Values{}

	if len(expand) > 0 {
		for _, e := range expand {
			params.Add("expand", e)
		}
	}
	if len(fields) > 0 {
		for _, f := range fields {
			params.Add("fields", f)
		}
	}
	if len(properties) > 0 {
		for _, p := range properties {
			params.Add("properties", p)
		}
	}

	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	req, err := s.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	issue := new(Issue)
	if err := s.do(req, issue); err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return issue, nil
}
