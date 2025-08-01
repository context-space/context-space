import type { BaseRemoteAPIResponse, Invocation, ProviderDetail } from "@/typings/index.js"
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"
import { z } from "zod"
import { serverLogger } from "@/lib/utils"
import { serverRemoteFetch } from "@/lib/utils/fetch/server"
import packageJSON from "../../../package.json"

const mcpLogger = serverLogger.withTag("mcp-server")

function paramsToZodShape(parameters: any[]) {
  const shape: Record<string, any> = {}
  if (!parameters) {
    return shape
  }
  parameters.forEach((param) => {
    let zodType
    switch (param.type) {
      case "string":
        zodType = z.string()
        break
      case "integer":
        zodType = z.number()
        break
      case "boolean":
        zodType = z.boolean()
        break
      default:
        zodType = z.any()
    }
    if (param.default) {
      // TODO
      // @ts-expect-error todo
      zodType = zodType.default(param.default)
    }
    if (!param.required) {
      zodType = zodType.optional()
    }

    let description = param.description
    // not support z.literal()
    if (param.enum?.length) {
      description = `${description}\n\nPossible values: ${param.enum.join(", ")}`
    }
    if (description) {
      zodType = zodType.describe(description)
    }

    shape[param.name] = zodType
  })

  return shape
}

export async function getServer(provider: string, authorization: string) {
  const server = new McpServer(
    {
      name: provider,
      version: packageJSON.version,
    },
    {
      capabilities: { logging: {} },
    },
  )

  try {
    const { data: { operations } } = (await serverRemoteFetch(`/providers/${provider}`)) as BaseRemoteAPIResponse<ProviderDetail>
    operations.forEach((operation) => {
      server.tool(
        operation.identifier,
        operation.description,
        paramsToZodShape(operation.parameters),
        async (parameters) => {
          try {
            const res = (await serverRemoteFetch(`/invocations/${provider}/${operation.identifier}`, {
              headers: {
                Authorization: authorization,
              },
              body: JSON.stringify({ parameters }),
              timeout: 60000,
              method: "POST",
            })) as BaseRemoteAPIResponse<Invocation>
            if (res.success) {
              return {
                content: [
                  {
                    type: "text",
                    text: JSON.stringify(res.data.response_data),
                  },
                ],
              }
            } else if (res?.data?.error) {
              return {
                content: [
                  {
                    type: "text",
                    text: JSON.stringify({
                      success: false,
                      message: res.data.error,
                    }),
                  },
                ],
              }
            }
            return {
              content: [
                {
                  type: "text",
                  text: JSON.stringify(res),
                },
              ],
            }
          } catch (error) {
            mcpLogger.error("MCP server error", { error })
            return {
              content: [
                {
                  type: "text",
                  text: JSON.stringify({
                    success: false,
                    message: error instanceof Error ? error.message : "Unknown error",
                  }),
                },
              ],
            }
          }
        },
      )
    })
  } catch (error) {
    mcpLogger.error("MCP server error", { error })
  }

  server.server.onerror = error => mcpLogger.error("MCP server error", { error })

  return server
}
