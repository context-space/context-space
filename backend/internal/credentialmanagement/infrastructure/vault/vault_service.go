package vault

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/hashicorp/vault/api"

	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
)

// VaultConfig contains configuration for the Vault service
type VaultConfig struct {
	Regions map[domain.VaultRegion]VaultRegionalConfig
}

// VaultRegionalConfig contains per-region Vault configuration
type VaultRegionalConfig struct {
	Address     string
	Token       string
	TransitPath string // Base transit path, will be extended with region and credential type
}

// RegionalVaultClient encapsulates a Vault client for a specific region
type RegionalVaultClient struct {
	client          *api.Client
	baseTransitPath string
	region          domain.VaultRegion
}

// VaultServiceImpl implements the domain.VaultService interface using HashiCorp Vault
type VaultServiceImpl struct {
	mu             sync.RWMutex
	clients        map[domain.VaultRegion]*RegionalVaultClient
	defaultRegion  domain.VaultRegion
	algorithmMap   map[domain.VaultRegion]domain.VaultAlgorithm
	keyNamePattern map[domain.CredentialType]string
}

// NewVaultService creates a new VaultService implementation
func NewVaultService(ctx context.Context, config *VaultConfig, defaultRegion domain.VaultRegion) (domain.VaultService, error) {
	if config == nil || len(config.Regions) == 0 {
		return nil, fmt.Errorf("vault configuration is required")
	}

	if _, ok := config.Regions[defaultRegion]; !ok {
		return nil, fmt.Errorf("default region %s must have a configuration", defaultRegion)
	}

	clients := make(map[domain.VaultRegion]*RegionalVaultClient)
	algorithmMap := map[domain.VaultRegion]domain.VaultAlgorithm{
		domain.RegionEU: domain.AlgorithmAESGCM,
		domain.RegionUS: domain.AlgorithmAESGCM,
		domain.RegionCN: domain.AlgorithmAESGCM, // Will support SM4GCM in the future
	}

	keyNamePattern := map[domain.CredentialType]string{
		domain.CredentialTypeOAuth:  "oauth-creds-%s-key",
		domain.CredentialTypeAPIKey: "apikey-creds-%s-key",
	}

	for region, cfg := range config.Regions {
		client, err := createVaultClient(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create Vault client for region %s: %w", region, err)
		}

		clients[region] = &RegionalVaultClient{
			client:          client,
			baseTransitPath: cfg.TransitPath,
			region:          region,
		}
	}

	return &VaultServiceImpl{
		clients:        clients,
		defaultRegion:  defaultRegion,
		algorithmMap:   algorithmMap,
		keyNamePattern: keyNamePattern,
	}, nil
}

// createVaultClient creates a new Vault client with the provided configuration
func createVaultClient(config VaultRegionalConfig) (*api.Client, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = config.Address

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	client.SetToken(config.Token)
	return client, nil
}

// getRegionalClient returns the client for the specified region
func (v *VaultServiceImpl) getRegionalClient(region domain.VaultRegion) (*RegionalVaultClient, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	client, ok := v.clients[region]
	if !ok {
		return nil, fmt.Errorf("no Vault client configured for region %s", region)
	}

	return client, nil
}

// getKeyName constructs the key name based on credential type and region
func (v *VaultServiceImpl) getKeyName(credType domain.CredentialType, region domain.VaultRegion) string {
	pattern, ok := v.keyNamePattern[credType]
	if !ok {
		// Default to a generic pattern if not found
		pattern = "generic-creds-%s-key"
	}

	return fmt.Sprintf(pattern, region)
}

// getTransitPath constructs the transit path for a specific region and credential type
func (v *VaultServiceImpl) getTransitPath(client *RegionalVaultClient, credType domain.CredentialType) string {
	// RefreshToken with the base transit path (usually "transit")
	basePath := client.baseTransitPath

	// Convert credential type to lowercase
	var credTypeStr string
	switch credType {
	case domain.CredentialTypeOAuth:
		credTypeStr = "oauth"
	case domain.CredentialTypeAPIKey:
		credTypeStr = "apikey"
	}

	// Construct the full transit path
	return fmt.Sprintf("%s-%s-%s", basePath, client.region, credTypeStr)
}

// EncryptData encrypts plaintext data for a specific region
func (v *VaultServiceImpl) EncryptData(ctx context.Context, plaintext string, region domain.VaultRegion, credentialType domain.CredentialType) (*domain.EncryptionMetadata, error) {
	client, err := v.getRegionalClient(region)
	if err != nil {
		return nil, err
	}

	keyName := v.getKeyName(credentialType, region)
	transitPath := v.getTransitPath(client, credentialType)
	path := fmt.Sprintf("%s/encrypt/%s", transitPath, keyName)

	data := map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
	}

	secret, err := client.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return nil, fmt.Errorf("vault encryption failed: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data returned from vault")
	}

	ciphertext, ok := secret.Data["ciphertext"].(string)
	if !ok || ciphertext == "" {
		return nil, fmt.Errorf("invalid ciphertext returned from vault")
	}

	// Extract key version from the ciphertext (format: vault:v{version}:{data})
	var keyVersion int
	_, err = fmt.Sscanf(ciphertext, "vault:v%d:", &keyVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key version from ciphertext: %w", err)
	}

	algorithm, ok := v.algorithmMap[region]
	if !ok {
		algorithm = domain.AlgorithmAESGCM // Default algorithm
	}

	metadata := &domain.EncryptionMetadata{
		Region:         region,
		KeyVersion:     keyVersion,
		CredentialType: credentialType,
		Algorithm:      algorithm,
		Ciphertext:     ciphertext,
	}

	return metadata, nil
}

// DecryptData decrypts data using metadata embedded within the provided structure
func (v *VaultServiceImpl) DecryptData(ctx context.Context, metadata *domain.EncryptionMetadata) (string, error) {
	if metadata == nil {
		return "", fmt.Errorf("metadata cannot be nil")
	}

	client, err := v.getRegionalClient(metadata.Region)
	if err != nil {
		return "", err
	}

	keyName := v.getKeyName(metadata.CredentialType, metadata.Region)
	transitPath := v.getTransitPath(client, metadata.CredentialType)
	path := fmt.Sprintf("%s/decrypt/%s", transitPath, keyName)

	data := map[string]interface{}{
		"ciphertext": metadata.Ciphertext,
	}

	secret, err := client.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return "", fmt.Errorf("vault decryption failed: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("no data returned from vault")
	}

	plaintextBase64, ok := secret.Data["plaintext"].(string)
	if !ok || plaintextBase64 == "" {
		return "", fmt.Errorf("invalid plaintext returned from vault")
	}

	plaintext, err := base64.StdEncoding.DecodeString(plaintextBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 plaintext: %w", err)
	}

	return string(plaintext), nil
}

// EncryptJSON encrypts a JSON-serializable structure for a specific region
func (v *VaultServiceImpl) EncryptJSON(ctx context.Context, data interface{}, region domain.VaultRegion, credentialType domain.CredentialType) (*domain.EncryptionMetadata, error) {
	jsonData, err := sonic.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return v.EncryptData(ctx, string(jsonData), region, credentialType)
}

// DecryptJSON decrypts data into the provided target structure
func (v *VaultServiceImpl) DecryptJSON(ctx context.Context, metadata *domain.EncryptionMetadata, target interface{}) error {
	plaintext, err := v.DecryptData(ctx, metadata)
	if err != nil {
		return err
	}

	if err := sonic.Unmarshal([]byte(plaintext), target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// RotateEncryptionKey triggers key rotation in the specified region for the given credential type
func (v *VaultServiceImpl) RotateEncryptionKey(ctx context.Context, region domain.VaultRegion, credentialType domain.CredentialType) error {
	client, err := v.getRegionalClient(region)
	if err != nil {
		return err
	}

	keyName := v.getKeyName(credentialType, region)
	transitPath := v.getTransitPath(client, credentialType)
	path := fmt.Sprintf("%s/keys/%s/rotate", transitPath, keyName)

	_, err = client.client.Logical().WriteWithContext(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("failed to rotate key %s in region %s: %w", keyName, region, err)
	}

	return nil
}

// ReWrapData re-encrypts the data using the latest key version in the specified region
func (v *VaultServiceImpl) ReWrapData(ctx context.Context, currentMetadata *domain.EncryptionMetadata) (*domain.EncryptionMetadata, error) {
	if currentMetadata == nil {
		return nil, fmt.Errorf("metadata cannot be nil")
	}

	client, err := v.getRegionalClient(currentMetadata.Region)
	if err != nil {
		return nil, err
	}

	keyName := v.getKeyName(currentMetadata.CredentialType, currentMetadata.Region)
	transitPath := v.getTransitPath(client, currentMetadata.CredentialType)
	path := fmt.Sprintf("%s/rewrap/%s", transitPath, keyName)

	data := map[string]interface{}{
		"ciphertext": currentMetadata.Ciphertext,
	}

	secret, err := client.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return nil, fmt.Errorf("vault rewrap failed: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data returned from vault")
	}

	newCiphertext, ok := secret.Data["ciphertext"].(string)
	if !ok || newCiphertext == "" {
		return nil, fmt.Errorf("invalid ciphertext returned from vault")
	}

	// Extract key version from the new ciphertext
	var keyVersion int
	_, err = fmt.Sscanf(newCiphertext, "vault:v%d:", &keyVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key version from ciphertext: %w", err)
	}

	// Create new metadata with updated ciphertext and key version
	newMetadata := &domain.EncryptionMetadata{
		Region:         currentMetadata.Region,
		KeyVersion:     keyVersion,
		CredentialType: currentMetadata.CredentialType,
		Algorithm:      currentMetadata.Algorithm,
		Ciphertext:     newCiphertext,
	}

	return newMetadata, nil
}

// CheckHealth checks the health of the Vault service for the specified region
func (v *VaultServiceImpl) CheckHealth(ctx context.Context, region domain.VaultRegion) error {
	client, err := v.getRegionalClient(region)
	if err != nil {
		return err
	}

	health, err := client.client.Sys().HealthWithContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to check Vault health for region %s: %w", region, err)
	}

	if !health.Initialized {
		return fmt.Errorf("vault in region %s is not initialized", region)
	}

	if health.Sealed {
		return fmt.Errorf("vault in region %s is sealed", region)
	}

	return nil
}
