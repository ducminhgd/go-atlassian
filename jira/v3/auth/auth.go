package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

// BasicAuth represents basic authentication using username and password/token
type BasicAuth struct {
	Username string
	Password string // can be either password or API token
}

// NewBasicAuth creates a new BasicAuth instance
func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{
		Username: username,
		Password: password,
	}
}

// AddAuthentication adds the basic authentication headers to the request
func (a *BasicAuth) AddAuthentication(req *http.Request) error {
	if a.Username == "" || a.Password == "" {
		return fmt.Errorf("username and password cannot be empty")
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.Username, a.Password)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	return nil
}

// Authenticator interface defines the methods required for authentication
type Authenticator interface {
	AddAuthentication(req *http.Request) error
}
