"use client"

import { LogOut, Settings, Shield, User } from "lucide-react"
import { useTranslations } from "next-intl"
import { useTransition } from "react"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useAccountModal } from "@/hooks/use-account-modal"
import { useAuth } from "@/hooks/use-auth"
import { Link } from "@/i18n/navigation"
import { clientLogger, cn } from "@/lib/utils"

const userMenuLogger = clientLogger.withTag("user-menu")

export function UserMenu() {
  const t = useTranslations()
  const { user, signOut, isAuthenticated, isAnonymous, loading } = useAuth()
  const [isPending, startTransition] = useTransition()
  const { isOpen: isAccountModalOpen, activeTab, openModal, closeModal } = useAccountModal()

  const handleSignOut = () => {
    startTransition(async () => {
      try {
        await signOut()
      } catch (error) {
        userMenuLogger.error("Failed to sign out", { error })
      }
    })
  }

  if (loading) {
    return (
      <div className="w-8 h-8 rounded-full bg-neutral-200 dark:bg-neutral-800 animate-pulse" />
    )
  }

  if (!isAuthenticated) {
    return (
      <button
        type="button"
        className={cn(
          "relative flex-shrink-0 h-10 px-4",
          "bg-primary/10 hover:bg-primary/20",
          "border border-primary/20 rounded-md",
          "text-sm font-medium",
          "transition-all duration-200 ease-out",
          "focus:outline-none focus:ring-2 focus:ring-primary/30 focus:ring-offset-2 dark:focus:ring-offset-black/20",
        )}
      >
        <Link href={`/login?from=${encodeURIComponent(window.location.href)}`}>
          {t("header.login")}
        </Link>
      </button>
    )
  }

  const userDisplayName = isAnonymous
    ? t("common.anonymous")
    : (
        user?.user_metadata?.full_name
        || user?.user_metadata?.name
        || user?.user_metadata?.user_name
        || user?.email?.split("@")[0]
        || t("common.user")
      )

  const userAvatar = isAnonymous
    ? "/logo.svg"
    : (
        user?.user_metadata?.avatar_url
        || user?.user_metadata?.picture
      )

  const userEmail = isAnonymous ? "anonymous@example.com" : user?.email

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          className="relative h-8 w-8 rounded-full"
          disabled={isPending}
        >
          <Avatar className="h-8 w-8">
            <AvatarImage src={userAvatar} alt={userDisplayName} />
            <AvatarFallback className="bg-neutral-100 dark:bg-neutral-800">
              {isAnonymous
                ? (
                    <User className="h-4 w-4" />
                  )
                : (
                    <span className="text-sm font-medium">
                      {userDisplayName.charAt(0).toUpperCase()}
                    </span>
                  )}
            </AvatarFallback>
          </Avatar>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className={cn(
          "w-56",
          "bg-white",
          "dark:bg-white/[0.02] dark:border-white/10 dark:backdrop-blur-xl",
          "shadow-none",
          "border border-primary/15",
        )}
        align="end"
        forceMount
      >
        <DropdownMenuLabel className="font-normal">
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-medium leading-none">
              {userDisplayName}
              {isAnonymous && (
                <Shield className="inline-block w-3 h-3 ml-1 opacity-60" />
              )}
            </p>
            <p className="text-xs leading-none text-muted-foreground">
              {userEmail}
            </p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />

        {!isAnonymous && (
          <>
            <DropdownMenuItem onClick={() => openModal("profile")}>
              <Settings className="mr-2 h-4 w-4" />
              {t("account.title")}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
          </>
        )}

        <DropdownMenuItem onClick={handleSignOut} disabled={isPending}>
          <LogOut className="mr-2 h-4 w-4" />
          {isPending ? t("common.loggingOut") : t("common.logout")}
        </DropdownMenuItem>
      </DropdownMenuContent>

    </DropdownMenu>
  )
}
