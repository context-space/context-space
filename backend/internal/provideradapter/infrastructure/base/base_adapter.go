package base

import (
	"context"
	"fmt"
	"reflect"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

// OperationDefinition describes an operation with its parameter schema
type OperationDefinition struct {
	OperationID string
	Schema      interface{} // Struct pointer that defines the parameters
}

// BaseAdapter provides common functionality for all adapters
type BaseAdapter struct {
	ProviderAdapterInfo *domain.ProviderAdapterInfo
	Config              *domain.AdapterConfig
	Operations          map[string]OperationDefinition
	validate            *validator.Validate
}

// NewBaseAdapter creates a new base adapter
func NewBaseAdapter(
	providerAdapterInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
) *BaseAdapter {
	return &BaseAdapter{
		ProviderAdapterInfo: providerAdapterInfo,
		Config:              config,
		Operations:          make(map[string]OperationDefinition),
		validate:            validator.New(),
	}
}

// GetProviderAdapterInfo returns information about this provider
func (a *BaseAdapter) GetProviderAdapterInfo() *domain.ProviderAdapterInfo {
	return &domain.ProviderAdapterInfo{
		Identifier:  a.ProviderAdapterInfo.Identifier,
		Name:        a.ProviderAdapterInfo.Name,
		Description: a.ProviderAdapterInfo.Description,
		AuthType:    a.ProviderAdapterInfo.AuthType,
	}
}

// Execute is a placeholder that should be overridden by concrete adapters
func (a *BaseAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{},
	credential interface{},
) (interface{}, error) {
	return nil, fmt.Errorf("Execute not implemented in base adapter")
}

// RegisterOperation registers an operation with its parameter schema
func (a *BaseAdapter) RegisterOperation(operationID string, schema interface{}) {
	a.Operations[operationID] = OperationDefinition{
		OperationID: operationID,
		Schema:      schema,
	}
}

// defaultValueHook is a mapstructure hook for setting default values
func defaultValueHook(
	_ reflect.Type,
	target reflect.Type,
	data interface{},
) (interface{}, error) {
	// If the data is nil or zero value, don't modify anything
	if data == nil {
		return data, nil
	}

	// For optional fields that have been omitted, the data will be nil
	return data, nil
}

// ProcessParams processes and validates all parameters for an operation
// It returns a properly typed struct containing the processed parameters
func (a *BaseAdapter) ProcessParams(operationID string, params map[string]interface{}) (interface{}, error) {
	opDef, exists := a.Operations[operationID]
	if !exists {
		return nil, fmt.Errorf("unknown operation: %s", operationID)
	}

	// Defensive: If schema is nil, this operation expects no parameters
	if opDef.Schema == nil {
		return nil, nil
	}

	// Create a new instance of the schema struct
	schemaType := reflect.TypeOf(opDef.Schema)
	if schemaType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("schema must be a pointer to a struct")
	}

	// Create a new instance of the schema struct
	result := reflect.New(schemaType.Elem()).Interface()

	// Configure mapstructure decoder
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           result,
		TagName:          "mapstructure",
		ErrorUnused:      false, // Don't error on unused fields in the input
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc("2006-01-02T15:04:05Z07:00"),
			mapstructure.StringToTimeDurationHookFunc(),
			defaultValueHook,
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create decoder: %w", err)
	}

	// Decode the parameters into the struct
	if err := decoder.Decode(params); err != nil {
		return nil, fmt.Errorf("failed to decode parameters: %w", err)
	}

	// Set defaults using struct tags
	a.setDefaults(result)

	// Validate the parameters
	if err := a.validate.Struct(result); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %w", err)
	}

	return result, nil
}

// setDefaults sets default values based on the "default" struct tag
func (a *BaseAdapter) setDefaults(obj interface{}) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Check if the field has a default tag and is zero value
		defaultValue := fieldType.Tag.Get("default")
		if defaultValue != "" && field.IsZero() {
			// Set the default value based on the field type
			switch field.Kind() {
			case reflect.String:
				field.SetString(defaultValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var intVal int64
				fmt.Sscanf(defaultValue, "%d", &intVal)
				field.SetInt(intVal)
			case reflect.Bool:
				field.SetBool(defaultValue == "true")
				// Add other type conversions as needed
			}
		}

		// Handle nested structs
		if field.Kind() == reflect.Struct {
			a.setDefaults(field.Addr().Interface())
		}
	}
}
