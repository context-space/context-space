import type { BaseRemoteAPIResponse, Credential } from "@/typings"
import { clientLogger } from "@/lib/utils"
import { clientRemoteFetch } from "@/lib/utils/fetch/client"

const credentialsLogger = clientLogger.withTag("credentials")
export async function connectAPIKey(identifier: string, apiKey: string) {
  return await clientRemoteFetch(`/credentials/auth/apikey/${identifier}`, {
    method: "POST",
    body: JSON.stringify({ api_key: apiKey }),
  }) as BaseRemoteAPIResponse<Credential>
}

async function openOAuthPopup(identifier: string, url: string): Promise<boolean> {
  const POPUP_WIDTH = 500
  const POPUP_HEIGHT = 700
  const width = POPUP_WIDTH
  const height = POPUP_HEIGHT
  const left = (window.screen.width - width) / 2
  const top = (window.screen.height - height) / 2

  return new Promise((resolve) => {
    const popup = window.open(
      `/mcp-oauth-callback?url=${encodeURIComponent(url)}`,
      "OAuth",
      `width=${width},height=${height},left=${left},top=${top},toolbar=no,menubar=no,scrollbars=yes,resizable=yes,location=no,status=no`,
    )

    if (!popup) {
      resolve(false)

      return
    }

    let resolved = false

    // Listen for messages from the popup (for error handling)
    const handleMessage = (event: MessageEvent) => {
      if (event.origin !== window.location.origin) {
        return
      }

      if (event.data.type === "oauth-error") {
        credentialsLogger.info("OAuth error received from popup", { error: event.data.error })
        resolved = true
        window.removeEventListener("message", handleMessage)
        resolve(false)
      }
    }

    window.addEventListener("message", handleMessage)

    // Check if popup is closed (primary method for success detection)
    const checkPopup = setInterval(async () => {
      if (popup.closed && !resolved) {
        clearInterval(checkPopup)
        window.removeEventListener("message", handleMessage)

        // Check if the connection was successful
        try {
          const isConnected = await checkCredential(identifier)
          resolve(isConnected)
        } catch (error) {
          // credentialsLogger.error("Failed to check credential", { error })
          resolve(false)
        }
      }
    }, 1000)

    // Timeout after 5 minutes
    setTimeout(() => {
      if (!popup.closed) {
        popup.close()
      }
      if (!resolved) {
        clearInterval(checkPopup)
        window.removeEventListener("message", handleMessage)
        resolve(false)
      }
    }, 5 * 60 * 1000)
  })
}

export async function connectOAuth(identifier: string, permissions: string[]) {
  const { data } = (await clientRemoteFetch(`/credentials/auth/oauth/${identifier}/auth-url`, {
    method: "POST",
    body: JSON.stringify({
      permissions,
      redirect_url: `${window.location.origin}/mcp-oauth-callback`,
    }),
    async onResponseError({ response }) {
      credentialsLogger.error("OAuth response error", { data: response._data })
    },
  }))
  const { auth_url } = data
  credentialsLogger.info("Generated auth URL", { auth_url })
  const isConnected = await openOAuthPopup(identifier, auth_url)
  if (!isConnected) {
    throw new Error("Failed to connect")
  }
}

export async function disconnectCredential(credentialId: string): Promise<void> {
  await clientRemoteFetch(`/credentials/${credentialId}`, {
    method: "DELETE",
  })
}

export async function checkCredential(identifier: string): Promise<boolean> {
  try {
    const { data } = (await clientRemoteFetch(`/credentials/auth/oauth/${identifier}`)) as BaseRemoteAPIResponse<Credential>
    return data.is_valid
  } catch (error) {
    return false
  }
}
