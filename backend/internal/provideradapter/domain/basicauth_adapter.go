package domain

import (
	"context"
)

// BasicAuthAdapter extends the base Adapter with basic auth-specific functionality
type BasicAuthAdapter interface {
	Adapter

	// ValidateBasicAuth validates username/password with the provider
	ValidateBasicAuth(ctx context.Context, username, password string) error
}
