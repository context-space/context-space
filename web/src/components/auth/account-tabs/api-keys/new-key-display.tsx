import { Check, Copy } from "lucide-react"
import { useTranslations } from "next-intl"
import { useCallback, useState } from "react"

import { Button } from "@/components/ui/button"
import { useCopyToClipboard } from "@/hooks/use-clipboard"

interface NewKeyDisplayProps {
  keyValue: string
  onDismiss: () => void
}

export function NewKeyDisplay({ keyValue, onDismiss }: NewKeyDisplayProps) {
  const t = useTranslations("account")
  const [isCopied, setIsCopied] = useState(false)
  const copyToClipboard = useCopyToClipboard()

  const handleCopy = useCallback(async () => {
    const success = await copyToClipboard(keyValue, t("apiKeyCopied"), t("copyApiKeyFailed"))
    if (success) {
      setIsCopied(true)
      setTimeout(() => setIsCopied(false), 2000)
    }
  }, [copyToClipboard, keyValue, t])

  return (
    <div
      className="p-4 border rounded-lg bg-green-50 dark:bg-green-950/20 border-green-200 dark:border-green-800"
      role="alert"
      aria-live="polite"
    >
      <h4 className="font-medium text-green-800 dark:text-green-200 mb-2">
        {t("newApiKeyCreated")}
      </h4>
      <p className="text-sm text-green-700 dark:text-green-300 mb-3">
        {t("newApiKeyCopyWarning")}
      </p>
      <div className="flex items-center gap-2 mb-3">
        <code
          className="flex-1 p-2 bg-white/20 dark:bg-gray-900/20 border rounded-md text-sm font-mono break-all"
          aria-label="Your new API key"
        >
          {keyValue}
        </code>
      </div>
      <div className="flex items-center gap-2 mb-3">

        <Button
          size="sm"
          variant="outline"
          onClick={handleCopy}
          aria-label="Copy API key to clipboard"
          className="h-8"
        >
          {isCopied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
        </Button>
        <Button
          size="sm"
          variant="outline"
          onClick={onDismiss}
        >
          {t("newApiKeySaved")}
        </Button>
      </div>
    </div>
  )
}
