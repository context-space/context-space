import type { MetadataRoute } from "next"
import { baseURL } from "@/config"

export default function robots(): MetadataRoute.Robots {
  if (baseURL.includes("dev")) {
    return {
      rules: [
        {
          userAgent: "*",
          disallow: "/",
        },
      ],
    }
  }

  return {
    rules: [
      {
        userAgent: "*",
        allow: [
          "/",
          "/en/",
          "/zh/",
          "/zh-TW/",
        ],
        disallow: [
          "/private/",
          "/admin/",
          "/api/",
          "/logout",
          "/login",
          "/auth-callback",
          "/provider-callback",
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
