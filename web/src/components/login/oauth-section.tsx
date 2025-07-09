import { SiGithub, SiGoogle } from "@icons-pack/react-simple-icons"
import { useTranslations } from "next-intl"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

interface OAuthSectionProps {
  isPending: boolean
  onOAuthLogin: (provider: string) => void
  onAnonymousLogin: () => void
}

export default function OAuthSection({ isPending, onOAuthLogin, onAnonymousLogin }: OAuthSectionProps) {
  const t = useTranslations()

  const oAuthButtonClass = cn("tracking-wide hover:bg-primary/20 dark:hover:bg-primary/10 border-neutral-200 dark:border-white/5")

  return (
    <div className="flex flex-col gap-4 flex-1">
      {/* <Button
        type="button"
        variant="outline"
        onClick={onAnonymousLogin}
        disabled={isPending}
        className={oAuthButtonClass}
      >
        <span>{t("login.continueAsGuest")}</span>
        <User className="w-4 h-4 opacity-70" />
      </Button> */}

      <Button
        type="button"
        variant="outline"
        onClick={() => onOAuthLogin("google")}
        disabled={isPending}
        className={oAuthButtonClass}
      >
        <span>{t("login.continueWith", { provider: "Google" })}</span>
        <SiGoogle className="w-4 h-4" />
      </Button>

      <Button
        type="button"
        variant="outline"
        onClick={() => onOAuthLogin("github")}
        disabled={isPending}
        className={oAuthButtonClass}
      >
        <span>{t("login.continueWith", { provider: "GitHub" })}</span>
        <SiGithub className="w-4 h-4" />
      </Button>
    </div>
  )
}
