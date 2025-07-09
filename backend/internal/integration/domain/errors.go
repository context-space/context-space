package domain

import (
	"errors"
)

// Common error definitions
var (
	// ErrProviderNotFound is returned when a provider cannot be found
	ErrProviderNotFound = errors.New("provider not found")

	// ErrOperationNotFound is returned when an operation cannot be found
	ErrOperationNotFound = errors.New("operation not found")

	// ErrInvocationNotFound is returned when an invocation cannot be found
	ErrInvocationNotFound = errors.New("invocation not found")

	// ErrValidation is returned when validation fails
	ErrValidation = errors.New("validation error")

	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrProviderUnavailable is returned when a provider is unavailable
	ErrProviderUnavailable = errors.New("provider unavailable")
)
