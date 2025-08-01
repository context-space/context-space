import type { IntegrationCollection } from "@/typings"
import { getTranslations } from "next-intl/server"
import { Featured, LoadingMore, Normal, Status } from "@/components/integrations"
import { SearchDialog, SearchTrigger } from "@/components/integrations/search"

interface IntegrationsPageContentProps extends IntegrationCollection { }

export async function IntegrationsPageContent({
  integrations,
  recommended_integrations,
  provider_statistics,
  hot_integrations,
}: IntegrationsPageContentProps) {
  const t = await getTranslations()

  return (
    <div className="max-w-7xl mx-auto w-full">
      <div id="integrations-page-content" className="flex flex-col items-center gap-8 py-16">
        <div className="flex flex-col items-center gap-3 relative">
          <h1 className="text-4xl font-bold tracking-tight text-neutral-900 dark:text-white sm:text-5xl">
            {t("integrations.pageTitle")}
          </h1>
          <p className="text-neutral-600 dark:text-gray-400 text-lg text-center">
            {t("integrations.pageSubtitle")}
          </p>
        </div>
        <div className="w-full max-w-2xl">
          <SearchTrigger />
        </div>
        <div className="flex items-center gap-6 pt-4">
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

      <div className="space-y-12">
        <Featured integrations={recommended_integrations} title={t("integrations.recommendedIntegrations")} />
        <Normal integrations={hot_integrations} title={t("integrations.hotIntegrations")} />
        <Normal integrations={integrations.slice(0, 12)} title={t("integrations.latestIntegrations")} />
      </div>

      <LoadingMore totalCount={integrations.length} />

      <SearchDialog />
    </div>
  )
}
