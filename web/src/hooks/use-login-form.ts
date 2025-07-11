import type { FormEvent, RefObject } from "react"
import type { CaptchaSectionRef } from "@/components/login/captcha-section"
import { useRouter } from "next/navigation"
import { useCallback, useEffect, useState } from "react"
import { loginAnonymously, loginWithEmail, loginWithOAuth, sendVerificationCode } from "@/app/[locale]/login/actions"
import { turnstile } from "@/config"

export function useLoginForm(from: string, captchaRef?: RefObject<CaptchaSectionRef | null>) {
  const router = useRouter()
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

    const formData = new FormData()
    formData.append("email", email)
    if (turnstile.enabled && captchaToken) {
      formData.append("captchaToken", captchaToken)
    }

    const result = await sendVerificationCode(formData)

    if (result.success) {
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
    } else {
      // Check if error is related to CAPTCHA verification
      const errorMessage = result.error || ""
      if (errorMessage.includes("captcha verification process failed") || errorMessage.includes("CAPTCHA")) {
        setCaptchaError(true)
        captchaRef?.current?.triggerErrorState()
      } else {
        setErrorWithTimeout(setVerificationCodeError)
      }
      // Reset CAPTCHA on error
      captchaRef?.current?.resetCaptcha()
      setCaptchaToken("")
    }

    setIsSendingCode(false)
  }, [clearErrors, email, termsAccepted, captchaToken, setErrorWithTimeout, captchaRef])

  const handleEmailLogin = useCallback(async (e?: FormEvent<HTMLFormElement>) => {
    e?.preventDefault()

    try {
      if (!email || !verificationCode) {
        return
      }

      setVerificationCodeError(false)

      const formData = new FormData()
      formData.append("email", email)
      formData.append("verificationCode", verificationCode)
      formData.append("from", from)

      const result = await loginWithEmail(formData)

      if (result.success) {
      // 登录成功，直接跳转
        router.push(result.redirectTo || "/")
      } else {
      // 登录失败，显示错误信息
        setErrorWithTimeout(setVerificationCodeError)
        setVerificationCode("")
      }
    } catch (error) {
      console.error("Error logging in with email", error)
      setErrorWithTimeout(setVerificationCodeError)
      setVerificationCode("")
    }
  }, [email, from, setErrorWithTimeout, verificationCode, router])

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

    const formData = new FormData()
    formData.append("provider", provider)
    formData.append("from", from)

    const result = await loginWithOAuth(formData)

    if (result.success) {
      router.push(result.redirectTo || "/")
    } else {
      setErrorWithTimeout(setTermsError)
    }
  }, [termsAccepted, from, setErrorWithTimeout, router])

  const handleAnonymousLogin = useCallback(async () => {
    if (!termsAccepted) {
      setErrorWithTimeout(setTermsError)
      return
    }

    setTermsError(false)

    const formData = new FormData()
    formData.append("from", from)

    const result = await loginAnonymously(formData)

    if (result.success) {
      // 匿名登录成功，直接跳转
      router.push(result.redirectTo || "/")
    } else {
      // 匿名登录失败，显示错误
      setErrorWithTimeout(setTermsError)
    }
  }, [from, router, setErrorWithTimeout, termsAccepted])

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
