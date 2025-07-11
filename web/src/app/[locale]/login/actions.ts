"use server"

import type { Provider } from "@supabase/supabase-js"
import type { LoginResult, SendCodeResult } from "@/typings/auth"
import { revalidatePath } from "next/cache"
import { z } from "zod"
import { baseURL, turnstile } from "@/config"
import { defaultLocale, locales } from "@/i18n/routing"
import { createClient } from "@/lib/supabase/server"
import { serverLogger } from "@/lib/utils"

const loginLogger = serverLogger.withTag("login")
function extractLocaleFromPath(redirectTo: string): string {
  const pathSegments = redirectTo.replace(/^\//, "").split("/")
  const firstSegment = pathSegments[0]

  if (locales.includes(firstSegment as typeof locales[number])) {
    return firstSegment
  }

  return defaultLocale
}

// 验证模式
const emailSchema = z.object({
  email: z.string().email("Invalid email address"),
})

const loginSchema = z.object({
  email: z.string().email("Invalid email address"),
  verificationCode: z.string().regex(/^\d{6}$/, "Verification code must be 6 digits"),
})

// 发送验证码
export async function sendVerificationCode(formData: FormData): Promise<SendCodeResult> {
  const email = formData.get("email") as string
  const captchaToken = formData.get("captchaToken") as string

  // 验证输入
  const validation = emailSchema.safeParse({ email })
  if (!validation.success) {
    return { success: false, error: "Invalid email address" }
  }

  try {
    const supabase = await createClient()

    loginLogger.info("Sending OTP verification code", {
      email: validation.data.email,
      hasCaptcha: !!captchaToken,
      timestamp: new Date().toISOString(),
    })

    const otpOptions: any = {
      shouldCreateUser: true,
    }

    // Only add captchaToken if it exists
    if (captchaToken) {
      otpOptions.captchaToken = captchaToken
    }

    const { error } = await supabase.auth.signInWithOtp({
      email: validation.data.email,
      options: otpOptions,
    })

    if (error) {
      loginLogger.error("Supabase OTP send error", {
        error: error.message,
        email: validation.data.email,
      })
      return { success: false, error: error.message }
    }

    loginLogger.info("OTP verification code sent successfully", {
      email: validation.data.email,
      timestamp: new Date().toISOString(),
    })

    return { success: true, message: "Verification code sent successfully" }
  } catch (error) {
    loginLogger.error("Send verification code error", { error })
    return { success: false, error: "Failed to send verification code" }
  }
}

// 邮箱验证码登录
export async function loginWithEmail(formData: FormData): Promise<LoginResult> {
  const email = formData.get("email") as string
  const verificationCode = formData.get("verificationCode") as string

  // 验证输入
  const validation = loginSchema.safeParse({ email, verificationCode })
  if (!validation.success) {
    const firstError = validation.error.errors[0]
    return { success: false, error: firstError.message }
  }

  try {
    const supabase = await createClient()

    loginLogger.info("Starting OTP verification", {
      email: validation.data.email,
      timestamp: new Date().toISOString(),
    })

    const { data, error } = await supabase.auth.verifyOtp({
      email: validation.data.email,
      token: validation.data.verificationCode,
      type: "email",
    })

    if (error) {
      loginLogger.error("OTP verification failed", {
        error: error.message,
        email: validation.data.email,
      })

      if (error.message.includes("expired")) {
        return { success: false, error: "验证码已过期，请重新获取验证码" }
      } else if (error.message.includes("invalid")) {
        return { success: false, error: "验证码无效，请检查输入是否正确" }
      } else {
        return { success: false, error: `验证失败: ${error.message}` }
      }
    }

    if (!data.user) {
      loginLogger.error("OTP verification succeeded but no user returned")
      return { success: false, error: "登录验证成功但用户信息获取失败" }
    }

    loginLogger.info("OTP verification successful", data)

    revalidatePath("/", "layout")

    // 返回成功状态和重定向URL，让客户端处理重定向
    const from = formData.get("from") as string || "/"
    return { success: true, redirectTo: from }
  } catch (error) {
    loginLogger.error("Email login error", { error })
    return { success: false, error: "登录失败，请重试" }
  }
}

// OAuth登录
export async function loginWithOAuth(formData: FormData): Promise<LoginResult> {
  const provider = formData.get("provider") as string
  const from = formData.get("from") as string || "/"

  const validProviders = ["google", "github"]
  if (!provider || !validProviders.includes(provider)) {
    return { success: false, error: "Invalid provider" }
  }

  try {
    const supabase = await createClient()

    const locale = extractLocaleFromPath(from)
    const callbackPath = locale === defaultLocale ? "/auth-callback" : `/${locale}/auth-callback`
    const callbackUrl = new URL(callbackPath, baseURL)
    if (from !== "/") {
      callbackUrl.searchParams.set("redirect_to", from)
    }

    const { data, error } = await supabase.auth.signInWithOAuth({
      provider: provider as Provider,
      options: {
        redirectTo: callbackUrl.toString(),
      },
    })

    if (error) {
      loginLogger.error("OAuth login error", { error })
      return { success: false, error: error.message }
    }

    if (data.url) {
      return { success: true, redirectTo: data.url }
    } else {
      return { success: false, error: "Failed to get OAuth URL" }
    }
  } catch (error) {
    loginLogger.error("OAuth login error", { error })
    return { success: false, error: "OAuth login failed" }
  }
}

// 匿名登录
export async function loginAnonymously(formData: FormData): Promise<LoginResult> {
  const from = formData.get("from") as string || "/"

  try {
    const supabase = await createClient()

    const { data, error } = await supabase.auth.signInAnonymously()

    if (error || !data.user) {
      loginLogger.error("Anonymous login error", { error })
      return { success: false, error: "Anonymous login failed" }
    }

    loginLogger.info("Anonymous login successful", {
      userId: data.user.id,
      timestamp: new Date().toISOString(),
    })

    revalidatePath("/", "layout")
    return { success: true, redirectTo: from }
  } catch (error) {
    loginLogger.error("Anonymous login error", { error })
    return { success: false, error: "Anonymous login failed" }
  }
}

// 修改：发送Magic Link操作
export async function sendMagicLinkAction(
  prevState: { success: boolean, sent: boolean, errors: Record<string, string>, email?: string },
  formData: FormData,
): Promise<{ success: boolean, sent: boolean, errors: Record<string, string>, email?: string }> {
  const email = formData.get("email") as string
  const captchaToken = formData.get("captchaToken") as string
  const termsAccepted = formData.get("termsAccepted") === "true"
  const from = formData.get("from") as string || "/"

  loginLogger.info("Magic Link request initiated", {
    email,
    hasEmail: !!email,
    hasCaptchaToken: !!captchaToken,
    termsAccepted,
    from,
    turnstileEnabled: turnstile.enabled,
    timestamp: new Date().toISOString(),
  })

  const errors: Record<string, string> = {}

  // 验证输入
  if (!email || !/^[^\s@]+@[^\s@][^\s.@]*\.[^\s@]+$/.test(email)) {
    loginLogger.warn("Invalid email address provided", { email: email ? "***@***" : null })
    errors.email = "Invalid email address"
  }

  if (!termsAccepted) {
    loginLogger.warn("Terms not accepted")
    errors.terms = "Please accept terms"
  }

  if (turnstile.enabled && !captchaToken) {
    loginLogger.warn("CAPTCHA token missing when required")
    errors.captcha = "CAPTCHA verification required"
  }

  if (Object.keys(errors).length > 0) {
    loginLogger.info("Magic Link request validation failed", {
      errors: Object.keys(errors),
      errorCount: Object.keys(errors).length,
    })
    return { success: false, sent: false, errors }
  }

  try {
    const supabase = await createClient()

    // 构建Magic Link重定向URL
    const locale = extractLocaleFromPath(from)
    const callbackPath = locale === defaultLocale ? "/auth-callback" : `/${locale}/auth-callback`
    const callbackUrl = new URL(callbackPath, baseURL)
    if (from !== "/") {
      callbackUrl.searchParams.set("redirect_to", from)
    }

    loginLogger.info("Sending Magic Link via Supabase", {
      email,
      redirectTo: callbackUrl.toString(),
      locale,
      callbackPath,
      hasCaptcha: !!captchaToken,
      baseURL,
      timestamp: new Date().toISOString(),
    })

    const otpOptions: any = {
      shouldCreateUser: true,
      emailRedirectTo: callbackUrl.toString(),
    }

    // Only add captchaToken if it exists
    if (captchaToken) {
      loginLogger.info("Adding CAPTCHA token to OTP options")
      otpOptions.captchaToken = captchaToken
    }

    loginLogger.info("Calling supabase.auth.signInWithOtp", {
      email,
      optionsKeys: Object.keys(otpOptions),
      hasCaptchaInOptions: !!otpOptions.captchaToken,
    })

    const { error } = await supabase.auth.signInWithOtp({
      email,
      options: otpOptions,
    })

    loginLogger.info("Supabase signInWithOtp response received", {
      hasError: !!error,
      errorCode: error?.code,
      errorMessage: error?.message,
      timestamp: new Date().toISOString(),
    })

    if (error) {
      loginLogger.error("Supabase Magic Link send error", {
        error: error.message,
        errorCode: error.code,
        email,
        fullError: error,
        timestamp: new Date().toISOString(),
      })

      const errorMessage = error.message || ""
      if (errorMessage.includes("captcha verification process failed") || errorMessage.includes("CAPTCHA")) {
        loginLogger.warn("CAPTCHA verification failed during Magic Link send")
        errors.captcha = "CAPTCHA verification failed"
      } else if (errorMessage.includes("over_email_send_rate_limit")) {
        loginLogger.warn("Email send rate limit exceeded", { email })
        errors.email = "Too many emails sent. Please wait before trying again."
      } else if (errorMessage.includes("signup_disabled")) {
        loginLogger.warn("Signup is disabled", { email })
        errors.email = "Account creation is currently disabled"
      } else {
        loginLogger.error("Unknown Magic Link send error", { errorMessage, errorCode: error.code })
        errors.email = "Failed to send magic link"
      }
      return { success: false, sent: false, errors }
    }

    loginLogger.info("Magic Link sent successfully", {
      email,
      redirectTo: callbackUrl.toString(),
      timestamp: new Date().toISOString(),
    })

    return {
      success: true,
      sent: true,
      errors: {},
      email,
    }
  } catch (error) {
    loginLogger.error("Send magic link error - unexpected exception", {
      error,
      errorMessage: error instanceof Error ? error.message : String(error),
      errorStack: error instanceof Error ? error.stack : undefined,
      email,
      timestamp: new Date().toISOString(),
    })
    errors.email = "Failed to send magic link"
    return { success: false, sent: false, errors }
  }
}

// 移除邮箱登录操作，因为Magic Link会自动处理登录
