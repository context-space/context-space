"use client"

import type { AccountUser } from "./types"
import { Calendar, Mail, Shield, User } from "lucide-react"
import { useTranslations } from "next-intl"
import { LogoutButton } from "@/components/auth/logout-button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

interface ProfileTabProps {
  user: AccountUser
}

export function ProfileTab({ user }: ProfileTabProps) {
  const t = useTranslations("account")
  const tCommon = useTranslations("common")

  return (
    <div className="space-y-6">
      {/* User Avatar and Basic Info */}
      <div className="flex items-center gap-4">
        <div className="relative shrink-0 group">
          <div className="absolute inset-0 bg-gradient-to-b from-primary/3 to-transparent blur-xl -z-10 opacity-0 group-hover:opacity-60 transition-opacity duration-300"></div>
          <Avatar className="size-16 bg-transparent rounded-full transition-transform duration-200 group-hover:scale-102">
            <AvatarImage src={user.avatar} alt={user.displayName} />
            <AvatarFallback className="bg-neutral-100 dark:bg-neutral-800 rounded-none">
              {user.isAnonymous
                ? (
                    <User className="h-8 w-8" />
                  )
                : (
                    <span className="text-xl font-medium">
                      {user.displayName.charAt(0).toUpperCase()}
                    </span>
                  )}
            </AvatarFallback>
          </Avatar>
        </div>
        <div className="space-y-1">
          <div className="flex items-center gap-3">
            <h2 className="text-xl font-semibold text-neutral-900 dark:text-white">
              {user.displayName}
            </h2>
            {user.isAnonymous && (
              <Badge variant="secondary" className="text-xs bg-green-50 dark:bg-green-500/10 text-green-700 dark:text-green-400 border-green-100 dark:border-green-500/20">
                <Shield className="h-3 w-3 mr-1" />
                {tCommon("anonymous")}
              </Badge>
            )}
          </div>
          <div className="flex items-center gap-2 text-neutral-600 dark:text-gray-400 text-sm">
            <Mail className="h-4 w-4" />
            {user.email}
          </div>
        </div>
      </div>

      {/* User Details */}
      <div className="space-y-3">
        <div className="flex items-center justify-between py-2 border-b border-neutral-200 dark:border-neutral-700">
          <span className="text-sm font-medium text-neutral-700 dark:text-neutral-300">{t("userId")}</span>
          <span className="text-xs text-muted-foreground font-mono truncate max-w-32">
            {user.id}
          </span>
        </div>

        <div className="flex items-center justify-between py-2 border-b border-neutral-200 dark:border-neutral-700">
          <span className="text-sm font-medium text-neutral-700 dark:text-neutral-300">{t("accountType")}</span>
          <Badge
            variant={user.isAnonymous ? "secondary" : "default"}
            className={cn(
              "text-xs",
              user.isAnonymous
                ? "bg-orange-50 dark:bg-orange-500/10 text-orange-700 dark:text-orange-400 border-orange-100 dark:border-orange-500/20"
                : "bg-blue-50 dark:bg-blue-500/10 text-blue-700 dark:text-blue-400 border-blue-100 dark:border-blue-500/20",
            )}
          >
            {user.isAnonymous ? t("accountTypes.anonymous") : t("accountTypes.registered")}
          </Badge>
        </div>

        <div className="flex items-center justify-between py-2 border-b border-neutral-200 dark:border-neutral-700">
          <span className="text-sm font-medium text-neutral-700 dark:text-neutral-300">{t("createdAt")}</span>
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <Calendar className="h-3 w-3" />
            {new Date(user.createdAt).toLocaleDateString()}
          </div>
        </div>

        {user.lastSignInAt && (
          <div className="flex items-center justify-between py-2 border-b border-neutral-200 dark:border-neutral-700">
            <span className="text-sm font-medium text-neutral-700 dark:text-neutral-300">{t("lastSignIn")}</span>
            <span className="text-xs text-muted-foreground">
              {new Date(user.lastSignInAt).toLocaleString()}
            </span>
          </div>
        )}
      </div>

      {/* OAuth Providers (if not anonymous) */}
      {!user.isAnonymous && user.providers.length > 0 && (
        <div className="space-y-3">
          <h3 className="text-sm font-semibold text-neutral-700 dark:text-neutral-300">{t("connectedAccounts")}</h3>
          <div className="flex flex-col gap-2">
            {user.providers.map((provider: string) => (
              <div key={provider} className="flex items-center justify-between py-2 border-b border-neutral-200 dark:border-neutral-700">
                <span className="text-sm font-medium capitalize text-neutral-700 dark:text-neutral-300">{provider}</span>
                <Badge variant="outline" className="text-xs bg-green-50 dark:bg-green-500/10 text-green-700 dark:text-green-400 border-green-100 dark:border-green-500/20">
                  {tCommon("connected")}
                </Badge>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Actions */}
      <div className="pt-4">
        <LogoutButton variant="outline" size="sm" />
      </div>
    </div>
  )
}
