"use client"

import { CheckCircle, ChevronDown, ChevronRight, Code, Loader2, XCircle } from "lucide-react"
import { useTranslations } from "next-intl"
import { useCallback, useMemo, useState } from "react"

import { Button } from "@/components/ui/button"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { cn } from "@/lib/utils"

type FunctionCallState = "loading" | "success" | "error" | "partial-call" | "call"

interface FunctionCallStatusProps {
  name: string
  args: string
  result?: string
  state?: FunctionCallState
}

const formatJsonData = (result: any): string => {
  try {
    const content = result?.content[0]?.text
    if (content) {
      const parsed = JSON.parse(content)
      return JSON.stringify(parsed, null, 2)
    }
  } catch { }
  return JSON.stringify(result, null, 2)
}

const isErrorResult = (result: any): boolean => {
  if (!result) return false
  try {
    if (typeof result === "object" && result !== null) {
      const content = result?.content[0]?.text
      if (content) {
        const parsed = JSON.parse(content)
        return parsed.success === false
      }
    }
  } catch { }
  return true
}

export function FunctionCallStatus({
  name,
  args,
  result,
  state = "success",
}: FunctionCallStatusProps) {
  const [isOpen, setIsOpen] = useState(false)

  const t = useTranslations("playground.functionCall")

  // Memoized values
  const hasResultError = useMemo(() => isErrorResult(result), [result])

  const actualState = useMemo((): FunctionCallState => {
    if (state === "error") return "error"
    if (hasResultError) return "error"
    return state
  }, [state, hasResultError])

  const statusIcon = useMemo(() => {
    switch (actualState) {
      case "loading":
      case "partial-call":
        return <Loader2 className="h-4 w-4 animate-spin text-blue-500" />
      case "success":
      case "call":
        return <CheckCircle className="h-4 w-4 text-green-500" />
      case "error":
        return <XCircle className="h-4 w-4 text-red-500" />
      default:
        return <CheckCircle className="h-4 w-4 text-green-500" />
    }
  }, [actualState])

  // Callbacks
  const handleToggle = useCallback((open: boolean) => {
    setIsOpen(open)
  }, [])

  return (
    <div className="flex justify-start">
      <div className={cn(
        "rounded-lg w-full",
        "bg-primary/8 dark:bg-primary/6",
        "border border-primary/15 shadow-none",
      )}
      >
        <div className="p-2">
          <Collapsible open={isOpen} onOpenChange={handleToggle}>
            <div className="flex items-center justify-between" onClick={() => setIsOpen(!isOpen)}>
              <div className="flex items-center gap-2">
                <Code className="h-4 w-4 text-muted-foreground" />
                <span className="font-medium text-sm">{name}</span>
              </div>
              <div className="flex items-center gap-2">
                {statusIcon}
                <CollapsibleTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-6 w-6 p-0"
                    aria-label={isOpen ? t("collapseDetails") : t("expandDetails")}
                  >
                    {isOpen
                      ? (
                          <ChevronDown className="h-3 w-3" />
                        )
                      : (
                          <ChevronRight className="h-3 w-3" />
                        )}
                  </Button>
                </CollapsibleTrigger>
              </div>
            </div>

            <CollapsibleContent className="mt-3">
              <div className="space-y-4">
                {args && (
                  <div>
                    <div className="text-xs font-medium mb-1 text-muted-foreground">
                      {t("parameters")}
                      :
                    </div>
                    <pre className="text-xs p-3 rounded-md font-mono leading-relaxed bg-background/30 border border-border/40 text-foreground/80 overflow-auto">
                      {String(formatJsonData(args))}
                    </pre>
                  </div>
                )}
                {result && (
                  <div>
                    <div className={cn(
                      "text-xs font-medium mb-1",
                      hasResultError ? "text-red-500 dark:text-red-400" : "text-muted-foreground",
                    )}
                    >
                      {t("result")}
                      :
                      {hasResultError && (
                        <span className="ml-1 text-xs bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 px-1.5 py-0.5 rounded-sm">
                          Error
                        </span>
                      )}
                    </div>
                    <pre className={cn(
                      "text-xs p-3 rounded-md font-mono leading-relaxed overflow-auto",
                      hasResultError
                        ? "bg-red-50/60 dark:bg-red-950/15 border border-red-200/60 dark:border-red-800/40 text-red-700 dark:text-red-300"
                        : "bg-background/30 border border-border/40 text-foreground/80",
                    )}
                    >
                      {formatJsonData(result) as string}
                    </pre>
                  </div>
                )}
              </div>
            </CollapsibleContent>
          </Collapsible>
        </div>
      </div>
    </div>
  )
}
