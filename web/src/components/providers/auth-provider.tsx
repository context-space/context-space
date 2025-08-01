"use client"

import type { Session, User } from "@supabase/supabase-js"
import { useRouter } from "next/navigation"
import { useCallback, useEffect, useMemo, useRef, useState } from "react"
import { AuthContext } from "@/hooks/use-auth"
import { clearAllPlaygroundData } from "@/lib/playground-storage"
import { createClient, getClientToken } from "@/lib/supabase/client"
import { clientLogger } from "@/lib/utils"

const authProviderLogger = clientLogger.withTag("auth-provider")

// 使用集中的存储清除函数

interface AuthProviderProps {
  children: React.ReactNode
  initialSession?: Session | null
}

export function AuthProvider({ children, initialSession }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(initialSession?.user ?? null)
  const [session, setSession] = useState<Session | null>(initialSession ?? null)
  const [loading, setLoading] = useState(!initialSession)
  const router = useRouter()

  // Use ref to track previous user state to avoid stale closure issues
  const previousUserRef = useRef<User | null>(initialSession?.user ?? null)

  const supabase = useMemo(() => createClient(), [])

  // Cache derived states separately
  const isAuthenticated = useMemo(() => !!user, [user])
  const isAnonymous = useMemo(() => user?.is_anonymous ?? false, [user?.is_anonymous])

  useEffect(() => {
    if (!initialSession) {
      const getSession = async () => {
        const { data: { session } } = await supabase.auth.getSession()
        setSession(session)
        setUser(session?.user ?? null)
        previousUserRef.current = session?.user ?? null
        setLoading(false)
      }
      getSession()
    }

    const { data: { subscription } } = supabase.auth.onAuthStateChange(async (event, session) => {
      authProviderLogger.info("Auth state changed", { event, userId: session?.user?.id })

      const newUser = session?.user ?? null

      setSession(session)
      setUser(newUser)
      setLoading(false)

      // Update ref for next comparison
      previousUserRef.current = newUser

      // Clear localStorage when user logs out
      if (event === "SIGNED_OUT") {
        authProviderLogger.info("User signed out, clearing data")
        clearAllPlaygroundData()
      } else if (event === "SIGNED_IN") {
        authProviderLogger.info("User signed in, auth state updated", { event, userId: newUser?.id })
      }
    })

    return () => {
      subscription.unsubscribe()
    }
  }, [supabase, initialSession, router])

  const signOut = useCallback(async (): Promise<void> => {
    try {
      authProviderLogger.info("Signing out user", { userId: user?.id })

      // 清理本地状态
      clearAllPlaygroundData()

      // 调用 Supabase 客户端登出 - 这会：
      // 1. 清除本地存储的会话信息
      // 2. 向 Supabase 发送登出请求使 token 失效
      // 3. 触发 auth state change 事件
      const { error } = await supabase.auth.signOut()

      if (error) {
        authProviderLogger.error("Supabase signOut error", { error })
        // 即使登出失败，也继续流程，因为我们已经清理了本地状态
      } else {
        authProviderLogger.info("User successfully signed out")
      }

      // 刷新页面，让 middleware 检测到会话失效并处理重定向
      setTimeout(() => {
        router.refresh()
      }, 100)
    } catch (error) {
      authProviderLogger.error("Error signing out", { error })
      // 即使出错也要清理状态并刷新
      setTimeout(() => {
        router.refresh()
      }, 100)
    }
  }, [user?.id, router, supabase])

  const getToken = useCallback(async () => {
    return await getClientToken()
  }, [])

  const value = useMemo(
    () => ({
      user,
      session,
      isAuthenticated,
      isAnonymous,
      loading,
      signOut,
      getClientToken: getToken,
    }),
    [user, session, isAuthenticated, isAnonymous, loading, signOut, getToken],
  )

  return <AuthContext value={value}>{children}</AuthContext>
}
