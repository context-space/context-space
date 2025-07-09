import { useTranslations } from "next-intl"
import { usePathname } from "next/navigation"
import { UserMenu } from "@/components/auth/user-menu"
import { Link } from "@/i18n/navigation"
import { cn } from "@/lib/utils"
import { GithubButton } from "../common/github-button"
import { LocaleSwitcher } from "./locale-switcher"
import { getNavigationItems, isActive } from "./navigation-items"
import { ThemeToggle } from "./theme-toggle"

interface MobileMenuProps {
  isOpen: boolean
  onClose: () => void
}

export function MobileMenu({ isOpen, onClose }: MobileMenuProps) {
  const t = useTranslations()
  const pathname = usePathname()
  const navigationItems = getNavigationItems(t)

  if (!isOpen) {
    return null
  }

  return (
    <div className="md:hidden">
      <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3 border-t border-neutral-200/60 dark:border-white/[0.05]">
        {navigationItems.map(item => (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "block px-3 py-2 rounded-md text-base font-medium transition-colors",
              isActive(pathname, item.href)
                ? "text-primary bg-primary/5 dark:bg-primary/10"
                : "text-neutral-700 dark:text-gray-300 hover:text-neutral-900 dark:hover:text-white hover:bg-neutral-50 dark:hover:bg-neutral-800",
            )}
            onClick={onClose}
          >
            {item.label}
          </Link>
        ))}
        <div className="pt-4 pb-3 px-3 border-t border-neutral-200/60 dark:border-white/[0.05]">
          <div className="flex justify-between items-center">
            <UserMenu />
            <div className="flex items-center gap-2">
              <ThemeToggle />
              <LocaleSwitcher />
              <GithubButton />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
