"use client"

import type { Session, User } from "@supabase/supabase-js"
import type { AuthContextType } from "@/hooks/use-auth"
import { useRouter } from "next/navigation"
import { useCallback, useEffect, useMemo, useRef, useState } from "react"
import { AuthContext } from "@/hooks/use-auth"
import { createClient, getClientToken } from "@/lib/supabase/client"
import { clientLogger } from "@/lib/utils"

const authProviderLogger = clientLogger.withTag("auth-provider")

// Import the clear function from playground
const clearUserChatData = () => {
  try {
    if (typeof window === "undefined") return

    const keysToRemove: string[] = []

    // Find all localStorage keys related to chat history for this user
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && (key.startsWith("chat_history_") || key.startsWith("input_history_"))) {
        keysToRemove.push(key)
      }
    }

    // Remove all found keys
    keysToRemove.forEach(key => localStorage.removeItem(key))
    authProviderLogger.info("Cleared chat data from localStorage", { keysRemoved: keysToRemove.length })
  } catch (error) {
    authProviderLogger.error("Failed to clear chat data", { error })
  }
}

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
        clearUserChatData()
      } else if (event === "SIGNED_IN") {
        authProviderLogger.info("User signed in, auth state updated", { event, userId: newUser?.id })
      }
    })

    return () => {
      subscription.unsubscribe()
    }
  }, [supabase, initialSession, router])

  const signOut = useCallback(async () => {
    try {
      authProviderLogger.info("Starting sign out process")

      const { error } = await supabase.auth.signOut()
      if (error) {
        authProviderLogger.error("Error signing out", { error })
        throw error
      }

      authProviderLogger.info("Sign out successful")

      // 清理本地状态
      clearUserChatData()

      // 刷新页面，让 middleware 处理后续的重定向逻辑
      setTimeout(() => {
        router.refresh()
      }, 100)
    } catch (error) {
      authProviderLogger.error("Error signing out", { error })
      throw error
    }
  }, [supabase, router])

  const value: AuthContextType = useMemo(() => ({
    user,
    session,
    loading,
    signOut,
    isAuthenticated: !!user,
    isAnonymous: user?.is_anonymous ?? false,
    getClientToken,
  }), [user, session, loading, signOut])

  return (
    <AuthContext value={value}>
      {children}
    </AuthContext>
  )
}
