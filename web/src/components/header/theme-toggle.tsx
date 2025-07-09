"use client"

import { Check, Monitor, Moon, Sun } from "lucide-react"
import { useTranslations } from "next-intl"
import { useTheme } from "next-themes"
import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { cn } from "@/lib/utils"

const THEME_CONFIG = {
  light: {
    icon: Sun,
    sunScale: "rotate-0 scale-100",
    moonScale: "rotate-90 scale-0",
    monitorScale: "rotate-0 scale-0",
  },
  dark: {
    icon: Moon,
    sunScale: "-rotate-90 scale-0",
    moonScale: "rotate-0 scale-100",
    monitorScale: "rotate-90 scale-0",
  },
  system: {
    icon: Monitor,
    sunScale: "rotate-0 scale-0",
    moonScale: "rotate-90 scale-0",
    monitorScale: "rotate-0 scale-100",
  },
} as const

const THEME_OPTIONS = [
  { value: "light", icon: Sun, key: "light" },
  { value: "dark", icon: Moon, key: "dark" },
  { value: "system", icon: Monitor, key: "system" },
] as const

export function ThemeToggle() {
  const { theme, setTheme } = useTheme()
  const t = useTranslations("theme")
  const [mounted, setMounted] = useState(false)

  // useEffect only runs on the client, so now we can safely show the UI
  useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted) {
    return (
      <Button variant="ghost" size="sm" className="h-8 w-8 px-0">
        <Sun className="h-4 w-4 opacity-50" />
        <span className="sr-only">Toggle theme</span>
      </Button>
    )
  }

  const currentConfig = THEME_CONFIG[theme as keyof typeof THEME_CONFIG] || THEME_CONFIG.light

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className="h-8 w-8 px-0">
          <Sun className={cn("h-4 w-4 transition-all", currentConfig.sunScale)} />
          <Moon className={cn("absolute h-4 w-4 transition-all", currentConfig.moonScale)} />
          <Monitor className={cn("absolute h-4 w-4 transition-all", currentConfig.monitorScale)} />
          <span className="sr-only">{t("toggle")}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        align="end"
        className={cn(
          "w-30 space-y-1",
          "bg-white",
          "dark:bg-white/[0.02] dark:border-white/10 backdrop-blur-xl",
          "shadow-none",
          "border border-primary/15",
        )}
      >
        {THEME_OPTIONS.map(({ value, icon: Icon, key }) => (
          <DropdownMenuItem
            key={value}
            onClick={() => setTheme(value)}
            className={cn(
              "flex items-center justify-between px-3",
              theme === value && "bg-accent text-accent-foreground",
            )}
          >
            <div className="flex items-center gap-2.5">
              <Icon className="h-4 w-4" />
              <span>{t(key)}</span>
            </div>
            {theme === value && <Check className="h-3.5 w-3.5 ml-2" />}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
