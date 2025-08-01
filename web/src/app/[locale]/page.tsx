import type { Locale } from "@/i18n/routing"
import { Features, Hero, Roadmap } from "@/components/landing"
import { HomeLayout } from "@/components/layouts"
import { StructuredData } from "@/components/seo/structured-data"
import {
  generateOrganizationStructuredData,
  generateSoftwareApplicationStructuredData,
  generateWebSiteStructuredData,
} from "@/lib/seo/structured-data"

interface Props {
  params: Promise<{ locale: Locale }>
}

async function Page({ params }: Props) {
  const structuredData = [
    generateWebSiteStructuredData(),
    generateOrganizationStructuredData(),
    generateSoftwareApplicationStructuredData(),
  ]

  return (
    <HomeLayout>
      <StructuredData data={structuredData} />
      <Hero />
      <Features />
      <Roadmap />
    </HomeLayout>
  )
}

export default Page
