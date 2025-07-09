package rest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
	"github.com/context-space/context-space/backend/internal/shared/utils"
)

// RESTConfig contains configuration specific to REST adapters
type RESTConfig struct {
	BaseURL     string
	Headers     map[string]string
	ContentType string
}

// RESTAdapter is an implementation of Adapter for REST APIs
type RESTAdapter struct {
	*base.BaseAdapter
	RestConfig *RESTConfig
	httpClient *http.Client
}

// NewRESTAdapter creates a new REST adapter
func NewRESTAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	restConfig *RESTConfig,
) *RESTAdapter {
	baseAdapter := base.NewBaseAdapter(
		providerInfo,
		config,
	)

	return &RESTAdapter{
		BaseAdapter: baseAdapter,
		RestConfig:  restConfig,
		httpClient: &http.Client{
			Timeout: config.Timeout,
			// TODO: Add transport for retries, circuit breaker, etc. later
		},
	}
}

// Execute executes a REST operation
func (a *RESTAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{},
	credential interface{},
) (interface{}, error) {
	// 1. Extract control parameters from params
	method, _ := params["method"].(string)
	pathTemplate, _ := params["path"].(string) // Get path template
	queryParams, _ := params["query_params"].(map[string]string)
	headers, _ := params["headers"].(map[string]string)
	body, bodyExists := params["body"]
	pathParams, pathParamsExist := params["path_params"].(map[string]string) // Get path params

	if method == "" {
		method = http.MethodGet // Default to GET
	}
	if pathTemplate == "" {
		return nil, fmt.Errorf("[%s] missing required parameter 'path' for REST execution", a.GetProviderAdapterInfo().Identifier)
	}

	// 1b. Substitute path parameters
	finalPath := pathTemplate
	if pathParamsExist {
		for key, value := range pathParams {
			placeholder := utils.StringsBuilder("{", key, "}")
			finalPath = strings.ReplaceAll(finalPath, placeholder, value)
		}
	}

	// 2. Construct the full URL
	baseURL, err := url.Parse(a.RestConfig.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("[%s] invalid base URL '%s': %w", a.GetProviderAdapterInfo().Identifier, a.RestConfig.BaseURL, err)
	}
	// Use the finalPath after substitution
	fullURL := baseURL.JoinPath(finalPath) // Handles slashes correctly

	// If full_url is provided, use it instead of the base URL and path
	if paramFullURL, ok := params["full_url"].(string); ok {
		fullURL, err = url.Parse(paramFullURL)
		if err != nil {
			return nil, fmt.Errorf("[%s] invalid full URL '%s': %w", a.GetProviderAdapterInfo().Identifier, paramFullURL, err)
		}
	}

	// 3. Prepare query parameters
	query := fullURL.Query()
	for k, v := range queryParams {
		query.Set(k, v)
	}
	fullURL.RawQuery = query.Encode()

	// 4. Prepare request body (if any)
	var reqBodyReader io.Reader = nil
	if bodyExists && body != nil {
		// Check if body is already a string (e.g., form-encoded data)
		if bodyStr, ok := body.(string); ok {
			// Body is already a string, use it directly
			reqBodyReader = strings.NewReader(bodyStr)
			// Set Content-Type for form data if not already specified
			if _, ok := headers["Content-Type"]; !ok {
				headers["Content-Type"] = "application/x-www-form-urlencoded"
			}
		} else {
			// Body is an object, serialize to JSON
			jsonData, err := sonic.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("[%s] failed to marshal request body for '%s': %w", a.GetProviderAdapterInfo().Identifier, operationID, err)
			}
			reqBodyReader = bytes.NewBuffer(jsonData)
			// Ensure Content-Type is set if body exists and not already specified
			if _, ok := headers["Content-Type"]; !ok {
				if a.RestConfig.ContentType != "" {
					headers["Content-Type"] = a.RestConfig.ContentType
				} else {
					headers["Content-Type"] = "application/json" // Default JSON
				}
			}
		}
	}

	// 5. Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), reqBodyReader)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to create request for '%s': %w", a.GetProviderAdapterInfo().Identifier, operationID, err)
	}

	// 6. Set headers (Defaults first, then specific ones)
	for k, v := range a.RestConfig.Headers {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 7. Send the request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		// TODO: Wrap with domain.AdapterError? How to get error code?
		return nil, fmt.Errorf("[%s] failed to execute request for '%s': %w", a.GetProviderAdapterInfo().Identifier, operationID, err)
	}
	defer resp.Body.Close()

	// 8. Read response body
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to read response body for '%s': %w", a.GetProviderAdapterInfo().Identifier, operationID, err)
	}

	// 9. Handle non-success status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to parse as JSON error, otherwise return raw body
		var errorResponse map[string]interface{}
		jsonErr := sonic.Unmarshal(respBodyBytes, &errorResponse)
		errorMessage := string(respBodyBytes)
		errorCode := "" // Placeholder for specific error code extraction if possible

		// Example: try to extract a common error message field
		if jsonErr == nil {
			if msg, ok := errorResponse["message"].(string); ok {
				errorMessage = msg
			}
			if code, ok := errorResponse["code"].(string); ok {
				errorCode = code
			}
		}

		finalErrorCode := errorCode
		if finalErrorCode == "" {
			// Map common HTTP statuses to domain errors, otherwise use a general one
			switch {
			case resp.StatusCode >= 500:
				finalErrorCode = domain.ErrProviderAPIError // Or a more specific server error if known
			case resp.StatusCode == 400:
				finalErrorCode = domain.ErrInvalidParameters
			case resp.StatusCode == 401 || resp.StatusCode == 403:
				finalErrorCode = domain.ErrCredentialError
			case resp.StatusCode == 404:
				finalErrorCode = domain.ErrOperationNotSupported // Or perhaps a more specific "not found"
			case resp.StatusCode == 429:
				finalErrorCode = domain.ErrProviderAPIError // Or a specific rate limit error if defined
			default:
				// For other 4xx errors or unhandled cases, stick to a general provider error
				finalErrorCode = domain.ErrProviderAPIError
			}
		}

		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			finalErrorCode, // Use the potentially mapped error code
			fmt.Sprintf("HTTP %d: %s", resp.StatusCode, errorMessage),
			resp.StatusCode,
		)
	}

	// 10. Parse successful response body (assume JSON)
	var result interface{}
	if len(respBodyBytes) > 0 {
		if err := sonic.Unmarshal(respBodyBytes, &result); err != nil {
			// If not JSON, return as raw string? Or error? Depends on expected content type.
			// For now, return error if JSON parsing fails for non-empty body
			return nil, fmt.Errorf("[%s] failed to unmarshal JSON response for '%s': %w (body: %s)", a.GetProviderAdapterInfo().Identifier, operationID, err, string(respBodyBytes))
		}
	} else {
		result = nil // No content response
	}

	return result, nil
}
