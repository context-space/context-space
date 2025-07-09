package domain

import (
	"context"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Get retrieves a user by ID
	Get(ctx context.Context, id string) (*User, error)

	// GetBySupID retrieves a user by sup id
	GetBySupID(ctx context.Context, supID string) (*User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// List retrieves users with pagination
	List(ctx context.Context, offset, limit int) ([]*User, error)

	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// Delete soft-deletes a user
	Delete(ctx context.Context, id string) error

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
}

// UserInfoRepository defines the interface for user identity data access
type UserInfoRepository interface {
	// Get gets a user info by id
	Get(ctx context.Context, id string) (*UserInfo, error)

	// GetByUserID gets a user info by user id
	GetByUserID(ctx context.Context, userID string) (*UserInfo, error)

	// Create creates a user info
	Create(ctx context.Context, info *UserInfo) error

	// Update updates a user info
	Update(ctx context.Context, info *UserInfo) error

	// Delete deletes a user info by id
	Delete(ctx context.Context, id string) error
}

// APIKeyRepository defines the interface for API key data access
type APIKeyRepository interface {
	// Get retrieves an API key by ID
	Get(ctx context.Context, id string) (*APIKey, error)

	// GetByKeyValue retrieves an API key by its value
	GetByKeyValue(ctx context.Context, value string) (*APIKey, error)

	// ListByUserID retrieves API keys for a user with pagination
	ListByUserID(ctx context.Context, userID string) ([]*APIKey, error)

	// Create creates a new API key
	Create(ctx context.Context, apiKey *APIKey) error

	// Update updates an existing API key
	Update(ctx context.Context, apiKey *APIKey) error

	// Delete soft-deletes an API key
	Delete(ctx context.Context, id string) error
}
