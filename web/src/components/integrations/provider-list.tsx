import type { Integration } from "@/typings"
import { useTranslations } from "next-intl"
import { ProviderCard } from "./provider-card"

interface Props {
  integrations: Integration[]
}

export function Latest({ integrations: providers }: Props) {
  const t = useTranslations()

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-xl font-semibold text-neutral-900 dark:text-white mb-6">
          {t("integrations.latestIntegrations")}
        </h2>
        <ul className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
          {providers.map(provider => (
            <li key={provider.identifier} className="p-2">
              <ProviderCard provider={provider} />
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}
