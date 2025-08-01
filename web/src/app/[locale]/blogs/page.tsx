import type { Metadata } from "next"
import { getTranslations } from "next-intl/server"
import { BaseLayout } from "@/components/layouts"
import { Breadcrumbs } from "@/components/seo/breadcrumbs"
import { StructuredData } from "@/components/seo/structured-data"
import { baseURL } from "@/config"
import { getAllBlogs } from "@/lib/blog"
import { generateBlogListStructuredData, generateBreadcrumbStructuredData } from "@/lib/seo/structured-data"
import { BlogsPageContent } from "./page-content"

export async function generateMetadata(): Promise<Metadata> {
  const canonicalUrl = `${baseURL}/blogs`

  return {
    title: "Blog | Context Space",
    description: "Stay updated with the latest news, tutorials, and insights from Context Space. Learn about AI engineering, context management, and the future of intelligent systems.",
    keywords: "blog, AI engineering, context engineering, tutorials, insights, news, Context Space, artificial intelligence, machine learning, automation",
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
      title: "Blog | Context Space",
      description: "Stay updated with the latest news, tutorials, and insights from Context Space.",
      type: "website",
      url: canonicalUrl,
      siteName: "Context Space",
      images: [
        {
          url: `/api/og?title=Blog&description=${encodeURIComponent("Latest updates and insights from Context Space")}`,
          width: 1200,
          height: 630,
          alt: "Context Space Blog",
        },
      ],
    },
    twitter: {
      card: "summary_large_image",
      title: "Blog | Context Space",
      description: "Stay updated with the latest news, tutorials, and insights from Context Space.",
      images: [`/api/og?title=Blog&description=${encodeURIComponent("Latest updates and insights from Context Space")}`],
      creator: "@contextspace",
    },
  }
}

export default async function BlogsPage() {
  const t = await getTranslations()
  const blogs = getAllBlogs()

  // Generate structured data for blog list
  const blogListStructuredData = generateBlogListStructuredData(
    blogs.map(blog => ({
      title: blog.title,
      description: blog.description,
      author: blog.author,
      publishedAt: blog.publishedAt,
      url: blog.url,
      image: blog.image,
    })),
  )

  const breadcrumbItems = [
    { name: t("nav.home"), url: `/` },
    { name: t("blogs.title"), url: `/blogs` },
  ]

  const breadcrumbStructuredData = generateBreadcrumbStructuredData(breadcrumbItems)

  return (
    <BaseLayout>
      <StructuredData data={blogListStructuredData} />
      <StructuredData data={breadcrumbStructuredData} />

      <div className="py-4">
        <Breadcrumbs
          items={breadcrumbItems}
          className="text-sm text-muted-foreground"
        />
      </div>
      <BlogsPageContent blogs={blogs} />
    </BaseLayout>
  )
}
