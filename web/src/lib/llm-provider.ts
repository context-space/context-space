import { createOpenAI } from "@ai-sdk/openai"

/**
 * Supported LLM provider identifiers.
 */
export type LLMProviderType = "openai" | "minimax"

/**
 * Model configuration for each provider.
 */
interface ProviderConfig {
  name: string
  defaultModel: string
  models: string[]
  baseURL: string
  apiKeyEnv: string
  temperatureRange: { min: number, max: number }
}

/**
 * Provider preset configurations.
 */
const PROVIDER_CONFIGS: Record<LLMProviderType, ProviderConfig> = {
  openai: {
    name: "OpenAI",
    defaultModel: "gpt-4.1",
    models: ["gpt-4.1", "gpt-4.1-mini", "gpt-4o", "gpt-4o-mini"],
    baseURL: "https://api.openai.com/v1",
    apiKeyEnv: "OPENAI_API_KEY",
    temperatureRange: { min: 0, max: 2 },
  },
  minimax: {
    name: "MiniMax",
    defaultModel: "MiniMax-M2.7",
    models: ["MiniMax-M2.7", "MiniMax-M2.7-highspeed"],
    baseURL: "https://api.minimax.io/v1",
    apiKeyEnv: "MINIMAX_API_KEY",
    temperatureRange: { min: 0.01, max: 1 },
  },
}

/**
 * Detects the active LLM provider from environment variables.
 * Priority: LLM_PROVIDER env var > auto-detect from API keys.
 */
export function detectProvider(): LLMProviderType {
  const explicit = process.env.LLM_PROVIDER?.toLowerCase()
  if (explicit && explicit in PROVIDER_CONFIGS) {
    return explicit as LLMProviderType
  }

  // Auto-detect: prefer MiniMax if its key is set and OpenAI's is not
  if (process.env.MINIMAX_API_KEY && !process.env.OPENAI_API_KEY) {
    return "minimax"
  }

  return "openai"
}

/**
 * Returns the provider configuration for the given provider type.
 */
export function getProviderConfig(provider?: LLMProviderType): ProviderConfig {
  return PROVIDER_CONFIGS[provider ?? detectProvider()]
}

/**
 * Clamps temperature to the provider's supported range.
 */
export function clampTemperature(temp: number, provider?: LLMProviderType): number {
  const config = getProviderConfig(provider)
  return Math.max(config.temperatureRange.min, Math.min(temp, config.temperatureRange.max))
}

/**
 * Creates an AI SDK model instance for the active LLM provider.
 * Uses @ai-sdk/openai since MiniMax exposes an OpenAI-compatible API.
 */
export function createLLMProvider(provider?: LLMProviderType) {
  const providerType = provider ?? detectProvider()
  const config = PROVIDER_CONFIGS[providerType]

  const apiKey = providerType === "minimax"
    ? process.env.MINIMAX_API_KEY
    : process.env.OPENAI_API_KEY

  const baseURL = providerType === "minimax"
    ? (process.env.MINIMAX_BASE_URL ?? config.baseURL)
    : (process.env.OPENAI_BASE_URL ?? config.baseURL)

  const client = createOpenAI({ apiKey, baseURL })
  const model = process.env.LLM_MODEL ?? config.defaultModel

  return { client, model, config, providerType }
}
