import { ProviderAvatar } from "../common/avatar"

interface HeaderProps {
  name: string
  description: string
  iconUrl?: string
}

export function Header({ name, description, iconUrl }: HeaderProps) {
  return (
    <div className="flex items-start gap-6 mt-5">
      <div className="relative shrink-0 group">
        <div className="absolute inset-0 bg-gradient-to-b from-primary/3 to-transparent blur-xl -z-10 opacity-0 group-hover:opacity-60 transition-opacity duration-300"></div>
        <ProviderAvatar
          src={iconUrl}
          alt={`${name} logo`}
          className="size-20"
        />
      </div>
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-bold tracking-tight text-neutral-900 dark:text-white">{name}</h1>
        <p className="text-neutral-600 dark:text-gray-400 text-sm max-w-2xl leading-relaxed">
          {description}
        </p>
      </div>
    </div>
  )
}
