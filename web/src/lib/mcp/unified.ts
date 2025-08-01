import type { BaseRemoteAPIResponse, Credential, Provider } from "@/typings"
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"
import { z } from "zod"
import { baseURL } from "@/config"
import { serverLogger } from "@/lib/utils"
import { remoteFetchWithErrorHandling, serverRemoteFetch } from "@/lib/utils/fetch/server"
import packageJSON from "../../../package.json"

const mcpLogger = serverLogger.withTag("mcp-server")

async function getConnectionStatusData(authorization: string) {
  const { credentials } = await remoteFetchWithErrorHandling<{ credentials: Credential[] }>(`/credentials`, {
    headers: {
      Authorization: authorization,
    },
  })

  const { providers: providersWithoutConnected } = await remoteFetchWithErrorHandling<{ providers: Provider[] }>("/providers/filter", {
    method: "POST",
    headers: {
      Authorization: authorization,
    },
    body: JSON.stringify({
      filters: {
        auth_type: "none",
      },
      pagination: {
        page: 1,
        page_size: 100,
      },
    }),
  })

  return { credentials, providersWithoutConnected }
}

function getProviderConnectionStatus(
  providerIdentifier: string,
  credentials: Credential[],
  providersWithoutConnected: Provider[],
): "ready_to_use" | "not_connected" {
  const is_connection_free = providersWithoutConnected.some(provider => provider.identifier === providerIdentifier)
  const is_connected = credentials.some(cred => cred.provider_identifier === providerIdentifier)

  if (is_connection_free || is_connected) {
    return "ready_to_use"
  } else {
    return "not_connected"
  }
}

export async function getServer(authorization: string) {
  const server = new McpServer(
    {
      name: "Context Space",
      version: packageJSON.version,
    },
    {
      capabilities: { logging: {} },
    },
  )

  server.registerTool("get_connection_status", {
    description: `Get the connection status of an integration if you don't know the status when you call the tool.
    if it needs connect and not connected, return "not_connected" and you should call the connect_to_integration tool;
    else if it is ready to use, return "ready_to_use"`,
    inputSchema: {
      provider_identifier: z.string().describe("The provider identifier, refer to the list_tools tool, empty means all providers"),
    },
  }, async (params) => {
    try {
      const { credentials, providersWithoutConnected } = await getConnectionStatusData(authorization)

      if (params.provider_identifier) {
        const status = getProviderConnectionStatus(params.provider_identifier, credentials, providersWithoutConnected)
        return {
          content: [{
            type: "text",
            text: status,
          }],
        }
      }

      const { tools } = await remoteFetchWithErrorHandling<{
        tools: Record<string, any>
      }>(`/mcp/list_tools`, {
        method: "POST",
        headers: {
          Authorization: authorization,
        },
      })

      const allProviderIdentifiers = Object.keys(tools)
      const connectionFreeProviders = providersWithoutConnected.map(provider => provider.identifier)
      const connectedProviders = credentials.map(cred => cred.provider_identifier)
      const readyToUseProviders = [...connectionFreeProviders, ...connectedProviders]
      const notConnectedProviders = allProviderIdentifiers.filter(
        identifier => !readyToUseProviders.includes(identifier),
      )

      return {
        content: [{
          type: "text",
          text: JSON.stringify({
            ready_to_use: readyToUseProviders,
            not_connected: notConnectedProviders,
          }),
        }],
      }
    } catch (error) {
      mcpLogger.error("MCP server error", { error })
      return {
        content: [{
          type: "text",
          text: error instanceof Error ? error.message : "Unknown error",
        }],
      }
    }
  })

  server.registerTool("list_tools", {
    description: "List all tools available in Context Space, if you don't find the tool you want, you could filter connection_status=not_connected to get all not connected tools",
    inputSchema: {
      connection_status: z.enum(["ready_to_use", "not_connected"]).optional().describe("The connection status filter. If not provided, return all ready to use tools"),
    },
  }, async (params) => {
    try {
      const { tools } = await remoteFetchWithErrorHandling<{
        tools: Record<string, any>
      }>(`/mcp/list_tools`, {
        method: "POST",
        headers: {
          Authorization: authorization,
        },
      })

      const { credentials, providersWithoutConnected } = await getConnectionStatusData(authorization)

      const categorizedTools = {
        ready_to_use: {} as Record<string, any>,
        not_connected: {} as Record<string, any>,
      }

      Object.keys(tools).forEach((providerKey) => {
        const connection_status = getProviderConnectionStatus(providerKey, credentials, providersWithoutConnected)
        categorizedTools[connection_status][providerKey] = tools[providerKey]
      })

      if (params?.connection_status) {
        return {
          content: [{
            type: "text",
            text: JSON.stringify({
              tools: categorizedTools[params.connection_status],
            }),
          }],
        }
      }

      return {
        content: [{
          type: "text",
          text: JSON.stringify({
            tools: categorizedTools.ready_to_use,
          }),
        }],
      }
    } catch (error) {
      mcpLogger.error("MCP server error", { error })
      return {
        content: [{
          type: "text",
          text: JSON.stringify({
            success: false,
            message: error instanceof Error ? error.message : "Unknown error",
          }),
        }],
      }
    }
  })

  server.registerTool("call_tool", {
    description: "Call a tool, refer to the format of the list_tools tool. If you are unsure whether this tool is connected, you should first call the get_connection_status tool to check if the connection is established.",
    inputSchema: {
      provider_identifier: z.string(),
      operation_identifier: z.string(),
      parameters: z.string().describe("The parameters to call the tool with, in JSON format"),
    },
  }, async (params) => {
    try {
      if (typeof params.parameters === "string" && params.parameters.startsWith("{")) {
        params.parameters = JSON.parse(params.parameters)
      }
      const res = (await serverRemoteFetch(`/mcp/call_tool/${params.provider_identifier}/${params.operation_identifier}`, {
        headers: {
          Authorization: authorization,
        },
        timeout: 60000,
        method: "POST",
        body: JSON.stringify(params.parameters),
      })) as BaseRemoteAPIResponse<{
        error?: string
        tool_result: any
      }>
      if (res.success) {
        return { content: [{ type: "text", text: JSON.stringify(res.data.tool_result) }] }
      } else if (res?.data?.error) {
        return {
          content: [{
            type: "text",
            text: JSON.stringify({
              success: false,
              message: res.data.error,
            }),
          }],
        }
      }
      return { content: [{ type: "text", text: JSON.stringify(res) }] }
    } catch (error) {
      mcpLogger.error("MCP server error", { error })
      return {
        content: [{
          type: "text",
          text: JSON.stringify({
            success: false,
            message: error instanceof Error ? error.message : "Unknown error",
          }),
        }],
      }
    }
  })

  server.registerTool("connect_to_integration", {
    description: "Connect to an integration if it is not connected",
    inputSchema: {
      provider_identifier: z.string(),
    },
  }, async (params) => {
    return {
      content: [{
        type: "text",
        text: JSON.stringify({
          success: true,
          message: `Tell the user to open the integration page in the browser: [click here](${baseURL}/integration/${params.provider_identifier})`,
        }),
      }],
    }
  })

  server.server.onerror = error => mcpLogger.error("MCP server error", { error })

  return server
}
