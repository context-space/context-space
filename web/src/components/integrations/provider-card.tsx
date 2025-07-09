import type { Integration } from "@/typings"
import { Link } from "@/i18n/navigation"
import { cn } from "@/lib/utils"
import { ProviderAvatar } from "../integration/provider-avatar"
import { Status } from "./status"

interface ProviderCardProps {
  provider: Integration
}

export function ProviderCard({ provider }: ProviderCardProps) {
  return (
    <Link
      href={`/integration/${provider.identifier}`}
      className="block h-full transition-transform hover:scale-[1.01] duration-300"
    >
      <div className={cn(
        "group relative flex flex-col h-full p-6 rounded-xl bg-white/[0.02]",
        "border border-primary/15 dark:border-white/10 hover:border-primary/30 dark:hover:border-primary/40",
        "transition-all duration-300",
      )}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="relative shrink-0">
              <ProviderAvatar
                src={provider.icon_url}
                alt={`${provider.name} logo`}
                className="size-10"
              />
            </div>
            <h3 className="font-medium text-neutral-900 dark:text-white truncate">
              {provider.name}
            </h3>
          </div>
          {provider.connection_status && <Status status={provider.connection_status} type="badge" />}
        </div>

        {provider.description && (
          <p className="mt-4 text-sm text-neutral-600 dark:text-gray-400 line-clamp-2">
            {provider.description}
          </p>
        )}
      </div>
    </Link>
  )
}
