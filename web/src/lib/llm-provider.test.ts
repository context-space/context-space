import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import { clampTemperature, createLLMProvider, detectProvider, getProviderConfig } from "./llm-provider"

describe("llm-provider", () => {
  const originalEnv = process.env

  beforeEach(() => {
    vi.resetModules()
    process.env = { ...originalEnv }
    // Clear LLM-related env vars
    delete process.env.LLM_PROVIDER
    delete process.env.LLM_MODEL
    delete process.env.OPENAI_API_KEY
    delete process.env.OPENAI_BASE_URL
    delete process.env.MINIMAX_API_KEY
    delete process.env.MINIMAX_BASE_URL
  })

  afterEach(() => {
    process.env = originalEnv
  })

  describe("detectProvider", () => {
    it("should default to openai when no env vars are set", () => {
      expect(detectProvider()).toBe("openai")
    })

    it("should respect explicit LLM_PROVIDER=minimax", () => {
      process.env.LLM_PROVIDER = "minimax"
      expect(detectProvider()).toBe("minimax")
    })

    it("should respect explicit LLM_PROVIDER=openai", () => {
      process.env.LLM_PROVIDER = "openai"
      expect(detectProvider()).toBe("openai")
    })

    it("should be case-insensitive for LLM_PROVIDER", () => {
      process.env.LLM_PROVIDER = "MiniMax"
      expect(detectProvider()).toBe("minimax")
    })

    it("should auto-detect minimax when only MINIMAX_API_KEY is set", () => {
      process.env.MINIMAX_API_KEY = "test-key"
      expect(detectProvider()).toBe("minimax")
    })

    it("should prefer openai when both API keys are set", () => {
      process.env.OPENAI_API_KEY = "openai-key"
      process.env.MINIMAX_API_KEY = "minimax-key"
      expect(detectProvider()).toBe("openai")
    })

    it("should prefer LLM_PROVIDER over auto-detection", () => {
      process.env.LLM_PROVIDER = "minimax"
      process.env.OPENAI_API_KEY = "openai-key"
      expect(detectProvider()).toBe("minimax")
    })

    it("should fallback to openai for unknown LLM_PROVIDER value", () => {
      process.env.LLM_PROVIDER = "unknown-provider"
      expect(detectProvider()).toBe("openai")
    })
  })

  describe("getProviderConfig", () => {
    it("should return openai config by default", () => {
      const config = getProviderConfig("openai")
      expect(config.name).toBe("OpenAI")
      expect(config.defaultModel).toBe("gpt-4.1")
      expect(config.baseURL).toBe("https://api.openai.com/v1")
    })

    it("should return minimax config", () => {
      const config = getProviderConfig("minimax")
      expect(config.name).toBe("MiniMax")
      expect(config.defaultModel).toBe("MiniMax-M2.7")
      expect(config.baseURL).toBe("https://api.minimax.io/v1")
    })

    it("should include MiniMax models", () => {
      const config = getProviderConfig("minimax")
      expect(config.models).toContain("MiniMax-M2.7")
      expect(config.models).toContain("MiniMax-M2.7-highspeed")
    })

    it("should have correct temperature range for MiniMax", () => {
      const config = getProviderConfig("minimax")
      expect(config.temperatureRange.min).toBe(0.01)
      expect(config.temperatureRange.max).toBe(1)
    })

    it("should have correct temperature range for OpenAI", () => {
      const config = getProviderConfig("openai")
      expect(config.temperatureRange.min).toBe(0)
      expect(config.temperatureRange.max).toBe(2)
    })
  })

  describe("clampTemperature", () => {
    it("should clamp low temperature for MiniMax", () => {
      expect(clampTemperature(0, "minimax")).toBe(0.01)
    })

    it("should clamp high temperature for MiniMax", () => {
      expect(clampTemperature(2.0, "minimax")).toBe(1)
    })

    it("should pass through valid MiniMax temperature", () => {
      expect(clampTemperature(0.7, "minimax")).toBe(0.7)
    })

    it("should allow zero temperature for OpenAI", () => {
      expect(clampTemperature(0, "openai")).toBe(0)
    })

    it("should clamp high temperature for OpenAI", () => {
      expect(clampTemperature(3.0, "openai")).toBe(2)
    })
  })

  describe("createLLMProvider", () => {
    it("should create openai provider with API key", () => {
      process.env.OPENAI_API_KEY = "test-openai-key"
      const { model, providerType, config } = createLLMProvider("openai")
      expect(providerType).toBe("openai")
      expect(model).toBe("gpt-4.1")
      expect(config.name).toBe("OpenAI")
    })

    it("should create minimax provider with API key", () => {
      process.env.MINIMAX_API_KEY = "test-minimax-key"
      const { model, providerType, config } = createLLMProvider("minimax")
      expect(providerType).toBe("minimax")
      expect(model).toBe("MiniMax-M2.7")
      expect(config.name).toBe("MiniMax")
    })

    it("should respect LLM_MODEL override", () => {
      process.env.OPENAI_API_KEY = "test-key"
      process.env.LLM_MODEL = "custom-model"
      const { model } = createLLMProvider("openai")
      expect(model).toBe("custom-model")
    })

    it("should use custom MINIMAX_BASE_URL", () => {
      process.env.MINIMAX_API_KEY = "test-key"
      process.env.MINIMAX_BASE_URL = "https://custom.minimax.io/v1"
      const { client } = createLLMProvider("minimax")
      expect(client).toBeDefined()
    })

    it("should auto-detect provider when none specified", () => {
      process.env.MINIMAX_API_KEY = "test-minimax-key"
      const { providerType } = createLLMProvider()
      expect(providerType).toBe("minimax")
    })
  })
})
