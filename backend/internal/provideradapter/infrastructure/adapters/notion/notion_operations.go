package notion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
)

// Constants define API paths relative to the BaseURL (https://api.notion.com/v1)
const (
	// Pages
	endpointCreatePage           = "pages"
	endpointGetPage              = "pages/{page_id}"
	endpointUpdatePageProperties = "pages/{page_id}"

	// Blocks
	endpointAppendBlockChildren = "blocks/{block_id}/children"

	// Databases
	endpointQueryDatabase = "databases/{database_id}/query"
	endpointGetDatabase   = "databases/{database_id}"
	// Note: Listing databases uses the search endpoint
	endpointListDatabasesViaSearch = "search"

	// Search
	endpointSearch = "search"

	// Users
	endpointListUsers = "users"
	endpointGetUserMe = "users/me" // For validation

)

// Define constants for operation IDs used by handlers.
const (
	opCreatePage           = "create_page"
	opSearch               = "search"
	opQueryDatabase        = "query_database"
	opAppendToBlock        = "append_to_block"
	opGetPage              = "get_page"
	opUpdatePageProperties = "update_page_properties"
	opListUsers            = "list_users"
	opGetDatabase          = "get_database"
	opListDatabases        = "list_databases"
)

// Create_pageParams defines parameters for the Create Page operation.
type Create_pageParams struct {

	// The ID of the parent page or database.
	Parent_id string `mapstructure:"parent_id" validate:"required"` // Required: true in json

	// The type of the parent, either 'page' or 'database'.
	Parent_type string `mapstructure:"parent_type" validate:"required,oneof=page database"` // Required: true in json

	// The title of the new page (required for both page and database entry).
	Title string `mapstructure:"title" validate:"required"` // Required: true in json

	// Optional. A JSON string representing an array of Notion block objects to add as content. See Notion API docs for block structure.
	Content_blocks_json string `mapstructure:"content_blocks_json" validate:"omitempty"` // Required: false in json

	// Optional. A JSON string representing the page properties when creating an entry in a database. Keys must match database schema. Title property is handled separately via the 'title' parameter. See Notion API docs.
	Database_properties_json string `mapstructure:"database_properties_json" validate:"omitempty"` // Required: false in json

}

// SearchParams defines parameters for the Search Pages/Databases operation.
type SearchParams struct {

	// The text to search for in page or database titles. Leave empty to list all accessible items.
	Query string `mapstructure:"query" validate:"omitempty"` // Required: false in json

	// Filter results to only 'page' or 'database'.
	Filter_type string `mapstructure:"filter_type" validate:"omitempty,oneof=page database ''"` // Required: false in json

	// Field to sort results by.
	Sort_by string `mapstructure:"sort_by" validate:"omitempty,oneof=last_edited_time ''"` // Required: false in json

	// Sort direction.
	Sort_direction string `mapstructure:"sort_direction" validate:"omitempty,oneof=ascending descending ''"` // Required: false in json

}

// Query_databaseParams defines parameters for the Query Database operation.
type Query_databaseParams struct {

	// The ID of the database to query.
	Database_id string `mapstructure:"database_id" validate:"required"` // Required: true in json

	// Optional. A JSON string representing the Notion filter object. See Notion API documentation for structure.
	Filter_json string `mapstructure:"filter_json" validate:"omitempty"` // Required: false in json

	// Optional. A JSON string representing the Notion sorts array (list of sort objects). See Notion API documentation for structure.
	Sorts_json string `mapstructure:"sorts_json" validate:"omitempty"` // Required: false in json

}

// Append_to_blockParams defines parameters for the Append Content operation.
type Append_to_blockParams struct {

	// The ID of the block (e.g., page) to append children to.
	Block_id string `mapstructure:"block_id" validate:"required"` // Required: true in json

	// A JSON string representing an array of Notion block objects to append. See Notion API docs for block structure.
	Content_blocks_json string `mapstructure:"content_blocks_json" validate:"required"` // Required: true in json

}

// Get_pageParams defines parameters for the Get Page Info operation.
type Get_pageParams struct {

	// The ID of the page to retrieve.
	Page_id string `mapstructure:"page_id" validate:"required"` // Required: true in json

}

// Update_page_propertiesParams defines parameters for the Update Page Properties operation.
type Update_page_propertiesParams struct {

	// The ID of the page to update.
	Page_id string `mapstructure:"page_id" validate:"required"` // Required: true in json

	// A JSON string representing the properties object to update. See Notion API docs for structure.
	Properties_update_json string `mapstructure:"properties_update_json" validate:"required"` // Required: true in json

	// Optional. Set to true to archive the page, false to unarchive.
	Archived bool `mapstructure:"archived" validate:"omitempty"` // Required: false in json

}

// List_usersParams defines parameters for the List Users operation.
type List_usersParams struct {
}

// Get_databaseParams defines parameters for the Get Database Info operation.
type Get_databaseParams struct {

	// The ID of the database to retrieve.
	Database_id string `mapstructure:"database_id" validate:"required"` // Required: true in json

}

// List_databasesParams defines parameters for the List Databases operation.
type List_databasesParams struct {
}

// OperationHandler defines the function signature for handling a specific API operation.
// It now receives context and processed parameters, and returns parameters for the REST adapter.
type OperationHandler func(ctx context.Context, params interface{}) (map[string]interface{}, error)

// OperationDefinition combines parameter schema and handler. (Response Type removed)
type OperationDefinition struct {
	Schema                interface{}      // Parameter schema (struct pointer)
	Handler               OperationHandler // Operation handler function
	PermissionIdentifiers []string         // List of internal permission identifiers required
}

// Operations maps operation IDs to their definitions.
type Operations map[string]OperationDefinition

// RegisterOperation registers the parameter schema and handler.
// Method and Path are no longer passed here. ResponseType is also removed.
func (a *NotionAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler, requiredPerms []string) {
	a.BaseAdapter.RegisterOperation(operationID, schema) // Register schema for validation with the base adapter
	if a.operations == nil {
		a.operations = make(Operations)
	}
	// Store the handler and permission requirements locally for execution and permission checks
	a.operations[operationID] = OperationDefinition{
		Schema:                schema,
		Handler:               handler,
		PermissionIdentifiers: requiredPerms,
	}
}

// registerOperations is called by the adapter constructor to register all supported operations.
// Method and path are now handled within each operation handler.
func (a *NotionAdapter) registerOperations() {
	a.RegisterOperation(
		opCreatePage,
		&Create_pageParams{},
		handleCreate_page,
		[]string{"access_notion"},
	)

	a.RegisterOperation(
		opSearch,
		&SearchParams{},
		handleSearch,
		[]string{"access_notion"},
	)

	a.RegisterOperation(
		opQueryDatabase,
		&Query_databaseParams{},
		handleQuery_database,
		[]string{"access_notion"},
	)

	a.RegisterOperation(
		opAppendToBlock,
		&Append_to_blockParams{},
		handleAppend_to_block,
		[]string{"access_notion"},
	)

	a.RegisterOperation(
		opGetPage,
		&Get_pageParams{},
		handleGet_page,
		[]string{"access_notion"},
	)

	a.RegisterOperation(
		opUpdatePageProperties,
		&Update_page_propertiesParams{},
		handleUpdate_page_properties,
		[]string{"access_notion"},
	)

	a.RegisterOperation(
		opListUsers,
		&List_usersParams{},
		handleList_users,
		[]string{"access_notion"},
	)

	a.RegisterOperation(
		opGetDatabase,
		&Get_databaseParams{},
		handleGet_database,
		[]string{"access_notion"},
	)

	a.RegisterOperation(
		opListDatabases,
		&List_databasesParams{},
		handleList_databases,
		[]string{"access_notion"},
	)
}

// handleCreate_page constructs parameters for the REST adapter to create a page.
func handleCreate_page(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	// 1. Cast params
	p, ok := params.(*Create_pageParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation create_page", params)
	}

	// 2. Construct request body using CreatePageRequest struct
	reqBody := &CreatePageRequest{
		Properties: make(map[string]interface{}), // Initialize properties map
	}

	// Set parent based on type
	if p.Parent_type == "page" {
		reqBody.Parent = Parent{Type: "page_id", PageID: p.Parent_id}
	} else if p.Parent_type == "database" {
		reqBody.Parent = Parent{Type: "database_id", DatabaseID: p.Parent_id}
	} else {
		return nil, fmt.Errorf("invalid parent_type: %s", p.Parent_type)
	}

	// Construct the title property value
	titleValue := []RichText{
		{
			Type: "text",
			Text: &TextContent{Content: p.Title},
			// Annotations and PlainText are omitted for request, Notion generates them
		},
	}

	// Handle Properties based on parent type
	if p.Parent_type == "page" {
		// For page parent, only 'title' property is allowed
		reqBody.Properties["title"] = map[string]interface{}{"title": titleValue}
	} else { // Parent is database
		// Expect database_properties_json to contain all properties, *including* the title property
		// keyed by the database's actual title column name/ID.
		if p.Database_properties_json != "" {
			var dbProps map[string]interface{}
			if err := sonic.Unmarshal([]byte(p.Database_properties_json), &dbProps); err == nil {
				reqBody.Properties = dbProps
			} else {
				return nil, fmt.Errorf("failed to parse database_properties_json: %w", err)
			}
		} else {
			// If database_properties_json is empty, we cannot guess the title column name.
			// Return an error or log a warning? Returning error is safer.
			// The user MUST provide the title property structure within database_properties_json.
			return nil, fmt.Errorf("database_properties_json is required and must include the title property when parent_type is 'database'")
		}
		// Ensure the provided title parameter matches the title in the JSON properties (optional check)
		// This requires knowing the title property key from the JSON, which is complex. Skip for now.
	}

	// Handle optional Content Blocks
	if p.Content_blocks_json != "" {
		var contentBlocks []map[string]interface{} // Keep as map for block creation flexibility
		if err := sonic.Unmarshal([]byte(p.Content_blocks_json), &contentBlocks); err == nil {
			reqBody.Children = &contentBlocks
		} else {
			return nil, fmt.Errorf("failed to parse content_blocks_json: %w", err)
		}
	}

	// 3. Prepare parameters for REST Adapter
	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointCreatePage, // Use constant
		"body":   reqBody,            // Use the structured request body
	}

	return restParams, nil
}

// handleSearch constructs parameters for the REST adapter to search pages/databases.
func handleSearch(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*SearchParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation search", params)
	}

	// Construct request body using SearchRequest struct
	body := &SearchRequest{}

	if p.Query != "" {
		body.Query = &p.Query
	}

	if p.Filter_type != "" {
		// Validate filter type if needed, e.g., oneof=page database
		body.Filter = &SearchFilter{
			Property: "object",
			Value:    p.Filter_type,
		}
	}

	if p.Sort_by != "" || p.Sort_direction != "" {
		sort := &SearchSort{}
		if p.Sort_by == "last_edited_time" { // Currently only last_edited_time is supported by API
			sort.Timestamp = "last_edited_time"
		} else if p.Sort_by != "" {
			// Log a warning or return error if unsupported sort field is provided?
			fmt.Printf("[WARN] handleSearch: Unsupported sort_by value '%s'. Only 'last_edited_time' is supported by Notion API.\n", p.Sort_by)
			// Defaulting to last_edited_time or omitting sort?
			// Let's default to last_edited_time if a direction is provided.
			if p.Sort_direction != "" {
				sort.Timestamp = "last_edited_time"
			}
		}

		if p.Sort_direction == "ascending" || p.Sort_direction == "descending" {
			sort.Direction = p.Sort_direction
		} else if p.Sort_direction != "" {
			// Log warning or return error for invalid direction?
			fmt.Printf("[WARN] handleSearch: Invalid sort_direction '%s'. Using default 'descending'.\n", p.Sort_direction)
			sort.Direction = "descending" // Default
		} else {
			// Default direction if only sort_by is provided
			if sort.Timestamp != "" {
				sort.Direction = "descending"
			}
		}

		// Only add sort object if valid fields were set
		if sort.Timestamp != "" && sort.Direction != "" {
			body.Sort = sort
		}
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointSearch,
		"body":   body,
	}

	return restParams, nil
}

// handleQuery_database constructs parameters for the REST adapter to query a database.
func handleQuery_database(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*Query_databaseParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation query_database", params)
	}

	// Construct request body using QueryDatabaseRequest struct
	body := &QueryDatabaseRequest{}

	if p.Filter_json != "" {
		var filter map[string]interface{}
		if err := sonic.Unmarshal([]byte(p.Filter_json), &filter); err != nil {
			return nil, fmt.Errorf("failed to parse filter_json: %w", err)
		}
		body.Filter = filter
	}

	if p.Sorts_json != "" {
		var sorts []map[string]interface{}
		if err := sonic.Unmarshal([]byte(p.Sorts_json), &sorts); err != nil {
			return nil, fmt.Errorf("failed to parse sorts_json: %w", err)
		}
		body.Sorts = sorts
	}

	// TODO: Add pagination (StartCursor, PageSize)
	// Example: if pageSize > 0 { body.PageSize = &pageSize }
	// Example: if startCursor != "" { body.StartCursor = &startCursor }

	restParams := map[string]interface{}{
		"method":      http.MethodPost,
		"path":        endpointQueryDatabase, // Use constant template
		"path_params": map[string]string{"database_id": p.Database_id},
		"body":        body,
	}

	return restParams, nil
}

// handleAppend_to_block constructs parameters for the REST adapter to append content.
func handleAppend_to_block(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*Append_to_blockParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation append_to_block", params)
	}

	// Construct request body using AppendBlockChildrenRequest struct
	body := &AppendBlockChildrenRequest{}

	var contentBlocks []map[string]interface{} // Keep as map for block creation flexibility
	if err := sonic.Unmarshal([]byte(p.Content_blocks_json), &contentBlocks); err != nil {
		return nil, fmt.Errorf("failed to parse content_blocks_json: %w", err)
	}
	body.Children = contentBlocks

	restParams := map[string]interface{}{
		"method":      http.MethodPatch,
		"path":        endpointAppendBlockChildren, // Use constant template
		"path_params": map[string]string{"block_id": p.Block_id},
		"body":        body,
	}

	return restParams, nil
}

// handleGet_page constructs parameters for the REST adapter to get page info.
func handleGet_page(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*Get_pageParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_page", params)
	}

	restParams := map[string]interface{}{
		"method":      http.MethodGet,
		"path":        endpointGetPage,
		"path_params": map[string]string{"page_id": p.Page_id},
	}

	return restParams, nil
}

// handleUpdate_page_properties constructs parameters for the REST adapter to update page properties.
func handleUpdate_page_properties(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*Update_page_propertiesParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation update_page_properties", params)
	}

	// Construct request body using UpdatePageRequest struct
	body := &UpdatePageRequest{}

	var propertiesUpdate map[string]interface{}
	if err := sonic.Unmarshal([]byte(p.Properties_update_json), &propertiesUpdate); err != nil {
		return nil, fmt.Errorf("failed to parse properties_update_json: %w", err)
	}
	body.Properties = propertiesUpdate

	// Use pointer for Archived to distinguish between false and not set
	// The parameter struct uses bool, so `false` is the default.
	// If the API requires explicit null/omit vs false, the param struct might need *bool
	// Assuming for now that setting it to false is intended if p.Archived is false.
	// Consider if a separate parameter `set_archived` is needed for PATCH semantics.
	// Let's include it always for now, as the struct field isn't omitempty.
	archivedValue := p.Archived
	body.Archived = &archivedValue

	// Icon and Cover updates are not handled by current parameters.
	// They would need additional parameters (e.g., icon_json, cover_json)
	// and logic to unmarshal into body.Icon and body.Cover.

	restParams := map[string]interface{}{
		"method":      http.MethodPatch,
		"path":        endpointUpdatePageProperties,
		"path_params": map[string]string{"page_id": p.Page_id},
		"body":        body,
	}

	return restParams, nil
}

// handleList_users constructs parameters for the REST adapter to list users.
func handleList_users(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	_, ok := params.(*List_usersParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation list_users", params)
	}

	queryParams := make(map[string]string)

	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointListUsers,
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}

	return restParams, nil
}

// handleGet_database constructs parameters for the REST adapter to get database info.
func handleGet_database(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*Get_databaseParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_database", params)
	}

	restParams := map[string]interface{}{
		"method":      http.MethodGet,
		"path":        endpointGetDatabase,
		"path_params": map[string]string{"database_id": p.Database_id},
	}

	return restParams, nil
}

// handleList_databases constructs parameters for the REST adapter to list databases using search.
func handleList_databases(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	_, ok := params.(*List_databasesParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation list_databases", params)
	}

	// Construct request body using SearchRequest struct
	body := &SearchRequest{
		Filter: &SearchFilter{
			Property: "object",
			Value:    "database",
		},
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointListDatabasesViaSearch,
		"body":   body,
	}

	return restParams, nil
}
