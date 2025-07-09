package application

import (
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
)

// AdapterTemplate provides a base implementation for a specific type of adapter
type AdapterTemplate interface {
	// CreateAdapter creates a new adapter instance from a configuration
	CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error)

	// ValidateConfig validates a configuration for this template
	ValidateConfig(provider *domain.ProviderAdapterConfig) error
}
