"use client"

import type { AccountUser, TabType } from "@/components/auth/account-tabs/types"
import { Key, Settings, User, Zap } from "lucide-react"
import { useTranslations } from "next-intl"
import { useCallback, useEffect, useState } from "react"
import { ApiKeysTab, McpTab, ProfileTab } from "@/components/auth/account-tabs"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { cn } from "@/lib/utils"

interface AccountModalProps {
  user: AccountUser
  open: boolean
  onOpenChange: (open: boolean) => void
}

const tabConfig = [
  {
    id: "profile" as TabType,
    label: "account.profile",
    icon: User,
  },
  {
    id: "api-keys" as TabType,
    label: "account.apiKeys",
    icon: Key,
  },
  {
    id: "mcp" as TabType,
    label: "account.mcp",
    icon: Zap,
  },
]

export function AccountModal({ user, open, onOpenChange, initialTab = "profile" }: AccountModalProps & { initialTab?: TabType }) {
  const t = useTranslations()
  const [activeTab, setActiveTab] = useState<TabType>(initialTab)

  // Sync internal state with external prop
  useEffect(() => {
    setActiveTab(initialTab)
  }, [initialTab])

  const renderContent = useCallback(() => {
    switch (activeTab) {
      case "api-keys":
        return <ApiKeysTab />
      case "mcp":
        return <McpTab onSwitchToApiKeys={() => setActiveTab("api-keys")} />
      default:
        return <ProfileTab user={user} />
    }
  }, [activeTab, user])

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sprinkle-primary max-w-7xl min-w-[50vw] p-0 gap-0 bg-white dark:bg-neutral-900 border">
        <DialogHeader className="p-4">
          <DialogTitle className="text-lg font-semibold flex items-center gap-2">
            <Settings className="h-5 w-5" />
            {t("account.title")}
          </DialogTitle>
        </DialogHeader>

        <div className="h-[600px] flex flex-col sm:grid sm:grid-cols-[auto_1fr] lg:grid-cols-[240px_1fr]">
          <nav className="order-2 sm:order-1 border-t sm:border-t-0 sm:border-r border-neutral-200 dark:border-neutral-700 bg-neutral-50/50 dark:bg-neutral-800/50 sm:bg-transparent sm:dark:bg-transparent">
            <div className="hidden sm:block p-2 space-y-2">
              {tabConfig.map((tab) => {
                const Icon = tab.icon
                return (
                  <Button
                    key={tab.id}
                    variant="ghost"
                    className={cn(
                      "flex flex-col md:flex-row w-full h-10 transition-all duration-200",
                      "justify-start gap-3 lg:justify-start lg:gap-3",
                      "md:w-10 md:justify-center md:gap-0 lg:w-full",
                      activeTab === tab.id
                        ? "bg-primary/10 text-primary border border-primary/20 hover:bg-primary/15"
                        : "text-neutral-600 dark:text-neutral-400 hover:bg-primary/5 hover:text-primary border border-transparent hover:border-primary/10",
                    )}
                    onClick={() => setActiveTab(tab.id)}
                    title={t(tab.label)}
                  >
                    <Icon className="h-4 w-4" />
                    <span className="hidden lg:inline">{t(tab.label)}</span>
                  </Button>
                )
              })}
            </div>

            {/* Mobile Navigation */}
            <div className="flex sm:hidden justify-around px-4 py-3">
              {tabConfig.map((tab) => {
                const Icon = tab.icon
                return (
                  <Button
                    key={tab.id}
                    variant="ghost"
                    className={cn(
                      "flex flex-col items-center gap-1 h-auto py-2 px-3 min-w-0 transition-all duration-200 rounded-lg",
                      activeTab === tab.id
                        ? "text-primary bg-primary/10 scale-105"
                        : "text-neutral-600 dark:text-neutral-400 hover:bg-primary/5 hover:text-primary hover:scale-102",
                    )}
                    onClick={() => setActiveTab(tab.id)}
                  >
                    <Icon className="h-5 w-5" />
                    <span className="text-xs truncate max-w-16">
                      {t(tab.label)}
                    </span>
                  </Button>
                )
              })}
            </div>
          </nav>

          {/* Content */}
          <div className="order-1 sm:order-2 flex-1 p-4 sm:p-6 overflow-auto">
            {renderContent()}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
