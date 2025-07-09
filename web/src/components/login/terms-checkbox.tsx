import { useTranslations } from "next-intl"
import { Checkbox } from "@/components/ui/checkbox"
import { Label } from "@/components/ui/label"
import { Link } from "@/i18n/navigation"
import { cn } from "@/lib/utils"

interface TermsCheckboxProps {
  checked: boolean
  onChange: (checked: boolean) => void
  error: boolean
}

export default function TermsCheckbox({ checked, onChange, error }: TermsCheckboxProps) {
  const t = useTranslations()

  return (
    <div className="flex items-center gap-3">
      <Checkbox
        id="terms"
        checked={checked}
        onCheckedChange={onChange}
        required
        className={cn(
          "h-4 w-4 rounded border text-neutral-900 focus:ring-0 mt-0.5 flex-shrink-0 transition-all duration-200",
          error
            ? "animate-shake border-red-500 data-[state=checked]:bg-red-500 data-[state=checked]:border-red-500"
            : "border-neutral-300",
        )}
      />
      <Label
        htmlFor="terms"
        className={cn(
          "text-sm leading-relaxed cursor-pointer transition-all duration-200 flex-1",
          error
            ? "text-red-600 dark:text-red-400 animate-shake"
            : "text-neutral-700 dark:text-gray-300",
        )}
      >
        <span className="block">
          {t("login.terms")}
          {" "}
          <Link href="/terms" className="font-medium text-neutral-900 dark:text-white hover:underline">
            {t("login.termsOfService")}
          </Link>
          {" "}
          and
          {" "}
          <Link href="/privacy" className="font-medium text-neutral-900 dark:text-white hover:underline">
            {t("login.privacyPolicy")}
          </Link>
        </span>
      </Label>
    </div>
  )
}
