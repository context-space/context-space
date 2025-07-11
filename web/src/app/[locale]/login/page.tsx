import type { Locale } from "@/i18n/routing"
import { FramelessLayout } from "@/components/layouts"
import { Login } from "@/components/login"

interface LoginPageProps {
  params: Promise<{ locale: Locale }>
  searchParams: Promise<{ from?: string, error?: string }>
}

export default async function LoginPage({ searchParams }: LoginPageProps) {
  const awaitedSearchParams = await searchParams
  const from = awaitedSearchParams?.from ? decodeURIComponent(awaitedSearchParams?.from) : "/"

  return (
    <FramelessLayout>
      <div className="min-h-[calc(100vh-20rem)] flex items-center justify-center px-4">
        <Login from={from} className="pt-16 w-[420px] max-w-[90vw]" />
      </div>
    </FramelessLayout>
  )
}
