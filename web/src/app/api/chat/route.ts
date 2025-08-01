import type { Message } from "ai"
import type { NextRequest } from "next/server"

import { createOpenAI } from "@ai-sdk/openai"
import { StreamableHTTPClientTransport } from "@modelcontextprotocol/sdk/client/streamableHttp.js"
import { experimental_createMCPClient, streamText } from "ai"
import { NextResponse } from "next/server"

import { baseURL } from "@/config"
import { getServerToken } from "@/lib/supabase/server"
import { serverLogger } from "@/lib/utils"

const chatApiLogger = serverLogger.withTag("chat-api")
export async function POST(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const provider = searchParams.get("id")
    const body = await request.json()
    const messages: Message[] = body.messages || []

    if (!messages.length) {
      return NextResponse.json(
        {
          type: "error",
          error: "No messages provided",
        },
        { status: 400 },
      )
    }

    const openai = createOpenAI({
      apiKey: process.env.OPENAI_API_KEY,
      baseURL: process.env.OPENAI_BASE_URL,
    })

    const tools = await (async function () {
      if (!provider) {
        return []
      }

      const accessToken = await getServerToken()

      const transport = new StreamableHTTPClientTransport(
        new URL(`/api/mcp/${provider}`, baseURL),
        {
          requestInit: {
            headers: {
              Authorization: `Bearer ${accessToken}`,
            },
          },
        },
      )

      const customClient = await experimental_createMCPClient({ transport })
      const tools = await customClient.tools()

      return tools
    })()

    const response = streamText({
      model: openai("gpt-4.1"),
      messages,
      tools,
      toolCallStreaming: true,
      maxSteps: 5,
    } as any)

    return response.toDataStreamResponse()
  } catch (error) {
    chatApiLogger.error("Error in chat API", { error })
    return NextResponse.json({ error: "Internal server error" }, { status: 500 })
  }
}
