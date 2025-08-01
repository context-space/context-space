import type { Blog } from "contentlayer/generated"
import { CalendarIcon, UserIcon } from "lucide-react"
import Image from "next/image"
import { Link } from "@/i18n/navigation"
import { formatDate } from "@/lib/blog"
import { cn } from "@/lib/utils"

interface BlogCardProps {
  blog: Blog
}

export function BlogCard({ blog }: BlogCardProps) {
  return (
    <Link
      href={blog.url}
      className="block h-full transition-transform hover:scale-[1.01] duration-300"
    >
      <div className={cn(
        "group relative flex flex-col h-full overflow-hidden rounded-xl bg-white/[0.02]",
        "border border-primary/15 dark:border-white/10 hover:border-primary/30 dark:hover:border-primary/40",
        "transition-all duration-300",
      )}
      >
        {/* Blog Image */}
        <div className="relative aspect-video w-full overflow-hidden">
          <Image
            src={blog.image}
            alt={blog.title}
            fill
            loading="lazy"
            quality={70}
            className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
        </div>

        {/* Content */}
        <div className="flex flex-col flex-grow p-6">
          <h2 className="text-lg font-semibold text-neutral-900 dark:text-white mb-3 group-hover:text-primary transition-colors line-clamp-2">
            {blog.title}
          </h2>

          <p className="text-sm text-neutral-600 dark:text-gray-400 mb-4 line-clamp-3 flex-grow">
            {blog.description}
          </p>

          <div className="flex items-center justify-between text-xs text-neutral-500 dark:text-gray-500 mt-auto pt-4 border-t border-primary/10 dark:border-white/5">
            <div className="flex items-center gap-3 justify-between w-full">
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
        </div>
      </div>
    </Link>
  )
}
