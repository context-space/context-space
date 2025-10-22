package http

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	"github.com/context-space/context-space/backend/internal/provideradapter/application"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/shared/apierrors"
	contractAdapter "github.com/context-space/context-space/backend/internal/shared/contract/provideradapter"
	httpapi "github.com/context-space/context-space/backend/internal/shared/interfaces/http"
	"github.com/context-space/context-space/backend/internal/shared/interfaces/http/middleware"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

// AdapterHandler handles HTTP requests for adapters
type AdapterHandler struct {
	adapterFactory         *application.AdapterFactory
	providerAdapterService *application.ProviderAdapterService
	providerLoaderService  *application.ProviderLoaderService
}

// NewAdapterHandler creates a new adapter handler
func NewAdapterHandler(
	adapterFactory *application.AdapterFactory,
	providerAdapterService *application.ProviderAdapterService,
	providerLoaderService *application.ProviderLoaderService,
) *AdapterHandler {
	return &AdapterHandler{
		adapterFactory:         adapterFactory,
		providerAdapterService: providerAdapterService,
		providerLoaderService:  providerLoaderService,
	}
}

// RegisterRoutes registers the routes for this handler
func (h *AdapterHandler) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	adapters := router.Group("/adapters")
	adapters.Use(middleware.I18nMiddleware())
	{
		// adapters.GET("", h.ListProviderAdapters)
		adapters.GET("/:identifier", h.GetProviderAdapterByIdentifier)
	}
}

// ListProviderAdaptersResponse represents the response for listing provider adapters
type ListProviderAdaptersResponse struct {
	Adapters []ProviderAdapterResponse `json:"adapters"`
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

// OperationResponse represents an operation in provider response
type OperationResponse struct {
	ID                  string               `json:"id"`
	Identifier          string               `json:"identifier"`
	Name                string               `json:"name"`
	Description         string               `json:"description"`
	Category            string               `json:"category"`
	RequiredPermissions []PermissionResponse `json:"required_permissions"`
	Parameters          []ParameterResponse  `json:"parameters,omitempty"`
}

// ParameterResponse represents a parameter in operation response
type ParameterResponse struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// PermissionResponse represents a permission in provider response
type PermissionResponse struct {
	Identifier  string `json:"identifier"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProviderAdapterResponse represents the provider adapter response
type ProviderAdapterResponse struct {
	ID          string               `json:"id"`
	Identifier  string               `json:"identifier"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	AuthType    string               `json:"auth_type"`
	Status      string               `json:"status"`
	IconURL     string               `json:"icon_url"`
	Categories  []string             `json:"categories"`
	Permissions []PermissionResponse `json:"permissions"`
	Operations  []OperationResponse  `json:"operations"`
	Loaded      bool                 `json:"loaded,omitempty"`
	Error       string               `json:"error,omitempty"`
}

// GetProviderAdapterByIdentifier godoc
// @Summary Get provider adapter by identifier
// @Description Gets a provider adapter by its unique identifier, including permissions from adapter table
// @Tags adapters
// @Accept json
// @Produce json
// @Param identifier path string true "Provider Identifier"
// @Success 200 {object} httpapi.Response{data=ProviderAdapterResponse} "Success response with provider adapter data"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /adapters/{identifier} [get]
func (h *AdapterHandler) GetProviderAdapterByIdentifier(c *gin.Context) {
	ctx := c.Request.Context()
	identifier := c.Param("identifier")
	preferredLang := getPreferredLanguage(c)

	provider, err := h.providerAdapterService.GetProviderAdapterByIdentifier(ctx, identifier, preferredLang)
	if err != nil {
		if apiErr, ok := err.(*apierrors.APIError); ok && apiErr.Code == apierrors.ErrNotFound {
			httpapi.NotFound(c, "Provider adapter not found")
		} else {
			httpapi.InternalServerError(c, "Failed to get provider adapter")
		}
		return
	}

	// Map to provider adapter response
	response := ProviderAdapterResponse{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
		Status:      provider.Status,
		IconURL:     provider.IconURL,
		Categories:  provider.Categories,
		Permissions: mapPermissionsToResponse(provider.Permissions),
		Operations:  mapOperationsToResponse(provider.Operations),
	}

	httpapi.OK(c, response, "Provider adapter retrieved successfully")
}

// mapOperationsToResponse maps operation DTOs to response format
func mapOperationsToResponse(operations []contractAdapter.OperationDTO) []OperationResponse {
	responses := make([]OperationResponse, len(operations))
	for i, op := range operations {
		responses[i] = OperationResponse{
			Identifier:          op.Identifier,
			Name:                op.Name,
			Description:         op.Description,
			Category:            op.Category,
			RequiredPermissions: mapPermissionsToResponse(op.RequiredPermissions),
			Parameters:          mapParametersToResponse(op.Parameters),
		}
	}
	return responses
}

// mapPermissionsToResponse maps permission DTOs to response format
func mapPermissionsToResponse(permissions []types.Permission) []PermissionResponse {
	responses := make([]PermissionResponse, len(permissions))
	for i, perm := range permissions {
		responses[i] = PermissionResponse{
			Identifier:  perm.Identifier,
			Name:        perm.Name,
			Description: perm.Description,
		}
	}
	return responses
}

// mapParametersToResponse maps parameter DTOs to response format
func mapParametersToResponse(parameters []contractAdapter.ParameterDTO) []ParameterResponse {
	responses := make([]ParameterResponse, len(parameters))
	for i, param := range parameters {
		responses[i] = ParameterResponse{
			Name:        param.Name,
			Type:        param.Type,
			Description: param.Description,
			Required:    param.Required,
			Enum:        param.Enum,
			Default:     param.Default,
		}
	}
	return responses
}

// getPreferredLanguage extracts the preferred language from context (set by middleware)
// Falls back to parsing Accept-Language header if middleware is not used
func getPreferredLanguage(c *gin.Context) language.Tag {
	// First try to get from context (set by middleware)
	if lang, exists := c.Get(string(middleware.PreferredLanguageKey)); exists {
		if langTag, ok := lang.(language.Tag); ok {
			return langTag
		}
		// If value from context is not language.Tag (e.g. due to misconfiguration or different middleware),
		// fall through to header parsing as a safety measure.
	}

	// Fallback: parse Accept-Language header directly
	acceptLang := c.GetHeader("Accept-Language")
	if acceptLang == "" {
		return language.English // Default to English if no header
	}

	tags, _, err := language.ParseAcceptLanguage(acceptLang)
	if err != nil || len(tags) == 0 {
		return language.English // Default to English if parsing fails
	}

	return tags[0] // Return the highest priority language
}
