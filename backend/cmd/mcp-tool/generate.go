package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/mcp"
)

type GenerateConfig struct {
	MCPToolConfig
	OutputDir         string
	Identifier        string // Unique identifier for the provider (also used as directory name)
	Name              string // Display name for the provider
	Description       string // Description of the provider
	AuthType          string
	Categories        []string
	EnableTranslation bool // Enable AI translation for i18n files
}

func runGenerateCommand(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf(`Usage: mcp-tool generate [options]

Generate provider adapter configuration from MCP server

Options:
`)
		fs.PrintDefaults()
		fmt.Printf(`
Environment Variables:
  OPENAI_API_KEY     OpenAI API key for AI translation (required when using -translate)
  OPENAI_BASE_URL    Custom OpenAI API base URL (optional, defaults to https://api.openai.com/v1)
  OPENAI_TIMEOUT     Timeout for OpenAI API calls in seconds (optional, defaults to 60)

Examples:
  # Generate adapter with AI translation (requires OPENAI_API_KEY environment variable)
  export OPENAI_API_KEY="your-api-key"
  mcp-tool generate -command npx \
                    -arg "-y" -arg "@modelcontextprotocol/server-brave-search" \
                    -identifier brave_search_i18n \
                    -name "Brave Search" \
                    -description "Multilingual web search powered by Brave" \
                    -auth apikey -categories "search,web" \
                    -translate
`)
	}

	var config GenerateConfig
	var categoriesStr string
	var cmdArgs []string
	var cmdEnvs []string

	// Parse all flags directly on config
	fs.StringVar(&config.Command, "command", "", "Base command to run MCP server (e.g., 'npx', 'uvx', './binary') (required)")

	// Support multiple arguments
	fs.Func("arg", "Command argument (can be specified multiple times)", func(s string) error {
		cmdArgs = append(cmdArgs, s)
		return nil
	})

	// Support multiple environment variables
	fs.Func("env", "Environment variable KEY=VALUE (can be specified multiple times)", func(s string) error {
		cmdEnvs = append(cmdEnvs, s)
		return nil
	})

	fs.StringVar(&config.OutputDir, "output", "configs/providers", "Output directory for provider configurations")
	fs.StringVar(&config.Identifier, "identifier", "", "Unique identifier for the provider (used as directory name) (required)")
	fs.StringVar(&config.Name, "name", "", "Display name for the provider (required)")
	fs.StringVar(&config.Description, "description", "", "Description of the provider (required)")
	fs.StringVar(&config.AuthType, "auth", "none", "Authentication type (oauth, apikey, basic, none)")
	fs.StringVar(&categoriesStr, "categories", "", "Comma-separated list of categories")
	fs.BoolVar(&config.EnableTranslation, "translate", false, "Enable AI translation for i18n files (requires OPENAI_API_KEY)")

	if err := fs.Parse(args); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	// Validate required flags
	if config.Command == "" {
		fmt.Fprintf(os.Stderr, "Error: -command is required\n")
		fs.Usage()
		os.Exit(1)
	}

	if config.Identifier == "" {
		fmt.Fprintf(os.Stderr, "Error: -identifier is required\n")
		fs.Usage()
		os.Exit(1)
	}

	if config.Name == "" {
		fmt.Fprintf(os.Stderr, "Error: -name is required\n")
		fs.Usage()
		os.Exit(1)
	}

	if config.Description == "" {
		fmt.Fprintf(os.Stderr, "Error: -description is required\n")
		fs.Usage()
		os.Exit(1)
	}

	// Use the collected arguments and environment variables directly
	config.Args = cmdArgs

	// Parse environment variables from KEY=VALUE format
	config.Envs = make(map[string]string)
	for _, envStr := range cmdEnvs {
		eqIndex := strings.Index(envStr, "=")
		if eqIndex == -1 {
			log.Fatalf("Invalid environment variable format: %s (expected KEY=VALUE)", envStr)
		}
		key := strings.TrimSpace(envStr[:eqIndex])
		value := envStr[eqIndex+1:]
		if key == "" {
			log.Fatalf("Empty environment variable key in: %s", envStr)
		}
		config.Envs[key] = value
	}

	// Parse categories
	if categoriesStr != "" {
		config.Categories = strings.Split(categoriesStr, ",")
		for i, cat := range config.Categories {
			config.Categories[i] = strings.TrimSpace(cat)
		}
	}

	// Generate the adapter
	if err := generateAdapter(config); err != nil {
		log.Fatalf("Failed to generate adapter: %v", err)
	}

	fmt.Printf("Successfully generated adapter for %s\n", config.Identifier)
}

func generateAdapter(config GenerateConfig) error {
	// Create MCP client configuration
	mcpConfig := mcp.MCPClientConfig{
		Command: config.Command,
		Args:    config.Args,
		Envs:    config.Envs,
		Timeout: 60 * time.Second,
	}

	// Get tools from MCP server
	tools, err := mcp.GetMCPServerTools(context.Background(), mcpConfig, config.Identifier)
	if err != nil {
		return fmt.Errorf("failed to get MCP server tools: %w", err)
	}

	if len(tools) == 0 {
		return fmt.Errorf("no tools found in MCP server")
	}

	fmt.Printf("Found %d tools from MCP server %s %v\n", len(tools), config.Command, config.Args)

	// Convert tools to operations and generate provider files
	operations := convertMCPToolsToOperations(tools)
	manifest := generateManifest(config.Identifier, config.Name, config.Description, config.AuthType, config.Categories, operations)

	// Create output directory structure
	providerDir := filepath.Join(config.OutputDir, config.Identifier)
	if err := os.MkdirAll(providerDir, 0755); err != nil {
		return fmt.Errorf("failed to create provider directory: %w", err)
	}

	// Write manifest.json
	manifestPath := filepath.Join(providerDir, "manifest.json")
	if err := writeJSONFile(manifestPath, manifest); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	// Create i18n directory and files
	i18nDir := filepath.Join(providerDir, "i18n")
	if err := os.MkdirAll(i18nDir, 0755); err != nil {
		return fmt.Errorf("failed to create i18n directory: %w", err)
	}

	// Generate i18n files for different languages
	languages := []string{"en", "zh-CN", "zh-TW"}
	for _, lang := range languages {
		i18nFile := filepath.Join(i18nDir, fmt.Sprintf("%s.json", lang))
		i18nData, err := generateI18nFile(config.Name, config.Description, operations, config.Categories, lang, config.EnableTranslation)
		if err != nil {
			return fmt.Errorf("failed to generate i18n file %s: %w", i18nFile, err)
		}
		if err := writeJSONFile(i18nFile, i18nData); err != nil {
			return fmt.Errorf("failed to write i18n file %s: %w", i18nFile, err)
		}
	}

	fmt.Printf("Output directory: %s\n", providerDir)
	fmt.Printf("Generated files:\n")
	fmt.Printf("  - %s\n", manifestPath)
	for _, lang := range languages {
		fmt.Printf("  - %s\n", filepath.Join(i18nDir, fmt.Sprintf("%s.json", lang)))
	}

	return nil
}
