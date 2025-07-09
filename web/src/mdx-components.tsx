import type { MDXComponents } from "mdx/types"
import type { ImageProps } from "next/image"
import Image from "next/image"
import { cn } from "@/lib/utils"

export function useMDXComponents(components: MDXComponents): MDXComponents {
  return {
    // Typography components with proper styling
    h1: ({ children, className, ...props }) => (
      <h1
        className={cn(
          "scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-8",
          className,
        )}
        {...props}
      >
        {children}
      </h1>
    ),
    h2: ({ children, className, ...props }) => (
      <h2
        className={cn(
          "scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0 mt-8 mb-4",
          className,
        )}
        {...props}
      >
        {children}
      </h2>
    ),
    h3: ({ children, className, ...props }) => (
      <h3
        className={cn(
          "scroll-m-20 text-2xl font-semibold tracking-tight mt-6 mb-3",
          className,
        )}
        {...props}
      >
        {children}
      </h3>
    ),
    p: ({ children, className, ...props }) => (
      <p
        className={cn(
          "leading-7 [&:not(:first-child)]:mt-6",
          className,
        )}
        {...props}
      >
        {children}
      </p>
    ),
    ul: ({ children, className, ...props }) => (
      <ul
        className={cn(
          "my-6 ml-6 list-disc [&>li]:mt-2",
          className,
        )}
        {...props}
      >
        {children}
      </ul>
    ),
    ol: ({ children, className, ...props }) => (
      <ol
        className={cn(
          "my-6 ml-6 list-decimal [&>li]:mt-2",
          className,
        )}
        {...props}
      >
        {children}
      </ol>
    ),
    li: ({ children, className, ...props }) => (
      <li className={cn("", className)} {...props}>
        {children}
      </li>
    ),
    blockquote: ({ children, className, ...props }) => (
      <blockquote
        className={cn(
          "mt-6 border-l-2 pl-6 italic",
          className,
        )}
        {...props}
      >
        {children}
      </blockquote>
    ),
    code: ({ children, className, ...props }) => (
      <code
        className={cn(
          "relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold",
          className,
        )}
        {...props}
      >
        {children}
      </code>
    ),
    pre: ({ children, className, ...props }) => (
      <pre
        className={cn(
          "mb-4 mt-6 overflow-x-auto rounded-lg border bg-zinc-950 py-4 dark:bg-zinc-900",
          className,
        )}
        {...props}
      >
        {children}
      </pre>
    ),
    img: props => (
      <Image
        sizes="100vw"
        loading="lazy"
        quality={70}
        style={{ width: "100%", height: "auto" }}
        {...(props as ImageProps)}
      />
    ),
    // Add wrapper for better layout
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
