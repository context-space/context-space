import { ArrowRight, Github, Rocket, Sparkles } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Link } from "@/i18n/navigation"
import { LogoWall } from "./logo-wall"
import { TypewriterText } from "./typewriter"

export function Hero() {
  return (
    <div id="hero" className="relative min-h-[85vh] sm:min-h-[90vh] overflow-hidden flex flex-col items-center justify-center">
      <div className="absolute top-1/4 left-1/4 w-48 h-48 sm:w-64 sm:h-64 md:w-96 md:h-96 bg-primary/20 dark:bg-primary/30 rounded-full blur-3xl opacity-20 animate-pulse"></div>
      <div className="absolute bottom-1/4 right-1/4 w-32 h-32 sm:w-48 sm:h-48 md:w-64 md:h-64 bg-primary/30 dark:bg-primary/40 rounded-full blur-2xl opacity-30 animate-pulse delay-1000"></div>

      <div className="container mx-auto px-4 py-12 sm:py-16 md:py-20 relative z-10">
        <div className="max-w-7xl mx-auto text-center space-y-8 sm:space-y-12 md:space-y-16">
          <a href="https://github.com/context-space/context-space" className="inline-flex items-center gap-2 sm:gap-2.5 px-4 sm:px-5 py-2 sm:py-2.5 rounded-full bg-primary/10 dark:bg-primary/20 border border-primary/20 dark:border-primary/30 hover:bg-primary/15 dark:hover:bg-primary/25 transition-colors duration-300">
            <Sparkles className="text-primary" size={14} />
            <span className="text-xs sm:text-sm font-semibold text-primary tracking-wide">Open Source</span>
            <img
              alt="GitHub stars badge"
              src="https://img.shields.io/github/stars/context-space/context-space?style=flat&labelColor=%238b8cac&color=%239899bd"
              className="opacity-90"
            />
          </a>

          <div className="space-y-6 sm:space-y-8">
            <h1
              className="font-bold tracking-tight leading-[1.4] md:leading-[1.2] text-neutral-900 dark:text-white"
              style={{ fontSize: "clamp(1rem, 8vw, 4rem)" }}
            >
              <span className="block">
                <span className="bg-gradient-to-r from-neutral-900 via-neutral-800 to-neutral-900 dark:from-white dark:via-gray-100 dark:to-white bg-clip-text text-transparent">
                  {"Tool-first "}
                </span>
                <span className="bg-gradient-to-r from-primary via-primary/80 to-primary bg-clip-text text-transparent font-extrabold truncate">Context Engineering</span>
              </span>
              <span className="block bg-gradient-to-r from-neutral-900 via-neutral-800 to-neutral-900 dark:from-white dark:via-gray-100 dark:to-white bg-clip-text text-transparent">
                Infrastructure
              </span>
            </h1>

            <div className="space-y-4 max-w-sm sm:max-w-2xl md:max-w-4xl mx-auto">
              <TypewriterText />
            </div>
          </div>

          <div className="flex flex-col sm:flex-row gap-4 sm:gap-6 justify-center items-center">
            <Button asChild variant="outline" className="h-12 sm:h-14 px-6 sm:px-8 text-base sm:text-lg font-semibold hover:scale-105 transform transition-all duration-300 w-full sm:w-auto" size="lg">
              <Link href="/integrations" className="group inline-flex items-center justify-center gap-2 sm:gap-3">
                Try Integrations
                <Rocket size={18} className="sm:w-[22px] sm:h-[22px] group-hover:translate-x-1 transition-transform duration-300" />
              </Link>
            </Button>
            <Button asChild size="lg" className="h-12 sm:h-14 px-6 sm:px-8 text-base sm:text-lg font-semibold hover:scale-105 transform transition-all duration-300 w-full sm:w-auto">
              <Link
                href="/github"
                target="_blank"
                rel="noopener noreferrer"
                className="group inline-flex items-center justify-center gap-2 sm:gap-3"
              >
                View on GitHub
                <Github size={18} className="sm:w-[22px] sm:h-[22px]" />
              </Link>
            </Button>
            <Button asChild variant="outline" className="h-12 sm:h-14 px-6 sm:px-8 text-base sm:text-lg font-semibold hover:scale-105 transform transition-all duration-300 w-full sm:w-auto" size="lg">
              <Link href="/blogs/context-is-the-new-engine" className="group inline-flex items-center justify-center gap-2 sm:gap-3">
                Read More
                <ArrowRight size={18} className="sm:w-[22px] sm:h-[22px] group-hover:translate-x-1 transition-transform duration-300" />
              </Link>
            </Button>
          </div>

          <LogoWall />
        </div>
      </div>

    </div>
  )
}
