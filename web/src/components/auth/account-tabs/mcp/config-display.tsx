import type { McpConfigDisplayProps } from "@/typings/mcp"
import { Check, Copy, Plus } from "lucide-react"
import { useTranslations } from "next-intl"

import { useCallback, useState } from "react"
import { Button } from "@/components/ui/button"
import { useMcpConfig } from "@/hooks/use-mcp"
import { cn } from "@/lib/utils"

// Loading Spinner Component
function LoadingSpinner({ className }: { className?: string }) {
  return (
    <div
      className={cn(
        "animate-spin h-4 w-4 border border-muted-foreground border-t-transparent rounded-full",
        className,
      )}
    />
  )
}

// Configuration Actions Component
interface ConfigurationActionsProps {
  onCopy: () => void
  onAddToCursor: () => void
  isConfigCopied: boolean
  isAddingToCursor: boolean
  keyLoading: boolean
  hasConfig: boolean
}

function ConfigurationActions({
  onCopy,
  onAddToCursor,
  isConfigCopied,
  isAddingToCursor,
  keyLoading,
  hasConfig,
}: ConfigurationActionsProps) {
  const t = useTranslations("account")

  return (
    <div className="flex items-center gap-2 flex-wrap">
      <Button
        variant="outline"
        size="sm"
        onClick={onAddToCursor}
        disabled={keyLoading || isAddingToCursor || !hasConfig}
        className="gap-2"
      >
        {isAddingToCursor
          ? (
              <>
                <LoadingSpinner />
                {t("addingToCursor")}
              </>
            )
          : (
              <>
                <Plus className="h-4 w-4" />
                {t("addToCursor")}
              </>
            )}
      </Button>
      <Button
        variant="outline"
        size="sm"
        onClick={onCopy}
        disabled={keyLoading}
        className="gap-2"
      >
        {keyLoading
          ? (
              <>
                <LoadingSpinner />
                {t("copyMcpConfig")}
              </>
            )
          : isConfigCopied
            ? (
                <>
                  <Check className="h-4 w-4 text-green-500" />
                  {t("copied")}
                </>
              )
            : (
                <>
                  <Copy className="h-4 w-4" />
                  {t("copyMcpConfig")}
                </>
              )}
      </Button>
    </div>
  )
}

// Configuration Display Component
interface ConfigurationDisplayProps {
  config: string
  keyLoading: boolean
}

function ConfigurationDisplay({
  config,
  keyLoading,
}: ConfigurationDisplayProps) {
  const t = useTranslations("account")

  return (
    <div className="relative">
      {keyLoading && (
        <div className="absolute inset-0 bg-background/50 backdrop-blur-[1px] rounded-lg flex items-center justify-center z-10">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <LoadingSpinner />
            {t("loadingApiKeyDetails")}
          </div>
        </div>
      )}
      <pre className={cn(
        "text-xs p-4 rounded-lg border bg-muted/30 overflow-x-auto",
        "font-mono leading-relaxed max-h-96 min-w-0 w-full",
        "whitespace-pre-wrap break-all",
        keyLoading && "opacity-50",
      )}
      >
        <code className="break-all">{config}</code>
      </pre>
    </div>
  )
}

// Command Display Component
interface CommandDisplayProps {
  command: string
  onCopy: () => void
  isCommandCopied: boolean
  keyLoading: boolean
}

function CommandDisplay({
  command,
  onCopy,
  isCommandCopied,
  keyLoading,
}: CommandDisplayProps) {
  const t = useTranslations("account")

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between flex-wrap gap-2">
        <h4 className="text-sm font-medium">{t("claudeInstallCommand")}</h4>
        <Button
          variant="outline"
          size="sm"
          onClick={onCopy}
          disabled={keyLoading}
          className="gap-2"
        >
          {keyLoading
            ? (
                <>
                  <LoadingSpinner />
                  {t("copyCommand")}
                </>
              )
            : isCommandCopied
              ? (
                  <>
                    <Check className="h-4 w-4 text-green-500" />
                    {t("copied")}
                  </>
                )
              : (
                  <>
                    <Copy className="h-4 w-4" />
                    {t("copyCommand")}
                  </>
                )}
        </Button>
      </div>
      <div className="relative">
        {keyLoading && (
          <div className="absolute inset-0 bg-background/50 backdrop-blur-[1px] rounded-lg flex items-center justify-center z-10">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <LoadingSpinner />
              {t("loadingApiKeyDetails")}
            </div>
          </div>
        )}
        <pre className={cn(
          "text-xs p-4 rounded-lg border bg-muted/30 overflow-x-auto",
          "font-mono leading-relaxed min-w-0 w-full",
          "whitespace-pre-wrap break-all",
          keyLoading && "opacity-50",
        )}
        >
          <code className="break-all">{command}</code>
        </pre>
      </div>
    </div>
  )
}

// Main MCP Configuration Display Component
export function McpConfigDisplay({
  apiKey,
  onCopy,
  onAddToCursor,
  isAddingToCursor,
  integrationId,
  integrationName,
}: McpConfigDisplayProps) {
  const t = useTranslations("account")
  const [isConfigCopied, setIsConfigCopied] = useState(false)
  const [isCommandCopied, setIsCommandCopied] = useState(false)
  const { config, keyLoading, command, originalConfig, originalCommand } = useMcpConfig(apiKey, integrationId, integrationName)

  const handleCopyConfig = useCallback(async () => {
    // Use original config with real API key for copying
    await onCopy(originalConfig)
    setIsConfigCopied(true)
    setTimeout(() => setIsConfigCopied(false), 2000)
  }, [originalConfig, onCopy])

  const handleAddToCursor = useCallback(() => {
    // Use original config with real API key for adding to Cursor
    if (originalConfig) {
      onAddToCursor(originalConfig)
    }
  }, [originalConfig, onAddToCursor])

  const handleCopyCommand = useCallback(async () => {
    // Use original command with real API key for copying
    await onCopy(originalCommand)
    setIsCommandCopied(true)
    setTimeout(() => setIsCommandCopied(false), 2000)
  }, [originalCommand, onCopy])

  return (
    <div className="space-y-6 min-w-0">
      <div className="space-y-4">
        <div className="flex items-center justify-between flex-wrap gap-2">
          <div>
            <h4 className="text-sm font-medium">{t("mcpConfiguration")}</h4>
          </div>
          <ConfigurationActions
            onCopy={handleCopyConfig}
            onAddToCursor={handleAddToCursor}
            isConfigCopied={isConfigCopied}
            isAddingToCursor={isAddingToCursor}
            keyLoading={keyLoading}
            hasConfig={!!originalConfig}
          />
        </div>

        <ConfigurationDisplay
          config={config}
          keyLoading={keyLoading}
        />
      </div>

      {command && (
        <CommandDisplay
          command={command}
          onCopy={handleCopyCommand}
          isCommandCopied={isCommandCopied}
          keyLoading={keyLoading}
        />
      )}
    </div>
  )
}
