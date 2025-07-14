package openweathermap

import (
	"fmt"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/registry"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/rest"
	providercore "github.com/context-space/context-space/backend/internal/providercore/domain"
)

const (
	identifier = "openweathermap"
	baseURL    = "https://api.openweathermap.org"
)

// 在包初始化时注册适配器模板
func init() {
	// 类型断言确保适配器实现了必要的接口
	var _ domain.APIKeyAdapter = (*OpenWeatherMapAdapter)(nil)

	template := &OpenWeatherMapAdapterTemplate{}
	registry.RegisterAdapterTemplate(identifier, template)
}

// OpenWeatherMapAdapterTemplate 实现 AdapterTemplate 接口
type OpenWeatherMapAdapterTemplate struct {
	// 如果需要，可以在这里添加此模板特定的配置
}

// CreateAdapter 从提供的配置创建新的适配器实例
func (t *OpenWeatherMapAdapterTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
	// 验证配置
	if err := t.ValidateConfig(provider); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// 创建提供者信息
	providerInfo := &domain.ProviderAdapterInfo{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
		Permissions: provider.Permissions,
		Operations:  provider.Operations,
		Status:      provider.Status,
	}

	// 创建适配器配置
	adapterConfig := &domain.AdapterConfig{
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryBackoff: 1 * time.Second,
	}

	// 创建 REST 配置
	restConfig := &rest.RESTConfig{
		BaseURL: baseURL,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// 创建底层 REST 适配器
	restAdapter := rest.NewRESTAdapter(providerInfo, adapterConfig, restConfig)

	// 创建主适配器
	adapter := NewOpenWeatherMapAdapter(
		providerInfo,
		adapterConfig,
		restAdapter,
	)

	return adapter, nil
}

// ValidateConfig 检查提供的配置是否包含必要的字段
func (t *OpenWeatherMapAdapterTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
	if provider == nil {
		return fmt.Errorf("provider model cannot be nil")
	}

	// 验证基本字段
	if provider.Identifier == "" {
		return fmt.Errorf("provider identifier is required")
	}

	if provider.Identifier != identifier {
		return fmt.Errorf("expected provider identifier '%s', got '%s'", identifier, provider.Identifier)
	}

	if provider.Name == "" {
		return fmt.Errorf("provider name is required")
	}

	if provider.AuthType == "" {
		return fmt.Errorf("auth type is required")
	}

	if provider.AuthType != providercore.AuthTypeAPIKey {
		return fmt.Errorf("expected auth type 'apikey', got '%s'", provider.AuthType)
	}

	// 验证 API Key 配置
	apiKey, ok := provider.CustomConfig["api_key"]
	if !ok {
		return fmt.Errorf("api_key is required")
	}

	if apiKey.(string) == "" {
		return fmt.Errorf("api_key is required")
	}

	// 验证操作
	if len(provider.Operations) == 0 {
		return fmt.Errorf("at least one operation must be defined")
	}

	// 验证必要的操作
	requiredOperations := []string{
		"get_current_weather",
		"get_weather_forecast",
		"get_one_call_weather",
		"get_air_pollution",
		"get_geocoding",
	}

	operationMap := make(map[string]bool)
	for _, op := range provider.Operations {
		operationMap[op.Identifier] = true
	}

	for _, required := range requiredOperations {
		if !operationMap[required] {
			return fmt.Errorf("required operation '%s' is missing", required)
		}
	}

	return nil
}

// GetIdentifier 返回适配器的标识符
func (t *OpenWeatherMapAdapterTemplate) GetIdentifier() string {
	return identifier
}

// GetName 返回适配器的名称
func (t *OpenWeatherMapAdapterTemplate) GetName() string {
	return "OpenWeatherMap"
}

// GetDescription 返回适配器的描述
func (t *OpenWeatherMapAdapterTemplate) GetDescription() string {
	return "Global weather data provider with current weather, forecasts, and historical data"
}

// GetAuthType 返回认证类型
func (t *OpenWeatherMapAdapterTemplate) GetAuthType() string {
	return "apikey"
}

// GetVersion 返回适配器版本
func (t *OpenWeatherMapAdapterTemplate) GetVersion() string {
	return "1.0.0"
}

// GetSupportedOperations 返回支持的操作列表
func (t *OpenWeatherMapAdapterTemplate) GetSupportedOperations() []string {
	return []string{
		"get_current_weather",
		"get_weather_forecast",
		"get_one_call_weather",
		"get_air_pollution",
		"get_geocoding",
	}
}

// GetRequiredCredentials 返回所需的凭证字段
func (t *OpenWeatherMapAdapterTemplate) GetRequiredCredentials() []string {
	return []string{
		"api_key",
	}
}

// GetCapabilities 返回适配器的能力
func (t *OpenWeatherMapAdapterTemplate) GetCapabilities() map[string]interface{} {
	return map[string]interface{}{
		"weather_data":       true,
		"forecasts":          true,
		"air_pollution":      true,
		"geocoding":          true,
		"historical_data":    true,
		"real_time_data":     true,
		"global_coverage":    true,
		"multiple_languages": true,
		"rate_limits": map[string]interface{}{
			"free_tier": map[string]interface{}{
				"requests_per_minute": 60,
				"requests_per_day":    1000,
			},
		},
	}
}

// GetTags 返回适配器的标签
func (t *OpenWeatherMapAdapterTemplate) GetTags() []string {
	return []string{
		"weather",
		"forecast",
		"climate",
		"environment",
		"api",
		"rest",
	}
}

// IsProduction 返回是否为生产就绪状态
func (t *OpenWeatherMapAdapterTemplate) IsProduction() bool {
	return true
}

// GetDocumentationURL 返回文档URL
func (t *OpenWeatherMapAdapterTemplate) GetDocumentationURL() string {
	return "https://openweathermap.org/api"
}

// GetHomepageURL 返回主页URL
func (t *OpenWeatherMapAdapterTemplate) GetHomepageURL() string {
	return "https://openweathermap.org"
}
