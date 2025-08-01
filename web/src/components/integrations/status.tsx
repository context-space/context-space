"use client"

import type { ConnectionStatus, ProviderStatus } from "@/typings"
import {
  AlertTriangle,
  CheckCircle,
  Circle,
  Gift,
  Plug,
  XCircle,
} from "lucide-react"
import { useTranslations } from "next-intl"
import { cn } from "@/lib/utils"
import { Badge } from "../ui/badge"

interface StatusProps {
  status?: ConnectionStatus | ProviderStatus
  count?: number
  type: "badge" | "text"
  className?: string
}

const statusConfig = {
  active: {
    color: "text-green-500",
    titleKey: "integrations.status.active",
    badgeClasses: "bg-green-50 dark:bg-green-500/5 text-green-700 dark:text-green-300 border-green-200 dark:border-green-500/10",
    icon: CheckCircle,
  },
  inactive: {
    color: "text-gray-500",
    titleKey: "integrations.status.inactive",
    badgeClasses: "bg-gray-50 dark:bg-gray-500/5 text-gray-700 dark:text-gray-300 border-gray-200 dark:border-gray-500/10",
    icon: Circle,
  },
  connected: {
    color: "text-green-500",
    titleKey: "integrations.status.connected",
    badgeClasses: "bg-green-50 dark:bg-green-500/5 text-green-700 dark:text-green-300 border-green-200 dark:border-green-500/10",
    icon: Plug,
  },
  free: {
    color: "text-green-500",
    titleKey: "integrations.status.free",
    badgeClasses: "bg-green-50 dark:bg-green-500/5 text-green-700 dark:text-green-300 border-green-200 dark:border-green-500/10",
    icon: Gift,
  },
  maintenance: {
    color: "text-orange-500",
    titleKey: "integrations.status.maintenance",
    badgeClasses: "bg-orange-50 dark:bg-orange-500/5 text-orange-700 dark:text-orange-300 border-orange-200 dark:border-orange-500/10",
    icon: AlertTriangle,
  },
  deprecated: {
    color: "text-red-500",
    titleKey: "integrations.status.deprecated",
    badgeClasses: "bg-red-50 dark:bg-red-500/5 text-red-700 dark:text-red-300 border-red-200 dark:border-red-500/10",
    icon: XCircle,
  },
} as const

export function Status({ status, count = 0, type, className }: StatusProps) {
  const t = useTranslations()

  if (!status || status === "unconnected") {
    return null
  }

  const { color: statusColor, titleKey, badgeClasses, icon: Icon } = statusConfig[status]
  const title = t(titleKey as any)

  if (type === "text") {
    return (
      <div className={cn("flex items-center gap-2 text-sm text-neutral-500 dark:text-gray-400", className)}>
        <Icon className={cn("w-3 h-3", statusColor)} />
        <span>{`${count} ${title}`}</span>
      </div>
    )
  }

  return (
    <Badge variant="secondary" className={cn(badgeClasses, className)}>
      <Icon className="w-3 h-3" />
      <span className="text-xs font-medium truncate">{title}</span>
    </Badge>
  )
}
