import { useTranslations } from "next-intl"
import { useCallback } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { InputOTP, InputOTPGroup, InputOTPSlot } from "@/components/ui/input-otp"
import { Label } from "@/components/ui/label"
import { cn } from "@/lib/utils"

interface EmailVerificationSectionProps {
  email: string
  verificationCode: string
  codeSent: boolean
  countdown: number
  emailError: boolean
  verificationCodeError: boolean
  isSendingCode: boolean
  isCaptchaRequired?: boolean
  onEmailChange: (value: string) => void
  onVerificationCodeChange: (value: string) => void
  onSendCode: () => void
}

export default function EmailVerificationSection({
  email,
  verificationCode,
  codeSent,
  countdown,
  emailError,
  verificationCodeError,
  isSendingCode,
  isCaptchaRequired = false,
  onEmailChange,
  onVerificationCodeChange,
  onSendCode,
}: EmailVerificationSectionProps) {
  const t = useTranslations()

  const getSlotClassName = useCallback(() => cn(
    "w-12 h-12 text-[15px] rounded-lg border text-neutral-900 dark:text-white focus:outline-none focus:ring-0 transition-all duration-200",
    verificationCodeError
      ? "border-red-500 dark:border-red-400 bg-red-50 dark:bg-red-950/20 focus:border-red-500 dark:focus:border-red-400"
      : "border-neutral-200/60 dark:border-white/[0.05] bg-white dark:bg-white/[0.02] focus:border-neutral-300 dark:focus:border-white/10",
  ), [verificationCodeError])

  return (
    <div className="flex flex-col gap-2.5">
      <div className="flex justify-between items-center gap-2.5">
        <Label
          htmlFor="email"
          className="text-sm font-medium tracking-wide text-neutral-700 dark:text-gray-300 block"
        >
          {t("login.emailLabel")}
        </Label>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={onSendCode}
          disabled={isSendingCode || countdown > 0 || isCaptchaRequired}
          className="text-sm font-medium px-2 py-1 h-auto rounded-md border border-neutral-200/80 dark:border-white/10 bg-white/80 dark:bg-white/[0.02] text-neutral-600 dark:text-gray-300 hover:bg-neutral-50 dark:hover:bg-white/[0.05] hover:border-neutral-300 dark:hover:border-white/20 hover:text-neutral-900 dark:hover:text-white transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-white/80 dark:disabled:hover:bg-white/[0.02]"
        >
          {countdown > 0
            ? (
                <span className="flex items-center gap-1">
                  <span className="w-3 h-3 border-2 border-current border-t-transparent rounded-full animate-spin opacity-70" />
                  {countdown}
                  s
                </span>
              )
            : isSendingCode
              ? (
                  <span className="flex items-center gap-1">
                    <span className="w-3 h-3 border-2 border-current border-t-transparent rounded-full animate-spin opacity-70" />
                    {t("login.sendingCode")}
                  </span>
                )
              : (
                  t("login.sendCode")
                )}
        </Button>
      </div>

      <Input
        type="email"
        name="email"
        id="email"
        required
        value={email}
        onChange={e => onEmailChange(e.target.value)}
        className={cn(
          "w-full rounded-lg border px-3.5 py-2.5 text-[15px] text-neutral-900 dark:text-white focus:outline-none focus:ring-0 sm:leading-6 transition-all duration-200",
          emailError
            ? "border-red-500 dark:border-red-400 bg-red-50 dark:bg-red-950/20 focus:border-red-500 dark:focus:border-red-400"
            : "border-neutral-200/60 dark:border-white/[0.05] bg-white dark:bg-white/[0.02] focus:border-neutral-300 dark:focus:border-white/10",
        )}
        placeholder={t("login.emailPlaceholder")}
      />

      {codeSent && (
        <>
          <div className="flex justify-between">
            <span className="text-sm text-neutral-500 dark:text-gray-400">
              {t("login.enterCode")}
            </span>
            {verificationCodeError && (
              <span className="text-sm text-red-600 dark:text-red-400">
                {t("login.invalidCode")}
              </span>
            )}
          </div>
          <InputOTP
            maxLength={6}
            value={verificationCode}
            onChange={onVerificationCodeChange}
          >
            <InputOTPGroup className="w-full justify-between">
              {Array.from({ length: 6 }).map((_, index) => (
                <InputOTPSlot
                  key={index.toFixed()}
                  index={index}
                  className={getSlotClassName()}
                />
              ))}
            </InputOTPGroup>
          </InputOTP>
        </>
      )}
    </div>
  )
}
