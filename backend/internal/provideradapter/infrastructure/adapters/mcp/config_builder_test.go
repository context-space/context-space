package mcp

import (
	"testing"
	"time"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
)

func TestMCPConfigBuilder(t *testing.T) {
	// Test configuration with correct argument mapping
	baseConfig := &MCPAdapterConfig{
		Command: "npx",
		Args:    []string{"-y", "test-package", "--target", "PLACEHOLDER"},
		Envs:    map[string]string{"BASE_ENV": "base_value"},
		Timeout: 30 * time.Second,
		CredentialMappings: map[string]string{
			"apikey": "env:API_KEY",
		},
		DummyCredentials: map[string]string{
			"apikey": "dummy_key",
		},
		ParameterMappings: map[string]string{
			"target": "arg:PLACEHOLDER", // Replace PLACEHOLDER with parameter value
		},
		DummyParameters: map[string]string{
			"target": "dummy_target",
		},
	}

	builder := NewMCPConfigBuilder(baseConfig)

	t.Run("BuildWithNoCredentialOrParams", func(t *testing.T) {
		config := builder.Build(nil, nil)

		// Should use dummy values
		if config.Envs["API_KEY"] != "dummy_key" {
			t.Errorf("Expected dummy credential, got: %s", config.Envs["API_KEY"])
		}

		// Should replace PLACEHOLDER with dummy parameter value
		found := false
		for _, arg := range config.Args {
			if arg == "dummy_target" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected PLACEHOLDER to be replaced with dummy_target, got: %v", config.Args)
		}
	})

	t.Run("BuildWithRealCredential", func(t *testing.T) {
		cred := &credDomain.APIKeyCredential{
			APIKey: "real_api_key",
		}

		config := builder.Build(cred, nil)

		// Should use real credential
		if config.Envs["API_KEY"] != "real_api_key" {
			t.Errorf("Expected real credential, got: %s", config.Envs["API_KEY"])
		}
	})

	t.Run("BuildWithRealParameters", func(t *testing.T) {
		params := map[string]interface{}{
			"target": "real_target",
		}

		config := builder.Build(nil, params)

		// Should replace PLACEHOLDER with real parameter value
		found := false
		for _, arg := range config.Args {
			if arg == "real_target" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected PLACEHOLDER to be replaced with real_target, got: %v", config.Args)
		}
	})
}

func TestValueResolver(t *testing.T) {
	dummies := &DummyValues{
		Credentials: map[string]string{
			"apikey": "dummy_key",
		},
		Parameters: map[string]string{
			"target": "dummy_target",
		},
	}

	t.Run("CredentialValueWithDummy", func(t *testing.T) {
		resolver := NewValueResolver(nil, nil, dummies)

		value, exists := resolver.GetCredentialValue("apikey")
		if !exists {
			t.Error("Expected credential value to exist")
		}
		if value != "dummy_key" {
			t.Errorf("Expected dummy_key, got: %s", value)
		}
	})

	t.Run("ParameterValueWithDummy", func(t *testing.T) {
		resolver := NewValueResolver(nil, nil, dummies)

		value, exists := resolver.GetParameterValue("target")
		if !exists {
			t.Error("Expected parameter value to exist")
		}
		if value != "dummy_target" {
			t.Errorf("Expected dummy_target, got: %s", value)
		}
	})

	t.Run("RealCredentialOverridesDummy", func(t *testing.T) {
		cred := &credDomain.APIKeyCredential{
			APIKey: "real_key",
		}
		resolver := NewValueResolver(cred, nil, dummies)

		value, exists := resolver.GetCredentialValue("apikey")
		if !exists {
			t.Error("Expected credential value to exist")
		}
		if value != "real_key" {
			t.Errorf("Expected real_key, got: %s", value)
		}
	})

	t.Run("RealParameterOverridesDummy", func(t *testing.T) {
		params := map[string]interface{}{
			"target": "real_target",
		}
		resolver := NewValueResolver(nil, params, dummies)

		value, exists := resolver.GetParameterValue("target")
		if !exists {
			t.Error("Expected parameter value to exist")
		}
		if value != "real_target" {
			t.Errorf("Expected real_target, got: %s", value)
		}
	})
}

func TestMappingApplier(t *testing.T) {
	config := &MCPClientConfig{
		Command: "test",
		Args:    []string{"arg1", "--replace-me", "arg3"},
		Envs:    make(map[string]string),
		Timeout: 30 * time.Second,
	}

	applier := NewMappingApplier(config)

	t.Run("EnvironmentMapping", func(t *testing.T) {
		applier.ApplyMapping("env:TEST_VAR", "test_value")

		if config.Envs["TEST_VAR"] != "test_value" {
			t.Errorf("Expected test_value, got: %s", config.Envs["TEST_VAR"])
		}
	})

	t.Run("ArgumentReplacement", func(t *testing.T) {
		applier.ApplyMapping("arg:--replace-me", "replaced")

		if config.Args[1] != "replaced" {
			t.Errorf("Expected replaced, got: %s", config.Args[1])
		}
	})

	t.Run("ArgumentAddition", func(t *testing.T) {
		applier.ApplyMapping("arg:--new-arg", "new_value")

		// Should be added as new arguments (--new-arg)
		if config.Args[len(config.Args)-1] != "--new-arg" {
			t.Errorf("Expected --new-arg to be added, got: %v", config.Args)
		}
	})
}
