package openweathermap

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// 定义 API 端点常量
const (
	endpointCurrentWeather  = "/data/2.5/weather"
	endpointWeatherForecast = "/data/2.5/forecast"
	endpointOneCallWeather  = "/data/2.5/onecall"
	endpointAirPollution    = "/data/2.5/air_pollution"
	endpointGeocoding       = "/geo/1.0/direct"
)

// 定义操作 ID 常量
const (
	operationIDGetCurrentWeather  = "get_current_weather"
	operationIDGetWeatherForecast = "get_weather_forecast"
	operationIDGetOneCallWeather  = "get_one_call_weather"
	operationIDGetAirPollution    = "get_air_pollution"
	operationIDGetGeocoding       = "get_geocoding"
)

// 参数结构体定义
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

// OperationHandler 定义操作处理函数的签名
type OperationHandler func(ctx context.Context, params interface{}) (map[string]interface{}, error)

// OperationDefinition 定义操作结构（API Key 认证版本）
type OperationDefinition struct {
	Schema  interface{}      // 参数结构体指针
	Handler OperationHandler // 操作处理函数
}

// Operations 映射操作 ID 到其定义
type Operations map[string]OperationDefinition

// RegisterOperation 注册操作（API Key 版本）
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

// registerOperations 注册所有操作
func (a *OpenWeatherMapAdapter) registerOperations() {
	// 注册当前天气操作
	a.RegisterOperation(operationIDGetCurrentWeather, &GetCurrentWeatherParams{}, a.handleGetCurrentWeather)

	// 注册天气预报操作
	a.RegisterOperation(operationIDGetWeatherForecast, &GetWeatherForecastParams{}, a.handleGetWeatherForecast)

	// 注册综合天气数据操作
	a.RegisterOperation(operationIDGetOneCallWeather, &GetOneCallWeatherParams{}, a.handleGetOneCallWeather)

	// 注册空气污染数据操作
	a.RegisterOperation(operationIDGetAirPollution, &GetAirPollutionParams{}, a.handleGetAirPollution)

	// 注册地理编码操作
	a.RegisterOperation(operationIDGetGeocoding, &GetGeocodingParams{}, a.handleGetGeocoding)
}

// handleGetCurrentWeather 处理当前天气查询
func (a *OpenWeatherMapAdapter) handleGetCurrentWeather(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var weatherParams GetCurrentWeatherParams
	if err := mapstructure.Decode(params, &weatherParams); err != nil {
		return nil, fmt.Errorf("参数解析失败: %v", err)
	}

	// 构建查询参数
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
		"url":          endpointCurrentWeather,
		"method":       "GET",
		"query_params": queryParams,
	}, nil
}

// handleGetWeatherForecast 处理天气预报查询
func (a *OpenWeatherMapAdapter) handleGetWeatherForecast(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var forecastParams GetWeatherForecastParams
	if err := mapstructure.Decode(params, &forecastParams); err != nil {
		return nil, fmt.Errorf("参数解析失败: %v", err)
	}

	// 构建查询参数
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
		"url":          endpointWeatherForecast,
		"method":       "GET",
		"query_params": queryParams,
	}, nil
}

// handleGetOneCallWeather 处理一次性天气数据查询
func (a *OpenWeatherMapAdapter) handleGetOneCallWeather(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var oneCallParams GetOneCallWeatherParams
	if err := mapstructure.Decode(params, &oneCallParams); err != nil {
		return nil, fmt.Errorf("参数解析失败: %v", err)
	}

	// 构建查询参数
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
		"url":          endpointOneCallWeather,
		"method":       "GET",
		"query_params": queryParams,
	}, nil
}

// handleGetAirPollution 处理空气污染数据查询
func (a *OpenWeatherMapAdapter) handleGetAirPollution(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var pollutionParams GetAirPollutionParams
	if err := mapstructure.Decode(params, &pollutionParams); err != nil {
		return nil, fmt.Errorf("参数解析失败: %v", err)
	}

	// 构建查询参数
	queryParams := make(map[string]string)
	queryParams["lat"] = fmt.Sprintf("%.6f", pollutionParams.Lat)
	queryParams["lon"] = fmt.Sprintf("%.6f", pollutionParams.Lon)

	return map[string]interface{}{
		"url":          endpointAirPollution,
		"method":       "GET",
		"query_params": queryParams,
	}, nil
}

// handleGetGeocoding 处理地理编码查询
func (a *OpenWeatherMapAdapter) handleGetGeocoding(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	var geocodingParams GetGeocodingParams
	if err := mapstructure.Decode(params, &geocodingParams); err != nil {
		return nil, fmt.Errorf("参数解析失败: %v", err)
	}

	// 构建查询参数
	queryParams := make(map[string]string)
	queryParams["q"] = geocodingParams.Q
	if geocodingParams.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(geocodingParams.Limit)
	}

	return map[string]interface{}{
		"url":          endpointGeocoding,
		"method":       "GET",
		"query_params": queryParams,
	}, nil
}

// buildQueryString 构建查询字符串
func buildQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	var parts []string
	for key, value := range params {
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(parts, "&")
}
