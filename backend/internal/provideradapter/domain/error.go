package domain

import (
	"fmt"
)

// Constants for standard adapter error codes (now string type).
const (
	// Credential related errors
	ErrCredentialError string = "CREDENTIAL_ERROR" // Error related to credentials (invalid, missing, expired)

	// Provider interaction errors
	ErrProviderAPIError string = "PROVIDER_API_ERROR" // General error during interaction with the provider API
	ErrLLMProviderError string = "LLM_PROVIDER_ERROR" // Specific error from the underlying LLM provider (e.g., OpenAI)
	ErrLLMEmptyResponse string = "LLM_EMPTY_RESPONSE" // LLM provider returned an empty or unusable response

	// Input/Output errors
	ErrInvalidParameters string = "INVALID_PARAMETERS" // Input parameters failed validation
	ErrDecodingError     string = "DECODING_ERROR"     // Error decoding provider response
	ErrEncodingError     string = "ENCODING_ERROR"     // Error encoding request to provider

	// Configuration errors
	ErrMissingRequiredField string = "MISSING_REQUIRED_FIELD" // A required configuration field is missing

	// General errors
	ErrOperationNotSupported string = "OPERATION_NOT_SUPPORTED" // The requested operation is not supported by the adapter
	ErrTimeout               string = "TIMEOUT"                 // Operation timed out
	ErrUnknownError          string = "UNKNOWN_ERROR"           // An unexpected or unclassified error occurred
	ErrInternal              string = "INTERNAL_ERROR"          // Internal error within the adapter logic
)

// AdapterError represents an error returned by an adapter
type AdapterError struct {
	ProviderIdentifier  string
	OperationIdentifier string
	ErrorCode           string
	ErrorMessage        string
	StatusCode          int
	Raw                 interface{}
}

// Error returns the error message
func (e *AdapterError) Error() string {
	return fmt.Sprintf("[%s] %s: %s (code: %s, status: %d)",
		e.ProviderIdentifier, e.OperationIdentifier, e.ErrorMessage, e.ErrorCode, e.StatusCode)
}

// NewAdapterError creates a new adapter error
func NewAdapterError(
	providerIdentifier string,
	operationIdentifier string,
	errorCode string,
	errorMessage string,
	statusCode int,
) *AdapterError {
	return &AdapterError{
		ProviderIdentifier:  providerIdentifier,
		OperationIdentifier: operationIdentifier,
		ErrorCode:           errorCode,
		ErrorMessage:        errorMessage,
		StatusCode:          statusCode,
	}
}
