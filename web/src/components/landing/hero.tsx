"use client"

import { ArrowRight, Github, Sparkles } from "lucide-react"
import { useTranslations } from "next-intl"
import { Button } from "@/components/ui/button"
import { Link } from "@/i18n/navigation"
import { TypewriterText } from "./typewriter"

export function Hero() {
  const t = useTranslations()

  return (
    <div id="hero" className="relative min-h-[90vh] flex items-center overflow-hidden">
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/20 dark:bg-primary/30 rounded-full blur-3xl opacity-20 animate-pulse"></div>
      <div className="absolute bottom-1/4 right-1/4 w-64 h-64 bg-primary/30 dark:bg-primary/40 rounded-full blur-2xl opacity-30 animate-pulse delay-1000"></div>

      <div className="container mx-auto px-4 py-20 relative z-10">
        <div className="max-w-7xl mx-auto text-center space-y-16">
          <a href="https://github.com/context-space/context-space" className="inline-flex items-center gap-2.5 px-5 py-2.5 rounded-full bg-primary/10 dark:bg-primary/20 border border-primary/20 dark:border-primary/30 hover:bg-primary/15 dark:hover:bg-primary/25 transition-colors duration-300">
            <Sparkles className="text-primary" size={16} />
            <span className="text-sm font-semibold text-primary tracking-wide">{t("hero.badge")}</span>
            {/* <img
              alt="GitHub stars badge"
              src="https://img.shields.io/github/stars/context-space/context-space?style=flat&labelColor=%238b8cac&color=%239899bd"
              className="opacity-90"
            /> */}
          </a>

          <div className="space-y-8">
            <h1 className="text-5xl md:text-6xl lg:text-7xl font-bold tracking-tight leading-[1.1] text-neutral-900 dark:text-white">
              <span className="block">
                <span className="bg-gradient-to-r from-neutral-900 via-neutral-800 to-neutral-900 dark:from-white dark:via-gray-100 dark:to-white bg-clip-text text-transparent">
                  {t("hero.title.part1")}
                  {" "}
                </span>
                <span className="bg-gradient-to-r from-primary via-primary/80 to-primary bg-clip-text text-transparent font-extrabold">{t("hero.title.part2")}</span>
              </span>
              <span className="block bg-gradient-to-r from-neutral-900 via-neutral-800 to-neutral-900 dark:from-white dark:via-gray-100 dark:to-white bg-clip-text text-transparent">
                {t("hero.title.part3")}
              </span>
            </h1>

            <div className="space-y-4 max-w-4xl mx-auto">
              <TypewriterText />
            </div>
          </div>

          <div className="flex flex-col sm:flex-row gap-6 justify-center items-center pt-12">
            <Button asChild size="lg" className="h-14 px-8 text-lg font-semibold border-2 hover:border-primary/40 dark:hover:border-primary/60 transition-all duration-300">
              <Link
                href="/github"
                target="_blank"
                rel="noopener noreferrer"
                className="group inline-flex items-center justify-center gap-3"
              >
                <Github size={22} />
                {t("hero.viewOnGithub")}
              </Link>
            </Button>
            <Button asChild variant="outline" className="h-14 px-8 text-lg font-semibold shadow-lg hover:shadow-xl transition-all duration-300" size="lg">
              <Link href="/blogs/what-is-context-engineering" className="group inline-flex items-center justify-center gap-3">
                {t("hero.readMore")}
                <ArrowRight size={22} className="group-hover:translate-x-1 transition-transform duration-300" />
              </Link>
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
