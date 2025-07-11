"use client"

import { useTranslations } from "next-intl"
import { useRouter, useSearchParams } from "next/navigation"
import { useCallback, useEffect, useState } from "react"
import Logo from "@/components/common/logo"
import { FramelessLayout } from "@/components/layouts"
import { Button } from "@/components/ui/button"
import { createClient } from "@/lib/supabase/client"
import { clientLogger } from "@/lib/utils"

type AuthStatus = "loading" | "success" | "error"

interface AuthState {
  status: AuthStatus
  message: string
  details?: string
}

export default function AuthCallbackPage() {
  const t = useTranslations()
  const router = useRouter()
  const searchParams = useSearchParams()
  const [authState, setAuthState] = useState<AuthState>({
    status: "loading",
    message: t("auth.oauth.processing"),
  })

  const redirectTo = searchParams.get("redirect_to") || searchParams.get("from") || "/"

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
  }, [searchParams, t])

  // Auto-redirect on success
  useEffect(() => {
    if (authState.status === "success") {
      const timeoutId = setTimeout(() => {
        router.push(redirectTo)
      }, 1500)

      return () => clearTimeout(timeoutId)
    }
  }, [authState.status, router, redirectTo])

  const handleRetryLogin = useCallback(() => {
    const loginPath = redirectTo !== "/"
      ? `/login?from=${encodeURIComponent(redirectTo)}`
      : "/login"
    router.push(loginPath)
  }, [router, redirectTo])

  const handleGoHome = useCallback(() => {
    router.push("/")
  }, [router])

  const handleContinueNow = useCallback(() => {
    router.push(redirectTo)
  }, [router, redirectTo])

  const getTitle = () => {
    switch (authState.status) {
      case "error":
        return t("auth.oauth.titles.authenticationError")
      case "success":
        return t("auth.oauth.titles.authenticationSuccessful")
      default:
        return t("auth.oauth.titles.processingAuthentication")
    }
  }

  const getSubtitle = () => {
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
  }

  return (
    <FramelessLayout>
      <div className="flex flex-col items-center space-y-6 max-w-md text-center">
        <Logo size={72} />

        <div className="space-y-2">
          <h1 className="text-xl font-medium text-neutral-900 dark:text-white">
            {getTitle()}
          </h1>
          <p className="text-sm text-neutral-500 dark:text-gray-400">
            {getSubtitle()}
          </p>
        </div>

        {/* Error State */}
        {authState.status === "error" && (
          <div className="w-full space-y-4">
            <div className="text-sm text-red-500 bg-red-50 dark:bg-red-900/20 p-4 rounded-lg border border-red-200 dark:border-red-800 space-y-3">
              <div className="font-medium">{authState.message}</div>
              {authState.details && (
                <div className="text-xs text-red-600 dark:text-red-400 leading-relaxed">
                  {authState.details}
                </div>
              )}
            </div>

            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={handleGoHome}
                className="flex-1"
              >
                {t("common.goHome")}
              </Button>
              <Button
                variant="default"
                size="sm"
                onClick={handleRetryLogin}
                className="flex-1"
              >
                {t("auth.oauth.retryLogin")}
              </Button>
            </div>
          </div>
        )}

        {/* Success State */}
        {authState.status === "success" && (
          <div className="w-full space-y-4">
            <div className="text-sm text-green-600 dark:text-green-400 bg-green-50 dark:bg-green-900/20 p-3 rounded-lg border border-green-200 dark:border-green-800">
              {authState.message}
            </div>
            <div className="text-xs text-neutral-400 dark:text-gray-500">
              {t("auth.oauth.redirectingAutomatically")}
            </div>
            <Button
              variant="default"
              size="sm"
              onClick={handleContinueNow}
              className="w-full"
            >
              {t("auth.oauth.continueNow")}
            </Button>
          </div>
        )}

        {/* Loading State */}
        {authState.status === "loading" && (
          <div className="w-full">
            <div className="flex items-center justify-center space-x-2 text-sm text-neutral-500 dark:text-gray-400">
              <div className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
              <span>{authState.message}</span>
            </div>
          </div>
        )}
      </div>
    </FramelessLayout>
  )
}
