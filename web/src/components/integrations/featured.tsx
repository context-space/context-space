"use client"

import type { Integration } from "@/typings"
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "../ui/carousel"
import { FeaturedCard } from "./featured-card"
import { Section } from "./section"

interface Props {
  integrations: Integration[]
  title: string
}

export function Featured({ integrations, title }: Props) {
  return (
    <Carousel opts={{ align: "center" }} className="pb-4">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold text-neutral-900 dark:text-white">
          {title}
        </h2>
        <div className="flex items-center gap-2">
          <CarouselPrevious className="relative left-0 top-0 translate-y-0" />
          <CarouselNext className="relative right-0 top-0 translate-y-0" />
        </div>
      </div>
      <CarouselContent className="p-2" style={{ width: "100px" }}>
        {integrations.map(integration => (
          <CarouselItem key={integration.identifier} className="basis-auto pr-6">
            <FeaturedCard provider={integration} className="w-[280px]" />
          </CarouselItem>
        ))}
      </CarouselContent>
    </Carousel>
  )
}

export async function _Featured({ integrations, title }: Props) {
  return (
    <Section title={title}>
      <ul className="grid grid-cols-1 xs:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-6">
        {integrations.map(integration => (
          <li key={integration.identifier} className="p-2">
            <FeaturedCard provider={integration} />
          </li>
        ))}
      </ul>
    </Section>
  )
}
