package http

import (
	"net/http"

	"github.com/context-space/context-space/backend/internal/identityaccess/domain"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/identityaccess/application"
	"github.com/context-space/context-space/backend/internal/shared/apierrors"
	httpapi "github.com/context-space/context-space/backend/internal/shared/interfaces/http"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user management
type UserHandler struct {
	userService *application.UserService
	obs         *observability.ObservabilityProvider
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *application.UserService, observabilityProvider *observability.ObservabilityProvider) *UserHandler {
	return &UserHandler{
		userService: userService,
		obs:         observabilityProvider,
	}
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	users := router.Group("/users")
	users.Use(requireAuth)
	{
		// Current user endpoint
		users.GET("/me", h.GetCurrentUser)
		users.DELETE("/me", h.DeleteCurrentUser)

		// API Key routes
		users.POST("/me/apikeys", h.CreateAPIKey)
		users.GET("/me/apikeys", h.ListAPIKeys)
		users.GET("/me/apikeys/:keyID", h.GetAPIKey)
		users.DELETE("/me/apikeys/:keyID", h.DeleteAPIKey)
	}
}

// UserResponse represents the user response
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Returns the current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httpapi.Response{data=UserResponse} "Success response with user data"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*domain.User)

	httpapi.OK(c, UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}, "User retrieved successfully")
}

// DeleteCurrentUser godoc
// @Summary Delete current user
// @Description Deletes the current user's account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 204 "No content success response"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request error response"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Router /users/me [delete]
func (h *UserHandler) DeleteCurrentUser(c *gin.Context) {
	// Not implemented yet
	httpapi.BadRequest(c, "Not implemented")
}

// CreateAPIKeyRequest represents the request to create an API key
type CreateAPIKeyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// APIKeyResponse represents the API key response
type APIKeyResponse struct {
	ID          string `json:"id"`
	KeyValue    string `json:"key_value,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

// CreateAPIKey godoc
// @Summary Create API key
// @Description Creates a new API key for the current user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAPIKeyRequest true "Create API key request"
// @Success 201 {object} httpapi.Response{data=APIKeyResponse} "Success response with created API key data"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request error response"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /users/me/apikeys [post]
func (h *UserHandler) CreateAPIKey(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*domain.User)

	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.BadRequest(c, utils.StringsBuilder("Invalid request: ", err.Error()))
		return
	}

	apiKey, err := h.userService.CreateAPIKey(ctx, user.ID, req.Name, req.Description)
	if err != nil {
		if err.(*apierrors.APIError).Code == apierrors.ErrNotFound {
			httpapi.NotFound(c, "User not found")
		} else {
			httpapi.InternalServerError(c, "Failed to create API key")
		}
		return
	}

	httpapi.Created(c, APIKeyResponse{
		ID:          apiKey.ID,
		KeyValue:    apiKey.KeyValue,
		Name:        apiKey.Name,
		Description: apiKey.Description,
		CreatedAt:   apiKey.CreatedAt.Format(http.TimeFormat),
	}, "API key created successfully")
}

// ListAPIKeysResponse represents the response for listing API keys
type ListAPIKeysResponse struct {
	APIKeys []APIKeyResponse `json:"api_keys"`
}

// ListAPIKeys godoc
// @Summary List API keys
// @Description Lists all API keys for the current user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httpapi.Response{data=ListAPIKeysResponse} "Success response with list of API keys"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /users/me/apikeys [get]
func (h *UserHandler) ListAPIKeys(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*domain.User)

	apiKeys, err := h.userService.ListAPIKeys(ctx, user.ID)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to list API keys")
		return
	}

	// Map to response
	apiKeyResponses := make([]APIKeyResponse, len(apiKeys))
	for i, apiKey := range apiKeys {
		apiKeyResponses[i] = APIKeyResponse{
			ID:          apiKey.ID,
			Name:        apiKey.Name,
			KeyValue:    apiKey.KeyValue[3:11],
			Description: apiKey.Description,
			CreatedAt:   apiKey.CreatedAt.Format(http.TimeFormat),
		}
	}

	httpapi.OK(c, ListAPIKeysResponse{
		APIKeys: apiKeyResponses,
	}, "API keys retrieved successfully")
}

// GetAPIKey godoc
// @Summary Get API key
// @Description Gets a specific API key by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyID path string true "API Key ID"
// @Success 200 {object} httpapi.Response{data=APIKeyResponse} "Success response with API key data"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 403 {object} httpapi.SwaggerErrorResponse "Forbidden error response"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Router /users/me/apikeys/{keyID} [get]
func (h *UserHandler) GetAPIKey(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*domain.User)

	keyID := c.Param("keyID")

	apiKey, err := h.userService.GetAPIKeyByID(ctx, keyID)
	if err != nil {
		httpapi.NotFound(c, "API key not found")
		return
	}

	// Verify that the API key belongs to the user
	if apiKey.UserID != user.ID {
		httpapi.Forbidden(c, "API key does not belong to this user")
		return
	}

	httpapi.OK(c, APIKeyResponse{
		ID:          apiKey.ID,
		KeyValue:    apiKey.KeyValue,
		Name:        apiKey.Name,
		Description: apiKey.Description,
		CreatedAt:   apiKey.CreatedAt.Format(http.TimeFormat),
	}, "API key retrieved successfully")
}

// DeleteAPIKey godoc
// @Summary Delete API key
// @Description Deletes a specific API key by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyID path string true "API Key ID"
// @Success 204 "No content success response"
// @Failure 401 {object} httpapi.Response "Unauthorized error response"
// @Failure 403 {object} httpapi.SwaggerErrorResponse "Forbidden error response"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /users/me/apikeys/{keyID} [delete]
func (h *UserHandler) DeleteAPIKey(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*domain.User)

	keyID := c.Param("keyID")

	// Verify that the API key exists and belongs to the user
	apiKey, err := h.userService.GetAPIKeyByID(ctx, keyID)
	if err != nil {
		if err.(*apierrors.APIError).Code == apierrors.ErrNotFound {
			httpapi.NotFound(c, "API key not found")
		} else {
			httpapi.InternalServerError(c, "Failed to delete API key")
		}
		return
	}

	if apiKey.UserID != user.ID {
		httpapi.Forbidden(c, "API key does not belong to this user")
		return
	}

	if err := h.userService.DeleteAPIKey(ctx, keyID); err != nil {
		httpapi.InternalServerError(c, "Failed to delete API key")
		return
	}

	httpapi.NoContent(c)
}
