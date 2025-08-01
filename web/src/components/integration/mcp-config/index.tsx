"use client"

import { ChevronDown, Globe } from "lucide-react"
import { useTranslations } from "next-intl"
import { useCallback, useEffect, useMemo, useState } from "react"
import { ApiKeySelector } from "@/components/auth/account-tabs/mcp/api-key-selector"
import { McpConfigDisplay } from "@/components/auth/account-tabs/mcp/config-display"
import { EmptyState, ErrorState, LoadingState } from "@/components/auth/account-tabs/mcp/states"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { useAccountModal } from "@/hooks/use-account-modal"
import { useApiKeys } from "@/hooks/use-api-keys"
import { useAuth } from "@/hooks/use-auth"
import { useCopyToClipboard } from "@/hooks/use-clipboard"
import { useAddToCursor } from "@/hooks/use-mcp"
import { useRouter } from "@/i18n/navigation"
import { cn } from "@/lib/utils"

interface McpConfigProps {
  onSwitchToApiKeys?: () => void
  authType?: "oauth" | "apikey" | "none"
  isConnected?: boolean
  integrationId?: string
  integrationName?: string
}

export function McpConfig({ onSwitchToApiKeys, authType, isConnected, integrationId, integrationName }: McpConfigProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [selectedApiKeyId, setSelectedApiKeyId] = useState<string>("")
  const [isReloading, setIsReloading] = useState(false)
  const [showUnifiedConfig, setShowUnifiedConfig] = useState(false)
  const t = useTranslations()
  const router = useRouter()

  const { isAuthenticated } = useAuth()
  const { openModal } = useAccountModal()
  const { apiKeys, loading, error, fetchApiKeys } = useApiKeys()
  const copyToClipboard = useCopyToClipboard()
  const { isAdding, addToCursor } = useAddToCursor()

  // Wrap copyToClipboard with translation messages
  const handleCopy = useCallback(async (text: string) => {
    return await copyToClipboard(text, t("account.mcpConfigCopied"), t("account.copyMcpConfigFailed"))
  }, [copyToClipboard, t])

  // Handle reload with loading state
  const handleReload = useCallback(async () => {
    setIsReloading(true)
    try {
      await fetchApiKeys()
    } finally {
      setIsReloading(false)
    }
  }, [fetchApiKeys])

  // Load API keys on mount
  useEffect(() => {
    fetchApiKeys()
  }, [fetchApiKeys])

  // Auto-select first API key when available
  useEffect(() => {
    if (apiKeys.length > 0 && !selectedApiKeyId) {
      setSelectedApiKeyId(apiKeys[0].id)
    }
  }, [apiKeys, selectedApiKeyId])

  const selectedApiKey = useMemo(() =>
    apiKeys.find(key => key.id === selectedApiKeyId), [apiKeys, selectedApiKeyId])

  const handleCreateApiKey = useCallback(() => {
    if (onSwitchToApiKeys) {
      onSwitchToApiKeys()
    } else {
      // Open account modal with API Keys tab
      openModal("api-keys")
    }
  }, [isAuthenticated, authType, isConnected, onSwitchToApiKeys, openModal, t])

  const hasApiKeys = apiKeys.length > 0

  const renderContent = () => {
    // Check authentication first
    if (!isAuthenticated) {
      return (
        <div className="text-center py-8">
          <div className="text-sm text-muted-foreground mb-4">
            {t("integrations.connect.loginRequiredDescription")}
          </div>
          <Button
            type="button"
            onClick={() => {
              const currentPath = encodeURIComponent(window.location.pathname)
              router.push(`/login?from=${currentPath}`)
            }}
          >
            {t("header.login")}
          </Button>
        </div>
      )
    }

    // Check connection if required
    const requiresConnection = authType && authType !== "none"
    if (requiresConnection && !isConnected) {
      return (
        <div className="text-center py-8">
          <div className="text-sm text-muted-foreground mb-2">
            {t("playground.connectionRequiredDescription")}
          </div>
        </div>
      )
    }

    // Error state
    if (error) {
      return (
        <div className="space-y-4">
          <ErrorState error={error} onRetry={fetchApiKeys} />
        </div>
      )
    }

    // Loading state (only for initial load, not for reload)
    if (loading && !isReloading) {
      return <LoadingState />
    }

    // Empty state
    if (!hasApiKeys) {
      return (
        <div className="space-y-4">
          <EmptyState
            onCreateApiKey={handleCreateApiKey}
            onReload={handleReload}
            isReloading={isReloading}
          />
        </div>
      )
    }

    return (
      <div className="space-y-4">
        <ApiKeySelector
          apiKeys={apiKeys}
          selectedApiKeyId={selectedApiKeyId}
          onSelect={setSelectedApiKeyId}
        />

        {/* Configuration Type Toggle - Only show for specific integrations */}
        {integrationId && integrationName && integrationName !== "Context Space" && (
          <div className="flex gap-2">
            <Button
              variant={showUnifiedConfig ? "outline" : "default"}
              size="sm"
              onClick={() => setShowUnifiedConfig(false)}
              className="flex items-center gap-2 flex-1 sm:flex-initial"
            >
              {integrationName || integrationId}
            </Button>
            <Button
              variant={showUnifiedConfig ? "default" : "outline"}
              size="sm"
              onClick={() => setShowUnifiedConfig(true)}
              className="flex items-center gap-2 flex-1 sm:flex-initial"
            >
              Context Space
            </Button>
          </div>
        )}

        {showUnifiedConfig && integrationId && integrationName && integrationName !== "Context Space" && (
          <Alert className={cn(
            "border-green-200/40 bg-green-50/40 dark:border-green-800/30 dark:bg-green-950/30",
            "text-green-700/70 dark:text-green-400/70",
          )}
          >
            <Globe className={cn(
              "h-4 w-4",
            )}
            />
            <AlertDescription className={cn(
              "text-green-800/80 dark:text-green-300/80",
            )}
            >
              <div className="space-y-1">
                <p className="font-medium">{t("account.unifiedMcpTitle")}</p>
                <p className="text-sm">{t("account.unifiedMcpSubtitle")}</p>
              </div>
            </AlertDescription>
          </Alert>
        )}

        {/* MCP Configuration Display */}
        {selectedApiKey && (
          <McpConfigDisplay
            apiKey={selectedApiKey}
            onCopy={handleCopy}
            onAddToCursor={addToCursor}
            isAddingToCursor={isAdding}
            integrationId={showUnifiedConfig ? undefined : integrationId}
            integrationName={showUnifiedConfig ? "Context Space" : integrationName}
          />
        )}
      </div>
    )
  }

  return (
    <Card data-testid="mcp-config-section" className="pb-0 bg-white/60 dark:bg-white/[0.02] border-base backdrop-blur-sm shadow-none">
      <Collapsible open={isOpen} onOpenChange={setIsOpen}>
        <CardHeader onClick={() => setIsOpen(!isOpen)}>
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <CardTitle className="text-lg font-semibold flex items-center justify-between">
                {showUnifiedConfig
                  ? t("account.unifiedMcpTitle")
                  : integrationName
                    ? t("account.mcpConfigurationForIntegration", { integrationName })
                    : t("account.mcp")}
                <div className="flex items-center gap-2">
                  <CollapsibleTrigger asChild>
                    <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                      <ChevronDown className={cn("h-4 w-4 transition-transform duration-200", isOpen && "rotate-180")} />
                      <span className="sr-only">{t("common.toggle")}</span>
                    </Button>
                  </CollapsibleTrigger>
                </div>
              </CardTitle>
              <div className="text-sm text-muted-foreground my-3">
                {showUnifiedConfig
                  ? t("account.unifiedMcpSubtitle")
                  : integrationName
                    ? t("account.mcpDescriptionForIntegration", { integrationName })
                    : t("account.mcpDescription")}
              </div>
            </div>
          </div>
        </CardHeader>
        <CollapsibleContent>
          <CardContent className="pt-0 px-0">
            <div className="border-t border-base p-5">
              {renderContent()}
            </div>
          </CardContent>
        </CollapsibleContent>
      </Collapsible>
    </Card>
  )
}
