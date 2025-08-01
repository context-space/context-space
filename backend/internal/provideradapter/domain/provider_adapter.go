package domain

import (
	"context"
	"time"

	providercore "github.com/context-space/context-space/backend/internal/providercore/domain"
)

// AdapterConfig defines common configuration for adapters
type AdapterConfig struct {
	Timeout        time.Duration
	MaxRetries     int
	RetryBackoff   time.Duration
	RateLimit      int
	RatePeriod     time.Duration
	CircuitBreaker CircuitBreakerConfig
}

// CircuitBreakerConfig holds configuration for a circuit breaker
type CircuitBreakerConfig struct {
	// FailureThreshold is the number of failures before opening the circuit
	FailureThreshold int
	// ResetTimeout is the time to wait before trying to close the circuit in seconds
	ResetTimeout int
	// HalfOpenMaxCalls is the number of calls allowed in half-open state
	HalfOpenMaxCalls int
}

// ProviderAdapterInfo contains the minimal provider information needed by adapters
type ProviderAdapterInfo struct {
	Identifier  string
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	AuthType    providercore.ProviderAuthType `json:"auth_type"`
	Permissions []providercore.Permission     `json:"permissions"`
	Operations  []providercore.Operation      `json:"operations"`
	Status      providercore.ProviderStatus   `json:"status"`
}

type ProviderAdapterConfig struct {
	ProviderAdapterInfo
	ID           string                 `json:"id"`
	OAuthConfig  *OAuthConfig           `json:"oauth_config"`
	CustomConfig map[string]interface{} `json:"custom_config"`
}

// Adapter is the interface for all provider adapters
type Adapter interface {
	// Execute an operation call to the provider
	Execute(
		ctx context.Context,
		operationID string,
		params map[string]interface{},
		credential interface{},
	) (interface{}, error)

	// GetProviderAdapterInfo returns information about this provider
	GetProviderAdapterInfo() *ProviderAdapterInfo
}

// ProviderAdapterLoader define the interface for provider adapter loader
type ProviderAdapterLoader interface {
	// LoadProvider load a single provider
	LoadProvider(config *ProviderAdapterConfig) error
	// GetLoadedProviders get all loaded providers
	GetLoadedProviders() []ProviderAdapterInfo
	// UnloadProvider unload a single provider
	UnloadProvider(identifier string) error
}
