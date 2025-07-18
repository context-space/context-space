package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/providercore/domain"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	lru "github.com/hashicorp/golang-lru/v2"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// JSON structures for deserializing translations JSONB field
type translationJSON struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Categories  []string     `json:"categories"`
	Permissions []permission `json:"permissions"`
	Operations  []operation  `json:"operations"`
}

type permission struct {
	Identifier  string `json:"identifier"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

// Cache configuration constants
const (
	cachePrefixID         = "provider-id:"
	cachePrefixIdentifier = "provider-identifier:"
	defaultCacheSize      = 10
)

// ProviderRepository implements the domain.ProviderRepository interface
type ProviderRepository struct {
	db    database.Database
	obs   *observability.ObservabilityProvider
	cache *lru.Cache[string, *domain.Provider]
}

// NewProviderRepository creates a new provider repository
func NewProviderRepository(db database.Database, observabilityProvider *observability.ObservabilityProvider) *ProviderRepository {
	cache, err := lru.New[string, *domain.Provider](defaultCacheSize)
	if err != nil {
		// This should never happen with a valid size, but handle it gracefully
		observabilityProvider.Logger.Fatal(context.Background(), "Failed to initialize provider cache", zap.Error(err))
	}

	return &ProviderRepository{
		db:    db,
		obs:   observabilityProvider,
		cache: cache,
	}
}

// getCacheKeyByID returns the cache key for a provider by ID
func (r *ProviderRepository) getCacheKeyByID(id string) string {
	return utils.StringsBuilder(cachePrefixID, id)
}

// getCacheKeyByIdentifier returns the cache key for a provider by identifier
func (r *ProviderRepository) getCacheKeyByIdentifier(identifier string) string {
	return utils.StringsBuilder(cachePrefixIdentifier, identifier)
}

// GetByID retrieves a provider by ID
func (r *ProviderRepository) GetByID(ctx context.Context, id string) (*domain.Provider, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.GetByID")
	defer span.End()

	// Check cache first
	cacheKey := r.getCacheKeyByID(id)
	if cachedProvider, found := r.cache.Get(cacheKey); found {
		return cachedProvider, nil
	}

	var model ProviderModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Get the provider without operations first
	provider, err := r.mapToDomain(&model)
	if err != nil {
		return nil, err
	}

	// Load operations for this provider
	operations, err := r.loadOperations(ctx, provider.ID)
	if err != nil {
		return nil, err
	}

	// Set operations on the provider
	provider.Operations = operations

	// Load translations for this provider
	if err := r.loadTranslations(ctx, provider); err != nil {
		r.obs.Logger.Warn(ctx, "Failed to load translations for provider, continuing without translations",
			zap.String("provider_id", provider.ID),
			zap.Error(err))
		// Continue without translations - GetTranslation will handle fallback
	}

	// Add to cache with both keys
	r.cache.Add(r.getCacheKeyByID(provider.ID), provider)
	r.cache.Add(r.getCacheKeyByIdentifier(provider.Identifier), provider)

	return provider, nil
}

// GetByIdentifier retrieves a provider by identifier
func (r *ProviderRepository) GetByIdentifier(ctx context.Context, identifier string) (*domain.Provider, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.GetByIdentifier")
	defer span.End()

	// Check cache first
	cacheKey := r.getCacheKeyByIdentifier(identifier)
	if cachedProvider, found := r.cache.Get(cacheKey); found {
		return cachedProvider, nil
	}

	var model ProviderModel
	result := r.db.WithContext(ctx).First(&model, "identifier = ?", identifier)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Get the provider without operations first
	provider, err := r.mapToDomain(&model)
	if err != nil {
		return nil, err
	}

	// Load operations for this provider
	operations, err := r.loadOperations(ctx, provider.ID)
	if err != nil {
		return nil, err
	}

	// Set operations on the provider
	provider.Operations = operations

	// Load translations for this provider
	if err := r.loadTranslations(ctx, provider); err != nil {
		r.obs.Logger.Warn(ctx, "Failed to load translations for provider, continuing without translations",
			zap.String("provider_identifier", provider.Identifier),
			zap.Error(err))
		// Continue without translations - GetTranslation will handle fallback
	}

	// Add to cache with both keys
	r.cache.Add(r.getCacheKeyByIdentifier(provider.Identifier), provider)
	r.cache.Add(r.getCacheKeyByID(provider.ID), provider)

	return provider, nil
}

// List returns all providers
func (r *ProviderRepository) List(ctx context.Context) ([]*domain.Provider, error) {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.List")
	defer span.End()

	var models []ProviderModel
	result := r.db.WithContext(ctx).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	providers := make([]*domain.Provider, 0, len(models))
	for i := range models {
		// Check cache first by ID
		cacheKey := r.getCacheKeyByID(models[i].ID)
		if cachedProvider, found := r.cache.Get(cacheKey); found {
			providers = append(providers, cachedProvider)
			continue
		}

		// Cache miss - load from database
		provider, err := r.mapToDomain(&models[i])
		if err != nil {
			return nil, err
		}

		// Load operations for this provider
		operations, err := r.loadOperations(ctx, provider.ID)
		if err != nil {
			return nil, err
		}

		// Set operations on the provider
		provider.Operations = operations

		// Load translations for this provider
		if err := r.loadTranslations(ctx, provider); err != nil {
			r.obs.Logger.Warn(ctx, "Failed to load translations for provider, continuing without translations",
				zap.String("provider_identifier", provider.Identifier),
				zap.Error(err))
			// Continue without translations - GetTranslation will handle fallback
		}

		// Add to cache with both keys for cache warming
		r.cache.Add(r.getCacheKeyByID(provider.ID), provider)
		r.cache.Add(r.getCacheKeyByIdentifier(provider.Identifier), provider)

		providers = append(providers, provider)
	}

	return providers, nil
}

// Create creates a new provider
func (r *ProviderRepository) Create(ctx context.Context, provider *domain.Provider) error {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.Create")
	defer span.End()

	model, err := r.mapToModel(provider)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Create(model)
	return result.Error
}

// Update updates a provider
func (r *ProviderRepository) Update(ctx context.Context, provider *domain.Provider) error {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.Update")
	defer span.End()

	model, err := r.mapToModel(provider)
	if err != nil {
		return err
	}

	// Use Where clause to update by ID to avoid unique constraint issues
	result := r.db.WithContext(ctx).Where("id = ?", provider.ID).Updates(model)
	if result.Error != nil {
		return result.Error
	}

	// Clear cache entries for the updated provider to ensure fresh data on next access
	r.cache.Remove(r.getCacheKeyByID(provider.ID))
	r.cache.Remove(r.getCacheKeyByIdentifier(provider.Identifier))

	return nil
}

// Delete soft-deletes a provider
func (r *ProviderRepository) Delete(ctx context.Context, id string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.Delete")
	defer span.End()

	var providerIdentifier string
	err := r.db.Transaction(ctx, func(tx *gorm.DB) error {
		// First, get the provider identifier for deleting translations
		var providerModel ProviderModel
		result := tx.Where("id = ?", id).First(&providerModel)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return fmt.Errorf("provider not found with id: %s", id)
			}
			return result.Error
		}

		// Store identifier for cache invalidation
		providerIdentifier = providerModel.Identifier

		// Delete operations associated with this provider
		result = tx.Where("provider_id = ?", id).Delete(&OperationModel{})
		if result.Error != nil {
			return result.Error
		}

		// Delete translations associated with this provider
		result = tx.Where("provider_identifier = ?", providerModel.Identifier).Delete(&ProviderTranslationModel{})
		if result.Error != nil {
			return result.Error
		}

		// Finally, delete the provider
		result = tx.Where("id = ?", id).Delete(&ProviderModel{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("provider not found with id: %s", id)
		}

		return nil
	})

	if err == nil {
		// Clear cache entries for the deleted provider
		r.cache.Remove(r.getCacheKeyByID(id))
		r.cache.Remove(r.getCacheKeyByIdentifier(providerIdentifier))
	}

	return err
}

// loadOperations loads all operations for a provider
func (r *ProviderRepository) loadOperations(ctx context.Context, providerID string) ([]domain.Operation, error) {
	// Create operation repository
	operationRepo := NewOperationRepository(r.db, r.obs)

	// Get operations for this provider
	operationPtrs, err := operationRepo.ListByProviderID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	// Convert from []*domain.Operation to []domain.Operation
	operations := make([]domain.Operation, len(operationPtrs))
	for i, op := range operationPtrs {
		operations[i] = *op
	}

	return operations, nil
}

// mapToDomain maps a provider model to a domain provider
func (r *ProviderRepository) mapToDomain(model *ProviderModel) (*domain.Provider, error) {
	var jsonAttributes struct {
		Categories  []string            `json:"categories"`
		Permissions []domain.Permission `json:"permissions"`
		Tags        []string            `json:"tags"`
	}

	if err := sonic.Unmarshal(model.JSONAttributes, &jsonAttributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider info metadata: %w", err)
	}

	authType := domain.ProviderAuthType(model.AuthType)
	status := domain.ProviderStatus(model.Status)

	return &domain.Provider{
		ID:          model.ID,
		Identifier:  model.Identifier,
		Name:        model.Name,
		Description: model.Description,
		AuthType:    authType,
		Status:      status,
		IconURL:     model.IconURL,
		Categories:  jsonAttributes.Categories,
		Permissions: jsonAttributes.Permissions,
		Tags:        jsonAttributes.Tags,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   parseGormDeletedAt(model.DeletedAt),
	}, nil
}

// mapToModel maps a domain provider to a provider model
func (r *ProviderRepository) mapToModel(provider *domain.Provider) (*ProviderModel, error) {
	jsonAttributes := struct {
		Categories  []string            `json:"categories"`
		Permissions []domain.Permission `json:"permissions"`
	}{
		Categories:  provider.Categories,
		Permissions: provider.Permissions,
	}

	jsonAttributesJSON, err := sonic.Marshal(jsonAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal provider info metadata: %w", err)
	}

	return &ProviderModel{
		ID:             provider.ID,
		Identifier:     provider.Identifier,
		Name:           provider.Name,
		Description:    provider.Description,
		AuthType:       string(provider.AuthType),
		Status:         string(provider.Status),
		IconURL:        provider.IconURL,
		JSONAttributes: jsonAttributesJSON,
		CreatedAt:      provider.CreatedAt,
		UpdatedAt:      provider.UpdatedAt,
		DeletedAt:      parseDomainDeletedAt(provider.DeletedAt),
	}, nil
}

// loadTranslations loads all translations for a provider
func (r *ProviderRepository) loadTranslations(ctx context.Context, provider *domain.Provider) error {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.loadTranslations")
	defer span.End()

	var translationModels []ProviderTranslationModel
	result := r.db.WithContext(ctx).Where("provider_identifier = ?", provider.Identifier).Find(&translationModels)
	if result.Error != nil {
		r.obs.Logger.Warn(ctx, "Failed to load translations for provider",
			zap.String("provider_identifier", provider.Identifier),
			zap.Error(result.Error))
		return result.Error
	}

	// Process each translation record
	for _, model := range translationModels {
		var keyTag language.Tag
		switch model.LanguageCode {
		case "en":
			keyTag = language.English
		case "zh-CN":
			keyTag = language.SimplifiedChinese
		case "zh-TW":
			keyTag = language.TraditionalChinese
		default:
			r.obs.Logger.Warn(ctx, "Unsupported language code in translation data, skipping",
				zap.String("language_code", model.LanguageCode),
				zap.String("provider_identifier", provider.Identifier))
			continue // Skip this translation if the code is not one of the three expected
		}

		// Unmarshal translation JSON
		var tj translationJSON
		if err := sonic.Unmarshal(model.Translations, &tj); err != nil {
			r.obs.Logger.Warn(ctx, "Failed to unmarshal translation JSON, skipping",
				zap.String("language_code", model.LanguageCode),
				zap.String("provider_identifier", provider.Identifier),
				zap.Error(err))
			continue
		}

		// Translate Permissions
		translatedPerms := make([]domain.Permission, len(provider.Permissions))
		copy(translatedPerms, provider.Permissions)

		// Create permission translation map for efficient lookup
		permTranslationMap := make(map[string]permission)
		for _, perm := range tj.Permissions {
			permTranslationMap[perm.Identifier] = perm
		}

		// Apply permission translations
		for i := range translatedPerms {
			if permTranslation, exists := permTranslationMap[translatedPerms[i].Identifier]; exists {
				translatedPerms[i].Name = permTranslation.Name
				translatedPerms[i].Description = permTranslation.Description
				// Deep copy OAuthScopes for data isolation
				clonedScopes := make([]string, len(translatedPerms[i].OAuthScopes))
				copy(clonedScopes, translatedPerms[i].OAuthScopes)
				translatedPerms[i].OAuthScopes = clonedScopes
			}
		}

		// Translate Operations and their Parameters
		translatedOps := make([]domain.Operation, len(provider.Operations))
		copy(translatedOps, provider.Operations)

		// Create operation translation map for efficient lookup
		opTranslationMap := make(map[string]operation)
		for _, op := range tj.Operations {
			opTranslationMap[op.Identifier] = op
		}

		// Apply operation translations
		for i := range translatedOps {
			// Deep copy Parameters
			clonedParams := make([]domain.Parameter, len(translatedOps[i].Parameters))
			copy(clonedParams, translatedOps[i].Parameters)
			translatedOps[i].Parameters = clonedParams

			if opTranslation, exists := opTranslationMap[translatedOps[i].Identifier]; exists {
				translatedOps[i].Name = opTranslation.Name
				translatedOps[i].Description = opTranslation.Description

				// Create parameter translation map for efficient lookup
				paramTranslationMap := make(map[string]parameter)
				for _, param := range opTranslation.Parameters {
					paramTranslationMap[param.Name] = param
				}

				// Apply parameter translations
				for j := range translatedOps[i].Parameters {
					if paramTranslation, exists := paramTranslationMap[translatedOps[i].Parameters[j].Name]; exists {
						translatedOps[i].Parameters[j].Description = paramTranslation.Description
					}
				}
			}
		}

		// Construct TranslatedProvider
		translatedProvider := domain.TranslatedProvider{
			ID:          provider.ID,
			Identifier:  provider.Identifier,
			Name:        tj.Name,
			Description: tj.Description,
			AuthType:    provider.AuthType,
			Status:      provider.Status,
			IconURL:     provider.IconURL,
			Categories:  tj.Categories,
			Tags:        provider.Tags,
			Permissions: translatedPerms,
			Operations:  translatedOps,
			CreatedAt:   provider.CreatedAt,
			UpdatedAt:   provider.UpdatedAt,
			DeletedAt:   provider.DeletedAt,
		}

		// Set translation on provider
		// The keyTag is now guaranteed to be one of language.English, language.SimplifiedChinese, or language.TraditionalChinese
		provider.SetTranslation(keyTag, translatedProvider)
	}

	return nil
}

// SyncTagsToProvider syncs tags to provider's json_attributes
func (r *ProviderRepository) SyncTagsToProvider(ctx context.Context, providerIdentifier string, tags []string) error {
	ctx, span := r.obs.Tracer.Start(ctx, "ProviderRepository.SyncTagsToProvider")
	defer span.End()

	// Get existing provider
	var model ProviderModel
	result := r.db.WithContext(ctx).First(&model, "identifier = ?", providerIdentifier)
	if result.Error != nil {
		return result.Error
	}

	// Parse existing json_attributes
	var jsonAttributes struct {
		Categories  []string            `json:"categories"`
		Permissions []domain.Permission `json:"permissions"`
		Tags        []string            `json:"tags"`
	}

	if err := sonic.Unmarshal(model.JSONAttributes, &jsonAttributes); err != nil {
		return fmt.Errorf("failed to unmarshal provider json_attributes: %w", err)
	}

	// Update tags field
	jsonAttributes.Tags = tags

	// Re-serialize
	updatedJSON, err := sonic.Marshal(jsonAttributes)
	if err != nil {
		return fmt.Errorf("failed to marshal provider json_attributes: %w", err)
	}

	// Update database
	result = r.db.WithContext(ctx).Model(&model).Where("identifier = ?", providerIdentifier).Update("json_attributes", updatedJSON)
	if result.Error != nil {
		return result.Error
	}

	// Clear cache
	r.cache.Remove(r.getCacheKeyByIdentifier(providerIdentifier))
	r.cache.Remove(r.getCacheKeyByID(model.ID))

	return nil
}
