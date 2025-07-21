package amap

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// Amap Basic API endpoint constants
const (
	// Geocoding/Reverse geocoding
	endpointGeocoding = "/v3/geocode/geo"   // Geocoding
	endpointRegeo     = "/v3/geocode/regeo" // Reverse geocoding

	// Route planning
	endpointDirectionWalking   = "/v3/direction/walking"            // Walking route planning
	endpointDirectionDriving   = "/v3/direction/driving"            // Driving route planning
	endpointDirectionTransit   = "/v3/direction/transit/integrated" // Transit route planning
	endpointDirectionBicycling = "/v4/direction/bicycling"          // Bicycling route planning
	endpointDistance           = "/v3/distance"                     // Distance measurement

	// Administrative district query
	endpointDistrict = "/v3/config/district" // Administrative district query

	// IP location
	endpointIP = "/v3/ip" // IP location

	// Coordinate conversion
	endpointConvert = "/v3/assistant/coordinate/convert" // Coordinate conversion

	// Amap Advanced API endpoint constants
	// Weather query
	endpointWeatherInfo = "/v3/weather/weatherInfo" // Weather query
)

// Basic API operation ID constants
const (
	// Geocoding
	operationIDGeocoding = "geocoding" // Geocoding
	operationIDRegeo     = "regeo"     // Reverse geocoding

	// Route planning
	operationIDDirectionWalking   = "direction_walking"            // Walking route planning
	operationIDDirectionDriving   = "direction_driving"            // Driving route planning
	operationIDDirectionTransit   = "direction_transit_integrated" // Transit route planning
	operationIDDirectionBicycling = "direction_bicycling"          // Bicycling route planning
	operationIDDistance           = "distance"                     // Distance measurement

	// Administrative district
	operationIDDistrict = "district" // Administrative district query

	// IP location
	operationIDIP = "ip_location" // IP location

	// Coordinate conversion
	operationIDConvert = "coordinate_convert" // Coordinate conversion

	// Advanced API operation ID constants
	// Weather query
	operationIDWeatherInfo = "weather_info" // Weather information
)

// Geocoding parameters structure
type GeocodingParams struct {
	Address string `mapstructure:"address" validate:"required"`
	City    string `mapstructure:"city" validate:"omitempty"`
}

// Reverse geocoding parameters structure
type RegeoParams struct {
	Location   string `mapstructure:"location" validate:"required"`
	Radius     *int   `mapstructure:"radius" validate:"omitempty,min=1,max=3000"`
	Extensions string `mapstructure:"extensions" validate:"omitempty,oneof=base all"`
	Poitype    string `mapstructure:"poitype" validate:"omitempty"`
	Roadlevel  *int   `mapstructure:"roadlevel" validate:"omitempty,oneof=0 1"`
}

// Driving route planning parameters structure
type DirectionDrivingParams struct {
	Origin      string `mapstructure:"origin" validate:"required"`
	Destination string `mapstructure:"destination" validate:"required"`
	Strategy    *int   `mapstructure:"strategy" validate:"omitempty,min=0,max=10"`
	Waypoints   string `mapstructure:"waypoints" validate:"omitempty"`
	Extensions  string `mapstructure:"extensions" validate:"omitempty,oneof=base all"`
}

// Walking route planning parameters structure
type DirectionWalkingParams struct {
	Origin      string `mapstructure:"origin" validate:"required"`
	Destination string `mapstructure:"destination" validate:"required"`
}

// Transit route planning parameters structure
type DirectionTransitParams struct {
	Origin      string `mapstructure:"origin" validate:"required"`
	Destination string `mapstructure:"destination" validate:"required"`
	City        string `mapstructure:"city" validate:"required"`
	Cityd       string `mapstructure:"cityd" validate:"omitempty"`
	Extensions  string `mapstructure:"extensions" validate:"omitempty,oneof=base all"`
	Strategy    *int   `mapstructure:"strategy" validate:"omitempty,min=0,max=5"`
}

// Bicycling route planning parameters structure
type DirectionBicyclingParams struct {
	Origin      string `mapstructure:"origin" validate:"required"`
	Destination string `mapstructure:"destination" validate:"required"`
}

// Distance measurement parameters structure
type DistanceParams struct {
	Origins     string `mapstructure:"origins" validate:"required"`
	Destination string `mapstructure:"destination" validate:"required"`
	Type        *int   `mapstructure:"type" validate:"omitempty,min=0,max=3"`
}

// Administrative district query parameters structure
type DistrictParams struct {
	Keywords    string `mapstructure:"keywords" validate:"omitempty"`
	Subdistrict *int   `mapstructure:"subdistrict" validate:"omitempty,min=0,max=3"`
	Extensions  string `mapstructure:"extensions" validate:"omitempty,oneof=base all"`
}

// IP location parameters structure
type IPLocationParams struct {
	IP string `mapstructure:"ip" validate:"omitempty"`
}

// Coordinate conversion parameters structure
type CoordinateConvertParams struct {
	Locations string `mapstructure:"locations" validate:"required"`
	Coordsys  string `mapstructure:"coordsys" validate:"omitempty,oneof=gps mapbar baidu autonavi"`
}

// Weather information parameters structure
type WeatherParams struct {
	City       string `mapstructure:"city" validate:"required"`
	Extensions string `mapstructure:"extensions" validate:"omitempty,oneof=base all"`
}

// OperationHandler defines the operation handler function signature
type OperationHandler func(ctx context.Context, params interface{}) (map[string]interface{}, error)

// OperationDefinition defines the operation definition structure (API Key version, no permissions required)
type OperationDefinition struct {
	Schema  interface{}      // Parameter schema (struct pointer)
	Handler OperationHandler // Operation handler function
}

// Operations maps operation IDs to their definitions
type Operations map[string]OperationDefinition

// RegisterOperation registers an operation (API Key version)
func (a *AmapAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
	a.BaseAdapter.RegisterOperation(operationID, schema)

	if a.operations == nil {
		a.operations = make(Operations)
	}

	a.operations[operationID] = OperationDefinition{
		Schema:  schema,
		Handler: handler,
	}
}

// registerOperations registers all Amap API operations
func (a *AmapAdapter) registerOperations() {
	// Basic API operations
	a.RegisterOperation(operationIDGeocoding, &GeocodingParams{}, handleGeocoding)
	a.RegisterOperation(operationIDRegeo, &RegeoParams{}, handleRegeo)
	a.RegisterOperation(operationIDDirectionWalking, &DirectionWalkingParams{}, handleDirectionWalking)
	a.RegisterOperation(operationIDDirectionDriving, &DirectionDrivingParams{}, handleDirectionDriving)
	a.RegisterOperation(operationIDDirectionTransit, &DirectionTransitParams{}, handleDirectionTransit)
	a.RegisterOperation(operationIDDirectionBicycling, &DirectionBicyclingParams{}, handleDirectionBicycling)
	a.RegisterOperation(operationIDDistance, &DistanceParams{}, handleDistance)
	a.RegisterOperation(operationIDDistrict, &DistrictParams{}, handleDistrict)
	a.RegisterOperation(operationIDIP, &IPLocationParams{}, handleIPLocation)
	a.RegisterOperation(operationIDConvert, &CoordinateConvertParams{}, handleCoordinateConvert)

	// Advanced API operations
	a.RegisterOperation(operationIDWeatherInfo, &WeatherParams{}, handleWeatherInfo)
}

// handleGeocoding handles geocoding operation
func handleGeocoding(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	geocodingParams, ok := params.(*GeocodingParams)
	if !ok {
		return nil, fmt.Errorf("invalid geocoding parameters")
	}

	queryParams := map[string]string{
		"address": geocodingParams.Address,
	}

	if geocodingParams.City != "" {
		queryParams["city"] = geocodingParams.City
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointGeocoding,
		"query_params": queryParams,
	}, nil
}

// handleRegeo handles reverse geocoding operation
func handleRegeo(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	regeoParams, ok := params.(*RegeoParams)
	if !ok {
		return nil, fmt.Errorf("invalid reverse geocoding parameters")
	}

	queryParams := map[string]string{
		"location": regeoParams.Location,
	}

	if regeoParams.Radius != nil {
		queryParams["radius"] = strconv.Itoa(*regeoParams.Radius)
	}

	if regeoParams.Extensions != "" {
		queryParams["extensions"] = regeoParams.Extensions
	}

	if regeoParams.Poitype != "" {
		queryParams["poitype"] = regeoParams.Poitype
	}

	if regeoParams.Roadlevel != nil {
		queryParams["roadlevel"] = strconv.Itoa(*regeoParams.Roadlevel)
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointRegeo,
		"query_params": queryParams,
	}, nil
}

// handleDirectionDriving handles driving route planning operation
func handleDirectionDriving(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	drivingParams, ok := params.(*DirectionDrivingParams)
	if !ok {
		return nil, fmt.Errorf("invalid driving route planning parameters")
	}

	queryParams := map[string]string{
		"origin":      drivingParams.Origin,
		"destination": drivingParams.Destination,
	}

	if drivingParams.Strategy != nil {
		queryParams["strategy"] = strconv.Itoa(*drivingParams.Strategy)
	}

	if drivingParams.Waypoints != "" {
		queryParams["waypoints"] = drivingParams.Waypoints
	}

	if drivingParams.Extensions != "" {
		queryParams["extensions"] = drivingParams.Extensions
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointDirectionDriving,
		"query_params": queryParams,
	}, nil
}

// handleDirectionWalking handles walking route planning operation
func handleDirectionWalking(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	walkingParams, ok := params.(*DirectionWalkingParams)
	if !ok {
		return nil, fmt.Errorf("invalid walking route planning parameters")
	}

	queryParams := map[string]string{
		"origin":      walkingParams.Origin,
		"destination": walkingParams.Destination,
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointDirectionWalking,
		"query_params": queryParams,
	}, nil
}

// handleDirectionTransit handles transit route planning operation
func handleDirectionTransit(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	transitParams, ok := params.(*DirectionTransitParams)
	if !ok {
		return nil, fmt.Errorf("invalid transit route planning parameters")
	}

	queryParams := map[string]string{
		"origin":      transitParams.Origin,
		"destination": transitParams.Destination,
		"city":        transitParams.City,
	}

	if transitParams.Cityd != "" {
		queryParams["cityd"] = transitParams.Cityd
	}

	if transitParams.Extensions != "" {
		queryParams["extensions"] = transitParams.Extensions
	}

	if transitParams.Strategy != nil {
		queryParams["strategy"] = strconv.Itoa(*transitParams.Strategy)
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointDirectionTransit,
		"query_params": queryParams,
	}, nil
}

// handleDirectionBicycling handles bicycling route planning operation
func handleDirectionBicycling(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	bicyclingParams, ok := params.(*DirectionBicyclingParams)
	if !ok {
		return nil, fmt.Errorf("invalid bicycling route planning parameters")
	}

	queryParams := map[string]string{
		"origin":      bicyclingParams.Origin,
		"destination": bicyclingParams.Destination,
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointDirectionBicycling,
		"query_params": queryParams,
	}, nil
}

// handleDistance handles distance measurement operation
func handleDistance(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	distanceParams, ok := params.(*DistanceParams)
	if !ok {
		return nil, fmt.Errorf("invalid distance measurement parameters")
	}

	queryParams := map[string]string{
		"origins":     distanceParams.Origins,
		"destination": distanceParams.Destination,
	}

	if distanceParams.Type != nil {
		queryParams["type"] = strconv.Itoa(*distanceParams.Type)
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointDistance,
		"query_params": queryParams,
	}, nil
}

// handleDistrict handles administrative district query operation
func handleDistrict(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	districtParams, ok := params.(*DistrictParams)
	if !ok {
		return nil, fmt.Errorf("invalid administrative district query parameters")
	}

	queryParams := make(map[string]string)

	if districtParams.Keywords != "" {
		queryParams["keywords"] = districtParams.Keywords
	}

	if districtParams.Subdistrict != nil {
		queryParams["subdistrict"] = strconv.Itoa(*districtParams.Subdistrict)
	}

	if districtParams.Extensions != "" {
		queryParams["extensions"] = districtParams.Extensions
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointDistrict,
		"query_params": queryParams,
	}, nil
}

// handleIPLocation handles IP location operation
func handleIPLocation(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	ipParams, ok := params.(*IPLocationParams)
	if !ok {
		return nil, fmt.Errorf("invalid IP location parameters")
	}

	queryParams := make(map[string]string)

	if ipParams.IP != "" {
		queryParams["ip"] = ipParams.IP
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointIP,
		"query_params": queryParams,
	}, nil
}

// handleCoordinateConvert handles coordinate conversion operation
func handleCoordinateConvert(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	convertParams, ok := params.(*CoordinateConvertParams)
	if !ok {
		return nil, fmt.Errorf("invalid coordinate conversion parameters")
	}

	queryParams := map[string]string{
		"locations": convertParams.Locations,
	}

	if convertParams.Coordsys != "" {
		queryParams["coordsys"] = convertParams.Coordsys
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointConvert,
		"query_params": queryParams,
	}, nil
}

// handleWeatherInfo handles weather information operation
func handleWeatherInfo(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	weatherParams, ok := params.(*WeatherParams)
	if !ok {
		return nil, fmt.Errorf("invalid weather parameters")
	}

	queryParams := map[string]string{
		"city": weatherParams.City,
	}

	if weatherParams.Extensions != "" {
		queryParams["extensions"] = weatherParams.Extensions
	}

	return map[string]interface{}{
		"method":       http.MethodGet,
		"path":         endpointWeatherInfo,
		"query_params": queryParams,
	}, nil
}
