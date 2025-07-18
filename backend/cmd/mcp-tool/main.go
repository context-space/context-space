package main

import (
	"fmt"
	"os"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate":
		runGenerateCommand(os.Args[2:])
	case "call":
		runCallCommand(os.Args[2:])
	case "serve":
		runServeCommand(os.Args[2:])
	case "version", "-v", "--version":
		fmt.Printf("mcp-tool version %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`mcp-tool - Model Context Protocol Tool v%s

Usage:
  mcp-tool <command> [options]

Available Commands:
  generate    Generate provider adapter from MCP server
  call        Call operations on MCP server
  serve       Start web UI server for mcp-tool commands
  version     Show version information
  help        Show this help message

Examples:
  # Generate adapter from filesystem server (Node.js)
  mcp-tool generate -command npx \
                    -arg "-y" -arg "@modelcontextprotocol/server-filesystem" -arg "/tmp" \
                    -identifier filesystem \
                    -name "Filesystem Provider" \
                    -description "Access to filesystem operations" \
                    -auth none -categories "storage,files"

  # Generate adapter for Brave Search
  mcp-tool generate -command npx \
                    -arg "-y" -arg "@modelcontextprotocol/server-brave-search" \
                    -identifier "brave_search" \
                    -name "Brave Search" \
                    -description "Web search powered by Brave" \
                    -auth apikey -categories "search,web"

  # Generate adapter from Python MCP server
  mcp-tool generate -command uvx \
                    -arg "python-mcp-server" -arg "--port" -arg "8080" \
                    -env "API_KEY=secret123" \
                    -identifier python_mcp \
                    -name "Python Provider" \
                    -description "Python-based MCP integration" \
                    -auth apikey -categories "python,tools"

  # Call operation on filesystem server  
  mcp-tool call -command npx \
                -arg "-y" -arg "@modelcontextprotocol/server-filesystem" -arg "/tmp" \
                -operation list_allowed_directories -params "{}"

  # Call operation on binary MCP server
  mcp-tool call -command ./custom-mcp-server \
                -arg "--config" -arg "/path/to/config.json" \
                -operation get_data -params '{"query": "test"}'

  # Start web UI server
  mcp-tool serve

  # Start web UI server on custom port
  mcp-tool serve -port 3000 -host 0.0.0.0

Use "mcp-tool <command> -h" for more information about a command.
`, version)
}

// Common server configuration structure
type MCPToolConfig struct {
	Command    string
	Args       []string
	Envs       map[string]string
	Operation  string
	Parameters map[string]interface{}
}
