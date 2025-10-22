package stripe

import (
	"context"
	"fmt"
	"net/http"

	"encoding/base64"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
	"github.com/context-space/context-space/backend/internal/shared/utils"
)

// StripeAdapter is an adapter for the Stripe API using API Key.
type StripeAdapter struct {
	*base.BaseAdapter
	restAdapter   domain.Adapter       // The underlying REST adapter instance
	operations    Operations           // Map of operation ID to definition defined in _operations.go
	permissionSet domain.PermissionSet // Permission set defined in providercore (used conceptually)
}

// NewStripeAdapter creates a new Stripe adapter.
func NewStripeAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	restAdapter domain.Adapter,
	permissions domain.PermissionSet,
) *StripeAdapter {
	baseAdapter := base.NewBaseAdapter(providerInfo, config)

	adapter := &StripeAdapter{
		BaseAdapter:   baseAdapter,
		restAdapter:   restAdapter,
		operations:    make(Operations),
		permissionSet: permissions,
	}

	adapter.registerOperations()

	return adapter
}

// Execute handles API calls based on the operationID using the REST adapter.
func (a *StripeAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{}, // User-provided parameters
	credential interface{}, // Expected to be *credDomain.APIKeyCredential
) (interface{}, error) {
	apiKeyCred, ok := credential.(*credDomain.APIKeyCredential)
	if !ok || apiKeyCred == nil || apiKeyCred.APIKey == "" {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrCredentialError, "invalid or missing API key credential", http.StatusUnauthorized)
	}

	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrOperationNotSupported, fmt.Sprintf("unknown operation ID: %s", operationID), http.StatusNotFound)
	}

	// Permission checks for API key based auth are usually handled by the provider based on the key's privileges.
	// The 'permissions' defined in stripe.json are more for user understanding and UI.
	// If specific pre-flight checks were needed based on `opDef.PermissionIdentifiers`, they could be added here.

	processedParams, err := a.ProcessParams(operationID, params) // From BaseAdapter
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInvalidParameters, fmt.Sprintf("parameter validation failed: %v", err), http.StatusBadRequest)
	}

	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInternal, fmt.Sprintf("operation handler failed: %v", err), http.StatusInternalServerError)
	}

	headers, _ := restParams["headers"].(map[string]string)
	if headers == nil {
		headers = make(map[string]string)
	}

	// Stripe uses HTTP Basic Auth with API key as username and empty password
	// Format: "Basic " + base64("apikey:")
	authString := utils.StringsBuilder(apiKeyCred.APIKey, ":")
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))
	headers[apikeyParamName] = utils.StringsBuilder("Basic ", encodedAuth)

	restParams["headers"] = headers

	rawResult, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		return nil, err
	}

	return rawResult, nil
}
