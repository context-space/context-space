package http

import (
	"github.com/gin-gonic/gin"

	"github.com/context-space/context-space/backend/internal/provideradapter/application"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	httpapi "github.com/context-space/context-space/backend/internal/shared/interfaces/http"
)

// AdapterHandler handles HTTP requests for adapters
type AdapterHandler struct {
	adapterFactory        *application.AdapterFactory
	providerLoaderService *application.ProviderLoaderService
}

// NewAdapterHandler creates a new adapter handler
func NewAdapterHandler(
	adapterFactory *application.AdapterFactory,
	providerLoaderService *application.ProviderLoaderService,
) *AdapterHandler {
	return &AdapterHandler{
		adapterFactory:        adapterFactory,
		providerLoaderService: providerLoaderService,
	}
}

// RegisterRoutes registers the routes for this handler
func (h *AdapterHandler) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	adapters := router.Group("/adapters")
	{
		adapters.GET("", h.ListProviderAdapters)
		adapters.POST("/reload/:identifier", h.ReloadProvider)
		adapters.POST("/reload_all", h.ReloadProviders)
	}
}

// ListProviderAdaptersResponse represents the response for listing provider adapters
type ListProviderAdaptersResponse struct {
	Adapters []ProviderAdapterResponse `json:"adapters"`
}

// ProviderAdapterResponse represents a provider adapter response
type ProviderAdapterResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	AuthType string `json:"auth_type"`
	Loaded   bool   `json:"loaded"`
	Error    string `json:"error,omitempty"`
}

// mapProviderAdapterToResponse maps a provider adapter to a provider adapter response
func mapProviderAdapterToResponse(providerAdapter domain.ProviderAdapterInfo) ProviderAdapterResponse {
	return ProviderAdapterResponse{
		ID:       providerAdapter.Identifier,
		Name:     providerAdapter.Name,
		AuthType: string(providerAdapter.AuthType),
		Loaded:   true,
		Error:    "",
	}
}

// ListProviderAdapters godoc
// @Summary Get all adapters
// @Description Returns all registered adapters
// @Tags adapters
// @Accept json
// @Produce json
// @Success 200 {object} httpapi.Response{data=ListProviderAdaptersResponse} "Success response with adapters data"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /adapters [get]
func (h *AdapterHandler) ListProviderAdapters(c *gin.Context) {
	providers := h.providerLoaderService.GetLoadedProviders()
	resp := ListProviderAdaptersResponse{}
	for _, provider := range providers {
		resp.Adapters = append(resp.Adapters, mapProviderAdapterToResponse(provider))
	}
	httpapi.OK(c, resp, "Adapters retrieved successfully")
}

func (h *AdapterHandler) ReloadProvider(c *gin.Context) {
	identifier := c.Param("identifier")
	if err := h.providerLoaderService.ReloadProvider(c, identifier); err != nil {
		httpapi.InternalServerError(c, err.Error())
		return
	}
	httpapi.OK(c, nil, "Provider reloaded successfully")
}

func (h *AdapterHandler) ReloadProviders(c *gin.Context) {
	if err := h.providerLoaderService.ReloadAllProviders(c); err != nil {
		httpapi.InternalServerError(c, err.Error())
		return
	}
	httpapi.OK(c, nil, "Providers reloaded successfully")
}
