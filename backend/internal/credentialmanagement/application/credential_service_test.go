package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/infrastructure/persistence"
	credentialmanagement_mocks "github.com/context-space/context-space/backend/internal/shared/testing/mocks/credentialmanagement"
	shared_mocks "github.com/context-space/context-space/backend/internal/shared/testing/mocks/shared"
)

// CredentialServiceTestSuite is the main test suite for CredentialService
type CredentialServiceTestSuite struct {
	suite.Suite

	// Service under test
	service *CredentialService

	// Mocks
	mockCredentialRepo      *credentialmanagement_mocks.MockCredentialRepository
	mockOAuthRepo           *credentialmanagement_mocks.MockOAuthCredentialRepository
	mockAPIKeyRepo          *credentialmanagement_mocks.MockAPIKeyCredentialRepository
	mockOAuthProvider       *credentialmanagement_mocks.MockOAuthProvider
	mockEventBus            *shared_mocks.MockEventBus
	mockObs                 *observability.ObservabilityProvider
	mockUnitOfWork          *shared_mocks.MockUnitOfWork
	mockUnitOfWorkFactory   *shared_mocks.MockUnitOfWorkFactory
	mockRedisClient         *shared_mocks.MockCache
	mockTokenRefreshService *credentialmanagement_mocks.MockTokenRefresh
	mockDB                  *shared_mocks.MockDatabase
	mockVaultService        *credentialmanagement_mocks.MockVaultService

	// Test data
	testUserID             string
	testProviderIdentifier string
	testOAuth2Token        *oauth2.Token
	testAPIKey             string
	testScopes             []string
	testContext            context.Context
}

// SetupSuite initializes the test suite with standard test data
func (suite *CredentialServiceTestSuite) SetupSuite() {
	// Create observability provider mock
	suite.mockObs, _, _ = observability.InitializeObservabilityProvider(context.Background(), &observability.LogConfig{
		Level:       observability.ParseLogLevel("debug"),
		Format:      observability.ParseLogFormat("json"),
		OutputPaths: []string{"stdout"},
		Development: true,
	}, &observability.TracingConfig{
		ServiceName:    "credential-management",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		Endpoint:       "http://localhost:8080/tracing",
		Enabled:        true,
		SamplingRate:   1.0,
	}, &observability.MetricsConfig{
		ServiceName:    "credential-management",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		Enabled:        true,
	})

	// Initialize test data
	suite.testUserID = "test-user-123"
	suite.testProviderIdentifier = "test-provider"
	suite.testOAuth2Token = &oauth2.Token{
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-123",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}
	suite.testAPIKey = "api-key-123"
	suite.testScopes = []string{"read", "write"}
	suite.testContext = context.Background()
}

// SetupTest creates fresh mock instances and service for each test
func (suite *CredentialServiceTestSuite) SetupTest() {
	// Create fresh mock instances for each test to avoid state sharing
	suite.mockCredentialRepo = &credentialmanagement_mocks.MockCredentialRepository{}
	suite.mockOAuthRepo = &credentialmanagement_mocks.MockOAuthCredentialRepository{}
	suite.mockAPIKeyRepo = &credentialmanagement_mocks.MockAPIKeyCredentialRepository{}
	suite.mockOAuthProvider = &credentialmanagement_mocks.MockOAuthProvider{}
	suite.mockEventBus = &shared_mocks.MockEventBus{}
	suite.mockUnitOfWork = &shared_mocks.MockUnitOfWork{}
	suite.mockUnitOfWorkFactory = &shared_mocks.MockUnitOfWorkFactory{}
	suite.mockRedisClient = &shared_mocks.MockCache{}
	suite.mockVaultService = &credentialmanagement_mocks.MockVaultService{}
	suite.mockDB = &shared_mocks.MockDatabase{}

	// Initialize credential factory
	credentialFactory := persistence.NewCredentialFactory(
		suite.mockCredentialRepo,
		suite.mockOAuthRepo,
		suite.mockAPIKeyRepo,
		suite.mockVaultService,
	)

	// Create service instance
	suite.service = NewCredentialService(
		suite.mockCredentialRepo,
		credentialFactory,
		suite.mockEventBus,
		suite.mockObs,
		suite.mockOAuthProvider,
		suite.mockUnitOfWorkFactory,
		"http://localhost:8080/callback",
		suite.mockRedisClient,
		suite.mockTokenRefreshService,
	)
}

// BeforeTest is called before each test method to reset mock state
// This ensures clean state between test methods and prevents mock pollution
func (suite *CredentialServiceTestSuite) BeforeTest(suiteName, testName string) {
	suite.resetMockState()
}

// resetMockState clears all mock expectations and call history
// This is used both by BeforeTest (for test method isolation) and within test cases (for test case isolation)
func (suite *CredentialServiceTestSuite) resetMockState() {
	if suite.mockCredentialRepo != nil {
		suite.mockCredentialRepo.ExpectedCalls = nil
		suite.mockCredentialRepo.Calls = nil
	}
	if suite.mockOAuthRepo != nil {
		suite.mockOAuthRepo.ExpectedCalls = nil
		suite.mockOAuthRepo.Calls = nil
	}
	if suite.mockAPIKeyRepo != nil {
		suite.mockAPIKeyRepo.ExpectedCalls = nil
		suite.mockAPIKeyRepo.Calls = nil
	}
	if suite.mockOAuthProvider != nil {
		suite.mockOAuthProvider.ExpectedCalls = nil
		suite.mockOAuthProvider.Calls = nil
	}
	if suite.mockEventBus != nil {
		suite.mockEventBus.ExpectedCalls = nil
		suite.mockEventBus.Calls = nil
	}
	if suite.mockUnitOfWork != nil {
		suite.mockUnitOfWork.ExpectedCalls = nil
		suite.mockUnitOfWork.Calls = nil
	}
	if suite.mockUnitOfWorkFactory != nil {
		suite.mockUnitOfWorkFactory.ExpectedCalls = nil
		suite.mockUnitOfWorkFactory.Calls = nil
	}
	if suite.mockRedisClient != nil {
		suite.mockRedisClient.ExpectedCalls = nil
		suite.mockRedisClient.Calls = nil
	}
	if suite.mockTokenRefreshService != nil {
		suite.mockTokenRefreshService.ExpectedCalls = nil
		suite.mockTokenRefreshService.Calls = nil
	}
	if suite.mockVaultService != nil {
		suite.mockVaultService.ExpectedCalls = nil
		suite.mockVaultService.Calls = nil
	}
}

// TearDownTest verifies all mock expectations were met and cleans up
func (suite *CredentialServiceTestSuite) TearDownTest() {
	// Assert expectations for all mocks
	if suite.mockCredentialRepo != nil {
		suite.mockCredentialRepo.AssertExpectations(suite.T())
	}
	if suite.mockOAuthRepo != nil {
		suite.mockOAuthRepo.AssertExpectations(suite.T())
	}
	if suite.mockAPIKeyRepo != nil {
		suite.mockAPIKeyRepo.AssertExpectations(suite.T())
	}
	if suite.mockOAuthProvider != nil {
		suite.mockOAuthProvider.AssertExpectations(suite.T())
	}
	if suite.mockEventBus != nil {
		suite.mockEventBus.AssertExpectations(suite.T())
	}
	if suite.mockUnitOfWork != nil {
		suite.mockUnitOfWork.AssertExpectations(suite.T())
	}
	if suite.mockUnitOfWorkFactory != nil {
		suite.mockUnitOfWorkFactory.AssertExpectations(suite.T())
	}
	if suite.mockRedisClient != nil {
		suite.mockRedisClient.AssertExpectations(suite.T())
	}
	if suite.mockTokenRefreshService != nil {
		suite.mockTokenRefreshService.AssertExpectations(suite.T())
	}
	if suite.mockVaultService != nil {
		suite.mockVaultService.AssertExpectations(suite.T())
	}
	if suite.mockDB != nil {
		suite.mockDB.AssertExpectations(suite.T())
	}
}

// TestCreateOAuthCredential tests OAuth credential creation functionality
func (suite *CredentialServiceTestSuite) TestCreateOAuthCredential() {
	testCases := []struct {
		name           string
		setupMocks     func(*CredentialServiceTestSuite)
		userID         string
		providerID     string
		token          *oauth2.Token
		scopes         []string
		expectedResult *domain.OAuthCredential
		expectedError  error
		assertions     func(*CredentialServiceTestSuite, *domain.OAuthCredential, error)
	}{
		{
			name: "successful_oauth_credential_creation",
			setupMocks: func(s *CredentialServiceTestSuite) {
				// Mock unit of work
				s.mockUnitOfWorkFactory.On("Create").Return(s.mockUnitOfWork)
				s.mockUnitOfWork.On("Begin", mock.Anything).Return(nil)
				s.mockUnitOfWork.On("Commit", mock.Anything).Return(nil)

				// Mock repository - no existing credential
				s.mockCredentialRepo.On("GetByUserAndProvider", mock.Anything, s.testUserID, s.testProviderIdentifier).Return(nil, nil)

				// Mock EncryptJSON for OAuth Token (structure encryption)
				s.mockVaultService.On("EncryptJSON", mock.Anything, mock.AnythingOfType("*oauth2.Token"), domain.RegionEU, domain.CredentialTypeOAuth).Return(&domain.EncryptionMetadata{
					Region:         domain.RegionEU,
					CredentialType: domain.CredentialTypeOAuth,
					Algorithm:      domain.AlgorithmAESGCM,
					Ciphertext:     "encrypted-oauth-token-ciphertext",
				}, nil)

				s.mockCredentialRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				s.mockOAuthRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

				// Mock event bus
				s.mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("events.Event")).Return(nil)
			},
			userID:     "test-user-123",
			providerID: "test-provider",
			token: &oauth2.Token{
				AccessToken:  "access-token-123",
				RefreshToken: "refresh-token-123",
				TokenType:    "Bearer",
				Expiry:       time.Now().Add(time.Hour),
			},
			scopes:        []string{"read", "write"},
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, result *domain.OAuthCredential, err error) {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), result)
				assert.Equal(s.T(), s.testUserID, result.UserID)
				assert.Equal(s.T(), s.testProviderIdentifier, result.ProviderIdentifier)
				assert.Equal(s.T(), domain.CredentialTypeOAuth, result.Type)
			},
		},
		{
			name: "replace_existing_oauth_credential",
			setupMocks: func(s *CredentialServiceTestSuite) {
				// Mock unit of work
				s.mockUnitOfWorkFactory.On("Create").Return(s.mockUnitOfWork)
				s.mockUnitOfWork.On("Begin", mock.Anything).Return(nil)
				s.mockUnitOfWork.On("Commit", mock.Anything).Return(nil)

				// Mock repository - existing credential
				existingCred := &domain.Credential{
					ID:                 "existing-cred-id",
					UserID:             s.testUserID,
					ProviderIdentifier: s.testProviderIdentifier,
					Type:               domain.CredentialTypeOAuth,
				}
				s.mockCredentialRepo.On("GetByUserAndProvider", mock.Anything, s.testUserID, s.testProviderIdentifier).Return(existingCred, nil)
				s.mockCredentialRepo.On("Delete", mock.Anything, "existing-cred-id").Return(nil)

				// Mock EncryptJSON for OAuth Token (structure encryption)
				s.mockVaultService.On("EncryptJSON", mock.Anything, mock.AnythingOfType("*oauth2.Token"), domain.RegionEU, domain.CredentialTypeOAuth).Return(&domain.EncryptionMetadata{
					Region:         domain.RegionEU,
					CredentialType: domain.CredentialTypeOAuth,
					Algorithm:      domain.AlgorithmAESGCM,
					Ciphertext:     "encrypted-oauth-token-ciphertext-2",
				}, nil)

				s.mockCredentialRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				s.mockOAuthRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				// Mock event bus
				s.mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("events.Event")).Return(nil)
			},
			userID:     "test-user-123",
			providerID: "test-provider",
			token: &oauth2.Token{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
				TokenType:    "Bearer",
			},
			scopes:        []string{"read", "write", "admin"},
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, result *domain.OAuthCredential, err error) {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), result)
				assert.NotEmpty(s.T(), result.ID)
				assert.Equal(s.T(), s.testUserID, result.UserID)
				assert.Equal(s.T(), s.testProviderIdentifier, result.ProviderIdentifier)
				assert.Equal(s.T(), domain.CredentialTypeOAuth, result.Type)
			},
		},
		{
			name: "transaction_failure_rollback",
			setupMocks: func(s *CredentialServiceTestSuite) {
				// Mock unit of work with failure
				s.mockUnitOfWorkFactory.On("Create").Return(s.mockUnitOfWork)
				s.mockUnitOfWork.On("Begin", mock.Anything).Return(errors.New("transaction begin failed"))
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			token:         &oauth2.Token{AccessToken: "token"},
			scopes:        []string{"read"},
			expectedError: errors.New("transaction begin failed"),
			assertions: func(s *CredentialServiceTestSuite, result *domain.OAuthCredential, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Contains(s.T(), err.Error(), "transaction begin failed")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			// Setup mocks for this test case
			tc.setupMocks(suite)

			// Execute the method under test
			result, err := suite.service.CreateOAuthCredential(
				suite.testContext,
				tc.userID,
				tc.providerID,
				tc.token,
				tc.scopes,
			)

			// Run assertions
			tc.assertions(suite, result, err)
		})
	}
}

// TestCreateAPIKeyCredential tests API key credential creation functionality
func (suite *CredentialServiceTestSuite) TestCreateAPIKeyCredential() {
	testCases := []struct {
		name           string
		setupMocks     func(*CredentialServiceTestSuite)
		userID         string
		providerID     string
		apiKey         string
		expectedResult *domain.APIKeyCredential
		expectedError  error
		assertions     func(*CredentialServiceTestSuite, *domain.APIKeyCredential, error)
	}{
		{
			name: "successful_apikey_credential_creation",
			setupMocks: func(s *CredentialServiceTestSuite) {
				// Mock unit of work
				s.mockUnitOfWorkFactory.On("Create").Return(s.mockUnitOfWork)
				s.mockUnitOfWork.On("Begin", mock.Anything).Return(nil)
				s.mockUnitOfWork.On("Commit", mock.Anything).Return(nil)

				// Mock repository - no existing credential
				s.mockCredentialRepo.On("GetByUserAndProvider", mock.Anything, s.testUserID, s.testProviderIdentifier).Return(nil, nil)

				// Mock EncryptData for API Key (string encryption)
				s.mockVaultService.On("EncryptData", mock.Anything, s.testAPIKey, domain.RegionEU, domain.CredentialTypeAPIKey).Return(&domain.EncryptionMetadata{
					Region:         domain.RegionEU,
					CredentialType: domain.CredentialTypeAPIKey,
					Algorithm:      domain.AlgorithmAESGCM,
					Ciphertext:     "encrypted-api-key-ciphertext",
				}, nil)

				s.mockCredentialRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				s.mockAPIKeyRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

				// Mock event bus
				s.mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("events.Event")).Return(nil)
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			apiKey:        "api-key-123",
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, result *domain.APIKeyCredential, err error) {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), result)
				assert.Equal(s.T(), s.testUserID, result.UserID)
				assert.Equal(s.T(), s.testProviderIdentifier, result.ProviderIdentifier)
				assert.Equal(s.T(), domain.CredentialTypeAPIKey, result.Type)
				assert.Equal(s.T(), s.testAPIKey, result.APIKey)
			},
		},
		{
			name: "replace_existing_apikey_credential",
			setupMocks: func(s *CredentialServiceTestSuite) {
				// Mock unit of work
				s.mockUnitOfWorkFactory.On("Create").Return(s.mockUnitOfWork)
				s.mockUnitOfWork.On("Begin", mock.Anything).Return(nil)
				s.mockUnitOfWork.On("Commit", mock.Anything).Return(nil)

				// Mock repository - existing credential
				existingCred := &domain.Credential{
					ID:                 "existing-apikey-id",
					UserID:             s.testUserID,
					ProviderIdentifier: s.testProviderIdentifier,
					Type:               domain.CredentialTypeAPIKey,
				}
				s.mockCredentialRepo.On("GetByUserAndProvider", mock.Anything, s.testUserID, s.testProviderIdentifier).Return(existingCred, nil)
				s.mockCredentialRepo.On("Delete", mock.Anything, "existing-apikey-id").Return(nil)

				// Mock EncryptData for new API Key (string encryption)
				s.mockVaultService.On("EncryptData", mock.Anything, "new-api-key", domain.RegionEU, domain.CredentialTypeAPIKey).Return(&domain.EncryptionMetadata{
					Region:         domain.RegionEU,
					CredentialType: domain.CredentialTypeAPIKey,
					Algorithm:      domain.AlgorithmAESGCM,
					Ciphertext:     "encrypted-new-api-key-ciphertext",
				}, nil)

				s.mockCredentialRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				s.mockAPIKeyRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				// Mock event bus
				s.mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("events.Event")).Return(nil)
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			apiKey:        "new-api-key",
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, result *domain.APIKeyCredential, err error) {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), result)
				assert.NotEmpty(s.T(), result.ID)
				assert.Equal(s.T(), "new-api-key", result.APIKey)
				assert.Equal(s.T(), s.testUserID, result.UserID)
				assert.Equal(s.T(), s.testProviderIdentifier, result.ProviderIdentifier)
			},
		},
		{
			name: "empty_api_key_validation_error",
			setupMocks: func(s *CredentialServiceTestSuite) {
				// Mock unit of work
				s.mockUnitOfWorkFactory.On("Create").Return(s.mockUnitOfWork)
				s.mockUnitOfWork.On("Begin", mock.Anything).Return(nil)
				s.mockUnitOfWork.On("Rollback", mock.Anything).Return(nil)

				s.mockCredentialRepo.On("GetByUserAndProvider", mock.Anything, s.testUserID, s.testProviderIdentifier).Return(nil, nil)

			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			apiKey:        "",
			expectedError: errors.New("API key is required"),
			assertions: func(s *CredentialServiceTestSuite, result *domain.APIKeyCredential, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Contains(s.T(), err.Error(), "API key is required")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			// Setup mocks for this test case
			tc.setupMocks(suite)

			// Execute the method under test
			result, err := suite.service.CreateAPIKeyCredential(
				suite.testContext,
				tc.userID,
				tc.providerID,
				tc.apiKey,
			)

			// Run assertions
			tc.assertions(suite, result, err)
		})
	}
}

// TestGetCredential tests credential retrieval by ID functionality
func (suite *CredentialServiceTestSuite) TestGetCredential() {
	testCases := []struct {
		name           string
		setupMocks     func(*CredentialServiceTestSuite)
		credentialID   string
		expectedResult interface{}
		expectedError  error
		assertions     func(*CredentialServiceTestSuite, interface{}, error)
	}{
		{
			name: "successful_credential_retrieval",
			setupMocks: func(s *CredentialServiceTestSuite) {
				s.mockCredentialRepo.On("GetByID", mock.Anything, "test-cred-id").Return(s.createTestCredential("test-cred-id", domain.CredentialTypeOAuth), nil)
				s.mockOAuthRepo.On("GetByCredentialID", mock.Anything, "test-cred-id").Return(s.createTestOAuthCredential("test-cred-id"), nil)
				s.mockVaultService.On("DecryptJSON", mock.Anything, mock.Anything, mock.AnythingOfType("*oauth2.Token")).Return(nil)
			},
			credentialID:  "test-cred-id",
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, result interface{}, err error) {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), result)

				oauthCred, ok := result.(*domain.OAuthCredential)
				require.True(s.T(), ok)
				assert.Equal(s.T(), "test-cred-id", oauthCred.ID)
				assert.Equal(s.T(), s.testUserID, oauthCred.UserID)
			},
		},
		{
			name: "credential_not_found",
			setupMocks: func(s *CredentialServiceTestSuite) {
				s.mockCredentialRepo.On("GetByID", mock.Anything, "non-existing-id").Return(nil, nil)
			},
			credentialID:  "non-existing-id",
			expectedError: ErrCredentialNotFound,
			assertions: func(s *CredentialServiceTestSuite, result interface{}, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Equal(s.T(), ErrCredentialNotFound, err)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			tc.setupMocks(suite)
			result, err := suite.service.GetCredential(suite.testContext, tc.credentialID)
			tc.assertions(suite, result, err)
		})
	}
}

// TestGetCredentialByUserAndProvider tests credential retrieval with locking mechanism
func (suite *CredentialServiceTestSuite) TestGetCredentialByUserAndProvider() {
	testCases := []struct {
		name           string
		setupMocks     func(*CredentialServiceTestSuite)
		userID         string
		providerID     string
		expectedResult interface{}
		expectedError  error
		assertions     func(*CredentialServiceTestSuite, interface{}, error)
	}{
		{
			name: "successful_retrieval_with_token_refresh",
			setupMocks: func(s *CredentialServiceTestSuite) {
				// Mock cache lock acquisition
				lockKey := "access_token_lock:test-provider:test-user-123"
				s.mockRedisClient.On("AcquireLock", mock.Anything, lockKey, time.Second*1).Return(true, nil)
				s.mockRedisClient.On("ReleaseLock", mock.Anything, lockKey).Return(nil)

				// Mock the credential factory call chain
				originalCred := s.createTestCredential("test-cred-id", domain.CredentialTypeOAuth)
				s.mockCredentialRepo.On("GetByUserAndProvider", mock.Anything, s.testUserID, s.testProviderIdentifier).Return(originalCred, nil)
				s.mockCredentialRepo.On("GetByID", mock.Anything, "test-cred-id").Return(originalCred, nil)

				originalToken := &oauth2.Token{
					AccessToken:  "original-access-token",
					RefreshToken: "original-refresh-token",
					TokenType:    "Bearer",
					Expiry:       time.Now().Add(time.Hour),
				}
				oauthCred := &domain.OAuthCredential{
					Credential: &domain.Credential{
						ID:                 "test-cred-id",
						UserID:             s.testUserID,
						ProviderIdentifier: s.testProviderIdentifier,
						Type:               domain.CredentialTypeOAuth,
						IsValid:            true,
					},
					Token:  originalToken,
					Scopes: s.testScopes,
				}

				s.mockOAuthRepo.On("GetByCredentialID", mock.Anything, "test-cred-id").Return(oauthCred, nil)

				s.mockVaultService.On("DecryptJSON", mock.Anything, mock.Anything, mock.AnythingOfType("*oauth2.Token")).
					Run(func(args mock.Arguments) {
						token := args.Get(2).(*oauth2.Token)
						token.AccessToken = "decrypted-access-token"
						token.RefreshToken = "decrypted-refresh-token"
					}).Return(nil)

				refreshedToken := &oauth2.Token{
					AccessToken:  "refreshed-access-token",
					RefreshToken: "refreshed-refresh-token",
					TokenType:    "Bearer",
					Expiry:       time.Now().Add(2 * time.Hour),
				}
				refreshedCred := &domain.OAuthCredential{
					Credential: &domain.Credential{
						ID:                 "test-cred-id",
						UserID:             s.testUserID,
						ProviderIdentifier: s.testProviderIdentifier,
						Type:               domain.CredentialTypeOAuth,
						IsValid:            true,
					},
					Token:  refreshedToken,
					Scopes: s.testScopes,
				}

				s.mockTokenRefreshService.On("RefreshAccessTokenIfNeeded", mock.Anything, s.testProviderIdentifier, mock.Anything).Return(refreshedCred, nil)
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, result interface{}, err error) {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), result)

				oauthCred, ok := result.(*domain.OAuthCredential)
				require.True(s.T(), ok)
				assert.Equal(s.T(), "refreshed-access-token", oauthCred.Token.AccessToken)
			},
		},
		{
			name: "lock_acquisition_failure",
			setupMocks: func(s *CredentialServiceTestSuite) {
				// Mock cache lock acquisition failure - return false means the lock is held by other processes
				lockKey := "access_token_lock:test-provider:test-user-123"

				s.mockRedisClient.On("AcquireLock", mock.Anything, lockKey, time.Second*1).Return(false, nil)
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			expectedError: ErrCredentialExpired,
			assertions: func(s *CredentialServiceTestSuite, result interface{}, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Contains(s.T(), err.Error(), "please re-authorize")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			tc.setupMocks(suite)
			result, err := suite.service.GetCredentialByUserAndProvider(suite.testContext, tc.userID, tc.providerID)
			tc.assertions(suite, result, err)
		})
	}
}

// TestDeleteCredential tests credential deletion functionality
func (suite *CredentialServiceTestSuite) TestDeleteCredential() {
	testCases := []struct {
		name          string
		setupMocks    func(*CredentialServiceTestSuite)
		credentialID  string
		expectedError error
		assertions    func(*CredentialServiceTestSuite, error)
	}{
		{
			name: "successful_credential_deletion",
			setupMocks: func(s *CredentialServiceTestSuite) {
				// Mock getting credential for event data
				existingCred := &domain.Credential{
					ID:                 "test-cred-id",
					UserID:             s.testUserID,
					ProviderIdentifier: s.testProviderIdentifier,
					Type:               domain.CredentialTypeOAuth,
				}
				s.mockCredentialRepo.On("GetByID", mock.Anything, "test-cred-id").Return(existingCred, nil)

				// Mock deletion
				s.mockCredentialRepo.On("Delete", mock.Anything, "test-cred-id").Return(nil)

				// Mock event publishing
				s.mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("events.Event")).Return(nil)
			},
			credentialID:  "test-cred-id",
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, err error) {
				require.NoError(s.T(), err)
			},
		},
		{
			name: "credential_not_found_for_deletion",
			setupMocks: func(s *CredentialServiceTestSuite) {
				s.mockCredentialRepo.On("GetByID", mock.Anything, "non-existing-id").Return(nil, nil)
			},
			credentialID:  "non-existing-id",
			expectedError: ErrCredentialNotFound,
			assertions: func(s *CredentialServiceTestSuite, err error) {
				require.Error(s.T(), err)
				assert.Equal(s.T(), ErrCredentialNotFound, err)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			tc.setupMocks(suite)
			err := suite.service.DeleteCredential(suite.testContext, tc.credentialID)
			tc.assertions(suite, err)
		})
	}
}

// TestGetOAuthURL tests OAuth URL generation functionality
func (suite *CredentialServiceTestSuite) TestGetOAuthURL() {
	testCases := []struct {
		name          string
		setupMocks    func(*CredentialServiceTestSuite)
		providerID    string
		state         string
		codeChallenge string
		permissionIDs []string
		expectedURL   string
		expectedError error
		assertions    func(*CredentialServiceTestSuite, string, error)
	}{
		{
			name: "successful_oauth_url_generation",
			setupMocks: func(s *CredentialServiceTestSuite) {
				s.mockOAuthProvider.On("GetScopesFromPermissions", mock.Anything, "test-provider", []string{"read", "write"}).Return([]string{"read_scope", "write_scope"}, nil)
				s.mockOAuthProvider.On("GenerateOAuthURL", mock.Anything, "test-provider", "http://localhost:8080/callback", "test-state", "test-challenge", []string{"read_scope", "write_scope"}).Return("https://oauth.provider.com/auth?client_id=test", nil)
			},
			providerID:    "test-provider",
			state:         "test-state",
			codeChallenge: "test-challenge",
			permissionIDs: []string{"read", "write"},
			expectedURL:   "https://oauth.provider.com/auth?client_id=test",
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, url string, err error) {
				require.NoError(s.T(), err)
				assert.Equal(s.T(), "https://oauth.provider.com/auth?client_id=test", url)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			tc.setupMocks(suite)
			result, err := suite.service.GetOAuthURL(suite.testContext, tc.providerID, tc.state, tc.codeChallenge, tc.permissionIDs)
			tc.assertions(suite, result, err)
		})
	}
}

// TestHandleOAuthCallback tests OAuth callback handling functionality
func (suite *CredentialServiceTestSuite) TestHandleOAuthCallback() {
	testCases := []struct {
		name           string
		setupMocks     func(*CredentialServiceTestSuite)
		code           string
		providerID     string
		userID         string
		permissions    []string
		codeVerifier   string
		expectedResult *domain.OAuthCredential
		expectedError  error
		assertions     func(*CredentialServiceTestSuite, *domain.OAuthCredential, error)
	}{
		{
			name: "successful_oauth_callback_handling",
			setupMocks: func(s *CredentialServiceTestSuite) {
				// Mock token exchange
				s.mockOAuthProvider.On("ExchangeCodeForToken", mock.Anything, "test-provider", "auth-code", "http://localhost:8080/callback", "code-verifier").Return(s.testOAuth2Token, nil)

				// Mock scopes conversion
				s.mockOAuthProvider.On("GetScopesFromPermissions", mock.Anything, "test-provider", []string{"read", "write"}).Return(s.testScopes, nil)

				// Mock credential creation (delegated to CreateOAuthCredential)
				s.mockUnitOfWorkFactory.On("Create").Return(s.mockUnitOfWork)
				s.mockUnitOfWork.On("Begin", mock.Anything).Return(nil)
				s.mockUnitOfWork.On("Commit", mock.Anything).Return(nil)
				s.mockCredentialRepo.On("GetByUserAndProvider", mock.Anything, s.testUserID, "test-provider").Return(nil, nil)

				// Mock EncryptJSON for OAuth Token
				s.mockVaultService.On("EncryptJSON", mock.Anything, mock.AnythingOfType("*oauth2.Token"), domain.RegionEU, domain.CredentialTypeOAuth).Return(&domain.EncryptionMetadata{
					Region:         domain.RegionEU,
					CredentialType: domain.CredentialTypeOAuth,
					Algorithm:      domain.AlgorithmAESGCM,
					Ciphertext:     "encrypted-callback-token-ciphertext",
				}, nil)

				s.mockCredentialRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				s.mockOAuthRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				s.mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("events.Event")).Return(nil)
			},
			code:          "auth-code",
			providerID:    "test-provider",
			userID:        "test-user-123",
			permissions:   []string{"read", "write"},
			codeVerifier:  "code-verifier",
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, result *domain.OAuthCredential, err error) {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), result)
				assert.NotEmpty(s.T(), result.ID)
				assert.Equal(s.T(), s.testUserID, result.UserID)
				assert.Equal(s.T(), "test-provider", result.ProviderIdentifier)
				assert.Equal(s.T(), domain.CredentialTypeOAuth, result.Type)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			tc.setupMocks(suite)
			result, err := suite.service.HandleOAuthCallback(suite.testContext, tc.code, tc.providerID, tc.userID, tc.permissions, tc.codeVerifier)
			tc.assertions(suite, result, err)
		})
	}
}

// TestGetAllCredentialsByUser tests batch credential retrieval functionality
func (suite *CredentialServiceTestSuite) TestGetAllCredentialsByUser() {
	testCases := []struct {
		name           string
		setupMocks     func(*CredentialServiceTestSuite)
		userID         string
		expectedResult []*domain.Credential
		expectedError  error
		assertions     func(*CredentialServiceTestSuite, []*domain.Credential, error)
	}{
		{
			name: "successful_batch_retrieval",
			setupMocks: func(s *CredentialServiceTestSuite) {
				expectedCreds := []*domain.Credential{
					{
						ID:                 "cred-1",
						UserID:             s.testUserID,
						ProviderIdentifier: "provider-1",
						Type:               domain.CredentialTypeOAuth,
						IsValid:            true,
					},
					{
						ID:                 "cred-2",
						UserID:             s.testUserID,
						ProviderIdentifier: "provider-2",
						Type:               domain.CredentialTypeAPIKey,
						IsValid:            true,
					},
				}
				s.mockCredentialRepo.On("ListByUser", mock.Anything, s.testUserID).Return(expectedCreds, nil)
			},
			userID:        "test-user-123",
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, result []*domain.Credential, err error) {
				require.NoError(s.T(), err)
				require.Len(s.T(), result, 2)
				assert.Equal(s.T(), "cred-1", result[0].ID)
				assert.Equal(s.T(), "cred-2", result[1].ID)
			},
		},
		{
			name: "empty_result_for_user",
			setupMocks: func(s *CredentialServiceTestSuite) {
				s.mockCredentialRepo.On("ListByUser", mock.Anything, "empty-user").Return([]*domain.Credential{}, nil)
			},
			userID:        "empty-user",
			expectedError: nil,
			assertions: func(s *CredentialServiceTestSuite, result []*domain.Credential, err error) {
				require.NoError(s.T(), err)
				require.Len(s.T(), result, 0)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			tc.setupMocks(suite)
			result, err := suite.service.GetAllCredentialsByUser(suite.testContext, tc.userID)
			tc.assertions(suite, result, err)
		})
	}
}

// TestCredentialServiceTestSuite runs the credential service test suite
func TestCredentialServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CredentialServiceTestSuite))
}

// Helper methods for creating test data
func (suite *CredentialServiceTestSuite) createTestCredential(id string, credType domain.CredentialType) *domain.Credential {
	return &domain.Credential{
		ID:                 id,
		UserID:             suite.testUserID,
		ProviderIdentifier: suite.testProviderIdentifier,
		Type:               credType,
		IsValid:            true,
	}
}

// createTestOAuthCredential creates a test OAuth credential
func (suite *CredentialServiceTestSuite) createTestOAuthCredential(id string) *domain.OAuthCredential {
	return &domain.OAuthCredential{
		Credential: &domain.Credential{
			ID:                 id,
			UserID:             suite.testUserID,
			ProviderIdentifier: suite.testProviderIdentifier,
			Type:               domain.CredentialTypeOAuth,
			IsValid:            true,
		},
		Token:  suite.testOAuth2Token,
		Scopes: suite.testScopes,
	}
}
