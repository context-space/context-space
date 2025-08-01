"use client"

import type { CaptchaSectionRef } from "./captcha-section"
import { useTranslations } from "next-intl"
import { useTheme } from "next-themes"
import Form from "next/form"
import { useActionState, useCallback, useRef, useState } from "react"
import { sendMagicLinkAction } from "@/app/[locale]/login/actions"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"
import { siteName, turnstile } from "@/config"
import { Link } from "@/i18n/navigation"
import { cn } from "@/lib/utils"
import { LogoAvatar } from "../common/avatar"
import CaptchaSection from "./captcha-section"
import OAuthSection from "./oauth-section"
import TermsCheckbox from "./terms-checkbox"

interface LoginFormProps {
  className?: string
  callbackURL: string
}

export function Login({ callbackURL, className }: LoginFormProps) {
  const t = useTranslations()
  const captchaRef = useRef<CaptchaSectionRef>(null)
  const { theme } = useTheme()

  // 客户端状态
  const [captchaToken, setCaptchaToken] = useState("")
  const [termsAccepted, setTermsAccepted] = useState(false)
  const [oauthTermsError, setOauthTermsError] = useState(false)

  // 使用 useActionState 管理发送Magic Link的状态
  const [magicLinkState, magicLinkAction, magicLinkPending] = useActionState(
    sendMagicLinkAction,
    { success: false, sent: false, errors: {} },
  )

  const handleOAuthLogin = useCallback (() => {
    if (!termsAccepted) {
      setOauthTermsError(true)
      // 3秒后清除错误状态
      setTimeout(() => setOauthTermsError(false), 1000)
      return false
    }
    setOauthTermsError(false)
    return true
  }, [termsAccepted])

  return (
    <Card className={cn("flex flex-col rounded-xl p-6 bg-white/50 dark:bg-white/5 border border-neutral-200 dark:border-white/10", className)}>
      <CardContent className="p-0">
        <Link href="/" className="flex flex-col items-center mb-4">
          <LogoAvatar alt={siteName} className="size-24" />
        </Link>

        {/* Magic Link 发送成功提示 */}
        {magicLinkState.sent
          ? (
              <div className="flex flex-col gap-4 text-center">
                <div className="bg-green-50 dark:bg-green-950/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
                  <h3 className="text-green-800 dark:text-green-200 font-medium mb-2">
                    {t("login.magicLinkSent")}
                  </h3>
                  <p className="text-green-700 dark:text-green-300 text-sm">
                    {t("login.checkEmail", { email: magicLinkState.email || "" })}
                  </p>
                </div>

                <Button
                  type="button"
                  variant="outline"
                  onClick={() => window.location.reload()}
                  className="w-full"
                >
                  {t("login.sendAnother")}
                </Button>
              </div>
            )
          : (
        /* Magic Link 发送表单 */
              <Form action={magicLinkAction} className="flex flex-col gap-4">
                <input type="hidden" name="callbackURL" value={callbackURL} />
                <input type="hidden" name="captchaToken" value={captchaToken} />
                <input type="hidden" name="termsAccepted" value={termsAccepted ? "true" : ""} />

                <div className="flex flex-col gap-2.5">
                  <Label htmlFor="email" className="text-sm font-medium tracking-wide text-neutral-700 dark:text-gray-300">
                    {t("login.emailLabel")}
                  </Label>
                  <Input
                    type="email"
                    name="email"
                    id="email"
                    required
                    className={cn(
                      "w-full rounded-lg border px-3.5 py-2.5 text-[15px] text-neutral-900 dark:text-white focus:outline-none focus:ring-0 sm:leading-6 transition-all duration-200",
                      magicLinkState.errors?.email
                        ? "border-red-500 dark:border-red-400 bg-red-50 dark:bg-red-950/20"
                        : "border-neutral-200/60 dark:border-white/[0.05] bg-white dark:bg-white/[0.02]",
                    )}
                    placeholder={t("login.emailPlaceholder")}
                  />
                  {magicLinkState.errors?.email && (
                    <span className="text-sm text-red-600 dark:text-red-400">
                      {magicLinkState.errors.email}
                    </span>
                  )}
                </div>

                {turnstile.enabled && (
                  <CaptchaSection
                    ref={captchaRef}
                    sitekey={turnstile.siteKey}
                    onVerify={token => setCaptchaToken(token)}
                    onExpire={() => setCaptchaToken("")}
                    className="my-4"
                    theme={theme === "dark" ? "dark" : "light"}
                    hasError={!!magicLinkState.errors?.captcha}
                  />
                )}

                <div className="flex flex-col gap-5">
                  <TermsCheckbox
                    checked={termsAccepted}
                    onChange={(checked) => {
                      setTermsAccepted(checked)
                      if (checked) {
                        setOauthTermsError(false)
                      }
                    }}
                    error={!!(magicLinkState.errors?.terms || oauthTermsError)}
                  />

                  <Button
                    type="submit"
                    disabled={magicLinkPending || (turnstile.enabled && !captchaToken)}
                    className="w-full rounded-lg tracking-wide hover:scale-[1.02] active:scale-[0.98] transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {magicLinkPending ? t("common.loading") : t("login.sendMagicLink")}
                  </Button>
                </div>
              </Form>
            )}

        <div className="relative my-8">
          <Separator className="border-neutral-200/60 dark:border-white/[0.05]" />
          <div className="absolute inset-0 flex items-center justify-center">
            <span className="px-4 text-neutral-500 dark:text-gray-400 text-[13px] tracking-wide">
              {t("login.or")}
            </span>
          </div>
        </div>

        <OAuthSection
          isPending={magicLinkPending}
          callbackURL={callbackURL}
          onOAuthLogin={handleOAuthLogin}
        />
      </CardContent>
    </Card>
  )
}
