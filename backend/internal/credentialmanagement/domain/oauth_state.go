package domain

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"

	"github.com/google/uuid"
)

type OAuthStateStatus string

const (
	OAuthStateStatusPending OAuthStateStatus = "pending"
	OAuthStateStatusSuccess OAuthStateStatus = "success"
	OAuthStateStatusFailed  OAuthStateStatus = "failed"
)

type OAuthStateData struct {
	ID                 string                 `json:"id"`
	State              string                 `json:"state"`
	CodeVerifier       string                 `json:"code_verifier"`
	CodeChallenge      string                 `json:"code_challenge"`
	Status             OAuthStateStatus       `json:"status"`
	UserID             string                 `json:"user_id"`
	ProviderIdentifier string                 `json:"provider_identifier"`
	RedirectURL        string                 `json:"redirect_url"`
	Permissions        []string               `json:"permissions"`
	UserData           map[string]interface{} `json:"user_data"`
	CallbackParams     map[string]interface{} `json:"callback_params"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	DeletedAt          *time.Time             `json:"deleted_at"`
}

func NewOAuthStateData(
	userID, providerIdentifier, redirectURL string,
	permissions []string,
	userData map[string]interface{},
) (*OAuthStateData, error) {
	state, err := generateRandomOAuthState()
	if err != nil {
		return nil, err
	}

	codeVerifier, err := generateRandomOAuthCodeVerifier()
	if err != nil {
		return nil, err
	}

	codeChallenge, err := generateRandomOAuthCodeChallenge(codeVerifier)
	if err != nil {
		return nil, err
	}

	return &OAuthStateData{
		ID:                 uuid.New().String(),
		State:              state,
		CodeVerifier:       codeVerifier,
		CodeChallenge:      codeChallenge,
		Status:             OAuthStateStatusPending,
		UserID:             userID,
		ProviderIdentifier: providerIdentifier,
		Permissions:        permissions,
		RedirectURL:        redirectURL,
		UserData:           userData,
		CallbackParams:     map[string]interface{}{},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}, nil
}

func (o *OAuthStateData) SetCallbackParams(callbackParams map[string]interface{}) {
	o.CallbackParams = callbackParams
	o.UpdatedAt = time.Now()
}

func (o *OAuthStateData) SetStatus(status OAuthStateStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
}

// generateRandomOAuthState generates a cryptographically secure random state string
func generateRandomOAuthState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// generateRandomOAuthCodeChallenge generates a cryptographically secure random code challenge string
func generateRandomOAuthCodeChallenge(codeVerifier string) (string, error) {
	var err error
	if codeVerifier == "" {
		codeVerifier, err = generateRandomOAuthCodeVerifier()
		if err != nil {
			return "", err
		}
	}
	hash := sha256.Sum256([]byte(codeVerifier))
	// use standard url safe base64 encoding, remove padding
	return strings.TrimRight(base64.URLEncoding.EncodeToString(hash[:]), "="), nil
}

// generateRandomOAuthCodeVerifier generates a cryptographically secure random code verifier string
func generateRandomOAuthCodeVerifier() (string, error) {
	b := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// use standard url safe base64 encoding, remove padding
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}

// OAuthStateRepository defines the interface for OAuth state operations
type OAuthStateRepository interface {
	// StoreStateData stores a state with associated data
	StoreStateData(ctx context.Context, data *OAuthStateData, expiration time.Duration) error

	// GetStateDataByState gets the data associated with a state
	GetStateDataByState(ctx context.Context, state string) (*OAuthStateData, error)

	// GetStateDataByID gets the data associated with a state by ID
	GetStateDataByID(ctx context.Context, id string) (*OAuthStateData, error)

	// UpdateStateData updates the data of a state
	UpdateStateData(ctx context.Context, data *OAuthStateData) error
}

// OAuthStateService handles OAuth state operations
type OAuthStateService interface {
	// StoreStateData stores it with associated data
	StoreStateData(ctx context.Context, data *OAuthStateData) error

	// GetByState gets by a state and returns associated data
	GetByState(ctx context.Context, state string) (*OAuthStateData, error)

	// GetByID gets by an id and returns associated data
	GetByID(ctx context.Context, id string) (*OAuthStateData, error)

	// UpdateStateData updates the data of a state
	UpdateStateData(ctx context.Context, data *OAuthStateData) error
}
