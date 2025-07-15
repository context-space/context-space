"use client"

import { Key, Shield } from "lucide-react"
import { useTranslations } from "next-intl"

import { memo, useEffect } from "react"

import { Alert, AlertDescription } from "@/components/ui/alert"
import { apikeyLogger, useApiKeys } from "@/hooks/use-api-keys"
import { ApiKeyForm, ApiKeyItem, NewKeyDisplay } from "./api-keys"

// Error boundary component for better error handling
function ErrorDisplay({ error, onRetry }: { error: string, onRetry: () => void }) {
  const t = useTranslations("account")

  return (
    <div className="text-center py-8" role="alert">
      <div className="p-4 border border-red-200/40 rounded-lg bg-red-50/40 dark:bg-red-950/30 dark:border-red-800/30">
        <p className="text-red-800/80 dark:text-red-200/80 mb-2">{t("failedToLoadApiKeys")}</p>
        <p className="text-sm text-red-600/70 dark:text-red-300/70 mb-4">{error}</p>
        <button
          onClick={onRetry}
          type="button"
          className="text-sm text-red-800/80 dark:text-red-200/80 underline hover:no-underline"
        >
          {t("tryAgain")}
        </button>
      </div>
    </div>
  )
}

// Loading state component
function LoadingDisplay({ message }: { message: string }) {
  return (
    <div className="text-center py-8" role="status" aria-live="polite">
      <div className="flex items-center justify-center space-x-2">
        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-current"></div>
        <p className="text-muted-foreground">{message}</p>
      </div>
    </div>
  )
}

// Empty state component
function EmptyDisplay({ title, description }: { title: string, description: string }) {
  return (
    <div className="text-center py-8">
      <Key className="h-12 w-12 mx-auto text-muted-foreground mb-4" aria-hidden="true" />
      <h4 className="text-lg font-medium mb-2">{title}</h4>
      <p className="text-muted-foreground mb-4">{description}</p>
    </div>
  )
}

// Main component
export const ApiKeysTab = memo(() => {
  const t = useTranslations("account")
  const {
    apiKeys,
    loading,
    error,
    creating,
    deleting,
    newlyCreatedKey,
    isAtLimit,
    apiKeyLimit,
    createApiKey,
    deleteApiKey,
    dismissNewKey,
    fetchApiKeys,
  } = useApiKeys()

  // Initial load
  useEffect(() => {
    apikeyLogger.info("Fetching API keys on mount")
    fetchApiKeys()
  }, [fetchApiKeys])

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h3 className="text-lg font-semibold">{t("apiKeys")}</h3>
        <div className="text-sm text-muted-foreground mt-1 space-y-1">
          <p>{t("apiKeysDescription")}</p>
        </div>
      </div>

      {/* Security Warning */}
      <Alert className="border-amber-200/40 bg-amber-50/40 dark:border-amber-800/30 dark:bg-amber-950/30">
        <Shield className="h-4 w-4 text-amber-700/70 dark:text-amber-400/70" />
        <AlertDescription className="text-amber-800/80 dark:text-amber-300/80">
          {t("apiKeySecurityWarning")}
        </AlertDescription>
      </Alert>

      {/* Newly created key display */}
      {newlyCreatedKey && (
        <NewKeyDisplay keyValue={newlyCreatedKey} onDismiss={dismissNewKey} />
      )}

      {/* Actions */}
      <div className="flex justify-between items-center">
        <ApiKeyForm
          onSubmit={createApiKey}
          isCreating={creating}
          isAtLimit={isAtLimit}
          apiKeyLimit={apiKeyLimit}
          currentCount={apiKeys.length}
        />
      </div>

      {/* Content */}
      {error
        ? (
            <ErrorDisplay error={error} onRetry={fetchApiKeys} />
          )
        : loading
          ? (
              <LoadingDisplay message={t("loadingApiKeys")} />
            )
          : apiKeys.length === 0
            ? (
                <EmptyDisplay
                  title={t("noApiKeys")}
                  description={t("createFirstApiKey")}
                />
              )
            : (
                <div className="space-y-4" role="list">
                  {apiKeys.map(apiKey => (
                    <ApiKeyItem
                      key={apiKey.id}
                      apiKey={apiKey}
                      onDelete={deleteApiKey}
                      isDeleting={deleting === apiKey.id}
                    />
                  ))}
                </div>
              )}
    </div>
  )
})
