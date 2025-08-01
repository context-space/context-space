package provideradapter

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	observability "github.com/context-space/cloud-observability"

	"github.com/context-space/context-space/backend/internal/provideradapter/application"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/acl"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/persistence"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/templates"
	"github.com/context-space/context-space/backend/internal/provideradapter/interfaces/http"
	providercore "github.com/context-space/context-space/backend/internal/providercore/application"
	"github.com/context-space/context-space/backend/internal/shared/infrastructure/database"
)

// Module encapsulates all provider adapter components
type Module struct {
	adapterFactory        *application.AdapterFactory
	providerLoaderService *application.ProviderLoaderService
	adapterHandler        *http.AdapterHandler
	obs                   *observability.ObservabilityProvider
}

// NewModule creates a new provider adapter module
func NewModule(
	db database.Database,
	observabilityProvider *observability.ObservabilityProvider,
	providerCoreService *providercore.ProviderService,
) (*Module, error) {
	// Initialize adapter factory
	adapterFactory := application.NewAdapterFactory()

	providerCoreACL := acl.NewProviderCoreACL(providerCoreService, observabilityProvider)

	// Create repositories
	providerRepo := persistence.NewAdapterRepository(db, observabilityProvider)

	providerLoader := registry.NewProviderLoader(adapterFactory)

	// Initialize provider loader
	providerLoaderService := application.NewProviderLoaderService(
		providerCoreACL,
		providerRepo,
		providerLoader,
		observabilityProvider,
	)

	// Initialize all provider templates
	// This imports the templates package which registers all provider templates
	// via their init() functions
	templates.Init()

	// Initialize HTTP handlers
	adapterHandler := http.NewAdapterHandler(
		adapterFactory,
		providerLoaderService,
	)

	return &Module{
		adapterFactory:        adapterFactory,
		providerLoaderService: providerLoaderService,
		adapterHandler:        adapterHandler,
		obs:                   observabilityProvider,
	}, nil
}

// RegisterRoutes registers all provider adapter HTTP routes
func (m *Module) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	// Now we can directly register with the router group
	m.adapterHandler.RegisterRoutes(router, requireAuth)
}

// Initialize loads all provider adapters from configuration
func (m *Module) Initialize(ctx context.Context) error {
	ctx, span := m.obs.Tracer.Start(ctx, "ProviderAdapterModule.Initialize")
	defer span.End()

	m.obs.Logger.Info(ctx, "Initializing Provider Adapter module")

	err := m.providerLoaderService.LoadAllProviders(ctx)
	if err != nil {
		m.obs.Logger.Error(ctx, "Failed to load providers", zap.Error(err))
		return fmt.Errorf("failed to initialize provider adapters: %w", err)
	}

	loadedProviders := m.providerLoaderService.GetLoadedProviders()

	m.obs.Logger.Info(ctx, "Provider Adapter module initialized successfully",
		zap.Int("total_providers", len(loadedProviders)))

	return nil
}

// GetAdapter returns an adapter for the given provider ID
func (m *Module) GetAdapter(providerID string) (domain.Adapter, error) {
	return m.adapterFactory.GetAdapter(providerID)
}

// GetAdapterFactory returns the adapter factory
func (m *Module) GetAdapterFactory() *application.AdapterFactory {
	return m.adapterFactory
}

func (m *Module) GetProviderLoaderService() *application.ProviderLoaderService {
	return m.providerLoaderService
}
