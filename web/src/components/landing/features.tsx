"use client"

import { Brain, Code, GitBranch, Plug, ShieldCheck, Users } from "lucide-react"
import { useTranslations } from "next-intl"
import { useMemo } from "react"
import { FeatureCard } from "./feature-card"

export function Features() {
  const t = useTranslations()

  const features = useMemo(() => [
    {
      key: "1",
      icon: Brain,
      title: t("hero.features.1.title"),
      description: t("hero.features.1.description"),
    },
    {
      key: "2",
      icon: Plug,
      title: t("hero.features.2.title"),
      description: t("hero.features.2.description"),
    },
    {
      key: "3",
      icon: GitBranch,
      title: t("hero.features.3.title"),
      description: t("hero.features.3.description"),
    },
    {
      key: "4",
      icon: ShieldCheck,
      title: t("hero.features.4.title"),
      description: t("hero.features.4.description"),
    },
    {
      key: "5",
      icon: Code,
      title: t("hero.features.5.title"),
      description: t("hero.features.5.description"),
    },
    {
      key: "6",
      icon: Users,
      title: t("hero.features.6.title"),
      description: t("hero.features.6.description"),
    },
  ], [t])

  return (
    <section id="features" className="relative py-24 overflow-hidden">
      {/* Background gradients */}
      <div className="absolute top-20 left-1/3 w-96 h-96 bg-primary/10 dark:bg-primary/20 rounded-full blur-3xl opacity-30"></div>
      <div className="absolute bottom-20 right-1/3 w-64 h-64 bg-primary/15 dark:bg-primary/25 rounded-full blur-2xl opacity-40"></div>

      <div className="container mx-auto px-4 relative z-10">
        <div className="max-w-6xl mx-auto">
          {/* Section header */}
          <div className="text-center space-y-6 mb-16">
            <div className="space-y-4">
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-neutral-900 dark:text-white mb-4">
                {t("landing.features.title")}
              </h2>
              <p className="text-xl text-neutral-600 dark:text-gray-400 max-w-3xl mx-auto leading-relaxed">
                {t("landing.features.subtitle")}
              </p>
            </div>
          </div>

          {/* Features grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {features.map(feature => (
              <div key={feature.key} className="group">
                <FeatureCard
                  Icon={feature.icon}
                  title={feature.title}
                  description={feature.description}
                />
              </div>
            ))}
          </div>

          {/* Bottom CTA */}
          <div className="text-center mt-16">
            <p className="text-neutral-500 dark:text-gray-400 text-sm">
              {t("landing.features.cta.text")}
              {" "}
              <a
                href="https://github.com/context-space/context-space"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:text-primary/80 transition-colors font-medium"
              >
                {t("landing.features.cta.link")}
              </a>
            </p>
          </div>
        </div>
      </div>
    </section>
  )
}
