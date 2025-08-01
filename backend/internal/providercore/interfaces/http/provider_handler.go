package http

import (
	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/providercore/application"
	"github.com/context-space/context-space/backend/internal/providercore/domain"
	"github.com/context-space/context-space/backend/internal/providercore/interfaces/http/middleware"
	"github.com/context-space/context-space/backend/internal/shared/apierrors"
	httpapi "github.com/context-space/context-space/backend/internal/shared/interfaces/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// ProviderHandler handles HTTP requests for provider management
type ProviderHandler struct {
	providerService *application.ProviderService
	obs             *observability.ObservabilityProvider
}

// NewProviderHandler creates a new provider handler
func NewProviderHandler(providerService *application.ProviderService, observabilityProvider *observability.ObservabilityProvider) *ProviderHandler {
	return &ProviderHandler{
		providerService: providerService,
		obs:             observabilityProvider,
	}
}

// RegisterRoutes registers the provider routes
func (h *ProviderHandler) RegisterRoutes(router *gin.RouterGroup) {
	providers := router.Group("/providers")
	providers.Use(middleware.I18nMiddleware()) // Add i18n middleware
	{
		providers.GET("", h.ListProviders)
		providers.GET("/categories", h.GetProviderCategories)
		providers.POST("/filter", h.FilterProviders)
		providers.GET("/:identifier", h.GetProviderByIdentifier)
	}
}

// getPreferredLanguage extracts the preferred language from context (set by middleware)
// Falls back to parsing Accept-Language header if middleware is not used
func getPreferredLanguage(c *gin.Context) language.Tag {
	// First try to get from context (set by middleware)
	if lang, exists := c.Get(string(domain.PreferredLanguageKey)); exists {
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

// ProviderResponse represents the provider response
type ProviderResponse struct {
	ID          string               `json:"id"`
	Identifier  string               `json:"identifier"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	AuthType    string               `json:"auth_type"`
	Status      string               `json:"status"`
	IconURL     string               `json:"icon_url"`
	ApiDocURL   string               `json:"api_doc_url"`
	Categories  []string             `json:"categories"`
	Permissions []PermissionResponse `json:"permissions"`
	Operations  []OperationResponse  `json:"operations"`
}

// BriefProviderResponse represents a brief provider response
type BriefProviderResponse struct {
	ID          string   `json:"id"`
	Identifier  string   `json:"identifier"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	AuthType    string   `json:"auth_type"`
	Status      string   `json:"status"`
	IconURL     string   `json:"icon_url"`
	ApiDocURL   string   `json:"api_doc_url"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
}

// mapProviderToResponse maps a domain provider to a provider response
func mapProviderToResponse(provider *domain.Provider, lang language.Tag) ProviderResponse {
	// Get translated provider
	translatedProvider := provider.GetTranslation(lang)
	// Map operations
	operations := make([]OperationResponse, len(translatedProvider.Operations))
	for i, op := range translatedProvider.Operations {
		parameters := make([]ParameterResponse, len(op.Parameters))
		for j, param := range op.Parameters {
			parameters[j] = ParameterResponse{
				Name:        param.Name,
				Type:        string(param.Type),
				Description: param.Description,
				Required:    param.Required,
				Enum:        param.Enum,
				Default:     param.Default,
			}
		}
		permissions := make([]PermissionResponse, len(op.RequiredPermissions))
		for j, permDef := range op.RequiredPermissions {
			// Find the translated permission by identifier
			var translatedPerm domain.Permission
			found := false
			for _, p := range translatedProvider.Permissions {
				if p.Identifier == permDef.Identifier {
					translatedPerm = p
					found = true
					break
				}
			}
			if !found {
				// Fallback to original permission if translation not found
				translatedPerm = permDef
			}

			permissions[j] = PermissionResponse{
				Identifier:  translatedPerm.Identifier,
				Name:        translatedPerm.Name,
				Description: translatedPerm.Description,
			}
		}
		operations[i] = OperationResponse{
			ID:                  op.ID,
			Identifier:          op.Identifier,
			Name:                op.Name,
			Description:         op.Description,
			Category:            op.Category,
			RequiredPermissions: permissions,
			Parameters:          parameters,
		}
	}

	// Map permissions
	permissions := make([]PermissionResponse, len(translatedProvider.Permissions))
	for i, perm := range translatedProvider.Permissions {
		permissions[i] = PermissionResponse{
			Identifier:  perm.Identifier,
			Name:        perm.Name,
			Description: perm.Description,
		}
	}

	return ProviderResponse{
		ID:          translatedProvider.ID,
		Identifier:  translatedProvider.Identifier,
		Name:        translatedProvider.Name,
		Description: translatedProvider.Description,
		AuthType:    string(translatedProvider.AuthType),
		Status:      string(translatedProvider.Status),
		IconURL:     translatedProvider.IconURL,
		ApiDocURL:   translatedProvider.ApiDocURL,
		Categories:  translatedProvider.Categories,
		Permissions: permissions,
		Operations:  operations,
	}
}

// mapProviderToBriefResponse maps a domain provider to a brief provider response
func mapProviderToBriefResponse(provider *domain.Provider, lang language.Tag) BriefProviderResponse {
	// Get translated provider
	translatedProvider := provider.GetTranslation(lang)

	return BriefProviderResponse{
		ID:          translatedProvider.ID,
		Identifier:  translatedProvider.Identifier,
		Name:        translatedProvider.Name,
		Description: translatedProvider.Description,
		AuthType:    string(translatedProvider.AuthType),
		Status:      string(translatedProvider.Status),
		IconURL:     translatedProvider.IconURL,
		ApiDocURL:   translatedProvider.ApiDocURL,
		Categories:  translatedProvider.Categories,
		Tags:        translatedProvider.Tags,
	}
}

// ListProvidersResponse represents the response for listing providers
type ListProvidersResponse struct {
	Providers []BriefProviderResponse `json:"providers"`
}

// ListProviders godoc
// @Summary List providers
// @Description Lists all available providers
// @Tags providers
// @Accept json
// @Produce json
// @Success 200 {object} httpapi.Response{data=ListProvidersResponse} "Success response with list of providers"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /providers [get]
func (h *ProviderHandler) ListProviders(c *gin.Context) {
	ctx := c.Request.Context()
	preferredLang := getPreferredLanguage(c)

	providers, err := h.providerService.ListProviders(ctx)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to list providers")
		return
	}

	// Map to response
	providerResponses := make([]BriefProviderResponse, 0, len(providers))
	for _, provider := range providers {
		if provider.Status == domain.ProviderStatusDeprecated {
			continue
		}
		providerResponses = append(providerResponses, mapProviderToBriefResponse(provider, preferredLang))
	}

	httpapi.OK(c, ListProvidersResponse{
		Providers: providerResponses,
	}, "Providers retrieved successfully")
}

// GetProviderByIdentifier godoc
// @Summary Get provider by identifier
// @Description Gets a provider by its unique identifier
// @Tags providers
// @Accept json
// @Produce json
// @Param identifier path string true "Provider Identifier"
// @Success 200 {object} httpapi.Response{data=ProviderResponse} "Success response with provider data"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /providers/{identifier} [get]
func (h *ProviderHandler) GetProviderByIdentifier(c *gin.Context) {
	ctx := c.Request.Context()
	identifier := c.Param("identifier")
	preferredLang := getPreferredLanguage(c)

	provider, err := h.providerService.GetProviderByIdentifier(ctx, identifier)
	if err != nil {
		if apiErr, ok := err.(*apierrors.APIError); ok && apiErr.Code == apierrors.ErrNotFound {
			httpapi.NotFound(c, "Provider not found")
		} else {
			httpapi.InternalServerError(c, "Failed to get provider")
		}
		return
	}

	httpapi.OK(c, mapProviderToResponse(provider, preferredLang), "Provider retrieved successfully")
}

// CategoriesResponse represents the response for listing provider categories
type CategoriesResponse struct {
	Categories []string `json:"categories"`
}

// GetProviderCategories godoc
// @Summary Get provider categories
// @Description Gets all available provider categories
// @Tags providers
// @Accept json
// @Produce json
// @Success 200 {object} httpapi.Response{data=CategoriesResponse} "Success response with list of categories"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /providers/categories [get]
func (h *ProviderHandler) GetProviderCategories(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all providers to extract categories
	providers, err := h.providerService.ListProviders(ctx)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to get provider categories")
		return
	}

	// Extract unique categories
	categoryMap := make(map[string]bool)
	for _, provider := range providers {
		for _, category := range provider.Categories {
			categoryMap[category] = true
		}
	}

	// Convert map to slice
	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}

	httpapi.OK(c, CategoriesResponse{
		Categories: categories,
	}, "Provider categories retrieved successfully")
}

// FilterProvidersRequest represents the request for filtering providers
type FilterProvidersRequest struct {
	Filters    FilterCriteria    `json:"filters"`
	Pagination PaginationRequest `json:"pagination"`
	Sort       SortRequest       `json:"sort"`
}

// FilterCriteria represents the filtering criteria for providers
type FilterCriteria struct {
	Tag          string `json:"tag,omitempty"`
	AuthType     string `json:"auth_type,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
}

// PaginationRequest represents the pagination parameters
type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// SortRequest represents the sorting parameters
type SortRequest struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// FilterProvidersResponse represents the response for filtering providers
type FilterProvidersResponse struct {
	Providers  []BriefProviderResponse `json:"providers"`
	Pagination PaginationResponse      `json:"pagination"`
}

// PaginationResponse represents the pagination metadata in response
type PaginationResponse struct {
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}

// FilterProviders godoc
// @Summary Filter providers
// @Description Filters providers based on various criteria with pagination and sorting
// @Tags providers
// @Accept json
// @Produce json
// @Param request body FilterProvidersRequest true "Filter parameters"
// @Success 200 {object} httpapi.Response{data=FilterProvidersResponse} "Success response with filtered providers"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /providers/filter [post]
func (h *ProviderHandler) FilterProviders(c *gin.Context) {
	ctx := c.Request.Context()
	preferredLang := getPreferredLanguage(c)

	var req FilterProvidersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.BadRequest(c, "Invalid request format")
		return
	}

	// Set default values
	if req.Pagination.Page <= 0 {
		req.Pagination.Page = 1
	}
	if req.Pagination.PageSize <= 0 {
		req.Pagination.PageSize = 20
	}
	if req.Pagination.PageSize > 100 {
		req.Pagination.PageSize = 100
	}
	if req.Sort.Field == "" {
		req.Sort.Field = "created_at"
	}
	if req.Sort.Order == "" {
		req.Sort.Order = "desc"
	}

	// Validate auth_type if provided
	if req.Filters.AuthType != "" {
		validAuthTypes := map[string]bool{
			"oauth":  true,
			"apikey": true,
			"basic":  true,
			"none":   true,
		}
		if !validAuthTypes[req.Filters.AuthType] {
			httpapi.BadRequest(c, "Invalid auth_type. Must be one of: oauth, apikey, basic, none")
			return
		}
	}

	// Validate sort field
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"name":       true,
		"identifier": true,
	}
	if !validSortFields[req.Sort.Field] {
		httpapi.BadRequest(c, "Invalid sort field. Must be one of: created_at, updated_at, name, identifier")
		return
	}

	// Validate sort order
	if req.Sort.Order != "asc" && req.Sort.Order != "desc" {
		httpapi.BadRequest(c, "Invalid sort order. Must be 'asc' or 'desc'")
		return
	}

	// Convert auth_type string to domain type
	var authType domain.ProviderAuthType
	if req.Filters.AuthType != "" {
		switch req.Filters.AuthType {
		case "oauth":
			authType = domain.AuthTypeOAuth
		case "apikey":
			authType = domain.AuthTypeAPIKey
		case "basic":
			authType = domain.AuthTypeBasic
		case "none":
			authType = domain.AuthTypeNone
		}
	}

	// Call service
	params := application.ProviderFilterParams{
		Tag:          req.Filters.Tag,
		AuthType:     authType,
		ProviderName: req.Filters.ProviderName,
		Page:         req.Pagination.Page,
		PageSize:     req.Pagination.PageSize,
		SortField:    req.Sort.Field,
		SortOrder:    req.Sort.Order,
	}

	result, err := h.providerService.FilterProviders(ctx, params)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to filter providers")
		return
	}

	// Map providers to response
	providerResponses := make([]BriefProviderResponse, 0, len(result.Providers))
	for _, provider := range result.Providers {
		if provider.Status == domain.ProviderStatusDeprecated {
			continue
		}
		providerResponses = append(providerResponses, mapProviderToBriefResponse(provider, preferredLang))
	}

	// Build response
	response := FilterProvidersResponse{
		Providers: providerResponses,
		Pagination: PaginationResponse{
			CurrentPage: result.CurrentPage,
			PageSize:    result.PageSize,
			TotalItems:  result.TotalCount,
			TotalPages:  result.TotalPages,
			HasNext:     result.HasNext,
			HasPrev:     result.HasPrev,
		},
	}

	httpapi.OK(c, response, "Providers filtered successfully")
}
