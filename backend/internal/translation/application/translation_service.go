package application

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	contract "github.com/context-space/context-space/backend/internal/shared/contract/providertranslation"
	"github.com/context-space/context-space/backend/internal/translation/domain"
	"golang.org/x/text/language"
)

type ProviderTranslationService struct {
	providerTranslationRepo domain.ProviderTranslationRepository
	obs                     *observability.ObservabilityProvider
}

func NewProviderTranslationService(
	providerTranslationRepo domain.ProviderTranslationRepository,
	obs *observability.ObservabilityProvider,
) *ProviderTranslationService {
	return &ProviderTranslationService{
		providerTranslationRepo: providerTranslationRepo,
		obs:                     obs,
	}
}

func (s *ProviderTranslationService) GetTranslation(ctx context.Context, providerIdentifier string, preferredLang language.Tag) (*contract.ProviderTranslationDTO, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderTranslationService.GetTranslation")
	defer span.End()

	var languageCode string
	switch preferredLang.Parent() {
	case language.English:
		languageCode = "en"
	case language.Chinese:
		languageCode = "zh-CN"
	case language.SimplifiedChinese:
		languageCode = "zh-CN"
	case language.TraditionalChinese:
		languageCode = "zh-TW"
	default:
		languageCode = "en" // Default to English for unsupported languages
	}
	translation, err := s.providerTranslationRepo.GetProviderTranslation(ctx, providerIdentifier, languageCode)
	if err != nil {
		return nil, err
	}

	translationDTO := &contract.ProviderTranslationDTO{
		Identifier:   translation.Identifier,
		LanguageCode: translation.LanguageCode,
		Name:         translation.Name,
		Description:  translation.Description,
		Categories:   translation.Categories,
		Operations:   make([]contract.OperationDTO, len(translation.Operations)),
		Permissions:  translation.Permissions,
	}

	for i, op := range translation.Operations {
		translationDTO.Operations[i] = contract.OperationDTO{
			Identifier:  op.Identifier,
			Name:        op.Name,
			Description: op.Description,
			Parameters:  make([]contract.ParameterDTO, len(op.Parameters)),
		}
		for j, param := range op.Parameters {
			translationDTO.Operations[i].Parameters[j] = contract.ParameterDTO{
				Name:        param.Name,
				Description: param.Description,
			}
		}

	}

	return translationDTO, nil
}
