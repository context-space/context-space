# MCP Tool - Model Context Protocol Tool

A unified command-line tool for working with MCP (Model Context Protocol) servers. This tool provides two main functionalities: generating provider adapters from MCP servers and calling operations on MCP servers.

## Overview

The MCP Tool combines the functionality of the previous separate `mcp-adapter-generator` and `mcp-operation-caller` tools into a single, cohesive command-line interface. It leverages the shared MCP library for consistent protocol handling and configuration parsing.

## Installation

To build the tool:

```bash
go build -o ./build/mcp-tool cmd/mcp-tool/*.go
```

Or run directly:

```bash
go run cmd/mcp-tool/*.go <command> [options]
```

## Commands

### generate

Generate provider adapter configuration from MCP server tools.

**Usage:**
```bash
mcp-tool generate [options]
```

**Required Options:**
- `-command`: Base command to run MCP server (e.g., 'npx', 'uvx', './binary')
- `-identifier`: Unique identifier for the provider (used as directory name)
- `-name`: Display name for the provider
- `-description`: Description of the provider

**Optional Options:**
- `-output`: Output directory for provider configurations (default: "configs/providers")
- `-auth`: Authentication type - oauth, apikey, basic, none (default: "none")
- `-categories`: Comma-separated list of categories
- `-arg`: Command argument (can be specified multiple times)
- `-env`: Environment variable KEY=VALUE (can be specified multiple times)
- `-translate`: Enable AI translation for i18n files (requires OPENAI_API_KEY)

**Examples:**

```bash
# Generate adapter for Node.js filesystem server
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
                  -description "Web search capabilities powered by Brave Search API" \
                  -auth apikey -categories "search,web"

# Generate adapter with custom OpenAI endpoint
export OPENAI_API_KEY="your-api-key"
export OPENAI_BASE_URL="https://your-custom-openai-endpoint.com/v1"
mcp-tool generate -command npx \
                  -arg "-y" -arg "@modelcontextprotocol/server-github" \
                  -name github_provider -auth apikey \
                  -categories "development,github" \
                  -translate
```

**Environment Variables for AI Translation:**
- `OPENAI_API_KEY`: Required for AI translation (when using `-translate`)
- `OPENAI_BASE_URL`: Custom OpenAI API base URL (optional, defaults to https://api.openai.com/v1)
- `OPENAI_TIMEOUT`: Timeout for OpenAI API calls in seconds (optional, defaults to 60)

**Output:**
The generate command creates:
- `manifest.json`: Provider configuration with operations and parameters
- `i18n/en.json`: Base English internationalization file
- `i18n/zh-CN.json`, `i18n/zh-TW.json`: Localized files (AI-translated if `-translate` is enabled)

### serve

Start HTTP server with Web UI for executing mcp-tool commands through a browser interface.

**Usage:**
```bash
mcp-tool serve [options]
```

**Optional Options:**
- `-host`: Host to serve on (default: "localhost")
- `-port`: Port to serve on (default: "8080")

**Examples:**

```bash
# Start server on default port (localhost:8080)
mcp-tool serve

# Start server on custom port
mcp-tool serve -port 3000

# Start server accessible from all interfaces
mcp-tool serve -host 0.0.0.0 -port 8080
```

**Features:**
- **Visual Interface**: Easy-to-use web forms for all mcp-tool commands
- **Real-time Output**: Live streaming of command output and errors
- **Job Management**: Track running, completed, and failed operations
- **Form Validation**: Built-in validation for required fields
- **Multi-line Support**: Textarea inputs for arguments and environment variables
- **Status Indicators**: Visual feedback for command execution status

**Web UI Sections:**
1. **Generate Provider Adapter**: Visual form for the `generate` command with all parameters
2. **Call MCP Operation**: Interface for the `call` command with parameter input
3. **Real-time Console**: Live output streaming with WebSocket connection
4. **Job Status**: Current and historical command execution status

**Usage Workflow:**
1. Start the server: `mcp-tool serve`
2. Open browser to `http://localhost:8080`
3. Fill in the form fields for your desired command
4. Click "Execute" to run the command
5. Watch real-time output in the console section
6. View results and status updates

**Development Notes:**
- The server automatically manages command execution as background processes
- WebSocket connections provide real-time output streaming
- All form data is validated before command execution
- Commands are executed using the same mcp-tool binary as subprocesses

### call

Call operations on MCP servers and return results.

**Usage:**
```bash
mcp-tool call [options]
```

**Required Options:**
- `-package`: NPM package name for the MCP server
- `-operation`: Operation/tool name to call

**Optional Options:**
- `-params`: JSON string containing operation parameters
- `-params-file`: JSON file containing operation parameters
- `-args`: Additional arguments for the MCP server (shell-style quoted)
- `-env`: Environment variables (comma-separated KEY=VALUE pairs)

**Examples:**

```bash
# Call operation with Python MCP server
mcp-tool call -command uvx \
              -arg "python-mcp-server" -arg "--port" -arg "8080" \
              -operation search \
              -params '{"query": "test"}'

# Call operation with environment variables
mcp-tool call -command npx \
              -arg "-y" -arg "@custom/mcp-server" \
              -env "API_KEY=secret123" -env "DEBUG=mcp:*" \
              -operation get_data \
              -params '{"query": "test"}'
```

**Output:**
The call command returns structured JSON with operation results:
```json
{
  "operation": "list_directory",
  "success": true,
  "content": [
    {
      "type": "text",
      "text": "Directory listing results..."
    }
  ]
}
```

## Advanced Features

### AI Translation for i18n Files

The `generate` command supports automatic translation of internationalization files using OpenAI's language models:

**Features:**
- Automatically translates provider names, descriptions, and parameter descriptions
- Maintains JSON structure and key names
- Supports Chinese Simplified (zh-CN) and Traditional Chinese (zh-TW)
- Falls back to English if translation fails or API key is missing
- Uses professional technical terminology suitable for developers

**Setup:**
1. Set your OpenAI API key: `export OPENAI_API_KEY="your-api-key"`
2. Optionally configure custom endpoint: `export OPENAI_BASE_URL="your-endpoint"`
3. Add the `-translate` flag to your generate command

**Example:**
```bash
export OPENAI_API_KEY="sk-xxx"
mcp-tool generate -command npx \
                  -arg "-y" -arg "@modelcontextprotocol/server-filesystem" \
                  -name filesystem -auth none \
                  -translate
```

This will create:
- `i18n/en.json`: Original English content
- `i18n/zh-CN.json`: AI-translated Simplified Chinese
- `i18n/zh-TW.json`: AI-translated Traditional Chinese

## Error Handling

The tool provides comprehensive error handling:

- **Argument parsing errors**: Invalid quotes, escape sequences
- **Environment variable errors**: Invalid KEY=VALUE format
- **MCP server errors**: Connection failures, protocol errors
- **Parameter validation**: JSON parsing, required fields

**Common Error Scenarios:**

1. **Server requires arguments:**
   ```
   Error: MCP server requires arguments. Server usage: filesystem <directory>
   ```
   **Solution:** Add the `-args` flag with required arguments.

2. **Invalid JSON parameters:**
   ```
   Error: failed to parse parameters from JSON string: invalid character...
   ```
   **Solution:** Check JSON syntax in `-params` or `-params-file`.

3. **Environment variable parsing:**
   ```
   Error: invalid environment variable format: INVALID (expected KEY=VALUE)
   ```
   **Solution:** Use proper KEY=VALUE format in `-env`.

## Integration with Context Space Backend

This tool is designed for seamless integration with the Context Space backend system:

1. **Provider Generation**: Creates manifests compatible with the adapter system
2. **Operation Testing**: Validates MCP server operations before integration
3. **Development Workflow**: Supports rapid prototyping of MCP integrations

## Dependencies

- Go 1.19 or later
- Node.js and npm (for running MCP servers)
- Access to MCP server packages

## Technical Details

- **Timeout**: Default 60-second timeout for MCP operations (configurable)
- **Process Management**: Automatic cleanup of MCP server processes
- **Protocol Support**: JSON-RPC 2.0 with MCP extensions
- **Memory Efficient**: Streaming JSON processing

## Troubleshooting

### Common Issues

1. **NPM package not found:**
   - Ensure the package exists: `npm view @package/name`
   - Check package spelling and version

2. **Server fails to start:**
   - Check if additional arguments are needed
   - Verify environment variables are set correctly
   - Check Node.js version compatibility

3. **Operation not found:**
   - List available operations: `mcp-tool generate -package @package/name -name temp`
   - Check operation name spelling

4. **Permission issues:**
   - Ensure proper file system permissions for output directory
   - Check if MCP server requires specific permissions

### Debug Mode

Enable debug output by setting environment variables:

```bash
export DEBUG=mcp:*
mcp-tool call -package @package/name -operation test
```
