"use client"

import type { McpTabProps } from "@/typings/mcp"
import { useTranslations } from "next-intl"
import { useCallback, useEffect, useMemo, useState } from "react"

import { toast } from "sonner"
import { useApiKeys } from "@/hooks/use-api-keys"
import { useCopyToClipboard } from "@/hooks/use-clipboard"
import { useAddToCursor } from "@/hooks/use-mcp"
import { ApiKeySelector } from "./mcp/api-key-selector"
import { McpConfigDisplay } from "./mcp/config-display"
import { EmptyState, ErrorState, LoadingState } from "./mcp/states"

export function McpTab({ onSwitchToApiKeys }: McpTabProps = {}) {
  const t = useTranslations("account")
  const [selectedApiKeyId, setSelectedApiKeyId] = useState<string>("")

  const { apiKeys, loading, error, fetchApiKeys } = useApiKeys()
  const copyToClipboard = useCopyToClipboard()
  const { isAdding, addToCursor } = useAddToCursor()

  // Wrap copyToClipboard with translation messages
  const handleCopy = useCallback(async (text: string) => {
    return await copyToClipboard(text, t("mcpConfigCopied"), t("copyMcpConfigFailed"))
  }, [copyToClipboard, t])

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
      toast.info("Please go to the API Keys tab to create an API key")
    }
  }, [onSwitchToApiKeys])

  // Loading state
  if (loading) {
    return <LoadingState />
  }

  // Error state
  if (error) {
    return (
      <div className="space-y-6">
        <ErrorState error={error} onRetry={fetchApiKeys} />
      </div>
    )
  }

  // Empty state
  if (apiKeys.length === 0) {
    return (
      <div className="space-y-6">
        <EmptyState onCreateApiKey={handleCreateApiKey} />
      </div>
    )
  }

  // Main content
  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h3 className="text-lg font-semibold">{t("mcp")}</h3>
        <p className="text-sm text-muted-foreground mt-1">
          {t("unifiedMcpDescription")}
        </p>
      </div>

      {/* API Key Selection */}
      <ApiKeySelector
        apiKeys={apiKeys}
        selectedApiKeyId={selectedApiKeyId}
        onSelect={setSelectedApiKeyId}
      />

      {/* MCP Configuration Display */}
      {selectedApiKey && (
        <McpConfigDisplay
          apiKey={selectedApiKey}
          onCopy={handleCopy}
          onAddToCursor={addToCursor}
          isAddingToCursor={isAdding}
          integrationName="Context Space"
        />
      )}
    </div>
  )
}
