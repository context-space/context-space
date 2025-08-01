import type { Metadata } from "next"

import type { ReactNode } from "react"
import { GoogleAnalytics } from "@next/third-parties/google"
import NextTopLoader from "nextjs-toploader"
import { MicrosoftClarity } from "@/components/clarity"
import { ScrollArea } from "@/components/ui/scroll-area"

import { baseURL, description, keywords, siteName, title } from "@/config"
import "@/app/globals.css"

export const metadata: Metadata = {
  title: {
    default: title,
    template: `%s | ${siteName}`,
  },
  description,
  keywords: keywords.join(", "),
  authors: [{ name: siteName }],
  creator: siteName,
  publisher: siteName,
  formatDetection: {
    email: false,
    address: false,
    telephone: false,
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      "index": true,
      "follow": true,
      "max-video-preview": -1,
      "max-image-preview": "large",
      "max-snippet": -1,
    },
  },
  icons: [
    // {
    //   rel: "icon",
    //   type: "image/png",
    //   url: "/icon",
    // },
    // {
    //   rel: "shortcut icon",
    //   type: "image/png",
    //   url: "/icon",
    // },
    // {
    //   rel: "apple-touch-icon",
    //   type: "image/png",
    //   url: "/icon",
    // },
    {
      rel: "icon",
      type: "image/svg",
      url: "/logo-color.svg",
    },
  ],
  manifest: "/manifest.json",
  alternates: {
    canonical: baseURL,
    languages: {
      "en": `${baseURL}/en`,
      "zh": `${baseURL}/zh`,
      "zh-TW": `${baseURL}/zh-TW`,
    },
  },
  metadataBase: new URL(baseURL),
  openGraph: {
    title,
    description,
    siteName,
    url: baseURL,
    type: "website",
    locale: "en_US",
    alternateLocale: ["zh_CN", "zh_TW"],
    images: [
      {
        url: "/api/og",
        width: 1200,
        height: 630,
        alt: `${siteName} Open Graph Image`,
        type: "image/png",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title,
    description,
    site: "@hi_contextspace",
    creator: "@hi_contextspace",
    images: ["/api/og"],
  },
  verification: {
    google: process.env.GOOGLE_SITE_VERIFICATION,
    other: {
      me: [baseURL],
    },
  },
  category: "technology",
}

function Layout({ children }: Readonly<{ children: ReactNode }>) {
  return (
    <html suppressHydrationWarning>
      <head>
        <meta name="theme-color" content="#8b8cac" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="default" />
        <meta name="mobile-web-app-capable" content="yes" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
        <link rel="alternate" type="application/rss+xml" title={`${siteName} Blog RSS Feed`} href={`${baseURL}/feed.xml`} />
        <link rel="sitemap" type="application/xml" href="/sitemap.xml" />
      </head>
      <body className="bg-background text-foreground overscroll-none antialiased sprinkle-primary">
        <NextTopLoader
          color="#9899BD"
          initialPosition={0.08}
          crawlSpeed={200}
          height={2}
          crawl={true}
          showSpinner={false}
          easing="ease"
          speed={200}
          shadow="0 0 10px #9899BD,0 0 5px #9899BD"
        />
        <div className="h-screen overflow-hidden">
          <ScrollArea className="h-full">
            {children}
          </ScrollArea>
        </div>
        {process.env.NODE_ENV === "production" && process.env.NEXT_PUBLIC_GOOGLE_ANALYTICS_ID && (
          <GoogleAnalytics gaId={process.env.NEXT_PUBLIC_GOOGLE_ANALYTICS_ID} />
        )}
        {process.env.NODE_ENV === "production" && process.env.NEXT_PUBLIC_MICROSOFT_CLARITY && (
          <MicrosoftClarity clarityId={process.env.NEXT_PUBLIC_MICROSOFT_CLARITY} />
        )}
      </body>
    </html>
  )
}

export default Layout
