"use client"

import { ChevronDown } from "lucide-react"
import { useTranslations } from "next-intl"
import { useRouter } from "next/navigation"
import { useCallback, useEffect, useState } from "react"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { connectAPIKey, connectOAuth, disconnectCredential } from "@/lib/client/credentials"
import { clientLogger, cn } from "@/lib/utils"
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
  const [isOpen, setIsOpen] = useState(true)
  const [connected, setConnected] = useState(isConnected)
  const router = useRouter()
  const t = useTranslations()

  useEffect(() => {
    setConnected(isConnected)
    // setIsOpen(!initialIsConnected && auth_type !== "none")
  }, [isConnected])

  const handleOAuthConnect = useCallback(async (selectedPermissions: string[]) => {
    try {
      await connectOAuth(providerId, selectedPermissions)
      setConnected(true)
      // Use router.refresh() instead of window.location.reload() to avoid flash
      router.refresh()
    } catch (error) {
      integrationConnectLogger.error("Failed to connect with OAuth", { error })
      setConnected(false)
    }
  }, [providerId, router])

  const handleApiKeyConnect = useCallback(async (apiKey: string) => {
    try {
      const response = await connectAPIKey(providerId, apiKey)
      if (response.success && response.data.is_valid) {
        setConnected(true)
        // Use router.refresh() instead of window.location.reload() to avoid flash
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
        await disconnectCredential(credentialId)
        setConnected(false)
        // Use router.refresh() instead of window.location.reload() to avoid flash
        router.refresh()
      }
    } catch (error) {
      integrationConnectLogger.error("Failed to disconnect", { error })
      setConnected(false)
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
                    <Badge variant={connected ? "default" : "secondary"} className="bg-green-50 dark:bg-green-500/10 text-green-700 dark:text-green-400 border-green-100 dark:border-green-500/20">
                      {connected ? t("integrations.connect.connected") : t("integrations.connect.noConnectionRequired")}
                    </Badge>
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
