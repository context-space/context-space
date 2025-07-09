package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Config represents the application configuration
type Config struct {
	Environment   string              `json:"environment"`
	Server        ServerConfig        `json:"server"`
	Supabase      SupabaseConfig      `json:"supabase"`
	Database      DatabaseConfig      `json:"database"`
	Redis         RedisConfig         `json:"redis"`
	Vault         VaultConfig         `json:"vault"`
	Logging       LoggingConfig       `json:"logging"`
	Provider      ProviderConfig      `json:"provider"`
	Observability ObservabilityConfig `json:"observability"`
	Security      SecurityConfig      `json:"security"`
}

// ServerConfig holds the server specific configuration
type ServerConfig struct {
	Address string `json:"address"`
}

// SupabaseConfig contains configuration for Supabase Auth
type SupabaseConfig struct {
	ProjectRef  string `json:"project_ref"`
	ServiceRole string `json:"service_role"`
	JWTSecret   string `json:"jwt_secret"`
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Database          string `json:"database"`
	SSLMode           string `json:"ssl_mode"`
	MigrationUsername string `json:"migration_username"`
	MigrationPassword string `json:"migration_password"`
}

// RedisConfig holds the Redis configuration
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// VaultRegionalConfig represents configuration for a regional Vault instance
type VaultRegionalConfig struct {
	Address     string `json:"address"`
	Token       string `json:"token"`
	TransitPath string `json:"transit_path"`
}

// VaultConfig defines configuration for Vault service
type VaultConfig struct {
	Regions       map[string]*VaultRegionalConfig `json:"regions"`
	DefaultRegion string                          `json:"default_region"`
}

// LoggingConfig holds the logging configuration
type LoggingConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	OutputPath string `json:"output_path"`
}

// ProviderConfig holds the provider configuration
type ProviderConfig struct {
	OAuthRedirectURL string `json:"oauth_redirect_url"`
}

// ObservabilityConfig holds the observability configuration
type ObservabilityConfig struct {
	Tracing TracingConfig `json:"tracing"`
	Metrics MetricsConfig `json:"metrics"`
}

// TracingConfig holds tracing specific configuration
type TracingConfig struct {
	Enabled      bool    `json:"enabled"`
	Endpoint     string  `json:"endpoint"`
	SamplingRate float64 `json:"sampling_rate"`
}

// MetricsConfig holds metrics specific configuration
type MetricsConfig struct {
	Enabled  bool   `json:"enabled"`
	Endpoint string `json:"endpoint"`
}

// SecurityConfig holds security specific configuration
type SecurityConfig struct {
	RedirectURLValidator RedirectURLValidatorConfig `json:"redirect_url_validator"`
	CORS                 CORSConfig                 `json:"cors"`
}

type RedirectURLValidatorConfig struct {
	AllowedDomains []string `json:"allowed_domains"`
	AllowedSchemes []string `json:"allowed_schemes"`
}

type CORSConfig struct {
	AllowedOrigins []string `json:"allowed_origins"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	// Load from configuration file
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	// Default config
	config := &Config{
		Environment: env,
		Server: ServerConfig{
			Address: ":8080",
		},
		Vault: VaultConfig{
			DefaultRegion: "eu",
			Regions:       make(map[string]*VaultRegionalConfig),
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			OutputPath: "stdout",
		},
		Security: SecurityConfig{
			RedirectURLValidator: RedirectURLValidatorConfig{
				AllowedDomains: []string{},
				AllowedSchemes: []string{},
			},
			CORS: CORSConfig{
				AllowedOrigins: []string{},
			},
		},
	}

	var configFile string
	switch env {
	case "local":
		configFile = "configs/local.json"
	case "development":
		configFile = "configs.dev/development.json"
	case "production":
		configFile = "configs.prod/production.json"
	default:
		return nil, fmt.Errorf("invalid environment: %s", env)
	}

	if _, err := os.Stat(configFile); err != nil {
		return nil, fmt.Errorf("error stat config file: %w", err)
	}
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}

	// Override with environment variables
	overrideWithEnv(config)

	return config, nil
}

// overrideWithEnv overrides config values with environment variables
func overrideWithEnv(config *Config) {
	// Server config
	if envVal := os.Getenv("SERVER_ADDRESS"); envVal != "" {
		config.Server.Address = envVal
	}

	// Supabase config
	if envVal := os.Getenv("SUPABASE_PROJECT_REF"); envVal != "" {
		config.Supabase.ProjectRef = envVal
	}
	if envVal := os.Getenv("SUPABASE_SERVICE_ROLE"); envVal != "" {
		config.Supabase.ServiceRole = envVal
	}
	if envVal := os.Getenv("SUPABASE_JWT_SECRET"); envVal != "" {
		config.Supabase.JWTSecret = envVal
	}

	// Database config
	if envVal := os.Getenv("DB_HOST"); envVal != "" {
		config.Database.Host = envVal
	}
	if envVal := os.Getenv("DB_PORT"); envVal != "" {
		fmt.Sscanf(envVal, "%d", &config.Database.Port)
	}
	if envVal := os.Getenv("DB_USERNAME"); envVal != "" {
		config.Database.Username = envVal
	}
	if envVal := os.Getenv("DB_PASSWORD"); envVal != "" {
		config.Database.Password = envVal
	}
	if envVal := os.Getenv("DB_DATABASE"); envVal != "" {
		config.Database.Database = envVal
	}
	if envVal := os.Getenv("DB_SSL_MODE"); envVal != "" {
		config.Database.SSLMode = envVal
	}
	if envVal := os.Getenv("DB_MIGRATION_USERNAME"); envVal != "" {
		config.Database.MigrationUsername = envVal
	}
	if envVal := os.Getenv("DB_MIGRATION_PASSWORD"); envVal != "" {
		config.Database.MigrationPassword = envVal
	}

	// Redis config
	if envVal := os.Getenv("REDIS_HOST"); envVal != "" {
		config.Redis.Host = envVal
	}
	if envVal := os.Getenv("REDIS_PORT"); envVal != "" {
		fmt.Sscanf(envVal, "%d", &config.Redis.Port)
	}
	if envVal := os.Getenv("REDIS_USERNAME"); envVal != "" {
		config.Redis.Username = envVal
	}
	if envVal := os.Getenv("REDIS_PASSWORD"); envVal != "" {
		config.Redis.Password = envVal
	}
	if envVal := os.Getenv("REDIS_DB"); envVal != "" {
		fmt.Sscanf(envVal, "%d", &config.Redis.DB)
	}

	// Vault config
	if envVal := os.Getenv("VAULT_DEFAULT_REGION"); envVal != "" {
		config.Vault.DefaultRegion = envVal
	}
	if envVal := os.Getenv("VAULT_DEFAULT_REGION_TOKEN"); envVal != "" {
		config.Vault.Regions[config.Vault.DefaultRegion].Token = envVal
	}

	// Logging config
	if envVal := os.Getenv("LOGGING_LEVEL"); envVal != "" {
		config.Logging.Level = envVal
	}
	if envVal := os.Getenv("LOGGING_FORMAT"); envVal != "" {
		config.Logging.Format = envVal
	}
	if envVal := os.Getenv("LOGGING_OUTPUT_PATH"); envVal != "" {
		config.Logging.OutputPath = envVal
	}

	// Observability config
	if envVal := os.Getenv("TRACING_ENABLED"); envVal != "" {
		config.Observability.Tracing.Enabled = strings.ToLower(envVal) == "true"
	}
	if envVal := os.Getenv("TRACING_ENDPOINT"); envVal != "" {
		config.Observability.Tracing.Endpoint = envVal
	}
	if envVal := os.Getenv("TRACING_SAMPLING_RATE"); envVal != "" {
		fmt.Sscanf(envVal, "%f", &config.Observability.Tracing.SamplingRate)
	}
	if envVal := os.Getenv("METRICS_ENABLED"); envVal != "" {
		config.Observability.Metrics.Enabled = strings.ToLower(envVal) == "true"
	}
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.Username, c.Database.Password, c.Database.Database, c.Database.SSLMode)
}

// GetRedisAddr returns the Redis connection string
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}
