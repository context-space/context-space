package acl

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/providercore/domain"
	contractTranslation "github.com/context-space/context-space/backend/internal/shared/contract/providertranslation"
	"golang.org/x/text/language"
)

type ProviderTranslationACL struct {
	providerTranslationReader contractTranslation.ProviderTranslationReader
	obs                       *observability.ObservabilityProvider
}

func NewProviderTranslationACL(
	providerTranslationReader contractTranslation.ProviderTranslationReader,
	obs *observability.ObservabilityProvider,
) domain.ProviderTranslationACL {
	return &ProviderTranslationACL{
		providerTranslationReader: providerTranslationReader,
		obs:                       obs,
	}
}

func (acl *ProviderTranslationACL) GetProviderTranslation(ctx context.Context, providerIdentifier string, preferredLang language.Tag) (*contractTranslation.ProviderTranslationDTO, error) {
	ctx, span := acl.obs.Tracer.Start(ctx, "ProviderTranslationACL.GetTranslation")
	defer span.End()

	return acl.providerTranslationReader.GetTranslation(ctx, providerIdentifier, preferredLang)
}
