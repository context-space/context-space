package openweathermap

import (
	"context"
	"fmt"
	"net/http"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
	"github.com/context-space/context-space/backend/internal/shared/utils"
)

// OpenWeatherMapAdapter implements OpenWeatherMap API adapter
type OpenWeatherMapAdapter struct {
	*base.BaseAdapter
	restAdapter domain.Adapter // Underlying REST adapter instance
	operations  Operations     // Mapping from operation ID to definition
}

// NewOpenWeatherMapAdapter creates a new adapter instance
func NewOpenWeatherMapAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	restAdapter domain.Adapter,
) *OpenWeatherMapAdapter {
	adapter := &OpenWeatherMapAdapter{
		BaseAdapter: base.NewBaseAdapter(providerInfo, config),
		restAdapter: restAdapter,
		operations:  make(Operations),
	}

	// Register all operations
	adapter.registerOperations()

	return adapter
}

// Execute executes the specified operation
func (a *OpenWeatherMapAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{},
	credential interface{},
) (interface{}, error) {
	// Validate API Key credential
	apiKeyCred, ok := credential.(*credDomain.APIKeyCredential)
	if !ok || apiKeyCred == nil || apiKeyCred.APIKey == "" {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrCredentialError,
			"invalid or missing API key credential",
			http.StatusUnauthorized,
		)
	}

	// Check if operation exists
	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrOperationNotSupported,
			fmt.Sprintf("unknown operation ID: %s", operationID),
			http.StatusNotFound,
		)
	}

	// Process and validate parameters
	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrInvalidParameters,
			fmt.Sprintf("parameter validation failed: %v", err),
			http.StatusBadRequest,
		)
	}

	// Call operation handler function
	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrInternal,
			fmt.Sprintf("operation handler failed: %v", err),
			http.StatusInternalServerError,
		)
	}

	// Add API Key to query parameters
	queryParams, _ := restParams["query_params"].(map[string]string)
	if queryParams == nil {
		queryParams = make(map[string]string)
	}
	queryParams["appid"] = apiKeyCred.APIKey
	restParams["query_params"] = queryParams

	// Set request headers
	headers, _ := restParams["headers"].(map[string]string)
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	headers["User-Agent"] = utils.StringsBuilder("OpenWeatherMap-Adapter/1.0")
	restParams["headers"] = headers

	// Call underlying REST adapter
	rawResult, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		return nil, a.handleRestError(operationID, err)
	}

	return rawResult, nil
}

// handleRestError handles REST adapter errors
func (a *OpenWeatherMapAdapter) handleRestError(operationID string, err error) error {
	if adapterErr, ok := err.(*domain.AdapterError); ok {
		// If already an adapter error, return directly
		return adapterErr
	}

	// Otherwise wrap as internal error
	return domain.NewAdapterError(
		a.GetProviderAdapterInfo().Identifier,
		operationID,
		domain.ErrInternal,
		fmt.Sprintf("REST adapter error: %v", err),
		http.StatusInternalServerError,
	)
}

// GetSupportedOperations returns list of supported operations
func (a *OpenWeatherMapAdapter) GetSupportedOperations() []string {
	operations := make([]string, 0, len(a.operations))
	for operationID := range a.operations {
		operations = append(operations, operationID)
	}
	return operations
}
