package acl

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	contractProvider "github.com/context-space/context-space/backend/internal/shared/contract/providercore"
	"golang.org/x/text/language"
)

type ProvidercoreACL struct {
	providerContract contractProvider.ProviderCoreReader
	obs              *observability.ObservabilityProvider
}

func NewProviderCoreACL(
	providerContract contractProvider.ProviderCoreReader,
	obs *observability.ObservabilityProvider,
) domain.ProvidercoreAcl {
	return &ProvidercoreACL{
		providerContract: providerContract,
		obs:              obs,
	}
}

func (acl *ProvidercoreACL) GetProvidercoreData(ctx context.Context, identifier string, preferredLang language.Tag) (*contractProvider.ProviderDTO, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "ProviderCoreACL.GetProviderCoreData")
	defer span.End()

	return acl.providerContract.GetProviderWithTranslation(ctx, identifier, preferredLang)
}

func (acl *ProvidercoreACL) ListProviders(ctx context.Context) ([]*contractProvider.ProviderDTO, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "ProviderCoreACL.ListProviders")
	defer span.End()

	return acl.providerContract.ListProviders(ctx)
}

func (acl *ProvidercoreACL) GetProvidercoreDataWithoutTranslation(ctx context.Context, identifier string) (*contractProvider.ProviderDTO, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "ProviderCoreACL.GetProviderCoreDataWithoutTranslation")
	defer span.End()

	return acl.providerContract.GetProviderByIdentifier(ctx, identifier)
}
