package acl

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	providerService "github.com/context-space/context-space/backend/internal/providercore/application"
	providerCore "github.com/context-space/context-space/backend/internal/providercore/domain"
)

type ProviderCoreACL struct {
	providerService *providerService.ProviderService
	obs             *observability.ObservabilityProvider
}

func NewProviderCoreACL(
	providerService *providerService.ProviderService,
	obs *observability.ObservabilityProvider,
) *ProviderCoreACL {
	return &ProviderCoreACL{
		providerService: providerService,
		obs:             obs,
	}
}

func (acl *ProviderCoreACL) GetProviderCoreData(ctx context.Context, identifier string) (*providerCore.Provider, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "ProviderCoreACL.GetProviderCoreData")
	defer span.End()

	return acl.providerService.GetProviderByIdentifier(ctx, identifier)
}

func (acl *ProviderCoreACL) ListProviders(ctx context.Context) ([]*providerCore.Provider, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "ProviderCoreACL.ListProviders")
	defer span.End()

	return acl.providerService.ListProviders(ctx)
}
