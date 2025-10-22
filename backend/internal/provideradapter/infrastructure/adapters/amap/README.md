# Amap Adapter

This module provides a unified adapter for the Amap (Gaode Map) API in the Context-Space platform, supporting geocoding, route planning, administrative district queries, weather, and other LBS features.

## Features
- Address and coordinate conversion (geocoding/reverse geocoding)
- Route planning (walking, driving, transit, bicycling)
- Distance measurement
- Administrative district queries
- IP location
- Static map generation
- Coordinate system conversion
- Weather information query

## Supported Operations
| Operation ID                 | Description           | Main Parameters                |
| ---------------------------- | --------------------- | ------------------------------ |
| geocoding                    | Geocoding             | address, city                  |
| regeo                        | Reverse Geocoding     | location, radius, extensions   |
| direction_walking            | Walking Route         | origin, destination            |
| direction_driving            | Driving Route         | origin, destination            |
| direction_transit_integrated | Transit Route         | origin, destination, city      |
| direction_bicycling          | Bicycling Route       | origin, destination            |
| distance                     | Distance Measurement  | origins, destination           |
| district                     | District Query        | keywords                       |
| ip_location                  | IP Location           | ip                             |

| coordinate_convert           | Coordinate Conversion | locations                      |
| weather_info                 | Weather Info          | city, extensions               |

## Quick Usage Example

```go
providerInfo := &domain.ProviderAdapterInfo{
    Identifier: "amap",
    Name:       "Amap API",
    AuthType:   types.AuthTypeApiKey,
    Status:     providercore.ProviderStatusActive,
}
adapterConfig := &domain.AdapterConfig{Timeout: 10 * time.Second}
restConfig := &rest.RESTConfig{BaseURL: "https://restapi.amap.com"}
restAdapter := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)
adapter := NewAmapAdapter(providerInfo, adapterConfig, restAdapter, "<YOUR_API_KEY>")

ctx := context.Background()
params := map[string]interface{}{"address": "Beijing Chaoyang District"}
result, err := adapter.Execute(ctx, "geocoding", params, nil)
if err != nil {
    // handle error
}
fmt.Println(result)
```

## API Key Configuration
- You must provide a valid Amap API Key when creating the adapter.
- It is recommended to store the API Key in a secure config file or environment variable.

## Testing

This module includes comprehensive unit tests covering all operations:

```bash
go test ./internal/provideradapter/infrastructure/adapters/amap -v
```

You can also run only the full integration test:

```bash
go test ./internal/provideradapter/infrastructure/adapters/amap -v -run TestAmapAdapter_AllOperations
```

## Notes
- All operations use GET requests. Parameters must conform to the official Amap API documentation.
- Some endpoints (such as static map, weather, etc.) have special parameter requirements. See manifest.json for details.
- The returned result is the original JSON structure from Amap. You may need to further parse it as needed.

## References
- [Amap Open Platform Official Documentation](https://lbs.amap.com/api/webservice/summary)
- See manifest.json and i18n/*.json for operation and parameter details
 