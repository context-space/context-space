"use client"

import { useTheme } from "next-themes"
import { Toaster as Sonner, ToasterProps } from "sonner"

const Toaster = ({ ...props }: ToasterProps) => {
  const { theme = "system" } = useTheme()

  return (
    <Sonner
      theme={theme as ToasterProps["theme"]}
      className="toaster group"
      toastOptions={{
        classNames: {
          toast: "group toast !bg-white/60 dark:!bg-white/[0.02] !border-neutral-200/60 dark:!border-white/[0.08] !backdrop-blur-sm !shadow-none !rounded-lg !p-4 !text-sm !border",
          description: "!text-muted-foreground",
          actionButton: "!bg-primary !text-primary-foreground hover:!bg-primary/90 !px-3 !py-1.5 !rounded-sm !text-xs !font-medium",
          cancelButton: "!bg-muted !text-muted-foreground hover:!bg-muted/80 !px-3 !py-1.5 !rounded-sm !text-xs !font-medium",
          success: "!bg-green-50/60 dark:!bg-green-500/10 !border-green-200/60 dark:!border-green-500/20 !text-green-700 dark:!text-green-400",
          error: "!bg-red-50/60 dark:!bg-red-900/20 !border-red-200/60 dark:!border-red-900/30 !text-red-700 dark:!text-red-300",
          warning: "!bg-yellow-50/60 dark:!bg-yellow-900/20 !border-yellow-200/60 dark:!border-yellow-900/30 !text-yellow-700 dark:!text-yellow-300",
          info: "!bg-blue-50/60 dark:!bg-blue-900/20 !border-blue-200/60 dark:!border-blue-900/30 !text-blue-700 dark:!text-blue-300",
        },
      }}
      {...props}
    />
  )
}

export { Toaster }
