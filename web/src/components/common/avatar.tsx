"use client"

import { Avatar } from "@/components/ui/avatar"
import { cn } from "@/lib/utils"

interface AvatarProps {
  src?: string
  alt: string
  className?: string
}

export function LogoAvatar({ src, alt, className }: AvatarProps) {
  return (
    <Avatar className={cn("dark:text-foreground text-primary/90 bg-transparent rounded-none transition-transform duration-200 group-hover:scale-102", className)}>
      <div
        className={cn(
          "w-full h-full bg-center bg-no-repeat",
          "bg-current",
        )}
        style={{
          mask: `url(${src || "/logo.svg"}) center/contain no-repeat`,
          WebkitMask: `url(${src || "/logo.svg"}) center/contain no-repeat`,
          backgroundColor: "currentColor",
        }}
        role="img"
        aria-label={alt}
      />
    </Avatar>
  )
}

export function ProviderAvatar({ src, alt, className }: AvatarProps) {
  return (
    <Avatar className={cn("bg-transparent rounded-none transition-transform duration-200 group-hover:scale-102", className)}>
      <div
        className={cn(
          "w-full h-full bg-center bg-no-repeat",
          src ? "bg-contain" : "bg-current",
        )}
        style={
          src
            ? { backgroundImage: `url(${src})` }
            : {
                mask: `url(/mcp.svg) center/contain no-repeat`,
                WebkitMask: `url(/mcp.svg) center/contain no-repeat`,
                backgroundColor: "currentColor",
              }
        }
        role="img"
        aria-label={alt}
      />
    </Avatar>
  )
}
