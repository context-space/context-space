package types

type ProviderAuthType string

const (
	AuthTypeNone   ProviderAuthType = "none"
	AuthTypeAPIKey ProviderAuthType = "apikey"
	AuthTypeOAuth  ProviderAuthType = "oauth"
	AuthTypeBasic  ProviderAuthType = "basic"
)

type ProviderStatus string

const (
	ProviderStatusActive      ProviderStatus = "active"
	ProviderStatusInactive    ProviderStatus = "inactive"
	ProviderStatusMaintenance ProviderStatus = "maintenance"
	ProviderStatusDeprecated  ProviderStatus = "deprecated"
)

type Permission struct {
	Identifier  string   `json:"identifier"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	OAuthScopes []string `json:"oauth_scopes"`
}
