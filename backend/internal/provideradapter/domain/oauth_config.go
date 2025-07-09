package domain

import "fmt"

// OAuthConfig defines the configuration information required for OAuth authentication
type OAuthConfig struct {
	// ClientID is the client ID of the OAuth application
	ClientID string `json:"client_id"`

	// ClientSecret is the client secret of the OAuth application
	ClientSecret string `json:"client_secret"`
}

// Validate checks if OAuthConfig contains all required fields
func (c *OAuthConfig) Validate() error {
	if c.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("client_secret is required")
	}
	return nil
}

// Clone returns a deep copy of OAuthConfig
func (c *OAuthConfig) Clone() *OAuthConfig {
	if c == nil {
		return nil
	}

	clone := &OAuthConfig{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
	}

	return clone
}
