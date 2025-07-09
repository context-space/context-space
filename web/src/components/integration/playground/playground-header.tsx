"use client"
import { Trash2 } from "lucide-react"
import { useTranslations } from "next-intl"
import { Button } from "@/components/ui/button"

interface PlaygroundHeaderProps {
  onClear: () => void
}

export function PlaygroundHeader({ onClear }: PlaygroundHeaderProps) {
  const t = useTranslations()

  return (
    <div>
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">{t("playground.title")}</h2>
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={onClear}
            className="h-8"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>
      <p className="text-sm text-muted-foreground my-3">
        {t("playground.description")}
      </p>
    </div>
  )
}
