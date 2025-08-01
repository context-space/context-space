interface Props {
  title: string
  children: React.ReactNode
}

export function Section({ title, children }: Props) {
  return (
    <div>
      <h2 className="text-2xl font-semibold text-neutral-900 dark:text-white mb-6">
        {title}
      </h2>
      {children}
    </div>
  )
}
