"use client"

import { useTranslations } from "next-intl"
import { useSearchParams } from "next/navigation"
import { useEffect, useState } from "react"
import Logo from "@/components/common/logo"
import { FramelessLayout } from "@/components/layouts"
import { Button } from "@/components/ui/button"
import { clientLogger } from "@/lib/utils"

const mcpOauthLogger = clientLogger.withTag("mcp-oauth")
export default function MCPOAuthCallbackPage() {
  const searchParams = useSearchParams()
  const t = useTranslations()

  // Get all possible parameters
  const provider = searchParams.get("provider")
  const url = searchParams.get("url")
  const success = searchParams.get("success")
  const error = searchParams.get("error")
  const errorDescription = searchParams.get("error_description")
  const code = searchParams.get("code")
  const oauthStateId = searchParams.get("oauth_state_id")

  const [processing, setProcessing] = useState(true)
  const [localError, setLocalError] = useState<string | null>(null)
  const [isClient, setIsClient] = useState(false)
  const [isPopup, setIsPopup] = useState(false)
  // Check if we're on the client side and if this is a popup
  useEffect(() => {
    setIsClient(true)
    setIsPopup(typeof window !== "undefined" && window.opener !== null)
  }, [])

  useEffect(() => {
    const handleCallback = async () => {
      try {
        // Case 1: Redirect to OAuth provider (when 'url' parameter is present)
        if (url && !success && !error && !code) {
          mcpOauthLogger.info("Redirecting to OAuth provider", { url })
          window.location.href = decodeURIComponent(url)

          return
        }

        // Case 2: Handle OAuth callback result
        // Check if success=true for successful authentication (error is just a message)
        if (success === "true") {
          mcpOauthLogger.info("OAuth authentication successful", { success, code, error })
          setProcessing(false)

          // If in popup, reload parent and close
          if (isPopup) {
            window.opener.location.reload()
            setTimeout(() => {
              window.close()
            }, 1000)
          } else {
            // If not in popup, redirect to integrations page
            window.location.href = "/integrations"
          }

          return
        }

        // Case 3: Handle errors - only if success is not 'true'
        if (success && success !== "true") {
          mcpOauthLogger.error("OAuth failed", { success, error, errorDescription })
          setLocalError(errorDescription || `Authentication failed: ${success}`)
          setProcessing(false)

          // If in popup, notify parent but DON'T close the popup
          if (isPopup) {
            window.opener.postMessage({
              type: "oauth-error",
              error: errorDescription || `Authentication failed: ${success}`,
            }, window.location.origin)
          }

          return
        }

        // Case 4: No valid parameters - should not happen normally
        if (!url && !success && !error && !code) {
          setLocalError(t("auth.oauth.mcp.invalidCallback"))
          setProcessing(false)
        }
      } catch (err: any) {
        mcpOauthLogger.error("OAuth callback error", { error: err })
        setLocalError(err.message || t("auth.oauth.mcp.authenticationFailed"))
        setProcessing(false)
      }
    }

    // Only execute if we're on the client side
    if (isClient) {
      handleCallback()
    }
  }, [url, success, error, errorDescription, code, oauthStateId, isClient, isPopup])

  // Only show error if success is not 'true' (when success=true, error is just a message)
  const displayError = success === "true" ? localError : (localError || errorDescription || error)

  const handleClosePopup = () => {
    if (isClient && isPopup) {
      window.close()
    }
  }

  const handleRetry = () => {
    if (isClient && isPopup) {
      window.opener.location.reload()
      window.close()
    } else if (isClient) {
      window.location.href = "/integrations"
    }
  }

  // Don't render popup-specific UI until we know we're on the client
  if (!isClient) {
    return (
      <div className="min-h-screen bg-neutral-50 dark:bg-neutral-950 flex items-center justify-center p-4">
        <div className="flex flex-col items-center space-y-6 max-w-md text-center">
          <Logo size={72} />
          <div className="w-8 h-8 border-2 border-neutral-200 dark:border-white/10 border-t-neutral-900 dark:border-t-white rounded-full animate-spin" />
        </div>
      </div>
    )
  }

  // Determine if this is a successful authentication
  const isSuccess = success === "true"

  return (
    <FramelessLayout>
      <div className="flex flex-col items-center space-y-6 max-w-md text-center">
        <Logo size={72} />

        <div className="space-y-2">
          <h1 className="text-xl font-medium text-neutral-900 dark:text-white">
            {displayError
              ? t("auth.oauth.titles.authenticationError")
              : isSuccess
                ? t("auth.oauth.titles.authenticationSuccessful")
                : t("auth.oauth.mcp.connectingWith", { provider: provider ? provider.charAt(0).toUpperCase() + provider.slice(1) : "Provider" })}
          </h1>
          <p className="text-sm text-neutral-500 dark:text-gray-400">
            {displayError
              ? t("auth.oauth.mcp.errorOccurred")
              : isSuccess
                ? t("auth.oauth.mcp.connectionSuccessful")
                : url
                  ? t("auth.oauth.mcp.redirectWait")
                  : t("auth.oauth.mcp.processingCallback")}
          </p>
        </div>

        {displayError
          ? (
              <div className="w-full space-y-4">
                <div className="text-sm text-red-500 bg-red-50 dark:bg-red-900/20 p-3 rounded-lg border border-red-200 dark:border-red-800">
                  {displayError}
                </div>
                {isPopup && (
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handleClosePopup}
                      className="flex-1"
                    >
                      {t("auth.oauth.mcp.close")}
                    </Button>
                    <Button
                      variant="default"
                      size="sm"
                      onClick={handleRetry}
                      className="flex-1"
                    >
                      {t("auth.oauth.mcp.tryAgain")}
                    </Button>
                  </div>
                )}
              </div>
            )
          : isSuccess
            ? (
                <div className="text-sm text-green-500 bg-green-50 dark:bg-green-900/20 p-3 rounded-lg border border-green-200 dark:border-green-800">
                  {t("auth.oauth.mcp.successMessage")}
                </div>
              )
            : processing
              ? (
                  <div className="w-8 h-8 border-2 border-neutral-200 dark:border-white/10 border-t-neutral-900 dark:border-t-white rounded-full animate-spin" />
                )
              : null}

        {isSuccess && isPopup && (
          <p className="text-xs text-neutral-400">
            {t("auth.oauth.mcp.autoCloseMessage")}
          </p>
        )}
      </div>
    </FramelessLayout>
  )
}
