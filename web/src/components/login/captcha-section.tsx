"use client"

import { Turnstile } from "@marsidev/react-turnstile"
import { useCallback, useEffect, useImperativeHandle, useRef, useState } from "react"
import { cn } from "@/lib/utils"

interface CaptchaSectionProps {
  sitekey: string
  onVerify: (token: string) => void
  onExpire?: () => void
  onError?: (error: string) => void
  className?: string
  theme?: "light" | "dark"
  showError?: boolean
  hasError?: boolean
}

export interface CaptchaSectionRef {
  resetCaptcha: () => void
  execute: () => void
  triggerErrorState: () => void
}

function CaptchaSection({
  ref,
  sitekey,
  onVerify,
  onExpire,
  onError,
  className,
  theme = "light",
  showError = false,
  hasError = false,
}: CaptchaSectionProps & { ref?: React.RefObject<CaptchaSectionRef | null> }) {
  const captchaRef = useRef<any>(null)
  const [isShaking, setIsShaking] = useState(false)

  useImperativeHandle(ref, () => ({
    resetCaptcha: () => {
      captchaRef.current?.reset()
    },
    execute: () => {
      captchaRef.current?.execute()
    },
    triggerErrorState: () => {
      setIsShaking(true)
      setTimeout(() => setIsShaking(false), 600)
    },
  }))

  // Trigger shake animation when hasError changes to true
  useEffect(() => {
    if (hasError) {
      setIsShaking(true)
      setTimeout(() => setIsShaking(false), 600)
    }
  }, [hasError])

  const handleVerify = useCallback((token: string) => {
    onVerify(token)
  }, [onVerify])

  const handleExpire = useCallback(() => {
    onExpire?.()
  }, [onExpire])

  const handleError = useCallback((error: string) => {
    onError?.(error)
  }, [onError])

  return (
    <div className={cn("flex justify-center", className)}>
      <div
        className={cn(
          "relative transition-all duration-200",
          hasError && "p-1 rounded-lg border-2 border-red-500 dark:border-red-400 bg-red-50 dark:bg-red-950/20",
          isShaking && "animate-shake",
        )}
      >
        <Turnstile
          ref={captchaRef}
          siteKey={sitekey}
          onSuccess={handleVerify}
          onExpire={handleExpire}
          onError={handleError}
          options={{
            theme,
            size: "normal",
          }}
        />
      </div>
    </div>
  )
}

CaptchaSection.displayName = "CaptchaSection"

export default CaptchaSection
