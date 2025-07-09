import type { NextRequest } from "next/server"
import type { Locale } from "@/i18n/routing"
import { cookies } from "next/headers"

import { NextResponse } from "next/server"
import { defaultLocale } from "@/i18n/routing"
import { getServerToken } from "@/lib/supabase/server"
import { serverLogger } from "@/lib/utils"
import { IntegrationService } from "@/services/integration"

const searchApiLogger = serverLogger.withTag("search-api")
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const query = searchParams.get("q")

    const cookieStore = await cookies()
    const accessToken = await getServerToken()
    const locale = (cookieStore.get("NEXT_LOCALE")?.value || defaultLocale) as Locale

    const integrationService = new IntegrationService(accessToken, locale)
    const { integrations } = await integrationService.getIntegrations()

    const filteredIntegrations = integrations.filter((integration) => {
      if (query) {
        return integration.name.toLowerCase().includes(query.toLowerCase())
      }

      return true
    })

    return NextResponse.json(filteredIntegrations, {
      headers: {
        "Content-Type": "application/json",
      },
    })
  } catch (error) {
    searchApiLogger.error("Error in search API", { error })
    return NextResponse.json({ error: "Internal server error" }, { status: 500 })
  }
}
