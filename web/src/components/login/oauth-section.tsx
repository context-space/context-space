import { SiGithub, SiGoogle } from "@icons-pack/react-simple-icons"
import { useTranslations } from "next-intl"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

interface OAuthSectionProps {
  isPending: boolean
  callbackURL: string
  onOAuthLogin?: (provider: string) => boolean
  onAnonymousLogin?: () => void
}

export default function OAuthSection({ isPending, callbackURL, onOAuthLogin }: OAuthSectionProps) {
  const t = useTranslations()

  const oAuthButtonClass = cn("tracking-wide hover:bg-primary/20 dark:hover:bg-primary/10 border-neutral-200 dark:border-white/5")

  const handleOAuthLogin = async (provider: string) => {
    if (onOAuthLogin && !onOAuthLogin(provider)) {
      return
    }
    const { loginWithOAuth } = await import("@/app/[locale]/login/actions")
    const result = await loginWithOAuth(provider, callbackURL)

    if (result.success && result.redirectTo) {
      window.location.href = result.redirectTo
    }
  }

  // const handleAnonymousLogin = async () => {
  //   const formData = new FormData()
  //   formData.append("from", from)

  //   // 调用服务器操作
  //   const { loginAnonymously } = await import("@/app/[locale]/login/actions")
  //   const result = await loginAnonymously(formData)

  //   if (result.success) {
  //     window.location.href = result.redirectTo || "/"
  //   }
  // }

  return (
    <div className="flex flex-col gap-4 flex-1">
      {/* <Button
        type="button"
        variant="outline"
        onClick={handleAnonymousLogin}
        disabled={isPending}
        className={oAuthButtonClass}
      >
        <span>{t("login.continueAsGuest")}</span>
        <User className="w-4 h-4 opacity-70" />
      </Button> */}

      <Button
        type="button"
        variant="outline"
        onClick={() => handleOAuthLogin("google")}
        disabled={isPending}
        className={oAuthButtonClass}
      >
        <span>{t("login.continueWith", { provider: "Google" })}</span>
        <SiGoogle className="w-4 h-4" />
      </Button>

      <Button
        type="button"
        variant="outline"
        onClick={() => handleOAuthLogin("github")}
        disabled={isPending}
        className={oAuthButtonClass}
      >
        <span>{t("login.continueWith", { provider: "GitHub" })}</span>
        <SiGithub className="w-4 h-4" />
      </Button>
    </div>
  )
}
