export interface APIKey {
  id: string
  name: string
  description?: string
  key_value?: string // Only included when creating
  created_at: string
}

export interface CreateAPIKeyRequest {
  name: string
  description?: string
}

export interface APIKeysResponse {
  api_keys: APIKey[]
}

export interface APIKeyResponse {
  id: string
  name: string
  description?: string
  key_value?: string
  created_at: string
}
