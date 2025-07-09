"use client"

import type { ConnectionStatus, ProviderStatus } from "@/typings"
import { useTranslations } from "next-intl"
import { cn } from "@/lib/utils"

interface StatusProps {
  status?: ConnectionStatus | ProviderStatus
  count?: number
  type: "badge" | "text"
  className?: string
}

const statusConfig = {
  active: {
    bg: "bg-green-500/10 dark:bg-green-500/5",
    color: "bg-green-500",
    titleKey: "integrations.status.active",
  },
  inactive: {
    bg: "bg-gray-500/10 dark:bg-gray-500/5",
    color: "bg-gray-500",
    titleKey: "integrations.status.inactive",
  },
  connected: {
    bg: "bg-primary/10 dark:bg-primary/5",
    color: "bg-primary",
    titleKey: "integrations.status.connected",
  },
  free: {
    bg: "bg-yellow-500/10 dark:bg-yellow-500/5",
    color: "bg-yellow-500",
    titleKey: "integrations.status.free",
  },
  maintenance: {
    bg: "bg-orange-500/10 dark:bg-orange-500/5",
    color: "bg-orange-500",
    titleKey: "integrations.status.maintenance",
  },
  deprecated: {
    bg: "bg-red-500/10 dark:bg-red-500/5",
    color: "bg-red-500",
    titleKey: "integrations.status.deprecated",
  },
} as const

export function Status({ status, count = 0, type, className }: StatusProps) {
  const t = useTranslations()

  if (!status || status === "unconnected") {
    return null
  }

  const { color: statusColor, bg: statusBg, titleKey } = statusConfig[status]
  const title = t(titleKey as any)

  if (type === "text") {
    return (
      <div className={cn("flex items-center gap-2 text-sm text-neutral-500 dark:text-gray-400", className)}>
        <span className={`w-1.5 h-1.5 rounded-full ${statusColor}`}></span>
        <span>{`${count} ${title}`}</span>
      </div>
    )
  }

  return (
    <div className={cn("flex items-center gap-2", className)}>
      <span className={cn("flex items-center gap-1.5 text-xs px-2 py-1 rounded-md", statusBg)}>
        <span className={cn("w-1 h-1 rounded-full", statusColor)}></span>
        {title}
      </span>
    </div>
  )
}
