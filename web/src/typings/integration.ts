export type ProviderStatus = "active" | "inactive" | "maintenance" | "deprecated"
export type ConnectionStatus = "connected" | "free" | "unconnected"
export type AuthType = "oauth" | "none" | "apikey"

export interface BaseRemoteAPIResponse<T> {
  success: boolean
  message: string
  code: number
  data: T
}

export interface Parameter {
  [k: string]: any
}

export interface Operation {
  id: string
  name: string
  description: string
  category: string
  parameters: Parameter[]
  identifier: string
  required_permissions: Permission[]
}

export interface Permission {
  identifier: string
  name: string
  description: string
}

export interface Provider {
  id: string
  identifier: string
  name: string
  icon_url: string
  description: string
  api_doc_url?: string
  auth_type: AuthType
  categories: string[]
  status: ProviderStatus
  tags: string[]
}

export interface ProviderDetail extends Provider {
  permissions: Permission[]
  operations: Operation[]
}

export interface Credential {
  created_at: string
  id: string
  is_valid: boolean
  provider_identifier: string
  permissions?: string[]
  type: string
  user_id: string
  api_doc_url?: string
}

interface AuthedInfo {
  credential?: Credential
  connection_status?: ConnectionStatus
}

export interface Integration extends Provider, AuthedInfo {}

export type ProviderStatistics = Record<ProviderStatus | ConnectionStatus, number | undefined>
export interface IntegrationCollection {
  integrations: Integration[]
  recommended_integrations: Integration[]
  hot_integrations: Integration[]
  provider_statistics: ProviderStatistics
}

export interface IntegrationDetail extends ProviderDetail, AuthedInfo {}

export interface OAuthConnectURL {
  auth_url: string
  permissions: string[]
  provider_identifier: string
  redirect_url: string
  state: string
}

export interface Invocation {
  error?: string
  completed_at: string
  created_at: string
  duration_ms: number
  error_message: string
  id: string
  operation_identifier: string
  parameters: Parameter
  provider_identifier: string
  response_data: any[]
  started_at: string
  status: string
  user_id: string
}
