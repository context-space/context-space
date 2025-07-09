export * from "./structured-data"

export function generateCanonicalUrl(path: string, baseUrl: string): string {
  const cleanPath = path.startsWith("/") ? path : `/${path}`
  const cleanBaseUrl = baseUrl.endsWith("/") ? baseUrl.slice(0, -1) : baseUrl
  return `${cleanBaseUrl}${cleanPath}`
}

export function generateLocalizedUrls(
  path: string,
  baseUrl: string,
  locales: string[],
): Record<string, string> {
  const cleanPath = path.replace(/^\/[a-z]{2}(-[A-Z]{2})?/, "")
  const cleanBaseUrl = baseUrl.endsWith("/") ? baseUrl.slice(0, -1) : baseUrl

  return locales.reduce((acc, locale) => {
    acc[locale] = `${cleanBaseUrl}/${locale}${cleanPath}`
    return acc
  }, {} as Record<string, string>)
}

export function truncateDescription(text: string, maxLength: number = 160): string {
  if (text.length <= maxLength) return text
  return `${text.substring(0, maxLength - 3).trim()}...`
}

export function generateMetaKeywords(keywords: string[]): string {
  return keywords.join(", ")
}

export function generateOpenGraphImage(
  title: string,
  description?: string,
  baseUrl: string = "",
): string {
  const params = new URLSearchParams()
  params.set("title", title)
  if (description) {
    params.set("description", description)
  }
  return `${baseUrl}/api/og?${params.toString()}`
}
