import type { ReactNode } from "react"
import { cn } from "@/lib/utils"
import { style } from "./base"

interface LayoutProps {
  children: ReactNode
  className?: string
}

export function FramelessLayout({ children, className }: LayoutProps) {
  return (
    <main className={cn(style, "flex flex-col items-center justify-center h-screen px-1", className)}>
      {children}
    </main>
  )
}
