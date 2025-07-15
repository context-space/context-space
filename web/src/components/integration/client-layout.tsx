"use client"

import type { Operation, Permission } from "@/typings"
import { Connect, McpConfig, Operations, Playground } from "./index"

interface ClientLayoutProps {
  operations: Operation[]
  provider: string
  authType: "oauth" | "apikey" | "none"
  permissions: Permission[]
  apiDocUrl: string
  isConnected: boolean
  providerId: string
  credentialId: string
  authorizedPermissions: string[]
  integrationName?: string
}

export function ClientLayout({
  operations,
  provider,
  authType,
  permissions,
  apiDocUrl,
  isConnected,
  providerId,
  credentialId,
  authorizedPermissions,
  integrationName,
}: ClientLayoutProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 min-h-[calc(100vh-14rem)]">
      <div className="flex flex-col gap-6">
        <Connect
          authType={authType}
          permissions={permissions}
          apiDocUrl={apiDocUrl}
          isConnected={isConnected}
          providerId={providerId}
          credentialId={credentialId}
          authorizedPermissions={authorizedPermissions}
        />
        <McpConfig
          authType={authType}
          isConnected={isConnected}
          integrationId={providerId}
          integrationName={integrationName || provider}
        />
        <Operations operations={operations} />
      </div>
      <Playground
        provider={provider}
        authType={authType}
        isConnected={isConnected}
        providerId={providerId}
      />
    </div>
  )
}
