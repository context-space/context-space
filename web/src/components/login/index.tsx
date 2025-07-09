"use client"

import type { CaptchaSectionRef } from "./captcha-section"
import { useTranslations } from "next-intl"
import { useTheme } from "next-themes"
import { useRef } from "react"
import Logo from "@/components/common/logo"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { turnstile } from "@/config"
import { useLoginForm } from "@/hooks/use-login-form"
import { Link } from "@/i18n/navigation"
import { cn } from "@/lib/utils"
import CaptchaSection from "./captcha-section"
import EmailVerificationSection from "./email-verification-section"
import OAuthSection from "./oauth-section"
import TermsCheckbox from "./terms-checkbox"

interface LoginFormProps {
  from: string
  className?: string
}

export function Login({ from, className }: LoginFormProps) {
  const t = useTranslations()
  const captchaRef = useRef<CaptchaSectionRef>(null)

  const {
    email,
    verificationCode,
    codeSent,
    countdown,
    termsAccepted,
    captchaToken,
    emailError,
    termsError,
    verificationCodeError,
    captchaError,
    isPending,
    isSendingCode,
    handleSendCode,
    handleEmailLogin,
    handleVerificationCodeChange,
    handleOAuthLogin,
    handleAnonymousLogin,
    handleEmailChange,
    handleTermsChange,
    handleCaptchaVerify,
    handleCaptchaExpire,
    resetCaptcha,
  } = useLoginForm(from, captchaRef)

  const { theme } = useTheme()

  // CAPTCHA is enabled based on environment configuration

  return (
    <Card className={cn("flex flex-col rounded-xl p-6 bg-white/50 dark:bg-white/5 border border-neutral-200 dark:border-white/10 w-[400px] max-w-[90vw]", className)}>
      <CardContent className="p-0">
        <Link href="/" className="flex flex-col items-center mb-4">
          <Logo size={72} />
        </Link>

        <form className="flex flex-col gap-4" onSubmit={handleEmailLogin}>
          <EmailVerificationSection
            email={email}
            verificationCode={verificationCode}
            codeSent={codeSent}
            countdown={countdown}
            emailError={emailError}
            verificationCodeError={verificationCodeError}
            isSendingCode={isSendingCode}
            onEmailChange={handleEmailChange}
            onVerificationCodeChange={handleVerificationCodeChange}
            onSendCode={handleSendCode}
          />

          {!codeSent && turnstile.enabled && (
            <CaptchaSection
              ref={captchaRef}
              sitekey={turnstile.siteKey}
              onVerify={handleCaptchaVerify}
              onExpire={handleCaptchaExpire}
              className="my-4"
              theme={theme === "dark" ? "dark" : "light"}
              hasError={captchaError}
            />
          )}

          <div className="flex flex-col gap-5">
            <TermsCheckbox
              checked={termsAccepted}
              onChange={handleTermsChange}
              error={termsError}
            />

            <Button
              type="submit"
              disabled={isPending || !termsAccepted || !codeSent || verificationCode.length !== 6}
              className="w-full rounded-lg tracking-wide hover:scale-[1.02] active:scale-[0.98] transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isPending ? t("common.loading") : t("login.loginWithEmail")}
            </Button>
          </div>
        </form>

        <div className="relative my-8">
          <Separator className="border-neutral-200/60 dark:border-white/[0.05]" />
          <div className="absolute inset-0 flex items-center justify-center">
            <span className="px-4 text-neutral-500 dark:text-gray-400 text-[13px] tracking-wide">
              {t("login.or")}
            </span>
          </div>
        </div>

        <OAuthSection
          isPending={isPending}
          onOAuthLogin={handleOAuthLogin}
          onAnonymousLogin={handleAnonymousLogin}
        />
      </CardContent>
    </Card>
  )
}
