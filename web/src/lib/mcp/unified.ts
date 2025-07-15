import type { BaseRemoteAPIResponse } from "@/typings/index.js"
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"
import { z } from "zod"
import { baseURL } from "@/config"
import { serverLogger } from "@/lib/utils"
import { serverRemoteFetch } from "@/lib/utils/fetch/server"
import packageJSON from "../../../package.json"

const mcpLogger = serverLogger.withTag("mcp-server")

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

  server.tool(
    "list_tools",
    "List all tools available in Context Space",
    async () => {
      try {
        const data = (await serverRemoteFetch(`/mcp/list_tools`, {
          method: "POST",
          headers: {
            Authorization: authorization,
          },
        })) as BaseRemoteAPIResponse<{
          tools: any[]
        }>
        if (data.success) {
          return { content: [{ type: "text", text: JSON.stringify(data.data.tools) }] }
        }
        return { content: [{ type: "text", text: JSON.stringify(data) }] }
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
    },
  )

  server.registerTool("call_tool", {
    description: "Call a tool, refer to the format of the list_tools tool",
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
    description: "Connect to an integration/provider, refer to the format of the list_tools tool",
    inputSchema: {
      provider_identifier: z.string(),
    },
  }, async (params) => {
    return {
      content: [{
        type: "text",
        text: JSON.stringify({
          success: true,
          message: `Open ${baseURL}/integration/${params.provider_identifier}`,
        }),
      }],
    }
  })

  server.server.onerror = error => mcpLogger.error("MCP server error", { error })

  return server
}
