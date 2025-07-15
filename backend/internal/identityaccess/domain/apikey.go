package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/context-space/context-space/backend/internal/shared/utils"
	"github.com/google/uuid"
)

// APIKey represents an API key for authentication
type APIKey struct {
	ID          string
	UserID      string
	KeyValue    string
	Name        string
	Description string
	LastUsed    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// NewAPIKey creates a new API key with default values
func NewAPIKey(userID, name, description string) *APIKey {
	return &APIKey{
		ID:          uuid.New().String(),
		UserID:      userID,
		KeyValue:    generateAPIKeyValue(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateLastUsed updates the last used timestamp
func (k *APIKey) UpdateLastUsed() {
	now := time.Now()
	k.LastUsed = &now
	k.UpdatedAt = now
}

// generateAPIKeyValue generates a random API key value
func generateAPIKeyValue() string {
	// Generate 32 random bytes (256 bits)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// In case of error, fallback to a UUID-based key
		return uuid.New().String()
	}

	// Encode to hex string
	return utils.StringsBuilder("cs-", hex.EncodeToString(bytes))
}
