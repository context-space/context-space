"use server"

import { revalidatePath } from "next/cache"
import { redirect } from "next/navigation"
import { createClient } from "@/lib/supabase/server"
import { serverLogger } from "@/lib/utils"

const logoutLogger = serverLogger.withTag("logout")
export async function logout() {
  try {
    const supabase = await createClient()

    const { error } = await supabase.auth.signOut()

    if (error) {
      logoutLogger.error("Logout error", { error })
      throw new Error("Failed to log out")
    }

    revalidatePath("/", "layout")
    redirect("/login")
  } catch (error) {
    logoutLogger.error("Logout error", { error })
    if (error instanceof Error && error.message === "NEXT_REDIRECT") {
      throw error // 让redirect正常工作
    }
    throw new Error("Logout failed")
  }
}
