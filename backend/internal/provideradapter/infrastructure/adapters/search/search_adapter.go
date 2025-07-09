package search

import (
	"context"
	"fmt"
	"net/http"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
)

// SearchAdapter is the adapter for Serper API
type SearchAdapter struct {
	*base.BaseAdapter
	restAdapter domain.Adapter // Use RESTAdapter for actual execution
	operations  Operations     // Operation ID mapping to OperationDefinition
	apiKey      string
}

// NewSearchAdapter creates a new SearchAdapter
func NewSearchAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	restAdapter domain.Adapter,
	apiKey string,
) *SearchAdapter {
	// 1. Create SearchAdapter instance
	baseAdapter := base.NewBaseAdapter(providerInfo, config)
	adapter := &SearchAdapter{
		BaseAdapter: baseAdapter,
		restAdapter: restAdapter,
		operations:  make(Operations),
		apiKey:      apiKey,
	}

	// Register specific Serper API operations
	adapter.registerOperations()

	return adapter
}

// Execute finds the appropriate handler for the operation, prepares parameters including API key,
// and delegates execution to the underlying RESTAdapter.
func (a *SearchAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{}, // Original user parameters
	credential interface{}, // This parameter is now possibly no longer needed, or used for other purposes
) (interface{}, error) {

	apiKey := a.apiKey
	if apiKey == "" {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrCredentialError, "invalid or missing API key credential", http.StatusUnauthorized)
	}

	// 1. Find the operation handler definition
	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "", fmt.Sprintf("Unknown operation ID: %s", operationID), http.StatusNotFound)
	}

	// 2. Process parameters using registered schema for decoding and validation
	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "PARAMETER_ERROR", fmt.Sprintf("Parameter validation failed: %s", err.Error()), http.StatusBadRequest)
	}

	// 3. Call handler with processed parameters
	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "HANDLER_ERROR", fmt.Sprintf("Handler execution failed: %s", err.Error()), http.StatusInternalServerError)
	}

	// 4. Inject API Key into restParams
	headers, _ := restParams["headers"].(map[string]string)
	if headers == nil {
		headers = make(map[string]string)
	}
	headers[apikeyParamName] = apiKey
	restParams["headers"] = headers

	// 5. Delegate execution to REST adapter
	result, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
