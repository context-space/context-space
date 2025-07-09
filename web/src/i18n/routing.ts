import { defineRouting } from "next-intl/routing"

export const locales = ["en", "zh", "zh-TW"] as const
export type Locale = (typeof locales)[number]
export const defaultLocale = "en"

export const routing = defineRouting({
  locales,
  defaultLocale,
  localeDetection: false,
  localePrefix: {
    mode: "as-needed",
    prefixes: {
      // 'en' will be available at `/` and `/en`
      // 'zh' will be available at `/zh`
      // 'zh-TW' will be available at `/zh-TW`
    },
  },
})
