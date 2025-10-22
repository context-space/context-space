package serper

import (
	"context"
	"fmt"
	"net/http"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
)

// SerperAdapter is an adapter for the Serper API
type SerperAdapter struct {
	*base.BaseAdapter
	restAdapter domain.Adapter // Uses RESTAdapter for actual execution
	permissions domain.PermissionSet
	operations  Operations // Map operation ID to OperationDefinition
}

// NewSerperAdapter creates a new SerperAdapter
func NewSerperAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	restAdapter domain.Adapter,
	permissions domain.PermissionSet,
) *SerperAdapter {

	baseAdapter := base.NewBaseAdapter(providerInfo, config)
	adapter := &SerperAdapter{
		BaseAdapter: baseAdapter,
		restAdapter: restAdapter,
		permissions: permissions,
		operations:  make(Operations),
	}

	adapter.registerOperations()

	return adapter
}

// Execute finds the appropriate handler for the operationID, prepares parameters
// including the API key, and delegates execution to the underlying RESTAdapter.
func (a *SerperAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{}, // Raw user parameters
	credential interface{},
) (interface{}, error) {
	cred, ok := credential.(*credDomain.APIKeyCredential)
	if !ok || cred == nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "", "invalid or missing API key credential", http.StatusUnauthorized)
	}
	apiKey := cred.APIKey
	if apiKey == "" {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "", "invalid or missing API key", http.StatusUnauthorized)
	}

	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "", fmt.Sprintf("unknown operation ID: %s", operationID), http.StatusNotFound)
	}

	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "PARAMETER_ERROR", fmt.Sprintf("parameter validation failed: %s", err.Error()), http.StatusBadRequest)
	}

	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, "HANDLER_ERROR", fmt.Sprintf("handler execution failed: %s", err.Error()), http.StatusInternalServerError)
	}

	headers, _ := restParams["headers"].(map[string]string)
	if headers == nil {
		headers = make(map[string]string)
	}
	headers[apikeyParamName] = apiKey
	restParams["headers"] = headers

	// For Serper, the 'scrape' operation has a full URL in its endpoint_path.
	// The restAdapter needs to know if the path is absolute.
	if operationID == operationIDScrape { // Assuming OperationIDScrape will be defined
		path, ok := restParams["path"].(string)
		if ok && (path == endpointScrape || path == "https://scrape.serper.dev") { // endpointScrape to be defined
			// If the path is absolute, we need to pass the full URL to the restAdapter
			restParams["full_url"] = path
		}
	}

	result, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
