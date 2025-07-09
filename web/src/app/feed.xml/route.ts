import { allBlogs } from "contentlayer/generated"
import { baseURL, siteName } from "@/config"

export async function GET() {
  // Sort blogs by publication date (newest first)
  const sortedBlogs = allBlogs.sort((a, b) =>
    new Date(b.publishedAt).getTime() - new Date(a.publishedAt).getTime(),
  )

  const rss = `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
<title>${siteName} Blog</title>
<description>Latest insights, tutorials, and updates about AI engineering and context management from ${siteName}</description>
<link>${baseURL}/blogs</link>
<language>en-us</language>
<lastBuildDate>${new Date().toUTCString()}</lastBuildDate>
<atom:link href="${baseURL}/feed.xml" rel="self" type="application/rss+xml" />
<ttl>60</ttl>
${sortedBlogs.map(blog => `
<item>
<title><![CDATA[${blog.title}]]></title>
<description><![CDATA[${blog.description}]]></description>
<pubDate>${new Date(blog.publishedAt).toUTCString()}</pubDate>
<link>${baseURL}${blog.url}</link>
<guid isPermaLink="true">${baseURL}${blog.url}</guid>
<author>${blog.author}</author>
${blog.category ? `<category>${blog.category}</category>` : ""}
</item>`).join("")}
</channel>
</rss>`

  return new Response(rss, {
    headers: {
      "Content-Type": "application/xml",
      "Cache-Control": "public, s-maxage=3600, stale-while-revalidate=86400",
    },
  })
}
