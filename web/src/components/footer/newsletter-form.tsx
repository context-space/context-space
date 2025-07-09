import { useTranslations } from "next-intl"
import { useCallback, useMemo, useState } from "react"
import { cn } from "@/lib/utils"

export function NewsletterForm() {
  const t = useTranslations()
  const [email, setEmail] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [status, setStatus] = useState<"idle" | "success" | "error">("idle")

  const handleSubmit = useCallback(async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setIsLoading(true)
    setStatus("idle")

    try {
      // TODO: Implement newsletter subscription API call
      await new Promise(resolve => setTimeout(resolve, 1000)) // Simulated API call
      setStatus("success")
      setEmail("")
    } catch {
      setStatus("error")
    } finally {
      setIsLoading(false)
    }
  }, [])

  const handleEmailChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value)
    if (status !== "idle") {
      setStatus("idle")
    }
  }, [status])

  const buttonText = useMemo(() => {
    if (isLoading) return t("common.loading")
    if (status === "success") return t("footer.newsletter.success")
    return t("footer.newsletter.button")
  }, [isLoading, status, t])

  return (
    <div className="flex flex-col mt-12 xl:mt-0 gap-4">
      <span className="font-medium tracking-wide">
        {t("footer.newsletter.title")}
      </span>
      <span className="text-sm text-neutral-600 dark:text-gray-400">
        {t("footer.newsletter.description")}
      </span>

      <form
        onSubmit={handleSubmit}
        className="flex flex-col space-y-3 sm:space-y-0 sm:flex-row sm:gap-3 sm:max-w-md"
        noValidate
      >
        <div className="flex-1">
          <label htmlFor="email-address" className="sr-only">
            {t("footer.newsletter.placeholder")}
          </label>
          <input
            type="email"
            name="email-address"
            id="email-address"
            autoComplete="email"
            required
            value={email}
            onChange={handleEmailChange}
            disabled={isLoading}
            aria-describedby={status === "error" ? "email-error" : undefined}
            className={cn(
              "relative w-full min-w-0 h-10 px-4 py-2",
              "bg-white/80 dark:bg-black/20 backdrop-blur-sm",
              "border border-neutral-200/60 dark:border-white/10",
              "rounded-md text-sm placeholder:text-neutral-500 dark:placeholder:text-neutral-400",
              "transition-all duration-200 ease-out",
              "focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary/30",
              "hover:border-neutral-300/80 dark:hover:border-white/20",
              "disabled:opacity-50 disabled:cursor-not-allowed",
              status === "error" && "border-red-500/50 focus:ring-red-500/20",
              status === "success" && "border-green-500/50 focus:ring-green-500/20",
            )}
            placeholder={t("footer.newsletter.placeholder")}
          />
        </div>

        <button
          type="submit"
          disabled={isLoading || !email.trim()}
          className={cn(
            "relative flex-shrink-0 h-10 px-4 min-w-[100px]",
            "bg-primary/10 hover:bg-primary/20",
            "border border-primary/20 rounded-md",
            "text-sm font-medium",
            "transition-all duration-200 ease-out",
            "focus:outline-none focus:ring-2 focus:ring-primary/30 focus:ring-offset-2 dark:focus:ring-offset-black/20",
            "disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-primary/10",
            status === "success" && "bg-green-500/10 border-green-500/20 text-green-700 dark:text-green-400",
          )}
        >
          <span className="relative flex items-center justify-center">
            {buttonText}
          </span>
        </button>
      </form>

      {/* Status messages */}
      {status === "error" && (
        <p id="email-error" className="text-sm text-red-600 dark:text-red-400" role="alert">
          {t("footer.newsletter.error")}
        </p>
      )}
    </div>
  )
}
