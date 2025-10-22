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
	credentialDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
	"github.com/context-space/context-space/backend/internal/integration/domain"
	providerAdapterDomain "github.com/context-space/context-space/backend/internal/provideradapter/domain"
	providerDomain "github.com/context-space/context-space/backend/internal/providercore/domain"
	integration_mocks "github.com/context-space/context-space/backend/internal/shared/testing/mocks/integration"
	provideradapter_mocks "github.com/context-space/context-space/backend/internal/shared/testing/mocks/provideradapter"
	shared_mocks "github.com/context-space/context-space/backend/internal/shared/testing/mocks/shared"
	"github.com/context-space/context-space/backend/internal/shared/types"
)

// InvocationServiceTestSuite is the main test suite for InvocationService
type InvocationServiceTestSuite struct {
	suite.Suite

	// Service under test
	service *InvocationService

	// Mocks
	mockProviderProvider    *integration_mocks.MockProviderProvider
	mockAdapterProvider     *integration_mocks.MockAdapterProvider
	mockCredentialProvider  *integration_mocks.MockCredentialProvider
	mockInvocationRepo      *integration_mocks.MockInvocationRepository
	mockEventBus            *shared_mocks.MockEventBus
	mockRedisClient         *shared_mocks.MockCache
	mockTokenRefreshService *integration_mocks.MockTokenRefreshProvider
	mockProviderAdapter     *provideradapter_mocks.MockAdapter
	mockObs                 *observability.ObservabilityProvider

	// Test data
	testUserID             string
	testProviderIdentifier string
	testOperationID        string
	testParams             map[string]interface{}
	testContext            context.Context
}

// SetupSuite initializes the test suite with standard test data
func (suite *InvocationServiceTestSuite) SetupSuite() {
	// Initialize test data
	suite.testUserID = "test-user-123"
	suite.testProviderIdentifier = "test-provider"
	suite.testOperationID = "test-operation"
	suite.testParams = map[string]interface{}{
		"param1": "value1",
		"param2": 123,
	}
	suite.testContext = context.Background()

	// Create observability provider mock
	suite.mockObs, _, _ = observability.InitializeObservabilityProvider(context.Background(), &observability.LogConfig{
		Level:       observability.ParseLogLevel("debug"),
		Format:      observability.ParseLogFormat("json"),
		OutputPaths: []string{"stdout"},
		Development: true,
	}, &observability.TracingConfig{
		ServiceName:    "integration",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		Endpoint:       "http://localhost:8080/tracing",
		Enabled:        true,
		SamplingRate:   1.0,
	}, &observability.MetricsConfig{
		ServiceName:    "integration",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		Enabled:        true,
	})
}

// SetupTest creates fresh mock instances and service for each test
func (suite *InvocationServiceTestSuite) SetupTest() {
	// Create fresh mock instances for each test to avoid state sharing
	suite.mockProviderProvider = &integration_mocks.MockProviderProvider{}
	suite.mockAdapterProvider = &integration_mocks.MockAdapterProvider{}
	suite.mockCredentialProvider = &integration_mocks.MockCredentialProvider{}
	suite.mockInvocationRepo = &integration_mocks.MockInvocationRepository{}
	suite.mockEventBus = &shared_mocks.MockEventBus{}
	suite.mockRedisClient = &shared_mocks.MockCache{}
	suite.mockTokenRefreshService = &integration_mocks.MockTokenRefreshProvider{}
	suite.mockProviderAdapter = &provideradapter_mocks.MockAdapter{}

	// Create service instance
	suite.service = NewInvocationService(
		suite.mockProviderProvider,
		suite.mockAdapterProvider,
		suite.mockCredentialProvider,
		suite.mockInvocationRepo,
		suite.mockEventBus,
		suite.mockObs,
		suite.mockRedisClient,
		suite.mockTokenRefreshService,
	)
}

// BeforeTest is called before each test method to reset mock state
// This ensures clean state between test methods and prevents mock pollution
func (suite *InvocationServiceTestSuite) BeforeTest(suiteName, testName string) {
	suite.resetMockState()
}

// resetMockState clears all mock expectations and call history
// This is used both by BeforeTest (for test method isolation) and within test cases (for test case isolation)
func (suite *InvocationServiceTestSuite) resetMockState() {
	if suite.mockProviderProvider != nil {
		suite.mockProviderProvider.ExpectedCalls = nil
		suite.mockProviderProvider.Calls = nil
	}
	if suite.mockAdapterProvider != nil {
		suite.mockAdapterProvider.ExpectedCalls = nil
		suite.mockAdapterProvider.Calls = nil
	}
	if suite.mockCredentialProvider != nil {
		suite.mockCredentialProvider.ExpectedCalls = nil
		suite.mockCredentialProvider.Calls = nil
	}
	if suite.mockInvocationRepo != nil {
		suite.mockInvocationRepo.ExpectedCalls = nil
		suite.mockInvocationRepo.Calls = nil
	}
	if suite.mockEventBus != nil {
		suite.mockEventBus.ExpectedCalls = nil
		suite.mockEventBus.Calls = nil
	}
	if suite.mockRedisClient != nil {
		suite.mockRedisClient.ExpectedCalls = nil
		suite.mockRedisClient.Calls = nil
	}
	if suite.mockTokenRefreshService != nil {
		suite.mockTokenRefreshService.ExpectedCalls = nil
		suite.mockTokenRefreshService.Calls = nil
	}
	if suite.mockProviderAdapter != nil {
		suite.mockProviderAdapter.ExpectedCalls = nil
		suite.mockProviderAdapter.Calls = nil
	}
}

// TearDownTest verifies all mock expectations were met
func (suite *InvocationServiceTestSuite) TearDownTest() {
	if suite.mockProviderProvider != nil {
		suite.mockProviderProvider.AssertExpectations(suite.T())
	}
	if suite.mockAdapterProvider != nil {
		suite.mockAdapterProvider.AssertExpectations(suite.T())
	}
	if suite.mockCredentialProvider != nil {
		suite.mockCredentialProvider.AssertExpectations(suite.T())
	}
	if suite.mockInvocationRepo != nil {
		suite.mockInvocationRepo.AssertExpectations(suite.T())
	}
	if suite.mockEventBus != nil {
		suite.mockEventBus.AssertExpectations(suite.T())
	}
	if suite.mockRedisClient != nil {
		suite.mockRedisClient.AssertExpectations(suite.T())
	}
	if suite.mockTokenRefreshService != nil {
		suite.mockTokenRefreshService.AssertExpectations(suite.T())
	}
	if suite.mockProviderAdapter != nil {
		suite.mockProviderAdapter.AssertExpectations(suite.T())
	}
}

// TestInvokeOperation tests the core operation invocation functionality
func (suite *InvocationServiceTestSuite) TestInvokeOperation() {
	testCases := []struct {
		name          string
		setupMocks    func(*InvocationServiceTestSuite)
		userID        string
		providerID    string
		operationID   string
		params        map[string]interface{}
		expectedError error
		assertions    func(*InvocationServiceTestSuite, *domain.Invocation, error)
	}{
		{
			name: "successful_oauth_operation_invocation",
			setupMocks: func(s *InvocationServiceTestSuite) {
				// Mock provider lookup
				provider := &providerDomain.Provider{
					ID:         "test-provider-id",
					Identifier: s.testProviderIdentifier,
					Name:       "Test Provider",
					Status:     types.ProviderStatusActive,
				}
				s.mockProviderProvider.On("FindProvider", mock.Anything, s.testProviderIdentifier).Return(provider, nil)

				// Mock adapter lookup with OAuth auth type
				adapterInfo := &providerAdapterDomain.ProviderAdapterInfo{
					Identifier: s.testProviderIdentifier,
					Name:       "Test Adapter",
					AuthType:   types.AuthTypeOAuth,
				}
				s.mockProviderAdapter.On("GetProviderAdapterInfo").Return(adapterInfo)
				s.mockAdapterProvider.On("GetAdapter", mock.Anything, s.testProviderIdentifier).Return(s.mockProviderAdapter, nil)

				// Mock OAuth lock acquisition
				lockKey := "access_token_lock:test-provider:test-user-123"
				s.mockRedisClient.On("AcquireLock", mock.Anything, lockKey, time.Second*1).Return(true, nil)
				s.mockRedisClient.On("ReleaseLock", mock.Anything, lockKey).Return(nil)

				// Mock credential retrieval
				oauthCred := &credentialDomain.OAuthCredential{
					Credential: &credentialDomain.Credential{
						ID:                 "test-cred-id",
						UserID:             s.testUserID,
						ProviderIdentifier: s.testProviderIdentifier,
						Type:               credentialDomain.CredentialTypeOAuth,
					},
					Token: &oauth2.Token{
						AccessToken:  "test-access-token",
						RefreshToken: "test-refresh-token",
						TokenType:    "Bearer",
						Expiry:       time.Now().Add(time.Hour),
					},
					Scopes: []string{"read", "write"},
				}
				s.mockCredentialProvider.On("GetCredentialByUserAndProvider", mock.Anything, s.testUserID, s.testProviderIdentifier).Return(oauthCred, nil)

				// Mock token refresh
				s.mockTokenRefreshService.On("RefreshAccessTokenIfNeeded", mock.Anything, s.testProviderIdentifier, mock.Anything).Return(oauthCred, nil)

				// Mock invocation repository
				s.mockInvocationRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Invocation")).Return(nil)
				s.mockInvocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Invocation")).Return(nil)

				// Mock adapter execution - use actual params from test case
				expectedResult := map[string]interface{}{
					"status": "success",
					"data":   "test-result",
				}
				s.mockProviderAdapter.On("Execute", mock.Anything, s.testOperationID, map[string]interface{}{"param1": "value1"}, mock.Anything).Return(expectedResult, nil)

				// Mock credential update
				s.mockCredentialProvider.On("UpdateCredentialLastUsedAt", mock.Anything, mock.Anything).Return(nil)

				// Mock event publishing
				s.mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("events.Event")).Return(nil)
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			operationID:   "test-operation",
			params:        map[string]interface{}{"param1": "value1"},
			expectedError: nil,
			assertions: func(s *InvocationServiceTestSuite, result *domain.Invocation, err error) {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), result)
				assert.Equal(s.T(), s.testUserID, result.UserID)
				assert.Equal(s.T(), s.testProviderIdentifier, result.ProviderIdentifier)
				assert.Equal(s.T(), s.testOperationID, result.OperationIdentifier)
				assert.Equal(s.T(), domain.InvocationStatusSuccess, result.Status)
				assert.NotEmpty(s.T(), result.ID)
			},
		},
		{
			name: "successful_none_auth_operation_invocation",
			setupMocks: func(s *InvocationServiceTestSuite) {
				// Mock provider lookup
				provider := &providerDomain.Provider{
					ID:         "test-provider-id",
					Identifier: s.testProviderIdentifier,
					Name:       "Test Provider",
					Status:     types.ProviderStatusActive,
				}
				s.mockProviderProvider.On("FindProvider", mock.Anything, s.testProviderIdentifier).Return(provider, nil)

				// Mock adapter lookup with None auth type
				adapterInfo := &providerAdapterDomain.ProviderAdapterInfo{
					Identifier: s.testProviderIdentifier,
					Name:       "Test Adapter",
					AuthType:   types.AuthTypeNone,
				}
				s.mockProviderAdapter.On("GetProviderAdapterInfo").Return(adapterInfo)
				s.mockAdapterProvider.On("GetAdapter", mock.Anything, s.testProviderIdentifier).Return(s.mockProviderAdapter, nil)

				// Mock none credential creation
				noneCred := &credentialDomain.NoneCredential{
					Credential: &credentialDomain.Credential{
						ID:                 "test-none-cred-id",
						UserID:             s.testUserID,
						ProviderIdentifier: s.testProviderIdentifier,
						Type:               credentialDomain.CredentialTypeNone,
					},
				}
				s.mockCredentialProvider.On("CreateNone", mock.Anything, s.testUserID, s.testProviderIdentifier).Return(noneCred, nil)

				// Mock invocation repository
				s.mockInvocationRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Invocation")).Return(nil)
				s.mockInvocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Invocation")).Return(nil)

				// Mock adapter execution - use actual params from test case
				expectedResult := map[string]interface{}{
					"status": "success",
					"data":   "public-data",
				}
				s.mockProviderAdapter.On("Execute", mock.Anything, s.testOperationID, map[string]interface{}{"param1": "value1"}, mock.Anything).Return(expectedResult, nil)

				// Mock credential update
				s.mockCredentialProvider.On("UpdateCredentialLastUsedAt", mock.Anything, mock.Anything).Return(nil)

				// Mock event publishing
				s.mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("events.Event")).Return(nil)
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			operationID:   "test-operation",
			params:        map[string]interface{}{"param1": "value1"},
			expectedError: nil,
			assertions: func(s *InvocationServiceTestSuite, result *domain.Invocation, err error) {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), result)
				assert.Equal(s.T(), s.testUserID, result.UserID)
				assert.Equal(s.T(), s.testProviderIdentifier, result.ProviderIdentifier)
				assert.Equal(s.T(), domain.InvocationStatusSuccess, result.Status)
			},
		},
		{
			name: "provider_not_found_error",
			setupMocks: func(s *InvocationServiceTestSuite) {
				s.mockProviderProvider.On("FindProvider", mock.Anything, s.testProviderIdentifier).Return(nil, errors.New("provider not found"))
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			operationID:   "test-operation",
			params:        map[string]interface{}{"param1": "value1"},
			expectedError: ErrProviderNotFound,
			assertions: func(s *InvocationServiceTestSuite, result *domain.Invocation, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Contains(s.T(), err.Error(), "provider not found")
			},
		},
		{
			name: "adapter_not_found_error",
			setupMocks: func(s *InvocationServiceTestSuite) {
				// Mock provider lookup
				provider := &providerDomain.Provider{
					ID:         "test-provider-id",
					Identifier: s.testProviderIdentifier,
					Name:       "Test Provider",
					Status:     types.ProviderStatusActive,
				}
				s.mockProviderProvider.On("FindProvider", mock.Anything, s.testProviderIdentifier).Return(provider, nil)

				// Mock adapter not found
				s.mockAdapterProvider.On("GetAdapter", mock.Anything, s.testProviderIdentifier).Return(nil, errors.New("adapter not found"))
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			operationID:   "test-operation",
			params:        map[string]interface{}{"param1": "value1"},
			expectedError: ErrProviderAdapterNotFound,
			assertions: func(s *InvocationServiceTestSuite, result *domain.Invocation, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Contains(s.T(), err.Error(), "provider adapter not found")
			},
		},
		{
			name: "oauth_lock_acquisition_failure",
			setupMocks: func(s *InvocationServiceTestSuite) {
				// Mock provider lookup
				provider := &providerDomain.Provider{
					ID:         "test-provider-id",
					Identifier: s.testProviderIdentifier,
					Name:       "Test Provider",
					Status:     types.ProviderStatusActive,
				}
				s.mockProviderProvider.On("FindProvider", mock.Anything, s.testProviderIdentifier).Return(provider, nil)

				// Mock adapter lookup with OAuth auth type
				adapterInfo := &providerAdapterDomain.ProviderAdapterInfo{
					Identifier: s.testProviderIdentifier,
					Name:       "Test Adapter",
					AuthType:   types.AuthTypeOAuth,
				}
				s.mockProviderAdapter.On("GetProviderAdapterInfo").Return(adapterInfo)
				s.mockAdapterProvider.On("GetAdapter", mock.Anything, s.testProviderIdentifier).Return(s.mockProviderAdapter, nil)

				// Mock OAuth lock acquisition failure
				lockKey := "access_token_lock:test-provider:test-user-123"
				s.mockRedisClient.On("AcquireLock", mock.Anything, lockKey, time.Second*1).Return(false, nil)
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			operationID:   "test-operation",
			params:        map[string]interface{}{"param1": "value1"},
			expectedError: ErrOperationNotFound,
			assertions: func(s *InvocationServiceTestSuite, result *domain.Invocation, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Contains(s.T(), err.Error(), "please try again later")
			},
		},
		{
			name: "credential_not_found_error",
			setupMocks: func(s *InvocationServiceTestSuite) {
				// Mock provider lookup
				provider := &providerDomain.Provider{
					ID:         "test-provider-id",
					Identifier: s.testProviderIdentifier,
					Name:       "Test Provider",
					Status:     types.ProviderStatusActive,
				}
				s.mockProviderProvider.On("FindProvider", mock.Anything, s.testProviderIdentifier).Return(provider, nil)

				// Mock adapter lookup with OAuth auth type
				adapterInfo := &providerAdapterDomain.ProviderAdapterInfo{
					Identifier: s.testProviderIdentifier,
					Name:       "Test Adapter",
					AuthType:   types.AuthTypeOAuth,
				}
				s.mockProviderAdapter.On("GetProviderAdapterInfo").Return(adapterInfo)
				s.mockAdapterProvider.On("GetAdapter", mock.Anything, s.testProviderIdentifier).Return(s.mockProviderAdapter, nil)

				// Mock OAuth lock acquisition
				lockKey := "access_token_lock:test-provider:test-user-123"
				s.mockRedisClient.On("AcquireLock", mock.Anything, lockKey, time.Second*1).Return(true, nil)
				s.mockRedisClient.On("ReleaseLock", mock.Anything, lockKey).Return(nil)

				// Mock credential not found
				s.mockCredentialProvider.On("GetCredentialByUserAndProvider", mock.Anything, s.testUserID, s.testProviderIdentifier).Return(nil, nil)
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			operationID:   "test-operation",
			params:        map[string]interface{}{"param1": "value1"},
			expectedError: ErrCredentialNotFound,
			assertions: func(s *InvocationServiceTestSuite, result *domain.Invocation, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Equal(s.T(), ErrCredentialNotFound, err)
			},
		},
		{
			name: "adapter_execution_failure",
			setupMocks: func(s *InvocationServiceTestSuite) {
				// Mock provider lookup
				provider := &providerDomain.Provider{
					ID:         "test-provider-id",
					Identifier: s.testProviderIdentifier,
					Name:       "Test Provider",
					Status:     types.ProviderStatusActive,
				}
				s.mockProviderProvider.On("FindProvider", mock.Anything, s.testProviderIdentifier).Return(provider, nil)

				// Mock adapter lookup with None auth type
				adapterInfo := &providerAdapterDomain.ProviderAdapterInfo{
					Identifier: s.testProviderIdentifier,
					Name:       "Test Adapter",
					AuthType:   types.AuthTypeNone,
				}
				s.mockProviderAdapter.On("GetProviderAdapterInfo").Return(adapterInfo)
				s.mockAdapterProvider.On("GetAdapter", mock.Anything, s.testProviderIdentifier).Return(s.mockProviderAdapter, nil)

				// Mock none credential creation
				noneCred := &credentialDomain.NoneCredential{
					Credential: &credentialDomain.Credential{
						ID:                 "test-none-cred-id",
						UserID:             s.testUserID,
						ProviderIdentifier: s.testProviderIdentifier,
						Type:               credentialDomain.CredentialTypeNone,
					},
				}
				s.mockCredentialProvider.On("CreateNone", mock.Anything, s.testUserID, s.testProviderIdentifier).Return(noneCred, nil)

				// Mock invocation repository
				s.mockInvocationRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Invocation")).Return(nil)
				s.mockInvocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Invocation")).Return(nil)

				// Mock adapter execution failure - use actual params from test case
				s.mockProviderAdapter.On("Execute", mock.Anything, s.testOperationID, map[string]interface{}{"param1": "value1"}, mock.Anything).Return(nil, errors.New("execution failed"))

				// Mock event publishing for both started and failed events
				s.mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("events.Event")).Return(nil)
			},
			userID:        "test-user-123",
			providerID:    "test-provider",
			operationID:   "test-operation",
			params:        map[string]interface{}{"param1": "value1"},
			expectedError: ErrAdapterExecuteFailed,
			assertions: func(s *InvocationServiceTestSuite, result *domain.Invocation, err error) {
				require.Error(s.T(), err)
				require.NotNil(s.T(), result) // Invocation record is created even on failure
				assert.Contains(s.T(), err.Error(), "adapter execution failed")
				assert.Equal(s.T(), domain.InvocationStatusFailed, result.Status)
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
			result, err := suite.service.InvokeOperation(
				suite.testContext,
				tc.userID,
				tc.providerID,
				tc.operationID,
				tc.params,
			)

			// Run assertions
			tc.assertions(suite, result, err)
		})
	}
}

// TestGetInvocationByID tests invocation retrieval by ID
func (suite *InvocationServiceTestSuite) TestGetInvocationByID() {
	testCases := []struct {
		name          string
		setupMocks    func(*InvocationServiceTestSuite)
		invocationID  string
		expectedError error
		assertions    func(*InvocationServiceTestSuite, *domain.Invocation, error)
	}{
		{
			name: "successful_invocation_retrieval",
			setupMocks: func(s *InvocationServiceTestSuite) {
				expectedInvocation := domain.NewInvocation(
					"test-invocation-id",
					s.testUserID,
					s.testProviderIdentifier,
					s.testOperationID,
					s.testParams,
				)
				expectedInvocation.SetSuccess([]byte(`{"result": "success"}`))

				s.mockInvocationRepo.On("GetByID", mock.Anything, "test-invocation-id").Return(expectedInvocation, nil)
			},
			invocationID:  "test-invocation-id",
			expectedError: nil,
			assertions: func(s *InvocationServiceTestSuite, result *domain.Invocation, err error) {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), result)
				assert.Equal(s.T(), "test-invocation-id", result.ID)
				assert.Equal(s.T(), s.testUserID, result.UserID)
				assert.Equal(s.T(), domain.InvocationStatusSuccess, result.Status)
			},
		},
		{
			name: "invocation_not_found",
			setupMocks: func(s *InvocationServiceTestSuite) {
				s.mockInvocationRepo.On("GetByID", mock.Anything, "non-existing-id").Return(nil, nil)
			},
			invocationID:  "non-existing-id",
			expectedError: ErrInvocationNotFound,
			assertions: func(s *InvocationServiceTestSuite, result *domain.Invocation, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Equal(s.T(), ErrInvocationNotFound, err)
			},
		},
		{
			name: "repository_error",
			setupMocks: func(s *InvocationServiceTestSuite) {
				s.mockInvocationRepo.On("GetByID", mock.Anything, "test-invocation-id").Return(nil, errors.New("database error"))
			},
			invocationID:  "test-invocation-id",
			expectedError: errors.New("failed to get invocation"),
			assertions: func(s *InvocationServiceTestSuite, result *domain.Invocation, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Contains(s.T(), err.Error(), "failed to get invocation")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			tc.setupMocks(suite)
			result, err := suite.service.GetInvocationByID(suite.testContext, tc.invocationID)
			tc.assertions(suite, result, err)
		})
	}
}

// TestListInvocationsByUserID tests listing invocations for a user
func (suite *InvocationServiceTestSuite) TestListInvocationsByUserID() {
	testCases := []struct {
		name          string
		setupMocks    func(*InvocationServiceTestSuite)
		userID        string
		limit         int
		offset        int
		expectedError error
		assertions    func(*InvocationServiceTestSuite, []*domain.Invocation, error)
	}{
		{
			name: "successful_invocation_listing",
			setupMocks: func(s *InvocationServiceTestSuite) {
				expectedInvocations := []*domain.Invocation{
					domain.NewInvocation("inv-1", s.testUserID, s.testProviderIdentifier, "op-1", s.testParams),
					domain.NewInvocation("inv-2", s.testUserID, s.testProviderIdentifier, "op-2", s.testParams),
				}
				s.mockInvocationRepo.On("ListByUserID", mock.Anything, s.testUserID, 10, 0).Return(expectedInvocations, nil)
			},
			userID:        "test-user-123",
			limit:         10,
			offset:        0,
			expectedError: nil,
			assertions: func(s *InvocationServiceTestSuite, result []*domain.Invocation, err error) {
				require.NoError(s.T(), err)
				require.Len(s.T(), result, 2)
				assert.Equal(s.T(), "inv-1", result[0].ID)
				assert.Equal(s.T(), "inv-2", result[1].ID)
			},
		},
		{
			name: "empty_result",
			setupMocks: func(s *InvocationServiceTestSuite) {
				s.mockInvocationRepo.On("ListByUserID", mock.Anything, "empty-user", 10, 0).Return([]*domain.Invocation{}, nil)
			},
			userID:        "empty-user",
			limit:         10,
			offset:        0,
			expectedError: nil,
			assertions: func(s *InvocationServiceTestSuite, result []*domain.Invocation, err error) {
				require.NoError(s.T(), err)
				require.Len(s.T(), result, 0)
			},
		},
		{
			name: "repository_error",
			setupMocks: func(s *InvocationServiceTestSuite) {
				s.mockInvocationRepo.On("ListByUserID", mock.Anything, s.testUserID, 10, 0).Return(nil, errors.New("database error"))
			},
			userID:        "test-user-123",
			limit:         10,
			offset:        0,
			expectedError: errors.New("failed to list invocations"),
			assertions: func(s *InvocationServiceTestSuite, result []*domain.Invocation, err error) {
				require.Error(s.T(), err)
				require.Nil(s.T(), result)
				assert.Contains(s.T(), err.Error(), "failed to list invocations")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			tc.setupMocks(suite)
			result, err := suite.service.ListInvocationsByUserID(suite.testContext, tc.userID, tc.limit, tc.offset)
			tc.assertions(suite, result, err)
		})
	}
}

// TestCountInvocationsByUserID tests counting invocations for a user
func (suite *InvocationServiceTestSuite) TestCountInvocationsByUserID() {
	testCases := []struct {
		name          string
		setupMocks    func(*InvocationServiceTestSuite)
		userID        string
		expectedError error
		assertions    func(*InvocationServiceTestSuite, int64, error)
	}{
		{
			name: "successful_invocation_count",
			setupMocks: func(s *InvocationServiceTestSuite) {
				s.mockInvocationRepo.On("CountByUserID", mock.Anything, s.testUserID).Return(int64(42), nil)
			},
			userID:        "test-user-123",
			expectedError: nil,
			assertions: func(s *InvocationServiceTestSuite, result int64, err error) {
				require.NoError(s.T(), err)
				assert.Equal(s.T(), int64(42), result)
			},
		},
		{
			name: "zero_count",
			setupMocks: func(s *InvocationServiceTestSuite) {
				s.mockInvocationRepo.On("CountByUserID", mock.Anything, "empty-user").Return(int64(0), nil)
			},
			userID:        "empty-user",
			expectedError: nil,
			assertions: func(s *InvocationServiceTestSuite, result int64, err error) {
				require.NoError(s.T(), err)
				assert.Equal(s.T(), int64(0), result)
			},
		},
		{
			name: "repository_error",
			setupMocks: func(s *InvocationServiceTestSuite) {
				s.mockInvocationRepo.On("CountByUserID", mock.Anything, s.testUserID).Return(int64(0), errors.New("database error"))
			},
			userID:        "test-user-123",
			expectedError: errors.New("failed to count invocations"),
			assertions: func(s *InvocationServiceTestSuite, result int64, err error) {
				require.Error(s.T(), err)
				assert.Equal(s.T(), int64(0), result)
				assert.Contains(s.T(), err.Error(), "failed to count invocations")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset mocks to ensure clean state for each test case
			suite.resetMockState()

			tc.setupMocks(suite)
			result, err := suite.service.CountInvocationsByUserID(suite.testContext, tc.userID)
			tc.assertions(suite, result, err)
		})
	}
}

// TestInvocationServiceTestSuite runs the invocation service test suite
func TestInvocationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InvocationServiceTestSuite))
}
