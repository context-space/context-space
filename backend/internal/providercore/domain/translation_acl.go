package domain

import (
	"context"

	contractTranslation "github.com/context-space/context-space/backend/internal/shared/contract/providertranslation"
	"golang.org/x/text/language"
)

type ProviderTranslationACL interface {
	GetProviderTranslation(ctx context.Context, providerIdentifier string, preferredLang language.Tag) (*contractTranslation.ProviderTranslationDTO, error)
}
