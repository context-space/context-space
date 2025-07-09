import type { ReactNode } from "react"
import Footer from "@/components/footer"
import HomeHeader from "@/components/header/home-header"
import { cn } from "@/lib/utils"
import { style } from "./base"

interface HomeLayoutProps {
  children: ReactNode
  className?: string
}

export function HomeLayout({ children, className }: HomeLayoutProps) {
  return (
    <>
      <HomeHeader className={style} />
      <main className={cn("flex flex-col", style, className)}>
        {children}
      </main>
      <Footer className={cn("mt-4", style)} />
    </>
  )
}
