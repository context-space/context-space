import type { MetadataRoute } from "next"
import { baseURL } from "@/config"

export default function robots(): MetadataRoute.Robots {
  return {
    rules: [
      {
        userAgent: "*",
        allow: [
          "/",
          "/en/",
          "/zh/",
          "/zh-TW/",
          "/integrations",
          "/documents",
          "/privacy",
          "/terms",
          "/api/og",
        ],
        disallow: [
          "/private/",
          "/admin/",
          "/api/",
          "/logout",
          "/login",
          "/account",
          "/oauth-callback",
          "/mcp-oauth-callback",
        ],
      },
      {
        userAgent: "GPTBot",
        disallow: "/",
      },
      {
        userAgent: "ChatGPT-User",
        disallow: "/",
      },
      {
        userAgent: "CCBot",
        disallow: "/",
      },
      {
        userAgent: "anthropic-ai",
        disallow: "/",
      },
    ],
    sitemap: `${baseURL}/sitemap.xml`,
    host: baseURL,
  }
}
