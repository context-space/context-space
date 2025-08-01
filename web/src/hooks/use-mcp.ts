import type { APIKey, APIKeyResponse } from "@/typings/api-keys"
import { encode } from "js-base64"
import { useTranslations } from "next-intl"
import { useCallback, useEffect, useState } from "react"

import { toast } from "sonner"
import { baseURL } from "@/config"
import { useApiKeyDetails } from "./use-api-keys"

// Constants
const CURSOR_DEEPLINK_BASE = "cursor://anysphere.cursor-deeplink/mcp/install"

// Utility functions
/**
 * Generate MCP configuration with API key (or placeholder if not provided)
 */
export const generateMcpConfig = (apiKeyValue?: string, hideKey = false, integrationId?: string, integrationName?: string): string => {
  const token = apiKeyValue || "<YOUR_API_KEY>"
  const displayToken = hideKey && apiKeyValue ? "*".repeat(32) : token

  // Use integration-specific name and URL if provided
  const serverName = integrationName || integrationId || "context-space"
  const serverUrl = integrationId ? `${baseURL}/api/mcp/${integrationId}` : `${baseURL}/api/mcp`

  const config = {
    [serverName]: {
      url: serverUrl,
      headers: {
        Authorization: `Bearer ${displayToken}`,
      },
    },
  }
  return JSON.stringify(config, null, 2)
}

/**
 * Generate Claude MCP server installation command
 */
export const generateClaudeCommand = (apiKeyValue?: string, hideKey = false, integrationId?: string, integrationName?: string): string => {
  const token = apiKeyValue || "<YOUR_API_KEY>"
  const displayToken = hideKey && apiKeyValue ? "*".repeat(32) : token
  const serverName = integrationName || integrationId || "context-space"
  const serverUrl = integrationId ? `${baseURL}/api/mcp/${integrationId}` : `${baseURL}/api/mcp`
  return `claude mcp add --transport http "${serverName.replaceAll(" ", "-")}" ${serverUrl} --header "Authorization: Bearer ${displayToken}"`
}

const generateCursorDeepLink = (serverName: string, serverConfig: unknown): string => {
  const configBase64 = encode(JSON.stringify(serverConfig))
  return `${CURSOR_DEEPLINK_BASE}?name=${encodeURIComponent(serverName)}&config=${encodeURIComponent(configBase64)}`
}

/**
 * Hook for managing MCP configuration loading and generation
 */
export const useMcpConfig = (apiKey: APIKey, integrationId?: string, integrationName?: string) => {
  const { fetchApiKeyDetails, loading: keyLoading } = useApiKeyDetails()
  const [actualApiKey, setActualApiKey] = useState<APIKeyResponse | null>(null)
  const [config, setConfig] = useState<string>("")
  const [command, setCommand] = useState<string>("")
  const [originalConfig, setOriginalConfig] = useState<string>("")
  const [originalCommand, setOriginalCommand] = useState<string>("")

  const loadApiKeyDetails = useCallback(async () => {
    if (!apiKey.id) return

    try {
      const keyDetails = await fetchApiKeyDetails(apiKey.id)
      if (keyDetails) {
        setActualApiKey(keyDetails)

        // Generate display versions (with hidden API key)
        const displayConfig = generateMcpConfig(keyDetails.key_value, true, integrationId, integrationName)
        const displayCommand = generateClaudeCommand(keyDetails.key_value, true, integrationId, integrationName)
        setConfig(displayConfig)
        setCommand(displayCommand)

        // Generate original versions (with real API key for copying)
        const fullConfig = generateMcpConfig(keyDetails.key_value, false, integrationId, integrationName)
        const fullCommand = generateClaudeCommand(keyDetails.key_value, false, integrationId, integrationName)
        setOriginalConfig(fullConfig)
        setOriginalCommand(fullCommand)
      }
    } catch (error) {
      console.error("Failed to load API key details:", error)
      setActualApiKey(null)
      setConfig("")
      setCommand("")
      setOriginalConfig("")
      setOriginalCommand("")
    }
  }, [apiKey.id, fetchApiKeyDetails, integrationId, integrationName])

  useEffect(() => {
    loadApiKeyDetails()
  }, [loadApiKeyDetails])

  return {
    config,
    command,
    originalConfig,
    originalCommand,
    actualApiKey,
    keyLoading,
    refetch: loadApiKeyDetails,
  }
}

/**
 * Hook for managing "Add to Cursor" functionality
 */
export const useAddToCursor = () => {
  const t = useTranslations("account")
  const [isAdding, setIsAdding] = useState(false)

  const addToCursor = useCallback(async (config: string) => {
    try {
      setIsAdding(true)

      const mcpConfig = JSON.parse(config)
      const serverName = Object.keys(mcpConfig)[0] // "github"
      const serverConfig = mcpConfig[serverName]

      const cursorUrl = generateCursorDeepLink(serverName, serverConfig)

      // Try to open the cursor:// protocol URL
      window.location.href = cursorUrl

      toast.success(t("addedToCursor"))
    } catch (error) {
      console.error("Failed to add to Cursor:", error)
      toast.error(t("addToCursorFailed"))
    } finally {
      setTimeout(() => {
        setIsAdding(false)
      }, 1000)
    }
  }, [t])

  return {
    isAdding,
    addToCursor,
  }
}
