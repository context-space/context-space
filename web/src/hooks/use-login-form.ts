import type { FormEvent, RefObject } from "react"
import type { CaptchaSectionRef } from "@/components/login/captcha-section"
import { useCallback, useEffect, useState, useTransition } from "react"
import { loginAnonymously, loginWithEmail, loginWithOAuth, sendVerificationCode } from "@/app/[locale]/login/actions"
import { turnstile } from "@/config"

export function useLoginForm(from: string, captchaRef?: RefObject<CaptchaSectionRef | null>) {
  const [email, setEmail] = useState("")
  const [verificationCode, setVerificationCode] = useState("")
  const [codeSent, setCodeSent] = useState(false)
  const [countdown, setCountdown] = useState(0)
  const [termsAccepted, setTermsAccepted] = useState(false)
  const [emailError, setEmailError] = useState(false)
  const [termsError, setTermsError] = useState(false)
  const [verificationCodeError, setVerificationCodeError] = useState(false)
  const [captchaToken, setCaptchaToken] = useState("")
  const [captchaError, setCaptchaError] = useState(false)

  const [isPending, startTransition] = useTransition()
  const [isSendingCode, setIsSendingCode] = useState(false)

  const clearErrors = useCallback(() => {
    setEmailError(false)
    setTermsError(false)
    setVerificationCodeError(false)
    setCaptchaError(false)
  }, [])

  const setErrorWithTimeout = useCallback((errorSetter: (value: boolean) => void) => {
    errorSetter(true)
    setTimeout(() => errorSetter(false), 1000)
  }, [])

  const handleSendCode = useCallback(async () => {
    clearErrors()

    if (!email || !/^[^\s@]+@[^\s@][^\s.@]*\.[^\s@]+$/.test(email)) {
      setErrorWithTimeout(setEmailError)
      return
    }

    if (!termsAccepted) {
      setErrorWithTimeout(setTermsError)
      return
    }

    if (turnstile.enabled && !captchaToken) {
      // CAPTCHA verification required when enabled
      setCaptchaError(true)
      captchaRef?.current?.triggerErrorState()
      return
    }

    setIsSendingCode(true)

    try {
      const formData = new FormData()
      formData.append("email", email)
      if (turnstile.enabled && captchaToken) {
        formData.append("captchaToken", captchaToken)
      }

      await sendVerificationCode(formData)

      setCodeSent(true)
      setCountdown(60)
      clearErrors()

      // Reset CAPTCHA after successful code send
      captchaRef?.current?.resetCaptcha()
      setCaptchaToken("")

      const timer = setInterval(() => {
        setCountdown((prev) => {
          if (prev <= 1) {
            clearInterval(timer)
            return 0
          }
          return prev - 1
        })
      }, 1000)
    } catch (error) {
      // Check if error is related to CAPTCHA verification
      const errorMessage = error instanceof Error ? error.message : String(error)
      if (errorMessage.includes("captcha verification process failed") || errorMessage.includes("CAPTCHA")) {
        setCaptchaError(true)
        captchaRef?.current?.triggerErrorState()
      } else {
        setErrorWithTimeout(setVerificationCodeError)
      }
      // Reset CAPTCHA on error
      captchaRef?.current?.resetCaptcha()
      setCaptchaToken("")
    } finally {
      setIsSendingCode(false)
    }
  }, [clearErrors, email, termsAccepted, captchaToken, setErrorWithTimeout, captchaRef])

  const handleEmailLogin = useCallback(async (e?: FormEvent<HTMLFormElement>) => {
    e?.preventDefault()

    if (!email || !verificationCode) {
      return
    }

    setVerificationCodeError(false)

    startTransition(async () => {
      try {
        const formData = new FormData()
        formData.append("email", email)
        formData.append("verificationCode", verificationCode)
        formData.append("from", from)

        await loginWithEmail(formData)
      } catch (error) {
        // 忽略 Next.js 重定向错误，这是正常的
        if (error instanceof Error && error.message === "NEXT_REDIRECT") {
          throw error
        }
        setErrorWithTimeout(setVerificationCodeError)
        setVerificationCode("")
      }
    })
  }, [email, from, setErrorWithTimeout, verificationCode])

  useEffect(() => {
    if (verificationCode.length === 6 && email && termsAccepted && codeSent) {
      handleEmailLogin()
    }
  }, [verificationCode, email, termsAccepted, codeSent, handleEmailLogin])

  const handleVerificationCodeChange = useCallback((value: string) => {
    setVerificationCode(value)
    if (verificationCodeError) setVerificationCodeError(false)
  }, [verificationCodeError])

  const handleOAuthLogin = useCallback(async (provider: string) => {
    if (!termsAccepted) {
      setErrorWithTimeout(setTermsError)
      return
    }

    setTermsError(false)

    startTransition(async () => {
      try {
        const formData = new FormData()
        formData.append("provider", provider)
        formData.append("from", from)

        await loginWithOAuth(formData)
      } catch (error) {
        // 忽略 Next.js 重定向错误，这是正常的
        if (error instanceof Error && error.message === "NEXT_REDIRECT") {
          throw error
        }
        // Handle other errors silently
      }
    })
  }, [termsAccepted, from, setErrorWithTimeout])

  const handleAnonymousLogin = useCallback(async () => {
    if (!termsAccepted) {
      setErrorWithTimeout(setTermsError)
      return
    }

    setTermsError(false)

    startTransition(async () => {
      try {
        const formData = new FormData()
        formData.append("from", from)

        await loginAnonymously(formData)
      } catch (error) {
        // 忽略 Next.js 重定向错误，这是正常的
        if (error instanceof Error && error.message === "NEXT_REDIRECT") {
          throw error
        }
        // Handle other errors silently
      }
    })
  }, [from, setErrorWithTimeout, termsAccepted])

  const handleEmailChange = useCallback((value: string) => {
    setEmail(value)
    if (emailError) setEmailError(false)
  }, [emailError])

  const handleTermsChange = useCallback((checked: boolean) => {
    setTermsAccepted(checked)
    if (termsError) setTermsError(false)
  }, [termsError])

  const handleCaptchaVerify = useCallback((token: string) => {
    setCaptchaToken(token)
    // Clear captcha error when successfully verified
    if (captchaError) {
      setCaptchaError(false)
    }
  }, [captchaError])

  const handleCaptchaExpire = useCallback(() => {
    setCaptchaToken("")
  }, [])

  const resetCaptcha = useCallback(() => {
    captchaRef?.current?.resetCaptcha()
    setCaptchaToken("")
    setCaptchaError(false)
  }, [captchaRef])

  return {
    // State
    email,
    verificationCode,
    codeSent,
    countdown,
    termsAccepted,
    captchaToken,
    emailError,
    termsError,
    verificationCodeError,
    captchaError,
    isPending,
    isSendingCode,

    // Actions
    handleSendCode,
    handleEmailLogin,
    handleVerificationCodeChange,
    handleOAuthLogin,
    handleAnonymousLogin,
    handleEmailChange,
    handleTermsChange,
    handleCaptchaVerify,
    handleCaptchaExpire,
    resetCaptcha,
  }
}
