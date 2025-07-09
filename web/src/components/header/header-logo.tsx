import Logo from "@/components/common/logo"
import { Link } from "@/i18n/navigation"

export function HeaderLogo() {
  return (
    <Link href="/" className="flex items-center gap-3 hover:opacity-80 transition-opacity">
      <Logo size={28} />
      <span className="text-xl font-semibold font-mono">
        Context Space
      </span>
    </Link>
  )
}
