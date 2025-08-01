"use client"

import type { Permission } from "./index"
import { Info } from "lucide-react"
import { useTranslations } from "next-intl"
import { useCallback, useEffect, useMemo, useState } from "react"
import { toast } from "sonner"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { ScrollArea } from "@/components/ui/scroll-area"
import { useAuth } from "@/hooks/use-auth"
import { useRouter } from "@/i18n/navigation"
import { clientLogger } from "@/lib/utils"

const oauthConnectLogger = clientLogger.withTag("oauth-connect")

interface OAuthConnectProps {
  permissions: Permission[]
  isConnected: boolean
  onConnect: (permissions: string[]) => Promise<boolean>
  onDisconnect: () => Promise<void>
  authorizedPermissions: string[]
}

export function OAuthConnect({
  permissions,
  isConnected,
  onConnect,
  onDisconnect,
  authorizedPermissions,
}: OAuthConnectProps) {
  const [checkedPermissions, setCheckedPermissions] = useState<string[]>(permissions.map(p => p.identifier))
  const [waiting, setWaiting] = useState(false)
  const [loading, setLoading] = useState(false)
  const t = useTranslations()
  const { isAuthenticated } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (isConnected && authorizedPermissions.length > 0) {
      setCheckedPermissions(authorizedPermissions)
    } else {
      setCheckedPermissions(permissions.map(p => p.identifier))
    }
  }, [isConnected, authorizedPermissions, permissions])

  const isSamePermissions = useMemo(() =>
    checkedPermissions.length === authorizedPermissions.length
    && checkedPermissions.every(p => authorizedPermissions.includes(p)), [checkedPermissions, authorizedPermissions])

  const handleSelectAll = useCallback(() => {
    setCheckedPermissions(permissions.map(p => p.identifier))
  }, [permissions])

  const handleCheckboxChange = useCallback((identifier: string, checked: boolean) => {
    setCheckedPermissions(prev =>
      checked ? [...prev, identifier] : prev.filter(pid => pid !== identifier),
    )
  }, [])

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
  }, [isAuthenticated, t, router])

  const handleConnect = useCallback(async () => {
    if (!checkAuthenticationAndShowToast()) {
      return
    }

    setWaiting(true)
    const success = await onConnect(checkedPermissions)
    if (!success) setWaiting(false)
  }, [checkedPermissions, onConnect, checkAuthenticationAndShowToast])

  const handleDisconnect = useCallback(async () => {
    if (!checkAuthenticationAndShowToast()) {
      return
    }

    try {
      setLoading(true)
      await onDisconnect()
    } catch (error) {
      oauthConnectLogger.error("Failed to disconnect", { error })
    } finally {
      setLoading(false)
    }
  }, [onDisconnect, checkAuthenticationAndShowToast])

  const handleUpdate = useCallback(async () => {
    if (!checkAuthenticationAndShowToast()) {
      return
    }

    setWaiting(true)
    const success = await onConnect(checkedPermissions)
    if (!success) setWaiting(false)
  }, [checkedPermissions, onConnect, checkAuthenticationAndShowToast])

  const wantUpdate = useMemo(() => {
    return isConnected && !isSamePermissions && checkedPermissions.length > 0
  }, [isConnected, isSamePermissions, checkedPermissions])

  const wantConnect = useMemo(() => {
    return !isConnected && checkedPermissions.length > 0
  }, [isConnected, checkedPermissions])

  const wantUnconnect = useMemo(() => {
    return isConnected && (checkedPermissions.length === 0 || isSamePermissions)
  }, [isConnected, checkedPermissions, isSamePermissions])

  return (
    <>
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-sm font-medium">{t("integrations.connect.oauth.requiredPermissions")}</h3>
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">
              {permissions.length}
              {" "}
              {t("integrations.connect.oauth.permissions")}
            </span>
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="h-6 px-2 text-xs"
              onClick={handleSelectAll}
            >
              {t("integrations.connect.oauth.selectAll")}
            </Button>
          </div>
        </div>

        <ScrollArea className="pr-4" style={{ height: (permissions.length > 4 ? 4 : permissions.length) * 36 + 10 }}>
          <div className="space-y-2">
            {permissions.map(permission => (
              <div key={permission.identifier} className="flex items-start space-x-3">
                <Checkbox
                  id={permission.identifier}
                  checked={checkedPermissions.includes(permission.identifier)}
                  onCheckedChange={checked =>
                    handleCheckboxChange(permission.identifier, checked as boolean)}
                  className="mt-0.5 cursor-pointer border-base"
                />
                <div className="space-y-1 leading-none flex-1">
                  <div className="flex items-center gap-3 mb-1">
                    <label
                      htmlFor={permission.identifier}
                      className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                    >
                      {permission.name}
                    </label>
                    <div className="h-px flex-1 border-b border-base" />
                    <Badge
                      variant="outline"
                      className="text-xs"
                    >
                      {permission.identifier}
                    </Badge>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    {permission.description}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </ScrollArea>

        <div className="flex items-center gap-2 pt-3 border-t border-base">
          <Info className="h-4 w-4 text-muted-foreground flex-shrink-0" />
          <span className="text-xs text-muted-foreground">
            {t("integrations.connect.oauth.permissionsNote")}
          </span>
        </div>
      </div>

      <Button
        className="w-full"
        variant={wantUnconnect ? "outline" : "default"}
        onClick={wantUnconnect ? handleDisconnect : wantUpdate ? handleUpdate : wantConnect ? handleConnect : undefined}
        disabled={loading || waiting || (!wantConnect && !wantUpdate && !wantUnconnect)}
      >
        {loading
          ? (t("common.loading"))
          : waiting
            ? (t("integrations.connect.oauth.waiting"))
            : (wantUnconnect ? t("integrations.connect.oauth.disconnect") : wantUpdate ? t("integrations.connect.oauth.updatePermissions") : wantConnect ? t("integrations.connect.oauth.connectWithOAuth") : t("common.connect"))}
      </Button>
    </>
  )
}
