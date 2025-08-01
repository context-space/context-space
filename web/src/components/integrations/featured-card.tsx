import type { Integration } from "@/typings"
import { Link } from "@/i18n/navigation"
import { cn } from "@/lib/utils"
import { ProviderAvatar } from "../common/avatar"
import { Button } from "../ui/button"
import { Status } from "./status"

interface FeaturedCardProps {
  provider: Integration
  className?: string
}

export function FeaturedCard({ provider, className }: FeaturedCardProps) {
  return (
    <Link
      href={`/integration/${provider.identifier}`}
      className="block h-full transition-transform hover:scale-[1.02] duration-300"
    >
      <div className={cn(
        "group relative flex flex-col h-full p-8 rounded-xl bg-white/[0.02] items-center gap-5",
        "border border-primary/15 dark:border-white/10 hover:border-primary/30 dark:hover:border-primary/40",
        "transition-all duration-300 w-full min-h-[300px]",
        className,
      )}
      >
        <ProviderAvatar
          src={provider.icon_url}
          alt={`${provider.name} logo`}
          className="size-20"
        />

        <h3 className="font-semibold text-neutral-900 dark:text-white text-lg truncate">
          {provider.name}
        </h3>

        {provider.description && (
          <p className="text-base text-neutral-600 dark:text-gray-400 text-center line-clamp-3 leading-relaxed">
            {provider.description}
          </p>
        )}

        {
          provider.connection_status
            ? (
                <div className="flex justify-center flex-1 items-end">
                  <Status status={provider.connection_status} type="badge" />
                </div>
              )
            : (
                <div className="flex justify-center flex-1 items-end">
                  <Button variant="outline" size="sm">Try it</Button>
                </div>
              )
        }
      </div>
    </Link>
  )
}
