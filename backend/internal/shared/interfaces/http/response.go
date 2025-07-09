package http

import (
	"net/http"

	"github.com/context-space/context-space/backend/internal/shared/apierrors"
	"github.com/gin-gonic/gin"
)

// Response represents a standardized API response format
type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SwaggerErrorResponse represents an error API response for Swagger docs
type SwaggerErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Error message"`
}

// SuccessResponse creates a success response with the given data and optional message
func SuccessResponse(code int, message string, data interface{}) Response {
	return Response{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response with the given message
func ErrorResponse(code int, message string) Response {
	return Response{
		Success: false,
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// RespondWithSuccess sends a success response with the given data and HTTP status code
func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(http.StatusOK, SuccessResponse(statusCode, message, data))
}

// RespondWithError sends an error response with the given message and HTTP status code
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(http.StatusOK, ErrorResponse(statusCode, message))
}

// RespondWithAPIError handles an APIError and sends the appropriate response
func RespondWithAPIError(c *gin.Context, err *apierrors.APIError) {
	c.JSON(err.HTTPCode, ErrorResponse(err.HTTPCode, err.Message))
}

// OK sends a 200 OK response with the given data
func OK(c *gin.Context, data interface{}, message string) {
	RespondWithSuccess(c, http.StatusOK, data, message)
}

// Created sends a 201 Created response with the given data
func Created(c *gin.Context, data interface{}, message string) {
	RespondWithSuccess(c, http.StatusCreated, data, message)
}

// NoContent sends a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a 400 Bad Request response with the given message
func BadRequest(c *gin.Context, message string) {
	RespondWithError(c, http.StatusBadRequest, message)
}

// Unauthorized sends a 401 Unauthorized response with the given message
func Unauthorized(c *gin.Context, message string) {
	RespondWithError(c, http.StatusUnauthorized, message)
}

// Forbidden sends a 403 Forbidden response with the given message
func Forbidden(c *gin.Context, message string) {
	RespondWithError(c, http.StatusForbidden, message)
}

// NotFound sends a 404 Not Found response with the given message
func NotFound(c *gin.Context, message string) {
	RespondWithError(c, http.StatusNotFound, message)
}

// TooManyRequests sends a 429 Too Many Requests response with the given message
func TooManyRequests(c *gin.Context, message string) {
	RespondWithError(c, http.StatusTooManyRequests, message)
}

// InternalServerError sends a 500 Internal Server Error response with the given message
func InternalServerError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusInternalServerError, message)
}
