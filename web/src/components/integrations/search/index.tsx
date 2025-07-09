"use client"

import { Loader2, Search as SearchIcon } from "lucide-react"
import { useTranslations } from "next-intl"
import { useCallback, useEffect, useMemo, useState } from "react"
import { ProviderAvatar } from "@/components/integration/provider-avatar"
import { CommandDialog, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command"
import { ScrollArea } from "@/components/ui/scroll-area"
import { useRouter } from "@/i18n/navigation"

import { clientLogger, cn } from "@/lib/utils"
import "./style.css"

const integrationsSearchLogger = clientLogger.withTag("integrations-search")

// Improved type definitions
interface Provider {
  id: string
  identifier: string
  name: string
  description?: string
  category?: string
  connectionStatus?: "connected" | "disconnected" | "pending"
  connection_status?: "connected" | "disconnected" | "pending"
  icon?: string
  icon_url?: string
}

interface GroupedProviders {
  [category: string]: Provider[]
}

// Custom hook for mobile detection
function useIsMobile(): boolean {
  const [isMobile, setIsMobile] = useState(false)

  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768)
    }

    checkMobile()
    window.addEventListener("resize", checkMobile)

    return () => window.removeEventListener("resize", checkMobile)
  }, [])

  return isMobile
}

// Custom hook for keyboard shortcuts
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

// Custom hook for providers data
function useProviders() {
  const [providers, setProviders] = useState<Provider[]>([])
  const [loading, setLoading] = useState(false)

  const fetchProviders = useCallback(async () => {
    setLoading(true)
    try {
      const response = await fetch("/api/search")
      if (response.ok) {
        const data = await response.json()
        setProviders(data)
      } else {
        integrationsSearchLogger.error("Failed to fetch providers")
        setProviders([])
      }
    } catch (error) {
      integrationsSearchLogger.error("Error fetching providers", { error })
      setProviders([])
    } finally {
      setLoading(false)
    }
  }, [])

  return { providers, loading, fetchProviders }
}

// Search trigger button component
interface SearchTriggerProps {
  onClick: () => void
  placeholderText: string
  isMobile: boolean
}

function SearchTrigger({ onClick, placeholderText, isMobile }: SearchTriggerProps) {
  return (
    <div className="relative w-full">
      <button
        type="button"
        onClick={onClick}
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
            <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
              <span className="text-xs">⌘</span>
              K
            </kbd>
          </div>
        )}
      </button>
    </div>
  )
}

interface LoadingStateProps {
  loadingText: string
}

function LoadingState({ loadingText }: LoadingStateProps) {
  return (
    <div className="flex items-center justify-center py-6">
      <Loader2 className="h-4 w-4 animate-spin" />
      <span className="ml-2 text-sm text-muted-foreground">{loadingText}</span>
    </div>
  )
}

interface ProviderItemProps {
  provider: Provider
  onSelect: (provider: Provider) => void
}

function ProviderItem({ provider, onSelect }: ProviderItemProps) {
  const handleSelect = useCallback(() => {
    onSelect(provider)
  }, [provider, onSelect])

  return (
    <CommandItem
      key={provider.id || provider.identifier}
      value={provider.name}
      onSelect={handleSelect}
    >
      <div className="mr-3 relative shrink-0">
        <ProviderAvatar
          src={provider.icon || provider.icon_url}
          alt={`${provider.name} logo`}
          className="size-8"
        />
      </div>
      <div className="flex flex-col min-w-0 flex-1">
        <span className="truncate">{provider.name}</span>
        {provider.description && (
          <span className="text-xs text-muted-foreground truncate">
            {provider.description}
          </span>
        )}
      </div>
    </CommandItem>
  )
}

export function Search() {
  const t = useTranslations()
  const router = useRouter()
  const [open, setOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState("")

  const isMobile = useIsMobile()
  const { providers, loading, fetchProviders } = useProviders()

  const handleToggle = useCallback(() => {
    setOpen(prev => !prev)
  }, [])

  useKeyboardShortcuts(handleToggle)

  const handleOpenChange = useCallback((newOpen: boolean) => {
    setOpen(newOpen)
    if (!newOpen) {
      setSearchQuery("")
    }
  }, [])

  const handleProviderSelect = useCallback((provider: Provider) => {
    integrationsSearchLogger.info("Provider selected", { name: provider.name })
    if (provider.identifier) {
      router.push(`/integration/${provider.identifier}`)
    }
    setOpen(false)
  }, [router])

  useEffect(() => {
    if (open && providers.length === 0) {
      fetchProviders()
    }
  }, [open, providers.length, fetchProviders])

  const placeholderText = useMemo(() => {
    return t("integrations.search.placeholder")
  }, [t])

  const filteredAndGroupedProviders = useMemo((): GroupedProviders => {
    let filteredProviders = providers
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase().trim()
      filteredProviders = providers.filter((provider: Provider) =>
        provider.name.toLowerCase().includes(query)
        || (provider.description?.toLowerCase().includes(query))
        || (provider.identifier?.toLowerCase().includes(query)),
      )
    }

    // Group filtered providers by connection status
    const groups: GroupedProviders = {}
    filteredProviders.forEach((provider: Provider) => {
      const groupKey = provider.connection_status || provider.connectionStatus || t("integrations.status.active")
      if (!groups[groupKey]) {
        groups[groupKey] = []
      }
      groups[groupKey].push(provider)
    })

    return groups
  }, [providers, searchQuery, t])

  return (
    <>
      <SearchTrigger
        onClick={handleToggle}
        placeholderText={placeholderText}
        isMobile={isMobile}
      />

      <CommandDialog
        open={open}
        onOpenChange={handleOpenChange}
        title={t("integrations.search.title")}
        description={t("integrations.search.description")}
        className="bg-white sprinkle-primary"
      >
        <CommandInput
          placeholder={t("integrations.search.placeholder")}
          value={searchQuery}
          onValueChange={setSearchQuery}
        />
        <CommandList className="max-h-none overflow-hidden">
          <ScrollArea className="h-[400px] px-1">
            {loading
              ? (
                  <LoadingState loadingText={t("common.loading")} />
                )
              : (
                  <>
                    <CommandEmpty>
                      <div className="flex flex-col min-w-0 gap-2 text-center">
                        <p className="text-muted-foreground">{t("integrations.search.empty")}</p>
                        <div className="flex flex-col gap-1 text-sm">
                          <a
                            href="https://github.com/context-space/context-space/issues/new"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-primary hover:text-primary/80 underline underline-offset-4 transition-colors"
                          >
                            {t("integrations.search.submitIssue")}
                          </a>
                          <span className="text-muted-foreground">{t("integrations.search.or")}</span>
                          <a
                            href="mailto:support@contextspace.ai"
                            className="text-primary hover:text-primary/80 underline underline-offset-4 transition-colors"
                          >
                            {t("integrations.search.contactUs")}
                          </a>
                        </div>
                      </div>
                    </CommandEmpty>
                    {Object.entries(filteredAndGroupedProviders).map(([category, categoryProviders]) => (
                      <CommandGroup key={category} heading={t(`integrations.status.${category}`)}>
                        {categoryProviders.map((provider: Provider) => (
                          <ProviderItem
                            key={provider.id || provider.identifier}
                            provider={provider}
                            onSelect={handleProviderSelect}
                          />
                        ))}
                      </CommandGroup>
                    )).sort((a) => {
                      return a.key === "connected" ? -1 : a.key === "free" ? -1 : 1
                    })}
                  </>
                )}
          </ScrollArea>
        </CommandList>
      </CommandDialog>
    </>
  )
}
