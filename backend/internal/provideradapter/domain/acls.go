package domain

import (
	"context"

	contractProvider "github.com/context-space/context-space/backend/internal/shared/contract/providercore"
	contractTranslation "github.com/context-space/context-space/backend/internal/shared/contract/providertranslation"
	"golang.org/x/text/language"
)

// ProvidercoreAcl defines the interface for accessing provider core data
// This follows DDD anti-corruption layer pattern
type ProvidercoreAcl interface {
	// Provider metadata access with translation support
	GetProvidercoreData(ctx context.Context, identifier string, preferredLang language.Tag) (*contractProvider.ProviderDTO, error)

	// List providers with translation support
	ListProviders(ctx context.Context) ([]*contractProvider.ProviderDTO, error)

	// Provider metadata access without translation support
	GetProvidercoreDataWithoutTranslation(ctx context.Context, identifier string) (*contractProvider.ProviderDTO, error)
}

type ProviderTranslationAcl interface {
	GetProviderTranslation(ctx context.Context, providerIdentifier string, preferredLang language.Tag) (*contractTranslation.ProviderTranslationDTO, error)
}
