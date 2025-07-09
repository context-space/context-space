package stripe

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/context-space/context-space/backend/internal/shared/utils"
)

// Define constants for API paths used by handlers.
const (
	endpointCreatePaymentIntent        = "/payment_intents"
	endpointRetrievePaymentIntent      = "/payment_intents/{payment_intent_id}"
	endpointConfirmPaymentIntent       = "/payment_intents/{payment_intent_id}/confirm"
	endpointCreateCustomer             = "/customers"
	endpointListCustomerPaymentMethods = "/payment_methods"
	endpointDeletePaymentMethod        = "/payment_methods/{payment_method_id}"
	endpointCreateSubscription         = "/subscriptions"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDCreatePaymentIntent        = "create_payment_intent"
	operationIDRetrievePaymentIntent      = "retrieve_payment_intent"
	operationIDConfirmPaymentIntent       = "confirm_payment_intent"
	operationIDCreateCustomer             = "create_customer"
	operationIDListCustomerPaymentMethods = "list_customer_payment_methods"
	operationIDDeletePaymentMethod        = "delete_payment_method"
	operationIDCreateSubscription         = "create_subscription"
)

// OperationHandler defines the function signature for handling a specific API operation.
type OperationHandler func(ctx context.Context, params interface{}) (map[string]interface{}, error)

// OperationDefinition combines parameter schema and handler.
type OperationDefinition struct {
	Schema                interface{}      // Parameter schema (struct pointer)
	Handler               OperationHandler // Operation handler function
	PermissionIdentifiers []string         // List of internal permission identifiers (conceptual for API key)
}

// Operations maps operation IDs to their definitions.
type Operations map[string]OperationDefinition

// RegisterOperation registers the parameter schema and handler.
// This is called by the auto-generated registerOperations in stripe_operations.go
func (a *StripeAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler, requiredPerms []string) {
	a.BaseAdapter.RegisterOperation(operationID, schema) // Register schema for validation with BaseAdapter
	if a.operations == nil {
		a.operations = make(Operations)
	}
	a.operations[operationID] = OperationDefinition{
		Schema:                schema,
		Handler:               handler,
		PermissionIdentifiers: requiredPerms, // Store for informational purposes or future use
	}
}

// Define a struct for each operation's parameters based on config.

// CreatePaymentIntentParams defines parameters for the Create Payment Intent operation.
type CreatePaymentIntentParams struct {
	// Define fields based on create_payment_intent parameters in config

	Amount             int           `mapstructure:"amount" validate:"required"`               // Payment amount in cents.
	Currency           string        `mapstructure:"currency" validate:"required"`             // Currency type (e.g., usd, eur).
	PaymentMethodTypes []interface{} `mapstructure:"payment_method_types" validate:"required"` // Payment method types (e.g., card).
	Customer           string        `mapstructure:"customer" validate:"omitempty"`            // Customer ID.

}

// RetrievePaymentIntentParams defines parameters for the Retrieve Payment Intent operation.
type RetrievePaymentIntentParams struct {
	// Define fields based on retrieve_payment_intent parameters in config

	PaymentIntentId string `mapstructure:"payment_intent_id" validate:"required"` // Payment Intent ID.

}

// ConfirmPaymentIntentParams defines parameters for the Confirm Payment Intent operation.
type ConfirmPaymentIntentParams struct {
	// Define fields based on confirm_payment_intent parameters in config

	PaymentIntentId string `mapstructure:"payment_intent_id" validate:"required"` // Payment Intent ID.
	PaymentMethod   string `mapstructure:"payment_method" validate:"omitempty"`   // Payment method ID.

}

// CreateCustomerParams defines parameters for the Create Customer operation.
type CreateCustomerParams struct {
	// Define fields based on create_customer parameters in config

	Email         string `mapstructure:"email" validate:"omitempty"`          // Customer's email address.
	Name          string `mapstructure:"name" validate:"omitempty"`           // Customer's name.
	PaymentMethod string `mapstructure:"payment_method" validate:"omitempty"` // Preset payment method ID.

}

// ListCustomerPaymentMethodsParams defines parameters for the List Customer Payment Methods operation.
type ListCustomerPaymentMethodsParams struct {
	// Define fields based on list_customer_payment_methods parameters in config

	Customer string `mapstructure:"customer" validate:"required"` // Customer ID.
	Type     string `mapstructure:"type" validate:"required"`     // Payment method type (e.g., card).
	Limit    int    `mapstructure:"limit" validate:"omitempty"`   // Limit the number of payment methods returned.

}

// DeletePaymentMethodParams defines parameters for the Delete Payment Method operation.
type DeletePaymentMethodParams struct {
	// Define fields based on delete_payment_method parameters in config

	PaymentMethodId string `mapstructure:"payment_method_id" validate:"required"` // Payment Method ID.

}

// CreateSubscriptionParams defines parameters for the Create Subscription operation.
type CreateSubscriptionParams struct {
	// Define fields based on create_subscription parameters in config

	Customer string        `mapstructure:"customer" validate:"required"` // Customer ID.
	Items    []interface{} `mapstructure:"items" validate:"required"`    // Subscription items, including service plan ID.

}

// registerOperations is called by the adapter constructor to register all supported operations.
func (a *StripeAdapter) registerOperations() {

	a.RegisterOperation(
		operationIDCreatePaymentIntent,
		&CreatePaymentIntentParams{},
		handleCreatePaymentIntent,
		[]string{"write_payments"},
	)

	a.RegisterOperation(
		operationIDRetrievePaymentIntent,
		&RetrievePaymentIntentParams{},
		handleRetrievePaymentIntent,
		[]string{"read_payments"},
	)

	a.RegisterOperation(
		operationIDConfirmPaymentIntent,
		&ConfirmPaymentIntentParams{},
		handleConfirmPaymentIntent,
		[]string{"write_payments"},
	)

	a.RegisterOperation(
		operationIDCreateCustomer,
		&CreateCustomerParams{},
		handleCreateCustomer,
		[]string{"write_customers"},
	)

	a.RegisterOperation(
		operationIDListCustomerPaymentMethods,
		&ListCustomerPaymentMethodsParams{},
		handleListCustomerPaymentMethods,
		[]string{"read_customers"},
	)

	a.RegisterOperation(
		operationIDDeletePaymentMethod,
		&DeletePaymentMethodParams{},
		handleDeletePaymentMethod,
		[]string{"write_customers"},
	)

	a.RegisterOperation(
		operationIDCreateSubscription,
		&CreateSubscriptionParams{},
		handleCreateSubscription,
		[]string{"write_subscriptions"},
	)
}

// Define a handler function for each operation.

// handleCreatePaymentIntent constructs parameters for the REST adapter for the Create Payment Intent operation.
func handleCreatePaymentIntent(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*CreatePaymentIntentParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for operation create_payment_intent")
	}

	headers := make(map[string]string)
	form := url.Values{}

	form.Set("amount", fmt.Sprintf("%d", p.Amount))
	form.Set("currency", p.Currency)
	if p.Customer != "" {
		form.Set("customer", p.Customer)
	}
	for _, pmType := range p.PaymentMethodTypes {
		if typeStr, ok := pmType.(string); ok {
			form.Add("payment_method_types[]", typeStr)
		} else {
			return nil, fmt.Errorf("invalid type for payment_method_type item: expected string, got %T", pmType)
		}
	}

	headers["Content-Type"] = "application/x-www-form-urlencoded"
	requestBodyString := form.Encode()

	restParams := map[string]interface{}{
		"method":  http.MethodPost,
		"path":    endpointCreatePaymentIntent,
		"headers": headers,
		"body":    requestBodyString,
	}
	return restParams, nil
}

// handleRetrievePaymentIntent constructs parameters for the REST adapter for the Retrieve Payment Intent operation.
func handleRetrievePaymentIntent(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*RetrievePaymentIntentParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for operation retrieve_payment_intent")
	}

	pathParams := make(map[string]string)
	pathParams["payment_intent_id"] = p.PaymentIntentId

	restParams := map[string]interface{}{
		"method":      http.MethodGet,
		"path":        endpointRetrievePaymentIntent,
		"path_params": pathParams,
	}
	return restParams, nil
}

// handleConfirmPaymentIntent constructs parameters for the REST adapter for the Confirm Payment Intent operation.
func handleConfirmPaymentIntent(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ConfirmPaymentIntentParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for operation confirm_payment_intent")
	}

	pathParams := make(map[string]string)
	pathParams["payment_intent_id"] = p.PaymentIntentId

	headers := make(map[string]string)
	form := url.Values{}

	if p.PaymentMethod != "" {
		form.Set("payment_method", p.PaymentMethod)
	}

	headers["Content-Type"] = "application/x-www-form-urlencoded"
	requestBodyString := form.Encode()

	restParams := map[string]interface{}{
		"method":      http.MethodPost,
		"path":        endpointConfirmPaymentIntent,
		"path_params": pathParams,
		"headers":     headers,
		"body":        requestBodyString,
	}
	return restParams, nil
}

// handleCreateCustomer constructs parameters for the REST adapter for the Create Customer operation.
func handleCreateCustomer(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*CreateCustomerParams)
	if !ok && params != nil {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation create_customer", params)
	}

	// Build form data for Stripe API
	form := url.Values{}

	if p != nil {
		if p.Email != "" {
			form.Set("email", p.Email)
		}
		if p.Name != "" {
			form.Set("name", p.Name)
		}
		if p.PaymentMethod != "" {
			form.Set("payment_method", p.PaymentMethod)
		}
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointCreateCustomer,
		"body":   form.Encode(),
	}
	return restParams, nil
}

// handleListCustomerPaymentMethods constructs parameters for the REST adapter for the List Customer Payment Methods operation.
func handleListCustomerPaymentMethods(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ListCustomerPaymentMethodsParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for operation list_customer_payment_methods")
	}

	queryParams := make(map[string]string)
	queryParams["customer"] = p.Customer
	queryParams["type"] = p.Type
	if p.Limit != 0 {
		queryParams["limit"] = fmt.Sprintf("%d", p.Limit)
	}

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointListCustomerPaymentMethods,
		"query_params": queryParams,
	}
	return restParams, nil
}

// handleDeletePaymentMethod constructs parameters for the REST adapter for the Delete Payment Method operation.
func handleDeletePaymentMethod(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*DeletePaymentMethodParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for operation delete_payment_method")
	}

	pathParams := make(map[string]string)
	pathParams["payment_method_id"] = p.PaymentMethodId

	restParams := map[string]interface{}{
		"method":      http.MethodDelete,
		"path":        endpointDeletePaymentMethod,
		"path_params": pathParams,
	}
	return restParams, nil
}

// handleCreateSubscription constructs parameters for the REST adapter for the Create Subscription operation.
func handleCreateSubscription(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*CreateSubscriptionParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for operation create_subscription")
	}

	headers := make(map[string]string)
	form := url.Values{}
	form.Set("customer", p.Customer)

	var itemsParams []string
	for i, itemInterface := range p.Items {
		itemMap, ok := itemInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type for item in 'items' array at index %d: expected map[string]interface{}, got %T", i, itemInterface)
		}
		priceInterface, ok := itemMap["price"]
		if !ok {
			return nil, fmt.Errorf("item in 'items' array at index %d is missing 'price' field", i)
		}
		priceStr, ok := priceInterface.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type for 'price' in item at index %d: expected string, got %T", i, priceInterface)
		}
		itemsParams = append(itemsParams, fmt.Sprintf("items[%d][price]=%s", i, url.QueryEscape(priceStr)))
	}

	requestBodyString := form.Encode()
	if len(itemsParams) > 0 {
		if requestBodyString != "" {
			requestBodyString = utils.StringsBuilder(requestBodyString, "&")
		}
		requestBodyString = utils.StringsBuilder(requestBodyString, strings.Join(itemsParams, "&"))
	}

	headers["Content-Type"] = "application/x-www-form-urlencoded"

	restParams := map[string]interface{}{
		"method":  http.MethodPost,
		"path":    endpointCreateSubscription,
		"headers": headers,
		"body":    requestBodyString,
	}
	return restParams, nil
}
