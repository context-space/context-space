package apierrors

import (
	"fmt"
	"net/http"
)

// ErrorType is the type of an error
type ErrorType string

const (
	// ErrorTypeUnknown is used when the error type is unknown
	ErrorTypeUnknown ErrorType = "UNKNOWN"
	// ErrorTypeValidation is used for validation errors
	ErrorTypeValidation ErrorType = "VALIDATION"
	// ErrorTypeNotFound is used when a resource is not found
	ErrorTypeNotFound ErrorType = "NOT_FOUND"
	// ErrorTypeUnauthorized is used for authentication errors
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	// ErrorTypeForbidden is used for authorization errors
	ErrorTypeForbidden ErrorType = "FORBIDDEN"
	// ErrorTypeInternal is used for internal server errors
	ErrorTypeInternal ErrorType = "INTERNAL"
	// ErrorTypeProvider is used for provider-specific errors
	ErrorTypeProvider ErrorType = "PROVIDER"
	// ErrorTypeRateLimit is used for rate limiting errors
	ErrorTypeRateLimit ErrorType = "RATE_LIMIT"
	// ErrorTypeProviderRateLimit is used for provider-specific rate limiting errors
	ErrorTypeProviderRateLimit ErrorType = "PROVIDER_RATE_LIMIT"
	// ErrorTypeCredential is used for credential-related errors
	ErrorTypeCredential ErrorType = "CREDENTIAL"
	// ErrorTypeCircuitOpen is used when a circuit breaker is open
	ErrorTypeCircuitOpen ErrorType = "CIRCUIT_OPEN"
	// Error codes
	ErrInvalidRequest = "INVALID_REQUEST"
	ErrNotFound       = "NOT_FOUND"
	ErrConflict       = "CONFLICT"
	ErrInternal       = "INTERNAL_ERROR"
	ErrUnauthorized   = "UNAUTHORIZED"
)

// APIError represents a standardized API error
type APIError struct {
	Type        ErrorType `json:"type"`
	Code        string    `json:"code"`
	Message     string    `json:"message"`
	Details     any       `json:"details,omitempty"`
	ProviderErr *APIError `json:"provider_error,omitempty"`
	Err         error     `json:"-"`
	HTTPCode    int       `json:"-"`
}

// Error returns the error message
func (e *APIError) Error() string {
	if e.ProviderErr != nil {
		return fmt.Sprintf("%s: %s (provider error: %s)", e.Type, e.Message, e.ProviderErr.Error())
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the wrapped error
func (e *APIError) Unwrap() error {
	return e.Err
}

// WithDetails adds details to the error
func (e *APIError) WithDetails(details any) *APIError {
	e.Details = details
	return e
}

// WithProviderError adds a provider error
func (e *APIError) WithProviderError(providerErr *APIError) *APIError {
	e.ProviderErr = providerErr
	return e
}

// NewValidationError creates a new validation error
func NewValidationError(message string, err error) *APIError {
	return &APIError{
		Type:     ErrorTypeValidation,
		Code:     "VALIDATION_ERROR",
		Message:  message,
		Err:      err,
		HTTPCode: http.StatusBadRequest,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, err error) *APIError {
	return &APIError{
		Type:     ErrorTypeNotFound,
		Code:     "NOT_FOUND",
		Message:  message,
		Err:      err,
		HTTPCode: http.StatusNotFound,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, err error) *APIError {
	return &APIError{
		Type:     ErrorTypeUnauthorized,
		Code:     "UNAUTHORIZED",
		Message:  message,
		Err:      err,
		HTTPCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string, err error) *APIError {
	return &APIError{
		Type:     ErrorTypeForbidden,
		Code:     "FORBIDDEN",
		Message:  message,
		Err:      err,
		HTTPCode: http.StatusForbidden,
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(message string, err error) *APIError {
	return &APIError{
		Type:     ErrorTypeInternal,
		Code:     "INTERNAL_ERROR",
		Message:  message,
		Err:      err,
		HTTPCode: http.StatusInternalServerError,
	}
}

// NewProviderError creates a new provider error
func NewProviderError(message string, err error) *APIError {
	return &APIError{
		Type:     ErrorTypeProvider,
		Code:     "PROVIDER_ERROR",
		Message:  message,
		Err:      err,
		HTTPCode: http.StatusBadGateway,
	}
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(message string, err error) *APIError {
	return &APIError{
		Type:     ErrorTypeRateLimit,
		Code:     "RATE_LIMIT",
		Message:  message,
		Err:      err,
		HTTPCode: http.StatusTooManyRequests,
	}
}

// NewProviderRateLimitError creates a new provider rate limit error
func NewProviderRateLimitError(message string, err error) *APIError {
	return &APIError{
		Type:     ErrorTypeProviderRateLimit,
		Code:     "PROVIDER_RATE_LIMIT",
		Message:  message,
		Err:      err,
		HTTPCode: http.StatusTooManyRequests,
	}
}

// NewCredentialError creates a new credential error
func NewCredentialError(message string, err error) *APIError {
	return &APIError{
		Type:     ErrorTypeCredential,
		Code:     "CREDENTIAL_ERROR",
		Message:  message,
		Err:      err,
		HTTPCode: http.StatusBadRequest,
	}
}

// NewCircuitOpenError creates a new circuit open error
func NewCircuitOpenError(message string, err error) *APIError {
	return &APIError{
		Type:     ErrorTypeCircuitOpen,
		Code:     "CIRCUIT_OPEN",
		Message:  message,
		Err:      err,
		HTTPCode: http.StatusServiceUnavailable,
	}
}

// NewAPIError creates a new API error with the given error code, message, and details
func NewAPIError(code string, message string, details any) *APIError {
	var errorType ErrorType
	var httpCode int

	switch code {
	case ErrInvalidRequest:
		errorType = ErrorTypeValidation
		httpCode = http.StatusBadRequest
	case ErrNotFound:
		errorType = ErrorTypeNotFound
		httpCode = http.StatusNotFound
	case ErrConflict:
		errorType = ErrorTypeValidation
		httpCode = http.StatusConflict
	case ErrInternal:
		errorType = ErrorTypeInternal
		httpCode = http.StatusInternalServerError
	case ErrUnauthorized:
		errorType = ErrorTypeUnauthorized
		httpCode = http.StatusUnauthorized
	default:
		errorType = ErrorTypeUnknown
		httpCode = http.StatusInternalServerError
	}

	return &APIError{
		Type:     errorType,
		Code:     code,
		Message:  message,
		Details:  details,
		HTTPCode: httpCode,
	}
}
