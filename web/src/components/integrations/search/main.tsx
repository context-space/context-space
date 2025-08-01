"use client"

import { useAtom } from "jotai"
import { Loader2 } from "lucide-react"
import { useTranslations } from "next-intl"
import { useCallback, useEffect, useMemo, useState } from "react"
import { ProviderAvatar } from "@/components/common/avatar"
import { CommandDialog, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command"
import { ScrollArea } from "@/components/ui/scroll-area"
import { useRouter } from "@/i18n/navigation"

import { clientLogger } from "@/lib/utils"
import { searchOpenAtom } from "./search-atom"
import "./style.css"

const integrationsSearchLogger = clientLogger.withTag("integrations-search")

interface FzfInstance {
  find: (query: string) => Array<{ item: Provider, score: number, positions: Set<number> }>
}

async function importFzf() {
  try {
    const { Fzf } = await import("fzf")
    return Fzf
  } catch (error) {
    integrationsSearchLogger.error("Failed to import fzf", { error })
    return null
  }
}

interface Provider {
  id: string
  identifier: string
  name: string
  description?: string
  category?: string
  connection_status?: "connected" | "disconnected" | "pending"
  icon?: string
  icon_url?: string
}

interface GroupedProviders {
  [category: string]: Provider[]
}

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

function useFzfSearch(providers: Provider[]) {
  const [fzf, setFzf] = useState<FzfInstance | null>(null)
  const [fzfLoading, setFzfLoading] = useState(false)

  useEffect(() => {
    if (providers.length > 0 && !fzf && !fzfLoading) {
      setFzfLoading(true)
      importFzf().then((FzfClass) => {
        if (FzfClass) {
          try {
            const fzfInstance = new FzfClass(providers, {
              selector: (provider: Provider) => {
                const searchText = [
                  provider.name,
                  // provider.description || '',
                  provider.identifier || "",
                  // provider.category || ''
                ].filter(Boolean).join(" ")
                return searchText
              },
            }) as FzfInstance
            setFzf(fzfInstance)
          } catch (error) {
            integrationsSearchLogger.error("Failed to initialize fzf", { error })
          }
        }
        setFzfLoading(false)
      })
    }
  }, [providers, fzf, fzfLoading])

  const search = useCallback((query: string): Provider[] => {
    if (!fzf || !query.trim()) {
      return providers
    }
    const results = fzf.find(query.trim())
    return results.map(result => result.item)
  }, [fzf, providers])

  return { search, fzfLoading }
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
      className="w-full"
    >
      <div className="mr-3 relative shrink-0">
        <ProviderAvatar
          src={provider.icon || provider.icon_url}
          alt={`${provider.name} logo`}
          className="size-8"
        />
      </div>
      <div className="flex flex-col w-full">
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

export function SearchDialog() {
  const t = useTranslations()
  const router = useRouter()
  const [searchOpen, setSearchOpen] = useAtom(searchOpenAtom)
  const [searchQuery, setSearchQuery] = useState("")

  const { providers, loading, fetchProviders } = useProviders()
  const { search: fzfSearch, fzfLoading } = useFzfSearch(providers)

  const handleClose = useCallback(() => {
    setSearchOpen(false)
    setSearchQuery("")
  }, [setSearchOpen])

  const handleProviderSelect = useCallback((provider: Provider) => {
    integrationsSearchLogger.info("Provider selected", { name: provider.name })
    if (provider.identifier) {
      router.push(`/integration/${provider.identifier}`)
    }
    handleClose()
  }, [router, handleClose])

  useEffect(() => {
    if (searchOpen && providers.length === 0) {
      fetchProviders()
    }
  }, [searchOpen, providers.length, fetchProviders])

  const filteredAndGroupedProviders = useMemo((): GroupedProviders => {
    const filteredProviders = fzfSearch(searchQuery)

    const groups: GroupedProviders = {}
    filteredProviders.forEach((provider: Provider) => {
      const groupKey = provider.connection_status || t("integrations.status.active")
      if (!groups[groupKey]) {
        groups[groupKey] = []
      }
      groups[groupKey].push(provider)
    })

    return groups
  }, [fzfSearch, searchQuery, t])

  return (
    <CommandDialog
      open={searchOpen}
      onOpenChange={open => !open && handleClose()}
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
        <ScrollArea className="max-h-[400px] h-[calc(100vh-100px)] px-1 pb-2">
          {(loading || fzfLoading)
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
  )
}