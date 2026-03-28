package llm

import (
	"os"
	"testing"
)

func clearLLMEnv() {
	os.Unsetenv("LLM_PROVIDER")
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("OPENAI_BASE_URL")
	os.Unsetenv("MINIMAX_API_KEY")
	os.Unsetenv("MINIMAX_BASE_URL")
}

func TestDetectProvider(t *testing.T) {
	t.Run("defaults to openai", func(t *testing.T) {
		clearLLMEnv()
		if got := DetectProvider(); got != ProviderOpenAI {
			t.Errorf("DetectProvider() = %v, want %v", got, ProviderOpenAI)
		}
	})

	t.Run("explicit LLM_PROVIDER=minimax", func(t *testing.T) {
		clearLLMEnv()
		os.Setenv("LLM_PROVIDER", "minimax")
		defer os.Unsetenv("LLM_PROVIDER")
		if got := DetectProvider(); got != ProviderMiniMax {
			t.Errorf("DetectProvider() = %v, want %v", got, ProviderMiniMax)
		}
	})

	t.Run("explicit LLM_PROVIDER=openai", func(t *testing.T) {
		clearLLMEnv()
		os.Setenv("LLM_PROVIDER", "openai")
		defer os.Unsetenv("LLM_PROVIDER")
		if got := DetectProvider(); got != ProviderOpenAI {
			t.Errorf("DetectProvider() = %v, want %v", got, ProviderOpenAI)
		}
	})

	t.Run("case insensitive LLM_PROVIDER", func(t *testing.T) {
		clearLLMEnv()
		os.Setenv("LLM_PROVIDER", "MiniMax")
		defer os.Unsetenv("LLM_PROVIDER")
		if got := DetectProvider(); got != ProviderMiniMax {
			t.Errorf("DetectProvider() = %v, want %v", got, ProviderMiniMax)
		}
	})

	t.Run("auto-detect minimax from API key", func(t *testing.T) {
		clearLLMEnv()
		os.Setenv("MINIMAX_API_KEY", "test-key")
		defer os.Unsetenv("MINIMAX_API_KEY")
		if got := DetectProvider(); got != ProviderMiniMax {
			t.Errorf("DetectProvider() = %v, want %v", got, ProviderMiniMax)
		}
	})

	t.Run("prefers openai when both keys set", func(t *testing.T) {
		clearLLMEnv()
		os.Setenv("OPENAI_API_KEY", "openai-key")
		os.Setenv("MINIMAX_API_KEY", "minimax-key")
		defer func() {
			os.Unsetenv("OPENAI_API_KEY")
			os.Unsetenv("MINIMAX_API_KEY")
		}()
		if got := DetectProvider(); got != ProviderOpenAI {
			t.Errorf("DetectProvider() = %v, want %v", got, ProviderOpenAI)
		}
	})

	t.Run("LLM_PROVIDER overrides auto-detection", func(t *testing.T) {
		clearLLMEnv()
		os.Setenv("LLM_PROVIDER", "minimax")
		os.Setenv("OPENAI_API_KEY", "openai-key")
		defer func() {
			os.Unsetenv("LLM_PROVIDER")
			os.Unsetenv("OPENAI_API_KEY")
		}()
		if got := DetectProvider(); got != ProviderMiniMax {
			t.Errorf("DetectProvider() = %v, want %v", got, ProviderMiniMax)
		}
	})
}

func TestGetProviderConfig(t *testing.T) {
	t.Run("openai config defaults", func(t *testing.T) {
		clearLLMEnv()
		cfg := GetProviderConfig(ProviderOpenAI, "", "")
		if cfg.Name != "OpenAI" {
			t.Errorf("Name = %v, want OpenAI", cfg.Name)
		}
		if cfg.DefaultModel != "gpt-4o" {
			t.Errorf("DefaultModel = %v, want gpt-4o", cfg.DefaultModel)
		}
		if cfg.BaseURL != "https://api.openai.com/v1" {
			t.Errorf("BaseURL = %v, want https://api.openai.com/v1", cfg.BaseURL)
		}
	})

	t.Run("minimax config defaults", func(t *testing.T) {
		clearLLMEnv()
		cfg := GetProviderConfig(ProviderMiniMax, "", "")
		if cfg.Name != "MiniMax" {
			t.Errorf("Name = %v, want MiniMax", cfg.Name)
		}
		if cfg.DefaultModel != "MiniMax-M2.7" {
			t.Errorf("DefaultModel = %v, want MiniMax-M2.7", cfg.DefaultModel)
		}
		if cfg.BaseURL != "https://api.minimax.io/v1" {
			t.Errorf("BaseURL = %v, want https://api.minimax.io/v1", cfg.BaseURL)
		}
	})

	t.Run("override API key", func(t *testing.T) {
		clearLLMEnv()
		cfg := GetProviderConfig(ProviderMiniMax, "custom-key", "")
		if cfg.APIKey != "custom-key" {
			t.Errorf("APIKey = %v, want custom-key", cfg.APIKey)
		}
	})

	t.Run("override base URL", func(t *testing.T) {
		clearLLMEnv()
		cfg := GetProviderConfig(ProviderMiniMax, "", "https://custom.api.io/v1")
		if cfg.BaseURL != "https://custom.api.io/v1" {
			t.Errorf("BaseURL = %v, want https://custom.api.io/v1", cfg.BaseURL)
		}
	})

	t.Run("env var API key fallback for minimax", func(t *testing.T) {
		clearLLMEnv()
		os.Setenv("MINIMAX_API_KEY", "env-minimax-key")
		defer os.Unsetenv("MINIMAX_API_KEY")
		cfg := GetProviderConfig(ProviderMiniMax, "", "")
		if cfg.APIKey != "env-minimax-key" {
			t.Errorf("APIKey = %v, want env-minimax-key", cfg.APIKey)
		}
	})

	t.Run("env var base URL fallback for minimax", func(t *testing.T) {
		clearLLMEnv()
		os.Setenv("MINIMAX_BASE_URL", "https://env.minimax.io/v1")
		defer os.Unsetenv("MINIMAX_BASE_URL")
		cfg := GetProviderConfig(ProviderMiniMax, "", "")
		if cfg.BaseURL != "https://env.minimax.io/v1" {
			t.Errorf("BaseURL = %v, want https://env.minimax.io/v1", cfg.BaseURL)
		}
	})

	t.Run("minimax temperature range", func(t *testing.T) {
		clearLLMEnv()
		cfg := GetProviderConfig(ProviderMiniMax, "", "")
		if cfg.MinTemperature != 0.01 {
			t.Errorf("MinTemperature = %v, want 0.01", cfg.MinTemperature)
		}
		if cfg.MaxTemperature != 1.0 {
			t.Errorf("MaxTemperature = %v, want 1.0", cfg.MaxTemperature)
		}
	})

	t.Run("openai temperature range", func(t *testing.T) {
		clearLLMEnv()
		cfg := GetProviderConfig(ProviderOpenAI, "", "")
		if cfg.MinTemperature != 0.0 {
			t.Errorf("MinTemperature = %v, want 0.0", cfg.MinTemperature)
		}
		if cfg.MaxTemperature != 2.0 {
			t.Errorf("MaxTemperature = %v, want 2.0", cfg.MaxTemperature)
		}
	})
}

func TestClampTemperature(t *testing.T) {
	tests := []struct {
		name     string
		temp     float64
		provider ProviderType
		want     float64
	}{
		{"minimax clamp low", 0.0, ProviderMiniMax, 0.01},
		{"minimax clamp high", 2.0, ProviderMiniMax, 1.0},
		{"minimax pass through", 0.7, ProviderMiniMax, 0.7},
		{"minimax at min boundary", 0.01, ProviderMiniMax, 0.01},
		{"minimax at max boundary", 1.0, ProviderMiniMax, 1.0},
		{"openai allow zero", 0.0, ProviderOpenAI, 0.0},
		{"openai clamp high", 3.0, ProviderOpenAI, 2.0},
		{"openai pass through", 1.5, ProviderOpenAI, 1.5},
		{"negative temperature minimax", -1.0, ProviderMiniMax, 0.01},
		{"negative temperature openai", -1.0, ProviderOpenAI, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearLLMEnv()
			got := ClampTemperature(tt.temp, tt.provider)
			if got != tt.want {
				t.Errorf("ClampTemperature(%v, %v) = %v, want %v", tt.temp, tt.provider, got, tt.want)
			}
		})
	}
}

func TestModelOrDefault(t *testing.T) {
	t.Run("returns custom model when set", func(t *testing.T) {
		clearLLMEnv()
		if got := ModelOrDefault("custom-model", ProviderMiniMax); got != "custom-model" {
			t.Errorf("ModelOrDefault() = %v, want custom-model", got)
		}
	})

	t.Run("returns minimax default when empty", func(t *testing.T) {
		clearLLMEnv()
		if got := ModelOrDefault("", ProviderMiniMax); got != "MiniMax-M2.7" {
			t.Errorf("ModelOrDefault() = %v, want MiniMax-M2.7", got)
		}
	})

	t.Run("returns openai default when empty", func(t *testing.T) {
		clearLLMEnv()
		if got := ModelOrDefault("", ProviderOpenAI); got != "gpt-4o" {
			t.Errorf("ModelOrDefault() = %v, want gpt-4o", got)
		}
	})
}

func TestNewClient(t *testing.T) {
	t.Run("creates openai client", func(t *testing.T) {
		clearLLMEnv()
		client, err := NewClient(ProviderOpenAI, "test-api-key", "")
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client == nil {
			t.Fatal("NewClient() returned nil")
		}
	})

	t.Run("creates minimax client", func(t *testing.T) {
		clearLLMEnv()
		client, err := NewClient(ProviderMiniMax, "test-api-key", "")
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client == nil {
			t.Fatal("NewClient() returned nil")
		}
	})

	t.Run("creates minimax client with custom base URL", func(t *testing.T) {
		clearLLMEnv()
		client, err := NewClient(ProviderMiniMax, "test-api-key", "https://custom.minimax.io/v1")
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client == nil {
			t.Fatal("NewClient() returned nil")
		}
	})

	t.Run("errors on missing API key", func(t *testing.T) {
		clearLLMEnv()
		_, err := NewClient(ProviderMiniMax, "", "")
		if err == nil {
			t.Fatal("NewClient() expected error for missing API key")
		}
	})

	t.Run("uses env var API key for minimax", func(t *testing.T) {
		clearLLMEnv()
		os.Setenv("MINIMAX_API_KEY", "env-key")
		defer os.Unsetenv("MINIMAX_API_KEY")
		client, err := NewClient(ProviderMiniMax, "", "")
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client == nil {
			t.Fatal("NewClient() returned nil")
		}
	})
}
