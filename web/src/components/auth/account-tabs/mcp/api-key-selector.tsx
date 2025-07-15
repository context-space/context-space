import type { APIKey } from "@/typings/api-keys"
import { Key } from "lucide-react"

import { useTranslations } from "next-intl"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

interface ApiKeySelectorProps {
  apiKeys: APIKey[]
  selectedApiKeyId: string
  onSelect: (id: string) => void
}

export function ApiKeySelector({
  apiKeys,
  selectedApiKeyId,
  onSelect,
}: ApiKeySelectorProps) {
  const t = useTranslations("account")

  // 如果只有一个 API key，不显示选择框
  if (apiKeys.length <= 1) {
    return null
  }

  return (
    <div className="flex justify-between items-center">
      <Label htmlFor="api-key-select">{t("selectApiKey")}</Label>
      <Select value={selectedApiKeyId} onValueChange={onSelect}>
        <SelectTrigger id="api-key-select" className="border border-base">
          <SelectValue placeholder={t("selectApiKeyPlaceholder")} />
        </SelectTrigger>
        <SelectContent>
          {apiKeys.map(apiKey => (
            <SelectItem key={apiKey.id} value={apiKey.id}>
              <div className="flex items-center gap-2">
                <Key className="h-4 w-4" />
                <span>{apiKey.name}</span>
                <span className="text-xs text-muted-foreground">
                  (
                  {apiKey.id.slice(-8)}
                  )
                </span>
              </div>
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  )
}
