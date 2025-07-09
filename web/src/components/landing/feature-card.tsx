interface FeatureCardProps {
  Icon: React.ElementType
  title: string
  description: string
}

export function FeatureCard({ Icon, title, description }: FeatureCardProps) {
  return (
    <div className="group relative p-7 rounded-lg bg-white/70 dark:bg-white/[0.02] border border-neutral-200/60 dark:border-white/[0.05] hover:border-neutral-300 dark:hover:border-white/10 transition-all duration-300">
      <div className="absolute inset-0 rounded-lg bg-gradient-to-b from-white/50 to-transparent dark:from-white/[0.02] dark:to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
      <div className="relative flex items-center gap-5">
        <Icon />
        <h3 className="text-[15px] font-medium tracking-wide text-neutral-900 dark:text-white">
          {title}
        </h3>
      </div>
      <div className="relative mt-5">
        <p className="text-[13.5px] text-neutral-500 dark:text-gray-400 leading-relaxed">
          {description}
        </p>
        <div className="absolute -inset-x-7 -inset-y-2 bg-gradient-to-r from-transparent via-white/50 to-transparent dark:via-white/[0.02] opacity-0 group-hover:opacity-100 transition-opacity"></div>
      </div>
    </div>
  )
}
