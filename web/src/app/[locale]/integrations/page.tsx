import type { Metadata } from "next"
import type { Locale } from "@/i18n/routing"
import { getTranslations } from "next-intl/server"
import { BaseLayout } from "@/components/layouts"
import { Breadcrumbs } from "@/components/seo/breadcrumbs"
import { getServerToken } from "@/lib/supabase/server"
import { IntegrationService } from "@/services/integration"
import { IntegrationsPageContent } from "./page-content"

interface IntegrationsPageProps {
  params: Promise<{ locale: Locale }>
}

export const metadata: Metadata = {
  title: "Integrations",
  description: "Discover and connect with hundreds of powerful integrations for your Context Space platform. Browse MCP-compatible tools and services.",
  keywords: "integrations, MCP, Model Context Protocol, API, tools, services, workflow automation",
  openGraph: {
    title: "Integrations | Context Space",
    description: "Discover and connect with hundreds of powerful integrations for your Context Space platform.",
    images: [
      {
        url: `/api/og?title=Integrations&description=${encodeURIComponent("Discover and connect with hundreds of powerful integrations")}`,
        width: 1200,
        height: 630,
        alt: "Context Space Integrations",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "Integrations | Context Space",
    description: "Discover and connect with hundreds of powerful integrations for your Context Space platform.",
    images: [`/api/og?title=Integrations&description=${encodeURIComponent("Discover and connect with hundreds of powerful integrations")}`],
  },
}

export default async function IntegrationsPage({ params }: IntegrationsPageProps) {
  const { locale } = await params
  const t = await getTranslations()
  const accessToken = await getServerToken()
  const integrationService = new IntegrationService(accessToken, locale)
  const { integrations, recommended_integrations, provider_statistics, hot_integrations } = await integrationService.getIntegrations()

  return (
    <BaseLayout>
      <div className="py-4">
        <Breadcrumbs
          items={[
            { name: t("nav.home"), url: "/" },
            { name: t("integrations.pageTitle"), url: "/integrations" },
          ]}
          className="text-sm text-muted-foreground"
        />
      </div>
      <IntegrationsPageContent
        integrations={integrations}
        recommended_integrations={recommended_integrations}
        hot_integrations={hot_integrations}
        provider_statistics={provider_statistics}
      />
    </BaseLayout>
  )
}
