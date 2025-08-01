"use client"

import type { ReactNode } from "react"
import { useTranslations } from "next-intl"
import { useCallback, useEffect } from "react"
import { FramelessLayout } from "@/components/layouts"
import { Button } from "@/components/ui/button"
import { siteName } from "@/config"
import { LogoAvatar } from "../common/avatar"

export type AuthStatus = "loading" | "success" | "error"

export interface AuthState {
  status: AuthStatus
  message: string
  details?: string
}

interface AuthCallbackLayoutProps {
  authState: AuthState
  redirectTo: string
  onGetTitle: () => string
  onGetSubtitle: () => string
  onRetry?: () => void
  onCancel?: () => void
  autoRedirectDelay?: number
  children?: ReactNode
  retryText: string
  cancelText: string
}

interface StatusMessageProps {
  status: AuthStatus
  message: string
  details?: string
  variant?: "default" | "detailed"
}

interface ActionButtonsProps {
  onRetry?: () => void
  onCancel?: () => void
  retryText: string
  cancelText: string
}

interface SuccessActionsProps {
  onContinueNow: () => void
  continueText: string
  redirectText: string
}

interface LoadingStateProps {
  message: string
}

// Status message component
function StatusMessage({ status, message, details, variant = "default" }: StatusMessageProps) {
  const baseClasses = "text-sm p-3 rounded-lg border"

  const getStatusClasses = () => {
    switch (status) {
      case "error":
        return "text-red-500 bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800"
      case "success":
        return "text-green-600 dark:text-green-400 bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800"
      default:
        return "text-neutral-500 bg-neutral-50 dark:bg-neutral-900/20 border-neutral-200 dark:border-neutral-800"
    }
  }

  if (variant === "detailed" && details) {
    return (
      <div className={`${baseClasses} ${getStatusClasses()} space-y-3`}>
        <div className="font-medium">{message}</div>
        <div className="text-xs text-red-600 dark:text-red-400 leading-relaxed">
          {details}
        </div>
      </div>
    )
  }

  return (
    <div className={`${baseClasses} ${getStatusClasses()}`}>
      {details || message}
    </div>
  )
}

// Action buttons component
function ActionButtons({ onRetry, onCancel, retryText, cancelText }: ActionButtonsProps) {
  return (
    <div className="flex gap-2">
      {onCancel && (
        <Button
          variant="outline"
          size="sm"
          onClick={onCancel}
          className="flex-1"
        >
          {cancelText}
        </Button>
      )}
      {onRetry && (
        <Button
          variant="default"
          size="sm"
          onClick={onRetry}
          className={onCancel ? "flex-1" : "w-full"}
        >
          {retryText}
        </Button>
      )}
    </div>
  )
}

// Success actions component
function SuccessActions({ onContinueNow, continueText, redirectText }: SuccessActionsProps) {
  return (
    <>
      <div className="text-xs text-neutral-400 dark:text-gray-500">
        {redirectText}
      </div>
      <Button
        variant="default"
        size="sm"
        onClick={onContinueNow}
        className="w-full"
      >
        {continueText}
      </Button>
    </>
  )
}

// Loading state component
function LoadingState({ message }: LoadingStateProps) {
  return (
    <div className="w-full">
      <div className="flex items-center justify-center space-x-2 text-sm text-neutral-500 dark:text-gray-400">
        <div className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
        <span>{message}</span>
      </div>
    </div>
  )
}

export default function AuthCallbackLayout({
  authState,
  redirectTo,
  onGetTitle,
  onGetSubtitle,
  onRetry,
  onCancel,
  autoRedirectDelay = 1000,
  children,
  retryText,
  cancelText,
}: AuthCallbackLayoutProps) {
  const t = useTranslations()

  // Auto-redirect on success
  useEffect(() => {
    if (authState.status === "success") {
      const timeoutId = setTimeout(() => {
        window.location.href = redirectTo
      }, autoRedirectDelay)
      return () => clearTimeout(timeoutId)
    }
  }, [authState.status, redirectTo, autoRedirectDelay])

  const handleContinueNow = useCallback(() => {
    window.location.href = redirectTo
  }, [redirectTo])

  return (
    <FramelessLayout>
      <div className="flex flex-col items-center space-y-6 max-w-md text-center">
        <LogoAvatar alt={siteName} className="size-24" />

        <div className="space-y-2">
          <h1 className="text-xl font-medium text-neutral-900 dark:text-white">
            {onGetTitle()}
          </h1>
          <p className="text-sm text-neutral-500 dark:text-gray-400">
            {onGetSubtitle()}
          </p>
        </div>

        {/* Custom children content */}
        {children}

        {/* Error State */}
        {authState.status === "error" && (
          <div className="w-full space-y-4">
            <StatusMessage
              status="error"
              message={authState.message}
              details={authState.details}
              variant={authState.details ? "detailed" : "default"}
            />
            <ActionButtons
              onRetry={onRetry}
              onCancel={onCancel}
              retryText={retryText}
              cancelText={cancelText}
            />
          </div>
        )}

        {/* Success State */}
        {authState.status === "success" && (
          <div className="w-full space-y-4">
            <StatusMessage
              status="success"
              message={authState.message}
            />
            <SuccessActions
              onContinueNow={handleContinueNow}
              continueText={t("auth.oauth.continueNow")}
              redirectText={t("auth.oauth.redirectingAutomatically")}
            />
          </div>
        )}

        {/* Loading State */}
        {authState.status === "loading" && (
          <LoadingState message={authState.message} />
        )}
      </div>
    </FramelessLayout>
  )
}

// Export sub-components for custom usage
export {
  ActionButtons,
  LoadingState,
  StatusMessage,
  SuccessActions,
}
