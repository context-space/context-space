"use client"

import { ChevronDown } from "lucide-react"
import { useLocale, useTranslations } from "next-intl"
import { useCallback, useEffect, useState } from "react"
import { toast } from "sonner"
import { Status } from "@/components/integrations"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { getPathname, useRouter } from "@/i18n/navigation"
import { clientLogger, cn } from "@/lib/utils"
import { credentialsService } from "@/services/credentials"
import { ApiKeyConnect } from "./api-key-connect"
import { OAuthConnect } from "./oauth-connect"

const integrationConnectLogger = clientLogger.withTag("integration-connect")

export interface Permission {
  identifier: string
  name: string
  description: string
}

interface ConnectSectionProps {
  authType: "oauth" | "apikey" | "none"
  permissions: Permission[]
  isConnected: boolean
  providerId: string
  apiDocUrl?: string
  credentialId?: string
  authorizedPermissions?: string[]
}

export function Connect({
  authType,
  permissions,
  isConnected,
  providerId,
  apiDocUrl,
  credentialId,
  authorizedPermissions,
}: ConnectSectionProps) {
  const [isOpen, setIsOpen] = useState(!isConnected && authType !== "none")
  const [connected, setConnected] = useState(isConnected)
  const router = useRouter()
  const locale = useLocale()
  const t = useTranslations()

  useEffect(() => {
    setConnected(isConnected)
    if (isConnected && authType !== "none") {
      setIsOpen(false)
    }
  }, [isConnected, authType])

  const handleOAuthConnect = useCallback(async (selectedPermissions: string[]) => {
    try {
      const redirectUrl = window.location.origin + getPathname({
        href: {
          pathname: "/provider-callback",
          query: {
            redirect_to: window.location.href,
          },
        },
        locale,
      })
      const url = await credentialsService.connectOAuth(providerId, selectedPermissions, redirectUrl)
      if (!url) {
        toast.error(t("integrations.connect.oauth.failedToGetOAuthUrl"))
        throw new Error(t("integrations.connect.oauth.failedToGetOAuthUrl"))
      }
      window.location.href = url
      return true
    } catch (error) {
      integrationConnectLogger.error("Failed to connect with OAuth", { error })
      return false
    }
  }, [providerId, locale, t])

  const handleApiKeyConnect = useCallback(async (apiKey: string) => {
    try {
      const response = await credentialsService.connectAPIKey(providerId, apiKey)
      if (response.success && response.data.is_valid) {
        setConnected(true)
        setIsOpen(false)
        router.refresh()
      } else {
        throw new Error("Invalid API key")
      }
    } catch (error) {
      integrationConnectLogger.error("Failed to connect with API key", { error })
      setConnected(false)
    }
  }, [providerId, router])

  const handleDisconnect = useCallback(async () => {
    try {
      if (credentialId) {
        await credentialsService.disconnectCredential(credentialId)
        setConnected(false)
        setIsOpen(true)
        router.refresh()
      }
    } catch (error) {
      integrationConnectLogger.error("Failed to disconnect", { error })
      setConnected(false)
      setIsOpen(true)
    }
  }, [credentialId, router])

  return (
    <Card data-testid="connect-section" className="pb-0 bg-white/60 dark:bg-white/[0.02] border-base backdrop-blur-sm shadow-none">
      <Collapsible open={isOpen} onOpenChange={setIsOpen}>
        <CardHeader onClick={() => setIsOpen(!isOpen)}>
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <CardTitle className="text-lg font-semibold flex items-center justify-between">
                {t("integrations.connect.title")}
                <div className="flex items-center gap-2">
                  {(authType === "none" || connected) && (
                    <Status status={connected ? "connected" : "free"} type="badge" />
                  )}
                  <CollapsibleTrigger asChild>
                    <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                      <ChevronDown className={cn("h-4 w-4 transition-transform duration-200", isOpen && "rotate-180")} />
                      <span className="sr-only">{t("common.toggle")}</span>
                    </Button>
                  </CollapsibleTrigger>
                </div>
              </CardTitle>
              <div className="text-sm text-muted-foreground my-3">
                {authType === "oauth" && (
                  t("integrations.connect.oauth.description")
                )}
                {authType === "apikey" && (
                  t("integrations.connect.apiKey.description")
                )}
                {authType === "none" && (
                  t("integrations.connect.noConnectionRequiredDesc")
                )}
              </div>
            </div>
          </div>
        </CardHeader>
        {(authType !== "none") && (
          <CollapsibleContent>
            <CardContent className="pt-0 px-0">
              <div className="border-t border-base p-5 space-y-4">

                {authType === "oauth" && (
                  <OAuthConnect
                    permissions={permissions}
                    isConnected={connected}
                    onConnect={handleOAuthConnect}
                    onDisconnect={handleDisconnect}
                    authorizedPermissions={authorizedPermissions ?? []}
                  />
                )}

                {authType === "apikey" && (
                  <ApiKeyConnect
                    isConnected={connected}
                    helpLink={apiDocUrl}
                    onConnect={handleApiKeyConnect}
                    onDisconnect={handleDisconnect}
                  />
                )}
              </div>
            </CardContent>
          </CollapsibleContent>
        )}
      </Collapsible>
    </Card>
  )
}
