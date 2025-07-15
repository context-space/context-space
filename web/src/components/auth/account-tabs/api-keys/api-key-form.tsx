import type { CreateAPIKeyRequest } from "@/typings/api-keys"
import { Plus } from "lucide-react"

import { useTranslations } from "next-intl"
import { useState } from "react"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

import { Textarea } from "@/components/ui/textarea"

interface ApiKeyFormProps {
  onSubmit: (request: CreateAPIKeyRequest) => Promise<void>
  isCreating: boolean
  isAtLimit: boolean
  apiKeyLimit: number
  currentCount: number
}

export function ApiKeyForm({ onSubmit, isCreating, isAtLimit, apiKeyLimit, currentCount }: ApiKeyFormProps) {
  const t = useTranslations()
  const [open, setOpen] = useState(false)
  const [name, setName] = useState("")
  const [description, setDescription] = useState("")

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!name.trim() || isAtLimit) return

    try {
      const request: CreateAPIKeyRequest = {
        name: name.trim(),
        description: description.trim() ?? "",
      }

      await onSubmit(request)

      // Show success message
      toast.success(t("account.createApiKeySuccess"))

      // Reset form and close dialog
      setName("")
      setDescription("")
      setOpen(false)
    } catch (error: any) {
      // Handle specific API key limit error
      if (error?.name === "API_KEY_LIMIT_REACHED") {
        toast.error(t("account.apiKeyLimitError", { limit: apiKeyLimit }))
      } else {
        // Handle other errors
        const errorMessage = error?.message || t("account.createApiKeyError")
        toast.error(errorMessage)
      }
    }
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!isCreating && !isAtLimit) {
      setOpen(newOpen)
      if (!newOpen) {
        // Reset form when closing
        setName("")
        setDescription("")
      }
    }
  }

  return (
    <div className="space-y-2">
      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogTrigger asChild>
          <Button disabled={isCreating || isAtLimit}>
            <Plus className="h-4 w-4 mr-2" />
            {t("account.createApiKey")}
          </Button>
        </DialogTrigger>
        <DialogContent>
          <form onSubmit={handleSubmit}>
            <DialogHeader>
              <DialogTitle>{t("account.createApiKey")}</DialogTitle>
              <DialogDescription>
                {t("account.createApiKeyDialogDescription")}
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="key-name" className="required">
                  {t("account.apiKeyName")}
                </Label>
                <Input
                  id="key-name"
                  placeholder={t("account.apiKeyNamePlaceholder")}
                  value={name}
                  onChange={e => setName(e.target.value)}
                  required
                  disabled={isCreating}
                  aria-describedby="key-name-description"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="key-description">{t("account.apiKeyDescription")}</Label>
                <Textarea
                  id="key-description"
                  placeholder={t("account.apiKeyDescriptionPlaceholder")}
                  value={description}
                  onChange={e => setDescription(e.target.value)}
                  rows={3}
                  disabled={isCreating}
                  aria-describedby="key-description-help"
                />
                <p id="key-description-help" className="text-xs text-muted-foreground">
                  {t("account.apiKeyDescriptionHelp")}
                </p>
              </div>
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => handleOpenChange(false)}
                disabled={isCreating}
              >
                {t("common.cancel")}
              </Button>
              <Button
                type="submit"
                disabled={isCreating || !name.trim() || isAtLimit}
              >
                {isCreating ? t("account.creating") : t("account.createApiKey")}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Limit information */}
      <div className="space-y-2">
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          {isAtLimit
            ? (
                <span className="text-amber-600 dark:text-amber-400 font-medium">
                  {t("account.apiKeyLimitReached", { count: currentCount, limit: apiKeyLimit })}
                </span>
              )
            : (
                <span>
                  {t("account.apiKeyLimit", { limit: apiKeyLimit })}
                  {" "}
                  (
                  {currentCount}
                  /
                  {apiKeyLimit}
                  )
                </span>
              )}
        </div>
      </div>
    </div>
  )
}
