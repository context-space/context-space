import { Link } from "@/i18n/navigation"
import providersWithIcons from "./logo-wall-data.json"

export async function LogoWall() {
  return (
    <div className="w-full">
      <div className="flex items-center justify-center flex-wrap gap-4 sm:gap-6 md:gap-8">
        {providersWithIcons.map(provider => (
          <Link
            key={provider.identifier}
            href={`/integration/${provider.identifier}`}
            className="flex-shrink-0 opacity-60 hover:opacity-100 hover:scale-110 transition-all duration-300"
            title={provider.name}
          >
            <img
              src={provider.icon_url}
              alt={provider.name}
              className="h-7 w-7 sm:h-9 sm:w-9 object-contain grayscale hover:grayscale-0 transition-all duration-300"
            />
          </Link>
        ))}
      </div>
    </div>
  )
}
