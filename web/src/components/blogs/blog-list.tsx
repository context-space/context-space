"use client"

import type { Blog } from "contentlayer/generated"
import { CalendarIcon, UserIcon } from "lucide-react"
import { useTranslations } from "next-intl"
import Image from "next/image"
import { Link } from "@/i18n/navigation"
import { formatDate } from "@/lib/blog"
import { cn } from "@/lib/utils"
import { BlogCard } from "./blog-card"

interface BlogListProps {
  blogs: Blog[]
}

export function FeaturedBlogs({ blogs }: BlogListProps) {
  if (blogs.length === 0) return null
  const largeBlog = blogs.find(blog => blog.featured === 1)
  const mediumBlogs = blogs.filter(blog => blog.featured === 2 || blog.featured === 3)

  return (
    <>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-16">
        {/* Left column - large featured blog */}
        {largeBlog && (
          <div className="lg:col-span-2">
            <FeaturedBlogCard blog={largeBlog} size="large" />
          </div>
        )}

        {/* Right column - two smaller blogs stacked vertically */}
        <div className="lg:col-span-1 flex flex-col gap-6">
          {mediumBlogs.map(blog => (
            <FeaturedBlogCard key={blog.id} blog={blog} size="medium" />
          ))}
        </div>
      </div>
    </>
  )
}

// Enhanced blog card for featured section with size variants
function FeaturedBlogCard({ blog, size }: { blog: Blog, size: "large" | "medium" }) {
  const isLarge = size === "large"

  if (isLarge) {
    // Large overlay card for left column
    return (
      <Link
        href={blog.url}
        className="block h-full transition-transform hover:scale-[1.01] duration-300"
      >
        <article className={cn(
          "group relative overflow-hidden rounded-xl bg-white/[0.02]",
          "border border-primary/15 dark:border-white/10 hover:border-primary/30 dark:hover:border-primary/40",
          "transition-all duration-300 h-full min-h-[400px]",
        )}
        >
          {/* Blog Image */}
          <div className="absolute inset-0">
            <Image
              src={blog.image}
              alt={blog.title}
              fill
              className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            />
            {/* Gradient overlay */}
            <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent" />
          </div>

          {/* Content */}
          <div className="relative z-10 flex flex-col justify-end h-full p-6">
            <h2 className="text-2xl font-semibold text-white mb-3 group-hover:text-primary/90 transition-colors line-clamp-2">
              {blog.title}
            </h2>

            <p className="text-gray-200 mb-4 line-clamp-2">
              {blog.description}
            </p>

            <div className="flex items-center gap-4 text-xs text-gray-300">
              <div className="flex items-center gap-1">
                <CalendarIcon className="w-3 h-3" />
                <span>{formatDate(blog.publishedAt)}</span>
              </div>
              <div className="flex items-center gap-1">
                <UserIcon className="w-3 h-3" />
                <span>{blog.author}</span>
              </div>
            </div>
          </div>
        </article>
      </Link>
    )
  }

  // Small square cards for right column
  return (
    <Link
      href={blog.url}
      className="block h-full transition-transform hover:scale-[1.01] duration-300"
    >
      <article className={cn(
        "group relative flex flex-col overflow-hidden rounded-xl bg-white/[0.02]",
        "border border-primary/15 dark:border-white/10 hover:border-primary/30 dark:hover:border-primary/40",
        "transition-all duration-300 h-full",
      )}
      >
        {/* Blog Image */}
        <div className="relative w-full aspect-video overflow-hidden">
          <Image
            src={blog.image}
            alt={blog.title}
            fill
            className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
        </div>

        {/* Content */}
        <div className="flex flex-col flex-grow p-4">
          <h3 className="text-base font-semibold text-neutral-900 dark:text-white mb-2 group-hover:text-primary transition-colors line-clamp-2">
            {blog.title}
          </h3>

          <div className="flex items-center gap-3 text-xs text-neutral-500 dark:text-gray-500 mt-auto">
            <div className="flex items-center gap-1">
              <CalendarIcon className="w-3 h-3" />
              <span>{formatDate(blog.publishedAt)}</span>
            </div>
          </div>
        </div>
      </article>
    </Link>
  )
}

export function RecentBlogs({ blogs }: BlogListProps) {
  const t = useTranslations()
  return (
    <>
      <h2 className="text-2xl font-bold text-neutral-900 dark:text-white mb-8">
        {t("blogs.recentBlogs")}
      </h2>
      <div className="grid gap-10 md:grid-cols-2 lg:grid-cols-3">
        {blogs.map(blog => (
          <BlogCard key={blog.id} blog={blog} />
        ))}
      </div>
    </>
  )
}
