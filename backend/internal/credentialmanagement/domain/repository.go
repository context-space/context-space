package domain

import (
	"context"
	"time"

	contractCredential "github.com/context-space/context-space/backend/internal/shared/contract/credentialmanagement"
	"golang.org/x/oauth2"
)

// CredentialRepository defines the interface for credential data access
type CredentialRepository interface {
	// GetByID retrieves a credential by ID
	GetByID(ctx context.Context, id string) (*Credential, error)

	// GetByUserAndProvider retrieves a credential by user ID and provider ID
	GetByUserAndProvider(ctx context.Context, userID, providerIdentifier string) (*Credential, error)

	// ListByUser retrieves all credentials for a user
	ListByUser(ctx context.Context, userID string) ([]*Credential, error)

	// ListByID retrieves credentials by IDs
	ListByID(ctx context.Context, ids []string) ([]*Credential, error)

	// Create creates a new credential
	Create(ctx context.Context, credential *Credential) error

	// Delete soft-deletes a credential
	Delete(ctx context.Context, id string) error

	// UpdateLastUsedAt updates the last used at time of a credential
	UpdateLastUsedAt(ctx context.Context, id string) error
}

// OAuthCredentialRepository defines the interface for OAuth credential data access
type OAuthCredentialRepository interface {
	// GetByCredentialID retrieves an OAuth credential by credential ID
	GetByCredentialID(ctx context.Context, credentialID string) (*OAuthCredential, error)

	// Create creates a new OAuth credential
	Create(ctx context.Context, credential *OAuthCredential) error

	// ListByExpiryWithin lists OAuth credentials by expiry within the given time
	ListByExpiryWithin(ctx context.Context, expiry time.Time) ([]*OAuthCredential, error)

	// Update updates an OAuth credential
	Update(ctx context.Context, credential *OAuthCredential) error
}

// APIKeyCredentialRepository defines the interface for API key credential data access
type APIKeyCredentialRepository interface {
	// GetByCredentialID retrieves an API key credential by credential ID
	GetByCredentialID(ctx context.Context, credentialID string) (*APIKeyCredential, error)

	// Create creates a new API key credential
	Create(ctx context.Context, credential *APIKeyCredential) error
}

// CredentialFactory can create and retrieve specialized credentials
type CredentialFactory interface {
	// CreateOAuth creates a new OAuth credential
	CreateOAuth(ctx context.Context, userID, providerIdentifier string, oauthToken *oauth2.Token, scopes []string) (*OAuthCredential, error)

	// CreateAPIKey creates a new API key credential
	CreateAPIKey(ctx context.Context, userID, providerIdentifier, apiKey string) (*APIKeyCredential, error)

	// CreateNone creates a new no-auth credential
	CreateNone(ctx context.Context, userID, providerIdentifier string) (*contractCredential.CredentialDTO, error)

	// GetCredential retrieves a credential by ID and converts it to the proper type
	GetCredential(ctx context.Context, id string) (interface{}, error)

	// GetCredentialByUserAndProvider retrieves a credential by user ID and provider identifier and converts it to the proper type
	GetCredentialByUserAndProvider(ctx context.Context, userID, providerIdentifier string) (interface{}, error)

	// UpdateCredentialLastUsedAt updates the last used at time of a credential
	UpdateCredentialLastUsedAt(ctx context.Context, credential interface{}) error
}

// VaultAlgorithm represents encryption algorithms used in Vault
type VaultAlgorithm string

const (
	// Algorithm constants
	AlgorithmAESGCM VaultAlgorithm = "aes-gcm" // Default algorithm
	AlgorithmSM4GCM VaultAlgorithm = "sm4-gcm" // China-specific algorithm
)

// VaultRegion represents a geographical or regulatory region.
type VaultRegion string

const (
	RegionEU VaultRegion = "eu"
	RegionUS VaultRegion = "us"
	RegionCN VaultRegion = "cn"
)

// EncryptionMetadata contains information about the encryption operation.
type EncryptionMetadata struct {
	Region         VaultRegion    `json:"region"`
	KeyVersion     int            `json:"key_version"`
	CredentialType CredentialType `json:"credential_type"`
	Algorithm      VaultAlgorithm `json:"algorithm"`
	Ciphertext     string         `json:"ciphertext"` // The actual base64 Vault ciphertext
}

// VaultService defines the interface for secure credential storage
type VaultService interface {
	// EncryptData encrypts plaintext data for a specific region.
	// It returns structured metadata including the ciphertext.
	EncryptData(ctx context.Context, plaintext string, region VaultRegion, credentialType CredentialType) (*EncryptionMetadata, error)

	// DecryptData decrypts data using metadata embedded within the provided structure.
	// The implementation routes the request to the correct regional Vault based on metadata.Region.
	DecryptData(ctx context.Context, metadata *EncryptionMetadata) (string, error)

	// EncryptJSON encrypts a JSON-serializable structure for a specific region.
	EncryptJSON(ctx context.Context, data interface{}, region VaultRegion, credentialType CredentialType) (*EncryptionMetadata, error)

	// DecryptJSON decrypts data into the provided target structure.
	DecryptJSON(ctx context.Context, metadata *EncryptionMetadata, target interface{}) error
}

type TokenRefresh interface {
	// RefreshAccessTokenIfNeeded refreshes the access token if needed
	RefreshAccessTokenIfNeeded(ctx context.Context, providerIdentifier string, credential interface{}) (interface{}, error)

	// RefreshAccessTokens refreshes all access tokens
	RefreshAccessTokens(ctx context.Context) error
}
