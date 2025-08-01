import type { NextConfig } from "next"
import initializeBundleAnalyzer from "@next/bundle-analyzer"
import withMDX from "@next/mdx"
import { withContentlayer } from "next-contentlayer2"
import createNextIntlPlugin from "next-intl/plugin"

const withNextIntl = createNextIntlPlugin("./src/i18n/request.ts")

const withBundleAnalyzer = initializeBundleAnalyzer({
  enabled: process.env.BUNDLE_ANALYZER_ENABLED === "true",
})

const nextConfig: NextConfig = {
  pageExtensions: ["js", "jsx", "md", "mdx", "ts", "tsx"],
  reactStrictMode: false,
  images: {
    disableStaticImages: true,
    remotePatterns: [
      {
        protocol: "https",
        hostname: "cdn-bucket.tos-cn-hongkong.volces.com",
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "img.shields.io",
        pathname: "/**",
      },
    ],
  },
}

const withMDXConfig = withMDX({
  extension: /\.mdx$/,
})

export default withNextIntl(withBundleAnalyzer(withContentlayer(withMDXConfig(nextConfig))))
