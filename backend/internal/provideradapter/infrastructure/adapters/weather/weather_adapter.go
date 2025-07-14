package openweathermap

import (
	"context"
	"fmt"
	"net/http"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/base"
	"github.com/context-space/context-space/backend/internal/shared/utils"
)

// OpenWeatherMapAdapter 是 OpenWeatherMap API 的适配器，使用 API Key 认证
type OpenWeatherMapAdapter struct {
	*base.BaseAdapter
	restAdapter domain.Adapter // 底层 REST 适配器实例
	operations  Operations     // 操作 ID 到定义的映射
}

// NewOpenWeatherMapAdapter 创建新的适配器实例
func NewOpenWeatherMapAdapter(
	providerInfo *domain.ProviderAdapterInfo,
	config *domain.AdapterConfig,
	restAdapter domain.Adapter,
) *OpenWeatherMapAdapter {
	baseAdapter := base.NewBaseAdapter(providerInfo, config)

	adapter := &OpenWeatherMapAdapter{
		BaseAdapter: baseAdapter,
		restAdapter: restAdapter,
		operations:  make(Operations),
	}

	// 注册操作
	adapter.registerOperations()
	return adapter
}

// Execute 根据 operationID 处理 API 调用，使用 REST 适配器
func (a *OpenWeatherMapAdapter) Execute(
	ctx context.Context,
	operationID string,
	params map[string]interface{},
	credential interface{},
) (interface{}, error) {
	// 验证 API Key 凭证
	apiKeyCred, ok := credential.(*credDomain.APIKeyCredential)
	if !ok || apiKeyCred == nil || apiKeyCred.APIKey == "" {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrCredentialError,
			"invalid or missing API key credential",
			http.StatusUnauthorized,
		)
	}

	// 检查操作是否存在
	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrOperationNotSupported,
			fmt.Sprintf("unknown operation ID: %s", operationID),
			http.StatusNotFound,
		)
	}

	// 处理和验证参数
	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrInvalidParameters,
			fmt.Sprintf("parameter validation failed: %v", err),
			http.StatusBadRequest,
		)
	}

	// 调用操作处理函数
	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(
			a.GetProviderAdapterInfo().Identifier,
			operationID,
			domain.ErrInternal,
			fmt.Sprintf("operation handler failed: %v", err),
			http.StatusInternalServerError,
		)
	}

	// 添加 API Key 到查询参数
	queryParams, _ := restParams["query_params"].(map[string]string)
	if queryParams == nil {
		queryParams = make(map[string]string)
	}
	queryParams["appid"] = apiKeyCred.APIKey
	restParams["query_params"] = queryParams

	// 设置请求头
	headers, _ := restParams["headers"].(map[string]string)
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	headers["User-Agent"] = utils.StringsBuilder("OpenWeatherMap-Adapter/1.0")
	restParams["headers"] = headers

	// 调用底层 REST 适配器
	rawResult, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		return nil, a.handleRestError(operationID, err)
	}

	return rawResult, nil
}

// handleRestError 处理 REST 适配器返回的错误
func (a *OpenWeatherMapAdapter) handleRestError(operationID string, err error) error {
	// 如果已经是适配器错误，直接返回
	if adapterErr, ok := err.(*domain.AdapterError); ok {
		return adapterErr
	}

	// 转换为适配器错误
	return domain.NewAdapterError(
		a.GetProviderAdapterInfo().Identifier,
		operationID,
		domain.ErrInternal,
		fmt.Sprintf("REST adapter error: %v", err),
		http.StatusInternalServerError,
	)
}

// GetSupportedOperations 返回支持的操作列表
func (a *OpenWeatherMapAdapter) GetSupportedOperations() []string {
	operations := make([]string, 0, len(a.operations))
	for operationID := range a.operations {
		operations = append(operations, operationID)
	}
	return operations
}

// ValidateCredential 验证凭证是否有效
func (a *OpenWeatherMapAdapter) ValidateCredential(credential interface{}) error {
	apiKeyCred, ok := credential.(*credDomain.APIKeyCredential)
	if !ok {
		return fmt.Errorf("invalid credential type, expected APIKeyCredential")
	}

	if apiKeyCred.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	return nil
}

// GetCredentialType 返回此适配器使用的凭证类型
func (a *OpenWeatherMapAdapter) GetCredentialType() string {
	return "apikey"
}

// HealthCheck 检查适配器和服务的健康状态
func (a *OpenWeatherMapAdapter) HealthCheck(ctx context.Context) error {
	// 可以实现一个简单的健康检查，例如调用一个轻量级的 API
	// 这里简单返回 nil，表示适配器本身是健康的
	// 实际的服务健康检查会在真正的 API 调用时进行
	return nil
}

// GetRateLimits 返回速率限制信息
func (a *OpenWeatherMapAdapter) GetRateLimits() map[string]interface{} {
	return map[string]interface{}{
		"requests_per_minute": 60,
		"requests_per_day":    1000,
		"note":                "Free tier limits, may vary based on subscription plan",
	}
}

// GetMetrics 返回适配器的度量信息
func (a *OpenWeatherMapAdapter) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"total_operations":     len(a.operations),
		"supported_operations": a.GetSupportedOperations(),
		"auth_type":            "apikey",
		"base_url":             "https://api.openweathermap.org",
	}
}
