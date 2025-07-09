"use client"

import Logo from "@/components/common/logo"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { cn } from "@/lib/utils"

interface ProviderAvatarProps {
  src?: string
  alt: string
  className?: string
}

export function ProviderAvatar({ src, alt, className }: ProviderAvatarProps) {
  return (
    <Avatar className={cn("bg-transparent rounded-none transition-transform duration-200 group-hover:scale-102", className)}>
      <AvatarImage src={src} alt={alt} />
      <AvatarFallback className="bg-transparent rounded-none">
        <Logo />
      </AvatarFallback>
    </Avatar>
  )
}
