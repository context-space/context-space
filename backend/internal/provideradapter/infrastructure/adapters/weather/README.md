# OpenWeatherMap Weather Adapter

This adapter provides integration with the OpenWeatherMap API for weather data retrieval.

## Features

- Current weather data
- Weather forecasts
- Air pollution data
- Geocoding services
- One Call weather API

## Operations

### 1. get_current_weather

Get current weather conditions for a specific location.

**Parameters:**

- `lat` (required): Latitude (-90 to 90)
- `lon` (required): Longitude (-180 to 180)
- `units` (optional): Units of measurement (metric, imperial, standard)
- `lang` (optional): Language code for weather descriptions

### 2. get_weather_forecast

Get weather forecast for a specific location.

**Parameters:**

- `lat` (required): Latitude (-90 to 90)
- `lon` (required): Longitude (-180 to 180)
- `units` (optional): Units of measurement (metric, imperial, standard)
- `lang` (optional): Language code for weather descriptions

### 3. get_geocoding

Convert location name to coordinates.

**Parameters:**

- `q` (required): City name, state code, and country code
- `limit` (optional): Number of results (1-5)

### 4. get_air_pollution

Get air pollution data for a specific location.

**Parameters:**

- `lat` (required): Latitude (-90 to 90)
- `lon` (required): Longitude (-180 to 180)

### 5. get_one_call_weather 

Get comprehensive weather data using One Call API （partially free）.

**Parameters:**

- `lat` (required): Latitude (-90 to 90)
- `lon` (required): Longitude (-180 to 180)
- `exclude` (optional): Exclude data blocks (minutely, hourly, daily, alerts)
- `units` (optional): Units of measurement (metric, imperial, standard)
- `lang` (optional): Language code for weather descriptions

### Configuration

The adapter is automatically registered during package initialization. It requires:

- **Authentication**: API Key authentication
- **Base URL**: https://api.openweathermap.org
- **Timeout**: 10 seconds
- **Max Retries**: 3
- **Retry Backoff**: 1 second

## Example Usage

### Current Weather Example

```go
params := map[string]interface{}{
    "lat":   39.9042,
    "lon":   116.4074,
    "units": "metric",
}
// Beijing: Temperature in Celsius
```

### Weather Forecast Example

```go
params := map[string]interface{}{
    "lat":   31.2304,
    "lon":   121.4737,
    "units": "metric",
}
// Shanghai: 5-day forecast
```

### Geocoding Example

params := map[string]interface{}{
    "q":     "Beijing",
    "limit": 1,
}
// Get coordinates for Beijing

### Air Pollution Example

```go
params := map[string]interface{}{
    "lat": 22.3193,
    "lon": 114.1694,
}
// Shenzhen: Air quality index and pollutant concentrations
```

## Error Handling

The adapter handles various error conditions:

- Invalid coordinates (outside valid ranges)
- Missing required parameters
- API rate limits
- Network timeouts
- Invalid API keys

## Development

### File Structure

- `template.go`: Adapter template and registration
- `weather_adapter.go`: Main adapter implementation
- `weather_operations.go`: Operation definitions and handlers
- `weather_adapter_test.go`: Unit tests
- `weather_integration_test.go`: Integration tests
- `README.md`: This documentation

### Adding New Operations

1. Define parameter struct in `weather_operations.go`
2. Add operation handler function
3. Register operation in `registerOperations()` method
4. Add to supported operations list in template
5. Add tests for new operation

## API Reference

For complete API documentation, visit: https://openweathermap.org/api

## Support

- OpenWeatherMap API documentation: https://openweathermap.org/api
- Get API key: https://openweathermap.org/api/statistics-api
