package tmdb

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// Define operation IDs as constants
const (
	// Genre operations
	operationIDGenreMovieList = "genre_movie_list"
	operationIDGenreTVList    = "genre_tv_list"

	// Movie operations
	operationIDMoviePopular    = "movie_popular"
	operationIDMovieNowPlaying = "movie_now_playing"
	operationIDMovieTopRated   = "movie_top_rated"
	operationIDMovieUpcoming   = "movie_upcoming"

	// TV operations
	operationIDTVPopular     = "tv_popular"
	operationIDTVOnTheAir    = "tv_on_the_air"
	operationIDTVTopRated    = "tv_top_rated"
	operationIDTVAiringToday = "tv_airing_today"

	// Discover operations
	operationIDDiscoverMovie = "discover_movie"
	operationIDDiscoverTV    = "discover_tv"
	operationIDTrending      = "trending"

	// Search operations
	operationIDSearchMovie  = "search_movie"
	operationIDSearchTV     = "search_tv"
	operationIDSearchPerson = "search_person"
	operationIDSearchMulti  = "search_multi"
)

// Define API paths as constants
const (
	// Genre API paths
	apiPathGenreMovieList = "/genre/movie/list"
	apiPathGenreTVList    = "/genre/tv/list"

	// Movie API paths
	apiPathMoviePopular    = "/movie/popular"
	apiPathMovieNowPlaying = "/movie/now_playing"
	apiPathMovieTopRated   = "/movie/top_rated"
	apiPathMovieUpcoming   = "/movie/upcoming"

	// TV Series API paths
	apiPathTVPopular     = "/tv/popular"
	apiPathTVOnTheAir    = "/tv/on_the_air"
	apiPathTVTopRated    = "/tv/top_rated"
	apiPathTVAiringToday = "/tv/airing_today"

	// Discover API paths
	// Find movies/tv series using over 30 filters and sort options.
	apiPathDiscoverMovie = "/discover/movie"
	apiPathDiscoverTV    = "/discover/tv"

	// Search API paths
	apiPathSearchMovie  = "/search/movie"
	apiPathSearchTV     = "/search/tv"
	apiPathSearchPerson = "/search/person"
	apiPathSearchMulti  = "/search/multi"

	// Trending API path format (需要动态拼接时间窗口)
	apiPathTrendingFormat = "/trending/all/%s"
)

// Parameter structs for each operation

// BasicParams for operations with common parameters
type BasicParams struct {
	Language string `mapstructure:"language" validate:"omitempty"`
	Page     int    `mapstructure:"page" validate:"omitempty,min=1,max=500"`
	Region   string `mapstructure:"region" validate:"omitempty"`
}

// TVParams for TV operations with timezone support
type TVParams struct {
	Language string `mapstructure:"language" validate:"omitempty"`
	Page     int    `mapstructure:"page" validate:"omitempty,min=1,max=500"`
	Timezone string `mapstructure:"timezone" validate:"omitempty"`
}

// DiscoverMovieParams for movie discovery
type DiscoverMovieParams struct {
	Language           string `mapstructure:"language" validate:"omitempty"`
	Page               int    `mapstructure:"page" validate:"omitempty,min=1,max=500"`
	SortBy             string `mapstructure:"sort_by" validate:"omitempty"`
	WithGenres         string `mapstructure:"with_genres" validate:"omitempty"`
	PrimaryReleaseYear string `mapstructure:"primary_release_year" validate:"omitempty"`
	Region             string `mapstructure:"region" validate:"omitempty"`
}

// DiscoverTVParams for TV discovery
type DiscoverTVParams struct {
	Language         string `mapstructure:"language" validate:"omitempty"`
	Page             int    `mapstructure:"page" validate:"omitempty,min=1,max=500"`
	SortBy           string `mapstructure:"sort_by" validate:"omitempty"`
	WithGenres       string `mapstructure:"with_genres" validate:"omitempty"`
	FirstAirDateYear string `mapstructure:"first_air_date_year" validate:"omitempty"`
}

// TrendingParams for trending content
type TrendingParams struct {
	TimeWindow string `mapstructure:"time_window" validate:"required,oneof=day week"`
	Language   string `mapstructure:"language" validate:"omitempty"`
}

// SearchMovieParams for movie search
type SearchMovieParams struct {
	Query              string `mapstructure:"query" validate:"required"`
	Language           string `mapstructure:"language" validate:"omitempty"`
	Page               int    `mapstructure:"page" validate:"omitempty,min=1,max=500"`
	IncludeAdult       bool   `mapstructure:"include_adult" validate:"omitempty"`
	Region             string `mapstructure:"region" validate:"omitempty"`
	Year               int    `mapstructure:"year" validate:"omitempty"`
	PrimaryReleaseYear int    `mapstructure:"primary_release_year" validate:"omitempty"`
}

// SearchTVParams for TV search
type SearchTVParams struct {
	Query            string `mapstructure:"query" validate:"required"`
	Language         string `mapstructure:"language" validate:"omitempty"`
	Page             int    `mapstructure:"page" validate:"omitempty,min=1,max=500"`
	IncludeAdult     bool   `mapstructure:"include_adult" validate:"omitempty"`
	FirstAirDateYear int    `mapstructure:"first_air_date_year" validate:"omitempty"`
}

// SearchPersonParams for person search
type SearchPersonParams struct {
	Query        string `mapstructure:"query" validate:"required"`
	Language     string `mapstructure:"language" validate:"omitempty"`
	Page         int    `mapstructure:"page" validate:"omitempty,min=1,max=500"`
	IncludeAdult bool   `mapstructure:"include_adult" validate:"omitempty"`
}

// SearchMultiParams for multi search
type SearchMultiParams struct {
	Query        string `mapstructure:"query" validate:"required"`
	Language     string `mapstructure:"language" validate:"omitempty"`
	Page         int    `mapstructure:"page" validate:"omitempty,min=1,max=500"`
	IncludeAdult bool   `mapstructure:"include_adult" validate:"omitempty"`
}

// OperationHandler defines the signature of operation handler functions
type OperationHandler func(ctx context.Context, params interface{}) (map[string]interface{}, error)

// OperationDefinition defines operation structure
type OperationDefinition struct {
	Schema  interface{}      // Parameter struct pointer
	Handler OperationHandler // Operation handler function
}

// Operations maps operation ID to its definition
type Operations map[string]OperationDefinition

// RegisterOperation registers an operation
func (a *TmdbAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
	// Register the parameter schema with the base adapter for validation/decoding
	a.BaseAdapter.RegisterOperation(operationID, schema)

	// Store the definition in the TmdbAdapter's operations map
	if a.operations == nil {
		a.operations = make(Operations)
	}
	a.operations[operationID] = OperationDefinition{
		Schema:  schema,
		Handler: handler,
	}
}

// registerOperations populates the operations map in the TmdbAdapter
func (a *TmdbAdapter) registerOperations() {
	// Register genre operations
	a.RegisterOperation(operationIDGenreMovieList, &BasicParams{}, handleGenreMovieList)
	a.RegisterOperation(operationIDGenreTVList, &BasicParams{}, handleGenreTVList)

	// Register movie operations
	a.RegisterOperation(operationIDMoviePopular, &BasicParams{}, handleMoviePopular)
	a.RegisterOperation(operationIDMovieNowPlaying, &BasicParams{}, handleMovieNowPlaying)
	a.RegisterOperation(operationIDMovieTopRated, &BasicParams{}, handleMovieTopRated)
	a.RegisterOperation(operationIDMovieUpcoming, &BasicParams{}, handleMovieUpcoming)

	// Register TV operations
	a.RegisterOperation(operationIDTVPopular, &BasicParams{}, handleTVPopular)
	a.RegisterOperation(operationIDTVOnTheAir, &TVParams{}, handleTVOnTheAir)
	a.RegisterOperation(operationIDTVTopRated, &BasicParams{}, handleTVTopRated)
	a.RegisterOperation(operationIDTVAiringToday, &TVParams{}, handleTVAiringToday)

	// Register discover operations
	a.RegisterOperation(operationIDDiscoverMovie, &DiscoverMovieParams{}, handleDiscoverMovie)
	a.RegisterOperation(operationIDDiscoverTV, &DiscoverTVParams{}, handleDiscoverTV)

	// Register trending operation
	a.RegisterOperation(operationIDTrending, &TrendingParams{}, handleTrending)

	// Register search operations
	a.RegisterOperation(operationIDSearchMovie, &SearchMovieParams{}, handleSearchMovie)
	a.RegisterOperation(operationIDSearchTV, &SearchTVParams{}, handleSearchTV)
	a.RegisterOperation(operationIDSearchPerson, &SearchPersonParams{}, handleSearchPerson)
	a.RegisterOperation(operationIDSearchMulti, &SearchMultiParams{}, handleSearchMulti)
}

// Helper function to build query parameters
func buildQueryParams(params interface{}) map[string]string {
	queryParams := make(map[string]string)

	switch p := params.(type) {
	case *BasicParams:
		if p.Language != "" {
			queryParams["language"] = p.Language
		}
		if p.Page > 0 {
			queryParams["page"] = strconv.Itoa(p.Page)
		}
		if p.Region != "" {
			queryParams["region"] = p.Region
		}
	case *TVParams:
		if p.Language != "" {
			queryParams["language"] = p.Language
		}
		if p.Page > 0 {
			queryParams["page"] = strconv.Itoa(p.Page)
		}
		if p.Timezone != "" {
			queryParams["timezone"] = p.Timezone
		}
	}
	return queryParams
}

// Operation handlers

func handleGenreMovieList(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*BasicParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDGenreMovieList)
	}

	queryParams := buildQueryParams(params)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathGenreMovieList,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleGenreTVList(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*BasicParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDGenreTVList)
	}

	queryParams := buildQueryParams(params)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathGenreTVList,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleMoviePopular(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*BasicParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDMoviePopular)
	}

	queryParams := buildQueryParams(params)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathMoviePopular,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleMovieNowPlaying(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*BasicParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDMovieNowPlaying)
	}

	queryParams := buildQueryParams(params)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathMovieNowPlaying,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleMovieTopRated(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*BasicParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDMovieTopRated)
	}

	queryParams := buildQueryParams(params)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathMovieTopRated,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleMovieUpcoming(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*BasicParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDMovieUpcoming)
	}

	queryParams := buildQueryParams(params)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathMovieUpcoming,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleTVPopular(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*BasicParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDTVPopular)
	}

	queryParams := buildQueryParams(params)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathTVPopular,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleTVOnTheAir(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*TVParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDTVOnTheAir)
	}

	queryParams := buildQueryParams(params)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathTVOnTheAir,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleTVTopRated(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*BasicParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDTVTopRated)
	}

	queryParams := buildQueryParams(params)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathTVTopRated,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleTVAiringToday(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*TVParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDTVAiringToday)
	}

	queryParams := buildQueryParams(params)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathTVAiringToday,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleDiscoverMovie(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*DiscoverMovieParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDDiscoverMovie)
	}

	queryParams := make(map[string]string)
	if params.Language != "" {
		queryParams["language"] = params.Language
	}
	if params.Page > 0 {
		queryParams["page"] = strconv.Itoa(params.Page)
	}
	if params.SortBy != "" {
		queryParams["sort_by"] = params.SortBy
	}
	if params.WithGenres != "" {
		queryParams["with_genres"] = params.WithGenres
	}
	if params.PrimaryReleaseYear != "" {
		queryParams["primary_release_year"] = params.PrimaryReleaseYear
	}
	if params.Region != "" {
		queryParams["region"] = params.Region
	}

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathDiscoverMovie,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleDiscoverTV(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*DiscoverTVParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDDiscoverTV)
	}

	queryParams := make(map[string]string)
	if params.Language != "" {
		queryParams["language"] = params.Language
	}
	if params.Page > 0 {
		queryParams["page"] = strconv.Itoa(params.Page)
	}
	if params.SortBy != "" {
		queryParams["sort_by"] = params.SortBy
	}
	if params.WithGenres != "" {
		queryParams["with_genres"] = params.WithGenres
	}
	if params.FirstAirDateYear != "" {
		queryParams["first_air_date_year"] = params.FirstAirDateYear
	}

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathDiscoverTV,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleTrending(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*TrendingParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDTrending)
	}

	queryParams := make(map[string]string)
	if params.Language != "" {
		queryParams["language"] = params.Language
	}

	path := fmt.Sprintf(apiPathTrendingFormat, params.TimeWindow)

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         path,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleSearchMovie(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*SearchMovieParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDSearchMovie)
	}

	queryParams := map[string]string{
		"query": params.Query,
	}
	if params.Language != "" {
		queryParams["language"] = params.Language
	}
	if params.Page > 0 {
		queryParams["page"] = strconv.Itoa(params.Page)
	}
	if params.IncludeAdult {
		queryParams["include_adult"] = "true"
	}
	if params.Region != "" {
		queryParams["region"] = params.Region
	}
	if params.Year > 0 {
		queryParams["year"] = strconv.Itoa(params.Year)
	}
	if params.PrimaryReleaseYear > 0 {
		queryParams["primary_release_year"] = strconv.Itoa(params.PrimaryReleaseYear)
	}

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathSearchMovie,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleSearchTV(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*SearchTVParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDSearchTV)
	}

	queryParams := map[string]string{
		"query": params.Query,
	}
	if params.Language != "" {
		queryParams["language"] = params.Language
	}
	if params.Page > 0 {
		queryParams["page"] = strconv.Itoa(params.Page)
	}
	if params.IncludeAdult {
		queryParams["include_adult"] = "true"
	}
	if params.FirstAirDateYear > 0 {
		queryParams["first_air_date_year"] = strconv.Itoa(params.FirstAirDateYear)
	}

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathSearchTV,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleSearchPerson(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*SearchPersonParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDSearchPerson)
	}

	queryParams := map[string]string{
		"query": params.Query,
	}
	if params.Language != "" {
		queryParams["language"] = params.Language
	}
	if params.Page > 0 {
		queryParams["page"] = strconv.Itoa(params.Page)
	}
	if params.IncludeAdult {
		queryParams["include_adult"] = "true"
	}

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathSearchPerson,
		"query_params": queryParams,
	}

	return restParams, nil
}

func handleSearchMulti(ctx context.Context, processedParams interface{}) (map[string]interface{}, error) {
	params, ok := processedParams.(*SearchMultiParams)
	if !ok {
		return nil, fmt.Errorf("internal error: unexpected parameter type for %s", operationIDSearchMulti)
	}

	queryParams := map[string]string{
		"query": params.Query,
	}
	if params.Language != "" {
		queryParams["language"] = params.Language
	}
	if params.Page > 0 {
		// int to string
		queryParams["page"] = strconv.Itoa(params.Page)
	}
	if params.IncludeAdult {
		queryParams["include_adult"] = "true"
	}

	restParams := map[string]interface{}{
		"method":       http.MethodGet,
		"path":         apiPathSearchMulti,
		"query_params": queryParams,
	}

	return restParams, nil
}
