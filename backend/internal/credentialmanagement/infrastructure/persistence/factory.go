package persistence

import (
	"context"
	"errors"

	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	contractCredential "github.com/context-space/context-space/backend/internal/shared/contract/credentialmanagement"
	"golang.org/x/oauth2"
)

// Errors returned by the credential factory
var (
	ErrCredentialTypeNotSupported = errors.New("credential type not supported")
	ErrCredentialNotFound         = errors.New("credential not found")
)

// CredentialFactoryImpl implements the CredentialFactory interface
type CredentialFactoryImpl struct {
	credentialRepo domain.CredentialRepository
	oauthRepo      domain.OAuthCredentialRepository
	apiKeyRepo     domain.APIKeyCredentialRepository
	vaultService   domain.VaultService
}

// NewCredentialFactory creates a new instance of the credential factory
func NewCredentialFactory(
	credentialRepo domain.CredentialRepository,
	oauthRepo domain.OAuthCredentialRepository,
	apiKeyRepo domain.APIKeyCredentialRepository,
	vaultService domain.VaultService,
) domain.CredentialFactory {
	return &CredentialFactoryImpl{
		credentialRepo: credentialRepo,
		oauthRepo:      oauthRepo,
		apiKeyRepo:     apiKeyRepo,
		vaultService:   vaultService,
	}
}

// CreateOAuth creates a new OAuth credential
func (f *CredentialFactoryImpl) CreateOAuth(
	ctx context.Context, userID, providerIdentifier string, oauth2Token *oauth2.Token, scopes []string,
) (*domain.OAuthCredential, error) {
	// Create the OAuth credential domain object
	oauthCred, err := domain.NewOAuthCredential(userID, providerIdentifier, oauth2Token, scopes)
	if err != nil {
		return nil, err
	}

	// Encrypt the OAuth token
	metadata, err := f.vaultService.EncryptJSON(ctx, oauth2Token, domain.RegionEU, domain.CredentialTypeOAuth)
	if err != nil {
		return nil, err
	}

	// Set encryption metadata
	oauthCred.EncryptionMetadata = metadata

	// Create the base credential in the database
	if err := f.credentialRepo.Create(ctx, oauthCred.Credential); err != nil {
		return nil, err
	}

	// Create the OAuth credential in the database
	if err := f.oauthRepo.Create(ctx, oauthCred); err != nil {
		// Try to clean up the base credential if oauth creation fails
		_ = f.credentialRepo.Delete(ctx, oauthCred.ID)
		return nil, err
	}

	return oauthCred, nil
}

// CreateAPIKey creates a new API key credential
func (f *CredentialFactoryImpl) CreateAPIKey(
	ctx context.Context,
	userID, providerIdentifier, apiKey string,
) (*domain.APIKeyCredential, error) {
	// Create the API key credential domain object
	apiKeyCred, err := domain.NewAPIKeyCredential(userID, providerIdentifier, apiKey)
	if err != nil {
		return nil, err
	}

	// Encrypt the OAuth token
	metadata, err := f.vaultService.EncryptData(ctx, apiKeyCred.APIKey, domain.RegionEU, domain.CredentialTypeAPIKey)
	if err != nil {
		return nil, err
	}

	// Set encryption metadata
	apiKeyCred.EncryptionMetadata = metadata

	// Create the base credential in the database
	if err := f.credentialRepo.Create(ctx, apiKeyCred.Credential); err != nil {
		return nil, err
	}

	// Create the API key credential in the database
	if err := f.apiKeyRepo.Create(ctx, apiKeyCred); err != nil {
		// Try to clean up the base credential if API key creation fails
		_ = f.credentialRepo.Delete(ctx, apiKeyCred.ID)
		return nil, err
	}

	return apiKeyCred, nil
}

// CreateNone creates a new no-auth credential
func (f *CredentialFactoryImpl) CreateNone(
	ctx context.Context,
	userID, providerIdentifier string,
) (*contractCredential.CredentialDTO, error) {
	// Create the none credential domain object
	noneCred, err := domain.NewNoneCredential(userID, providerIdentifier)
	if err != nil {
		return nil, err
	}

	return &contractCredential.CredentialDTO{
		ID:                 noneCred.ID,
		UserID:             noneCred.UserID,
		ProviderIdentifier: noneCred.ProviderIdentifier,
		Type:               string(noneCred.Type),
		IsValid:            noneCred.IsValid,
	}, nil
}

// GetCredential retrieves a credential by ID and returns the specific type
func (f *CredentialFactoryImpl) GetCredential(ctx context.Context, id string) (interface{}, error) {
	// Retrieve the base credential to determine its type
	baseCred, err := f.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// If the base credential is not found, return nil
	if baseCred == nil {
		return nil, nil
	}

	return f.loadCredentialDetails(ctx, baseCred)
}

// GetCredentialByUserAndProvider retrieves a credential by user ID and provider identifier
func (f *CredentialFactoryImpl) GetCredentialByUserAndProvider(ctx context.Context, userID, providerIdentifier string) (interface{}, error) {
	// Retrieve the base credential to determine its type
	baseCred, err := f.credentialRepo.GetByUserAndProvider(ctx, userID, providerIdentifier)
	if err != nil {
		return nil, err
	}

	// If the base credential is not found, return nil
	if baseCred == nil {
		return nil, nil
	}

	return f.loadCredentialDetails(ctx, baseCred)
}

// loadCredentialDetails loads and decrypts the complete credential information based on the base credential
func (f *CredentialFactoryImpl) loadCredentialDetails(ctx context.Context, baseCred *domain.Credential) (interface{}, error) {
	switch baseCred.Type {
	case domain.CredentialTypeOAuth:
		oauthCred, err := f.oauthRepo.GetByCredentialID(ctx, baseCred.ID)
		if err != nil {
			return nil, err
		}
		token := &oauth2.Token{}
		if err := f.vaultService.DecryptJSON(ctx, oauthCred.EncryptionMetadata, token); err != nil {
			return nil, err
		}
		oauthCred.Token = token
		return oauthCred, nil

	case domain.CredentialTypeAPIKey:
		apiKeyCred, err := f.apiKeyRepo.GetByCredentialID(ctx, baseCred.ID)
		if err != nil {
			return nil, err
		}
		apiKey, err := f.vaultService.DecryptData(ctx, apiKeyCred.EncryptionMetadata)
		if err != nil {
			return nil, err
		}
		apiKeyCred.APIKey = apiKey
		return apiKeyCred, nil

	default:
		return nil, ErrCredentialTypeNotSupported
	}
}

func (f *CredentialFactoryImpl) UpdateCredentialLastUsedAt(ctx context.Context, credential interface{}) error {
	switch cred := credential.(type) {
	case *domain.OAuthCredential:
		return f.credentialRepo.UpdateLastUsedAt(ctx, cred.Credential.ID)
	case *domain.APIKeyCredential:
		return f.credentialRepo.UpdateLastUsedAt(ctx, cred.Credential.ID)
	}

	return nil
}
