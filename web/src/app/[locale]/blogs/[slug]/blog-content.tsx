"use client"

import type { Blog } from "contentlayer/generated"
import { CalendarIcon, UserIcon } from "lucide-react"
import { getMDXComponent } from "next-contentlayer2/hooks"
import Image from "next/image"
import { formatDate } from "@/lib/blog"
import { useMDXComponents } from "@/mdx-components"

interface BlogContentProps {
  blog: Blog
}

export function BlogContent({ blog }: BlogContentProps) {
  const Content = getMDXComponent(blog.body.code)
  const mdxComponents = useMDXComponents({})

  return (
    <article className="max-w-4xl mx-auto">

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
      <div className="relative aspect-video w-full overflow-hidden rounded-xl">
        <Image
          src={blog.image}
          alt={blog.title}
          fill
          className="w-full h-full object-cover"
        />
      </div>
      <Content components={mdxComponents} />
    </article>
  )
}
