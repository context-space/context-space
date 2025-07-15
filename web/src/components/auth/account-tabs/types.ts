export interface AccountUser {
  id: string
  email: string
  displayName: string
  avatar?: string
  isAnonymous: boolean
  createdAt: string
  lastSignInAt?: string
  providers: string[]
}

export type TabType = "profile" | "api-keys" | "mcp"
