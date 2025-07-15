package fetch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
)

// FetchAdapter is an adapter for web content fetching operations
type FetchAdapter struct {
	*base.BaseAdapter
	restAdapter domain.Adapter     // Uses RESTAdapter for actual execution
	operations  Operations         // Map operation ID to OperationDefinition
	defaults    *OperationDefaults // Operation default values
	apiKey      string             // API key for authentication
}

// NewFetchAdapter creates a new FetchAdapter
func NewFetchAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	restAdapter *rest.RESTAdapter,
	defaults *OperationDefaults,
	apiKey string,
) *FetchAdapter {

	// Create FetchAdapter instance
	baseAdapter := base.NewBaseAdapter(providerInfo, config)
	adapter := &FetchAdapter{
		BaseAdapter: baseAdapter,
		restAdapter: restAdapter,
		operations:  make(Operations),
		defaults:    defaults,
		apiKey:      apiKey,
	}

	// Register specific fetch operations
	adapter.registerOperations()

	return adapter
}

// Execute finds the appropriate handler for the operationID, prepares parameters
// including the API key, and delegates execution to the underlying RESTAdapter.
func (a *FetchAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{}, // Original user parameters
	credential interface{}, // This parameter is not used, authentication is handled via apiKey
) (interface{}, error) {

	apiKey := a.apiKey
	if apiKey == "" {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrCredentialError, "invalid or missing API key credential", http.StatusUnauthorized)
	}

	// Find the operation handler definition
	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "", fmt.Sprintf("unknown operation ID: %s", operationID), http.StatusNotFound)
	}

	// Process and validate parameters
	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "PARAMETER_ERROR", fmt.Sprintf("parameter validation failed: %s", err.Error()), http.StatusBadRequest)
	}

	restParams, err := opDef.Handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "HANDLER_ERROR", fmt.Sprintf("handler execution failed: %s", err.Error()), http.StatusInternalServerError)
	}

	// Inject API Key into restParams based on configuration
	headers, _ := restParams["headers"].(map[string]string)
	if headers == nil {
		headers = make(map[string]string)
	}
	headers[apikeyParamName] = apiKey
	restParams["headers"] = headers

	// Delegate execution to REST adapter
	// Pass nil credential to restAdapter since authentication is handled here
	result, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		// Error should already be wrapped by RESTAdapter or domain.AdapterError
		return nil, err
	}

	return result, nil
}
