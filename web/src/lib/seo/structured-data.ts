import { baseURL, description, siteName } from "@/config"

export interface Person {
  "@type": "Person"
  "name": string
  "url"?: string
}

export interface Organization {
  "@type": "Organization"
  "name": string
  "url": string
  "logo"?: string
  "sameAs"?: string[]
}

export interface WebSite {
  "@context": "https://schema.org"
  "@type": "WebSite"
  "name": string
  "description": string
  "url": string
  "potentialAction"?: {
    "@type": "SearchAction"
    "target": {
      "@type": "EntryPoint"
      "urlTemplate": string
    }
    "query-input": string
  }
}

export interface SoftwareApplication {
  "@context": "https://schema.org"
  "@type": "SoftwareApplication"
  "name": string
  "description": string
  "url": string
  "applicationCategory": string
  "operatingSystem": string[]
  "offers": {
    "@type": "Offer"
    "price": string
    "priceCurrency": string
  }
  "creator": Organization
  "featureList": string[]
}

export function generateWebSiteStructuredData(): WebSite {
  return {
    "@context": "https://schema.org",
    "@type": "WebSite",
    "name": siteName,
    description,
    "url": baseURL,
    "potentialAction": {
      "@type": "SearchAction",
      "target": {
        "@type": "EntryPoint",
        "urlTemplate": `${baseURL}/integrations?search={search_term_string}`,
      },
      "query-input": "required name=search_term_string",
    },
  }
}

export function generateOrganizationStructuredData(): Organization {
  return {
    "@type": "Organization",
    "name": siteName,
    "url": baseURL,
    "logo": `${baseURL}/logo-color.svg`,
    "sameAs": [
      "https://github.com/context-space",
      "https://x.com/hi_contextspace",
    ],
  }
}

export function generateSoftwareApplicationStructuredData(): SoftwareApplication {
  return {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    "name": siteName,
    description,
    "url": baseURL,
    "applicationCategory": "DeveloperApplication",
    "operatingSystem": ["Windows", "macOS", "Linux", "Web"],
    "offers": {
      "@type": "Offer",
      "price": "0",
      "priceCurrency": "USD",
    },
    "creator": generateOrganizationStructuredData(),
    "featureList": [
      "MCP Integration Management",
      "Workflow Automation",
      "Enterprise Security",
      "AI-Powered Recommendations",
      "Real-time Collaboration",
      "API Management",
      "Cloud-Native Architecture",
    ],
  }
}

export function generateBreadcrumbStructuredData(items: Array<{ name: string, url: string }>) {
  return {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    "itemListElement": items.map((item, index) => ({
      "@type": "ListItem",
      "position": index + 1,
      "name": item.name,
      "item": item.url,
    })),
  }
}

export function generateArticleStructuredData({
  headline,
  description,
  author,
  datePublished,
  dateModified,
  url,
  image,
}: {
  headline: string
  description: string
  author: string
  datePublished: string
  dateModified?: string
  url: string
  image?: string
}) {
  return {
    "@context": "https://schema.org",
    "@type": "Article",
    headline,
    description,
    "author": {
      "@type": "Person",
      "name": author,
    },
    "publisher": generateOrganizationStructuredData(),
    datePublished,
    "dateModified": dateModified || datePublished,
    url,
    "image": image || `${baseURL}/api/og?title=${encodeURIComponent(headline)}`,
    "mainEntityOfPage": {
      "@type": "WebPage",
      "@id": url,
    },
  }
}

export function generateBlogListStructuredData(blogs: Array<{
  title: string
  description: string
  author: string
  publishedAt: string
  url: string
  image: string
}>) {
  return {
    "@context": "https://schema.org",
    "@type": "Blog",
    "name": `${siteName} Blog`,
    "description": "Latest insights, tutorials, and updates about AI engineering and context management",
    "url": `${baseURL}/blogs`,
    "publisher": generateOrganizationStructuredData(),
    "blogPost": blogs.map(blog => ({
      "@type": "BlogPosting",
      "headline": blog.title,
      "description": blog.description,
      "author": {
        "@type": "Person",
        "name": blog.author,
      },
      "datePublished": blog.publishedAt,
      "url": `${baseURL}${blog.url}`,
      "image": blog.image,
      "publisher": generateOrganizationStructuredData(),
    })),
  }
}
