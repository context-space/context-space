package translation

import (
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	"github.com/context-space/context-space/backend/internal/translation/application"
	"github.com/context-space/context-space/backend/internal/translation/infrastructure/persistence"
)

type Module struct {
	providerTranslationService *application.ProviderTranslationService
	obs                        *observability.ObservabilityProvider
}

func NewModule(
	db database.Database,
	observabilityProvider *observability.ObservabilityProvider,
) (*Module, error) {
	providerTranslationRepo := persistence.NewTranslationRepository(db, observabilityProvider)
	providerTranslationService := application.NewProviderTranslationService(providerTranslationRepo, observabilityProvider)

	return &Module{
		providerTranslationService: providerTranslationService,
		obs:                        observabilityProvider,
	}, nil
}

func (m *Module) GetProviderTranslationService() *application.ProviderTranslationService {
	return m.providerTranslationService
}
