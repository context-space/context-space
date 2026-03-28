// Package llm provides a configurable LLM provider abstraction.
// It supports OpenAI and MiniMax (via OpenAI-compatible API) as backends.
package llm

import (
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// ProviderType identifies which LLM backend to use.
type ProviderType string

const (
	ProviderOpenAI  ProviderType = "openai"
	ProviderMiniMax ProviderType = "minimax"
)

// ProviderConfig holds the configuration for a specific LLM provider.
type ProviderConfig struct {
	Name             string
	DefaultModel     string
	BaseURL          string
	APIKey           string
	MinTemperature   float64
	MaxTemperature   float64
}

// DetectProvider determines the active LLM provider from environment variables.
// Priority: LLM_PROVIDER env var > auto-detect from API keys > fallback to OpenAI.
func DetectProvider() ProviderType {
	if explicit := os.Getenv("LLM_PROVIDER"); explicit != "" {
		switch strings.ToLower(explicit) {
		case "minimax":
			return ProviderMiniMax
		case "openai":
			return ProviderOpenAI
		}
	}

	// Auto-detect: prefer MiniMax if its key is set and OpenAI's is not
	if os.Getenv("MINIMAX_API_KEY") != "" && os.Getenv("OPENAI_API_KEY") == "" {
		return ProviderMiniMax
	}

	return ProviderOpenAI
}

// GetProviderConfig returns the configuration for the given provider type.
// If overrideAPIKey or overrideBaseURL are non-empty, they take precedence.
func GetProviderConfig(provider ProviderType, overrideAPIKey, overrideBaseURL string) ProviderConfig {
	switch provider {
	case ProviderMiniMax:
		apiKey := overrideAPIKey
		if apiKey == "" {
			apiKey = os.Getenv("MINIMAX_API_KEY")
		}
		baseURL := overrideBaseURL
		if baseURL == "" {
			baseURL = envOrDefault("MINIMAX_BASE_URL", "https://api.minimax.io/v1")
		}
		return ProviderConfig{
			Name:           "MiniMax",
			DefaultModel:   "MiniMax-M2.7",
			BaseURL:        baseURL,
			APIKey:         apiKey,
			MinTemperature: 0.01,
			MaxTemperature: 1.0,
		}
	default: // openai
		apiKey := overrideAPIKey
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		baseURL := overrideBaseURL
		if baseURL == "" {
			baseURL = envOrDefault("OPENAI_BASE_URL", "https://api.openai.com/v1")
		}
		return ProviderConfig{
			Name:           "OpenAI",
			DefaultModel:   "gpt-4o",
			BaseURL:        baseURL,
			APIKey:         apiKey,
			MinTemperature: 0.0,
			MaxTemperature: 2.0,
		}
	}
}

// ClampTemperature constrains a temperature value to the provider's supported range.
func ClampTemperature(temp float64, provider ProviderType) float64 {
	cfg := GetProviderConfig(provider, "", "")
	if temp < cfg.MinTemperature {
		return cfg.MinTemperature
	}
	if temp > cfg.MaxTemperature {
		return cfg.MaxTemperature
	}
	return temp
}

// NewClient creates an openai.Client configured for the given provider.
// MiniMax exposes an OpenAI-compatible API so the same SDK works for both.
func NewClient(provider ProviderType, apiKey, baseURL string) (*openai.Client, error) {
	cfg := GetProviderConfig(provider, apiKey, baseURL)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("llm: API key is required for provider %s", cfg.Name)
	}

	opts := []option.RequestOption{
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
	}

	client := openai.NewClient(opts...)
	return &client, nil
}

// ModelOrDefault returns the model name if non-empty, otherwise the provider's default.
func ModelOrDefault(model string, provider ProviderType) string {
	if model != "" {
		return model
	}
	cfg := GetProviderConfig(provider, "", "")
	return cfg.DefaultModel
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
