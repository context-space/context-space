import type { MetadataRoute } from "next"
import { allBlogs } from "contentlayer/generated"
import { baseURL } from "@/config"
import { defaultLocale, locales } from "@/i18n/routing"
import { IntegrationService } from "@/services/integration"

// 定义静态页面路径
const staticPaths = [
  { path: "/integrations", priority: 0.9, changeFrequency: "daily" as const },
  { path: "/blogs", priority: 0.8, changeFrequency: "daily" as const },
  // { path: "/docs", priority: 0.8, changeFrequency: "weekly" as const },
  { path: "/privacy", priority: 0.3, changeFrequency: "monthly" as const },
  { path: "/terms", priority: 0.3, changeFrequency: "monthly" as const },
]

/**
 * 为给定路径生成所有语言版本的 URL
 * 根据 Next.js 国际化配置，默认语言可以不带语言前缀访问
 */
function generateLocalizedUrls(
  path: string,
  options: {
    priority: number
    changeFrequency: MetadataRoute.Sitemap[0]["changeFrequency"]
    lastModified?: Date
  },
): MetadataRoute.Sitemap {
  const { priority, changeFrequency, lastModified = new Date() } = options

  return locales.flatMap((locale) => {
    const urls: MetadataRoute.Sitemap = []

    if (locale === defaultLocale) {
      // 为默认语言生成两个 URL：不带语言前缀（主要）和带语言前缀（备用）
      urls.push({
        url: `${baseURL}${path}`,
        lastModified,
        changeFrequency,
        priority,
      })
      // 带语言前缀的版本优先级稍低
      urls.push({
        url: `${baseURL}/${locale}${path}`,
        lastModified,
        changeFrequency,
        priority: Math.max(priority - 0.1, 0.1),
      })
    } else {
      // 非默认语言只生成带语言前缀的 URL
      urls.push({
        url: `${baseURL}/${locale}${path}`,
        lastModified,
        changeFrequency,
        priority,
      })
    }

    return urls
  })
}

/**
 * 为动态页面生成多语言 URL
 */
function generateDynamicLocalizedUrls(
  pathTemplate: (locale: string) => string,
  options: {
    priority: number
    changeFrequency: MetadataRoute.Sitemap[0]["changeFrequency"]
    lastModified?: Date
  },
): MetadataRoute.Sitemap {
  const { priority, changeFrequency, lastModified = new Date() } = options

  return locales.flatMap((locale) => {
    const urls: MetadataRoute.Sitemap = []

    if (locale === defaultLocale) {
      // 为默认语言生成不带语言前缀的 URL
      const path = pathTemplate("")
      urls.push({
        url: `${baseURL}${path}`,
        lastModified,
        changeFrequency,
        priority,
      })
      // 也生成带语言前缀的版本
      const localizedPath = pathTemplate(`/${locale}`)
      urls.push({
        url: `${baseURL}${localizedPath}`,
        lastModified,
        changeFrequency,
        priority: Math.max(priority - 0.1, 0.1),
      })
    } else {
      // 非默认语言生成带语言前缀的 URL
      const localizedPath = pathTemplate(`/${locale}`)
      urls.push({
        url: `${baseURL}${localizedPath}`,
        lastModified,
        changeFrequency,
        priority,
      })
    }

    return urls
  })
}

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const now = new Date()

  // 生成根页面
  const rootPages: MetadataRoute.Sitemap = [
    {
      url: baseURL,
      lastModified: now,
      changeFrequency: "daily",
      priority: 1.0,
    },
    // 为默认语言生成带语言前缀的根页面
    {
      url: `${baseURL}/${defaultLocale}`,
      lastModified: now,
      changeFrequency: "daily",
      priority: 0.9,
    },
    // 为其他语言生成根页面
    ...locales
      .filter(locale => locale !== defaultLocale)
      .map(locale => ({
        url: `${baseURL}/${locale}`,
        lastModified: now,
        changeFrequency: "daily" as const,
        priority: 1.0,
      })),
  ]

  // 生成静态页面的多语言版本
  const staticPages = staticPaths.flatMap(({ path, priority, changeFrequency }) =>
    generateLocalizedUrls(path, { priority, changeFrequency, lastModified: now }),
  )

  // 生成博客页面的多语言版本
  const blogPages = allBlogs.flatMap(blog =>
    generateDynamicLocalizedUrls(
      localePrefix => `${localePrefix}/blogs/${blog.id}`,
      {
        priority: 0.7,
        changeFrequency: "weekly",
        lastModified: new Date(blog.publishedAt),
      },
    ),
  )

  // 生成集成页面的多语言版本
  const integrationService = new IntegrationService()
  const allIntegrations = await integrationService.getIntegrations()
  const integrationPages = allIntegrations.integrations.flatMap(integration =>
    generateDynamicLocalizedUrls(
      localePrefix => `${localePrefix}/integration/${integration.identifier}`,
      {
        priority: 0.7,
        changeFrequency: "weekly",
        lastModified: now,
      },
    ),
  )

  return [...rootPages, ...staticPages, ...blogPages, ...integrationPages]
}
