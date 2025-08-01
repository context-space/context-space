"use client"

import type { AuthState } from "@/components/auth"
import { useTranslations } from "next-intl"
import { useRouter, useSearchParams } from "next/navigation"
import { useCallback, useEffect, useState } from "react"
import { AuthCallbackLayout } from "@/components/auth"

export default function OAuthCallbackPage() {
  const t = useTranslations()
  const router = useRouter()
  const searchParams = useSearchParams()

  const [authState, setAuthState] = useState<AuthState>({
    status: "loading",
    message: t("auth.oauth.mcp.processingCallback"),
  })

  const success = searchParams.get("success")
  const message = searchParams.get("message") || "Connection failed"
  const code = searchParams.get("code")
  const redirectTo = searchParams.get("redirect_to") ?? "/integrations"

  // Determine auth state based on URL parameters
  useEffect(() => {
    // Handle successful authentication
    if (success === "true") {
      setAuthState({
        status: "success",
        message: t("auth.oauth.mcp.connectionSuccessful"),
      })
      return
    }

    // Determine error state
    const hasError = success && success !== "true"
    const displayError = hasError ? message : ""

    // Handle invalid callback - no valid parameters
    const isInvalidCallback = !success && !message && !code
    const finalError = isInvalidCallback ? t("auth.oauth.mcp.invalidCallback") : displayError

    if (finalError) {
      setAuthState({
        status: "error",
        message: t("auth.oauth.titles.authenticationError"),
        details: finalError,
      })
    }
  }, [success, message, code, t])

  const handleRetry = useCallback(() => {
    router.push(redirectTo)
  }, [router, redirectTo])

  const getTitle = useCallback(() => {
    switch (authState.status) {
      case "error":
        return t("auth.oauth.titles.authenticationError")
      case "success":
        return t("auth.oauth.titles.authenticationSuccessful")
      default:
        return t("auth.oauth.mcp.connectingWith", { provider: "Provider" })
    }
  }, [authState.status, t])

  const getSubtitle = useCallback(() => {
    switch (authState.status) {
      case "loading":
        return authState.message
      case "success":
        return t("auth.oauth.redirecting")
      case "error":
        return t("auth.oauth.mcp.errorOccurred")
      default:
        return ""
    }
  }, [authState.status, authState.message, t])

  const handleCancel = useCallback(() => {
    router.push(redirectTo)
  }, [router, redirectTo])

  return (
    <AuthCallbackLayout
      authState={authState}
      redirectTo={redirectTo}
      onGetTitle={getTitle}
      onGetSubtitle={getSubtitle}
      onRetry={handleRetry}
      onCancel={handleCancel}
      autoRedirectDelay={1000}
      retryText={t("auth.oauth.mcp.retryConnection")}
      cancelText={t("auth.oauth.cancelLogin")}
    />
  )
}
