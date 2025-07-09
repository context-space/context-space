"use server"

import type { Provider } from "@supabase/supabase-js"
import { revalidatePath } from "next/cache"
import { redirect } from "next/navigation"
import { z } from "zod"
import { baseURL } from "@/config"
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
export async function sendVerificationCode(formData: FormData) {
  const email = formData.get("email") as string
  const captchaToken = formData.get("captchaToken") as string

  // 验证输入
  const validation = emailSchema.safeParse({ email })
  if (!validation.success) {
    throw new Error("Invalid email address")
  }

  try {
    const supabase = await createClient()

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
      loginLogger.error("Supabase auth error", { error })
      throw new Error(error.message)
    }

    return { success: true, message: "Verification code sent successfully" }
  } catch (error) {
    loginLogger.error("Send verification code error", { error })
    throw new Error("Failed to send verification code")
  }
}

// 邮箱验证码登录
export async function loginWithEmail(formData: FormData) {
  const email = formData.get("email") as string
  const verificationCode = formData.get("verificationCode") as string

  // 验证输入
  const validation = loginSchema.safeParse({ email, verificationCode })
  if (!validation.success) {
    const firstError = validation.error.errors[0]
    throw new Error(firstError.message)
  }

  try {
    const supabase = await createClient()

    const { data, error } = await supabase.auth.verifyOtp({
      email: validation.data.email,
      token: validation.data.verificationCode,
      type: "email",
    })

    if (error || !data.user) {
      loginLogger.error("OTP verification error", { error })
      throw new Error("Invalid verification code")
    }

    revalidatePath("/", "layout")

    // 获取重定向URL
    const from = formData.get("from") as string || "/"
    redirect(from)
  } catch (error) {
    if (error instanceof Error && error.message === "NEXT_REDIRECT") {
      throw error // 让redirect正常工作
    }
    loginLogger.error("Email login error", { error })
    throw new Error("Login failed. Please try again.")
  }
}

// OAuth登录
export async function loginWithOAuth(formData: FormData) {
  const provider = formData.get("provider") as string
  const from = formData.get("from") as string || "/"

  const validProviders = ["google", "github"]
  if (!provider || !validProviders.includes(provider)) {
    throw new Error("Invalid provider")
  }

  try {
    const supabase = await createClient()

    const locale = extractLocaleFromPath(from)
    const callbackPath = locale === defaultLocale ? "/oauth-callback" : `/${locale}/oauth-callback`
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
      throw new Error(error.message)
    }

    if (data.url) {
      redirect(data.url)
    }
  } catch (error) {
    if (error instanceof Error && error.message === "NEXT_REDIRECT") {
      throw error // 让redirect正常工作
    }
    loginLogger.error("OAuth login error", { error })
    throw new Error("OAuth login failed")
  }
}

// 匿名登录
export async function loginAnonymously(formData: FormData) {
  const from = formData.get("from") as string || "/"

  try {
    const supabase = await createClient()

    const { data, error } = await supabase.auth.signInAnonymously()

    if (error || !data.user) {
      loginLogger.error("Anonymous login error", { error })
      throw new Error("Anonymous login failed")
    }

    revalidatePath("/", "layout")
    redirect(from)
  } catch (error) {
    if (error instanceof Error && error.message === "NEXT_REDIRECT") {
      throw error // 让redirect正常工作
    }
    loginLogger.error("Anonymous login error", { error })
    throw new Error("Anonymous login failed")
  }
}
