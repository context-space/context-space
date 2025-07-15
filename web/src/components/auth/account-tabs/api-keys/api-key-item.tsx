import type { APIKey } from "@/typings/api-keys"
import { Key, Trash2 } from "lucide-react"
import { useTranslations } from "next-intl"
import { useCallback, useMemo } from "react"

import { toast } from "sonner"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"

interface ApiKeyItemProps {
  apiKey: APIKey
  onDelete: (keyId: string) => Promise<void>
  isDeleting: boolean
}

export function ApiKeyItem({ apiKey, onDelete, isDeleting }: ApiKeyItemProps) {
  const t = useTranslations()

  const handleDelete = useCallback(async () => {
    try {
      await onDelete(apiKey.id)
      toast.success(t("account.deleteApiKeySuccess"))
    } catch (error: any) {
      const errorMessage = error?.message || t("account.deleteApiKeyError")
      toast.error(errorMessage)
    }
  }, [apiKey.id, onDelete, t])

  const formatDate = useCallback((dateString: string) => {
    return new Date(dateString).toLocaleDateString()
  }, [])

  const apikey = useMemo(() => apiKey.key_value?.replace("cs-", "").slice(0, 8), [apiKey.key_value])

  return (
    <div
      className="p-4 border border-base rounded-lg hover:bg-muted/50 transition-colors"
      role="article"
      aria-label={`API key: ${apiKey.name}`}
    >
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <Key className="h-4 w-4 text-muted-foreground" aria-hidden="true" />
            <h4 className="font-medium">{apiKey.name}</h4>
            {apikey && (
              <Badge variant="secondary" className="text-xs">
                {apikey}
              </Badge>
            )}
          </div>
          {apiKey.description && (
            <p className="text-sm text-muted-foreground mb-2">
              {apiKey.description}
            </p>
          )}
          <p className="text-xs text-muted-foreground">
            {t("account.createdLabel")}
            {" "}
            <time dateTime={apiKey.created_at}>{formatDate(apiKey.created_at)}</time>
          </p>
        </div>
        <div className="flex items-center gap-2">
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button
                size="sm"
                variant="destructive"
                disabled={isDeleting}
                aria-label={`Delete API key ${apiKey.name}`}
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>{t("account.deleteApiKey")}</AlertDialogTitle>
                <AlertDialogDescription>
                  {t("account.deleteApiKeyConfirm")}
                  {" "}
                  {t("account.deleteApiKeyConfirmAction")}
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>{t("common.cancel")}</AlertDialogCancel>
                <AlertDialogAction
                  onClick={handleDelete}
                  className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                  disabled={isDeleting}
                >
                  {isDeleting ? t("common.deleting") : t("common.delete")}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </div>
    </div>
  )
}
