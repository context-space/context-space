package openweathermap

//  Operations Implementation
import (
	"context"
	"fmt"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

// Define API endpoint constants
const (
	endpointCurrentWeather  = "/data/2.5/weather"
	endpointWeatherForecast = "/data/2.5/forecast"
	endpointOneCallWeather  = "/data/2.5/onecall"
	endpointAirPollution    = "/data/2.5/air_pollution"
	endpointGeocoding       = "/geo/1.0/direct"
)

// Define operation ID constants
const (
	operationIDGetCurrentWeather  = "get_current_weather"
	operationIDGetWeatherForecast = "get_weather_forecast"
	operationIDGetOneCallWeather  = "get_one_call_weather"
	operationIDGetAirPollution    = "get_air_pollution"
	operationIDGetGeocoding       = "get_geocoding"
)

// Parameter struct definitions
type GetCurrentWeatherParams struct {
	Lat   float64 `mapstructure:"lat" validate:"required,min=-90,max=90"`
	Lon   float64 `mapstructure:"lon" validate:"required,min=-180,max=180"`
	Units string  `mapstructure:"units" validate:"omitempty,oneof=metric imperial standard"`
	Lang  string  `mapstructure:"lang" validate:"omitempty"`
}

type GetWeatherForecastParams struct {
	Lat   float64 `mapstructure:"lat" validate:"required,min=-90,max=90"`
	Lon   float64 `mapstructure:"lon" validate:"required,min=-180,max=180"`
	Units string  `mapstructure:"units" validate:"omitempty,oneof=metric imperial standard"`
	Lang  string  `mapstructure:"lang" validate:"omitempty"`
}

type GetOneCallWeatherParams struct {
	Lat     float64 `mapstructure:"lat" validate:"required,min=-90,max=90"`
	Lon     float64 `mapstructure:"lon" validate:"required,min=-180,max=180"`
	Exclude string  `mapstructure:"exclude" validate:"omitempty"`
	Units   string  `mapstructure:"units" validate:"omitempty,oneof=metric imperial standard"`
	Lang    string  `mapstructure:"lang" validate:"omitempty"`
}

type GetAirPollutionParams struct {
	Lat float64 `mapstructure:"lat" validate:"required,min=-90,max=90"`
	Lon float64 `mapstructure:"lon" validate:"required,min=-180,max=180"`
}

type GetGeocodingParams struct {
	Q     string `mapstructure:"q" validate:"required"`
	Limit int    `mapstructure:"limit" validate:"omitempty,min=1,max=5"`
}

// OperationHandler defines the signature of operation handler functions
type OperationHandler func(ctx context.Context, params interface{}) (map[string]interface{}, error)

// OperationDefinition defines operation structure (API Key authentication version)
type OperationDefinition struct {
	Schema  interface{}      // Parameter struct pointer
	Handler OperationHandler // Operation handler function
}

// Operations maps operation ID to its definition
type Operations map[string]OperationDefinition

// RegisterOperation registers operation (API Key version)
func (a *OpenWeatherMapAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
	a.BaseAdapter.RegisterOperation(operationID, schema)
	if a.operations == nil {
		a.operations = make(Operations)
	}
	a.operations[operationID] = OperationDefinition{
		Schema:  schema,
		Handler: handler,
	}
}

// registerOperations registers all operations
func (a *OpenWeatherMapAdapter) registerOperations() {
	// Register current weather operation
	a.RegisterOperation(operationIDGetCurrentWeather, &GetCurrentWeatherParams{}, a.handleGetCurrentWeather)

	// Register weather forecast operation
	a.RegisterOperation(operationIDGetWeatherForecast, &GetWeatherForecastParams{}, a.handleGetWeatherForecast)

	// Register one call weather data operation
	a.RegisterOperation(operationIDGetOneCallWeather, &GetOneCallWeatherParams{}, a.handleGetOneCallWeather)

	// Register air pollution data operation
	a.RegisterOperation(operationIDGetAirPollution, &GetAirPollutionParams{}, a.handleGetAirPollution)

	// Register geocoding operation
	a.RegisterOperation(operationIDGetGeocoding, &GetGeocodingParams{}, a.handleGetGeocoding)
}

// handleGetCurrentWeather handles current weather query
func (a *OpenWeatherMapAdapter) handleGetCurrentWeather(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var weatherParams GetCurrentWeatherParams
	if err := mapstructure.Decode(params, &weatherParams); err != nil {
		return nil, fmt.Errorf("parameter parsing failed: %v", err)
	}

	// Build query parameters
	queryParams := make(map[string]string)
	queryParams["lat"] = fmt.Sprintf("%.6f", weatherParams.Lat)
	queryParams["lon"] = fmt.Sprintf("%.6f", weatherParams.Lon)
	if weatherParams.Units != "" {
		queryParams["units"] = weatherParams.Units
	}
	if weatherParams.Lang != "" {
		queryParams["lang"] = weatherParams.Lang
	}

	return map[string]interface{}{
		"path":         endpointCurrentWeather,
		"method":       "GET",
		"query_params": queryParams,
	}, nil
}

// handleGetWeatherForecast handles weather forecast query
func (a *OpenWeatherMapAdapter) handleGetWeatherForecast(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var forecastParams GetWeatherForecastParams
	if err := mapstructure.Decode(params, &forecastParams); err != nil {
		return nil, fmt.Errorf("parameter parsing failed: %v", err)
	}

	// Build query parameters
	queryParams := make(map[string]string)
	queryParams["lat"] = fmt.Sprintf("%.6f", forecastParams.Lat)
	queryParams["lon"] = fmt.Sprintf("%.6f", forecastParams.Lon)
	if forecastParams.Units != "" {
		queryParams["units"] = forecastParams.Units
	}
	if forecastParams.Lang != "" {
		queryParams["lang"] = forecastParams.Lang
	}

	return map[string]interface{}{
		"path":         endpointWeatherForecast,
		"method":       "GET",
		"query_params": queryParams,
	}, nil
}

// handleGetOneCallWeather handles one call weather data query
func (a *OpenWeatherMapAdapter) handleGetOneCallWeather(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var oneCallParams GetOneCallWeatherParams
	if err := mapstructure.Decode(params, &oneCallParams); err != nil {
		return nil, fmt.Errorf("parameter parsing failed: %v", err)
	}

	// Build query parameters
	queryParams := make(map[string]string)
	queryParams["lat"] = fmt.Sprintf("%.6f", oneCallParams.Lat)
	queryParams["lon"] = fmt.Sprintf("%.6f", oneCallParams.Lon)
	if oneCallParams.Exclude != "" {
		queryParams["exclude"] = oneCallParams.Exclude
	}
	if oneCallParams.Units != "" {
		queryParams["units"] = oneCallParams.Units
	}
	if oneCallParams.Lang != "" {
		queryParams["lang"] = oneCallParams.Lang
	}

	return map[string]interface{}{
		"path":         endpointOneCallWeather,
		"method":       "GET",
		"query_params": queryParams,
	}, nil
}

// handleGetAirPollution handles air pollution data query
func (a *OpenWeatherMapAdapter) handleGetAirPollution(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var pollutionParams GetAirPollutionParams
	if err := mapstructure.Decode(params, &pollutionParams); err != nil {
		return nil, fmt.Errorf("parameter parsing failed: %v", err)
	}

	// Build query parameters
	queryParams := make(map[string]string)
	queryParams["lat"] = fmt.Sprintf("%.6f", pollutionParams.Lat)
	queryParams["lon"] = fmt.Sprintf("%.6f", pollutionParams.Lon)

	return map[string]interface{}{
		"path":         endpointAirPollution,
		"method":       "GET",
		"query_params": queryParams,
	}, nil
}

// handleGetGeocoding handles geocoding query
func (a *OpenWeatherMapAdapter) handleGetGeocoding(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var geocodingParams GetGeocodingParams
	if err := mapstructure.Decode(params, &geocodingParams); err != nil {
		return nil, fmt.Errorf("parameter parsing failed: %v", err)
	}

	// Build query parameters
	queryParams := make(map[string]string)
	queryParams["q"] = geocodingParams.Q
	if geocodingParams.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(geocodingParams.Limit)
	}

	return map[string]interface{}{
		"path":         endpointGeocoding,
		"method":       "GET",
		"query_params": queryParams,
	}, nil
}
