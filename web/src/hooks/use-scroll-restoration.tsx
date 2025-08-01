"use client"

import type { ReadonlyURLSearchParams } from "next/navigation"
import { usePathname, useSearchParams } from "next/navigation"
import { useCallback, useEffect, useMemo, useRef } from "react"

// Types
interface ScrollPosition {
  top: number
  left: number
  timestamp: number
}

interface ScrollPositions {
  [key: string]: ScrollPosition
}

export interface ScrollRestorationManagerProps {
  /**
   * 滚动容器的选择器
   * 默认为 Radix ScrollArea 的 viewport 选择器
   */
  scrollContainerSelector?: string
  /**
   * 是否启用滚动恢复
   */
  enabled?: boolean
  /**
   * 滚动恢复的延迟时间（毫秒）
   */
  restoreDelay?: number
  /**
   * 存储键前缀
   */
  storageKeyPrefix?: string
  /**
   * 是否在页面卸载时保存滚动位置
   */
  saveOnUnload?: boolean
  /**
   * 防抖延迟时间（毫秒）
   */
  debounceDelay?: number
  /**
   * 位置清理的最大存储时间（毫秒）
   */
  maxAge?: number
}

// Constants
const DEFAULT_CONTAINER_SELECTOR = "[data-radix-scroll-area-viewport]"
const DEFAULT_STORAGE_PREFIX = "scroll_position"
const DEFAULT_DEBOUNCE_DELAY = 150
const DEFAULT_MAX_AGE = 24 * 60 * 60 * 1000 // 24 hours
const RESTORE_RETRY_DELAY = 5
const SCROLL_POSITION_TOLERANCE = 1

// Utility functions
const getStorageKey = (prefix: string, pathname: string, searchParams: string): string => {
  return `${prefix}_${pathname}${searchParams ? `_${searchParams}` : ""}`
}

const parseStorageData = (data: string | null): ScrollPositions => {
  try {
    return data ? JSON.parse(data) : {}
  } catch {
    return {}
  }
}

const saveToStorage = (data: ScrollPositions): void => {
  try {
    sessionStorage.setItem("scrollPositions", JSON.stringify(data))
  } catch (error) {
    console.warn("Failed to save to sessionStorage:", error)
  }
}

// Custom hook for scroll position management
function useScrollPositionStorage(
  pathname: string,
  searchParams: URLSearchParams | ReadonlyURLSearchParams,
  storagePrefix: string,
  maxAge: number,
) {
  const storageKey = useMemo(() =>
    getStorageKey(storagePrefix, pathname, searchParams.toString()), [storagePrefix, pathname, searchParams])

  const savePosition = useCallback((position: Omit<ScrollPosition, "timestamp">) => {
    const positions = parseStorageData(sessionStorage.getItem("scrollPositions"))
    positions[storageKey] = {
      ...position,
      timestamp: Date.now(),
    }
    saveToStorage(positions)
  }, [storageKey])

  const getPosition = useCallback((): ScrollPosition | null => {
    const positions = parseStorageData(sessionStorage.getItem("scrollPositions"))
    return positions[storageKey] || null
  }, [storageKey])

  const cleanupOldPositions = useCallback(() => {
    const positions = parseStorageData(sessionStorage.getItem("scrollPositions"))
    const now = Date.now()

    const cleanedPositions = Object.fromEntries(
      Object.entries(positions).filter(([, position]) =>
        now - position.timestamp < maxAge,
      ),
    )

    saveToStorage(cleanedPositions)
  }, [maxAge])

  return { savePosition, getPosition, cleanupOldPositions }
}

// Custom hook for scroll container management
function useScrollContainer(containerSelector: string) {
  const containerRef = useRef<Element | null>(null)

  const getContainer = useCallback((): Element | null => {
    if (containerRef.current) {
      return containerRef.current
    }

    const container = document.querySelector(containerSelector)
    if (container) {
      containerRef.current = container
      return container
    }

    return null
  }, [containerSelector])

  const resetContainer = useCallback(() => {
    containerRef.current = null
  }, [])

  return { getContainer, resetContainer }
}

// Custom hook for scroll restoration logic
export function useScrollRestoration({
  scrollContainerSelector = DEFAULT_CONTAINER_SELECTOR,
  enabled = true,
  debounceDelay = DEFAULT_DEBOUNCE_DELAY,
  storageKeyPrefix = DEFAULT_STORAGE_PREFIX,
  saveOnUnload = true,
  maxAge = DEFAULT_MAX_AGE,
}: Omit<ScrollRestorationManagerProps, "restoreDelay">) {
  const pathname = usePathname()
  const searchParams = useSearchParams()

  const { savePosition, getPosition, cleanupOldPositions } = useScrollPositionStorage(
    pathname,
    searchParams,
    storageKeyPrefix,
    maxAge,
  )

  const { getContainer, resetContainer } = useScrollContainer(scrollContainerSelector)

  const isRestoredRef = useRef(false)
  const scrollSaveTimerRef = useRef<NodeJS.Timeout | undefined>(undefined)

  const saveScrollPosition = useCallback(() => {
    if (!enabled) return

    try {
      const container = getContainer()
      if (container) {
        savePosition({
          top: container.scrollTop,
          left: container.scrollLeft,
        })
      } else {
        // Fallback to window scroll
        savePosition({
          top: window.scrollY,
          left: window.scrollX,
        })
      }
    } catch (error) {
      console.warn("Failed to save scroll position:", error)
    }
  }, [enabled, getContainer, savePosition])

  const restoreScrollPosition = useCallback(() => {
    if (!enabled || isRestoredRef.current) return

    try {
      const savedPosition = getPosition()
      if (!savedPosition) return

      const container = getContainer()

      if (container) {
        // Use ScrollArea container
        container.scrollTop = savedPosition.top || 0
        container.scrollLeft = savedPosition.left || 0

        // Verify scroll position was set correctly
        requestAnimationFrame(() => {
          if (Math.abs(container.scrollTop - (savedPosition.top || 0)) > SCROLL_POSITION_TOLERANCE) {
            container.scrollTop = savedPosition.top || 0
          }
        })
      } else {
        // Fallback to window scroll
        window.scrollTo(savedPosition.left || 0, savedPosition.top || 0)
      }

      isRestoredRef.current = true
    } catch (error) {
      console.warn("Failed to restore scroll position:", error)
    }
  }, [enabled, getPosition, getContainer])

  const handleScroll = useCallback(() => {
    if (scrollSaveTimerRef.current) {
      clearTimeout(scrollSaveTimerRef.current)
    }
    scrollSaveTimerRef.current = setTimeout(saveScrollPosition, debounceDelay)
  }, [saveScrollPosition, debounceDelay])

  const attemptRestore = useCallback(() => {
    const container = getContainer()
    const savedPosition = getPosition()

    if (container && savedPosition) {
      container.scrollTop = savedPosition.top || 0
      container.scrollLeft = savedPosition.left || 0
      isRestoredRef.current = true
    } else if (savedPosition) {
      // Retry with short delay if container not ready
      setTimeout(() => {
        if (!isRestoredRef.current) {
          restoreScrollPosition()
        }
      }, RESTORE_RETRY_DELAY)
    }
  }, [getContainer, getPosition, restoreScrollPosition])

  useEffect(() => {
    if (!enabled) return

    // Reset restoration state
    isRestoredRef.current = false
    resetContainer()

    // Cleanup old positions
    cleanupOldPositions()

    // Attempt immediate restoration
    attemptRestore()

    // Retry in next frame if needed
    const restoreTimer = requestAnimationFrame(() => {
      if (!isRestoredRef.current) {
        attemptRestore()
      }
    })

    // Set up scroll listener
    const container = getContainer()

    if (container) {
      container.addEventListener("scroll", handleScroll, { passive: true })
    } else {
      window.addEventListener("scroll", handleScroll, { passive: true })
    }

    // Set up unload listener
    const handleBeforeUnload = () => {
      if (saveOnUnload) {
        if (scrollSaveTimerRef.current) {
          clearTimeout(scrollSaveTimerRef.current)
        }
        saveScrollPosition()
      }
    }

    const handleVisibilityChange = () => {
      if (document.visibilityState === "hidden") {
        saveScrollPosition()
      }
    }

    if (saveOnUnload) {
      window.addEventListener("beforeunload", handleBeforeUnload)
    }

    document.addEventListener("visibilitychange", handleVisibilityChange)

    // Cleanup
    return () => {
      cancelAnimationFrame(restoreTimer)

      if (scrollSaveTimerRef.current) {
        clearTimeout(scrollSaveTimerRef.current)
      }

      if (container) {
        container.removeEventListener("scroll", handleScroll)
      } else {
        window.removeEventListener("scroll", handleScroll)
      }

      if (saveOnUnload) {
        window.removeEventListener("beforeunload", handleBeforeUnload)
      }

      document.removeEventListener("visibilitychange", handleVisibilityChange)
    }
  }, [
    pathname,
    searchParams,
    enabled,
    saveOnUnload,
    attemptRestore,
    getContainer,
    resetContainer,
    cleanupOldPositions,
    handleScroll,
    saveScrollPosition,
  ])

  return {
    saveScrollPosition,
    restoreScrollPosition,
    isRestored: isRestoredRef.current,
  }
}
