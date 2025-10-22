package provider

import (
	"context"

	"golang.org/x/text/language"
)

// ProviderCoreReader defines the contract for Provider query service
type ProviderCoreReader interface {
	// GetProviderByIdentifier retrieves Provider information by identifier
	GetProviderByIdentifier(ctx context.Context, identifier string) (*ProviderDTO, error)

	// GetProviderWithTranslation retrieves Provider information by identifier with translation
	GetProviderWithTranslation(ctx context.Context, identifier string, preferredLang language.Tag) (*ProviderDTO, error)

	// ListProviders retrieves all providers
	ListProviders(ctx context.Context) ([]*ProviderDTO, error)

	// ListProvidersByIDs retrieves Provider information by ID list
	ListProvidersByIDs(ctx context.Context, ids []string) ([]*ProviderDTO, error)

	// ListOperationsByIDs retrieves Operation information by ID list
	ListOperationsByIDs(ctx context.Context, ids []string) ([]*OperationDTO, error)
}
