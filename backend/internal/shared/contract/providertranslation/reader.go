package providertranslation

import (
	"context"

	"golang.org/x/text/language"
)

type ProviderTranslationReader interface {
	GetTranslation(ctx context.Context, providerIdentifier string, preferredLang language.Tag) (*ProviderTranslationDTO, error)
}
