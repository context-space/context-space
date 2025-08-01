import type { NextRequest } from "next/server"
import { NextResponse } from "next/server"
import { serverLogger } from "@/lib/utils"
import { serverRemoteFetch } from "@/lib/utils/fetch/server"
import { getServerToken } from "./supabase/server"

export interface ProxyRouteParams {
  params: Promise<{
    path: string[]
  }>
}

export interface ProxyOptions {
  basePathPrefix: string
  loggerTag: string
}

export class APIProxy {
  private logger: any
  private basePathPrefix: string

  constructor(options: ProxyOptions) {
    this.logger = serverLogger.withTag(options.loggerTag)
    this.basePathPrefix = options.basePathPrefix
  }

  async handleRequest(request: NextRequest, { params }: ProxyRouteParams) {
    try {
      const { path } = await params
      const pathString = path.join("/")

      const url = `${this.basePathPrefix}/${pathString}`

      const searchParams = request.nextUrl.searchParams
      const queryString = searchParams.toString()
      const fullUrl = queryString ? `${url}?${queryString}` : url

      const options: any = {
        method: request.method,
        headers: {
          "Authorization": `Bearer ${await getServerToken()}`,
          "Content-Type": "application/json",
        },
      }

      if (request.method !== "GET" && request.method !== "HEAD") {
        const body = await request.text()
        if (body) {
          options.body = body
        }
      }

      this.logger.info("Forwarding API request", {
        method: request.method,
        url: fullUrl,
      })

      // 转发请求到远程API
      const response = await serverRemoteFetch(fullUrl, options)

      return NextResponse.json(response)
    } catch (error) {
      this.logger.error("Error forwarding API request", { error })
      return NextResponse.json(
        { error: "Internal server error", message: String(error) },
        { status: 500 },
      )
    }
  }
}

export function createProxyHandlers(options: ProxyOptions) {
  const proxy = new APIProxy(options)

  const handler = (request: NextRequest, context: ProxyRouteParams) =>
    proxy.handleRequest(request, context)

  return {
    GET: handler,
    POST: handler,
    PUT: handler,
    DELETE: handler,
    PATCH: handler,
  }
}