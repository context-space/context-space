package airtable

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bytedance/sonic"
)

// Define constants for API paths used by handlers.
const (
	endpointListBases                   = "/meta/bases"
	endpointGetBaseSchema               = "/meta/bases/{baseId}/tables"
	endpointListRecords                 = "/{baseId}/{tableIdOrName}"
	endpointGetRecord                   = "/{baseId}/{tableIdOrName}/{recordId}"
	endpointCreateRecords               = "/{baseId}/{tableIdOrName}"
	endpointUpdateRecordsPatch          = "/{baseId}/{tableIdOrName}"
	endpointUpdateRecordsPut            = "/{baseId}/{tableIdOrName}"
	endpointDeleteRecords               = "/{baseId}/{tableIdOrName}"
	endpointListWebhooks                = "/bases/{baseId}/webhooks"
	endpointCreateWebhook               = "/bases/{baseId}/webhooks"
	endpointDeleteWebhook               = "/bases/{baseId}/webhooks/{webhookId}"
	endpointManageWebhookPayloadSigning = "/bases/{baseId}/webhooks/{webhookId}/enablePayloadSigning"
	endpointRefreshWebhookPii           = "/bases/{baseId}/webhooks/{webhookId}/refresh"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDListBases                   = "list_bases"
	operationIDGetBaseSchema               = "get_base_schema"
	operationIDListRecords                 = "list_records"
	operationIDGetRecord                   = "get_record"
	operationIDCreateRecords               = "create_records"
	operationIDUpdateRecordsPatch          = "update_records_patch"
	operationIDUpdateRecordsPut            = "update_records_put"
	operationIDDeleteRecords               = "delete_records"
	operationIDListWebhooks                = "list_webhooks"
	operationIDCreateWebhook               = "create_webhook"
	operationIDDeleteWebhook               = "delete_webhook"
	operationIDManageWebhookPayloadSigning = "manage_webhook_payload_signing"
	operationIDRefreshWebhookPii           = "refresh_webhook_pii"
)

// Define a struct for each operation's parameters based on config.

// ListBasesParams defines parameters for the List Bases operation.
type ListBasesParams struct {
	// Define fields based on list_bases parameters in config

	Offset string `mapstructure:"offset" validate:"omitempty"` // Pagination offset to retrieve the next page of bases.

}

// GetBaseSchemaParams defines parameters for the Get Base Schema operation.
type GetBaseSchemaParams struct {
	// Define fields based on get_base_schema parameters in config

	Baseid string `mapstructure:"baseId" validate:"required"` // The ID of the Airtable base.

}

// ListRecordsSortParams defines the structure for a sort object in ListRecordsParams.
type ListRecordsSortParams struct {
	Field     string `mapstructure:"field"`     // Name of the field to sort by.
	Direction string `mapstructure:"direction"` // Sort direction ("asc" or "desc").
}

// ListRecordsParams defines parameters for the List Records operation.
type ListRecordsParams struct {
	// Define fields based on list_records parameters in config

	Baseid                string                  `mapstructure:"baseId" validate:"required"`                                  // The ID of the Airtable base.
	Tableidorname         string                  `mapstructure:"tableIdOrName" validate:"required"`                           // The ID or name of the Airtable table.
	Timezone              string                  `mapstructure:"timeZone" validate:"omitempty"`                               // The timezone for formula fields and createdTime. (e.g., 'America/Los_Angeles').
	Userlocale            string                  `mapstructure:"userLocale" validate:"omitempty"`                             // The user locale for formula fields. (e.g., 'en-US').
	Pagesize              int                     `mapstructure:"pageSize" validate:"omitempty,gt=0,lte=100"`                  // The number of records to return per page (max 100).
	Maxrecords            int                     `mapstructure:"maxRecords" validate:"omitempty,gt=0"`                        // The maximum total number of records to return.
	Offset                string                  `mapstructure:"offset" validate:"omitempty"`                                 // Pagination offset to retrieve the next page of records.
	View                  string                  `mapstructure:"view" validate:"omitempty"`                                   // The name or ID of a view in the table. If set, records will be returned in that view's order.
	Sort                  []ListRecordsSortParams `mapstructure:"sort" validate:"omitempty,dive"`                              // Specifies how the records will be sorted. Array of sort objects.
	Filterbyformula       string                  `mapstructure:"filterByFormula" validate:"omitempty"`                        // A formula used to filter records.
	Cellformat            string                  `mapstructure:"cellFormat" validate:"omitempty,oneof=json string"`           // Specifies how cell values are returned ('json' or 'string').
	Fields                []string                `mapstructure:"fields" validate:"omitempty,dive,required"`                   // Only return specified field Rames or field IDs. Cannot be used with returnFieldsByFieldId=true.
	Returnfieldsbyfieldid bool                    `mapstructure:"returnFieldsByFieldId" validate:"omitempty"`                  // Use field IDs in the response 'fields' key instead of field names. Cannot be used with the fields parameter.
	Recordmetadata        []string                `mapstructure:"recordMetadata" validate:"omitempty,dive,oneof=commentCount"` // An array of strings for desired record metadata. Currently the only supported option is 'commentCount'.
}

// GetRecordParams defines parameters for the Get Record operation.
type GetRecordParams struct {
	// Define fields based on get_record parameters in config

	Baseid                string `mapstructure:"baseId" validate:"required"`                        // The ID of the Airtable base.
	Tableidorname         string `mapstructure:"tableIdOrName" validate:"required"`                 // The ID or name of the Airtable table.
	Recordid              string `mapstructure:"recordId" validate:"required"`                      // The ID of the record to retrieve.
	Cellformat            string `mapstructure:"cellFormat" validate:"omitempty,oneof=json string"` // Specifies how cell values are returned ('json' or 'string').
	Returnfieldsbyfieldid bool   `mapstructure:"returnFieldsByFieldId" validate:"omitempty"`        // Use field IDs in the response 'fields' key instead of field names.

}

// CreateRecordItemParams defines a single record item for creation.
type CreateRecordItemParams struct {
	Fields map[string]interface{} `mapstructure:"fields" validate:"required"`
}

// CreateRecordsParams defines parameters for the Create Records operation.
type CreateRecordsParams struct {
	Baseid                    string                   `mapstructure:"baseId" validate:"required"`                      // The ID of the Airtable base.
	Tableidorname             string                   `mapstructure:"tableIdOrName" validate:"required"`               // The ID or name of the Airtable table.
	RecordsBody               []CreateRecordItemParams `mapstructure:"records_body" validate:"required,min=1,dive"`     // Array of record objects to create.
	TypecastBody              *bool                    `mapstructure:"typecast_body" validate:"omitempty"`              // Attempt to convert string values to Airtable cell values. Defaults to false.
	ReturnFieldsByFieldIdBody *bool                    `mapstructure:"returnFieldsByFieldId_body" validate:"omitempty"` // Use field IDs in the response 'fields' key instead of field names. Defaults to false.
}

// UpdateRecordItemParams defines a single record item for update (PATCH/PUT).
type UpdateRecordItemParams struct {
	ID     string                 `mapstructure:"id" validate:"required"`
	Fields map[string]interface{} `mapstructure:"fields" validate:"required"`
}

// UpdateRecordsPatchParams defines parameters for the Update Records (PATCH) operation.
type UpdateRecordsPatchParams struct {
	Baseid                    string                   `mapstructure:"baseId" validate:"required"`                      // The ID of the Airtable base.
	Tableidorname             string                   `mapstructure:"tableIdOrName" validate:"required"`               // The ID or name of the Airtable table.
	RecordsBody               []UpdateRecordItemParams `mapstructure:"records_body" validate:"required,min=1,dive"`     // Array of record objects to update.
	TypecastBody              *bool                    `mapstructure:"typecast_body" validate:"omitempty"`              // Attempt to convert string values to Airtable cell values.
	ReturnFieldsByFieldIdBody *bool                    `mapstructure:"returnFieldsByFieldId_body" validate:"omitempty"` // Use field IDs in the response 'fields' key.
}

// UpdateRecordsPutParams defines parameters for the Update Records (PUT) operation.
type UpdateRecordsPutParams struct {
	Baseid                    string                   `mapstructure:"baseId" validate:"required"`
	Tableidorname             string                   `mapstructure:"tableIdOrName" validate:"required"`
	RecordsBody               []UpdateRecordItemParams `mapstructure:"records_body" validate:"required,min=1,dive"`
	TypecastBody              *bool                    `mapstructure:"typecast_body" validate:"omitempty"`
	ReturnFieldsByFieldIdBody *bool                    `mapstructure:"returnFieldsByFieldId_body" validate:"omitempty"`
}

// DeleteRecordsParams defines parameters for the Delete Records operation.
type DeleteRecordsParams struct {
	// Define fields based on delete_records parameters in config

	Baseid        string   `mapstructure:"baseId" validate:"required"`        // The ID of the Airtable base.
	Tableidorname string   `mapstructure:"tableIdOrName" validate:"required"` // The ID or name of the Airtable table.
	Records       []string `mapstructure:"records" validate:"required,min=1"` // An array of record IDs to delete.

}

// ListWebhooksParams defines parameters for the List Webhooks operation.
type ListWebhooksParams struct {
	// Define fields based on list_webhooks parameters in config

	Baseid string `mapstructure:"baseId" validate:"required"`  // The ID of the Airtable base.
	Cursor string `mapstructure:"cursor" validate:"omitempty"` // Cursor for pagination to retrieve the next set of webhooks.

}

// CreateWebhookSpecificationOptionsFiltersParams ...
type CreateWebhookSpecificationOptionsFiltersParams struct {
	DataTypes               []string `mapstructure:"dataTypes" json:"dataTypes,omitempty" validate:"omitempty,dive,oneof=tableData tableFields tableMetadata"`
	RecordChangeScope       string   `mapstructure:"recordChangeScope" json:"recordChangeScope,omitempty" validate:"omitempty"`
	WatchDataInFieldIds     []string `mapstructure:"watchDataInFieldIds" json:"watchDataInFieldIds,omitempty" validate:"omitempty,dive"`
	WatchSchemasForTableIds []string `mapstructure:"watchSchemasForTableIds" json:"watchSchemasForTableIds,omitempty" validate:"omitempty,dive"`
}

// CreateWebhookSpecificationOptionsIncludesParams ...
type CreateWebhookSpecificationOptionsIncludesParams struct {
	Data               *bool `mapstructure:"data" json:"data,omitempty" validate:"omitempty"`
	Schema             *bool `mapstructure:"schema" json:"schema,omitempty" validate:"omitempty"`
	PreviousData       *bool `mapstructure:"previousData" json:"previousData,omitempty" validate:"omitempty"`
	CellValuesAsParent *bool `mapstructure:"cellValuesAsParent" json:"cellValuesAsParent,omitempty" validate:"omitempty"`
}

// CreateWebhookSpecificationOptionsParams ...
type CreateWebhookSpecificationOptionsParams struct {
	Filters  *CreateWebhookSpecificationOptionsFiltersParams  `mapstructure:"filters" json:"filters,omitempty" validate:"omitempty"`
	Includes *CreateWebhookSpecificationOptionsIncludesParams `mapstructure:"includes" json:"includes,omitempty" validate:"omitempty"`
}

// CreateWebhookSpecificationParams ...
type CreateWebhookSpecificationParams struct {
	Options *CreateWebhookSpecificationOptionsParams `mapstructure:"options" json:"options" validate:"required"` // 'options' is required by Airtable
}

// CreateWebhookParams defines parameters for the Create Webhook operation.
type CreateWebhookParams struct {
	Baseid              string                            `mapstructure:"baseId" json:"-"` // Path param, not in body
	NotificationURLBody string                            `mapstructure:"notificationUrl_body" json:"notificationUrl" validate:"required,url"`
	SpecificationBody   *CreateWebhookSpecificationParams `mapstructure:"specification_body" json:"specification" validate:"required"`
}

// DeleteWebhookParams defines parameters for the Delete Webhook operation.
type DeleteWebhookParams struct {
	// Define fields based on delete_webhook parameters in config

	Baseid    string `mapstructure:"baseId" validate:"required"`    // The ID of the Airtable base.
	Webhookid string `mapstructure:"webhookId" validate:"required"` // The ID of the webhook to delete.

}

// ManageWebhookPayloadSigningParams defines parameters for the Enable/Disable Webhook Payload Signing operation.
type ManageWebhookPayloadSigningParams struct {
	Baseid     string `mapstructure:"baseId" validate:"required"`
	Webhookid  string `mapstructure:"webhookId" validate:"required"`
	EnableBody *bool  `mapstructure:"enable_body" validate:"required"` // Pointer to allow explicit true/false
}

// RefreshWebhookPiiParams defines parameters for the Refresh Webhook PII operation.
type RefreshWebhookPiiParams struct {
	// Define fields based on refresh_webhook_pii parameters in config

	Baseid    string `mapstructure:"baseId" validate:"required"`    // The ID of the Airtable base.
	Webhookid string `mapstructure:"webhookId" validate:"required"` // The ID of the webhook to refresh.

}

// OperationHandler defines the function signature for handling a specific API operation.
type OperationHandler func(ctx context.Context, params interface{}) (map[string]interface{}, error)

// OperationDefinition combines parameter schema and handler.
type OperationDefinition struct {
	Schema                interface{}      // Parameter schema (struct pointer)
	Handler               OperationHandler // Operation handler function
	PermissionIdentifiers []string         // List of internal permission identifiers required
}

// Operations maps operation IDs to their definitions.
type Operations map[string]OperationDefinition

// RegisterOperation registers the parameter schema and handler.
func (a *AirtableAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler, requiredPerms []string) {
	a.BaseAdapter.RegisterOperation(operationID, schema) // Register schema for validation
	if a.operations == nil {
		a.operations = make(Operations)
	}
	a.operations[operationID] = OperationDefinition{
		Schema:                schema,
		Handler:               handler,
		PermissionIdentifiers: requiredPerms,
	}
}

// registerOperations is called by the adapter constructor to register all supported operations.
func (a *AirtableAdapter) registerOperations() {
	a.RegisterOperation(
		operationIDListBases,
		&ListBasesParams{},
		handleListBases,
		[]string{"read_schema"},
	)

	a.RegisterOperation(
		operationIDGetBaseSchema,
		&GetBaseSchemaParams{},
		handleGetBaseSchema,
		[]string{"read_schema"},
	)

	a.RegisterOperation(
		operationIDListRecords,
		&ListRecordsParams{},
		handleListRecords,
		[]string{"read_records"},
	)

	a.RegisterOperation(
		operationIDGetRecord,
		&GetRecordParams{},
		handleGetRecord,
		[]string{"read_records"},
	)

	a.RegisterOperation(
		operationIDCreateRecords,
		&CreateRecordsParams{},
		handleCreateRecords,
		[]string{"write_records"},
	)

	a.RegisterOperation(
		operationIDUpdateRecordsPatch,
		&UpdateRecordsPatchParams{},
		handleUpdateRecordsPatch,
		[]string{"write_records"},
	)

	a.RegisterOperation(
		operationIDUpdateRecordsPut,
		&UpdateRecordsPutParams{},
		handleUpdateRecordsPut,
		[]string{"write_records"},
	)

	a.RegisterOperation(
		operationIDDeleteRecords,
		&DeleteRecordsParams{},
		handleDeleteRecords,
		[]string{"write_records"},
	)

	a.RegisterOperation(
		operationIDListWebhooks,
		&ListWebhooksParams{},
		handleListWebhooks,
		[]string{"manage_webhooks"},
	)

	a.RegisterOperation(
		operationIDCreateWebhook,
		&CreateWebhookParams{},
		handleCreateWebhook,
		[]string{"manage_webhooks"},
	)

	a.RegisterOperation(
		operationIDDeleteWebhook,
		&DeleteWebhookParams{},
		handleDeleteWebhook,
		[]string{"manage_webhooks"},
	)

	a.RegisterOperation(
		operationIDManageWebhookPayloadSigning,
		&ManageWebhookPayloadSigningParams{},
		handleManageWebhookPayloadSigning,
		[]string{"manage_webhooks"},
	)

	a.RegisterOperation(
		operationIDRefreshWebhookPii,
		&RefreshWebhookPiiParams{},
		handleRefreshWebhookPii,
		[]string{"manage_webhooks"},
	)
}

// handleListBases constructs parameters for the REST adapter for the List Bases operation.
func handleListBases(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ListBasesParams)
	if !ok {
		if params == nil {
			p = &ListBasesParams{}
		} else {
			return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation list_bases", params)
		}
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)

	if p.Offset != "" {
		queryParams.Set("offset", p.Offset)
	}

	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointListBases,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}

	return restParams, nil
}

// handleGetBaseSchema constructs parameters for the REST adapter for the Get Base Schema operation.
func handleGetBaseSchema(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*GetBaseSchemaParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_base_schema, or missing required parameters", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)

	if p.Baseid == "" {
		return nil, fmt.Errorf("missing required path parameter: baseId for get_base_schema")
	}
	pathParams["baseId"] = p.Baseid

	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetBaseSchema,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}

	return restParams, nil
}

// handleListRecords constructs parameters for the REST adapter for the List Records operation.
func handleListRecords(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ListRecordsParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation list_records, or missing required parameters", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)

	if p.Baseid == "" || p.Tableidorname == "" {
		return nil, fmt.Errorf("missing required path parameters for list_records: baseId and/or tableIdOrName")
	}
	pathParams["baseId"] = p.Baseid
	pathParams["tableIdOrName"] = p.Tableidorname

	if p.Timezone != "" {
		queryParams.Set("timeZone", p.Timezone)
	}
	if p.Userlocale != "" {
		queryParams.Set("userLocale", p.Userlocale)
	}
	if p.Pagesize > 0 {
		queryParams.Set("pageSize", strconv.Itoa(p.Pagesize))
	}
	if p.Maxrecords > 0 {
		queryParams.Set("maxRecords", strconv.Itoa(p.Maxrecords))
	}
	if p.Offset != "" {
		queryParams.Set("offset", p.Offset)
	}
	if p.View != "" {
		queryParams.Set("view", p.View)
	}
	if p.Filterbyformula != "" {
		queryParams.Set("filterByFormula", p.Filterbyformula)
	}
	if p.Cellformat != "" {
		queryParams.Set("cellFormat", p.Cellformat)
	}
	if p.Returnfieldsbyfieldid {
		queryParams.Set("returnFieldsByFieldId", "true")
	}

	for _, field := range p.Fields {
		queryParams.Add("fields[]", field)
	}

	for i, sortObj := range p.Sort {
		if sortObj.Field != "" {
			queryParams.Set(fmt.Sprintf("sort[%d][field]", i), sortObj.Field)
			if sortObj.Direction != "" {
				queryParams.Set(fmt.Sprintf("sort[%d][direction]", i), sortObj.Direction)
			} else { // Airtable defaults to asc if direction is omitted but field is present
				queryParams.Set(fmt.Sprintf("sort[%d][direction]", i), "asc")
			}
		}
	}

	for _, meta := range p.Recordmetadata {
		queryParams.Add("recordMetadata[]", meta)
	}

	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointListRecords,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}

	return restParams, nil
}

// handleGetRecord constructs parameters for the REST adapter for the Get Record operation.
func handleGetRecord(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*GetRecordParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_record, or missing required parameters", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)

	if p.Baseid == "" || p.Tableidorname == "" || p.Recordid == "" {
		return nil, fmt.Errorf("missing required path parameters for get_record")
	}
	pathParams["baseId"] = p.Baseid
	pathParams["tableIdOrName"] = p.Tableidorname
	pathParams["recordId"] = p.Recordid

	if p.Cellformat != "" {
		queryParams.Set("cellFormat", p.Cellformat)
	}
	if p.Returnfieldsbyfieldid {
		queryParams.Set("returnFieldsByFieldId", "true")
	}

	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetRecord,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}

	return restParams, nil
}

// handleCreateRecords constructs parameters for the REST adapter for the Create Records operation.
func handleCreateRecords(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*CreateRecordsParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation create_records, or missing required parameters", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)
	requestBodyMap := make(map[string]interface{})

	if p.Baseid == "" || p.Tableidorname == "" {
		return nil, fmt.Errorf("missing required path parameters for create_records")
	}
	pathParams["baseId"] = p.Baseid
	pathParams["tableIdOrName"] = p.Tableidorname

	// Populate requestBodyMap from p
	if p.RecordsBody != nil {
		// Explicitly build the records structure for the JSON body
		recordsForJSON := make([]map[string]interface{}, len(p.RecordsBody))
		for i, rb := range p.RecordsBody {
			recordsForJSON[i] = map[string]interface{}{
				"fields": rb.Fields,
			}
		}
		requestBodyMap["records"] = recordsForJSON
	}
	if p.TypecastBody != nil {
		requestBodyMap["typecast"] = *p.TypecastBody
	}
	if p.ReturnFieldsByFieldIdBody != nil {
		requestBodyMap["returnFieldsByFieldId"] = *p.ReturnFieldsByFieldIdBody
	}

	if len(requestBodyMap) <= 0 {
		return nil, fmt.Errorf("request body for update_records_put is empty, 'records_body' parameter is required")
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointCreateRecords,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}
	restParams["body"] = requestBodyMap

	return restParams, nil
}

// handleUpdateRecordsPatch constructs parameters for the REST adapter for the Update Records (PATCH) operation.
func handleUpdateRecordsPatch(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*UpdateRecordsPatchParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation update_records_patch", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)
	requestBodyMap := make(map[string]interface{})

	if p.Baseid == "" || p.Tableidorname == "" {
		return nil, fmt.Errorf("missing required path parameters for update_records_patch")
	}
	pathParams["baseId"] = p.Baseid
	pathParams["tableIdOrName"] = p.Tableidorname

	// Populate requestBodyMap from p
	if p.RecordsBody != nil {
		// Explicitly build the records structure for the JSON body
		recordsForJSON := make([]map[string]interface{}, len(p.RecordsBody))
		for i, rb := range p.RecordsBody {
			recordsForJSON[i] = map[string]interface{}{
				"id":     rb.ID,
				"fields": rb.Fields,
			}
		}
		requestBodyMap["records"] = recordsForJSON
	}
	if p.TypecastBody != nil {
		requestBodyMap["typecast"] = *p.TypecastBody
	}
	if p.ReturnFieldsByFieldIdBody != nil {
		requestBodyMap["returnFieldsByFieldId"] = *p.ReturnFieldsByFieldIdBody
	}

	if len(requestBodyMap) <= 0 {
		return nil, fmt.Errorf("request body for update_records_put is empty, 'records_body' parameter is required")
	}

	restParams := map[string]interface{}{
		"method": http.MethodPatch,
		"path":   endpointUpdateRecordsPatch,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}
	restParams["body"] = requestBodyMap

	return restParams, nil
}

// handleUpdateRecordsPut constructs parameters for the REST adapter for the Update Records (PUT) operation.
func handleUpdateRecordsPut(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*UpdateRecordsPutParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation update_records_put", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)
	requestBodyMap := make(map[string]interface{})

	if p.Baseid == "" || p.Tableidorname == "" {
		return nil, fmt.Errorf("missing required path parameters for update_records_put")
	}
	pathParams["baseId"] = p.Baseid
	pathParams["tableIdOrName"] = p.Tableidorname

	// Populate requestBodyMap from p
	if p.RecordsBody != nil {
		// Explicitly build the records structure for the JSON body
		recordsForJSON := make([]map[string]interface{}, len(p.RecordsBody))
		for i, rb := range p.RecordsBody {
			recordsForJSON[i] = map[string]interface{}{
				"id":     rb.ID,
				"fields": rb.Fields,
			}
		}
		requestBodyMap["records"] = recordsForJSON
	}
	if p.TypecastBody != nil {
		requestBodyMap["typecast"] = *p.TypecastBody
	}
	if p.ReturnFieldsByFieldIdBody != nil {
		requestBodyMap["returnFieldsByFieldId"] = *p.ReturnFieldsByFieldIdBody
	}

	if len(requestBodyMap) <= 0 {
		return nil, fmt.Errorf("request body for update_records_put is empty, 'records_body' parameter is required")
	}

	restParams := map[string]interface{}{
		"method": http.MethodPut,
		"path":   endpointUpdateRecordsPut,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}
	restParams["body"] = requestBodyMap

	return restParams, nil
}

// handleDeleteRecords constructs parameters for the REST adapter for the Delete Records operation.
func handleDeleteRecords(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*DeleteRecordsParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation delete_records", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)

	if p.Baseid == "" || p.Tableidorname == "" {
		return nil, fmt.Errorf("missing required path parameters for delete_records")
	}
	pathParams["baseId"] = p.Baseid
	pathParams["tableIdOrName"] = p.Tableidorname

	if len(p.Records) <= 0 {
		return nil, fmt.Errorf("records parameter is required for delete_records and was empty")
	}

	for _, recordID := range p.Records {
		queryParams.Add("records[]", recordID) // Use key "records[]" as per Airtable docs
	}

	restParams := map[string]interface{}{
		"method": http.MethodDelete,
		"path":   endpointDeleteRecords,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}

	return restParams, nil
}

// handleListWebhooks constructs parameters for the REST adapter for the List Webhooks operation.
func handleListWebhooks(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ListWebhooksParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation list_webhooks", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)

	if p.Baseid == "" {
		return nil, fmt.Errorf("missing required path parameter: baseId for list_webhooks")
	}
	pathParams["baseId"] = p.Baseid

	if p.Cursor != "" {
		queryParams.Set("cursor", p.Cursor)
	}

	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointListWebhooks,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}

	return restParams, nil
}

// handleCreateWebhook constructs parameters for the REST adapter for the Create Webhook operation.
func handleCreateWebhook(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*CreateWebhookParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation create_webhook", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)
	requestBodyMap := make(map[string]interface{})

	if p.Baseid == "" {
		return nil, fmt.Errorf("missing required path parameter: baseId for create_webhook")
	}
	pathParams["baseId"] = p.Baseid

	// Populate requestBodyMap from p
	requestBodyMap["notificationUrl"] = p.NotificationURLBody
	if p.SpecificationBody != nil {
		requestBodyMap["specification"] = p.SpecificationBody // This is a pointer to a struct, will be marshalled to JSON object by the HTTP client
	}

	// For debugging, print the requestBodyMap
	b, _ := sonic.MarshalIndent(requestBodyMap, "", "  ")
	fmt.Printf("DEBUG: CreateWebhook requestBodyMap JSON:\n%s\n", string(b))

	if len(requestBodyMap) <= 0 {
		return nil, fmt.Errorf("request body for create_webhook is empty, 'notificationUrl_body' and 'specification_body' are required")
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointCreateWebhook,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}
	restParams["body"] = requestBodyMap

	return restParams, nil
}

// handleDeleteWebhook constructs parameters for the REST adapter for the Delete Webhook operation.
func handleDeleteWebhook(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*DeleteWebhookParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation delete_webhook", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)

	if p.Baseid == "" || p.Webhookid == "" {
		return nil, fmt.Errorf("missing required path parameters for delete_webhook")
	}
	pathParams["baseId"] = p.Baseid
	pathParams["webhookId"] = p.Webhookid

	restParams := map[string]interface{}{
		"method": http.MethodDelete,
		"path":   endpointDeleteWebhook,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}

	return restParams, nil
}

// handleManageWebhookPayloadSigning constructs parameters for the REST adapter for the Enable/Disable Webhook Payload Signing operation.
func handleManageWebhookPayloadSigning(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ManageWebhookPayloadSigningParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation manage_webhook_payload_signing", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)
	requestBodyMap := make(map[string]interface{})

	if p.Baseid == "" || p.Webhookid == "" {
		return nil, fmt.Errorf("missing required path parameters for manage_webhook_payload_signing")
	}
	pathParams["baseId"] = p.Baseid
	pathParams["webhookId"] = p.Webhookid

	// Populate requestBodyMap from p
	if p.EnableBody == nil {
		// This field is required by validation, so this case implies a validation bypass or issue.
		return nil, fmt.Errorf("missing required body parameter: enable_body for manage_webhook_payload_signing")
	}
	requestBodyMap["enable"] = *p.EnableBody

	if len(requestBodyMap) <= 0 {
		return nil, fmt.Errorf("request body for manage_webhook_payload_signing is empty, 'enable_body' parameter is required")
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointManageWebhookPayloadSigning,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}

	restParams["body"] = requestBodyMap

	return restParams, nil
}

// handleRefreshWebhookPii constructs parameters for the REST adapter for the Refresh Webhook PII operation.
func handleRefreshWebhookPii(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*RefreshWebhookPiiParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation refresh_webhook_pii", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headers := make(map[string]string)

	if p.Baseid == "" || p.Webhookid == "" {
		return nil, fmt.Errorf("missing required path parameters for refresh_webhook_pii")
	}
	pathParams["baseId"] = p.Baseid
	pathParams["webhookId"] = p.Webhookid

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointRefreshWebhookPii,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	if len(headers) > 0 {
		restParams["headers"] = headers
	}

	return restParams, nil
}
