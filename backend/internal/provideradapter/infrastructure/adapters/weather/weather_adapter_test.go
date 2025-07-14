package openweathermap

import (
	"testing"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	providercore "github.com/context-space/context-space/backend/internal/providercore/domain"
	"github.com/stretchr/testify/assert"
)

// TestAdapterRegistration 测试适配器是否正确注册
func TestAdapterRegistration(t *testing.T) {
	// 验证适配器模板是否已注册
	template, exists := registry.GetAdapterTemplate("openweathermap")
	assert.True(t, exists, "OpenWeatherMap adapter template should be registered")
	assert.NotNil(t, template, "OpenWeatherMap adapter template should not be nil")

	// 验证模板是否是正确的类型
	weatherTemplate, ok := template.(*OpenWeatherMapAdapterTemplate)
	assert.True(t, ok, "Template should be of type OpenWeatherMapAdapterTemplate")
	assert.NotNil(t, weatherTemplate, "Weather template should not be nil")

	// 验证适配器的基本信息（通过具体实现类型来测试）
	assert.Equal(t, "openweathermap", weatherTemplate.GetIdentifier(), "Identifier should match")
	assert.Equal(t, "OpenWeatherMap", weatherTemplate.GetName(), "Name should match")
	assert.Equal(t, "apikey", weatherTemplate.GetAuthType(), "Auth type should be apikey")
	assert.True(t, weatherTemplate.IsProduction(), "Adapter should be production ready")

	// 验证支持的操作
	operations := weatherTemplate.GetSupportedOperations()
	assert.NotEmpty(t, operations, "Should have supported operations")

	expectedOperations := []string{
		"get_current_weather",
		"get_weather_forecast",
		"get_one_call_weather",
		"get_air_pollution",
		"get_geocoding",
	}

	for _, expectedOp := range expectedOperations {
		assert.Contains(t, operations, expectedOp, "Should support operation: %s", expectedOp)
	}

	// 验证所需的凭证
	credentials := weatherTemplate.GetRequiredCredentials()
	assert.Contains(t, credentials, "api_key", "Should require api_key credential")

	// 验证文档和主页URL
	assert.Equal(t, "https://openweathermap.org/api", weatherTemplate.GetDocumentationURL(), "Documentation URL should match")
	assert.Equal(t, "https://openweathermap.org", weatherTemplate.GetHomepageURL(), "Homepage URL should match")
}

// TestAdapterTemplateValidation 测试适配器模板的配置验证
func TestAdapterTemplateValidation(t *testing.T) {
	template := &OpenWeatherMapAdapterTemplate{}

	// 测试空配置
	t.Run("Nil Provider Config", func(t *testing.T) {
		err := template.ValidateConfig(nil)
		assert.Error(t, err, "Should return error for nil provider config")
		assert.Contains(t, err.Error(), "provider model cannot be nil")
	})

	// 测试错误的标识符
	t.Run("Invalid Identifier", func(t *testing.T) {
		config := &domain.ProviderAdapterConfig{
			ProviderAdapterInfo: domain.ProviderAdapterInfo{
				Identifier: "invalid",
				Name:       "Test",
				AuthType:   providercore.AuthTypeAPIKey,
			},
			CustomConfig: map[string]interface{}{
				"api_key": "test_key",
			},
		}

		err := template.ValidateConfig(config)
		assert.Error(t, err, "Should return error for invalid identifier")
		assert.Contains(t, err.Error(), "expected provider identifier 'openweathermap'")
	})

	// 测试缺少API Key
	t.Run("Missing API Key", func(t *testing.T) {
		config := &domain.ProviderAdapterConfig{
			ProviderAdapterInfo: domain.ProviderAdapterInfo{
				Identifier: "openweathermap",
				Name:       "OpenWeatherMap",
				AuthType:   providercore.AuthTypeAPIKey,
			},
			CustomConfig: map[string]interface{}{},
		}

		err := template.ValidateConfig(config)
		assert.Error(t, err, "Should return error for missing api_key")
		assert.Contains(t, err.Error(), "api_key is required")
	})

	// 测试空的API Key
	t.Run("Empty API Key", func(t *testing.T) {
		config := &domain.ProviderAdapterConfig{
			ProviderAdapterInfo: domain.ProviderAdapterInfo{
				Identifier: "openweathermap",
				Name:       "OpenWeatherMap",
				AuthType:   providercore.AuthTypeAPIKey,
			},
			CustomConfig: map[string]interface{}{
				"api_key": "",
			},
		}

		err := template.ValidateConfig(config)
		assert.Error(t, err, "Should return error for empty api_key")
		assert.Contains(t, err.Error(), "api_key is required")
	})
}

// TestAdapterCreation 测试适配器创建
func TestAdapterCreation(t *testing.T) {
	template := &OpenWeatherMapAdapterTemplate{}

	// 测试有效配置的适配器创建
	t.Run("Valid Configuration", func(t *testing.T) {
		config := &domain.ProviderAdapterConfig{
			ProviderAdapterInfo: domain.ProviderAdapterInfo{
				Identifier:  "openweathermap",
				Name:        "OpenWeatherMap",
				Description: "Weather API provider",
				AuthType:    providercore.AuthTypeAPIKey,
				Status:      providercore.ProviderStatusActive,
				Operations: []providercore.Operation{
					{Identifier: "get_current_weather"},
					{Identifier: "get_weather_forecast"},
					{Identifier: "get_one_call_weather"},
					{Identifier: "get_air_pollution"},
					{Identifier: "get_geocoding"},
				},
			},
			CustomConfig: map[string]interface{}{
				"api_key": "test_api_key",
			},
		}

		adapter, err := template.CreateAdapter(config)
		assert.NoError(t, err, "Should create adapter without error")
		assert.NotNil(t, adapter, "Created adapter should not be nil")

		// 验证创建的适配器类型
		weatherAdapter, ok := adapter.(*OpenWeatherMapAdapter)
		assert.True(t, ok, "Created adapter should be OpenWeatherMapAdapter")
		assert.NotNil(t, weatherAdapter, "Weather adapter should not be nil")

		// 验证适配器的基本信息
		providerInfo := weatherAdapter.GetProviderAdapterInfo()
		assert.Equal(t, "openweathermap", providerInfo.Identifier)
		assert.Equal(t, "OpenWeatherMap", providerInfo.Name)
		assert.Equal(t, providercore.AuthTypeAPIKey, providerInfo.AuthType)
	})

	// 测试无效配置的适配器创建
	t.Run("Invalid Configuration", func(t *testing.T) {
		config := &domain.ProviderAdapterConfig{
			ProviderAdapterInfo: domain.ProviderAdapterInfo{
				Identifier: "invalid",
				Name:       "Test",
				AuthType:   providercore.AuthTypeAPIKey,
			},
			CustomConfig: map[string]interface{}{
				"api_key": "test_key",
			},
		}

		adapter, err := template.CreateAdapter(config)
		assert.Error(t, err, "Should return error for invalid configuration")
		assert.Nil(t, adapter, "Adapter should be nil on error")
	})
}

// TestAdapterCapabilities 测试适配器能力
func TestAdapterCapabilities(t *testing.T) {
	template := &OpenWeatherMapAdapterTemplate{}

	capabilities := template.GetCapabilities()
	assert.NotEmpty(t, capabilities, "Should have capabilities")

	// 验证关键能力
	assert.True(t, capabilities["weather_data"].(bool), "Should support weather data")
	assert.True(t, capabilities["forecasts"].(bool), "Should support forecasts")
	assert.True(t, capabilities["air_pollution"].(bool), "Should support air pollution data")
	assert.True(t, capabilities["geocoding"].(bool), "Should support geocoding")
	assert.True(t, capabilities["real_time_data"].(bool), "Should support real-time data")
	assert.True(t, capabilities["global_coverage"].(bool), "Should have global coverage")

	// 验证速率限制信息
	rateLimits, ok := capabilities["rate_limits"].(map[string]interface{})
	assert.True(t, ok, "Should have rate limits information")
	assert.NotEmpty(t, rateLimits, "Rate limits should not be empty")
}

// TestAdapterTags 测试适配器标签
func TestAdapterTags(t *testing.T) {
	template := &OpenWeatherMapAdapterTemplate{}

	tags := template.GetTags()
	assert.NotEmpty(t, tags, "Should have tags")

	expectedTags := []string{"weather", "forecast", "climate", "environment", "api", "rest"}
	for _, expectedTag := range expectedTags {
		assert.Contains(t, tags, expectedTag, "Should contain tag: %s", expectedTag)
	}
}

// TestListAllAdapters 测试列出所有已注册的适配器
func TestListAllAdapters(t *testing.T) {
	adapters := registry.ListAdapterTemplates()
	assert.NotEmpty(t, adapters, "Should have registered adapters")
	assert.Contains(t, adapters, "openweathermap", "Should contain openweathermap adapter")

	t.Logf("已注册的适配器: %v", adapters)
}
