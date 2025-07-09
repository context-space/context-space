package zoom

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// Define constants for API paths used by handlers.
const (
	endpointGetUserInfo       = "/users/me"
	endpointCreateMeeting     = "/users/me/meetings"
	endpointListMeetings      = "/users/me/meetings"
	endpointGetMeetingDetails = "/meetings/{meetingId}"
	endpointListRecordings    = "/users/me/recordings"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDGetUserInfo       = "get_user_info"
	operationIDCreateMeeting     = "create_meeting"
	operationIDListMeetings      = "list_meetings"
	operationIDGetMeetingDetails = "get_meeting_details"
	operationIDListRecordings    = "list_recordings"
)

// Define a struct for each operation's parameters based on config.

// CreateMeetingParams defines parameters for the Create Meeting operation.
type CreateMeetingParams struct {
	// Define fields based on create_meeting parameters in config

	Topic string `mapstructure:"topic" validate:"required"` // Meeting topic.

	Type int `mapstructure:"type" validate:"required"` // Meeting type (1: Instant, 2: Scheduled, 3: Recurring no fixed time, 8: Recurring fixed time).

	StartTime string `mapstructure:"start_time" validate:"omitempty"` // Meeting start time (ISO 8601 format). Required for scheduled meetings.

	Duration int `mapstructure:"duration" validate:"omitempty"` // Meeting duration in minutes.

	Timezone string `mapstructure:"timezone" validate:"omitempty"` // Timezone for the meeting.

	Settings map[string]interface{} `mapstructure:"settings" validate:"omitempty"` // Meeting settings object.

}

// ListMeetingsParams defines parameters for the List Meetings operation.
type ListMeetingsParams struct {
	// Define fields based on list_meetings parameters in config

	Type string `mapstructure:"type" validate:"omitempty"` // The type of meeting to list (scheduled, live, upcoming). Defaults to 'live'.

	PageSize int `mapstructure:"page_size" validate:"omitempty"` // Number of results per page.

	NextPageToken string `mapstructure:"next_page_token" validate:"omitempty"` // Token for next page results.

}

// GetMeetingDetailsParams defines parameters for the Get Meeting Details operation.
type GetMeetingDetailsParams struct {
	// Define fields based on get_meeting_details parameters in config

	Meetingid string `mapstructure:"meetingId" validate:"required"` // The ID of the meeting.

}

// ListRecordingsParams defines parameters for the List Recordings operation.
type ListRecordingsParams struct {
	// Define fields based on list_recordings parameters in config

	PageSize int `mapstructure:"page_size" validate:"omitempty"` // Number of results per page.

	NextPageToken string `mapstructure:"next_page_token" validate:"omitempty"` // Token for next page results.

	From string `mapstructure:"from" validate:"omitempty"` // Start date for recordings query (YYYY-MM-DD).

	To string `mapstructure:"to" validate:"omitempty"` // End date for recordings query (YYYY-MM-DD).

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
func (a *ZoomAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler, requiredPerms []string) {
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
func (a *ZoomAdapter) registerOperations() {
	a.RegisterOperation(
		operationIDGetUserInfo,
		&struct{}{},
		handleGetUserInfo,
		[]string{"read_user_info"},
	)

	a.RegisterOperation(
		operationIDCreateMeeting,
		&CreateMeetingParams{},
		handleCreateMeeting,
		[]string{"manage_meetings"},
	)

	a.RegisterOperation(
		operationIDListMeetings,
		&ListMeetingsParams{},
		handleListMeetings,
		[]string{"read_meetings"},
	)

	a.RegisterOperation(
		operationIDGetMeetingDetails,
		&GetMeetingDetailsParams{},
		handleGetMeetingDetails,
		[]string{"read_meetings"},
	)

	a.RegisterOperation(
		operationIDListRecordings,
		&ListRecordingsParams{},
		handleListRecordings,
		[]string{"read_recordings"},
	)
}

// handleGetUserInfo constructs parameters for the REST adapter for the Get User Information operation.
func handleGetUserInfo(ctx context.Context, params interface{}) (map[string]interface{}, error) {

	// No parameters defined for this operation in the config.
	// The 'params' argument in the handler will be nil.
	// Ensure the handler logic doesn't expect a non-nil params struct.

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)
	headers := make(map[string]string)

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetUserInfo,
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

// handleCreateMeeting constructs parameters for the REST adapter for the Create Meeting operation.
func handleCreateMeeting(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var p *CreateMeetingParams
	if params != nil {
		var ok bool
		p, ok = params.(*CreateMeetingParams)
		if !ok {
			return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation create_meeting", params)
		}
	} else {
		// Parameters are required for this operation as per zoom.json (topic, type are required)
		return nil, fmt.Errorf("internal error: params are required for operation create_meeting but were nil")
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)
	headers := make(map[string]string)
	var requestBody interface{}

	if p != nil {
		requestBody = p
	}

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointCreateMeeting,
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
	if requestBody != nil {
		restParams["body"] = requestBody
	}

	return restParams, nil
}

// handleListMeetings constructs parameters for the REST adapter for the List Meetings operation.
func handleListMeetings(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var p *ListMeetingsParams
	if params != nil {
		var ok bool
		p, ok = params.(*ListMeetingsParams)
		if !ok {
			return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation list_meetings", params)
		}
	} else {
		// Params are optional, initialize to empty struct if nil
		p = &ListMeetingsParams{}
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)
	headers := make(map[string]string)

	if p != nil {
		if p.Type != "" {
			queryParams["type"] = p.Type
		}
		if p.PageSize > 0 {
			queryParams["page_size"] = strconv.Itoa(p.PageSize)
		}
		if p.NextPageToken != "" {
			queryParams["next_page_token"] = p.NextPageToken
		}
	}

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointListMeetings,
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

// handleGetMeetingDetails constructs parameters for the REST adapter for the Get Meeting Details operation.
func handleGetMeetingDetails(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var p *GetMeetingDetailsParams
	if params != nil {
		var ok bool
		p, ok = params.(*GetMeetingDetailsParams)
		if !ok {
			return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation get_meeting_details", params)
		}
	} else {
		// meetingId is a required path parameter, so params cannot be nil.
		return nil, fmt.Errorf("internal error: params are required for operation get_meeting_details but were nil")
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)
	headers := make(map[string]string)

	if p != nil {
		if p.Meetingid != "" {
			pathParams["meetingId"] = p.Meetingid
		} else {
			// This case should ideally be caught by validation if Meetingid is marked 'required:true' in params struct tag
			return nil, fmt.Errorf("meetingId is required for get_meeting_details")
		}
	}

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetMeetingDetails,
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

// handleListRecordings constructs parameters for the REST adapter for the List Recordings operation.
func handleListRecordings(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var p *ListRecordingsParams
	if params != nil {
		var ok bool
		p, ok = params.(*ListRecordingsParams)
		if !ok {
			return nil, fmt.Errorf("internal error: unexpected parameter type %T for operation list_recordings", params)
		}
	} else {
		// Params are optional, initialize to empty struct if nil
		p = &ListRecordingsParams{}
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)
	headers := make(map[string]string)

	if p != nil {
		if p.PageSize > 0 {
			queryParams["page_size"] = strconv.Itoa(p.PageSize)
		}
		if p.NextPageToken != "" {
			queryParams["next_page_token"] = p.NextPageToken
		}
		if p.From != "" {
			queryParams["from"] = p.From
		}
		if p.To != "" {
			queryParams["to"] = p.To
		}
	}

	// Construct the map to return to the REST adapter
	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointListRecordings,
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
