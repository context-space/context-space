import type { Message } from "@ai-sdk/react"
import { clientLogger } from "./utils"

const storageLogger = clientLogger.withTag("playground-storage")

// 统一的playground数据结构 - 按provider组织
interface PlaygroundProviderData {
  chatHistory: Message[]
  inputHistory: string[]
  lastUpdated: number
}

// 全局playground数据结构 - 包含所有provider的数据
interface GlobalPlaygroundData {
  providers: Record<string, PlaygroundProviderData>
  lastUpdated: number
}

// 获取全局存储键
const getGlobalPlaygroundKey = (userId?: string) => {
  return userId ? `playground_data_${userId}` : `playground_data_anonymous`
}

// 加载全局playground数据
const loadGlobalPlaygroundData = (userId?: string): GlobalPlaygroundData => {
  try {
    const key = getGlobalPlaygroundKey(userId)
    const stored = sessionStorage.getItem(key)

    if (!stored) {
      return { providers: {}, lastUpdated: Date.now() }
    }

    const data = JSON.parse(stored)
    // 验证数据结构
    if (data && typeof data === "object" && typeof data.providers === "object") {
      return data
    }

    return { providers: {}, lastUpdated: Date.now() }
  } catch (error) {
    storageLogger.error("Failed to load global playground data", { error, userId })
    return { providers: {}, lastUpdated: Date.now() }
  }
}

// 保存全局playground数据
const saveGlobalPlaygroundData = (data: GlobalPlaygroundData, userId?: string) => {
  try {
    const key = getGlobalPlaygroundKey(userId)
    data.lastUpdated = Date.now()
    sessionStorage.setItem(key, JSON.stringify(data))
  } catch (error) {
    storageLogger.error("Failed to save global playground data", { error, userId })
  }
}

// 加载指定provider的playground数据
export const loadPlaygroundData = (provider: string, userId?: string): PlaygroundProviderData | null => {
  try {
    const globalData = loadGlobalPlaygroundData(userId)
    return globalData.providers[provider] || null
  } catch (error) {
    storageLogger.error("Failed to load playground data", { error, provider, userId })
    return null
  }
}

// 保存指定provider的playground数据
export const savePlaygroundData = (
  provider: string,
  chatHistory: Message[],
  inputHistory: string[],
  userId?: string,
) => {
  try {
    const globalData = loadGlobalPlaygroundData(userId)

    globalData.providers[provider] = {
      chatHistory,
      inputHistory,
      lastUpdated: Date.now(),
    }

    saveGlobalPlaygroundData(globalData, userId)
  } catch (error) {
    storageLogger.error("Failed to save playground data", { error, provider, userId })
  }
}

// 清除指定provider的数据
export const clearPlaygroundData = (provider: string, userId?: string) => {
  try {
    const globalData = loadGlobalPlaygroundData(userId)

    if (globalData.providers[provider]) {
      delete globalData.providers[provider]
      saveGlobalPlaygroundData(globalData, userId)
    }
  } catch (error) {
    storageLogger.error("Failed to clear playground data", { error, provider, userId })
  }
}

// 清除用户的所有playground数据
export const clearAllPlaygroundData = (userId?: string) => {
  try {
    if (typeof window === "undefined") return

    const key = getGlobalPlaygroundKey(userId)
    sessionStorage.removeItem(key)
  } catch (error) {
    storageLogger.error("Failed to clear all playground data", { error, userId })
  }
}

// 获取所有provider的数据统计
export const getPlaygroundDataStats = (userId?: string) => {
  try {
    const globalData = loadGlobalPlaygroundData(userId)
    return {
      totalProviders: Object.keys(globalData.providers).length,
      providers: Object.entries(globalData.providers).map(([provider, data]) => ({
        provider,
        chatCount: data.chatHistory.length,
        inputCount: data.inputHistory.length,
        lastUpdated: data.lastUpdated,
      })),
      globalLastUpdated: globalData.lastUpdated,
    }
  } catch (error) {
    storageLogger.error("Failed to get playground data stats", { error, userId })
    return {
      totalProviders: 0,
      providers: [],
      globalLastUpdated: Date.now(),
    }
  }
}
