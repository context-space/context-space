import type { ReactNode } from "react"
import type { Locale } from "@/i18n/routing"
import { hasLocale, NextIntlClientProvider } from "next-intl"
import { getMessages } from "next-intl/server"
import { ThemeProvider } from "next-themes"
import { notFound } from "next/navigation"
import { GlobalAccountModal } from "@/components/auth/global-account-modal"
import { ScrollRestorationManager } from "@/components/common/scroll-restoration-manager"
import { AuthProvider } from "@/components/providers/auth-provider"
import { QueryProvider } from "@/components/providers/query-provider"
import { Toaster } from "@/components/ui/sonner"
import { routing } from "@/i18n/routing"
import { createClient } from "@/lib/supabase/server"

interface Props {
  children: ReactNode
  params: Promise<{ locale: Locale }>
}

export default async function LocaleLayout({ children, params }: Props) {
  const { locale } = await params
  if (!hasLocale(routing.locales, locale)) {
    notFound()
  }

  const messages = await getMessages()

  const supabase = await createClient()
  const { data: { session } } = await supabase.auth.getSession()

  return (
    <NextIntlClientProvider messages={messages}>
      <ThemeProvider attribute="class" defaultTheme="system" storageKey="context-space-theme">
        <QueryProvider>
          <AuthProvider initialSession={session}>
            {children}
            <GlobalAccountModal />
            <ScrollRestorationManager />
            <Toaster
              offset={{
                top: 5,
              }}
              mobileOffset={{
                top: 5,
              }}
              richColors
              position="top-center"
              duration={3000}
            />
          </AuthProvider>
        </QueryProvider>
      </ThemeProvider>
    </NextIntlClientProvider>
  )
}
