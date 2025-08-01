import type { MDXComponents } from "mdx/types"

export function useMDXComponents(components: MDXComponents): MDXComponents {
  return {
    img: props => (
      <img {...props} className="rounded-xl" />
      // <Image
      //   loading="lazy"
      //   quality={70}
      //   width={1000}
      //   height={1000}
      //   className="rounded-xl"
      //   {...(props as ImageProps)}
      // />
    ),
    wrapper: ({ children }) => (
      <div className="prose prose-slate dark:prose-invert max-w-none">
        <div className="max-w-4xl mx-auto py-12 px-4">
          {children}
        </div>
      </div>
    ),
    ...components,
  }
}
