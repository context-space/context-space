"use client"

import { useSetAtom } from "jotai"
import { ChevronDown } from "lucide-react"
import { useTranslations } from "next-intl"
import { useState } from "react"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { ScrollArea } from "@/components/ui/scroll-area"
import { cn } from "@/lib/utils"
import { selectedOperationAtom, triggerInputUpdateAtom } from "../store"

export interface Operation {
  identifier: string
  name: string
  description: string
  required_permissions: Array<{ identifier: string }>
}

interface OperationsSectionProps {
  operations: Operation[]
}

export function Operations({ operations }: OperationsSectionProps) {
  const [isOpen, setIsOpen] = useState(false)
  const t = useTranslations()
  const setSelectedOperation = useSetAtom(selectedOperationAtom)
  const setTriggerUpdate = useSetAtom(triggerInputUpdateAtom)

  const handleOperationClick = (operation: Operation) => {
    setSelectedOperation(`${operation.name} `)
    // Trigger an update to notify playground to update its input
    setTriggerUpdate(prev => prev + 1)
  }

  return (
    <Card className="pb-0 bg-white/60 dark:bg-white/[0.02] border-base backdrop-blur-sm shadow-none">
      <Collapsible open={isOpen} onOpenChange={setIsOpen}>
        <CardHeader onClick={() => setIsOpen(!isOpen)}>
          <div className="space-y-1">
            <CardTitle className="text-lg font-semibold flex items-center justify-between">
              <span>
                {t("operations.title")}
              </span>
              <CollapsibleTrigger asChild>
                <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                  <ChevronDown className={cn("h-4 w-4 transition-transform duration-200", isOpen && "rotate-180")} />
                  <span className="sr-only">{t("common.toggle")}</span>
                </Button>
              </CollapsibleTrigger>
            </CardTitle>
            <p className="text-sm text-muted-foreground my-3">
              {t("operations.availableCount", { count: operations.length })}
              {", "}
              {t("operations.clickToExecute")}
            </p>
          </div>
        </CardHeader>
        <CollapsibleContent>
          <CardContent className="pt-0 px-2  border-t border-base">
            <ScrollArea style={{ height: (operations.length > 4 ? 4 : operations.length) * 125 }}>
              <div className="space-y-3 px-5 py-5">
                {operations.map(operation => (
                  <Card
                    key={operation.identifier}
                    className="py-4 bg-neutral-50/50 dark:bg-white/[0.02] hover:bg-neutral-100/50 dark:hover:bg-white/[0.03] transition-all duration-300 border-base shadow-none cursor-pointer"
                    onClick={() => handleOperationClick(operation)}
                  >
                    <CardContent>
                      <div className="flex items-center gap-3 mb-2">
                        <span className="font-medium text-sm">
                          {operation.name}
                        </span>
                        <div className="h-px flex-1 border-b border-base" />
                        <div className="flex gap-1 flex-wrap">
                          {operation.required_permissions.map(permission => (
                            <Badge
                              key={permission.identifier}
                              variant="outline"
                              className="text-xs"
                            >
                              {permission.identifier}
                            </Badge>
                          ))}
                        </div>
                      </div>
                      <p className="text-muted-foreground text-xs leading-relaxed">
                        {operation.description}
                      </p>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </ScrollArea>
          </CardContent>
        </CollapsibleContent>
      </Collapsible>
    </Card>
  )
}
