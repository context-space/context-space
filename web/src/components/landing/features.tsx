"use client"

import { Brain, Code, GitBranch, Plug, ShieldCheck, Users } from "lucide-react"
import { FeatureCard } from "./feature-card"

const features = [
  {
    key: "1",
    icon: Brain,
    title: "Context Engineering Infrastructure",
    description: "Comprehensive context management beyond prompt engineering – the right information, at the right time, in the right format.",
  },
  {
    key: "2",
    icon: Plug,
    title: "14+ Service Integrations",
    description: "Production-ready integrations with GitHub, Slack, Airtable, HubSpot, Notion, Figma, Spotify, Stripe, and more.",
  },
  {
    key: "3",
    icon: GitBranch,
    title: "Enhanced MCP Support",
    description: "Production-grade Model Context Protocol implementation with built-in OAuth flows and persistent credential management.",
  },
  {
    key: "4",
    icon: ShieldCheck,
    title: "Enterprise Security",
    description: "HashiCorp Vault integration, automatic token rotation, encrypted credential storage, and OAuth UI flows.",
  },
  {
    key: "5",
    icon: Code,
    title: "RESTful API",
    description: "Clean HTTP endpoints with comprehensive documentation at api.context.space that actually work reliably in production.",
  },
  {
    key: "6",
    icon: Users,
    title: "Open Source Community",
    description: "AGPL v3 license transitioning to Apache 2.0, with active Discord community and welcoming contribution process.",
  },
]

export function Features() {
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
                Tame Chaos with Context
              </h2>
              <p className="text-xl text-neutral-600 dark:text-gray-400 max-w-3xl mx-auto leading-relaxed">
                Build more useful AI systems with seamless integrations, secure credential flows, and production-ready APIs
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
              Want to see more? Check out our roadmap or
              {" "}
              <a
                href="https://github.com/context-space/context-space"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:text-primary/80 transition-colors font-medium"
              >
                contribute on GitHub
              </a>
            </p>
          </div>
        </div>
      </div>
    </section>
  )
}
