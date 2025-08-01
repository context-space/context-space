"use client"

import { useLocale } from "next-intl"
import { useSearchParams } from "next/navigation"
import { useMemo } from "react"
import { FramelessLayout } from "@/components/layouts"
import { Login } from "@/components/login"
import { getPathname } from "@/i18n/navigation"

export default function LoginPage() {
  const searchParams = useSearchParams()
  const from = searchParams.get("from")
  const locale = useLocale()

  const callbackURL = useMemo(() => window.location.origin + getPathname({
    href: {
      pathname: "/auth-callback",
      query: {
        redirect_to: from ?? "/",
      },
    },
    locale,
  }), [from, locale])

  return (
    <FramelessLayout>
      <div className="min-h-[calc(100vh-20rem)] flex items-center justify-center px-4">
        <Login callbackURL={callbackURL} className="pt-16 w-[420px] max-w-[90vw]" />
      </div>
    </FramelessLayout>
  )
}
