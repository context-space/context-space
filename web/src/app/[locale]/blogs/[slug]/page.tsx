import type { Metadata } from "next"
import type { Locale } from "@/i18n/routing"
import { allBlogs } from "contentlayer/generated"
import { getTranslations } from "next-intl/server"
import { notFound } from "next/navigation"
import { BaseLayout } from "@/components/layouts"
import { Breadcrumbs } from "@/components/seo/breadcrumbs"
import { StructuredData } from "@/components/seo/structured-data"
import { baseURL } from "@/config"
import { generateArticleStructuredData, generateBreadcrumbStructuredData } from "@/lib/seo/structured-data"
import { BlogContent } from "./blog-content"

interface BlogPageProps {
  params: Promise<{ locale: Locale, slug: string }>
}

// export async function generateStaticParams() {
//   const { locales } = await import("@/i18n/routing")

//   return locales.flatMap(locale =>
//     allBlogs.map(blog => ({
//       locale,
//       slug: blog.id,
//     })),
//   )
// }

// export const dynamic = "force-static"

export async function generateMetadata({ params }: BlogPageProps): Promise<Metadata> {
  const { slug } = await params
  const blog = allBlogs.find(blog => blog._raw.flattenedPath === slug)

  if (!blog) {
    return {
      title: "Blog Not Found",
    }
  }

  const canonicalUrl = `${baseURL}/blogs/${slug}`
  const ogImageUrl = `/api/og?title=${encodeURIComponent(blog.title)}&description=${encodeURIComponent(blog.description)}`

  // Generate keywords from title, description and category
  const keywords = [
    blog.category || "AI engineering",
    "context management",
    "artificial intelligence",
    "Context Space",
    "blog",
    "tutorial",
    blog.author.toLowerCase(),
    ...blog.title.toLowerCase().split(" ").filter(word => word.length > 3),
  ].join(", ")

  return {
    title: `${blog.title} | Context Space Blog`,
    description: blog.description,
    keywords,
    authors: [{ name: blog.author }],
    category: blog.category || "Technology",
    alternates: {
      canonical: canonicalUrl,
    },
    robots: {
      index: true,
      follow: true,
      googleBot: {
        "index": true,
        "follow": true,
        "max-image-preview": "large",
        "max-snippet": -1,
      },
    },
    openGraph: {
      title: `${blog.title} | Context Space Blog`,
      description: blog.description,
      type: "article",
      url: canonicalUrl,
      publishedTime: blog.publishedAt,
      authors: [blog.author],
      section: blog.category || "Technology",
      tags: [blog.category || "AI engineering", "context management", "artificial intelligence"],
      images: [
        {
          url: ogImageUrl,
          width: 1200,
          height: 630,
          alt: blog.title,
        },
      ],
      siteName: "Context Space",
    },
    twitter: {
      card: "summary_large_image",
      title: `${blog.title} | Context Space Blog`,
      description: blog.description,
      images: [ogImageUrl],
      creator: "@contextspace",
    },
    other: {
      "article:published_time": blog.publishedAt,
      "article:author": blog.author,
      "article:section": blog.category || "Technology",
    },
  }
}

export default async function BlogPage({ params }: BlogPageProps) {
  const { slug } = await params
  const t = await getTranslations()

  const blog = allBlogs.find(blog => blog.id === slug)

  if (!blog) {
    notFound()
  }

  const canonicalUrl = `${baseURL}/blogs/${slug}`

  // Generate structured data
  const articleStructuredData = generateArticleStructuredData({
    headline: blog.title,
    description: blog.description,
    author: blog.author,
    datePublished: blog.publishedAt,
    url: canonicalUrl,
    image: blog.image,
  })

  const breadcrumbItems = [
    { name: t("nav.home"), url: `/` },
    { name: t("blogs.title"), url: `/blogs` },
    { name: blog.title, url: `/blogs/${slug}` },
  ]

  const breadcrumbStructuredData = generateBreadcrumbStructuredData(breadcrumbItems)

  return (
    <BaseLayout>
      <StructuredData data={articleStructuredData} />
      <StructuredData data={breadcrumbStructuredData} />

      <div className="py-4">
        <Breadcrumbs
          items={breadcrumbItems}
          className="text-sm text-muted-foreground"
        />
      </div>
      <BlogContent blog={blog} />
    </BaseLayout>
  )
}
