package slack

import (
	"context"
	"fmt"
	"net/http"
)

// Define constants for API paths used by handlers.
const (
	endpointPostMessage         = "/chat.postMessage"
	endpointReplyInThread       = "/chat.postMessage"
	endpointListThreadMessages  = "/conversations.replies"
	endpointListChannels        = "/conversations.list"
	endpointListDirectMessages  = "/conversations.list"
	endpointOpenDirectMessage   = "/conversations.open"
	endpointListChannelMembers  = "/conversations.members"
	endpointListChannelMessages = "/conversations.history"
	endpointSendDirectMessage   = "/chat.postMessage"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDPostMessage         = "post_message"
	operationIDReplyInThread       = "reply_in_thread"
	operationIDListThreadMessages  = "list_thread_messages"
	operationIDListChannels        = "list_channels"
	operationIDListDirectMessages  = "list_direct_messages"
	operationIDOpenDirectMessage   = "open_direct_message"
	operationIDListChannelMembers  = "list_channel_members"
	operationIDListChannelMessages = "list_channel_messages"
	operationIDSendDirectMessage   = "send_direct_message"
)

// PostMessageParams defines parameters for the Post Message operation.
type PostMessageParams struct {
	// Define fields based on post_message parameters in config

	Channel string `mapstructure:"channel" validate:"required"` // Conversation ID (channel, group or DM)

	Text string `mapstructure:"text" validate:"required"` // Plain-text message content

	Blocks []interface{} `mapstructure:"blocks" validate:"omitempty"` // Block Kit rich message blocks (JSON)

}

// ReplyInThreadParams defines parameters for the Reply In Thread operation.
type ReplyInThreadParams struct {
	// Define fields based on reply_in_thread parameters in config

	Channel string `mapstructure:"channel" validate:"required"` // Conversation ID that contains the thread

	ThreadTs string `mapstructure:"thread_ts" validate:"required"` // Timestamp of the parent message to reply to

	Text string `mapstructure:"text" validate:"required"` // Plain-text reply content

	ReplyBroadcast bool `mapstructure:"reply_broadcast" validate:"omitempty"` // Whether to also send the reply to the channel

	Blocks []interface{} `mapstructure:"blocks" validate:"omitempty"` // Block Kit rich message blocks (JSON)

}

// ListThreadMessagesParams defines parameters for the List Thread Messages operation.
type ListThreadMessagesParams struct {
	// Define fields based on list_thread_messages parameters in config

	Channel string `mapstructure:"channel" validate:"required"` // Conversation ID that contains the thread

	Ts string `mapstructure:"ts" validate:"required"` // Timestamp of the parent message

	Cursor string `mapstructure:"cursor" validate:"omitempty"` // Cursor for pagination

	Limit int `mapstructure:"limit" validate:"omitempty"` // Maximum number of messages to return

}

// ListChannelsParams defines parameters for the List Channels operation.
type ListChannelsParams struct {
	// Define fields based on list_channels parameters in config

	Types string `mapstructure:"types" validate:"omitempty"` // Filter conversation types (comma-separated)

	ExcludeArchived bool `mapstructure:"exclude_archived" validate:"omitempty"` // Whether to exclude archived conversations

	Cursor string `mapstructure:"cursor" validate:"omitempty"` // Cursor for pagination

	Limit int `mapstructure:"limit" validate:"omitempty"` // Maximum number of results to return

}

// ListDirectMessagesParams defines parameters for the List Direct Messages operation.
type ListDirectMessagesParams struct {
	// Define fields based on list_direct_messages parameters in config

	Cursor string `mapstructure:"cursor" validate:"omitempty"` // Cursor for pagination

	Limit int `mapstructure:"limit" validate:"omitempty"` // Maximum number of results to return

}

// OpenDirectMessageParams defines parameters for the Open Direct Message operation.
type OpenDirectMessageParams struct {
	// Define fields based on open_direct_message parameters in config

	UserId string `mapstructure:"user_id" validate:"required"` // Slack User ID to open a DM with

	ReturnIm bool `mapstructure:"return_im" validate:"omitempty"` // Set to true to return the DM conversation if it already exists

}

// ListChannelMembersParams defines parameters for the List Channel Members operation.
type ListChannelMembersParams struct {
	// Define fields based on list_channel_members parameters in config

	Channel string `mapstructure:"channel" validate:"required"` // Channel ID

	Cursor string `mapstructure:"cursor" validate:"omitempty"` // Cursor for pagination

	Limit int `mapstructure:"limit" validate:"omitempty"` // Maximum number of results to return

}

// ListChannelMessagesParams defines parameters for the List Channel Messages operation.
type ListChannelMessagesParams struct {
	// Define fields based on list_channel_messages parameters in config

	Channel string `mapstructure:"channel" validate:"required"` // Conversation ID to fetch history for

	Cursor string `mapstructure:"cursor" validate:"omitempty"` // Cursor for pagination

	Limit int `mapstructure:"limit" validate:"omitempty"` // Maximum number of messages to return

	Oldest string `mapstructure:"oldest" validate:"omitempty"` // Only include messages after this timestamp

	Latest string `mapstructure:"latest" validate:"omitempty"` // Only include messages before this timestamp

}

// SendDirectMessageParams defines parameters for the Send Direct Message operation.
type SendDirectMessageParams struct {
	// Define fields based on send_direct_message parameters in config

	UserId string `mapstructure:"user_id" validate:"required"` // Slack User ID to send message to

	Text string `mapstructure:"text" validate:"required"` // Plain-text message content

	Blocks []interface{} `mapstructure:"blocks" validate:"omitempty"` // Block Kit rich message blocks (JSON)

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
func (a *SlackAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler, requiredPerms []string) {
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
func (a *SlackAdapter) registerOperations() {

	a.RegisterOperation(
		operationIDPostMessage,
		&PostMessageParams{},
		handlePostMessage,
		[]string{"send_message"},
	)

	a.RegisterOperation(
		operationIDReplyInThread,
		&ReplyInThreadParams{},
		handleReplyInThread,
		[]string{"send_message"},
	)

	a.RegisterOperation(
		operationIDListThreadMessages,
		&ListThreadMessagesParams{},
		handleListThreadMessages,
		[]string{"read_messages"},
	)

	a.RegisterOperation(
		operationIDListChannels,
		&ListChannelsParams{},
		handleListChannels,
		[]string{"read_conversations"},
	)

	a.RegisterOperation(
		operationIDListDirectMessages,
		&ListDirectMessagesParams{},
		handleListDirectMessages,
		[]string{"read_conversations"},
	)

	a.RegisterOperation(
		operationIDOpenDirectMessage,
		&OpenDirectMessageParams{},
		handleOpenDirectMessage,
		[]string{"write_im"},
	)

	a.RegisterOperation(
		operationIDListChannelMembers,
		&ListChannelMembersParams{},
		handleListChannelMembers,
		[]string{"read_conversations"},
	)

	a.RegisterOperation(
		operationIDListChannelMessages,
		&ListChannelMessagesParams{},
		handleListChannelMessages,
		[]string{"read_messages"},
	)

	a.RegisterOperation(
		operationIDSendDirectMessage,
		&SendDirectMessageParams{},
		handleSendDirectMessage,
		[]string{"write_im", "send_message"},
	)

}

// handlePostMessage constructs parameters for the REST adapter for the Post Message operation.
func handlePostMessage(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*PostMessageParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for post_message")
	}
	requestBody := map[string]interface{}{
		"channel": p.Channel,
		"text":    p.Text,
	}
	if len(p.Blocks) > 0 {
		requestBody["blocks"] = p.Blocks
	}
	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointPostMessage,
		"body":   requestBody,
	}
	return restParams, nil
}

// handleReplyInThread constructs parameters for the REST adapter for the Reply In Thread operation.
func handleReplyInThread(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ReplyInThreadParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for reply_in_thread")
	}
	requestBody := map[string]interface{}{
		"channel":   p.Channel,
		"text":      p.Text,
		"thread_ts": p.ThreadTs,
	}
	if len(p.Blocks) > 0 {
		requestBody["blocks"] = p.Blocks
	}
	if p.ReplyBroadcast {
		requestBody["reply_broadcast"] = p.ReplyBroadcast
	}
	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointReplyInThread,
		"body":   requestBody,
	}
	return restParams, nil
}

// handleListThreadMessages constructs parameters for the REST adapter for the List Thread Messages operation.
func handleListThreadMessages(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ListThreadMessagesParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for list_thread_messages")
	}
	queryParams := map[string]string{
		"channel": p.Channel,
		"ts":      p.Ts,
	}
	if p.Cursor != "" {
		queryParams["cursor"] = p.Cursor
	}
	if p.Limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", p.Limit)
	}
	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointListThreadMessages,
		"query_params": queryParams,
	}
	return restParams, nil
}

// handleListChannels constructs parameters for the REST adapter for the List Channels operation.
func handleListChannels(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ListChannelsParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for list_channels")
	}
	queryParams := map[string]string{}
	if p.Types != "" {
		queryParams["types"] = p.Types
	}
	if p.ExcludeArchived {
		queryParams["exclude_archived"] = "true"
	}
	if p.Cursor != "" {
		queryParams["cursor"] = p.Cursor
	}
	if p.Limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", p.Limit)
	}
	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointListChannels,
		"query_params": queryParams,
	}
	return restParams, nil
}

// handleListDirectMessages constructs parameters for the REST adapter for the List Direct Messages operation.
func handleListDirectMessages(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ListDirectMessagesParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for list_direct_messages")
	}
	queryParams := map[string]string{
		"types": "im",
	}
	if p.Cursor != "" {
		queryParams["cursor"] = p.Cursor
	}
	if p.Limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", p.Limit)
	}
	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointListDirectMessages,
		"query_params": queryParams,
	}
	return restParams, nil
}

// handleOpenDirectMessage constructs parameters for the REST adapter for the Open Direct Message operation.
func handleOpenDirectMessage(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*OpenDirectMessageParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for open_direct_message")
	}
	requestBody := map[string]interface{}{
		"users": p.UserId,
	}
	if p.ReturnIm {
		requestBody["return_im"] = true
	}
	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointOpenDirectMessage,
		"body":   requestBody,
	}
	return restParams, nil
}

// handleListChannelMembers constructs parameters for the REST adapter for the List Channel Members operation.
func handleListChannelMembers(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ListChannelMembersParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for list_channel_members")
	}
	queryParams := map[string]string{
		"channel": p.Channel,
	}
	if p.Cursor != "" {
		queryParams["cursor"] = p.Cursor
	}
	if p.Limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", p.Limit)
	}
	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointListChannelMembers,
		"query_params": queryParams,
	}
	return restParams, nil
}

// handleListChannelMessages constructs parameters for the REST adapter for the List Channel Messages operation.
func handleListChannelMessages(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*ListChannelMessagesParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for list_channel_messages")
	}
	queryParams := map[string]string{
		"channel": p.Channel,
	}
	if p.Cursor != "" {
		queryParams["cursor"] = p.Cursor
	}
	if p.Limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", p.Limit)
	}
	if p.Oldest != "" {
		queryParams["oldest"] = p.Oldest
	}
	if p.Latest != "" {
		queryParams["latest"] = p.Latest
	}
	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointListChannelMessages,
		"query_params": queryParams,
	}
	return restParams, nil
}

// handleSendDirectMessage constructs parameters for the REST adapter for the Send Direct Message operation.
func handleSendDirectMessage(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*SendDirectMessageParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for send_direct_message")
	}
	// To send a DM, you need the channel ID (from conversations.open), but here we assume user_id is provided and the logic elsewhere will open the DM and get the channel ID.
	requestBody := map[string]interface{}{
		// "channel": <should be set by the caller after opening DM>,
		"text": p.Text,
	}
	if len(p.Blocks) > 0 {
		requestBody["blocks"] = p.Blocks
	}
	// Note: The actual sending of a DM requires the channel ID, not user_id. This handler may need to be adapted depending on how the integration is structured.
	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointSendDirectMessage,
		"body":   requestBody,
	}
	return restParams, nil
}
