import { siteName } from "@/config"
import { Link } from "@/i18n/navigation"
import { LogoAvatar } from "../common/avatar"

export function HeaderLogo() {
  return (
    <Link href="/" className="flex items-center gap-3 dark:text-foreground text-primary/90">
      <LogoAvatar alt={siteName} className="size-8" />
      <span className="text-xl font-semibold font-mono">
        {siteName}
      </span>
    </Link>
  )
}
