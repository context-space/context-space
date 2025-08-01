"use client"

import { Menu, X } from "lucide-react"
import { useTranslations } from "next-intl"
import { usePathname } from "next/navigation"
import { useCallback, useState } from "react"
import { UserMenu } from "@/components/auth/user-menu"
import { useScrollToSection } from "@/hooks/use-scroll-to-section"
import { Link } from "@/i18n/navigation"
import { cn } from "@/lib/utils"
import { GithubButton } from "../common/github-button"
import { HeaderLogo } from "./header-logo"
import { LocaleSwitcher } from "./locale-switcher"
import { getNavigationItems, isActive } from "./navigation-items"
import { ThemeToggle } from "./theme-toggle"

// Constants
const SCROLL_BEHAVIOR = "smooth" as const

// Shared styles
const styles = {
  button: {
    base: "font-medium transition-colors text-neutral-700 dark:text-gray-300 hover:text-neutral-900 dark:hover:text-white",
    mobile: "block w-full text-left px-3 py-2 rounded-md text-base hover:bg-neutral-50 dark:hover:bg-neutral-800",
    desktop: "text-[15px]",
  },
  link: {
    base: "block px-3 py-2 rounded-md text-base font-medium transition-colors",
    active: "text-primary bg-primary/5 dark:bg-primary/10",
    inactive: "text-neutral-700 dark:text-gray-300 hover:text-neutral-900 dark:hover:text-white hover:bg-neutral-50 dark:hover:bg-neutral-800",
  },
  separator: "border-t border-neutral-200/60 dark:border-white/[0.05]",
  divider: "w-px h-4 bg-neutral-200 dark:bg-neutral-700",
} as const

interface HomeHeaderProps {
  className?: string
}

interface ScrollToSectionItem {
  id: string
  labelKey: string
}

// Home page specific sections
const scrollSections: ScrollToSectionItem[] = [
  { id: "features", labelKey: "nav.features" },
  { id: "roadmap", labelKey: "nav.roadmap" },
]

// Shared scroll button component
interface ScrollButtonProps {
  sectionId: string
  children: React.ReactNode
  className?: string
  onClick?: () => void
}

function ScrollButton({ sectionId, children, className, onClick }: ScrollButtonProps) {
  const scrollToSection = useScrollToSection({
    behavior: SCROLL_BEHAVIOR,
    additionalOffset: () => {
      // Add small buffer for better visual spacing
      return 20
    },
  })

  const handleClick = useCallback(() => {
    // Close mobile menu first if it exists
    if (onClick) {
      onClick()
    }

    // Small delay to ensure menu closes before scrolling
    setTimeout(() => {
      scrollToSection(sectionId)
    }, onClick ? 150 : 0)
  }, [scrollToSection, sectionId, onClick])

  return (
    <button type="button" onClick={handleClick} className={className}>
      {children}
    </button>
  )
}

// Mobile menu component
interface HomeMobileMenuProps {
  isOpen: boolean
  onClose: () => void
}

function HomeMobileMenu({ isOpen, onClose }: HomeMobileMenuProps) {
  const t = useTranslations()
  const pathname = usePathname()
  const navigationItems = getNavigationItems(t)

  if (!isOpen) return null

  return (
    <div className="md:hidden">
      <div className={cn("px-2 pt-2 pb-3 space-y-1 sm:px-3", styles.separator)}>
        {/* Home page sections */}
        {scrollSections.map(({ id, labelKey }) => (
          <ScrollButton
            key={id}
            sectionId={id}
            className={cn(styles.button.mobile, styles.button.base)}
            onClick={onClose}
          >
            {t(labelKey)}
          </ScrollButton>
        ))}

        <div className={cn("my-2", styles.separator)} />

        {/* Regular navigation items */}
        {navigationItems.map(item => (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              styles.link.base,
              isActive(pathname, item.href)
                ? styles.link.active
                : styles.link.inactive,
            )}
            onClick={onClose}
          >
            {item.label}
          </Link>
        ))}

        {/* Mobile actions */}
        <div className={cn("p-3", styles.separator)}>
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

// Desktop navigation component
function DesktopNavigation() {
  const t = useTranslations()
  const pathname = usePathname()
  const navigationItems = getNavigationItems(t)

  return (
    <nav className="hidden md:flex items-center gap-6">
      {/* Home sections */}
      {scrollSections.map(({ id, labelKey }) => (
        <ScrollButton
          key={id}
          sectionId={id}
          className={cn(
            "text-[15px] font-medium transition-colors relative",
            "after:absolute after:left-0 after:bottom-[-8px] after:h-0.5 after:w-0 after:bg-primary after:transition-all after:duration-300",
            "text-neutral-700 dark:text-gray-300 hover:text-neutral-900 dark:hover:text-white hover:after:w-full",
          )}
        >
          {t(labelKey)}
        </ScrollButton>
      ))}

      <div className={styles.divider} />

      {/* Regular navigation items */}
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

// Desktop actions component
function DesktopActions() {
  return (
    <div className="hidden md:flex items-center gap-3">
      <div className={styles.divider} />
      <GithubButton />
      <UserMenu />
    </div>
  )
}

// Mobile menu toggle button component
interface MobileMenuToggleProps {
  isOpen: boolean
  onToggle: () => void
}

function MobileMenuToggle({ isOpen, onToggle }: MobileMenuToggleProps) {
  const t = useTranslations()

  return (
    <div className="md:hidden flex items-center space-x-2">
      <button
        type="button"
        className="inline-flex items-center justify-center rounded-md p-2 text-neutral-400 hover:bg-neutral-100 dark:hover:bg-neutral-800 hover:text-neutral-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary/50"
        aria-expanded={isOpen}
        onClick={onToggle}
        aria-label={isOpen ? t("header.closeMobileMenu") : t("header.openMobileMenu")}
      >
        {isOpen
          ? (
              <X className="h-6 w-6" aria-hidden="true" />
            )
          : (
              <Menu className="h-6 w-6" aria-hidden="true" />
            )}
      </button>
    </div>
  )
}

export default function HomeHeader({ className }: HomeHeaderProps) {
  const [isOpen, setIsOpen] = useState(false)

  const toggleMenu = useCallback(() => {
    setIsOpen(prev => !prev)
  }, [])

  const closeMenu = useCallback(() => {
    setIsOpen(false)
  }, [])

  return (
    <header className="sticky top-0 z-50 w-full backdrop-blur-lg">
      <div className={cn(styles.separator, className)}>
        <div className="flex h-16 items-center justify-between">
          <HeaderLogo />

          <div className="flex items-center gap-6">
            <DesktopNavigation />
            <DesktopActions />
            <MobileMenuToggle isOpen={isOpen} onToggle={toggleMenu} />
          </div>
        </div>

        <HomeMobileMenu isOpen={isOpen} onClose={closeMenu} />
      </div>
    </header>
  )
}
