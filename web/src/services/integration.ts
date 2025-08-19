import type { Locale } from "@/i18n/routing"
import type { ConnectionStatus, Credential, IntegrationCollection, IntegrationDetail, Provider, ProviderDetail } from "@/typings"
import defu from "defu"
import { defaultLocale } from "@/i18n/routing"
import { genAuthorization } from "@/lib/utils/fetch"
import { remoteFetchWithErrorHandling } from "@/lib/utils/fetch/server"

interface FilterOption {
  filters?: {
    auth_type?: string
    provider_name?: string
    tag?: "hot" | "recommend"
  }
  pagination?: {
    page?: number
    page_size?: number
  }
  sort?: {
    field?: string
    order?: string
  }
}

const defaultOptions: FilterOption = {
  // filters: {
  //   tag: "featured",
  // },
  pagination: {
    page: 1,
    page_size: 100,
  },
  sort: {
    field: "created_at",
    order: "desc",
  },
}

export class IntegrationService {
  #token: string | undefined
  #locale: Locale

  constructor(token?: string, locale: Locale = defaultLocale) {
    this.#token = token
    this.#locale = locale
  }

  async getProviders(options?: FilterOption) {
    try {
      const mergedOptions = defu(options, defaultOptions)
      const res = await remoteFetchWithErrorHandling<{ providers: Provider[] }>("/providers/filter", {
        method: "POST",
        body: JSON.stringify(mergedOptions),
        headers: {
          "Accept-Language": this.#locale,
        },
      })
      return res.providers
    } catch {
      return []
    }
  }

  async #getCredentials() {
    try {
      const res = await remoteFetchWithErrorHandling<{ credentials: Credential[] }>("/credentials", {
        headers: genAuthorization(this.#token),
      })
      return res.credentials
    } catch {
      return []
    }
  }

  async #getProvider(id: string) {
    const res = await remoteFetchWithErrorHandling<ProviderDetail>(`/providers/${id}`, {
      headers: {
        "Accept-Language": this.#locale,
      },
    })
    return res
  }

  async #getCredential(id: string) {
    try {
      const res = await remoteFetchWithErrorHandling<Credential>(`/credentials/provider/${id}`, {
        headers: genAuthorization(this.#token),
      })
      return res
    } catch {
    }
  }

  #getConnectionStatus(provider: Provider, credential?: Credential): ConnectionStatus {
    return provider.auth_type === "none" ? "free" : credential?.is_valid ? "connected" : "unconnected"
  }

  async getIntegrations(): Promise<IntegrationCollection> {
    const [providers, credentials] = await Promise.all([this.getProviders(), this.#getCredentials()])
    const integrations = providers.map((provider) => {
      const credential = credentials.find(cred => cred.provider_identifier === provider.identifier)
      const connection_status = this.#getConnectionStatus(provider, credential)
      return {
        ...provider,
        credential,
        connection_status,
      }
    })
    const provider_statistics = integrations.reduce((acc, integration) => {
      acc[integration.connection_status] = (acc[integration.connection_status] || 0) + 1
      acc[integration.status] = (acc[integration.status] || 0) + 1
      return acc
    }, {} as IntegrationCollection["provider_statistics"])
    return {
      integrations,
      recommended_integrations: integrations.filter(integration => integration.tags?.includes("recommend")),
      hot_integrations: integrations.filter(integration => integration.tags?.includes("hot")),
      provider_statistics,
    }
  }

  async getIntegration(id: string): Promise<IntegrationDetail> {
    const [provider, credential] = await Promise.all([this.#getProvider(id), this.#getCredential(id)])
    const connection_status = this.#getConnectionStatus(provider, credential)
    return {
      ...provider,
      credential,
      connection_status,
    }
  }
}
