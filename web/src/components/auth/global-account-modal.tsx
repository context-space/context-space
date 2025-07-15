"use client"

import { AccountModal } from "@/components/auth/account-modal"
import { useAccountModal } from "@/hooks/use-account-modal"
import { useAuth } from "@/hooks/use-auth"

export function GlobalAccountModal() {
  const { user, isAnonymous } = useAuth()
  const { isOpen, activeTab, closeModal } = useAccountModal()

  // Only render if user is authenticated and not anonymous
  if (!user || isAnonymous) {
    return null
  }

  const userDisplayName = user.user_metadata?.full_name
    || user.user_metadata?.name
    || user.email?.split("@")[0]
    || "User"

  const userAvatar = user.user_metadata?.avatar_url
    || user.user_metadata?.picture
    || undefined

  return (
    <AccountModal
      user={{
        id: user.id,
        email: user.email || "",
        displayName: userDisplayName,
        avatar: userAvatar,
        isAnonymous,
        createdAt: user.created_at,
        lastSignInAt: user.last_sign_in_at,
        providers: user.app_metadata?.providers || [],
      }}
      open={isOpen}
      onOpenChange={closeModal}
      initialTab={activeTab}
    />
  )
}
