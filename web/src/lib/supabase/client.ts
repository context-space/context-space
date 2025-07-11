import { createBrowserClient } from "@supabase/ssr"
import { clientLogger } from "../utils"

const supabaseLogger = clientLogger.withTag("supabase")

export function createClient() {
  const client = createBrowserClient(process.env.NEXT_PUBLIC_SUPABASE_URL!, process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!, {
    auth: {
      autoRefreshToken: true,
      persistSession: true,
      detectSessionInUrl: true,
      storage: {
        getItem: (key) => {
          const value = localStorage.getItem(key)
          return value
        },
        setItem: (key, value) => {
          localStorage.setItem(key, value)
        },
        removeItem: (key) => {
          localStorage.removeItem(key)
        },
      },
    },
  })

  client.auth.onAuthStateChange((event, session) => {
    supabaseLogger.info("Auth state changed", { event, session })
  })

  return client
}

export async function getClientToken() {
  const supabase = createClient()
  const { data: { session } } = await supabase.auth.getSession()
  return session?.access_token
}
