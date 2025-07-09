import { UserMenu } from "@/components/auth/user-menu"
import { GithubButton } from "@/components/common/github-button"
import { LocaleSwitcher } from "./locale-switcher"
import { ThemeToggle } from "./theme-toggle"

export function DesktopActions() {
  return (
    <div className="hidden md:flex justify-center items-center space-x-2">
      <ThemeToggle />
      <LocaleSwitcher />
      <GithubButton />
      <div className="w-px h-4 bg-neutral-200 dark:bg-neutral-700 mr-4" />
      <UserMenu />
    </div>
  )
}
