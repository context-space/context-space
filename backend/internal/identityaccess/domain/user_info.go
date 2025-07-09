package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserInfo represents a user's info
type UserInfo struct {
	ID           string
	UserID       string
	InfoMetadata map[string]interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func NewUserInfo(userID string, infoMetadata map[string]interface{}) *UserInfo {
	return &UserInfo{
		ID:           uuid.New().String(),
		UserID:       userID,
		InfoMetadata: infoMetadata,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
