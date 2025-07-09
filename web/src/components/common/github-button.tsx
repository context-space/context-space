import { Github } from "lucide-react"
import { useTranslations } from "next-intl"
import { Button } from "@/components/ui/button"

export function GithubButton() {
  const t = useTranslations()
  return (
    <Button
      variant="ghost"
      size="sm"
      className="h-9 w-9"
      asChild
    >
      <a
        href="https://github.com/context-space/context-space"
        target="_blank"
        rel="noopener noreferrer"
        aria-label={t("common.github")}
      >
        <Github className="h-4 w-4" />
      </a>
    </Button>
  )
}
