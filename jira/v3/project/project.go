package project

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ducminhgd/go-atlassian/jira/v3/auth"
	"github.com/ducminhgd/go-atlassian/jira/v3/responsetypes"
)

// Service handles communication with the project related methods
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

// do makes a request and decodes the response into v
func (s *Service) do(req *http.Request, v interface{}) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response from API: status=%d, body=%s", resp.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}

// GetAll returns all projects visible to the user
func (s *Service) GetAll(ctx context.Context, opts ProjectGetAllOpts) ([]responsetypes.Project, error) {
	path := PROJECT_LIST_ENDPOINT
	params := url.Values{}

	if opts.Expand != "" {
		params.Add("expand", opts.Expand)
	}
	if opts.Recent > 0 {
		params.Add("recent", strconv.Itoa(opts.Recent))
	}
	if len(opts.Properties) > 0 {
		params.Add("properties", strings.Join(opts.Properties, ","))
	}

	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	req, err := s.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	var projects []responsetypes.Project
	if err := s.do(req, &projects); err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return projects, nil
}

// Get returns the project details for a project
// See: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-projects/#api-rest-api-3-project-projectidorkey-get
func (s *Service) Get(ctx context.Context, projectIDOrKey string, opts ProjectQueryOpts) (*responsetypes.Project, error) {
	path := fmt.Sprintf(PROJECT_DETAIL_ENDPOINT, projectIDOrKey)
	params := url.Values{}

	if opts.Expand != "" {
		params.Add("expand", opts.Expand)
	}
	if len(opts.Properties) > 0 {
		params.Add("properties", strings.Join(opts.Properties, ","))
	}

	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	req, err := s.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	project := new(responsetypes.Project)
	if err := s.do(req, project); err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return project, nil
}

// Create creates a new project
func (s *Service) Create(ctx context.Context, project *responsetypes.Project) (*responsetypes.Project, error) {
	req, err := s.newRequest(ctx, http.MethodPost, PROJECT_LIST_ENDPOINT, project)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	createdProject := new(responsetypes.Project)
	if err := s.do(req, createdProject); err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return createdProject, nil
}

// Update updates an existing project
// See: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-projects/#api-rest-api-3-project-projectidorkey-put
func (s *Service) Update(ctx context.Context, project *responsetypes.Project, opts ProjectUpdateOpts) (*responsetypes.Project, error) {
	// Use key if provided; otherwise fall back to ID
	idOrKey := project.Key
	if idOrKey == "" {
		idOrKey = project.ID
	}
	path := fmt.Sprintf(PROJECT_DETAIL_ENDPOINT, idOrKey)
	if opts.Expand != "" {
		path = fmt.Sprintf("%s?expand=%s", path, url.QueryEscape(opts.Expand))
	}

	req, err := s.newRequest(ctx, http.MethodPut, path, project)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	updatedProject := new(responsetypes.Project)
	if err := s.do(req, updatedProject); err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return updatedProject, nil
}

// Delete deletes a project
func (s *Service) Delete(ctx context.Context, projectIDOrKey string, async bool) error {
	var path, method string
	if async {
		path = fmt.Sprintf(PROJECT_DELETE_ENDPOINT, projectIDOrKey)
		method = http.MethodPost
	} else {
		path = fmt.Sprintf(PROJECT_DETAIL_ENDPOINT, projectIDOrKey)
		method = http.MethodDelete
	}
	req, err := s.newRequest(ctx, method, path, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	if err := s.do(req, nil); err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	return nil
}

// Archive archives a project
func (s *Service) Archive(ctx context.Context, projectIDOrKey string) error {
	path := fmt.Sprintf(PROJECT_ARCHIVE_ENDPOINT, projectIDOrKey)
	req, err := s.newRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	if err := s.do(req, nil); err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	return nil
}

// Search returns a paginated list of projects
func (s *Service) Search(ctx context.Context, startAt, maxResults int, query string) (*responsetypes.ProjectListResponse, error) {
	params := url.Values{}
	params.Add("startAt", strconv.Itoa(startAt))
	params.Add("maxResults", strconv.Itoa(maxResults))
	if query != "" {
		params.Add("query", query)
	}

	path := fmt.Sprintf(PROJECT_SEARCH_ENDPOINT, params.Encode())
	req, err := s.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	response := new(responsetypes.ProjectListResponse)
	if err := s.do(req, response); err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return response, nil
}

// GetRecent returns a list of up to 20 projects recently viewed by the user
func (s *Service) GetRecent(ctx context.Context) ([]responsetypes.Project, error) {
	req, err := s.newRequest(ctx, http.MethodGet, PROJECT_RECENT_ENDPOINT, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	var projects []responsetypes.Project
	if err := s.do(req, &projects); err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return projects, nil
}

// GetAllStatuses returns all statuses associated with a project
func (s *Service) GetAllStatuses(ctx context.Context, projectIDOrKey string) ([]responsetypes.IssueType, error) {
	path := fmt.Sprintf(PROJECT_STATUS_ENDPOINT, projectIDOrKey)
	req, err := s.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	var issueTypes []responsetypes.IssueType
	if err := s.do(req, &issueTypes); err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return issueTypes, nil
}
