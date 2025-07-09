"use client"

import type { Blog } from "contentlayer/generated"
import { CalendarIcon, UserIcon } from "lucide-react"
import { getMDXComponent } from "next-contentlayer2/hooks"
import Image from "next/image"
import { formatDate } from "@/lib/blog"

interface BlogContentProps {
  blog: Blog
}

export function BlogContent({ blog }: BlogContentProps) {
  const Content = getMDXComponent(blog.body.code)

  return (
    <article className="max-w-4xl mx-auto">
      {/* Hero Image */}

      <header className="my-4">
        <h1 className="text-3xl lg:text-4xl font-bold text-neutral-900 dark:text-white mb-4 leading-tight">
          {blog.title}
        </h1>

        <div className="flex flex-wrap items-center gap-4 text-sm text-neutral-500 dark:text-gray-500 pb-4 border-b border-base">
          <div className="flex items-center gap-2">
            <UserIcon className="w-4 h-4" />
            <span>{blog.author}</span>
          </div>
          <div className="flex items-center gap-2">
            <CalendarIcon className="w-4 h-4" />
            <span>{formatDate(blog.publishedAt)}</span>
          </div>
        </div>
      </header>
      <div className="relative aspect-video w-full overflow-hidden rounded-xl mb-6">
        <Image
          src={blog.image}
          alt={blog.title}
          fill
          className="w-full h-full object-cover"
        />
      </div>
      <div className="prose prose-neutral dark:prose-invert max-w-none text-base leading-7 prose-h1:text-2xl prose-h1:font-bold prose-h1:mb-4 prose-h2:text-xl prose-h2:font-semibold prose-h2:mb-3 prose-h3:text-lg prose-h3:font-semibold prose-h3:mb-3 prose-h4:text-base prose-h4:font-medium prose-h4:mb-2 prose-h5:text-sm prose-h5:font-medium prose-h5:mb-2 prose-h6:text-sm prose-h6:font-medium prose-h6:mb-2 prose-p:text-neutral-800 prose-p:font-normal prose-p:leading-7 prose-p:mb-4 dark:prose-p:text-neutral-200 prose-strong:font-semibold prose-strong:text-neutral-900 dark:prose-strong:text-white prose-ul:mb-4 prose-ol:mb-4 prose-li:mb-1">
        <Content />
      </div>
    </article>
  )
}
