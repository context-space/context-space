import type { Blog } from "contentlayer/generated"
import { allBlogs } from "contentlayer/generated"

export function getAllBlogs(): Blog[] {
  return allBlogs.sort((a, b) => new Date(b.publishedAt).getTime() - new Date(a.publishedAt).getTime())
}

export function getFeaturedBlogs(): Blog[] {
  return allBlogs.filter(blog => blog.featured).sort((a, b) => new Date(b.publishedAt).getTime() - new Date(a.publishedAt).getTime())
}

export function formatDate(date: string): string {
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  })
}
