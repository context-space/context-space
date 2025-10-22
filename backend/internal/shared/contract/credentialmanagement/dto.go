package credentialmanagement

import "time"

type CredentialDTO struct {
	ID                 string
	UserID             string
	ProviderIdentifier string
	Type               string
	IsValid            bool
	LastUsedAt         time.Time
}

// NoneCredentialDTO represents a no-authentication credential for contract communication
type NoneCredentialDTO struct {
	UserID             string `json:"user_id"`
	ProviderIdentifier string `json:"provider_identifier"`
	CreatedAt          int64  `json:"created_at"`
	UpdatedAt          int64  `json:"updated_at"`
	LastUsedAt         *int64 `json:"last_used_at,omitempty"`
}

// CredentialResponse represents a generic credential response
type CredentialResponse struct {
	Type       string      `json:"type"`       // "none", "api_key", "oauth", "basic"
	Credential interface{} `json:"credential"` // Actual credential data
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
}
