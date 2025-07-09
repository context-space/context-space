package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// CredentialType indicates the type of credential
type CredentialType string

const (
	// CredentialTypeOAuth represents OAuth credentials
	CredentialTypeOAuth CredentialType = "oauth"
	// CredentialTypeAPIKey represents API key credentials
	CredentialTypeAPIKey CredentialType = "apikey"
	// CredentialTypeBasicAuth represents basic authentication credentials
	CredentialTypeBasicAuth CredentialType = "basicauth"
	// CredentialTypeNone represents no credentials (public API)
	CredentialTypeNone CredentialType = "none"
)

// Credential defines the common attributes for all credentials
type Credential struct {
	ID                 string
	UserID             string
	ProviderIdentifier string
	Type               CredentialType
	IsValid            bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
	LastUsedAt         time.Time
	DeletedAt          *time.Time
}

// NewCredential creates a new credential with default values
func NewCredential(userID, providerIdentifier string, credType CredentialType) (*Credential, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if providerIdentifier == "" {
		return nil, errors.New("provider identifier is required")
	}

	return &Credential{
		ID:                 uuid.New().String(),
		UserID:             userID,
		ProviderIdentifier: providerIdentifier,
		Type:               credType,
		IsValid:            true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}, nil
}

// OAuthCredential represents OAuth credentials
type OAuthCredential struct {
	*Credential
	// Encryption metadata
	EncryptionMetadata *EncryptionMetadata
	// In-memory only
	Token *oauth2.Token
	// Scopes granted to this credential
	Scopes []string
}

// NewOAuthCredential creates a new OAuth credential
func NewOAuthCredential(userID, providerIdentifier string, oauthToken *oauth2.Token, scopes []string) (*OAuthCredential, error) {
	cred, err := NewCredential(userID, providerIdentifier, CredentialTypeOAuth)
	if err != nil {
		return nil, err
	}

	return &OAuthCredential{
		Credential: cred,
		Token:      oauthToken,
		Scopes:     scopes,
	}, nil
}

// APIKeyCredential represents API key credentials
type APIKeyCredential struct {
	*Credential
	// Encryption metadata
	EncryptionMetadata *EncryptionMetadata
	// In-memory only
	APIKey string
}

// NewAPIKeyCredential creates a new API key credential
func NewAPIKeyCredential(userID, providerIdentifier, apiKey string) (*APIKeyCredential, error) {
	if apiKey == "" {
		return nil, errors.New("API key is required")
	}

	cred, err := NewCredential(userID, providerIdentifier, CredentialTypeAPIKey)
	if err != nil {
		return nil, err
	}

	return &APIKeyCredential{
		Credential: cred,
		APIKey:     apiKey,
	}, nil
}

// BasicAuthCredential represents basic authentication credentials
type BasicAuthCredential struct {
	*Credential
	// Encryption metadata
	EncryptionMetadata *EncryptionMetadata
	// In-memory only
	Username string
	Password string
}

// NewBasicAuthCredential creates a new basic authentication credential
func NewBasicAuthCredential(userID, providerIdentifier, username, password string) (*BasicAuthCredential, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	cred, err := NewCredential(userID, providerIdentifier, CredentialTypeBasicAuth)
	if err != nil {
		return nil, err
	}

	return &BasicAuthCredential{
		Credential: cred,
		Username:   username,
		Password:   password,
	}, nil
}

// NoneCredential represents no credentials for public APIs
type NoneCredential struct {
	*Credential
}

// NewNoneCredential creates a credential for providers that don't require authentication
func NewNoneCredential(userID, providerIdentifier string) (*NoneCredential, error) {
	cred, err := NewCredential(userID, providerIdentifier, CredentialTypeNone)
	if err != nil {
		return nil, err
	}

	return &NoneCredential{
		Credential: cred,
	}, nil
}
