"use client"

import type { ScrollRestorationManagerProps } from "@/hooks/use-scroll-restoration"
import { useScrollRestoration } from "@/hooks/use-scroll-restoration"

export function ScrollRestorationManager(props: ScrollRestorationManagerProps) {
  useScrollRestoration(props)
  return null
}
