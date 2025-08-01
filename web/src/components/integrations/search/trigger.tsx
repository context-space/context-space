"use client"

import { useWindowSize } from "@uidotdev/usehooks"
import { useAtom } from "jotai"
import { SearchIcon } from "lucide-react"
import { useTranslations } from "next-intl"
import { useCallback, useEffect, useMemo } from "react"
import { cn } from "@/lib/utils"
import { searchOpenAtom } from "./search-atom"

function useIsMobile(): boolean {
  const { width } = useWindowSize()
  const isMobile = useMemo(() => !!width && width < 768, [width])

  return isMobile
}

function useKeyboardShortcuts(onToggle: () => void) {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        onToggle()
      }
    }

    document.addEventListener("keydown", handleKeyDown)
    return () => document.removeEventListener("keydown", handleKeyDown)
  }, [onToggle])
}

export function SearchTrigger() {
  const t = useTranslations()
  const [, setSearchOpen] = useAtom(searchOpenAtom)
  const isMobile = useIsMobile()

  const handleToggle = useCallback(() => {
    setSearchOpen(prev => !prev)
  }, [setSearchOpen])

  useKeyboardShortcuts(handleToggle)

  const placeholderText = useMemo(() => {
    return t("integrations.search.placeholder")
  }, [t])

  return (
    <div className="relative w-full">
      <button
        type="button"
        onClick={handleToggle}
        className={cn(
          "relative flex items-center w-full transition-all duration-300 ease-out",
          "bg-white dark:bg-white/[0.02] border border-neutral-200/60 dark:border-white/[0.05]",
          "rounded-xl text-left cursor-pointer",
          isMobile ? "h-12" : "h-14",
          "hover:ring-1 hover:ring-primary/20 dark:hover:ring-primary/30",
          "focus:outline-none focus:ring-2 focus:ring-primary/50",
        )}
        aria-label="Open search dialog"
      >
        <div className="absolute left-4 pointer-events-none">
          <SearchIcon
            size={isMobile ? 18 : 20}
            className="text-neutral-400 dark:text-gray-500"
            aria-hidden="true"
          />
        </div>

        <span
          className={cn(
            "w-full bg-transparent border-0 outline-none pointer-events-none",
            "text-neutral-400 dark:text-gray-500",
            isMobile ? "pl-11 pr-16 text-base" : "pl-12 pr-20 text-lg",
          )}
        >
          {placeholderText}
        </span>

        {!isMobile && (
          <div className="absolute right-4 flex items-center gap-1 text-xs text-neutral-400 dark:text-gray-500">
            <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-primary/10 px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
              <span className="text-xs">⌘</span>
              K
            </kbd>
          </div>
        )}
      </button>
    </div>
  )
}

export function LoadingMoreSearchTrigger() {
  const t = useTranslations()
  const [, setSearchOpen] = useAtom(searchOpenAtom)

  const handleClick = useCallback(() => {
    setSearchOpen(true)
  }, [setSearchOpen])

  return (
    <button
      onClick={handleClick}
      className={cn(
        "flex items-center gap-3 px-6 py-3 bg-white/50 dark:bg-white/5 rounded-full border border-neutral-200/60 dark:border-white/10",
        "hover:bg-white/70 dark:hover:bg-white/10 transition-all duration-300",
        "focus:outline-none focus:ring-2 focus:ring-primary/50",
      )}
    >
      <SearchIcon className="w-5 h-5 text-primary" />
      <span className="text-sm font-medium text-neutral-700 dark:text-gray-300">
        {t("integrations.loadingMore.searchHint")}
      </span>

      <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-primary/10 px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
        <span className="text-xs">⌘</span>
        K
      </kbd>
    </button>
  )
}