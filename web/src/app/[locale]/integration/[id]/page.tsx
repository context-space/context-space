import type { Metadata } from "next"
import type { Locale } from "@/i18n/routing"
import { getTranslations } from "next-intl/server"
import { notFound } from "next/navigation"
import { cache } from "react"
import { ClientLayout, Header } from "@/components/integration"
import { BaseLayout } from "@/components/layouts"
import { Breadcrumbs } from "@/components/seo/breadcrumbs"
import { keywords as baseKeywords, baseURL } from "@/config"
import { getServerToken } from "@/lib/supabase/server"
import { serverLogger } from "@/lib/utils"
import { IntegrationService } from "@/services/integration"

const integrationPageLogger = serverLogger.withTag("integration-page")

interface IntegrationPageProps {
  params: Promise<{ id: string, locale: Locale }>
}

const getCachedIntegration = cache(async (id: string, locale: Locale, accessToken: string | undefined) => {
  const integrationService = new IntegrationService(accessToken, locale)
  return await integrationService.getIntegration(id)
})

export async function generateMetadata({ params }: IntegrationPageProps): Promise<Metadata> {
  const { id, locale } = await params

  try {
    const accessToken = await getServerToken()
    const integration = await getCachedIntegration(id, locale, accessToken)

    const title = `${integration.name} | Integration`
    const description = integration.description
      ? `Connect and integrate ${integration.name}. ${integration.description}`
      : `Connect and integrate ${integration.name} with Context Space.`

    const keywords = [
      integration.name,
      ...(integration.categories || []),
      ...baseKeywords,
      integration.auth_type === "oauth" ? "OAuth" : integration.auth_type,
    ].filter(Boolean).join(", ")

    // Generate OG image URL with integration data passed directly
    const ogParams = new URLSearchParams({
      title,
      description,
      integration_name: integration.name,
      integration_icon: integration.icon_url,
      integration_categories: integration.categories.join(","),
      auth_type: integration.auth_type,
      connection_status: integration.connection_status || "unconnected",
    })
    const ogImageUrl = `${baseURL}/api/og?${ogParams.toString()}`

    return {
      title,
      description,
      keywords,
      openGraph: {
        title,
        description,
        type: "website",
        images: [
          {
            url: ogImageUrl,
            width: 1200,
            height: 630,
            alt: `${integration.name} Integration - Context Space`,
          },
        ],
      },
      twitter: {
        card: "summary_large_image",
        title,
        description,
        images: [ogImageUrl],
      },
      robots: {
        index: integration.status === "active",
        follow: integration.status === "active",
      },
      alternates: {
        canonical: `/integration/${id}`,
        languages: {
          "en": `/en/integration/${id}`,
          "zh": `/zh/integration/${id}`,
          "zh-TW": `/zh-TW/integration/${id}`,
        },
      },
    }
  } catch (error) {
    integrationPageLogger.error("Error generating metadata", { error, id })

    // Fallback metadata
    return {
      title: `Integration | Context Space`,
      description: `Explore integrations on Context Space platform.`,
      robots: {
        index: false,
        follow: false,
      },
    }
  }
}

export default async function IntegrationPage({ params }: IntegrationPageProps) {
  const { id, locale } = await params
  const t = await getTranslations()
  try {
    const accessToken = await getServerToken()
    const integration = await getCachedIntegration(id, locale, accessToken)
    return (
      <BaseLayout>
        <div className="py-4">
          <Breadcrumbs
            items={[
              { name: t("nav.home"), url: "/" },
              { name: t("integrations.pageTitle"), url: "/integrations" },
              { name: integration.name, url: `/integration/${id}` },
            ]}
            className="text-sm text-muted-foreground"
          />
        </div>
        <div className="space-y-6 pb-8 pt-8 max-w-7xl mx-auto w-full">
          <Header
            name={integration.name}
            description={integration.description || ""}
            iconUrl={integration.icon_url}
          />
          <ClientLayout
            operations={integration.operations}
            provider={id}
            authType={integration.auth_type}
            permissions={integration.permissions}
            apiDocUrl={integration.api_doc_url || ""}
            isConnected={integration.connection_status === "connected"}
            providerId={id}
            credentialId={integration.credential?.id ?? ""}
            authorizedPermissions={integration.credential?.permissions ?? []}
          />
        </div>
      </BaseLayout>
    )
  } catch (error) {
    integrationPageLogger.error("Error fetching provider details", { error })
    return notFound()
  }
}
