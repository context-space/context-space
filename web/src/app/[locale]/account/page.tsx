import { redirect } from "next/navigation"
import { BaseLayout } from "@/components/layouts"
import { createClient } from "@/lib/supabase/server"
import { AccountPageContent } from "./page-content"

export default async function AccountPage() {
  const supabase = await createClient()

  const { data: { user }, error } = await supabase.auth.getUser()

  if (error || !user) {
    redirect("/login")
  }

  const isAnonymous = user.is_anonymous ?? false
  const userDisplayName = isAnonymous
    ? "Anonymous"
    : (
        user.user_metadata?.full_name
        || user.user_metadata?.name
        || user.user_metadata?.user_name
        || user.email?.split("@")[0]
        || "User"
      )

  const userAvatar = isAnonymous
    ? "/logo.svg"
    : (
        user.user_metadata?.avatar_url
        || user.user_metadata?.picture
      )

  const userEmail = isAnonymous ? "anonymous@example.com" : user.email

  return (
    <BaseLayout>
      <AccountPageContent
        user={{
          id: user.id,
          email: userEmail || "",
          displayName: userDisplayName,
          avatar: userAvatar,
          isAnonymous,
          createdAt: user.created_at,
          lastSignInAt: user.last_sign_in_at,
          providers: user.app_metadata?.providers || [],
        }}
      />
    </BaseLayout>
  )
}
