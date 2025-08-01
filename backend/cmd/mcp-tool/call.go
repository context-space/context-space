package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/mcp"
)

type CallConfig struct {
	MCPToolConfig
}

func runCallCommand(args []string) {
	fs := flag.NewFlagSet("call", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf(`Usage: mcp-tool call [options]

Call operations on MCP server

Options:
`)
		fs.PrintDefaults()
		fmt.Printf(`
Examples:
  # Call operation with environment variables
  mcp-tool call -command npx \
                -arg "-y" -arg "@custom/mcp-server" \
                -env "API_KEY=secret123" -env "DEBUG=mcp:*" \
                -operation get_data \
                -params '{"query": "test"}'
`)
	}

	var config CallConfig
	var paramStr string
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

	fs.StringVar(&config.Operation, "operation", "", "Operation/tool name to call (required)")
	fs.StringVar(&paramStr, "params", "", "Operation parameters (JSON format)")

	if err := fs.Parse(args); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	// Validate required flags
	if config.Command == "" {
		fmt.Fprintf(os.Stderr, "Error: -command is required\n")
		fs.Usage()
		os.Exit(1)
	}

	if config.Operation == "" {
		fmt.Fprintf(os.Stderr, "Error: -operation is required\n")
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

	// Parse parameters
	if paramStr != "" {
		if err := json.Unmarshal([]byte(paramStr), &config.Parameters); err != nil {
			log.Fatalf("Failed to parse parameters JSON: %v", err)
		}
	}

	// Create MCP client configuration
	mcpConfig := mcp.MCPClientConfig{
		Command: config.Command,
		Args:    config.Args,
		Envs:    config.Envs,
		Timeout: 60 * time.Second,
	}

	// Call the operation
	adapterIdentifier := ""
	if mcpConfig.Command == "npx" {
		adapterIdentifier = mcpConfig.Args[1]
	} else {
		adapterIdentifier = mcpConfig.Args[0]
	}
	result, err := mcp.CallMCPOperation(context.Background(), mcpConfig, config.Operation, config.Parameters, adapterIdentifier)
	if err != nil {
		log.Fatalf("Failed to call MCP operation: %v", err)
	}

	// Output the result as JSON
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}

	fmt.Println(string(resultJSON))
}
