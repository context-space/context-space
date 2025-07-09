package apierrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorResponse is the standard error response format for the API
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// HandleError handles an error and returns an appropriate response
func HandleError(c *gin.Context, err error, message string, statusCode ...int) {
	// Default to internal server error
	code := http.StatusInternalServerError

	// If status code is provided, use it
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	// Check if the error is an APIError
	var apiErr *APIError
	if apiError, ok := err.(*APIError); ok {
		apiErr = apiError
		code = apiError.HTTPCode
	} else {
		// Create a new internal error
		apiErr = NewInternalError(message, err)
	}

	// Log the error
	logger, exists := c.Get("logger")
	if exists {
		if zapLogger, ok := logger.(*zap.Logger); ok {
			zapLogger.Error("API Error",
				zap.String("type", string(apiErr.Type)),
				zap.String("code", apiErr.Code),
				zap.String("message", apiErr.Message),
				zap.Any("details", apiErr.Details),
				zap.Error(err),
			)
		}
	}

	// Create the response
	response := ErrorResponse{
		Error: *apiErr,
	}

	// Return the error response
	c.JSON(code, response)
}
