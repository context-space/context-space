package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

var skipPaths = []string{"/health"}

const (
	// 请求体记录的最大大小限制（64KB）
	maxRequestBodySize = 64 * 1024
)

// 支持记录的Content-Type列表
var loggableContentTypes = []string{
	"application/json",
	"application/xml",
	"text/xml",
	"text/plain",
	"application/x-www-form-urlencoded",
}

// RequestLoggingMiddleware logs all HTTP requests
func RequestLoggingMiddleware(obs *observability.ObservabilityProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if this path should be skipped
		for _, path := range skipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		start := time.Now()
		requestID := uuid.New().String()

		// Add request ID to context and headers
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Read request body for logging
		requestBody := readRequestBody(c)

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)
		responseBody := blw.body.String()

		// Build log fields
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("duration_ms", duration.Milliseconds()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("size", c.Writer.Size()),
			zap.String("request_body", requestBody),
			zap.String("response_body", responseBody),
		}

		// Add error information
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.String()))
		}

		// Choose log level based on status code
		if c.Writer.Status() >= 400 {
			obs.Logger.Error(c.Request.Context(), "HTTP Request Error", fields...)
		} else {
			obs.Logger.Info(c.Request.Context(), "HTTP Request", fields...)
		}
	}
}

// isLoggableContentType 检查Content-Type是否可记录
func isLoggableContentType(contentType string) bool {
	for _, loggableType := range loggableContentTypes {
		if strings.Contains(strings.ToLower(contentType), loggableType) {
			return true
		}
	}
	return false
}

// readRequestBody 安全地读取并重建请求体
func readRequestBody(c *gin.Context) string {
	// 检查Content-Type是否可记录
	contentType := c.GetHeader("Content-Type")
	if !isLoggableContentType(contentType) {
		return ""
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}

	// 重建请求体供后续handler使用
	c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

	// 限制记录的body大小
	if len(body) > maxRequestBodySize {
		return string(body[:maxRequestBodySize]) + "...truncated"
	}

	return string(body)
}

// TracingMiddleware adds distributed tracing to all requests
func TracingMiddleware(obs *observability.ObservabilityProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract trace context from headers
		ctx := c.Request.Context()
		carrier := propagation.HeaderCarrier(c.Request.Header)
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		// RefreshToken a new span for this request
		operationName := utils.StringsBuilder(c.Request.Method, " ", c.FullPath())
		if operationName == " " { // For 404 routes
			operationName = utils.StringsBuilder(c.Request.Method, " ", c.Request.URL.Path)
		}

		ctx, span := obs.Tracer.Start(ctx, operationName)
		defer span.End()

		// Add request details to span
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.client_ip", c.ClientIP()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
			attribute.String("http.request_id", c.GetString("request_id")),
		)

		// Store context in request
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Add response details to span
		statusCode := c.Writer.Status()
		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.Int("http.response_size", c.Writer.Size()),
		)

		// Set span status based on HTTP status code
		if statusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	}
}

// MetricsMiddleware collects metrics for all HTTP requests
func MetricsMiddleware(obs *observability.ObservabilityProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := c.Request.Context()

		// Get clean route pattern (e.g., /users/:id instead of /users/123)
		path := c.FullPath()
		method := c.Request.Method

		// Increment active requests counter at the start
		obs.Metrics.IncrementCounter(ctx, "http_requests_active", 1,
			attribute.String("method", method),
			attribute.String("path", path),
		)

		// Process request
		c.Next()

		// Get status code after request is processed
		statusCode := c.Writer.Status()

		// Record request duration
		duration := time.Since(start).Seconds()
		obs.Metrics.RecordHistogram(ctx, "http_request_duration_seconds", duration,
			attribute.String("method", method),
			attribute.String("path", path),
			attribute.Int("status", statusCode),
		)

		// Increment request counter
		obs.Metrics.IncrementCounter(ctx, "http_requests_total", 1,
			attribute.String("method", method),
			attribute.String("path", path),
			attribute.Int("status", statusCode),
		)

		// Decrement active requests counter at the end
		obs.Metrics.IncrementCounter(ctx, "http_requests_active", -1,
			attribute.String("method", method),
			attribute.String("path", path),
		)
	}
}

// RegisterObservabilityMiddleware registers all observability middleware with the Gin router
func RegisterObservabilityMiddleware(router *gin.Engine, obs *observability.ObservabilityProvider) {
	// Apply middleware in the correct order:
	// 1. Request logging (to capture timing accurately)
	// 2. Tracing (to set up spans)
	// 3. Metrics (to count and time requests)
	router.Use(
		RequestLoggingMiddleware(obs),
		TracingMiddleware(obs),
		MetricsMiddleware(obs),
	)
}

// bodyLogWriter is a custom ResponseWriter used to capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
