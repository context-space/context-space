"use client"

import { useState } from "react"
import { cn } from "@/lib/utils"
import { DesktopActions } from "./desktop-actions"
import { DesktopNavigation } from "./desktop-navigation"
import { HeaderLogo } from "./header-logo"
import { MobileMenu } from "./mobile-menu"
import { MobileMenuButton } from "./mobile-menu-button"

interface HeaderProps {
  className?: string
}

export default function Header({ className }: HeaderProps) {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <header className="sticky top-0 z-50 w-full backdrop-blur-2xl">
      <div className={cn("mx-auto max-w-[1600px] px-6", className)}>
        <div className="flex h-16 items-center justify-between">
          <HeaderLogo />
          <div className="flex items-center gap-4">
            <DesktopNavigation />
            <div className="w-px h-4 bg-neutral-200 dark:bg-neutral-700" />
            <DesktopActions />
            <MobileMenuButton
              isOpen={isOpen}
              onToggle={() => setIsOpen(!isOpen)}
            />
          </div>
        </div>
        <MobileMenu
          isOpen={isOpen}
          onClose={() => setIsOpen(false)}
        />
      </div>
    </header>
  )
}
