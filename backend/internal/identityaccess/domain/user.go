package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID          string
	SupID       string
	Email       string
	IsAnonymous bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func NewUser(supID, email string, isAnonymous bool) *User {
	return &User{
		ID:          uuid.New().String(),
		SupID:       supID,
		Email:       email,
		IsAnonymous: isAnonymous,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
