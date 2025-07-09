import type { ReactNode } from "react"
import Footer from "@/components/footer"
import Header from "@/components/header"
import { cn } from "@/lib/utils"

interface LayoutProps {
  children: ReactNode
  className?: string
}

export const style = "max-w-[1600px] mx-auto px-6"
export function BaseLayout({ children, className }: LayoutProps) {
  return (
    <>
      <Header className={style} />
      <main className={cn("flex flex-col", style, className)}>
        {children}
      </main>
      <Footer className={cn("mt-4", style)} />
    </>
  )
}
