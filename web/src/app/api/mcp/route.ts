import type { NextRequest } from "next/server"
import { StreamableHTTPServerTransport } from "@modelcontextprotocol/sdk/server/streamableHttp.js"

import { toFetchResponse, toReqRes } from "fetch-to-node"
import { NextResponse } from "next/server"
import { getServer } from "@/lib/mcp/unified"
import { serverLogger } from "@/lib/utils"

const mcpApiLogger = serverLogger.withTag("mcp-api")
interface RouteParams {
  params: Promise<{
    id: string
  }>
}

async function handleMCPRequest(request: NextRequest, { params }: RouteParams) {
  try {
    const authorization = request.headers.get("Authorization")

    const { req, res } = toReqRes(request)
    const server = await getServer(authorization ?? "")
    const transport = new StreamableHTTPServerTransport({
      sessionIdGenerator: undefined,
    })
    transport.onerror = error => mcpApiLogger.error("MCP transport error", { error })
    await server.connect(transport)
    const requestBody = request.method === "POST" ? await request.json() : null
    await transport.handleRequest(req, res, requestBody)

    return toFetchResponse(res)
  } catch (error) {
    mcpApiLogger.error("Error in MCP API", { error })
    return NextResponse.json({ error: "Internal server error" }, { status: 500 })
  }
}

export async function GET(request: NextRequest, context: RouteParams) {
  return handleMCPRequest(request, context)
}

export async function POST(request: NextRequest, context: RouteParams) {
  return handleMCPRequest(request, context)
}
