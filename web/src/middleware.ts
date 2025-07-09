import type { NextRequest } from "next/server"
import type { Locale } from "./i18n/routing"

import createMiddleware from "next-intl/middleware"
import { NextResponse } from "next/server"
import { serverLogger } from "@/lib/utils"
import { locales, routing } from "./i18n/routing"
import { updateSession } from "./lib/supabase/middleware"

const middlewareLogger = serverLogger.withTag("middleware")

/**
 * page route 由于需要获取用户信息，尽管部分页面不需要认证，也需要在middleware中处理
 * api route 如果不需要认证，也就不需要通过 middleware 处理
 */

interface RouteConfig {
  protectedRoutes: string[]
  publicAPIRoutes: string[]
}

// Route configuration
const routeConfig: RouteConfig = {
  // Page routes that require authentication
  protectedRoutes: ["/account"],
  // API routes that do not require authentication
  publicAPIRoutes: ["/api/auth/**", "/api/search", "/api/mcp/**", "/api/og"],
}

/**
 * 检查路由是否匹配给定的模式
 * 支持通配符 * 和 **
 */
function isRouteMatch(pathname: string, routes: string[]): boolean {
  return routes.some((route) => {
    if (route === pathname) {
      return true
    }

    if (route.endsWith("**")) {
      const prefix = route.slice(0, -2)

      return pathname.startsWith(prefix)
    }

    if (route.endsWith("*")) {
      const prefix = route.slice(0, -1)

      return pathname.startsWith(prefix) && !pathname.slice(prefix.length).includes("/")
    }

    return false
  })
}

/**
 * 处理未认证的请求
 */
function handleUnauthenticated(request: NextRequest, pathname: string): NextResponse {
  if (pathname.startsWith("/api/")) {
    middlewareLogger.warn("Unauthorized API access", { pathname })

    return NextResponse.json(
      {
        error: "Unauthorized",
        message: "Authentication required",
      },
      { status: 401 },
    )
  } else {
    const redirectUrl = encodeURIComponent(pathname)
    middlewareLogger.warn("Redirecting unauthenticated user to login", { pathname })
    const loginUrl = new URL(`/login?from=${redirectUrl}`, request.url)

    return NextResponse.redirect(loginUrl)
  }
}

/**
 * 获取不带语言前缀的路径
 */
function getPathnameWithoutLocale(pathname: string): string {
  // 移除语言前缀 (如 /zh, /en)
  const segments = pathname.split("/").filter(Boolean)
  if (segments.length > 0 && locales.includes(segments[0] as Locale)) {
    return `/${segments.slice(1).join("/")}`
  }
  return pathname
}

const handleI18nRouting = createMiddleware(routing)

/**
 * 主中间件函数
 */
export async function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname
  const pathnameWithoutLocale = getPathnameWithoutLocale(pathname)

  // API路由跳过国际化处理，直接处理认证
  if (pathname.startsWith("/api/")) {
    const { response: authResponse, user } = await updateSession(request, NextResponse.next())

    // 检查是否需要认证
    const isPublicAPIRoute = isRouteMatch(pathnameWithoutLocale, routeConfig.publicAPIRoutes)

    if (!isPublicAPIRoute && !user) {
      return handleUnauthenticated(request, pathnameWithoutLocale)
    }

    return authResponse
  }

  const i18nResponse = handleI18nRouting(request)

  // 然后处理认证
  const { response, user } = await updateSession(request, i18nResponse)

  const isProtectedRoute = isRouteMatch(pathnameWithoutLocale, routeConfig.protectedRoutes)

  if (isProtectedRoute && !user) {
    return handleUnauthenticated(request, pathnameWithoutLocale)
  }

  return response
}

export const config = {
  matcher: [
    // 匹配所有路径，排除静态文件和图片
    "/((?!_next/static|_next/image|favicon.ico|robots.txt|sitemap.xml|manifest.json|.*\\.(?:png|svg|jpg|jpeg|gif|webp|ico)$).*)",
  ],
}
