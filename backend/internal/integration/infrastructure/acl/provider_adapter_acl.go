package acl

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/integration/domain"
	contractAdapter "github.com/context-space/context-space/backend/internal/shared/contract/provideradapter"
	"go.uber.org/zap"
)

// ProviderAdapterACL implements the Anti-Corruption Layer pattern for provider adapter operations
// This isolates the integration domain from external provideradapter module dependencies
type ProviderAdapterACL struct {
	providerAdapterContract contractAdapter.ProviderAdapterContract
	obs                     *observability.ObservabilityProvider
}

// NewProviderAdapterACL creates a new provider adapter ACL
func NewProviderAdapterACL(
	providerAdapterContract contractAdapter.ProviderAdapterContract,
	obs *observability.ObservabilityProvider,
) domain.AdapterProvider {
	return &ProviderAdapterACL{
		providerAdapterContract: providerAdapterContract,
		obs:                     obs,
	}
}

// GetAdapter gets an adapter through the contract layer - simplified without wrapper
func (acl *ProviderAdapterACL) GetAdapterByProviderIdentifier(ctx context.Context, providerIdentifier string) (contractAdapter.AdapterContract, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "ProviderAdapterACL.GetAdapterByProviderIdentifier")
	defer span.End()

	acl.obs.Logger.Debug(ctx, "Getting provider adapter through ACL",
		zap.String("provider_identifier", providerIdentifier))

	// Direct call through contract layer - no unnecessary wrapper
	adapter, err := acl.providerAdapterContract.GetAdapterContract(ctx, providerIdentifier)
	if err != nil {
		acl.obs.Logger.Error(ctx, "Failed to get adapter through contract",
			zap.String("provider_identifier", providerIdentifier),
			zap.Error(err))
		return nil, acl.translateError(ctx, "get_adapter", err)
	}

	acl.obs.Logger.Debug(ctx, "Successfully got adapter through ACL",
		zap.String("provider_identifier", providerIdentifier))

	return adapter, nil
}

// translateError translates external errors into domain-appropriate errors
// This is a key function of the ACL - protecting the domain from external error formats
func (acl *ProviderAdapterACL) translateError(ctx context.Context, operation string, err error) error {
	// Error translation can be performed here as needed
	// For example: convert specific errors from external modules to domain errors
	// Or add additional context information

	acl.obs.Logger.Warn(ctx, "Translating external error through ACL",
		zap.String("operation", operation),
		zap.Error(err))

	return err
}
