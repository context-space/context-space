import React from "react"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Link } from "@/i18n/navigation"
import { generateBreadcrumbStructuredData } from "@/lib/seo/structured-data"
import { StructuredData } from "./structured-data"

interface BreadcrumbsProps {
  items: { name: string, url: string }[]
  className?: string
}

export function Breadcrumbs({ items, className = "" }: BreadcrumbsProps) {
  const structuredData = generateBreadcrumbStructuredData(items)

  return (
    <>
      <StructuredData data={structuredData} />
      <Breadcrumb className={className}>
        <BreadcrumbList>
          {items.map((item, index) => (
            <React.Fragment key={item.url}>
              <BreadcrumbItem>
                {index === items.length - 1
                  ? (
                      <BreadcrumbPage>{item.name}</BreadcrumbPage>
                    )
                  : (
                      <BreadcrumbLink asChild>
                        <Link href={item.url}>{item.name}</Link>
                      </BreadcrumbLink>
                    )}
              </BreadcrumbItem>
              {index < items.length - 1 && <BreadcrumbSeparator />}
            </React.Fragment>
          ))}
        </BreadcrumbList>
      </Breadcrumb>
    </>
  )
}
