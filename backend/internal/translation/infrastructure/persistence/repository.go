package persistence

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/cache"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"github.com/context-space/context-space/backend/internal/shared/serviceerrors"
	"github.com/context-space/context-space/backend/internal/shared/types"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	"github.com/context-space/context-space/backend/internal/translation/domain"
	"gorm.io/gorm"
)

type translationJSON struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Categories  []string           `json:"categories"`
	Permissions []types.Permission `json:"permissions"`
	Operations  []operation        `json:"operations"`
}

type operation struct {
	Identifier  string      `json:"identifier"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  []parameter `json:"parameters"`
}

type parameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TranslationRepository struct {
	db        database.Database
	obs       *observability.ObservabilityProvider
	cache     *cache.LRUCache[string, *domain.ProviderTranslation]
	cacheOnce sync.Once
}

func NewTranslationRepository(db database.Database, observabilityProvider *observability.ObservabilityProvider) *TranslationRepository {
	repo := &TranslationRepository{
		db:  db,
		obs: observabilityProvider,
	}
	repo.cacheOnce.Do(func() {
		repo.cache = cache.NewLRUCache[string, *domain.ProviderTranslation](50, 1*time.Hour)
	})
	return repo
}

func (r *TranslationRepository) GetProviderTranslation(ctx context.Context, providerIdentifier string, languageCode string) (*domain.ProviderTranslation, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "TranslationRepository.GetProviderTranslation")
	defer span.End()

	cacheKey := utils.StringsBuilder(providerIdentifier, ":", languageCode)
	if translation, ok := r.cache.Get(cacheKey); ok {
		return translation, nil
	}

	var translationModel TranslationModel
	result := r.db.WithContext(ctx).Where("provider_identifier = ? AND language_code = ?", providerIdentifier, languageCode).First(&translationModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, serviceerrors.ErrTranslationNotFound
		}
		return nil, result.Error
	}

	providerTranslation, err := r.mapToDomain(&translationModel)
	if err != nil {
		return nil, err
	}

	r.cache.Set(cacheKey, providerTranslation)

	return providerTranslation, nil
}

func (r *TranslationRepository) mapToDomain(model *TranslationModel) (*domain.ProviderTranslation, error) {
	var translations translationJSON
	err := sonic.Unmarshal(model.Translations, &translations)
	if err != nil {
		return nil, err
	}

	translatedOperations := make([]domain.Operation, len(translations.Operations))

	for i, op := range translations.Operations {
		translatedOperations[i] = domain.Operation{
			Identifier:  op.Identifier,
			Name:        op.Name,
			Description: op.Description,
			Parameters:  make([]domain.Parameter, len(op.Parameters)),
		}
		for j, param := range op.Parameters {
			translatedOperations[i].Parameters[j] = domain.Parameter{
				Name:        param.Name,
				Description: param.Description,
			}
		}
	}

	return &domain.ProviderTranslation{
		Identifier:   model.ProviderIdentifier,
		LanguageCode: model.LanguageCode,
		Name:         translations.Name,
		Description:  translations.Description,
		Categories:   translations.Categories,
		Operations:   translatedOperations,
		Permissions:  translations.Permissions,
	}, nil
}
