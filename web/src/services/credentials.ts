import type { BaseRemoteAPIResponse, Credential } from "@/typings"
import { clientLogger } from "@/lib/utils"
import { clientRemoteFetch } from "@/lib/utils/fetch/client"

const credentialsLogger = clientLogger.withTag("credentials")

class CredentialsService {
  /**
   * Connect using API Key authentication
   */
  async connectAPIKey(identifier: string, apiKey: string): Promise<BaseRemoteAPIResponse<Credential>> {
    return await clientRemoteFetch(`/credentials/auth/apikey/${identifier}`, {
      method: "POST",
      body: JSON.stringify({ api_key: apiKey }),
    }) as BaseRemoteAPIResponse<Credential>
  }

  /**
   * Connect using OAuth authentication
   */
  async connectOAuth(identifier: string, permissions: string[], redirectUrl: string): Promise<string> {
    try {
      const res = await clientRemoteFetch(`/credentials/auth/oauth/${identifier}/auth-url`, {
        method: "POST",
        body: JSON.stringify({
          permissions,
          redirect_url: redirectUrl,
        }),
      })
      return res.data.auth_url
    } catch (error) {
      credentialsLogger.error("Failed to connect with OAuth", { error })
      throw error
    }
  }

  /**
   * Disconnect a credential
   */
  async disconnectCredential(credentialId: string): Promise<void> {
    await clientRemoteFetch(`/credentials/${credentialId}`, {
      method: "DELETE",
    })
  }
}

export const credentialsService = new CredentialsService()
