package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
)

func TestBasicAuth_AddAuthentication(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid credentials",
			username: "testuser",
			password: "valid-password",
			wantErr:  false,
		},
		{
			name:     "With API token",
			username: "user@example.com",
			password: "api-token-123",
			wantErr:  false,
		},
		{
			name:     "Empty username",
			username: "",
			password: "valid-password",
			wantErr:  true,
		},
		{
			name:     "Empty password",
			username: "testuser",
			password: "",
			wantErr:  true,
		},
		{
			name:     "Both empty",
			username: "",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewBasicAuth(tt.username, tt.password)
			req, _ := http.NewRequest("GET", "http://example.com", nil)

			err := a.AddAuthentication(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BasicAuth.AddAuthentication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got := req.Header.Get("Authorization")
				if got == "" {
					t.Error("BasicAuth.AddAuthentication() did not set Authorization header")
				}

				// Verify the auth header format
				expectedAuth := fmt.Sprintf("Basic %s",
					base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", tt.username, tt.password))))
				if got != expectedAuth {
					t.Errorf("BasicAuth.AddAuthentication() = %v, want %v", got, expectedAuth)
				}
			}
		})
	}
}
