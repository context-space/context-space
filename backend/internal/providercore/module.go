package providercore

import (
	"context"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/providercore/application"
	"github.com/context-space/context-space/backend/internal/providercore/infrastructure/acl"
	"github.com/context-space/context-space/backend/internal/providercore/infrastructure/persistence"
	pchttp "github.com/context-space/context-space/backend/internal/providercore/interfaces/http"
	"github.com/context-space/context-space/backend/internal/shared/events"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
	translation "github.com/context-space/context-space/backend/internal/translation/application"
	"github.com/gin-gonic/gin"
)

// Module encapsulates all provider components
type Module struct {
	providerService *application.ProviderService
	providerHandler *pchttp.ProviderHandler
	obs             *observability.ObservabilityProvider
}

// NewModule creates a new provider module
func NewModule(
	db database.Database,
	eventBus *events.Bus,
	observabilityProvider *observability.ObservabilityProvider,
	providerTranslationService *translation.ProviderTranslationService,
) (*Module, error) {
	// Create repositories
	providerRepo := persistence.NewProviderRepository(db, observabilityProvider)
	operationRepo := persistence.NewOperationRepository(db, observabilityProvider)
	providerTranslationACL := acl.NewProviderTranslationACL(providerTranslationService, observabilityProvider)

	// Create application service
	providerService := application.NewProviderService(
		providerRepo,
		operationRepo,
		providerTranslationACL,
		eventBus,
		observabilityProvider,
	)

	// Create HTTP handler
	providerHandler := pchttp.NewProviderHandler(providerService, observabilityProvider)

	return &Module{
		providerService: providerService,
		providerHandler: providerHandler,
		obs:             observabilityProvider,
	}, nil
}

// Initialize initializes the provider module
func (m *Module) Initialize(ctx context.Context) error {
	// Initialize any providers or operations from configuration or external sources
	return nil
}

// RegisterRoutes registers all Provider HTTP routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	// Register provider routes
	m.providerHandler.RegisterRoutes(router)
}

// GetProviderService returns the provider service instance
func (m *Module) GetProviderService() *application.ProviderService {
	return m.providerService
}
