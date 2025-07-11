"use client"

import { LogOut } from "lucide-react"
import { useTranslations } from "next-intl"
import { useTransition } from "react"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/hooks/use-auth"
import { clientLogger } from "@/lib/utils"

interface LogoutButtonProps {
  variant?: "default" | "destructive" | "outline" | "secondary" | "ghost" | "link"
  size?: "default" | "sm" | "lg" | "icon"
  className?: string
  children?: React.ReactNode
}

const logoutButtonLogger = clientLogger.withTag("logout-button")

export function LogoutButton({
  variant = "outline",
  size = "sm",
  className,
  children,
}: LogoutButtonProps) {
  const t = useTranslations()
  const { signOut } = useAuth()
  const [isPending, startTransition] = useTransition()

  const handleLogout = () => {
    startTransition(async () => {
      try {
        logoutButtonLogger.info("User initiated logout")
        await signOut()
      } catch (error) {
        logoutButtonLogger.error("Failed to logout", { error })
      }
    })
  }

  return (
    <Button
      variant={variant}
      size={size}
      onClick={handleLogout}
      disabled={isPending}
      className={className}
    >
      {children || (
        <>
          <LogOut className="mr-2 h-4 w-4" />
          {isPending ? t("common.loggingOut") : t("common.logout")}
        </>
      )}
    </Button>
  )
}
