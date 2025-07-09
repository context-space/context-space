package http

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	observability "github.com/context-space/cloud-observability"
	identityDomain "github.com/context-space/context-space/backend/internal/identityaccess/domain"
	"github.com/context-space/context-space/backend/internal/integration/application"
	"github.com/context-space/context-space/backend/internal/integration/domain"
	httpapi "github.com/context-space/context-space/backend/internal/shared/interfaces/http"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	"github.com/gin-gonic/gin"
)

// InvocationHandler handles HTTP requests for invocations
type InvocationHandler struct {
	invocationService *application.InvocationService
	obs               *observability.ObservabilityProvider
}

// NewInvocationHandler creates a new invocation handler
func NewInvocationHandler(
	invocationService *application.InvocationService,
	observabilityProvider *observability.ObservabilityProvider,
) *InvocationHandler {
	return &InvocationHandler{
		invocationService: invocationService,
		obs:               observabilityProvider,
	}
}

// RegisterRoutes registers the routes for this handler
func (h *InvocationHandler) RegisterRoutes(router *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	// Base invocation routes
	invocations := router.Group("/invocations")
	invocations.Use(requireAuth)
	{
		invocations.GET("", h.ListInvocationsByUser)
		invocations.GET("/:invocation_id", h.GetInvocation)
		invocations.POST("/:provider_identifier/:operation_identifier", h.InvokeOperation)
	}
}

// InvocationResponse represents an invocation in responses
type InvocationResponse struct {
	ID                  string                 `json:"id"`
	UserID              string                 `json:"user_id"`
	ProviderIdentifier  string                 `json:"provider_identifier"`
	OperationIdentifier string                 `json:"operation_identifier"`
	Status              string                 `json:"status"`
	Parameters          map[string]interface{} `json:"parameters"`
	ResponseData        json.RawMessage        `json:"response_data,omitempty"`
	ErrorMessage        string                 `json:"error_message,omitempty"`
	Duration            int64                  `json:"duration_ms"`
	StartedAt           string                 `json:"started_at,omitempty"`
	CompletedAt         string                 `json:"completed_at,omitempty"`
	CreatedAt           string                 `json:"created_at"`
}

// ListInvocationsResponse represents the response for listing invocations
type ListInvocationsResponse struct {
	Invocations []InvocationResponse `json:"invocations"`
	Total       int64                `json:"total"`
}

// mapInvocationToResponse maps a domain invocation to a response
func mapInvocationToResponse(invocation *domain.Invocation, withResponseData bool) (InvocationResponse, error) {
	// Format timestamp strings
	startedAt := ""
	if invocation.StartedAt != nil {
		startedAt = invocation.StartedAt.Format(time.RFC3339)
	}
	completedAt := ""
	if invocation.CompletedAt != nil {
		completedAt = invocation.CompletedAt.Format(time.RFC3339)
	}
	createdAt := invocation.CreatedAt.Format(time.RFC3339)

	responseData := json.RawMessage{}
	if withResponseData {
		responseData = invocation.ResponseData
	}

	return InvocationResponse{
		ID:                  invocation.ID,
		UserID:              invocation.UserID,
		ProviderIdentifier:  invocation.ProviderIdentifier,
		OperationIdentifier: invocation.OperationIdentifier,
		Status:              string(invocation.Status),
		Parameters:          invocation.Parameters,
		ResponseData:        responseData,
		ErrorMessage:        invocation.ErrorMessage,
		Duration:            invocation.Duration,
		StartedAt:           startedAt,
		CompletedAt:         completedAt,
		CreatedAt:           createdAt,
	}, nil
}

// GetInvocation godoc
// @Summary Get invocation by ID
// @Description Gets an invocation by its ID
// @Tags invocation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param invocation_id path string true "Invocation ID"
// @Success 200 {object} httpapi.Response{data=InvocationResponse} "Success response with invocation data"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /invocations/{invocation_id} [get]
func (h *InvocationHandler) GetInvocation(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*identityDomain.User)

	invocationID := c.Param("invocation_id")

	invocation, err := h.invocationService.GetInvocationByID(ctx, invocationID)
	if err != nil {
		if errors.Is(err, application.ErrInvocationNotFound) {
			httpapi.NotFound(c, "Invocation not found")
		} else {
			httpapi.InternalServerError(c, "Failed to get invocation")
		}
		return
	}

	if invocation.UserID != user.ID {
		httpapi.Forbidden(c, "You are not allowed to access this invocation")
		return
	}

	// Map to response
	response, err := mapInvocationToResponse(invocation, true)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to format response")
		return
	}

	httpapi.OK(c, response, "Invocation retrieved successfully")
}

// ListInvocations godoc
// @Summary List invocations
// @Description Lists invocations for the authenticated user
// @Tags invocation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit (default: 20)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} httpapi.Response{data=ListInvocationsResponse} "Success response with list of invocations"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /invocations [get]
func (h *InvocationHandler) ListInvocationsByUser(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user ID from context
	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*identityDomain.User)

	// Get pagination parameters
	limit := 20 // Default limit
	offset := 0 // Default offset

	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100 // Cap at 100
			}
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get invocations for the user
	invocations, err := h.invocationService.ListInvocationsByUserID(ctx, user.ID, limit, offset)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to list invocations")
		return
	}

	// Get total count
	total, err := h.invocationService.CountInvocationsByUserID(ctx, user.ID)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to count invocations")
		return
	}

	// Map to response
	invocationResponses := make([]InvocationResponse, 0, len(invocations))
	for _, invocation := range invocations {
		response, err := mapInvocationToResponse(invocation, false)
		if err != nil {
			continue // Skip this one on error
		}
		invocationResponses = append(invocationResponses, response)
	}

	httpapi.OK(c, ListInvocationsResponse{
		Invocations: invocationResponses,
		Total:       total,
	}, "Invocations retrieved successfully")
}

// InvokeRequest represents the request body for invoking an operation
type InvokeRequest struct {
	Parameters map[string]interface{} `json:"parameters"`
}

// InvokeOperation godoc
// @Summary Invoke provider operation
// @Description Executes an operation on a provider
// @Tags invocation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider_identifier path string true "Provider Identifier"
// @Param operation_identifier path string true "Operation Identifier"
// @Param request body InvokeRequest true "Invocation parameters"
// @Success 200 {object} httpapi.Response{data=InvocationResponse} "Success response with invocation result"
// @Failure 400 {object} httpapi.SwaggerErrorResponse "Bad request error response"
// @Failure 401 {object} httpapi.SwaggerErrorResponse "Unauthorized error response"
// @Failure 404 {object} httpapi.SwaggerErrorResponse "Not found error response"
// @Failure 429 {object} httpapi.SwaggerErrorResponse "Rate limit exceeded error response"
// @Failure 500 {object} httpapi.SwaggerErrorResponse "Internal server error response"
// @Router /invocations/{provider_identifier}/{operation_identifier} [post]
func (h *InvocationHandler) InvokeOperation(c *gin.Context) {
	ctx := c.Request.Context()

	userI, exists := c.Get("user")
	if !exists {
		httpapi.Unauthorized(c, "Authentication required")
		return
	}
	user := userI.(*identityDomain.User)

	// Get parameters from the URL
	providerIdentifier := c.Param("provider_identifier")
	operationIdentifier := c.Param("operation_identifier")

	// Parse request body
	var req InvokeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.BadRequest(c, utils.StringsBuilder("Invalid request format: ", err.Error()))
		return
	}

	// Invoke the operation
	invocation, err := h.invocationService.InvokeOperation(
		ctx,
		user.ID,
		providerIdentifier,
		operationIdentifier,
		req.Parameters,
	)

	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, application.ErrProviderNotFound):
			httpapi.NotFound(c, "Provider not found")
		case errors.Is(err, application.ErrProviderAdapterNotFound):
			httpapi.NotFound(c, "Provider adapter not found")
		case errors.Is(err, application.ErrOperationNotFound):
			httpapi.NotFound(c, "Operation not found")
		case errors.Is(err, application.ErrInvalidParameters):
			httpapi.BadRequest(c, utils.StringsBuilder("Invalid parameters: ", err.Error()))
		case errors.Is(err, application.ErrCredentialNotFound):
			httpapi.Unauthorized(c, "Provider authentication required")
		case errors.Is(err, application.ErrRateLimitExceeded):
			httpapi.TooManyRequests(c, "Rate limit exceeded")
		case errors.Is(err, application.ErrAdapterExecuteFailed):
			httpapi.InternalServerError(c, utils.StringsBuilder("Failed to invoke operation: ", err.Error()))
		default:
			httpapi.InternalServerError(c, utils.StringsBuilder("Failed to invoke operation: ", err.Error()))
		}
		return
	}

	// Map to response
	response, err := mapInvocationToResponse(invocation, true)
	if err != nil {
		httpapi.InternalServerError(c, "Failed to format response")
		return
	}

	httpapi.OK(c, response, "Operation invoked successfully")
}
