"use server"

import type { Provider } from "@supabase/supabase-js"
import type { LoginResult, SendCodeResult } from "@/typings/auth"
import { revalidatePath } from "next/cache"
import { z } from "zod"
import { baseURL, turnstile } from "@/config"
import { createClient } from "@/lib/supabase/server"
import { serverLogger } from "@/lib/utils"

const loginLogger = serverLogger.withTag("login")

// 验证模式
const emailSchema = z.object({
  email: z.string().email("Invalid email address"),
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

export async function loginWithOAuth(provider: string, callbackUrl: string): Promise<LoginResult> {
  const validProviders = ["google", "github"]
  if (!provider || !validProviders.includes(provider)) {
    return { success: false, error: "Invalid provider" }
  }

  try {
    const supabase = await createClient()
    const { data, error } = await supabase.auth.signInWithOAuth({
      provider: provider as Provider,
      options: {
        redirectTo: callbackUrl,
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

export async function sendMagicLinkAction(
  prevState: { success: boolean, sent: boolean, errors: Record<string, string>, email?: string },
  formData: FormData,
): Promise<{ success: boolean, sent: boolean, errors: Record<string, string>, email?: string }> {
  const email = formData.get("email") as string
  const captchaToken = formData.get("captchaToken") as string
  const termsAccepted = formData.get("termsAccepted") === "true"
  const callbackURL = formData.get("callbackURL") as string

  loginLogger.info("Magic Link request initiated", {
    email,
    hasEmail: !!email,
    hasCaptchaToken: !!captchaToken,
    termsAccepted,
    callbackURL,
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

    loginLogger.info("Sending Magic Link via Supabase", {
      email,
      redirectTo: callbackURL,
      hasCaptcha: !!captchaToken,
      baseURL,
      timestamp: new Date().toISOString(),
    })

    const otpOptions: any = {
      shouldCreateUser: true,
      emailRedirectTo: callbackURL,
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
      redirectTo: callbackURL,
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
