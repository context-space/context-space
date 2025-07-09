import type { Locale } from "@/i18n/routing"
import { notFound } from "next/navigation"
import { locales } from "@/i18n/routing"

interface Props {
  params: Promise<{ locale: Locale }>
}

// 生成静态参数，启用静态站点生成 (SSG)
export async function generateStaticParams() {
  return locales.map(locale => ({ locale }))
}

// 可选：启用静态导出
export const dynamic = "force-static"

export default async function MDXPage({ params }: Props) {
  const { locale } = await params

  try {
    const Content = (await import(`./content.${locale}.mdx`)).default
    return <Content />
  } catch {
    notFound()
  }
}
