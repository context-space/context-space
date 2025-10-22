package acl

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/integration/domain"
	contractProvider "github.com/context-space/context-space/backend/internal/shared/contract/providercore"
	"go.uber.org/zap"
)

// ProviderACL implements the Anti-Corruption Layer pattern for provider operations
// This isolates the integration domain from external providercore module dependencies
type ProviderACL struct {
	providerContract contractProvider.ProviderCoreReader
	obs              *observability.ObservabilityProvider
}

// NewProviderACL creates a new provider ACL
func NewProviderACL(providerContract contractProvider.ProviderCoreReader, obs *observability.ObservabilityProvider) domain.ProviderProvider {
	return &ProviderACL{
		providerContract: providerContract,
		obs:              obs,
	}
}

// GetProviderByIdentifier adapts the contract interface method signature
func (acl *ProviderACL) GetProviderByIdentifier(ctx context.Context, providerIdentifier string) (*contractProvider.ProviderDTO, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "ProviderACL.GetProviderByIdentifier")
	defer span.End()

	acl.obs.Logger.Debug(ctx, "Getting provider by identifier through ACL",
		zap.String("provider_identifier", providerIdentifier))

	// Call contract layer with default language (English)
	provider, err := acl.providerContract.GetProviderByIdentifier(ctx, providerIdentifier)
	if err != nil {
		acl.obs.Logger.Error(ctx, "Failed to get provider through contract",
			zap.String("provider_identifier", providerIdentifier),
			zap.Error(err))
		return nil, acl.translateError(ctx, "get_provider", err)
	}

	acl.obs.Logger.Debug(ctx, "Successfully got provider through ACL",
		zap.String("provider_identifier", providerIdentifier))

	return provider, nil
}

// translateError translates external errors into domain-appropriate errors
// This is a key function of the ACL - protecting the domain from external error formats
func (acl *ProviderACL) translateError(ctx context.Context, operation string, err error) error {
	// Error translation can be performed here as needed
	// For example: convert specific errors from external modules to domain errors
	// Or add additional context information

	acl.obs.Logger.Warn(ctx, "Translating external error through ACL",
		zap.String("operation", operation),
		zap.Error(err))

	return err
}
