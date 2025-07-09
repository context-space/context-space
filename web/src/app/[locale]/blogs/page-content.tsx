import type { Blog } from "contentlayer/generated"
import { getTranslations } from "next-intl/server"
import { FeaturedBlogs, RecentBlogs } from "@/components/blogs"

interface BlogsPageContentProps {
  blogs: Blog[]
}

export async function BlogsPageContent({ blogs }: BlogsPageContentProps) {
  const t = await getTranslations()
  const featuredBlogs = blogs.filter(blog => blog.featured)
  const recentBlogs = blogs.filter(blog => !blog.featured).sort((a, b) => new Date(b.publishedAt).getTime() - new Date(a.publishedAt).getTime())

  return (
    <div className="max-w-7xl mx-auto w-full">
      <div className="py-16">
        <div className="flex flex-col items-center gap-8">
          <div className="flex flex-col items-center gap-3 relative">
            <h1 className="text-4xl font-bold tracking-tight text-neutral-900 dark:text-white sm:text-5xl">
              {t("blogs.pageTitle")}
            </h1>
            <p className="text-neutral-600 dark:text-gray-400 text-lg text-center">
              {t("blogs.pageSubtitle")}
            </p>
          </div>
        </div>
      </div>

      <main className="pb-16">
        {featuredBlogs.length > 0 && <FeaturedBlogs blogs={featuredBlogs} />}
        <RecentBlogs blogs={recentBlogs} />
      </main>
    </div>
  )
}
