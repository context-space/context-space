package domain

import "context"

type ProviderTranslationRepository interface {
	GetProviderTranslation(ctx context.Context, providerIdentifier string, languageCode string) (*ProviderTranslation, error)
}
