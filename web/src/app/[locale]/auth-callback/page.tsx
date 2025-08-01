"use client"

import type { AuthState } from "@/components/auth"
import { useTranslations } from "next-intl"
import { useRouter, useSearchParams } from "next/navigation"
import { useCallback, useEffect, useState } from "react"
import { AuthCallbackLayout } from "@/components/auth"
import { createClient } from "@/lib/supabase/client"
import { clientLogger } from "@/lib/utils"

export default function AuthCallbackPage() {
  const t = useTranslations()
  const router = useRouter()
  const searchParams = useSearchParams()
  const [authState, setAuthState] = useState<AuthState>({
    status: "loading",
    message: t("auth.oauth.processing"),
  })

  const redirectTo = searchParams.get("redirect_to") ?? "/"

  // Handle authentication
  useEffect(() => {
    const supabase = createClient()

    // Check for URL errors first
    const error = searchParams.get("error")
    const errorCode = searchParams.get("error_code")
    const errorDescription = searchParams.get("error_description")

    if (error) {
      let errorMessage = ""
      let details = ""

      switch (errorCode || error) {
        case "otp_expired":
        case "expired":
          errorMessage = t("auth.oauth.magicLink.expired")
          details = t("auth.oauth.magicLink.expiredDetails")
          break
        case "otp_disabled":
          errorMessage = t("auth.oauth.magicLink.disabled")
          details = t("auth.oauth.magicLink.disabledDetails")
          break
        case "invalid_request":
          errorMessage = t("auth.oauth.magicLink.invalidRequest")
          details = t("auth.oauth.magicLink.invalidRequestDetails")
          break
        case "over_email_send_rate_limit":
          errorMessage = t("auth.oauth.magicLink.rateLimited")
          details = t("auth.oauth.magicLink.rateLimitedDetails")
          break
        default:
          errorMessage = t("auth.oauth.authenticationFailed")
          details = errorDescription || t("auth.oauth.magicLink.genericDetails")
      }

      setAuthState({
        status: "error",
        message: errorMessage,
        details,
      })
      return
    }

    // Handle magic link verification
    const handleAuth = async () => {
      try {
        const tokenHash = searchParams.get("token_hash")
        const type = searchParams.get("type")

        if (tokenHash && type === "email") {
          // Verify magic link
          const { data, error } = await supabase.auth.verifyOtp({
            token_hash: tokenHash,
            type: "email",
          })

          if (error) {
            clientLogger.error("Magic link verification failed", error)
            setAuthState({
              status: "error",
              message: t("auth.oauth.authenticationFailed"),
              details: t("auth.oauth.magicLink.verificationFailed"),
            })
          } else if (data.user) {
            setAuthState({
              status: "success",
              message: t("auth.oauth.loginSuccess"),
            })
          } else {
            setAuthState({
              status: "error",
              message: t("auth.oauth.authenticationFailed"),
              details: t("auth.oauth.magicLink.verificationFailed"),
            })
          }
        } else {
          // Check existing session for OAuth callbacks
          const { data: { session } } = await supabase.auth.getSession()

          if (session) {
            setAuthState({
              status: "success",
              message: t("auth.oauth.loginSuccess"),
            })
          } else {
            setAuthState({
              status: "error",
              message: t("auth.oauth.authenticationFailed"),
              details: t("auth.oauth.magicLink.verificationFailed"),
            })
          }
        }
      } catch (err) {
        clientLogger.error("Auth callback error", err)
        setAuthState({
          status: "error",
          message: t("auth.oauth.authenticationFailed"),
          details: t("auth.oauth.magicLink.verificationFailed"),
        })
      }
    }

    // Listen for auth state changes
    const { data: { subscription } } = supabase.auth.onAuthStateChange((event, session) => {
      if (session && authState.status === "loading") {
        setAuthState({
          status: "success",
          message: t("auth.oauth.loginSuccess"),
        })
      }
    })

    handleAuth()

    return () => {
      subscription.unsubscribe()
    }
  }, [searchParams, t, authState.status])

  const handleRetryLogin = useCallback(() => {
    const loginPath = redirectTo !== "/"
      ? `/login?from=${encodeURIComponent(redirectTo)}`
      : "/login"
    router.push(loginPath)
  }, [router, redirectTo])

  const handleCancel = useCallback(() => {
    router.push(redirectTo)
  }, [router, redirectTo])

  const getTitle = useCallback(() => {
    switch (authState.status) {
      case "error":
        return t("auth.oauth.titles.authenticationError")
      case "success":
        return t("auth.oauth.titles.authenticationSuccessful")
      default:
        return t("auth.oauth.titles.processingAuthentication")
    }
  }, [authState.status, t])

  const getSubtitle = useCallback(() => {
    switch (authState.status) {
      case "loading":
        return t("auth.oauth.processing")
      case "success":
        return t("auth.oauth.redirecting")
      case "error":
        return t("auth.oauth.tryAgain")
      default:
        return ""
    }
  }, [authState.status, t])

  return (
    <AuthCallbackLayout
      authState={authState}
      redirectTo={redirectTo}
      onGetTitle={getTitle}
      onGetSubtitle={getSubtitle}
      onRetry={handleRetryLogin}
      onCancel={handleCancel}
      autoRedirectDelay={1000}
      retryText={t("auth.oauth.retryLogin")}
      cancelText={t("auth.oauth.cancelLogin")}
    />
  )
}
