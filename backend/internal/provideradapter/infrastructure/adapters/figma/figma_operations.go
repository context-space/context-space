package figma

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// Define constants for API paths used by handlers.
const (
	endpointGetMe                = "/me"
	endpointGetTeamStyles        = "/teams/{team_id}/styles"
	endpointGetFileNodes         = "/files/{file_key}/nodes"
	endpointGetFile              = "/files/{file_key}"
	endpointGetImageRenders      = "/images/{file_key}"
	endpointGetTeamComponents    = "/teams/{team_id}/components"
	endpointGetTeamComponentSets = "/teams/{team_id}/component_sets"
	endpointGetFileComponentSets = "/files/{file_key}/component_sets"
	endpointGetFileStyles        = "/files/{file_key}/styles"
	endpointGetFileComponents    = "/files/{file_key}/components"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDGetMe                = "get_me"
	operationIDGetTeamStyles        = "get_team_styles"
	operationIDGetFileNodes         = "get_file_nodes"
	operationIDGetFile              = "get_file"
	operationIDGetImageRenders      = "get_image_renders"
	operationIDGetTeamComponents    = "get_team_components"
	operationIDGetTeamComponentSets = "get_team_component_sets"
	operationIDGetFileComponentSets = "get_file_component_sets"
	operationIDGetFileStyles        = "get_file_styles"
	operationIDGetFileComponents    = "get_file_components"
)

// Define a struct for each operation's parameters based on config.

// GetTeamStylesParams defines parameters for the Get Team Styles operation.
type GetTeamStylesParams struct {
	// Define fields based on get_team_styles parameters in config

	TeamId string `mapstructure:"team_id" validate:"required"` // ID of the team.

	PageSize int `mapstructure:"page_size" validate:"omitempty"` // Number of items per page.

	After int `mapstructure:"after" validate:"omitempty"` // Cursor for pagination.

}

// GetFileNodesParams defines parameters for the Get File Nodes operation.
type GetFileNodesParams struct {
	// Define fields based on get_file_nodes parameters in config

	FileKey string `mapstructure:"file_key" validate:"required"` // Key of the file.

	Ids string `mapstructure:"ids" validate:"required"` // Comma-separated list of node IDs to retrieve.

	Depth int `mapstructure:"depth" validate:"omitempty"` // Depth of the node tree to return.

	Geometry string `mapstructure:"geometry" validate:"omitempty"` // Include geometry data ('paths' or 'bounds').

}

// GetFileParams defines parameters for the Get File operation.
type GetFileParams struct {
	// Define fields based on get_file parameters in config

	FileKey string `mapstructure:"file_key" validate:"required"` // Key of the file.

	Geometry string `mapstructure:"geometry" validate:"omitempty"` // Include geometry data ('paths' or 'bounds').

	Version string `mapstructure:"version" validate:"omitempty"` // Get a specific version of the file.

	PluginData string `mapstructure:"plugin_data" validate:"omitempty"` // Specifies a list of plugin IDs to include in the response.
	BranchData bool   `mapstructure:"branch_data" validate:"omitempty"` // If true, the file's branches will be included in the response.
}

// GetImageRendersParams defines parameters for the Get Image Renders operation.
type GetImageRendersParams struct {
	// Define fields based on get_image_renders parameters in config

	FileKey string `mapstructure:"file_key" validate:"required"` // Key of the file.

	Ids string `mapstructure:"ids" validate:"required"` // Comma-separated list of node IDs to render.

	Scale float64 `mapstructure:"scale" validate:"omitempty"` // Image scale (0.01 to 4).

	Format string `mapstructure:"format" validate:"omitempty,oneof=jpg png svg pdf"` // Image format ('jpg', 'png', 'svg', 'pdf').

	Version string `mapstructure:"version" validate:"omitempty"` // Render nodes from a specific version.

	UseAbsoluteBounds bool `mapstructure:"use_absolute_bounds" validate:"omitempty"` // Use absolute bounds for rendering.

}

// GetTeamComponentsParams defines parameters for the Get Team Components operation.
type GetTeamComponentsParams struct {
	// Define fields based on get_team_components parameters in config

	TeamId string `mapstructure:"team_id" validate:"required"` // ID of the team.

	PageSize int `mapstructure:"page_size" validate:"omitempty"` // Number of items per page.

	After int `mapstructure:"after" validate:"omitempty"` // Cursor for pagination.

}

// GetTeamComponentSetsParams defines parameters for the Get Team Component Sets operation.
type GetTeamComponentSetsParams struct {
	// Define fields based on get_team_component_sets parameters in config

	TeamId string `mapstructure:"team_id" validate:"required"` // ID of the team.

	PageSize int `mapstructure:"page_size" validate:"omitempty"` // Number of items per page.

	After int `mapstructure:"after" validate:"omitempty"` // Cursor for pagination.

}

// GetFileComponentSetsParams defines parameters for the Get File Component Sets operation.
type GetFileComponentSetsParams struct {
	// Define fields based on get_file_component_sets parameters in config

	FileKey string `mapstructure:"file_key" validate:"required"` // Key of the file.

}

// GetFileStylesParams defines parameters for the Get File Styles operation.
type GetFileStylesParams struct {
	// Define fields based on get_file_styles parameters in config

	FileKey string `mapstructure:"file_key" validate:"required"` // Key of the file.

}

// GetFileComponentsParams defines parameters for the Get File Components operation.
type GetFileComponentsParams struct {
	// Define fields based on get_file_components parameters in config

	FileKey string `mapstructure:"file_key" validate:"required"` // Key of the file.

}

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
func (a *FigmaAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler, requiredPerms []string) {
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
func (a *FigmaAdapter) registerOperations() {
	a.RegisterOperation(
		operationIDGetMe,
		&struct{}{},
		handleGetMe,
		[]string{},
	)

	a.RegisterOperation(
		operationIDGetTeamStyles,
		&GetTeamStylesParams{},
		handleGetTeamStyles,
		[]string{},
	)

	a.RegisterOperation(
		operationIDGetFileNodes,
		&GetFileNodesParams{},
		handleGetFileNodes,
		[]string{},
	)

	a.RegisterOperation(
		operationIDGetFile,
		&GetFileParams{},
		handleGetFile,
		[]string{},
	)

	a.RegisterOperation(
		operationIDGetImageRenders,
		&GetImageRendersParams{},
		handleGetImageRenders,
		[]string{},
	)

	a.RegisterOperation(
		operationIDGetTeamComponents,
		&GetTeamComponentsParams{},
		handleGetTeamComponents,
		[]string{},
	)

	a.RegisterOperation(
		operationIDGetTeamComponentSets,
		&GetTeamComponentSetsParams{},
		handleGetTeamComponentSets,
		[]string{},
	)

	a.RegisterOperation(
		operationIDGetFileComponentSets,
		&GetFileComponentSetsParams{},
		handleGetFileComponentSets,
		[]string{},
	)

	a.RegisterOperation(
		operationIDGetFileStyles,
		&GetFileStylesParams{},
		handleGetFileStyles,
		[]string{},
	)

	a.RegisterOperation(
		operationIDGetFileComponents,
		&GetFileComponentsParams{},
		handleGetFileComponents,
		[]string{},
	)
}

// handleGetMe constructs parameters for the REST adapter for the Get Authenticated User operation.
func handleGetMe(ctx context.Context, params interface{}) (map[string]interface{}, error) {

	// No parameters defined for this operation in the config.
	// The 'params' argument in the handler will be nil.
	// Ensure the handler logic doesn't expect a non-nil params struct.

	// pathParams, queryParams, headers, requestBody are not needed as per figma.json
	// No specific path parameters, query parameters, or body for /me endpoint.

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		// Use the pre-calculated Go constant string
		"method": http.MethodGet,
		"path":   endpointGetMe, // Use generated constant
	}

	return restParams, nil
}

// handleGetTeamStyles constructs parameters for the REST adapter for the Get Team Styles operation.
func handleGetTeamStyles(ctx context.Context, params interface{}) (map[string]interface{}, error) {

	// Cast params if they exist
	p, ok := params.(*GetTeamStylesParams)
	if !ok {
		// Handle case where params is expected but nil or wrong type
		if params == nil {
			// This operation requires team_id, so params should not be nil.
			// Validation should catch this earlier if team_id is marked as required.
			return nil, fmt.Errorf("internal error: parameters are required for operation get_team_styles")
		}
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_team_styles", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	// Map parameters from p to REST parameters
	// Path parameters
	if p.TeamId == "" { // Already validated by struct tag, but good practice to check
		return nil, fmt.Errorf("team_id is required for get_team_styles")
	}
	pathParams["team_id"] = p.TeamId

	// Query parameters
	if p.PageSize > 0 { // Figma API might treat 0 as default or ignore it. Explicitly check for > 0.
		queryParams["page_size"] = strconv.Itoa(p.PageSize)
	}
	if p.After > 0 { // Similar to PageSize, Figma's 'after' is an integer cursor.
		queryParams["after"] = strconv.Itoa(p.After)
	}

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		// Use the pre-calculated Go constant string
		"method": http.MethodGet,
		"path":   endpointGetTeamStyles, // Use generated constant
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}
	return restParams, nil
}

// handleGetFileNodes constructs parameters for the REST adapter for the Get File Nodes operation.
func handleGetFileNodes(ctx context.Context, params interface{}) (map[string]interface{}, error) {

	// Cast params if they exist
	p, ok := params.(*GetFileNodesParams)
	if !ok {
		// Handle case where params is expected but nil or wrong type
		if params == nil {
			return nil, fmt.Errorf("internal error: parameters are required for operation get_file_nodes")
		}
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_file_nodes", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	// Path parameters
	if p.FileKey == "" {
		return nil, fmt.Errorf("file_key is required for get_file_nodes")
	}
	pathParams["file_key"] = p.FileKey

	// Query parameters
	if p.Ids == "" {
		return nil, fmt.Errorf("ids parameter is required for get_file_nodes")
	}
	queryParams["ids"] = p.Ids

	if p.Depth > 0 { // Figma API: A specific depth to traverse to. Defaults to 1, meaning only direct children of the specified nodes.
		queryParams["depth"] = strconv.Itoa(p.Depth)
	}
	if p.Geometry != "" {
		queryParams["geometry"] = p.Geometry
	}

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetFileNodes, // Use generated constant
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}

	return restParams, nil
}

// handleGetFile constructs parameters for the REST adapter for the Get File operation.
func handleGetFile(ctx context.Context, params interface{}) (map[string]interface{}, error) {

	// Cast params if they exist
	p, ok := params.(*GetFileParams)
	if !ok {
		// Handle case where params is expected but nil or wrong type
		if params == nil {
			return nil, fmt.Errorf("internal error: parameters are required for operation get_file")
		}
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_file", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	// Path parameters
	if p.FileKey == "" {
		return nil, fmt.Errorf("file_key is required for get_file")
	}
	pathParams["file_key"] = p.FileKey

	// Query parameters
	if p.Geometry != "" {
		queryParams["geometry"] = p.Geometry
	}
	if p.Version != "" {
		queryParams["version"] = p.Version
	}
	if p.PluginData != "" {
		queryParams["plugin_data"] = p.PluginData
	}
	if p.BranchData {
		queryParams["branch_data"] = "true" // Figma API typically expects "true" as a string for boolean flags
	}

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetFile, // Use generated constant
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}

	return restParams, nil
}

// handleGetImageRenders constructs parameters for the REST adapter for the Get Image Renders operation.
func handleGetImageRenders(ctx context.Context, params interface{}) (map[string]interface{}, error) {

	// Cast params if they exist
	p, ok := params.(*GetImageRendersParams)
	if !ok {
		// Handle case where params is expected but nil or wrong type
		if params == nil {
			return nil, fmt.Errorf("internal error: parameters are required for operation get_image_renders")
		}
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_image_renders", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	// Path parameters
	if p.FileKey == "" {
		return nil, fmt.Errorf("file_key is required for get_image_renders")
	}
	pathParams["file_key"] = p.FileKey

	// Query parameters
	if p.Ids == "" {
		return nil, fmt.Errorf("ids parameter is required for get_image_renders")
	}
	queryParams["ids"] = p.Ids

	if p.Scale != 0 { // Default is 1 according to Figma docs, 0 might be invalid or treated as default. Consider omitempty behavior.
		// Figma API specifies scale: number between 0.01 and 4. Default: 1.
		// Let's assume mapstructure correctly sets default if not provided, or validation catches it.
		// If p.Scale is its zero value (0.0) and it's not omitempty, it means it was explicitly set to 0.
		// We should only send it if it's a valid Figma API value.
		// For simplicity here, if it's not the Go zero value, we send it.
		queryParams["scale"] = strconv.FormatFloat(p.Scale, 'f', -1, 64)
	}

	if p.Format != "" { // Already validated by oneof, but good practice.
		queryParams["format"] = p.Format
	}

	if p.Version != "" {
		queryParams["version"] = p.Version
	}

	// use_absolute_bounds is a boolean. If omitempty and false, it won't be included.
	// If it's present (e.g. explicitly false or true), we add it.
	// The mapstructure tag `omitempty` means if `UseAbsoluteBounds` is `false` (the zero value for bool), it won't be set from input unless explicitly provided.
	// We should only send it if it's explicitly provided as true OR false in the parameters. The struct field isn't a pointer.
	// However, the Figma API docs imply it's a boolean query param. If not present, it defaults to false.
	// So, we only need to send it if it's true.
	if p.UseAbsoluteBounds {
		queryParams["use_absolute_bounds"] = strconv.FormatBool(p.UseAbsoluteBounds)
	}

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetImageRenders, // Use generated constant
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}

	return restParams, nil
}

// handleGetTeamComponents constructs parameters for the REST adapter for the Get Team Components operation.
func handleGetTeamComponents(ctx context.Context, params interface{}) (map[string]interface{}, error) {

	// Cast params if they exist
	p, ok := params.(*GetTeamComponentsParams)
	if !ok {
		if params == nil {
			return nil, fmt.Errorf("internal error: parameters are required for operation get_team_components")
		}
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_team_components", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	// Path parameters
	if p.TeamId == "" {
		return nil, fmt.Errorf("team_id is required for get_team_components")
	}
	pathParams["team_id"] = p.TeamId

	// Query parameters
	if p.PageSize > 0 {
		queryParams["page_size"] = strconv.Itoa(p.PageSize)
	}
	if p.After > 0 {
		queryParams["after"] = strconv.Itoa(p.After)
	}

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetTeamComponents, // Use generated constant
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}

	return restParams, nil
}

// handleGetTeamComponentSets constructs parameters for the REST adapter for the Get Team Component Sets operation.
func handleGetTeamComponentSets(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*GetTeamComponentSetsParams)
	if !ok {
		if params == nil {
			return nil, fmt.Errorf("internal error: parameters are required for operation get_team_component_sets")
		}
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_team_component_sets", params)
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	// Path parameters
	if p.TeamId == "" {
		return nil, fmt.Errorf("team_id is required for get_team_component_sets")
	}
	pathParams["team_id"] = p.TeamId

	// Query parameters
	if p.PageSize > 0 {
		queryParams["page_size"] = strconv.Itoa(p.PageSize)
	}
	if p.After > 0 {
		queryParams["after"] = strconv.Itoa(p.After)
	}

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetTeamComponentSets, // Use generated constant
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}

	return restParams, nil
}

// handleGetFileComponentSets constructs parameters for the REST adapter for the Get File Component Sets operation.
func handleGetFileComponentSets(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*GetFileComponentSetsParams)
	if !ok {
		if params == nil {
			return nil, fmt.Errorf("internal error: parameters are required for operation get_file_component_sets")
		}
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_file_component_sets", params)
	}

	pathParams := make(map[string]string)

	// Path parameters
	if p.FileKey == "" {
		return nil, fmt.Errorf("file_key is required for get_file_component_sets")
	}
	pathParams["file_key"] = p.FileKey

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetFileComponentSets,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}

	return restParams, nil
}

// handleGetFileStyles constructs parameters for the REST adapter for the Get File Styles operation.
func handleGetFileStyles(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*GetFileStylesParams)
	if !ok {
		if params == nil {
			return nil, fmt.Errorf("internal error: parameters are required for operation get_file_styles")
		}
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_file_styles", params)
	}

	pathParams := make(map[string]string)

	// Path parameters
	if p.FileKey == "" {
		return nil, fmt.Errorf("file_key is required for get_file_styles")
	}
	pathParams["file_key"] = p.FileKey

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetFileStyles,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}

	return restParams, nil
}

// handleGetFileComponents constructs parameters for the REST adapter for the Get File Components operation.
func handleGetFileComponents(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*GetFileComponentsParams)
	if !ok {
		if params == nil {
			return nil, fmt.Errorf("internal error: parameters are required for operation get_file_components")
		}
		return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_file_components", params)
	}

	pathParams := make(map[string]string)

	// Path parameters
	if p.FileKey == "" {
		return nil, fmt.Errorf("file_key is required for get_file_components")
	}
	pathParams["file_key"] = p.FileKey

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetFileComponents,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}

	return restParams, nil
}
