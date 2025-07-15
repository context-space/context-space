import type { BaseRemoteAPIResponse } from "@/typings"
import type { APIKey, APIKeyResponse, CreateAPIKeyRequest } from "@/typings/api-keys"
import { useTranslations } from "next-intl"
import { useCallback, useState } from "react"

import { logger } from "@/lib/utils"
import { clientRemoteFetch } from "@/lib/utils/fetch/client"

export const apikeyLogger = logger.withTag("apikey")

// Constants
const API_KEY_LIMIT = 3

// Export the constant for use in other components if needed
export const API_KEY_MAXIMUM_LIMIT = API_KEY_LIMIT

class ApiKeyLimitError extends Error {
  constructor(limit: number, message: string) {
    super(message)
    this.name = "API_KEY_LIMIT_REACHED"
  }
}

// Custom hook for fetching individual API key details
export function useApiKeyDetails() {
  const t = useTranslations("account")
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchApiKeyDetails = useCallback(async (keyId: string): Promise<APIKeyResponse | null> => {
    try {
      setLoading(true)
      setError(null)

      apikeyLogger.info("Fetching API key details...", { keyId })
      const response = await clientRemoteFetch(`/users/me/apikeys/${keyId}`) as BaseRemoteAPIResponse<APIKeyResponse>

      if (response.data) {
        apikeyLogger.info("API key details fetched successfully:", response)
        return response.data
      } else {
        throw new Error(t("failedToLoadApiKeys"))
      }
    } catch (err: any) {
      const errorMessage = err?.message || t("failedToLoadApiKeys")
      apikeyLogger.error("Error loading API key details:", err)
      setError(errorMessage)
      return null
    } finally {
      setLoading(false)
    }
  }, [t])

  return {
    fetchApiKeyDetails,
    loading,
    error,
  }
}

export function useApiKeys() {
  const t = useTranslations("account")

  // State
  const [apiKeys, setApiKeys] = useState<APIKey[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [deleting, setDeleting] = useState<string | null>(null)
  const [newlyCreatedKey, setNewlyCreatedKey] = useState<string | null>(null)

  // Computed properties
  const isAtLimit = apiKeys.length >= API_KEY_LIMIT

  // Fetch API keys
  const fetchApiKeys = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)

      apikeyLogger.info("Fetching API keys...")
      const response = await clientRemoteFetch("/users/me/apikeys") as BaseRemoteAPIResponse<{ api_keys: APIKey[] }>

      if (response.data?.api_keys) {
        setApiKeys(response.data.api_keys)
      } else {
        setApiKeys([])
      }
      apikeyLogger.info("API keys fetched successfully:", response)
    } catch (err: any) {
      const errorMessage = err?.message || t("createApiKeyError")
      apikeyLogger.error("Error loading API keys:", err)

      // Only show error toast if it's not a 404 (which might be normal for users with no API keys)
      if (!errorMessage.includes("404") && err?.status !== 404) {
        setError(errorMessage)
      }
      setApiKeys([])
    } finally {
      setLoading(false)
    }
  }, [t])

  // Create API key
  const createApiKey = useCallback(async (request: CreateAPIKeyRequest) => {
    // Check limit before attempting to create
    if (apiKeys.length >= API_KEY_LIMIT) {
      throw new ApiKeyLimitError(API_KEY_LIMIT, t("apiKeyLimitError", { limit: API_KEY_LIMIT }))
    }

    try {
      setCreating(true)
      setError(null)

      const response = await clientRemoteFetch("/users/me/apikeys", {
        method: "POST",
        body: request,
      }) as BaseRemoteAPIResponse<APIKey>

      if (response.data?.key_value) {
        setNewlyCreatedKey(response.data.key_value)
      }

      // Add the new API key to the list (without key_value for security)
      const newApiKey: APIKey = response.data

      if (newApiKey.id && newApiKey.name) {
        setApiKeys(prev => [newApiKey, ...prev])
        // Don't show success toast here, let component handle it
      } else {
        throw new Error(t("createApiKeyError"))
      }
    } catch (err: any) {
      const errorMessage = err?.message || t("createApiKeyError")
      apikeyLogger.error("Error creating API key:", err)
      setError(errorMessage)
      throw err
    } finally {
      setCreating(false)
    }
  }, [apiKeys.length, t])

  // Delete API key
  const deleteApiKey = useCallback(async (keyId: string) => {
    try {
      setDeleting(keyId)
      setError(null)

      await clientRemoteFetch(`/users/me/apikeys/${keyId}`, {
        method: "DELETE",
      })

      // Remove from local state
      setApiKeys(prev => prev.filter(key => key.id !== keyId))
    } catch (err: any) {
      const errorMessage = err?.message || t("deleteApiKeyError")
      apikeyLogger.error("Error deleting API key:", err)
      setError(errorMessage)
      throw err
    } finally {
      setDeleting(null)
    }
  }, [t])

  // Dismiss newly created key
  const dismissNewKey = useCallback(() => {
    setNewlyCreatedKey(null)
  }, [])

  return {
    apiKeys,
    loading,
    error,
    creating,
    deleting,
    newlyCreatedKey,
    isAtLimit,
    apiKeyLimit: API_KEY_LIMIT,
    createApiKey,
    deleteApiKey,
    dismissNewKey,
    fetchApiKeys,
  }
}
