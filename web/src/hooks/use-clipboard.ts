import { useCallback } from "react"
import { toast } from "sonner"

/**
 * Hook for clipboard copy functionality
 * Returns a copy function - state management is left to the consuming component
 */
export const useCopyToClipboard = () => {
  const copyToClipboard = useCallback(async (text: string, successMessage?: string, errorMessage?: string) => {
    try {
      await navigator.clipboard.writeText(text)

      if (successMessage) {
        toast.success(successMessage)
      }

      return true // Return success
    } catch (error) {
      console.error("Failed to copy to clipboard:", error)

      if (errorMessage) {
        toast.error(errorMessage)
      }

      return false // Return failure
    }
  }, [])

  return copyToClipboard
}
