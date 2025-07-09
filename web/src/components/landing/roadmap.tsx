"use client"

import { Check, Clock, Lightbulb } from "lucide-react"
import { useTranslations } from "next-intl"
import { cn } from "@/lib/utils"

const roadmapItems = [
  {
    quarter: "Phase 1",
    status: "completed",
    items: ["productionInfrastructure", "serviceIntegrations", "oauthFlows", "enterpriseSecurity"],
  },
  {
    quarter: "Phase 2",
    status: "current",
    items: ["nativeMcpSupport", "contextMemory", "smartAggregation", "apiDocumentation"],
  },
  {
    quarter: "Phase 3",
    status: "planned",
    items: ["semanticRetrieval", "contextOptimization", "realtimeUpdates", "advancedMcp"],
  },
  {
    quarter: "Phase 4",
    status: "planned",
    items: ["contextSynthesis", "predictiveLoading", "aiContextReasoning", "fullContextEngineering"],
  },
]

const StatusIcon = ({ status }: { status: string }) => {
  switch (status) {
    case "completed":
      return <Check className="text-green-500" size={16} />
    case "current":
      return <Clock className="text-primary" size={16} />
    case "planned":
      return <Lightbulb className="text-neutral-400" size={16} />
    default:
      return null
  }
}

export function Roadmap() {
  const t = useTranslations()

  return (
    <section id="roadmap" className="py-24">
      <div className="container mx-auto px-4">
        <div className="max-w-4xl mx-auto">
          {/* Header */}
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-neutral-900 dark:text-white mb-4">
              {t("landing.roadmap.title")}
            </h2>
            <p className="text-lg text-neutral-600 dark:text-gray-400 max-w-2xl mx-auto">
              {t("landing.roadmap.subtitle")}
            </p>
          </div>

          {/* Timeline */}
          <div className="relative">
            {/* Timeline Line */}
            <div className="absolute left-8 top-0 bottom-0 w-px bg-gradient-to-b from-primary/20 via-primary/40 to-primary/20 dark:from-primary/30 dark:via-primary/60 dark:to-primary/30"></div>

            {/* Timeline Items */}
            <div className="space-y-12">
              {roadmapItems.map(item => (
                <div key={item.quarter} className="relative">
                  {/* Timeline Dot */}
                  <div className={cn(
                    "absolute left-6 w-4 h-4 rounded-full border-2 flex items-center justify-center",
                    item.status === "completed" && "bg-green-500 border-green-500",
                    item.status === "current" && "bg-primary border-primary",
                    item.status === "planned" && "bg-neutral-200 dark:bg-neutral-700 border-neutral-300 dark:border-neutral-600",
                  )}
                  >
                    {item.status === "completed" && <Check size={8} className="text-white" />}
                    {item.status === "current" && <div className="w-1 h-1 bg-white rounded-full" />}
                  </div>

                  {/* Content */}
                  <div className="ml-20">
                    <div className={cn(
                      "group relative p-6 rounded-xl border transition-all duration-300",
                      item.status === "completed" && "bg-green-50/50 dark:bg-green-900/10 border-green-200/60 dark:border-green-800/30",
                      item.status === "current" && "bg-primary/5 dark:bg-primary/10 border-primary/20 dark:border-primary/30",
                      item.status === "planned" && "bg-neutral-50/50 dark:bg-neutral-800/20 border-neutral-200/60 dark:border-neutral-700/30",
                    )}
                    >
                      {/* Quarter Header */}
                      <div className="flex items-center gap-3 mb-4">
                        <StatusIcon status={item.status} />
                        <h3 className="text-lg font-semibold text-neutral-900 dark:text-white">
                          {item.quarter}
                        </h3>
                        <span className={cn(
                          "px-2 py-1 text-xs font-medium rounded-full",
                          item.status === "completed" && "bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300",
                          item.status === "current" && "bg-primary/20 text-primary dark:text-primary-300",
                          item.status === "planned" && "bg-neutral-100 dark:bg-neutral-800 text-neutral-600 dark:text-neutral-400",
                        )}
                        >
                          {t(`landing.roadmap.status.${item.status}`)}
                        </span>
                      </div>

                      {/* Features List */}
                      <div className="space-y-2">
                        {item.items.map(feature => (
                          <div key={feature} className="flex items-start gap-2">
                            <div className="mt-1.5 w-1 h-1 bg-neutral-400 dark:bg-neutral-500 rounded-full flex-shrink-0" />
                            <span className="text-sm text-neutral-600 dark:text-gray-400">
                              {t(`landing.roadmap.items.${feature}`)}
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
