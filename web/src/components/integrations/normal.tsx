import type { Integration } from "@/typings"
import { ProviderCard } from "./provider-card"
import { Section } from "./section"

interface Props {
  integrations: Integration[]
  title: string
}

export async function Normal({ integrations, title }: Props) {
  return (
    <Section title={title}>
      <ul className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
        {integrations.map(provider => (
          <li key={provider.identifier} className="p-2">
            <ProviderCard provider={provider} />
          </li>
        ))}
      </ul>
    </Section>
  )
}
