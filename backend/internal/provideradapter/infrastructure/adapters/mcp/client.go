package mcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCPClientConfig holds configuration for MCP client operations
type MCPClientConfig struct {
	Command string            `json:"command"` // Base command: "npx", "uvx", "./binary"
	Args    []string          `json:"args"`    // Complete argument list
	Envs    map[string]string `json:"envs"`    // Environment variables
	Timeout time.Duration     `json:"timeout"` // Timeout setting
}

// MCPClient wraps mcp-go client with adapter configuration
type MCPClient struct {
	config            MCPClientConfig
	adapterIdentifier string
	client            *client.Client
	initializeResult  *mcp.InitializeResult
}

// NewMCPGoClient creates a new mcp-go based client
func NewMCPGoClient(ctx context.Context, config MCPClientConfig, adapterIdentifier string) (*MCPClient, error) {
	client := &MCPClient{
		config:            config,
		adapterIdentifier: adapterIdentifier,
	}

	if err := client.connect(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

// Connect establishes connection to MCP server
func (c *MCPClient) connect(ctx context.Context) error {
	// Build command using the new configuration structure
	if c.config.Command == "" {
		return fmt.Errorf("command is required in MCP client config")
	}

	// Build environment variables
	envs := make([]string, 0, len(c.config.Envs))
	for k, v := range c.config.Envs {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	stdioTransport := transport.NewStdio(c.config.Command, envs, c.config.Args...)

	c.client = client.NewClient(stdioTransport)

	if err := c.client.Start(ctx); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	// Set up logging for stderr if available
	if stderr, ok := client.GetStderr(c.client); ok {
		go func() {
			buf := make([]byte, 4096)
			for {
				n, err := stderr.Read(buf)
				if err != nil {
					// Don't log EOF or "file already closed" errors as they are expected
					// when the MCP server process terminates
					if err != io.EOF && !errors.Is(err, os.ErrClosed) && !strings.Contains(err.Error(), "file already closed") {
						log.Printf("Error reading stderr: %v", err)
					}
					return
				}
				if n > 0 {
					fmt.Fprintf(os.Stderr, "[MCP Server (%s)] %s", c.adapterIdentifier, buf[:n])
				}
			}
		}()
	}

	initRequest := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			Capabilities:    mcp.ClientCapabilities{},
			ClientInfo: mcp.Implementation{
				Name:    fmt.Sprintf("Context Space MCP Client (%s)", c.adapterIdentifier),
				Version: "1.0.0",
			},
		},
	}

	initializeResult, err := c.client.Initialize(ctx, initRequest)
	if err != nil {
		return fmt.Errorf("failed to initialize MCP server: %w", err)
	}

	c.initializeResult = initializeResult

	return nil
}

// ListTools retrieves the list of tools from the MCP server
func (c *MCPClient) ListTools(ctx context.Context) (*mcp.ListToolsResult, error) {
	result, err := c.client.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	return result, nil
}

// CallTool calls a specific tool with the given arguments
func (c *MCPClient) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	callToolRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: arguments,
		},
	}
	result, err := c.client.CallTool(ctx, callToolRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool %s: %w", toolName, err)
	}

	return result, nil
}

// Close closes the client connection and terminates the server process
func (c *MCPClient) Close() error {
	return c.client.Close()
}

// GetMCPServerTools retrieves the list of tools from an MCP server using mcp-go
func GetMCPServerTools(ctx context.Context, config MCPClientConfig, adapterIdentifier string) ([]MCPTool, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	// Create and connect client
	mcpClient, err := NewMCPGoClient(ctx, config, adapterIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer mcpClient.Close()

	// Get tools using mcp-go
	toolsResult, err := mcpClient.ListTools(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	// Convert mcp.Tool to MCPTool
	mcpTools := make([]MCPTool, len(toolsResult.Tools))
	for i, tool := range toolsResult.Tools {
		mcpTools[i] = MCPTool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: convertInputSchema(tool.InputSchema),
		}
	}

	return mcpTools, nil
}

// CallMCPOperation calls a specific operation on an MCP server using mcp-go
func CallMCPOperation(ctx context.Context, config MCPClientConfig, operationName string, parameters map[string]interface{}, adapterIdentifier string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	// Create and connect client
	mcpClient, err := NewMCPGoClient(ctx, config, adapterIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer mcpClient.Close()

	// Call tool using mcp-go
	result, err := mcpClient.CallTool(ctx, operationName, parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool: %w", err)
	}

	// Convert result to our expected format
	operationResult := map[string]interface{}{
		"operation": operationName,
		"success":   !result.IsError,
		"content":   convertContent(result.Content),
	}

	if result.IsError {
		operationResult["error"] = true
	}

	return operationResult, nil
}

// convertInputSchema converts mcp.ToolInputSchema to MCPInputSchema
func convertInputSchema(schema mcp.ToolInputSchema) MCPInputSchema {
	properties := make(map[string]MCPProperty)

	// Convert properties from mcp schema
	for key, value := range schema.Properties {
		if propMap, ok := value.(map[string]interface{}); ok {
			properties[key] = convertProperty(propMap)
		}
	}

	return MCPInputSchema{
		Type:       schema.Type,
		Properties: properties,
		Required:   schema.Required,
	}
}

// convertProperty converts a property map to MCPProperty
func convertProperty(propMap map[string]interface{}) MCPProperty {
	prop := MCPProperty{}

	if typ, ok := propMap["type"].(string); ok {
		prop.Type = typ
	}

	if desc, ok := propMap["description"].(string); ok {
		prop.Description = desc
	}

	if enum, ok := propMap["enum"].([]interface{}); ok {
		prop.Enum = enum
	}

	if def, ok := propMap["default"]; ok {
		prop.Default = def
	}

	if items, ok := propMap["items"].(map[string]interface{}); ok {
		itemsProp := convertProperty(items)
		prop.Items = &itemsProp
	}

	return prop
}

// convertContent converts mcp.Content slice to our MCPContent format
func convertContent(content []mcp.Content) []MCPContent {
	result := make([]MCPContent, len(content))

	for i, c := range content {
		// Try to cast to TextContent (most common case)
		if textContent, ok := mcp.AsTextContent(c); ok {
			result[i] = MCPContent{
				Type: "text",
				Text: textContent.Text,
			}
		} else {
			// Handle other content types if needed
			result[i] = MCPContent{
				Type: "unknown",
				Text: fmt.Sprintf("%v", c),
			}
		}
	}

	return result
}
