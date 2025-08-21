# Contributing to Context Space Backend

Thank you for your interest in contributing to Context Space Backend! This guide will help you understand how to contribute new third-party service adapters to our integration platform.

## Table of Contents
- [Getting Started](#getting-started)
- [Adding a New Provider Adapter](#adding-a-new-provider-adapter)
- [Code Standards](#code-standards)
- [Documentation Requirements](#documentation-requirements)
- [Submission Process](#submission-process)
- [Community Guidelines](#community-guidelines)
- [Security Considerations](#security-considerations)
- [Advanced Topics](#advanced-topics)

## Getting Started

### Prerequisites
- Go 1.24 or higher
- Docker & Docker Compose
- Git
- Basic understanding of OAuth 2.0 and REST APIs
- Familiarity with the third-party service you want to integrate

### Development Environment Setup
1. **Fork and Clone the Repository from GitHub**

https://github.com/context-space/context-space

2. **Enter Project Directory**

```
cd backend
```

3. **Set Up Local Development**
   ```bash
   # Install dependencies
   make deps
   
   # Start development environment
   make docker-db
   sleep 5
   make docker-migrate
   
   # Create development configuration file
   mkdir -p configs.dev
   cat > configs.dev/development.json << 'EOF'
{
    "environment": "development",
    "database": {
        "host": "postgres",
        "port": 5432,
        "username": "postgres",
        "password": "postgres",
        "database": "contextspace",
        "ssl_mode": "disable",
        "migration_username": "postgres",
        "migration_password": "postgres"
    },
    "vault": {
        "default_region": "eu",
        "regions": {
            "eu": {
                "address": "http://vault:8200",
                "token": "dev-root-token",
                "transit_path": "transit"
            }
        }
    },
    "security": {
        "redirect_url_validator": {
            "allowed_domains": [],
            "allowed_schemes": []
        },
        "cors": {
            "allowed_origins": ["http://localhost:4321", "*"]
        }
    }
}
EOF
   
   # Load providers into database (CRITICAL STEP)
   DB_HOST=localhost DB_PORT=5433 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=contextspace DB_SSL_MODE=disable go run cmd/load_providers/main.go --all
   
   # Start backend API
   make docker-api
   ```

4. **Install Required Tools**
   ```bash
   # Install mockery for generating test mocks
   go install github.com/vektra/mockery/v2@latest
   
   # Install swag for API documentation generation
   go install github.com/swaggo/swag/cmd/swag@latest
   
   # Add Go bin to PATH (add to your ~/.zshrc or ~/.bashrc)
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

5. **Verify Setup**
   
   ```bash
   # Generate mocks
   make mocks

   # Run tests to ensure everything works
   make test
   
   # Check backend API is running
   curl http://localhost:8080/v1/swagger.json
   ```

## Adding a New Provider Adapter

### Step 1: Research and Planning

Before starting development, ensure you have:
- [ ] Access to the target service's API documentation
- [ ] Understanding of the service's authentication method (OAuth, API Key, etc.)
- [ ] List of core operations you want to support
- [ ] Knowledge of required permissions/scopes

### Step 2: Create Provider Structure

1. **Create Directory Structure**
   ```bash
   # Replace 'newservice' with your service name (lowercase, no spaces)
   mkdir -p internal/provideradapter/infrastructure/adapters/newservice
   mkdir -p configs/providers/newservice/i18n
   ```

2. **Understanding Configuration Files**

   **Two types of configuration files are required for each provider:**

   #### **manifest.json - Technical Configuration**
   - **Purpose**: Defines the technical specifications and capabilities of your adapter
   - **Usage**: Used by the backend system to understand how to interact with your service
   - **Content**: API endpoints, authentication methods, operation definitions, parameter schemas
   - **Language**: Single file, English only (technical identifiers)
   - **Processing**: Directly loaded and processed by the Go backend

   #### **i18n/*.json - User Interface Translations**
   - **Purpose**: Provides human-readable labels and descriptions for the UI
   - **Usage**: Used by the frontend to display service information to users
   - **Content**: Display names, descriptions, help text
   - **Language**: Multiple files for different locales (en.json, zh-CN.json, zh-TW.json)
   - **Processing**: Served to frontend based on user's language preference

   #### **Key Differences**

   | Aspect | manifest.json | i18n/*.json |
   |--------|---------------|-------------|
   | **Audience** | Backend system | Frontend users |
   | **Content** | Technical definitions | Human-readable text |
   | **Language** | English (identifiers) | Multiple languages |
   | **Structure** | Complete API schema | UI display information only |
   | **Fields** | All technical fields | Only user-facing fields |

3. **Create Configuration Files**
   
   **Create `configs/providers/newservice/manifest.json`:**
   ```json
   {
     "identifier": "newservice",
     "name": "New Service",
     "description": "Description of the new service integration",
     "auth_type": "oauth",
     "status": "active",
     "icon_url": "",
     "categories": ["category1", "category2"],
     "permissions": [
       {
         "identifier": "read_data",
         "name": "Read Data",
         "description": "Read access to service data",
         "oauth_scopes": ["scope1", "scope2"]
       }
     ],
     "operations": [
       {
         "identifier": "list_items",
         "name": "List Items",
         "description": "List all items from the service",
         "category": "data",
         "required_permissions": ["read_data"],
         "http_method": "GET",
         "endpoint_path":"list_items",
         "parameters": [
           {
             "name": "limit",
             "type": "integer",
             "required": false,
             "description": "Maximum number of items to return"
           }
         ]
       }
     ],
     "oauth_config": {
       "client_id": "your_client_id",
       "client_secret": "your_client_secret"
     }
   }
   ```

4. **Create Translation Files**

   The i18n files contain **only the user-facing text** that will be displayed in the frontend UI. Note how these files contain **fewer fields** than manifest.json and focus only on **display information**.
   
   **Create `configs/providers/newservice/i18n/en.json`:**
   ```json
   {
     "name": "New Service",
     "description": "Description of the new service integration",
     "categories": ["Category 1", "Category 2"],
     "permissions": [
       {
         "identifier": "read_data",
         "name": "Read Data",
         "description": "Read access to service data"
       }
     ],
     "operations": [
       {
         "identifier": "list_items",
         "name": "List Items",
         "description": "List all items from the service",
         "parameters": [
           {
             "name": "limit",
             "description": "Maximum number of items to return"
           }
         ]
       }
     ]
   }
   ```
   
   **Create similar files for `zh-CN.json` and `zh-TW.json` with appropriate translations.**

   #### **Important Notes About Translation Files:**

   ✅ **Fields to Include in i18n:**
   - `name` - Service display name
   - `description` - Service description  
   - `categories` - Category display names
   - `permissions[].name` - Permission display names
   - `permissions[].description` - Permission descriptions
   - `operations[].name` - Operation display names
   - `operations[].description` - Operation descriptions
   - `operations[].parameters[].description` - Parameter descriptions

   ❌ **Fields NOT to Include in i18n:**
   - `identifier` - Technical IDs (always English)
   - `auth_type` - Technical values
   - `status` - Technical status
   - `oauth_scopes` - API-specific scopes
   - `required_permissions` - Technical references
   - `parameters[].name` - API parameter names
   - `parameters[].type` - Data types
   - `oauth_config` - Technical configuration

   **Translation Guidelines:**
   - Keep technical terms consistent across languages
   - Use clear, user-friendly language
   - Ensure all three language files have the same structure
   - Test translations with native speakers when possible

   #### **Example: Same Operation in Both File Types**

   **In manifest.json (Technical Definition):**
   ```json
   {
     "identifier": "list_items",
     "name": "List Items",
     "description": "List all items from the service",
     "category": "data",
     "required_permissions": ["read_data"],
     "parameters": [
       {
         "name": "limit",
         "type": "integer",
         "required": false,
         "description": "Maximum number of items to return"
       }
     ]
   }
   ```

   **In i18n/en.json (User Display):**
   ```json
   {
     "identifier": "list_items",
     "name": "List Items",
     "description": "List all items from the service",
     "parameters": [
       {
         "name": "limit",
         "description": "Maximum number of items to return"
       }
     ]
   }
   ```

   **In i18n/zh-CN.json (Chinese Translation):**
   ```json
   {
     "identifier": "list_items",
     "name": "列出项目",
     "description": "列出服务中的所有项目",
     "parameters": [
       {
         "name": "limit",
         "description": "返回的最大项目数量"
       }
     ]
   }
   ```

   Notice how i18n files **omit technical fields** like `category`, `required_permissions`, `type`, and `required`.

---

### Step 3: Implement the Adapter

**The project supports two authentication methods: OAuth 2.0 and API Key. Please choose the appropriate implementation method based on the service you want to integrate.**

#### How to Choose Authentication Method?

**Choose OAuth 2.0 when:**
- Service needs to act on behalf of users
- Service provides fine-grained permission control (scopes)
- Need to access user's private data
- Service supports token refresh mechanism
- Enterprise-level security requirements

**Choose API Key when:**
- Service provides public data or enterprise APIs
- No user personal authorization required
- Service uses simple key authentication
- Quick integration and testing
- Service doesn't support OAuth

#### Detailed Comparison of Supported Authentication Methods

| Authentication Method | Use Case | Configuration Complexity | Security Level | User Experience | Example Services |
|---------------------|----------|--------------------------|----------------|-----------------|------------------|
| **OAuth 2.0** | Services requiring user authorization | High | High | Requires user authorization flow | GitHub, Slack, Zoom, Google APIs |
| **API Key** | Services using direct key access | Low | Medium | No user interaction needed | Weather API, Currency API, News API |

#### Technical Implementation Differences

| Feature | OAuth 2.0 | API Key |
|---------|-----------|---------|
| **Authentication Method** | Bearer Token | API Key Header |
| **Token Refresh** | Supported | Not Applicable |
| **User Authorization** | Required | Not Required |
| **Configuration File** | Requires `permissions` and `oauth_config` | Basic configuration only |
| **Error Handling** | Permission checks + OAuth errors | Basic API errors |

---

#### Configuration File Differences

##### 1. OAuth Configuration Example

**`configs/providers/newservice/manifest.json` (OAuth):**
```json
{
  "identifier": "newservice",
  "name": "New Service",
  "description": "Description of the new service integration",
  "auth_type": "oauth",
  "status": "active",
  "icon_url": "",
  "categories": ["category1", "category2"],
  "permissions": [
    {
      "identifier": "read_data",
      "name": "Read Data", 
      "description": "Read access to service data",
      "oauth_scopes": ["read:items", "read:user"]
    },
    {
      "identifier": "write_data",
      "name": "Write Data",
      "description": "Write access to service data", 
      "oauth_scopes": ["write:items"]
    }
  ],
  "operations": [
    {
      "identifier": "list_items",
      "name": "List Items",
      "description": "List all items from the service",
      "category": "data",
      "required_permissions": ["read_data"],
      "parameters": [
        {
          "name": "limit",
          "type": "integer",
          "required": false,
          "description": "Maximum number of items to return"
        }
      ]
    }
  ],
  "oauth_config": {
    "client_id": "your_client_id",
    "client_secret": "your_client_secret"
  }
}
```

##### 2. API Key Configuration Example

**`configs/providers/newservice/manifest.json` (API Key):**
```json
{
  "identifier": "newservice",
  "name": "New Service", 
  "description": "Description of the new service integration",
  "auth_type": "apikey",
  "status": "active",
  "icon_url": "",
  "categories": ["category1", "category2"],
  "operations": [
    {
      "identifier": "list_items",
      "name": "List Items",
      "description": "List all items from the service",
      "category": "data",
      "parameters": [
        {
          "name": "limit",
          "type": "integer", 
          "required": false,
          "description": "Maximum number of items to return"
        }
      ]
    }
  ]
}
```

> **Note: API Key authentication does not require `permissions` and `oauth_config` fields.**

---

#### OAuth 2.0 Authentication Implementation

> Reference: `internal/provideradapter/infrastructure/adapters/airtable`

##### 1. **Create OAuth Adapter Implementation**
   
   **Create `internal/provideradapter/infrastructure/adapters/newservice/newservice_adapter.go`:**
```go
    package newservice

    import (
        // ...
    )

    // NewServiceAdapter is an adapter for the New Service API using OAuth2
    type NewServiceAdapter struct {
        *base.BaseAdapter
        oauthConfig   *domain.OAuthConfig
        restAdapter   domain.Adapter             // The underlying REST adapter instance
        operations    Operations                 // Map of operation ID to definition
        permissionSet providercore.PermissionSet // Permission set defined in providercore
    }

    // NewNewServiceAdapter creates a new adapter instance
    func NewNewServiceAdapter(
        providerInfo *domain.ProviderAdapterInfo,
        config *domain.AdapterConfig,
        oauthConfig *domain.OAuthConfig,
        restAdapter domain.Adapter,
        permissions providercore.PermissionSet,
    ) *NewServiceAdapter {
        baseAdapter := base.NewBaseAdapter(providerInfo, config)    
        adapter := &NewServiceAdapter{
    	   BaseAdapter: baseAdapter,
            // ...
        }   
        adapter.registerOperations()
        return adapter
    }

    // createOAuth2Config creates the oauth2.Config object
    func (a *NewServiceAdapter) createOAuth2Config(redirectURL string, scopes []string) *oauth2.Config {
        endpoint := oauth2.Endpoint{
            AuthURL:   authURL,
            TokenURL:  tokenURL,
            AuthStyle: oauth2.AuthStyleInHeader,
        }   
        return &oauth2.Config{
            ClientID:     a.oauthConfig.ClientID,
            ClientSecret: a.oauthConfig.ClientSecret,
            RedirectURL:  redirectURL,
            Scopes:       scopes,
            Endpoint:     endpoint,
        }
    }

    // Implement OAuthAdapter interface methods

    // ShouldRefreshToken checks if the token should be refreshed
    func (a *NewServiceAdapter) ShouldRefreshToken(oldToken *oauth2.Token) bool {
       // ...
    }

    // RefreshOAuthToken refreshes an OAuth token
    func (a *NewServiceAdapter) RefreshOAuthToken(ctx context.Context, oldToken *oauth2.Token) (*oauth2.Token, error) {
       // ...
    }

    // GenerateOAuthURL generates an OAuth authorization URL
    func (a *NewServiceAdapter) GenerateOAuthURL(
        ctx context.Context,
        redirectURL, state, codeChallenge string,
        scopes []string,
    ) (string, error) {
        // ...
    }

    // ExchangeCodeForTokens exchanges an authorization code for tokens
    func (a *NewServiceAdapter) ExchangeCodeForTokens(ctx context.Context, code, redirectURL, codeVerifier string) (*oauth2.Token, error) {
        // ...
    }

    // CheckMissingPermissions checks if required permissions are present in authorized scopes
    func (a *NewServiceAdapter) CheckMissingPermissions(operationIdentifier string, authorizedScopes []string) (bool, []string, error) {
        // ...
    }

    // GetScopesFromPermissions translates internal permission identifiers to required OAuth scopes
    func (a *NewServiceAdapter) GetScopesFromPermissions(permissions []string) ([]string, error) {
        // ...
    }

    // GetPermissionIdentifiersFromScopes translates OAuth scopes to internal permission identifiers
    func (a *NewServiceAdapter) GetPermissionIdentifiersFromScopes(scopes []string) ([]string, error) {
        // ...
    }

    // Implement Adapter interface methods
    // Execute handles API calls based on the operationID using the REST adapter
    func (a *NewServiceAdapter) Execute(
        ctx context.Context,
        operationID string,
        params map[string]interface{},
        credential interface{},
    ) (interface{}, error) {
        // ...
    }
   ```

##### 2. **Create OAuth Adapter Template**
   **Create `internal/provideradapter/infrastructure/adapters/newservice/template.go`:**
```go
    package newservice

    import (
        // ...
    )

    const (
        identifier = "newservice"
        baseURL    = "https://api.newservice.com/v1"
        
        // OAuth endpoints
        authURL  = "https://newservice.com/oauth/authorize"
        tokenURL = "https://newservice.com/oauth/token"
    )

    // Register the adapter template during package initialization
    func init() {
        // Type assertion to ensure the adapter implements the necessary interfaces
        var _ domain.OAuthAdapter = (*NewServiceAdapter)(nil)
        
        template := &NewServiceAdapterTemplate{}
        registry.RegisterAdapterTemplate(identifier, template)
    }

    // NewServiceAdapterTemplate implements the AdapterTemplate interface
    type NewServiceAdapterTemplate struct {
        // Configuration specific to this template could be added here if needed
    }

    // CreateAdapter creates a new adapter instance from the provided configuration
    func (t *NewServiceAdapterTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
        // ...

        // Create the main adapter
        adapter := NewNewServiceAdapter(
            // ...
        )
        return adapter, nil
    }

    // ValidateConfig checks if the provided configuration contains the necessary fields
    func (t *NewServiceAdapterTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
        // ...
    }
   ```

---

#### API Key Authentication Implementation

> Reference: `internal/provideradapter/infrastructure/adapters/serper`

##### 1. **Create API Key Adapter Implementation**
   **Create `internal/provideradapter/infrastructure/adapters/newservice/newservice_adapter.go`:**

```go
    package newservice

    import (
        // ...
    )

    // NewServiceAdapter is an adapter for the New Service API using API Key
    type NewServiceAdapter struct {
        *base.BaseAdapter
        restAdapter domain.Adapter // The underlying REST adapter instance
        operations  Operations     // Map of operation ID to definition
    }

    // NewNewServiceAdapter creates a new adapter instance for API Key authentication
    func NewNewServiceAdapter(
        providerInfo *domain.ProviderAdapterInfo,
        config *domain.AdapterConfig,
        restAdapter domain.Adapter,
    ) *NewServiceAdapter {
        baseAdapter := base.NewBaseAdapter(providerInfo, config)

        adapter := &NewServiceAdapter{
            BaseAdapter: baseAdapter,
            restAdapter: restAdapter,
            operations:  make(Operations),
        }

        adapter.registerOperations()
        return adapter
    }

    // Execute handles API calls based on the operationID using the REST adapter
    func (a *NewServiceAdapter) Execute(
        ctx context.Context,
        operationID string,
        params map[string]interface{},
        credential interface{},
    ) (interface{}, error) {
        // ...
    }
   ```

##### 2. **Create API Key Adapter Template**
   **Create `internal/provideradapter/infrastructure/adapters/newservice/template.go`:**
```go
    package newservice

    import (
        // ...
    )

    const (
        identifier = "newservice"
        baseURL    = "https://api.newservice.com/v1"
    )

    // Register the adapter template during package initialization
    func init() {
        // Type assertion to ensure the adapter implements the necessary interfaces
        var _ domain.APIKeyAdapter = (*NewServiceAdapter)(nil)

        template := &NewServiceAdapterTemplate{}
        registry.RegisterAdapterTemplate(identifier, template)
    }

    // NewServiceAdapterTemplate implements the AdapterTemplate interface
    type NewServiceAdapterTemplate struct {
        // Configuration specific to this template could be added here if needed
    }

    // CreateAdapter creates a new adapter instance from the provided configuration
    func (t *NewServiceAdapterTemplate) CreateAdapter(provider *domain.ProviderAdapterConfig) (domain.Adapter, error) {
        // ...

        // Create the main adapter (no OAuth config needed for API Key)
        adapter := NewNewServiceAdapter(
            // ...
        )

        return adapter, nil
    }

    // ValidateConfig checks if the provided configuration contains the necessary fields
    func (t *NewServiceAdapterTemplate) ValidateConfig(provider *domain.ProviderAdapterConfig) error {
        // ...
    }
   ```

---

#### **Create Operations Implementation**
   
   **Create `internal/provideradapter/infrastructure/adapters/newservice/newservice_operations.go`:**
   
   This file is shared between OAuth and API Key implementations, with slight differences in the operation registration:

```go
   package newservice

   import (
       "context"
       "fmt"
       "net/http"
   )

   // Define constants for API endpoints
   const (
       endpointListItems    = "/items"
       endpointGetItem      = "/items/{itemId}"
       endpointCreateItem   = "/items"
   )

   // Define constants for operation IDs
   const (
       operationIDListItems  = "list_items"
       operationIDGetItem    = "get_item"
       operationIDCreateItem = "create_item"
   )

   // Parameter structs for each operation
   type ListItemsParams struct {
       Limit  int    `mapstructure:"limit" validate:"omitempty,min=1,max=100"`
       Offset int    `mapstructure:"offset" validate:"omitempty,min=0"`
       Filter string `mapstructure:"filter" validate:"omitempty"`
   }

   type GetItemParams struct {
       ItemID string `mapstructure:"itemId" validate:"required"`
   }

   type CreateItemParams struct {
       Name        string `mapstructure:"name" validate:"required"`
       Description string `mapstructure:"description" validate:"omitempty"`
   }

   // OperationHandler defines the function signature for handling operations
   type OperationHandler func(ctx context.Context, params interface{}) (map[string]interface{}, error)

   // OperationDefinition structure differs based on authentication method:
   
   // For OAuth implementations:
   type OperationDefinition struct {
       Schema                interface{}      // Parameter schema (struct pointer)
       Handler               OperationHandler // Operation handler function
       PermissionIdentifiers []string         // Required permission identifiers (OAuth only)
   }

   // For API Key implementations:
   // type OperationDefinition struct {
   //     Schema  interface{}      // Parameter schema (struct pointer)
   //     Handler OperationHandler // Operation handler function
   // }

   // Operations maps operation IDs to their definitions
   type Operations map[string]OperationDefinition

   // RegisterOperation - OAuth version (with permissions)
   func (a *NewServiceAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler, requiredPerms []string) {
       a.BaseAdapter.RegisterOperation(operationID, schema)
       if a.operations == nil {
           a.operations = make(Operations)
       }
       a.operations[operationID] = OperationDefinition{
           Schema:                schema,
           Handler:               handler,
           PermissionIdentifiers: requiredPerms,
       }
   }

   // RegisterOperation - API Key version (without permissions)
   // func (a *NewServiceAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
   //     a.BaseAdapter.RegisterOperation(operationID, schema)
   //     if a.operations == nil {
   //         a.operations = make(Operations)
   //     }
   //     a.operations[operationID] = OperationDefinition{
   //         Schema:  schema,
   //         Handler: handler,
   //     }
   // }

   // registerOperations - OAuth version
   func (a *NewServiceAdapter) registerOperations() {
       a.RegisterOperation(operationIDListItems, &ListItemsParams{}, handleListItems, []string{"read_items"})
       a.RegisterOperation(operationIDGetItem, &GetItemParams{}, handleGetItem, []string{"read_items"})
       a.RegisterOperation(operationIDCreateItem, &CreateItemParams{}, handleCreateItem, []string{"write_items"})
   }

   // registerOperations - API Key version
   // func (a *NewServiceAdapter) registerOperations() {
   //     a.RegisterOperation(operationIDListItems, &ListItemsParams{}, handleListItems)
   //     a.RegisterOperation(operationIDGetItem, &GetItemParams{}, handleGetItem)
   //     a.RegisterOperation(operationIDCreateItem, &CreateItemParams{}, handleCreateItem)
   // }

   // Handler functions implementation
   // Implement handleListItems, handleGetItem, handleCreateItem functions
   // Each handler should:
   // 1. Validate and cast input parameters
   // 2. Build REST request parameters (method, path, query params, headers, body)
   // 3. Return map[string]interface{} for the REST adapter

   func handleListItems(ctx context.Context, params interface{}) (map[string]interface{}, error) {
       // Implementation: Build GET request to /items with query parameters
       // Return map with "method", "path", "query_params" fields
   }

   func handleGetItem(ctx context.Context, params interface{}) (map[string]interface{}, error) {
       // Implementation: Build GET request to /items/{itemId} with path parameters
       // Return map with "method", "path", "path_params" fields
   }

   func handleCreateItem(ctx context.Context, params interface{}) (map[string]interface{}, error) {
       // Implementation: Build POST request to /items with request body
       // Return map with "method", "path", "body" fields
   }
   ```

---

#### Alternative Authentication Methods
```http
# Bearer token format
Authorization: Bearer your_api_key_here

# API Key format  
Authorization: ApiKey your_api_key_here

# Custom header (check service documentation)
X-RapidAPI-Key: your_api_key_here
```

#### Configuration
No additional configuration needed in manifest.json for API key authentication.

#### Security Best Practices
- **Never commit API keys** to version control
- **Use environment variables** for API key storage
- **Rotate API keys regularly** (if supported by the service)
- **Monitor API key usage** for unusual activity
- **Restrict API key permissions** to minimum required scopes


#### Troubleshooting
- **401 Unauthorized**: Check if API key is valid and not expired
- **403 Forbidden**: Verify API key has required permissions
- **429 Too Many Requests**: You've hit the rate limit, wait before retrying


### Step 4: Verify Adapter Registration

The adapter registration is handled automatically through the template system. When your adapter package is imported, the `init()` function in `template.go` will automatically register your adapter template with the system.

To ensure your adapter is properly registered, you can verify it by checking the logs when the system starts, or by running a simple test:

```go
// Create a simple test to verify your adapter is registered
func TestAdapterRegistration(t *testing.T) {
    // This test should pass if your adapter is properly registered
    template := registry.GetAdapterTemplate("newservice")
    assert.NotNil(t, template, "newservice adapter template should be registered")
}
```

Make sure your adapter package is imported in `internal/provideradapter/infrastructure/templates/init.go`

```go
import (
    // ... other imports
    _ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/newservice"
)
```

### Step 5: Generate SQL and Database Updates

```bash
# Generate SQL files from your provider configuration
go run cmd/load_providers/main.go -path configs/providers -sql -sql-output generated_sql

# Load the generated SQL into your development database
psql -h localhost -U postgres -d contextspace -f generated_sql/providers_inserts.sql
psql -h localhost -U postgres -d contextspace -f generated_sql/operations_inserts.sql
psql -h localhost -U postgres -d contextspace -f generated_sql/provider_adapters_inserts.sql
psql -h localhost -U postgres -d contextspace -f generated_sql/provider_translations_inserts.sql
```

### Step 6: Verify Service Functionality

After completing the technical configuration, you must verify that your new service integration works correctly. This step includes two essential test scenarios:

#### 1. Service Authorization Test

**For OAuth Services:**
```bash
# Start the development server
make docker-api

# Test OAuth authorization flow
curl -X POST "http://localhost:8080/api/v1/oauth/authorize" \
  -H "Content-Type: application/json" \
  -d '{
    "provider_id": "newservice",
    "redirect_uri": "http://localhost:3000/callback",
    "permissions": ["read_data"]
  }'

# Follow the returned authorization URL in browser
# After callback, verify token exchange worked correctly
```

**For API Key Services:**
```bash
# Test API key credential storage
curl -X POST "http://localhost:8080/api/v1/credentials/apikey" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_USER_TOKEN" \
  -d '{
    "provider_id": "newservice",
    "api_key": "your_test_api_key"
  }'

# Verify credential was stored successfully (should return 200 OK)
```

#### 2. Operation Functionality Test

Test each operation defined in your manifest.json:

```bash
# Test list_items operation
curl -X POST "http://localhost:8080/api/v1/invocations/invoke" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_USER_TOKEN" \
  -d '{
    "provider_identifier": "newservice",
    "operation_identifier": "list_items",
    "parameters": {
      "limit": 10
    }
  }'

# Expected response: 200 OK with data from the service
# Verify the response contains expected fields and data structure
```

**Additional Operation Tests:**
- Test all operations defined in your manifest.json
- Verify parameter validation works correctly
- Test error handling for invalid parameters
- Confirm permission requirements are enforced

#### 3. Verification Checklist

Before proceeding, ensure:
- [ ] OAuth authorization flow completes successfully (if applicable)
- [ ] API key authentication works correctly (if applicable)
- [ ] All defined operations return expected responses
- [ ] Parameter validation works as expected
- [ ] Error responses are properly formatted
- [ ] Required permissions are correctly enforced
- [ ] Service responses match expected data structures

#### 4. Common Issues and Troubleshooting

**Authorization Issues:**
- Verify OAuth client credentials are correct
- Check redirect URI matches registered value
- Ensure required scopes are properly configured

**Operation Issues:**
- Verify API endpoints are accessible
- Check request/response format matches service documentation
- Confirm parameter mapping is correct
- Validate error handling covers common scenarios

**Database Issues:**
- Ensure SQL files were applied correctly
- Check provider registration in database
- Verify operation definitions are loaded

---

## Code Standards

This section outlines the coding standards and best practices that all contributors must follow to maintain code quality and consistency.

### Go Code Guidelines
- Follow Go conventions and best practices
- Use meaningful variable and function names
- Include comprehensive error handling
- Add proper logging with OpenTelemetry tracing
- Document all public functions with GoDoc comments

### Example Code Structure
```go
// Execute handles API calls based on the operationID using the REST adapter
func (a *NewServiceAdapter) Execute(
    ctx context.Context,
    operationID string,
    params map[string]interface{},
    credential interface{},
) (interface{}, error) {
    	oauthCred, ok := credential.(*credDomain.OAuthCredential)
	if !ok || oauthCred == nil || oauthCred.Token == nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrCredentialError, "invalid or missing OAuth credential", http.StatusUnauthorized)
	}

	opDef, exists := a.operations[operationID]
	if !exists {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrOperationNotSupported, fmt.Sprintf("unknown operation ID: %s", operationID), http.StatusNotFound)
	}

	authorizedScopes := oauthCred.Scopes

	allScopesPresent, missingIDs, err := a.CheckMissingPermissions(operationID, authorizedScopes)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInternal, fmt.Sprintf("error checking permissions: %v", err), http.StatusInternalServerError)
	}
	if !allScopesPresent {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrCredentialError, fmt.Sprintf("missing required permissions: %v", missingIDs), http.StatusForbidden)
	}

	processedParams, err := a.ProcessParams(operationID, params)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInvalidParameters, fmt.Sprintf("parameter validation failed: %v", err), http.StatusBadRequest)
	}

	handler := opDef.Handler
	restParams, err := handler(ctx, processedParams)
	if err != nil {
		return nil, domain.NewAdapterError(a.GetProviderAdapterInfo().Identifier, operationID, domain.ErrInternal, fmt.Sprintf("operation handler failed: %v", err), http.StatusInternalServerError)
	}

	headers, _ := restParams["headers"].(map[string]string)
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Authorization"] = utils.StringsBuilder("Bearer ", oauthCred.Token.AccessToken)
	restParams["headers"] = headers

	rawResult, err := a.restAdapter.Execute(ctx, operationID, restParams, nil)
	if err != nil {
		return nil, err
	}

	return rawResult, nil
}
```

---

## Documentation Requirements

Proper documentation is essential for maintainability and community adoption. This section covers all documentation standards and requirements.

### 1. Code Documentation
- All public functions must have GoDoc comments
- Complex algorithms should have inline comments
- Include examples in documentation when helpful

### 2. Setup Guide
- Clear step-by-step instructions for creating applications
- Screenshots for complex UI steps (optional but helpful)
- Common troubleshooting issues and solutions

### 3. API Coverage
Document which APIs are supported and any limitations:

```markdown
## Supported Operations
- ✅ List Items - Full support
- ✅ Get Item - Full support  
- ❌ Delete Item - Not supported (requires admin permissions)
- 🔄 Update Item - Partial support (some fields read-only)
```

---

## Submission Process

Follow these steps to properly submit your adapter implementation for review and integration into the project.

### 1. Pre-submission Checklist
- [ ] All tests pass: `make test`
- [ ] Code passes linting: `make lint`
- [ ] Mocks are updated: `make mocks`
- [ ] Documentation is complete
- [ ] SQL files are generated and tested
- [ ] Multi-language translations are provided

### 2. Create Pull Request
1. **Create Feature Branch**
   ```bash
   git checkout -b feature/add-newservice-adapter
   git add .
   git commit -m "feat: add New Service adapter with OAuth support"
   git push origin feature/add-newservice-adapter
   ```

2. **Pull Request Template**
   Use this template for your PR description:
   
   ```markdown
   ## Summary
   Brief description of the new adapter and its capabilities.
   
   ## Service Details
   - **Service Name**: New Service
   - **Authentication Type**: OAuth 2.0
   - **Supported Operations**: List items, Get item details
   - **API Documentation**: https://docs.newservice.com
   
   ## Testing
   - [ ] Unit tests added and passing
   - [ ] Integration tests with real API
   - [ ] Error handling tested
   - [ ] Multi-language support verified
   
   ## Documentation
   - [ ] Setup guide created
   - [ ] API operations documented
   - [ ] Code properly commented
   
   ## Breaking Changes
   None / List any breaking changes
   
   ## Additional Notes
   Any additional context or considerations
   ```

3. **Review Process**
   - Maintainers will review your code for quality and security
   - You may be asked to make changes or provide additional tests
   - Once approved, your adapter will be merged and included in the next release

---

## Community Guidelines

Our community values collaboration, respect, and helping each other build better integrations. Please follow these guidelines to create a positive experience for everyone.

### Code of Conduct
- Be respectful and inclusive in all interactions
- Provide constructive feedback during code reviews
- Help newcomers understand the codebase and contribution process

### Getting Help
- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Documentation**: Check existing docs before asking questions

### Recognition
Contributors will be recognized in:
- Release notes for significant contributions
- Contributors section in the main README
- GitHub contributor graphs

---

## Security Considerations

Security is paramount when dealing with third-party integrations. This section covers essential security practices and requirements.

### Handling Credentials
- Never commit real API keys or secrets
- Use environment variables for testing
- Ensure proper validation of user inputs
- Follow OAuth security best practices

### API Security
- Implement proper rate limiting respect
- Validate all responses from third-party APIs
- Use HTTPS for all API communications
- Handle authentication errors gracefully

---

## Advanced Topics

For complex integrations that go beyond standard REST APIs, this section provides guidance on advanced implementation patterns and architectural considerations.

### Custom Authentication Types
If your service uses a non-standard authentication method, you may need to:
1. Extend the authentication interfaces
2. Add new credential types
3. Update the credential management system

### Complex API Patterns
For APIs with unique patterns (GraphQL, WebSockets, etc.):
1. Document the approach in your adapter
2. Consider creating reusable utilities
3. Discuss with maintainers for architectural guidance

### Performance Optimization
- Implement connection pooling for high-traffic scenarios
- Add caching where appropriate
- Consider pagination for large datasets
- Monitor and optimize API call patterns

---

Thank you for contributing to Context Space Backend! Your efforts help make integration easier for developers worldwide. 🚀