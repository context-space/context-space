"use client"

import { useTranslations } from "next-intl"
import { useRef } from "react"
import { cn } from "@/lib/utils"
import { LogoAvatar } from "../common"
import { LoadingMoreSearchTrigger } from "./search"

interface LoadingMoreProps {
  totalCount: number
  className?: string
}

export function LoadingMore({ totalCount, className }: LoadingMoreProps) {
  const t = useTranslations()
  const sectionRef = useRef<HTMLDivElement>(null)

  return (
    <div
      ref={sectionRef}
      className={cn(
        "py-16 text-center transition-all duration-700 ease-out",
        className,
      )}
    >
      <div className="max-w-2xl mx-auto space-y-6">
        <div className="flex justify-center">
          <div className="relative">
            <div className="absolute inset-0 bg-primary/20 rounded-full blur-lg animate-pulse"></div>
            <div className="relative bg-primary/10 p-4 rounded-full">
              <LogoAvatar alt="Context Space" className="size-10 text-primary" />
            </div>
          </div>
        </div>

        <div className="space-y-3">
          <h3 className="text-2xl md:text-3xl font-bold text-neutral-900 dark:text-white">
            {t("integrations.loadingMore.title", { count: totalCount })}
          </h3>
          <p className="text-lg text-neutral-600 dark:text-gray-400">
            {t("integrations.loadingMore.description")}
          </p>
        </div>

        <div className="flex flex-col sm:flex-row items-center justify-center gap-4 sm:gap-6">
          <LoadingMoreSearchTrigger />
        </div>

      </div>
    </div>
  )
}