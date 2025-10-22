package credentialmanagement

import "context"

// CredentialManagementContract defines the contract interface for credential management
// This is used for cross-module communication through the contract layer
//
// Design Principles:
//  1. Stable Interface: This contract provides a stable interface that won't change
//     due to internal implementation changes in credentialmanagement module
//  2. Bounded Context Isolation: Prevents direct dependencies between bounded contexts
//  3. Anti-Corruption Layer Support: Enables ACL implementation without coupling
//  4. Version Management: Can be versioned for backward compatibility
//
// Contract Version: v1.0
type CredentialManagementContract interface {
	// GetCredentialByUserAndProviderContract retrieves a credential by user ID and provider ID
	// Returns the raw credential object (domain entity)
	GetCredentialByUserAndProviderContract(ctx context.Context, userID, providerIdentifier string) (interface{}, error)

	// CreateNoneCredentialContract creates a new no-auth credential
	// Returns a standardized DTO for cross-module communication
	CreateNoneCredentialContract(ctx context.Context, userID, providerIdentifier string) (*CredentialDTO, error)

	// UpdateCredentialLastUsedAtContract updates the last used at time of a credential
	// Accepts any credential type (interface{}) for flexibility
	UpdateCredentialLastUsedAtContract(ctx context.Context, credential interface{}) error

	// RefreshAccessTokenContract refreshes OAuth access token if needed
	// Returns updated credential or original if refresh not needed
	RefreshAccessTokenContract(ctx context.Context, providerIdentifier string, credential interface{}) (interface{}, error)
}
