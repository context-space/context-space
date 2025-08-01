package spotify

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// Define constants for API paths used by handlers.
const (
	endpointCreatePlaylist          = "/users/{user_id}/playlists"
	endpointAddTracksToPlaylist     = "/playlists/{playlist_id}/tracks"
	endpointGetPlaylistTracks       = "/playlists/{playlist_id}/tracks"
	endpointRemovePlaylistTracks    = "/playlists/{playlist_id}/tracks"
	endpointGetCurrentUserPlaylists = "/me/playlists"
	endpointSearch                  = "/search"
	endpointStartPlayback           = "/me/player/play"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDCreatePlaylist          = "create_playlist"
	operationIDAddTracksToPlaylist     = "add_tracks_to_playlist"
	operationIDGetPlaylistTracks       = "get_playlist_tracks"
	operationIDRemovePlaylistTracks    = "remove_playlist_tracks"
	operationIDGetCurrentUserPlaylists = "get_current_user_playlists"
	operationIDSearch                  = "search"
	operationIDStartPlayback           = "start_playback"
)

// Define a struct for each operation's parameters based on config.

// CreatePlaylistParams defines parameters for the Create Playlist operation.
type CreatePlaylistParams struct {
	// Define fields based on create_playlist parameters in config

	UserId string `mapstructure:"user_id" validate:"required"` // The user's Spotify user ID.

	Name string `mapstructure:"name" validate:"required"` // The name for the new playlist.

	Public bool `mapstructure:"public" validate:"omitempty"` // Defaults to true. If true the playlist will be public, if false it will be private.

	Collaborative bool `mapstructure:"collaborative" validate:"omitempty"` // Defaults to false. If true the playlist will be collaborative.

	Description string `mapstructure:"description" validate:"omitempty"` // Value for playlist description as displayed in Spotify Clients and in the Web API.

}

// AddTracksToPlaylistParams defines parameters for the Add Tracks to Playlist operation.
type AddTracksToPlaylistParams struct {
	// Define fields based on add_tracks_to_playlist parameters in config

	PlaylistId string `mapstructure:"playlist_id" validate:"required"` // The Spotify ID of the playlist.

	Uris []interface{} `mapstructure:"uris" validate:"required"` // A JSON array of the Spotify track URIs to add.

	Position int `mapstructure:"position" validate:"omitempty"` // The position to insert the tracks, a zero-based index.

}

// GetPlaylistTracksParams defines parameters for the Get Playlist's Tracks operation.
type GetPlaylistTracksParams struct {
	// Define fields based on get_playlist_tracks parameters in config

	PlaylistId string `mapstructure:"playlist_id" validate:"required"` // The Spotify ID of the playlist.

	Fields string `mapstructure:"fields" validate:"omitempty"` // A comma-separated list of fields to return.

	Limit int `mapstructure:"limit" validate:"omitempty"` // The maximum number of items to return. Default: 100. Minimum: 1. Maximum: 100.

	Offset int `mapstructure:"offset" validate:"omitempty"` // The index of the first item to return. Default: 0 (the first object).

	Market string `mapstructure:"market" validate:"omitempty"` // An ISO 3166-1 alpha-2 country code.

}

// RemovePlaylistTracksParams defines parameters for the Remove Playlist Tracks operation.
type RemovePlaylistTracksParams struct {
	// Define fields based on remove_playlist_tracks parameters in config

	PlaylistId string `mapstructure:"playlist_id" validate:"required"` // The Spotify ID of the playlist.

	Tracks []interface{} `mapstructure:"tracks" validate:"required"` // An array of objects containing Spotify URIs of the tracks to remove.

	SnapshotId string `mapstructure:"snapshot_id" validate:"omitempty"` // The playlist's snapshot ID.

}

// GetCurrentUserPlaylistsParams defines parameters for the Get Current User's Playlists operation.
type GetCurrentUserPlaylistsParams struct {
	Limit  int `mapstructure:"limit" validate:"omitempty,min=1,max=50"`      // The maximum number of items to return (Default: 20, Min: 1, Max: 50)
	Offset int `mapstructure:"offset" validate:"omitempty,min=0,max=100000"` // The index of the first playlist to return (Default: 0, Min: 0, Max: 100,000)
}

// SearchParams defines parameters for the Search operation.
type SearchParams struct {
	Q      string `mapstructure:"q" validate:"required"`              // Search query keywords
	Type   string `mapstructure:"type" validate:"required"`           // Types to search (comma-separated: track,artist,album,playlist)
	Limit  int    `mapstructure:"limit" validate:"omitempty,min=1,max=50"` // Maximum number of results (1-50, default 20)
	Offset int    `mapstructure:"offset" validate:"omitempty,min=0"`       // The index of the first result to return
	Market string `mapstructure:"market" validate:"omitempty"`        // ISO 3166-1 alpha-2 country code
}

// StartPlaybackParams defines parameters for the Start/Resume Playback operation.
type StartPlaybackParams struct {
	DeviceId   string      `mapstructure:"device_id" validate:"omitempty"`   // The id of the device this command is targeting
	ContextUri string      `mapstructure:"context_uri" validate:"omitempty"` // Spotify URI of the context to play
	Uris       []string    `mapstructure:"uris" validate:"omitempty"`        // Array of Spotify track URIs to play
	Offset     interface{} `mapstructure:"offset" validate:"omitempty"`      // Indicates from where in the context playback should start
	PositionMs int         `mapstructure:"position_ms" validate:"omitempty"` // The position in milliseconds to start playback
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
func (a *SpotifyAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler, requiredPerms []string) {
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
func (a *SpotifyAdapter) registerOperations() {
	a.RegisterOperation(
		operationIDCreatePlaylist,
		&CreatePlaylistParams{},
		handleCreatePlaylist,
		[]string{"modify_public_playlists", "modify_private_playlists"},
	)

	a.RegisterOperation(
		operationIDAddTracksToPlaylist,
		&AddTracksToPlaylistParams{},
		handleAddTracksToPlaylist,
		[]string{"modify_public_playlists", "modify_private_playlists"},
	)

	a.RegisterOperation(
		operationIDGetPlaylistTracks,
		&GetPlaylistTracksParams{},
		handleGetPlaylistTracks,
		[]string{"read_private_playlists"},
	)

	a.RegisterOperation(
		operationIDRemovePlaylistTracks,
		&RemovePlaylistTracksParams{},
		handleRemovePlaylistTracks,
		[]string{"modify_public_playlists", "modify_private_playlists"},
	)

	a.RegisterOperation(
		operationIDGetCurrentUserPlaylists,
		&GetCurrentUserPlaylistsParams{},
		handleGetCurrentUserPlaylists,
		[]string{"read_private_playlists"},
	)

	a.RegisterOperation(
		operationIDSearch,
		&SearchParams{},
		handleSearch,
		[]string{}, // No special permissions needed for search
	)

	a.RegisterOperation(
		operationIDStartPlayback,
		&StartPlaybackParams{},
		handleStartPlayback,
		[]string{"modify_playback_state"},
	)

}

// Define a handler function for each operation.

// handleCreatePlaylist constructs parameters for the REST adapter for the Create Playlist operation.
func handleCreatePlaylist(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*CreatePlaylistParams)
	if !ok || p == nil { // Also check if p is nil
		return nil, fmt.Errorf("invalid or missing parameters for create_playlist")
	}

	pathParams := make(map[string]string)

	requestBodyData := make(map[string]interface{})

	// Populate pathParams
	if p.UserId != "" { // Already validated by 'required' tag, but good practice
		pathParams["user_id"] = p.UserId
	}

	// Populate requestBodyData
	if p.Name != "" { // Already validated by 'required' tag
		requestBodyData["name"] = p.Name
	}
	// omitempty for boolean means if it's the zero value (false), it might be omitted by some JSON encoders.
	// Spotify API might expect 'false' explicitly or might default if not present.
	// Assuming the API defaults correctly if not present for optional booleans.
	// If explicit false is needed, then add: requestBodyData["public"] = p.Public
	if p.Public { // Only add if true, if false, let API default (or ensure explicit false if API requires)
		requestBodyData["public"] = p.Public
	}
	if p.Collaborative {
		requestBodyData["collaborative"] = p.Collaborative
	}
	if p.Description != "" {
		requestBodyData["description"] = p.Description
	}

	var requestBody interface{}
	if len(requestBodyData) > 0 {
		requestBody = requestBodyData
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointCreatePlaylist,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}

	if requestBody != nil {
		restParams["body"] = requestBody
	}

	return restParams, nil
}

// handleAddTracksToPlaylist constructs parameters for the REST adapter for the Add Tracks to Playlist operation.
func handleAddTracksToPlaylist(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*AddTracksToPlaylistParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for add_tracks_to_playlist")
	}

	pathParams := make(map[string]string)

	requestBodyData := make(map[string]interface{})

	if p.PlaylistId != "" {
		pathParams["playlist_id"] = p.PlaylistId
	}

	if len(p.Uris) > 0 { // Uris is required
		requestBodyData["uris"] = p.Uris
	}
	// For optional int, 0 is its zero value. API might interpret 0 as a valid position.
	// If 0 is not a valid position or should be omitted if not set, more logic is needed.
	// Assuming if Position is 0, it's either a valid desired position or API handles default if not specified.
	// The 'omitempty' on mapstructure tag suggests it's fine to omit if zero.
	// However, the request body might still need the field if it's truly optional at API level and not just at input.
	// Spotify docs usually clarify if 0 is a special value or if omission implies a default (e.g., append).
	// For now, including if not zero, assuming 0 is a valid position. If not, an explicit check `if p.Position != 0` might be needed.
	// Update: Spotify's Add Items to Playlist expects "position" to be an integer, so 0 is valid.
	// If it were optional and omitting it has a different meaning than sending 0, then `if p.Position != 0 || fieldIsSet(p, "Position")`
	// But with mapstructure and omitempty, if it's not in the input, it will be 0.
	// So, we will always send it if the API expects it. If it was a pointer *int, we could check for nil.
	requestBodyData["position"] = p.Position // Send 0 if that's what p.Position is.

	var requestBody interface{}
	if len(requestBodyData) > 0 {
		requestBody = requestBodyData
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointAddTracksToPlaylist,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}

	if requestBody != nil {
		restParams["body"] = requestBody
	}

	return restParams, nil
}

// handleGetPlaylistTracks constructs parameters for the REST adapter for the Get Playlist's Tracks operation.
func handleGetPlaylistTracks(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*GetPlaylistTracksParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for get_playlist_tracks")
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	if p.PlaylistId != "" {
		pathParams["playlist_id"] = p.PlaylistId
	}

	if p.Fields != "" {
		queryParams["fields"] = p.Fields
	}
	if p.Limit > 0 { // Spotify API: Default 100, Min 1, Max 100. Let API handle default if 0.
		queryParams["limit"] = strconv.Itoa(p.Limit)
	}
	if p.Offset >= 0 { // Spotify API: Default 0.
		queryParams["offset"] = strconv.Itoa(p.Offset)
	}
	if p.Market != "" {
		queryParams["market"] = p.Market
	}

	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetPlaylistTracks,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}
	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}

	return restParams, nil
}

// handleRemovePlaylistTracks constructs parameters for the REST adapter for the Remove Playlist Tracks operation.
func handleRemovePlaylistTracks(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*RemovePlaylistTracksParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for remove_playlist_tracks")
	}

	pathParams := make(map[string]string)
	requestBodyData := make(map[string]interface{})

	if p.PlaylistId != "" {
		pathParams["playlist_id"] = p.PlaylistId
	}

	if len(p.Tracks) > 0 { // Tracks is required
		// Ensure tracks are in the format: {"tracks": [{"uri": "spotify:track:xxxx"}, ...]}
		requestBodyData["tracks"] = p.Tracks
	}
	if p.SnapshotId != "" {
		requestBodyData["snapshot_id"] = p.SnapshotId
	}

	var requestBody interface{}
	if len(requestBodyData) > 0 {
		requestBody = requestBodyData
	}

	restParams := map[string]interface{}{
		"method": http.MethodDelete,
		"path":   endpointRemovePlaylistTracks,
	}

	if len(pathParams) > 0 {
		restParams["path_params"] = pathParams
	}

	if requestBody != nil {
		restParams["body"] = requestBody
	}

	return restParams, nil
}

// handleGetCurrentUserPlaylists constructs parameters for the REST adapter for the Get Current User's Playlists operation.
func handleGetCurrentUserPlaylists(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*GetCurrentUserPlaylistsParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for get_current_user_playlists")
	}

	// Validate parameter ranges according to Spotify API limits
	if p.Limit != 0 && (p.Limit < 1 || p.Limit > 50) {
		return nil, fmt.Errorf("limit parameter must be between 1 and 50 (or 0 to use default 20), got %d", p.Limit)
	}
	if p.Offset < 0 || p.Offset > 100000 {
		return nil, fmt.Errorf("offset parameter must be between 0 and 100,000, got %d", p.Offset)
	}

	queryParams := make(map[string]string)

	if p.Limit > 0 { // Spotify API: Default 20, Min 1, Max 50. Let API handle default if 0.
		queryParams["limit"] = strconv.Itoa(p.Limit)
	}
	if p.Offset >= 0 { // Spotify API: Default 0. 0 is a valid starting index.
		queryParams["offset"] = strconv.Itoa(p.Offset)
	}

	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointGetCurrentUserPlaylists,
	}

	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}

	return restParams, nil
}

// handleSearch constructs parameters for the REST adapter for the Search operation.
func handleSearch(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*SearchParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for search")
	}

	queryParams := make(map[string]string)

	// Required parameters
	queryParams["q"] = p.Q
	queryParams["type"] = p.Type

	// Optional parameters
	if p.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(p.Limit)
	}
	if p.Offset > 0 {
		queryParams["offset"] = strconv.Itoa(p.Offset)
	}
	if p.Market != "" {
		queryParams["market"] = p.Market
	}

	restParams := map[string]interface{}{
		"method": http.MethodGet,
		"path":   endpointSearch,
	}

	if len(queryParams) > 0 {
		restParams["query_params"] = queryParams
	}

	return restParams, nil
}

// handleStartPlayback constructs parameters for the REST adapter for the Start/Resume Playback operation.
func handleStartPlayback(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*StartPlaybackParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for start_playback")
	}

	restParams := map[string]interface{}{
		"method": http.MethodPut,
		"path":   endpointStartPlayback,
	}

	// Add device_id as query parameter if provided
	if p.DeviceId != "" {
		restParams["query_params"] = map[string]string{
			"device_id": p.DeviceId,
		}
	}

	// Build request body
	requestBody := make(map[string]interface{})

	// At least one of context_uri or uris should be provided
	if p.ContextUri != "" {
		requestBody["context_uri"] = p.ContextUri
	}
	if len(p.Uris) > 0 {
		requestBody["uris"] = p.Uris
	}
	if p.Offset != nil {
		requestBody["offset"] = p.Offset
	}
	if p.PositionMs > 0 {
		requestBody["position_ms"] = p.PositionMs
	}

	// Only add body if there's content
	if len(requestBody) > 0 {
		restParams["body"] = requestBody
	}

	return restParams, nil
}
