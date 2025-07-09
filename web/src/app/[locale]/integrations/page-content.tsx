"use client"

import type { IntegrationCollection } from "@/typings"
import { useTranslations } from "next-intl"
import { Latest, Search, Status } from "@/components/integrations"

interface IntegrationsPageContentProps extends IntegrationCollection { }

export function IntegrationsPageContent({
  integrations,
  recommended_integrations,
  provider_statistics,
}: IntegrationsPageContentProps) {
  const t = useTranslations()

  return (
    <div className="max-w-7xl mx-auto w-full">
      <div className="py-16">
        <div className="flex flex-col items-center gap-8">
          <div className="flex flex-col items-center gap-3 relative">
            <h1 className="text-4xl font-bold tracking-tight text-neutral-900 dark:text-white sm:text-5xl">
              {t("integrations.pageTitle")}
            </h1>
            <p className="text-neutral-600 dark:text-gray-400 text-lg text-center">
              {t("integrations.pageSubtitle")}
            </p>
          </div>
          <div className="w-full max-w-2xl">
            <Search />
          </div>
          <div className="flex items-center gap-6 pt-4">
            <Status
              status="active"
              count={provider_statistics.active}
              type="text"
              className="hidden sm:flex"
            />
            <div className="hidden sm:block h-4 w-px bg-neutral-200 dark:bg-white/10"></div>
            <Status
              status="free"
              count={provider_statistics.free}
              type="text"
            />
            <div className="h-4 w-px bg-neutral-200 dark:bg-white/10"></div>
            <Status
              status="connected"
              count={provider_statistics.connected}
              type="text"
            />
          </div>
        </div>
      </div>

      <main className="pb-16">
        {/* <Latest integrations={recommended_integrations} /> */}
        <Latest integrations={integrations} />
      </main>
    </div>
  )
}
