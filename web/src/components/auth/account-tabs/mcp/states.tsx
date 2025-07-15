import { Key, RefreshCw, Zap } from "lucide-react"
import { useTranslations } from "next-intl"

import { Button } from "@/components/ui/button"

interface EmptyStateProps {
  onCreateApiKey: () => void
  onReload?: () => void
  isReloading?: boolean
}

export function EmptyState({ onCreateApiKey, onReload, isReloading = false }: EmptyStateProps) {
  const t = useTranslations()

  return (
    <div className="text-center py-8">
      <Key className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
      <h3 className="text-lg font-semibold mb-2">{t("account.mcpNoApiKeys")}</h3>
      <p className="text-muted-foreground text-sm mb-4">
        {t("account.mcpNoApiKeysDescription")}
      </p>
      <div className="flex flex-col sm:flex-row gap-3 justify-center items-center">
        <Button
          className="min-w-[140px]"
          variant="outline"
          onClick={onCreateApiKey}
        >
          {t("account.createApiKeyFirst")}
        </Button>
        {onReload && (
          <Button
            className="min-w-[120px]"
            variant="ghost"
            size="sm"
            onClick={onReload}
            disabled={isReloading}
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${isReloading ? "animate-spin" : ""}`} />
            {isReloading ? t("common.loading") : t("common.reload")}
          </Button>
        )}
      </div>
    </div>
  )
}

interface ErrorStateProps {
  error: string
  onRetry: () => void
}

export function ErrorState({ error, onRetry }: ErrorStateProps) {
  const t = useTranslations("account")

  return (
    <div className="text-center py-8">
      <Zap className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
      <h3 className="text-lg font-semibold mb-2">{t("mcp")}</h3>
      <p className="text-muted-foreground text-sm mb-4">
        {t("unifiedMcpDescription")}
      </p>
      <div className="text-sm text-red-500 bg-red-50 dark:bg-red-900/20 p-3 rounded-lg border border-red-200 dark:border-red-800">
        {t("failedToLoadApiKeys")}
        :
        {error}
      </div>
      <Button
        className="mt-4"
        variant="outline"
        onClick={onRetry}
      >
        {t("tryAgain")}
      </Button>
    </div>
  )
}

export function LoadingState() {
  const t = useTranslations("account")

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold">{t("mcp")}</h3>
        <p className="text-sm text-muted-foreground mt-1">
          {t("unifiedMcpDescription")}
        </p>
      </div>

      <div className="flex items-center justify-center py-4">
        <div className="animate-pulse flex items-center gap-2">
          <div className="h-4 w-4 bg-muted rounded-full"></div>
          <div className="h-4 bg-muted rounded w-32"></div>
        </div>
      </div>
    </div>
  )
}
