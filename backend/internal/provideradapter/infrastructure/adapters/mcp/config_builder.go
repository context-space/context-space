package mcp

import (
	"fmt"
	"strings"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
)

// DummyValues holds dummy values for credentials and parameters
type DummyValues struct {
	Credentials map[string]string
	Parameters  map[string]string
}

// ValueResolver handles extraction and resolution of credential and parameter values
type ValueResolver struct {
	credentials map[string]string
	parameters  map[string]string
	dummies     *DummyValues
}

// NewValueResolver creates a new value resolver
func NewValueResolver(credential interface{}, parameters map[string]interface{}, dummies *DummyValues) *ValueResolver {
	resolver := &ValueResolver{
		credentials: extractCredentialValues(credential),
		parameters:  convertParametersToStrings(parameters),
		dummies:     dummies,
	}
	return resolver
}

// GetCredentialValue returns credential value with dummy fallback
func (vr *ValueResolver) GetCredentialValue(field string) (string, bool) {
	// Try real credential first
	if value, exists := vr.credentials[field]; exists {
		return value, true
	}

	// Fall back to dummy value
	if vr.dummies != nil && vr.dummies.Credentials != nil {
		if dummyValue, exists := vr.dummies.Credentials[field]; exists {
			return dummyValue, true
		}
	}

	return "", false
}

// GetParameterValue returns parameter value with dummy fallback
func (vr *ValueResolver) GetParameterValue(field string) (string, bool) {
	// Try real parameter first
	if value, exists := vr.parameters[field]; exists {
		return value, true
	}

	// Fall back to dummy value
	if vr.dummies != nil && vr.dummies.Parameters != nil {
		if dummyValue, exists := vr.dummies.Parameters[field]; exists {
			return dummyValue, true
		}
	}

	return "", false
}

// MappingApplier handles applying mappings to MCP client configuration
type MappingApplier struct {
	config *MCPClientConfig
}

// NewMappingApplier creates a new mapping applier
func NewMappingApplier(config *MCPClientConfig) *MappingApplier {
	return &MappingApplier{config: config}
}

// ApplyMapping applies a target mapping with the given value
func (ma *MappingApplier) ApplyMapping(target, value string) {
	if strings.HasPrefix(target, "env:") {
		ma.applyEnvironmentMapping(target, value)
	} else if strings.HasPrefix(target, "arg:") {
		ma.applyArgumentMapping(target, value)
	}
}

// applyEnvironmentMapping applies environment variable mapping
func (ma *MappingApplier) applyEnvironmentMapping(target, value string) {
	envVar := strings.TrimPrefix(target, "env:")
	ma.config.Envs[envVar] = value
}

// applyArgumentMapping applies argument mapping with smart parsing
func (ma *MappingApplier) applyArgumentMapping(target, value string) {
	argPattern := strings.TrimPrefix(target, "arg:")

	// Handle placeholder replacement
	result := strings.ReplaceAll(argPattern, "${value}", value)

	// Try to replace existing argument first
	if ma.replaceExistingArgument(argPattern, value) {
		return
	}

	// Add new argument(s) with smart parsing
	ma.addNewArguments(result)
}

// replaceExistingArgument tries to replace an existing argument
func (ma *MappingApplier) replaceExistingArgument(pattern, value string) bool {
	for i, arg := range ma.config.Args {
		if arg == pattern {
			ma.config.Args[i] = value
			return true
		}
	}
	return false
}

// addNewArguments adds new arguments with smart parsing
func (ma *MappingApplier) addNewArguments(result string) {
	if strings.Contains(result, "=") {
		// Key-value format: --config=/path/to/config
		ma.config.Args = append(ma.config.Args, result)
	} else if strings.Contains(result, " ") {
		// Parameter pair format: --target-store store_value
		parts := strings.Fields(result)
		ma.config.Args = append(ma.config.Args, parts...)
	} else {
		// Single argument format: /path/to/file
		ma.config.Args = append(ma.config.Args, result)
	}
}

// MCPConfigBuilder handles building MCP client configuration
type MCPConfigBuilder struct {
	baseConfig *MCPAdapterConfig
}

// NewMCPConfigBuilder creates a new configuration builder
func NewMCPConfigBuilder(baseConfig *MCPAdapterConfig) *MCPConfigBuilder {
	return &MCPConfigBuilder{baseConfig: baseConfig}
}

// Build creates MCP client configuration with credential and parameter injection
func (builder *MCPConfigBuilder) Build(credential interface{}, parameters map[string]interface{}) MCPClientConfig {
	// Start with base configuration
	config := builder.createBaseConfig()

	// Create dummy values
	dummies := &DummyValues{
		Credentials: builder.baseConfig.DummyCredentials,
		Parameters:  builder.baseConfig.DummyParameters,
	}

	// Create value resolver
	resolver := NewValueResolver(credential, parameters, dummies)

	// Create mapping applier
	applier := NewMappingApplier(&config)

	// Apply credential mappings
	builder.applyCredentialMappings(resolver, applier)

	// Apply parameter mappings
	builder.applyParameterMappings(resolver, applier)

	return config
}

// createBaseConfig creates the base configuration by copying from adapter config
func (builder *MCPConfigBuilder) createBaseConfig() MCPClientConfig {
	config := MCPClientConfig{
		Command: builder.baseConfig.Command,
		Args:    make([]string, len(builder.baseConfig.Args)),
		Envs:    make(map[string]string),
		Timeout: builder.baseConfig.Timeout,
	}

	// Copy server args
	copy(config.Args, builder.baseConfig.Args)

	// Copy server environment
	for k, v := range builder.baseConfig.Envs {
		config.Envs[k] = v
	}

	return config
}

// applyCredentialMappings applies all credential mappings
func (builder *MCPConfigBuilder) applyCredentialMappings(resolver *ValueResolver, applier *MappingApplier) {
	for credField, target := range builder.baseConfig.CredentialMappings {
		if value, exists := resolver.GetCredentialValue(credField); exists {
			applier.ApplyMapping(target, value)
		}
	}
}

// applyParameterMappings applies all parameter mappings
func (builder *MCPConfigBuilder) applyParameterMappings(resolver *ValueResolver, applier *MappingApplier) {
	for paramField, target := range builder.baseConfig.ParameterMappings {
		if value, exists := resolver.GetParameterValue(paramField); exists {
			applier.ApplyMapping(target, value)
		}
	}
}

// Helper functions

// extractCredentialValues extracts values from credential object
func extractCredentialValues(credential interface{}) map[string]string {
	values := make(map[string]string)

	if credential == nil {
		return values
	}

	switch cred := credential.(type) {
	case *credDomain.APIKeyCredential:
		values["apikey"] = cred.APIKey
	case *credDomain.OAuthCredential:
		if cred.Token != nil {
			values["access_token"] = cred.Token.AccessToken
			values["refresh_token"] = cred.Token.RefreshToken
		}
	case *credDomain.BasicAuthCredential:
		values["username"] = cred.Username
		values["password"] = cred.Password
	}

	return values
}

// convertParametersToStrings converts parameters to string values
func convertParametersToStrings(parameters map[string]interface{}) map[string]string {
	result := make(map[string]string)

	if parameters == nil {
		return result
	}

	for key, value := range parameters {
		result[key] = fmt.Sprintf("%v", value)
	}

	return result
}
