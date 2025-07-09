"use client"

import { useTranslations } from "next-intl"
import { useRouter, useSearchParams } from "next/navigation"
import { useCallback, useEffect, useMemo, useState } from "react"
import Logo from "@/components/common/logo"
import { FramelessLayout } from "@/components/layouts"
import { Button } from "@/components/ui/button"
import { createClient } from "@/lib/supabase/client"
import { clientLogger } from "@/lib/utils"

interface OAuthCallbackState {
  status: "loading" | "success" | "error"
  message: string
}

interface OAuthParams {
  error?: string | null
  success?: string | null
  redirectTo: string
}

// Constants
const REDIRECT_DELAY = 1500
const DEFAULT_REDIRECT = "/"

// Logger
const oauthCallbackLogger = clientLogger.withTag("oauth-callback")

// Custom hook for OAuth status management with auth state detection
function useOAuthStatus(): OAuthCallbackState {
  const t = useTranslations()
  const searchParams = useSearchParams()
  const [authState, setAuthState] = useState<"checking" | "authenticated" | "unauthenticated">("checking")

  // Monitor auth state changes
  useEffect(() => {
    const supabase = createClient()

    // Check current session
    const checkSession = async () => {
      try {
        const { data: { session } } = await supabase.auth.getSession()
        if (session) {
          oauthCallbackLogger.info("User session detected", { userId: session.user.id })
          setAuthState("authenticated")
        } else {
          setAuthState("unauthenticated")
        }
      } catch (error) {
        oauthCallbackLogger.error("Error checking session", { error })
        setAuthState("unauthenticated")
      }
    }

    // Listen for auth state changes
    const { data: { subscription } } = supabase.auth.onAuthStateChange((event, session) => {
      oauthCallbackLogger.info("Auth state changed", { event, hasSession: !!session })

      if (session) {
        setAuthState("authenticated")
      } else if (event === "SIGNED_OUT") {
        setAuthState("unauthenticated")
      }
    })

    checkSession()

    return () => subscription.unsubscribe()
  }, [])

  return useMemo((): OAuthCallbackState => {
    const error = searchParams.get("error")

    if (error) {
      oauthCallbackLogger.warn("OAuth error detected", { error })

      const errorMessage = (() => {
        switch (error) {
          case "Authentication failed":
            return t("auth.oauth.authenticationFailed")
          case "Server error occurred":
            return t("auth.oauth.serverError")
          default:
            return t("auth.oauth.unknownError")
        }
      })()

      return {
        status: "error",
        message: errorMessage,
      }
    }

    // Check auth state
    if (authState === "authenticated") {
      oauthCallbackLogger.info("OAuth success confirmed")
      return {
        status: "success",
        message: t("auth.oauth.loginSuccess"),
      }
    }

    if (authState === "unauthenticated") {
      return {
        status: "error",
        message: t("auth.oauth.authenticationFailed"),
      }
    }

    // Still checking
    return {
      status: "loading",
      message: t("auth.oauth.processing"),
    }
  }, [searchParams, t, authState])
}

// Custom hook for redirect parameters
function useOAuthParams(): OAuthParams {
  const searchParams = useSearchParams()

  return useMemo((): OAuthParams => {
    const redirectTo = searchParams.get("redirect_to") || searchParams.get("from") || DEFAULT_REDIRECT

    return {
      error: searchParams.get("error"),
      success: searchParams.get("success"),
      redirectTo,
    }
  }, [searchParams])
}

// Custom hook for navigation actions
function useOAuthNavigation() {
  const router = useRouter()
  const { redirectTo } = useOAuthParams()

  const redirectToDestination = useCallback(() => {
    oauthCallbackLogger.info("Redirecting user", { destination: redirectTo })
    router.push(redirectTo)
  }, [router, redirectTo])

  const redirectToLogin = useCallback(() => {
    const loginPath = redirectTo !== DEFAULT_REDIRECT
      ? `/login?from=${encodeURIComponent(redirectTo)}`
      : "/login"

    oauthCallbackLogger.info("Redirecting to login", { loginPath })
    router.push(loginPath)
  }, [router, redirectTo])

  const redirectToHome = useCallback(() => {
    oauthCallbackLogger.info("Redirecting to home")
    router.push(DEFAULT_REDIRECT)
  }, [router])

  return {
    redirectToDestination,
    redirectToLogin,
    redirectToHome,
  }
}

// Auto-redirect effect hook
function useAutoRedirect(status: OAuthCallbackState["status"], onRedirect: () => void) {
  // Auto-redirect on success
  useEffect(() => {
    if (status === "success") {
      const timeoutId = setTimeout(() => {
        onRedirect()
      }, REDIRECT_DELAY)

      return () => clearTimeout(timeoutId)
    }
  }, [status, onRedirect])
}

// Main component
export default function OAuthCallbackPage() {
  const t = useTranslations()
  const { status, message } = useOAuthStatus()
  const { redirectToDestination, redirectToLogin, redirectToHome } = useOAuthNavigation()

  // Auto-redirect on success
  useAutoRedirect(status, redirectToDestination)

  const handleManualRedirect = useCallback(() => {
    redirectToDestination()
  }, [redirectToDestination])

  const subtitleText = useMemo(() => {
    switch (status) {
      case "loading":
        return t("auth.oauth.processing")
      case "success":
        return t("auth.oauth.redirecting")
      case "error":
        return t("auth.oauth.tryAgain")
      default:
        return ""
    }
  }, [status, t])

  return (
    <FramelessLayout>
      <div className="flex flex-col items-center space-y-6 max-w-md text-center">
        <Logo size={72} />

        <div className="space-y-2">
          <h1 className="text-xl font-medium text-neutral-900 dark:text-white">
            {status === "error"
              ? t("auth.oauth.titles.authenticationError")
              : status === "success"
                ? t("auth.oauth.titles.authenticationSuccessful")
                : t("auth.oauth.titles.processingAuthentication")}
          </h1>
          <p className="text-sm text-neutral-500 dark:text-gray-400">
            {subtitleText}
          </p>
        </div>

        {/* Status display */}
        {status === "error" && (
          <div className="w-full space-y-4">
            <div className="text-sm text-red-500 bg-red-50 dark:bg-red-900/20 p-3 rounded-lg border border-red-200 dark:border-red-800">
              {message}
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={redirectToHome}
                className="flex-1"
              >
                {t("common.goHome")}
              </Button>
              <Button
                variant="default"
                size="sm"
                onClick={redirectToLogin}
                className="flex-1"
              >
                {t("auth.oauth.retryLogin")}
              </Button>
            </div>
          </div>
        )}

        {status === "success" && (
          <div className="w-full space-y-4">
            <div className="text-sm text-green-500 bg-green-50 dark:bg-green-900/20 p-3 rounded-lg border border-green-200 dark:border-green-800">
              {message}
            </div>
            <p className="text-sm text-neutral-500 dark:text-gray-400">
              {t("auth.oauth.noAutoRedirect")}
              <button
                type="button"
                onClick={handleManualRedirect}
                className="text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 underline cursor-pointer ml-1"
              >
                {t("auth.oauth.clickHere")}
              </button>
            </p>
            <Button
              variant="outline"
              size="sm"
              onClick={redirectToDestination}
              className="w-full"
            >
              {t("auth.oauth.continueNow")}
            </Button>
          </div>
        )}

        {status === "loading" && (
          <div className="w-8 h-8 border-2 border-neutral-200 dark:border-white/10 border-t-neutral-900 dark:border-t-white rounded-full animate-spin" />
        )}

        {status === "success" && (
          <p className="text-xs text-neutral-400">
            {t("auth.oauth.redirectingAutomatically")}
          </p>
        )}
      </div>
    </FramelessLayout>
  )
}
