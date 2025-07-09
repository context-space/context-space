"use client"

import { useTranslations } from "next-intl"
import { useMemo } from "react"
import { Link } from "@/i18n/navigation"
import { cn } from "@/lib/utils"
import { NewsletterForm } from "./newsletter-form"

interface FooterProps {
  className?: string
}

interface SocialSection {
  title: string
  links: Array<{
    name: string
    href: string
  }>
}

interface SocialProps {
  section: SocialSection
  index: number
}

function Social({ section, index }: SocialProps) {
  return (
    <div className={index > 1 ? "mt-12 md:mt-0" : ""}>
      <span className="font-medium tracking-wide">
        {section.title}
      </span>
      <ul role="list" className="mt-5 space-y-3">
        {section.links.map(link => (
          <li key={link.name}>
            <Link
              href={link.href}
              className="text-sm text-neutral-600 dark:text-gray-400 hover:text-neutral-900 dark:hover:text-white transition-colors"
            >
              {link.name}
            </Link>
          </li>
        ))}
      </ul>
    </div>
  )
}

const currentYear = new Date().getFullYear()

export default function Footer({ className }: FooterProps) {
  const t = useTranslations()

  const SocialSections = useMemo((): SocialSection[] => [
    // {
    //   title: t("footer.features"),
    //   links: [
    //     { name: t("nav.integrations"), href: "/integrations" },
    //     { name: t("nav.workflows"), href: "/workflows" },
    //   ],
    // },
    // {
    //   title: t("footer.info.title"),
    //   links: [
    //     { name: t("footer.info.aboutUs"), href: "#" },
    //     { name: t("footer.info.contact"), href: "mailto:hi@context.space" },
    //   ],
    // },
    {
      title: t("footer.community"),
      links: [
        { name: "GitHub", href: "https://github.com/context-space/context-space" },
        { name: "Twitter", href: "https://x.com/hi_contextspace" },
        { name: "Discord", href: "https://discord.gg/Q74Ta5Xv" },
      ],
    },
    {
      title: t("footer.resources"),
      links: [
        { name: t("footer.resourcesItems.docs"), href: "/docs" },
        { name: t("footer.resourcesItems.blogs"), href: "/blogs" },
        { name: t("footer.resourcesItems.roadmap"), href: "/#roadmap" },
      ],
    },
    {
      title: t("footer.legal"),
      links: [
        { name: t("footer.terms"), href: "/terms" },
        { name: t("footer.privacy"), href: "/privacy" },
      ],
    },
  ], [t])

  return (
    <footer className={cn("border-t border-neutral-200/60 dark:border-white/[0.05] py-8", className)}>
      <div className="xl:grid xl:grid-cols-3 xl:gap-12">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 xl:col-span-2">
          {SocialSections.map((section, index) => (
            <Social
              key={section.title}
              section={section}
              index={index}
            />
          ))}
        </div>

        <NewsletterForm />
      </div>

      <div className="border-neutral-200/60 dark:border-white/[0.05] mt-8">
        <p className="text-sm text-neutral-600 dark:text-gray-400">
          {t("footer.copyright", { year: currentYear })}
          <span>
            {" "}
            {t("footer.contact.text")}
            <a href="mailto:hi@context.space">
              {" "}
              {t("footer.contact.email")}
            </a>
          </span>
        </p>
      </div>
    </footer>
  )
}
