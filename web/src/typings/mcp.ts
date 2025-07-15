import type { APIKey } from "./api-keys"

export interface McpTabProps {
  onSwitchToApiKeys?: () => void
}

export interface McpConfigDisplayProps {
  apiKey: APIKey
  onCopy: (config: string) => Promise<boolean>
  onAddToCursor: (config: string) => void
  isAddingToCursor: boolean
  integrationId?: string
  integrationName?: string
}
