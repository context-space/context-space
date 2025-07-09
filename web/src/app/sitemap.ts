import type { MetadataRoute } from "next"
import { allBlogs } from "contentlayer/generated"
import { baseURL } from "@/config"

const locales = ["en", "zh", "zh-TW"] as const

export default function sitemap(): MetadataRoute.Sitemap {
  const now = new Date()

  const staticPages: MetadataRoute.Sitemap = [
    {
      url: baseURL,
      lastModified: now,
      changeFrequency: "daily",
      priority: 1.0,
    },
    {
      url: `${baseURL}/en`,
      lastModified: now,
      changeFrequency: "daily",
      priority: 1.0,
    },
    {
      url: `${baseURL}/zh`,
      lastModified: now,
      changeFrequency: "daily",
      priority: 1.0,
    },
    {
      url: `${baseURL}/zh-TW`,
      lastModified: now,
      changeFrequency: "daily",
      priority: 1.0,
    },
    {
      url: `${baseURL}/en/integrations`,
      lastModified: now,
      changeFrequency: "daily",
      priority: 0.9,
    },
    {
      url: `${baseURL}/zh/integrations`,
      lastModified: now,
      changeFrequency: "daily",
      priority: 0.9,
    },
    {
      url: `${baseURL}/zh-TW/integrations`,
      lastModified: now,
      changeFrequency: "daily",
      priority: 0.9,
    },
    // Blog list pages
    {
      url: `${baseURL}/en/blogs`,
      lastModified: now,
      changeFrequency: "daily",
      priority: 0.8,
    },
    {
      url: `${baseURL}/zh/blogs`,
      lastModified: now,
      changeFrequency: "daily",
      priority: 0.8,
    },
    {
      url: `${baseURL}/zh-TW/blogs`,
      lastModified: now,
      changeFrequency: "daily",
      priority: 0.8,
    },
    {
      url: `${baseURL}/en/documents`,
      lastModified: now,
      changeFrequency: "weekly",
      priority: 0.8,
    },
    {
      url: `${baseURL}/zh/documents`,
      lastModified: now,
      changeFrequency: "weekly",
      priority: 0.8,
    },
    {
      url: `${baseURL}/zh-TW/documents`,
      lastModified: now,
      changeFrequency: "weekly",
      priority: 0.8,
    },
    {
      url: `${baseURL}/en/privacy`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
    {
      url: `${baseURL}/zh/privacy`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
    {
      url: `${baseURL}/zh-TW/privacy`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
    {
      url: `${baseURL}/en/terms`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
    {
      url: `${baseURL}/zh/terms`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
    {
      url: `${baseURL}/zh-TW/terms`,
      lastModified: now,
      changeFrequency: "monthly",
      priority: 0.3,
    },
  ]

  // Generate blog pages for all locales
  const blogPages: MetadataRoute.Sitemap = allBlogs.flatMap(blog =>
    locales.map(locale => ({
      url: `${baseURL}/${locale}/blogs/${blog.id}`,
      lastModified: new Date(blog.publishedAt),
      changeFrequency: "weekly" as const,
      priority: 0.7,
    })),
  )

  return [...staticPages, ...blogPages]
}
