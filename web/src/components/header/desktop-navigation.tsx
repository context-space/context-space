import { useTranslations } from "next-intl"
import { usePathname } from "next/navigation"
import { Link } from "@/i18n/navigation"
import { cn } from "@/lib/utils"
import { getNavigationItems, isActive } from "./navigation-items"

export function DesktopNavigation() {
  const t = useTranslations()
  const pathname = usePathname()
  const navigationItems = getNavigationItems(t)

  return (
    <nav className="hidden md:flex justify-center items-center space-x-8">
      {navigationItems.map(item => (
        <Link
          key={item.href}
          href={item.href}
          className={cn(
            "text-[15px] font-medium transition-colors relative",
            "after:absolute after:left-0 after:bottom-[-8px] after:h-0.5 after:w-0 after:bg-primary after:transition-all after:duration-300",
            isActive(pathname, item.href)
              ? "text-primary hover:text-primary/80"
              : "text-neutral-700 dark:text-gray-300 hover:text-neutral-900 dark:hover:text-white hover:after:w-full",
          )}
        >
          {item.label}
        </Link>
      ))}
    </nav>
  )
}
