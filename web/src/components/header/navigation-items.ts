export interface NavigationItem {
  href: string
  labelKey: string
}

export function getNavigationItems(t: (key: string) => string): Array<{ href: string, label: string }> {
  return [
    { href: "/integrations", label: t("nav.integrations") },
    { href: "/docs", label: t("nav.documents") },
    { href: "/blogs", label: t("nav.blogs") },
  ]
}

export const isActive = (pathname: string, href: string) => pathname.startsWith(href)
