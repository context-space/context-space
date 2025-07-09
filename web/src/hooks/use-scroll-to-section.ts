import { useCallback } from "react"

interface UseScrollToSectionOptions {
  offset?: number
  behavior?: ScrollBehavior
  headerSelector?: string
  additionalOffset?: () => number
}

/**
 * Hook for scrolling to sections within a page
 * Supports both regular window scroll and ScrollArea components
 * Dynamically calculates header height for accurate positioning
 */
export function useScrollToSection(options: UseScrollToSectionOptions = {}) {
  const {
    offset,
    behavior = "smooth",
    headerSelector = "header",
    additionalOffset,
  } = options

  return useCallback((sectionId: string) => {
    const element = document.getElementById(sectionId)
    if (!element) return

    // Dynamically calculate offset
    let dynamicOffset = offset || 0

    // If no fixed offset provided, calculate from header
    if (!offset) {
      const headerElement = document.querySelector(headerSelector)
      if (headerElement) {
        dynamicOffset = headerElement.getBoundingClientRect().height
      }
    }

    // Add any additional offset (e.g., for mobile menu)
    if (additionalOffset) {
      dynamicOffset += additionalOffset()
    }

    const scrollContainer = document.querySelector("[data-radix-scroll-area-viewport]")
    const elementPosition = element.getBoundingClientRect().top

    if (scrollContainer) {
      // Scroll within ScrollArea component
      const currentScrollTop = scrollContainer.scrollTop
      const targetScrollTop = currentScrollTop + elementPosition - dynamicOffset

      scrollContainer.scrollTo({
        top: targetScrollTop,
        behavior,
      })
    } else {
      // Fallback to window scroll
      const offsetPosition = elementPosition + window.pageYOffset - dynamicOffset
      window.scrollTo({
        top: offsetPosition,
        behavior,
      })
    }
  }, [offset, behavior, headerSelector, additionalOffset])
}
