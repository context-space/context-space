import { Menu, X } from "lucide-react"
import { useTranslations } from "next-intl"

interface MobileMenuButtonProps {
  isOpen: boolean
  onToggle: () => void
}

export function MobileMenuButton({ isOpen, onToggle }: MobileMenuButtonProps) {
  const t = useTranslations()
  return (
    <div className="md:hidden flex items-center space-x-2">
      <button
        type="button"
        className="inline-flex items-center justify-center rounded-md p-2 text-neutral-400 hover:bg-neutral-100 dark:hover:bg-neutral-800 hover:text-neutral-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary/50"
        aria-expanded="false"
        onClick={onToggle}
      >
        <span className="sr-only">{isOpen ? t("header.closeMobileMenu") : t("header.openMobileMenu")}</span>
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
