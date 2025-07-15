package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/shared/utils"

	volcenginetypes "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/knowledgebase/volcengine/types"

	"github.com/bytedance/sonic"
	volcBase "github.com/volcengine/volc-sdk-golang/base"
)

const (
	defaultAPIBaseURL = "https://api-knowledgebase.mlp.cn-beijing.volces.com"
	defaultService    = "air"
	defaultRegion     = "cn-north-1"
)

// VolcengineClient is an internal client responsible for making signed HTTP requests
// to the Volcengine API endpoints. It does not implement the domain.Adapter interface.
type VolcengineClient struct {
	service            string
	region             string
	apiBaseURL         string
	httpClient         *http.Client
	providerIdentifier string
	credentials        *volcenginetypes.VolcengineCredential
}

// NewVolcengineClient creates a new internal Volcengine client.
func NewVolcengineClient(
	providerIdentifier string, // ID of the provider using this client (for error reporting)
	apiBaseURL string, // Base URL for the API endpoint
	service string, // e.g., "air"
	region string, // e.g., "cn-north-1"
	timeout time.Duration, // HTTP client timeout
	creds *volcenginetypes.VolcengineCredential, // New parameter for storing credentials
) *VolcengineClient {
	if service == "" {
		service = defaultService
	}
	if region == "" {
		region = defaultRegion
	}
	if apiBaseURL == "" {
		apiBaseURL = defaultAPIBaseURL
	}

	client := &VolcengineClient{
		providerIdentifier: providerIdentifier,
		service:            service,
		region:             region,
		apiBaseURL:         strings.TrimRight(apiBaseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
		credentials: creds, // Store the credentials
	}
	return client
}

// Execute performs a signed Volcengine API request.
// It takes the necessary HTTP method, path, query parameters, request body,
// and credentials to perform the call.
func (c *VolcengineClient) Execute(
	ctx context.Context,
	operationID string, // Added operationID for more context in errors
	method string,
	path string,
	query url.Values,
	body interface{}, // Accepts struct to be marshaled
) (interface{}, error) {

	credToUse := c.credentials
	// Check if we have credentials to use
	if credToUse == nil {
		adapterErr := domain.NewAdapterError(
			c.providerIdentifier,
			operationID,
			"CREDENTIAL_ERROR",
			"no credentials provided and no credentials stored in client",
			http.StatusUnauthorized,
		)
		return nil, adapterErr
	}

	// 1. Serialize Request Body (if provided)
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = sonic.Marshal(body)
		if err != nil {
			adapterErr := domain.NewAdapterError(
				c.providerIdentifier,
				operationID,
				"INTERNAL_ERROR", // Or a more specific code?
				fmt.Sprintf("failed to marshal request body: %v", err),
				http.StatusInternalServerError, // Assuming this is an internal issue
			)
			adapterErr.Raw = err
			return nil, adapterErr
		}
	}

	// 2. Build and Sign the HTTP request using provided credentials
	fullURL := utils.StringsBuilder(c.apiBaseURL, path)
	// Pass credentials to the signing function
	signedReq, err := c.prepareAndSignRequest(ctx, method, fullURL, query, bodyBytes, credToUse)
	if err != nil {
		adapterErr := domain.NewAdapterError(
			c.providerIdentifier,
			operationID,
			"AUTHENTICATION_FAILED", // Specific code for signing issue
			fmt.Sprintf("failed to sign request: %v", err),
			http.StatusInternalServerError, // Treat signing issues as internal error for now
		)
		adapterErr.Raw = err
		return nil, adapterErr
	}

	// 3. Execute the request
	httpResp, err := c.httpClient.Do(signedReq.WithContext(ctx)) // Ensure context is passed
	if err != nil {
		// Network errors, timeouts, etc.
		adapterErr := domain.NewAdapterError(
			c.providerIdentifier,
			operationID,
			"NETWORK_ERROR",
			fmt.Sprintf("HTTP request failed: %v", err),
			http.StatusServiceUnavailable, // 503 Service Unavailable seems appropriate
		)
		adapterErr.Raw = err
		return nil, adapterErr
	}
	defer httpResp.Body.Close()

	// 4. Handle the response
	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		adapterErr := domain.NewAdapterError(
			c.providerIdentifier,
			operationID,
			"PROVIDER_ERROR", // Error reading response from provider
			fmt.Sprintf("failed to read response body: %v", err),
			httpResp.StatusCode, // Use actual status code
		)
		adapterErr.Raw = err
		return nil, adapterErr
	}

	// Check for non-2xx status codes first
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		// Attempt to parse Volcengine error response
		// parseVolcengineError returns *domain.AdapterError directly
		// Pass providerIdentifier and operationID for context
		apiErr := parseVolcengineError(c.providerIdentifier, operationID, respBodyBytes, httpResp.StatusCode)
		return nil, apiErr
	}

	// 5. Parse successful response (assuming JSON)
	var result interface{}
	if len(respBodyBytes) > 0 {
		// Unmarshal into a generic interface{} for now.
		// Future: Could accept a target struct pointer for unmarshaling.
		err = sonic.Unmarshal(respBodyBytes, &result)
		if err != nil {
			// If successful status code but invalid JSON response
			adapterErr := domain.NewAdapterError(
				c.providerIdentifier,
				operationID,
				"PROVIDER_ERROR", // Provider returned success status but bad body
				fmt.Sprintf("failed to parse successful JSON response: %v. Body: %s", err, string(respBodyBytes)),
				httpResp.StatusCode,
			)
			adapterErr.Raw = err
			return nil, adapterErr
		}
	} else {
		// Handle empty successful response if applicable
		result = nil
	}

	return result, nil
}

// prepareAndSignRequest is a helper to create and sign a Volcengine API request.
// It now accepts VolcengineCredential to get AK/SK and optional AccountID.
func (c *VolcengineClient) prepareAndSignRequest(
	ctx context.Context,
	method string,
	fullURL string,
	query url.Values,
	bodyBytes []byte,
	cred *volcenginetypes.VolcengineCredential, // Accept credentials (use type from types package)
) (*http.Request, error) {
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL '%s': %w", fullURL, err)
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	// Host is set automatically

	// Use credentials passed into the function
	volcCredential := volcBase.Credentials{
		AccessKeyID:     cred.AccessKeyID,
		SecretAccessKey: cred.SecretAccessKey,
		Service:         c.service, // Use service/region from client config
		Region:          c.region,
	}

	// Sign modifies the request in place
	signedReq := volcCredential.Sign(req)

	return signedReq, nil
}

// VolcengineAPIError represents the structure of a typical Volcengine error response.
type VolcengineAPIError struct {
	Code       int    `json:"code"`       // Specific Volcengine error code
	Message    string `json:"message"`    // Error message from Volcengine
	RequestID  string `json:"request_id"` // Request ID for tracing
	HTTPStatus int    // Added field to store HTTP status
	RawBody    string // Added field to store raw error body
}

func (e *VolcengineAPIError) Error() string {
	return fmt.Sprintf("Volcengine API Error: Code=%d, Message='%s', HTTPStatus=%d, RequestID='%s'", e.Code, e.Message, e.HTTPStatus, e.RequestID)
}

// IsAuthError checks if the error code suggests an authentication/authorization issue.
func (e *VolcengineAPIError) IsAuthError() bool {
	knownAuthCodes := map[int]bool{
		401: true, // Unauthorized (General HTTP)
		403: true, // Forbidden (General HTTP)
		// Add specific Volcengine codes if documented, e.g.:
		// 10001: true, // Invalid AccessKeyId
		// 10004: true, // SignatureDoesNotMatch
	}
	// Check internal code first, then HTTP status
	return knownAuthCodes[e.Code] || knownAuthCodes[e.HTTPStatus]
}

// parseVolcengineError attempts to parse a Volcengine error response body
// and returns a populated domain.AdapterError.
// Takes providerIdentifier and operationID for context.
func parseVolcengineError(providerID, operationID string, bodyBytes []byte, httpStatus int) *domain.AdapterError {
	var volcError VolcengineAPIError
	err := sonic.Unmarshal(bodyBytes, &volcError)

	rawBody := string(bodyBytes)
	parsedError := &volcError // Keep reference even if parsing fails
	parsedError.HTTPStatus = httpStatus
	parsedError.RawBody = rawBody

	if err != nil {
		// Failed to parse JSON error structure
		errorCode := "PROVIDER_ERROR"
		errorMessage := fmt.Sprintf("Volcengine request failed with status %d: %s", httpStatus, rawBody)
		// Refine error code based on HTTP status if possible
		if httpStatus >= 400 && httpStatus < 500 {
			if httpStatus == 401 || httpStatus == 403 {
				errorCode = "AUTHENTICATION_FAILED"
			} else {
				errorCode = "INVALID_INPUT" // Assume other 4xx are input errors
			}
		}

		adapterErr := domain.NewAdapterError(
			providerID,
			operationID,
			errorCode,
			errorMessage,
			httpStatus,
		)
		adapterErr.Raw = rawBody // Store raw body as Raw
		return adapterErr
	}

	// Successfully parsed the Volcengine error structure
	adapterErrorCode := "PROVIDER_ERROR" // Default for server errors
	if httpStatus >= 400 && httpStatus < 500 {
		if parsedError.IsAuthError() { // Check parsed code and status
			adapterErrorCode = "AUTHENTICATION_FAILED"
		} else {
			adapterErrorCode = "INVALID_INPUT" // Assume other 4xx are input errors
		}
	}

	adapterErr := domain.NewAdapterError(
		providerID,
		operationID,
		adapterErrorCode,
		volcError.Message, // Use parsed message from API
		httpStatus,
	)
	// Store the parsed VolcengineAPIError struct as Raw for more details
	adapterErr.Raw = parsedError
	return adapterErr
}
