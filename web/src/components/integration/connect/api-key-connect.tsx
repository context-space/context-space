"use client"

import { Eye, EyeOff, ShieldCheck } from "lucide-react"
import { useTranslations } from "next-intl"
import { useCallback, useState } from "react"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useAuth } from "@/hooks/use-auth"
import { useRouter } from "@/i18n/navigation"
import { clientLogger } from "@/lib/utils"

const apiKeyConnectLogger = clientLogger.withTag("api-key-connect")
interface ApiKeyConnectProps {
  isConnected: boolean
  onConnect: (apiKey: string) => Promise<void>
  onDisconnect: () => Promise<void>
  helpLink?: string
}

export function ApiKeyConnect({ isConnected, onConnect, onDisconnect, helpLink = "#" }: ApiKeyConnectProps) {
  const [apiKey, setApiKey] = useState("")
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const t = useTranslations()
  const { isAuthenticated } = useAuth()
  const router = useRouter()

  const checkAuthenticationAndShowToast = useCallback(() => {
    if (!isAuthenticated) {
      toast.error(t("integrations.connect.loginRequired"), {
        description: t("integrations.connect.loginRequiredDescription"),
        action: {
          label: t("header.login"),
          onClick: () => {
            const currentPath = encodeURIComponent(window.location.pathname)
            router.push(`/login?from=${currentPath}`)
          },
        },
      })
      return false
    }
    return true
  }, [isAuthenticated, t])

  const handleSubmit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault()
    if (!apiKey.trim()) {
      return
    }

    if (!checkAuthenticationAndShowToast()) {
      return
    }

    try {
      setLoading(true)
      await onConnect(apiKey)
      setApiKey("")
    } catch (error) {
      apiKeyConnectLogger.error("Failed to connect API key", { error })
    } finally {
      setLoading(false)
    }
  }, [apiKey, onConnect, checkAuthenticationAndShowToast])

  const handleDisconnect = useCallback(async () => {
    if (!checkAuthenticationAndShowToast()) {
      return
    }

    try {
      setLoading(true)
      await onDisconnect()
    } catch (error) {
      apiKeyConnectLogger.error("Failed to disconnect", { error })
    } finally {
      setLoading(false)
    }
  }, [onDisconnect, checkAuthenticationAndShowToast])

  if (isConnected) {
    return (
      <div className="space-y-4">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <ShieldCheck className="h-4 w-4 text-green-500" />
          {t("integrations.connect.apiKey.secureNote")}
        </div>
        <Button
          className="w-full"
          variant="outline"
          onClick={handleDisconnect}
          disabled={loading}
        >
          {loading ? t("integrations.connect.apiKey.disconnecting") : t("integrations.connect.oauth.disconnect")}
        </Button>
      </div>
    )
  }

  return (
    <>
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-medium">{t("integrations.connect.apiKey.title")}</h3>
        <a href={helpLink} target="_blank" className="text-xs text-primary hover:text-primary/80">
          {t("integrations.connect.apiKey.howToGet")}
        </a>
      </div>

      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <ShieldCheck className="h-4 w-4 text-green-500" />
        {t("integrations.connect.apiKey.secureNote")}
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="relative">
          <Input
            type={showPassword ? "text" : "password"}
            value={apiKey}
            onChange={e => setApiKey(e.target.value)}
            placeholder={t("integrations.connect.apiKey.placeholder")}
            autoComplete="off"
            className="pr-10"
          />
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword
              ? (
                  <EyeOff className="h-4 w-4 text-muted-foreground" />
                )
              : (
                  <Eye className="h-4 w-4 text-muted-foreground" />
                )}
            <span className="sr-only">
              {showPassword ? t("integrations.connect.apiKey.hide") : t("integrations.connect.apiKey.show")}
              API key
            </span>
          </Button>
        </div>

        <Button
          type="submit"
          className="w-full"
          disabled={loading || !apiKey.trim()}
        >
          {loading ? t("integrations.connect.apiKey.connecting") : t("integrations.connect.apiKey.saveApiKey")}
        </Button>
      </form>
    </>
  )
}
