package project

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ducminhgd/go-atlassian/jira/v3/auth"
	"github.com/ducminhgd/go-atlassian/jira/v3/responsetypes"
)

func TestGetAll(t *testing.T) {
	tests := []struct {
		name       string
		opts       ProjectGetAllOpts
		wantURL    string
		wantErr    bool
		statusCode int
	}{
		{
			name:       "no options",
			opts:       ProjectGetAllOpts{},
			wantURL:    "/rest/api/3/project",
			statusCode: http.StatusOK,
		},
		{
			name:       "server error",
			opts:       ProjectGetAllOpts{},
			wantURL:    "/rest/api/3/project",
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "unauthorized",
			opts:       ProjectGetAllOpts{},
			wantURL:    "/rest/api/3/project",
			wantErr:    true,
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "with options",
			opts: ProjectGetAllOpts{
				Expand:     "description,lead",
				Properties: []string{"prop1", "prop2"},
				Recent:     10,
			},
			wantURL:    "/rest/api/3/project?expand=description%2Clead&properties=prop1%2Cprop2&recent=10",
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				gotURL := r.URL.Path
				if r.URL.RawQuery != "" {
					gotURL += "?" + r.URL.RawQuery
				}
				if gotURL != tt.wantURL {
					t.Errorf("URL = %v, want %v", gotURL, tt.wantURL)
				}
				if r.Method != "GET" {
					t.Errorf("Method = %v, want GET", r.Method)
				}

				// Set content type header
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				if tt.statusCode >= 400 {
					if err := json.NewEncoder(w).Encode(map[string]interface{}{
						"errorMessages": []string{"Error retrieving projects"},
						"errors": map[string]string{
							"base": "Server error occurred",
						},
					}); err != nil {
						t.Errorf("Failed to encode error response: %v", err)
					}
					return
				}

				// Return mock response
				projects := []responsetypes.Project{
					{
						ID:   "10000",
						Key:  "TEST",
						Name: "Test Project",
					},
				}
				err := json.NewEncoder(w).Encode(projects)
				if err != nil {
					t.Errorf("Failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			// Create client
			client := &http.Client{}
			basicAuth := auth.NewBasicAuth("testuser", "secret123")
			service := NewService(client, server.URL, basicAuth)

			// Make request
			projects, err := service.GetAll(context.Background(), tt.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify response
				if len(projects) != 1 {
					t.Errorf("GetAll() got %d projects, want 1", len(projects))
				}
				if projects[0].Key != "TEST" {
					t.Errorf("GetAll() project key = %v, want TEST", projects[0].Key)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		opts       ProjectQueryOpts
		projectKey string
		wantURL    string
		wantErr    bool
		statusCode int
	}{
		{
			name:       "success - no options",
			opts:       ProjectQueryOpts{},
			projectKey: "TEST",
			wantURL:    "/rest/api/3/project/TEST",
			statusCode: http.StatusOK,
		},
		{
			name: "success - with expand",
			opts: ProjectQueryOpts{
				Expand: "description,lead",
			},
			projectKey: "TEST",
			wantURL:    "/rest/api/3/project/TEST?expand=description%2Clead",
			statusCode: http.StatusOK,
		},
		{
			name: "success - with properties",
			opts: ProjectQueryOpts{
				Properties: []string{"prop1", "prop2"},
			},
			projectKey: "TEST",
			wantURL:    "/rest/api/3/project/TEST?properties=prop1%2Cprop2",
			statusCode: http.StatusOK,
		},
		{
			name:       "error - project not found",
			opts:       ProjectQueryOpts{},
			projectKey: "NOTFOUND",
			wantURL:    "/rest/api/3/project/NOTFOUND",
			wantErr:    true,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "error - invalid project key",
			opts:       ProjectQueryOpts{},
			projectKey: "!@#$%",
			wantURL:    "/rest/api/3/project/!@#$%",
			wantErr:    true,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "error - unauthorized",
			opts:       ProjectQueryOpts{},
			projectKey: "TEST",
			wantURL:    "/rest/api/3/project/TEST",
			wantErr:    true,
			statusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request URL
				gotURL := r.URL.Path
				if r.URL.RawQuery != "" {
					gotURL += "?" + r.URL.RawQuery
				}
				if gotURL != tt.wantURL {
					t.Errorf("URL = %v, want %v", gotURL, tt.wantURL)
				}

				// Verify HTTP method
				if r.Method != "GET" {
					t.Errorf("Method = %v, want GET", r.Method)
				}

				// Verify auth header presence
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					t.Error("Missing Authorization header")
				}

				// Set content type header
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				if tt.statusCode >= 400 {
					if err := json.NewEncoder(w).Encode(map[string]interface{}{
						"errorMessages": []string{"Project not found"},
						"errors": map[string]string{
							"projectKey": "The project key is invalid",
						},
					}); err != nil {
						t.Errorf("Failed to encode error response: %v", err)
					}
					return
				}

				project := responsetypes.Project{
					ID:   "10000",
					Key:  tt.projectKey,
					Name: "Test Project",
				}
				err := json.NewEncoder(w).Encode(project)
				if err != nil {
					t.Errorf("Failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			client := &http.Client{}
			basicAuth := auth.NewBasicAuth("testuser", "secret123")
			service := NewService(client, server.URL, basicAuth)

			project, err := service.Get(context.Background(), tt.projectKey, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && project != nil {
				if project.Key != tt.projectKey {
					t.Errorf("Get() project key = %v, want %v", project.Key, tt.projectKey)
				}

				// Basic field validations for successful responses
				if project.ID == "" {
					t.Error("Get() project ID is empty")
				}
				if project.Name == "" {
					t.Error("Get() project name is empty")
				}
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name       string
		project    *responsetypes.Project
		wantErr    bool
		statusCode int
	}{
		{
			name: "success - basic project",
			project: &responsetypes.Project{
				Key:  "TEST",
				Name: "Test Project",
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "success - full project",
			project: &responsetypes.Project{
				Key:         "FULL",
				Name:        "Full Project",
				Description: "A test project with full details",
				URL:         "https://example.com/project",
				Lead: responsetypes.User{
					AccountID: "user123",
				},
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "error - missing key",
			project: &responsetypes.Project{
				Name: "Test Project",
			},
			wantErr:    true,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error - missing name",
			project: &responsetypes.Project{
				Key: "TEST",
			},
			wantErr:    true,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error - duplicate key",
			project: &responsetypes.Project{
				Key:  "TEST",
				Name: "Test Project",
			},
			wantErr:    true,
			statusCode: http.StatusConflict,
		},
		{
			name:       "error - nil project",
			project:    nil,
			wantErr:    true,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request URL and method
				if r.URL.Path != "/rest/api/3/project" {
					t.Errorf("Expected path '/rest/api/3/project', got %s", r.URL.Path)
				}
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				// Verify content type header
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
				}

				// Verify auth header presence
				if r.Header.Get("Authorization") == "" {
					t.Error("Missing Authorization header")
				}

				// Set response status and headers
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				if tt.statusCode >= 400 {
					if err := json.NewEncoder(w).Encode(map[string]interface{}{
						"errorMessages": []string{"Invalid project data"},
						"errors": map[string]string{
							"key":  "Project key is required",
							"name": "Project name is required",
						},
					}); err != nil {
						t.Errorf("Failed to encode error response: %v", err)
					}
					return
				}

				// For successful response, echo back the project with an ID
				response := *tt.project // Create a copy
				response.ID = "10000"   // Add generated ID
				err := json.NewEncoder(w).Encode(response)
				if err != nil {
					t.Errorf("Failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			client := &http.Client{}
			basicAuth := auth.NewBasicAuth("testuser", "secret123")
			service := NewService(client, server.URL, basicAuth)

			project, err := service.Create(context.Background(), tt.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && project != nil {
				if project.ID == "" {
					t.Error("Create() project ID is empty")
				}
				if tt.project.Key != "" && project.Key != tt.project.Key {
					t.Errorf("Create() project key = %v, want %v", project.Key, tt.project.Key)
				}
				if tt.project.Name != "" && project.Name != tt.project.Name {
					t.Errorf("Create() project name = %v, want %v", project.Name, tt.project.Name)
				}
				if tt.project.Description != "" && project.Description != tt.project.Description {
					t.Errorf("Create() project description = %v, want %v", project.Description, tt.project.Description)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name       string
		project    *responsetypes.Project
		opts       ProjectUpdateOpts
		wantURL    string
		wantBody   map[string]interface{}
		wantErr    bool
		statusCode int
	}{
		{
			name: "success - basic update",
			project: &responsetypes.Project{
				Key:         "TEST",
				Name:        "Updated Project",
				Description: "Updated description",
			},
			opts:    ProjectUpdateOpts{},
			wantURL: "/rest/api/3/project/TEST",
			wantBody: map[string]interface{}{
				"name":        "Updated Project",
				"key":         "TEST",
				"description": "Updated description",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success - with expand",
			project: &responsetypes.Project{
				Key:         "TEST",
				Name:        "Updated Project",
				Description: "Updated description",
			},
			opts: ProjectUpdateOpts{
				Expand: "description,lead",
			},
			wantURL: "/rest/api/3/project/TEST?expand=description%2Clead",
			wantBody: map[string]interface{}{
				"name":        "Updated Project",
				"key":         "TEST",
				"description": "Updated description",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success - full update",
			project: &responsetypes.Project{
				Key:         "TEST",
				Name:        "Updated Project",
				Description: "Updated description",
				URL:         "https://example.com",
				Lead: responsetypes.User{
					AccountID:   "123456",
					DisplayName: "User Name",
					Active:      true,
				},
				Style:      "classic",
				Simplified: false,
				AvatarURLs: responsetypes.AvatarUrls{
					Size_48x48: "https://example.com/avatar.png",
				},
			},
			opts:    ProjectUpdateOpts{},
			wantURL: "/rest/api/3/project/TEST",
			wantBody: map[string]interface{}{
				"name":        "Updated Project",
				"key":         "TEST",
				"description": "Updated description",
				"url":         "https://example.com",
				"style":       "classic",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "error - project not found",
			project: &responsetypes.Project{
				Key:  "NOTFOUND",
				Name: "Updated Project",
			},
			opts:       ProjectUpdateOpts{},
			wantURL:    "/rest/api/3/project/NOTFOUND",
			wantErr:    true,
			statusCode: http.StatusNotFound,
		},
		{
			name: "error - invalid key format",
			project: &responsetypes.Project{
				Key:  "test",
				Name: "Updated Project",
			},
			opts:       ProjectUpdateOpts{},
			wantURL:    "/rest/api/3/project/test",
			wantErr:    true,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error - unauthorized",
			project: &responsetypes.Project{
				Key:  "TEST",
				Name: "Updated Project",
			},
			opts:       ProjectUpdateOpts{},
			wantURL:    "/rest/api/3/project/TEST",
			wantErr:    true,
			statusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Read request body first
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read request body: %v", err)
				}
				r.Body.Close()

				// Verify URL
				gotURL := r.URL.Path
				if r.URL.RawQuery != "" {
					gotURL += "?" + r.URL.RawQuery
				}
				if gotURL != tt.wantURL {
					t.Errorf("URL = %v, want %v", gotURL, tt.wantURL)
				}

				// Verify method
				if r.Method != "PUT" {
					t.Errorf("Method = %v, want PUT", r.Method)
				}

				// Verify content type header
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
				}

				// Verify auth header presence
				if r.Header.Get("Authorization") == "" {
					t.Error("Missing Authorization header")
				}

				// Verify request body for non-error cases
				if tt.wantBody != nil {
					var gotBody map[string]interface{}
					if err := json.Unmarshal(bodyBytes, &gotBody); err != nil {
						t.Fatalf("Failed to decode request body: %v", err)
					}

					// Compare fields including nested objects
					for k, want := range tt.wantBody {
						got, exists := gotBody[k]
						if !exists {
							t.Errorf("Body missing key %q", k)
							continue
						}

						// Special handling for nested objects like Lead and avatarUrls
						if k == "lead" {
							wantLead, ok1 := want.(map[string]interface{})
							gotLead, ok2 := got.(map[string]interface{})
							if !ok1 || !ok2 {
								t.Errorf("Body[%q] type mismatch: got %T, want map[string]interface{}", k, got)
								continue
							}
							if !reflect.DeepEqual(gotLead, wantLead) {
								t.Errorf("Body[%q] = %v, want %v", k, gotLead, wantLead)
							}
							continue
						}

						// Special handling for avatarUrls
						if k == "avatarUrls" {
							// Convert both to map[string]interface{} for comparison
							wantUrls, ok1 := want.(map[string]string)
							gotUrls, ok2 := got.(map[string]interface{})
							if !ok1 {
								t.Errorf("Body[%q] type mismatch: want should be map[string]string", k)
								continue
							}
							if !ok2 {
								t.Errorf("Body[%q] type mismatch: got %T, want map[string]interface{}", k, got)
								continue
							}

							// Convert wantUrls to map[string]interface{} for comparison
							wantUrlsInterface := make(map[string]interface{})
							for k, v := range wantUrls {
								wantUrlsInterface[k] = v
							}

							if !reflect.DeepEqual(gotUrls, wantUrlsInterface) {
								t.Errorf("Body[%q] = %v, want %v", k, gotUrls, wantUrlsInterface)
							}
							continue
						}

						if !reflect.DeepEqual(got, want) {
							t.Errorf("Body[%q] = %v, want %v", k, got, want)
						}
					}
				}

				// Set response headers and status
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				if tt.statusCode >= 400 {
					if err := json.NewEncoder(w).Encode(map[string]interface{}{
						"errorMessages": []string{"Error updating project"},
						"errors": map[string]string{
							"projectKey": "Invalid project key format",
						},
					}); err != nil {
						t.Errorf("Failed to encode error response: %v", err)
					}
					return
				}

				// For successful response, return updated project
				updatedProject := *tt.project // Create a copy
				updatedProject.ID = "10000"   // Add generated ID
				if err := json.NewEncoder(w).Encode(updatedProject); err != nil {
					t.Errorf("Failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			client := &http.Client{}
			basicAuth := auth.NewBasicAuth("testuser", "secret123")
			service := NewService(client, server.URL, basicAuth)

			project, err := service.Update(context.Background(), tt.project, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && project != nil {
				if project.ID == "" {
					t.Error("Update() project ID is empty")
				}
				if tt.project.Key != "" && project.Key != tt.project.Key {
					t.Errorf("Update() project key = %v, want %v", project.Key, tt.project.Key)
				}
				if tt.project.Name != "" && project.Name != tt.project.Name {
					t.Errorf("Update() project name = %v, want %v", project.Name, tt.project.Name)
				}
				if tt.project.Description != "" && project.Description != tt.project.Description {
					t.Errorf("Update() project description = %v, want %v", project.Description, tt.project.Description)
				}
				if tt.project.Lead.AccountID != "" && !reflect.DeepEqual(project.Lead, tt.project.Lead) {
					t.Errorf("Update() project lead = %v, want %v", project.Lead, tt.project.Lead)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name       string
		projectKey string
		async      bool
		wantURL    string
		wantMethod string
		wantErr    bool
		statusCode int
	}{
		{
			name:       "success - sync delete",
			projectKey: "TEST",
			async:      false,
			wantURL:    "/rest/api/3/project/TEST",
			wantMethod: http.MethodDelete,
			statusCode: http.StatusNoContent,
		},
		{
			name:       "success - async delete",
			projectKey: "TEST",
			async:      true,
			wantURL:    "/rest/api/3/project/TEST/delete",
			wantMethod: http.MethodPost,
			statusCode: http.StatusAccepted,
		},
		{
			name:       "error - project not found",
			projectKey: "NOTFOUND",
			async:      false,
			wantURL:    "/rest/api/3/project/NOTFOUND",
			wantMethod: http.MethodDelete,
			wantErr:    true,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "error - unauthorized",
			projectKey: "TEST",
			async:      false,
			wantURL:    "/rest/api/3/project/TEST",
			wantMethod: http.MethodDelete,
			wantErr:    true,
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "error - server error",
			projectKey: "TEST",
			async:      false,
			wantURL:    "/rest/api/3/project/TEST",
			wantMethod: http.MethodDelete,
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request URL
				if r.URL.Path != tt.wantURL {
					t.Errorf("URL = %v, want %v", r.URL.Path, tt.wantURL)
				}

				// Verify HTTP method
				if r.Method != tt.wantMethod {
					t.Errorf("Method = %v, want %v", r.Method, tt.wantMethod)
				}

				// Verify auth header presence
				if r.Header.Get("Authorization") == "" {
					t.Error("Missing Authorization header")
				}

				w.WriteHeader(tt.statusCode)

				if tt.statusCode >= 400 {
					if err := json.NewEncoder(w).Encode(map[string]interface{}{
						"errorMessages": []string{"Error deleting project"},
						"errors": map[string]string{
							"projectKey": "Project not found",
						},
					}); err != nil {
						t.Errorf("Failed to encode error response: %v", err)
					}
				}
			}))
			defer server.Close()

			client := &http.Client{}
			basicAuth := auth.NewBasicAuth("testuser", "secret123")
			service := NewService(client, server.URL, basicAuth)

			err := service.Delete(context.Background(), tt.projectKey, tt.async)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
