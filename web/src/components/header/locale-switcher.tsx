"use client"

import type { Locale } from "@/i18n/routing"
import { Check, Languages } from "lucide-react"
import { useLocale, useTranslations } from "next-intl"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { usePathname, useRouter } from "@/i18n/navigation"
import { cn } from "@/lib/utils"

const locales = [
  { code: "en", name: "English" },
  { code: "zh", name: "简体中文" },
  { code: "zh-TW", name: "繁體中文" },
] as { code: Locale, name: string }[]

export function LocaleSwitcher() {
  const locale = useLocale()
  const router = useRouter()
  const pathname = usePathname()
  const t = useTranslations("language")

  const handleLocaleChange = (newlocale: Locale) => {
    router.replace(pathname, { locale: newlocale })
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className="h-8 w-8 px-0">
          <Languages className="h-4 w-4" />
          <span className="sr-only">{t("toggle")}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        align="end"
        className={cn(
          "w-32 space-y-1",
          "bg-white",
          "dark:bg-white/[0.02] dark:border-white/10 dark:backdrop-blur-xl",
          "shadow-none",
          "border border-primary/15",
        )}
      >
        {locales.map(loc => (
          <DropdownMenuItem
            key={loc.code}
            onClick={() => handleLocaleChange(loc.code)}
            className={cn(
              "flex items-center justify-between px-3",
              locale === loc.code && "bg-accent text-accent-foreground",
            )}
          >
            <div className="flex items-center gap-2.5">
              <span>{loc.name}</span>
            </div>
            {locale === loc.code && <Check className="h-3.5 w-3.5 ml-2" />}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
