package provideradapter

import "context"

type ProviderAdapterReloader interface {
	// ReloadAllProviders reloads all providers
	ReloadAllProviders(ctx context.Context) error

	// ReloadProvider reloads a specific provider
	ReloadProvider(ctx context.Context, identifier string) error
}
