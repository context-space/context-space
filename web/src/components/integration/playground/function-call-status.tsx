"use client"

import { Check, CheckCircle, ChevronDown, ChevronRight, Code, Copy, Loader2, XCircle } from "lucide-react"
import { useTranslations } from "next-intl"
import { useTheme } from "next-themes"
import { useCallback, useMemo, useState } from "react"

import { CodeBlock, CodeBlockCode } from "@/components/integration/playground/prompt-kit/code-block"
import { Button } from "@/components/ui/button"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { clientLogger, cn } from "@/lib/utils"

const COPY_SUCCESS_TIMEOUT = 2000
const ERROR_KEYWORDS = ["error", "failed"] as const

type FunctionCallState = "loading" | "success" | "error" | "partial-call" | "call"

interface FunctionCallStatusProps {
  name: string
  args: unknown
  result?: unknown
  state?: FunctionCallState
}

interface JsonDisplayProps {
  data: unknown
  label: string
  name: string
  isErrorResult: boolean
  formattedData: string
  isStringData: boolean
  onCopy: (text: string, key: string) => Promise<void>
  isCopied: boolean
  theme?: string
}

interface ErrorResult {
  isError: boolean
  [key: string]: unknown
}

// Utility functions
const formatJsonData = (data: unknown): string => {
  if (typeof data === "string") {
    try {
      const parsed = JSON.parse(data)
      return JSON.stringify(parsed, null, 2)
    } catch {
      return data
    }
  }
  return JSON.stringify(data, null, 2)
}

const isErrorResult = (result: unknown): result is ErrorResult => {
  if (!result) return false

  if (typeof result === "object" && result !== null && "isError" in result) {
    return (result as ErrorResult).isError === true
  }

  if (typeof result === "string") {
    try {
      const parsed = JSON.parse(result) as ErrorResult
      return parsed.isError === true
    } catch {
      return ERROR_KEYWORDS.some(keyword =>
        result.toLowerCase().includes(keyword),
      )
    }
  }

  return false
}

const isStringData = (data: unknown): data is string => {
  return typeof data === "string"
    && !data.startsWith("{")
    && !data.startsWith("[")
}

// Components
function JsonDisplay({
  data,
  label,
  name,
  isErrorResult: isError,
  formattedData,
  isStringData: isString,
  onCopy,
  isCopied,
  theme,
}: JsonDisplayProps) {
  const copyKey = `${name}-${label}`

  const handleCopy = useCallback(() => {
    onCopy(formattedData, copyKey)
  }, [formattedData, copyKey, onCopy])

  return (
    <div>
      <div className="flex items-center justify-between mb-2">
        <div className={cn(
          "text-xs font-medium mb-1",
          isError ? "text-red-500 dark:text-red-400" : "text-muted-foreground",
        )}
        >
          {label}
          :
          {isError && (
            <span className="ml-1 text-xs bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 px-1.5 py-0.5 rounded-sm">
              Error
            </span>
          )}
        </div>
        <Button
          variant="ghost"
          size="sm"
          className="h-6 w-6 p-0 opacity-60 hover:opacity-100 transition-all"
          onClick={handleCopy}
          aria-label={`Copy ${label.toLowerCase()}`}
        >
          {isCopied
            ? (
                <Check className="h-3 w-3 text-green-500" />
              )
            : (
                <Copy className="h-3 w-3" />
              )}
        </Button>
      </div>

      {isString
        ? (
            <div className={cn(
              "text-xs p-3 rounded-md font-mono leading-relaxed",
              isError
                ? "bg-red-50/60 dark:bg-red-950/15 border border-red-200/60 dark:border-red-800/40 text-red-700 dark:text-red-300"
                : "bg-background/30 border border-border/40 text-foreground/80",
            )}
            >
              {formattedData}
            </div>
          )
        : (
            <CodeBlock className={cn(
              "text-xs max-h-96 shadow-none",
              isError && "border-red-200/60 dark:border-red-800/40",
            )}
            >
              <CodeBlockCode
                code={formattedData}
                theme={theme === "dark" ? "vitesse-dark" : "vitesse-light"}
                language="json"
                className="text-xs [&>pre]:py-3 [&>pre]:px-3 [&>pre]:leading-relaxed"
              />
            </CodeBlock>
          )}
    </div>
  )
}

const functionCallLogger = clientLogger.withTag("function-call")

export function FunctionCallStatus({
  name,
  args,
  result,
  state = "success",
}: FunctionCallStatusProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [copiedStates, setCopiedStates] = useState<Record<string, boolean>>({})

  const t = useTranslations("playground.functionCall")
  const { theme } = useTheme()

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
  const copyToClipboard = useCallback(async (text: string, key: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopiedStates(prev => ({ ...prev, [key]: true }))

      setTimeout(() => {
        setCopiedStates(prev => ({ ...prev, [key]: false }))
      }, COPY_SUCCESS_TIMEOUT)
    } catch (error) {
      functionCallLogger.error("Failed to copy text", { error })
    }
  }, [])

  const handleToggle = useCallback((open: boolean) => {
    setIsOpen(open)
  }, [])

  // Prepare data for JsonDisplay components
  const argsData = useMemo(() => {
    if (!args) return null

    return {
      formattedData: formatJsonData(args),
      isStringData: isStringData(args),
      isErrorResult: false,
    }
  }, [args])

  const resultData = useMemo(() => {
    if (!result) return null

    return {
      formattedData: formatJsonData(result),
      isStringData: isStringData(result),
      isErrorResult: hasResultError,
    }
  }, [result, hasResultError])

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
            <div className="flex items-center justify-between">
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
                    aria-label={isOpen ? "Collapse details" : "Expand details"}
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
                {argsData && (
                  <JsonDisplay
                    data={args}
                    label={t("parameters")}
                    name={name}
                    isErrorResult={argsData.isErrorResult}
                    formattedData={argsData.formattedData}
                    isStringData={argsData.isStringData}
                    onCopy={copyToClipboard}
                    isCopied={copiedStates[`${name}-${t("parameters")}`] ?? false}
                    theme={theme}
                  />
                )}
                {resultData && (
                  <JsonDisplay
                    data={result}
                    label={t("result")}
                    name={name}
                    isErrorResult={resultData.isErrorResult}
                    formattedData={resultData.formattedData}
                    isStringData={resultData.isStringData}
                    onCopy={copyToClipboard}
                    isCopied={copiedStates[`${name}-${t("result")}`] ?? false}
                    theme={theme}
                  />
                )}
              </div>
            </CollapsibleContent>
          </Collapsible>
        </div>
      </div>
    </div>
  )
}
