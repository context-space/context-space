"use client"

import type { Message } from "@ai-sdk/react"
import { useChat } from "@ai-sdk/react"
import { useAtom, useSetAtom } from "jotai"
import { Loader2, Square, Triangle } from "lucide-react"
import { useTranslations } from "next-intl"
import React, { useCallback, useEffect, useMemo, useRef, useState } from "react"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { useAuth } from "@/hooks/use-auth"
import { useRouter } from "@/i18n/navigation"
import { clientLogger, cn } from "@/lib/utils"
import { selectedOperationAtom, triggerInputUpdateAtom } from "../store"
import { MessageItem } from "./message-item"
import { PlaygroundHeader } from "./playground-header"
import { ChatContainerContent, ChatContainerRoot, ChatContainerScrollAnchor } from "./prompt-kit/chat-container"
import {
  PromptInput,
  PromptInputAction,
  PromptInputActions,
  PromptInputTextarea,
} from "./prompt-kit/prompt-input"
import { ScrollButton } from "./prompt-kit/scroll-button"

interface PlaygroundProps {
  provider: string
  authType?: "oauth" | "apikey" | "none"
  isConnected?: boolean
  providerId?: string
}

const createId = () => Math.random().toString(36).slice(2) + Date.now()

const playgroundLogger = clientLogger.withTag("playground")

// localStorage utility functions
const getChatHistoryKey = (provider: string, userId?: string) => {
  return userId ? `chat_history_${provider}_${userId}` : `chat_history_${provider}_anonymous`
}

const getInputHistoryKey = (provider: string, userId?: string) => {
  return userId ? `input_history_${provider}_${userId}` : `input_history_${provider}_anonymous`
}

const saveChatHistory = (provider: string, messages: Message[], userId?: string) => {
  try {
    const key = getChatHistoryKey(provider, userId)
    localStorage.setItem(key, JSON.stringify(messages))
  } catch (error) {
    playgroundLogger.error("Failed to save chat history", { error })
  }
}

const loadChatHistory = (provider: string, userId?: string): Message[] | null => {
  try {
    const key = getChatHistoryKey(provider, userId)
    const stored = localStorage.getItem(key)
    return stored ? JSON.parse(stored) : null
  } catch (error) {
    playgroundLogger.error("Failed to load chat history", { error })
    return null
  }
}

const saveInputHistory = (provider: string, inputHistory: string[], userId?: string) => {
  try {
    const key = getInputHistoryKey(provider, userId)
    localStorage.setItem(key, JSON.stringify(inputHistory))
  } catch (error) {
    playgroundLogger.error("Failed to save input history", { error })
  }
}

const loadInputHistory = (provider: string, userId?: string): string[] => {
  try {
    const key = getInputHistoryKey(provider, userId)
    const stored = localStorage.getItem(key)
    return stored ? JSON.parse(stored) : []
  } catch (error) {
    playgroundLogger.error("Failed to load input history", { error })
    return []
  }
}

const clearUserChatData = (userId?: string) => {
  try {
    if (typeof window === "undefined") return

    const keysToRemove: string[] = []

    // Find all localStorage keys related to chat history for this user
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && (
        (userId && (key.includes(`_${userId}`) || key.includes("_anonymous")))
        || (!userId && key.includes("_anonymous"))
      )) {
        if (key.startsWith("chat_history_") || key.startsWith("input_history_")) {
          keysToRemove.push(key)
        }
      }
    }

    // Remove all found keys
    keysToRemove.forEach(key => localStorage.removeItem(key))
    playgroundLogger.info("Cleared chat data from localStorage", { keysRemoved: keysToRemove.length })
  } catch (error) {
    playgroundLogger.error("Failed to clear chat data", { error })
  }
}

export function Playground({ provider, authType, isConnected, providerId }: PlaygroundProps) {
  const t = useTranslations()
  const { user, isAuthenticated } = useAuth()
  const [selectedOperation] = useAtom(selectedOperationAtom)
  const [triggerUpdate] = useAtom(triggerInputUpdateAtom)
  const setSelectedOperation = useSetAtom(selectedOperationAtom)
  const setTriggerUpdate = useSetAtom(triggerInputUpdateAtom)
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const router = useRouter()

  // Input history state
  const [inputHistory, setInputHistory] = useState<string[]>([])
  const [historyIndex, setHistoryIndex] = useState(-1)

  // Check if the integration requires connection
  const requiresConnection = authType && authType !== "none"

  const checkAuthAndConnectionBeforeSubmit = useCallback(() => {
    // Check authentication first
    if (!isAuthenticated) {
      toast.error(t("integrations.connect.loginRequired"), {
        description: t("integrations.connect.loginRequiredDescription"),
        action: {
          label: t("header.login"),
          onClick: () => {
            const currentPath = encodeURIComponent(window.location.pathname)
            router.push(`/login?from=${currentPath}`)
          },
        },
      })
      return false
    }

    // Check connection if required
    if (requiresConnection && !isConnected) {
      toast.error(t("playground.connectionRequired"), {
        description: t("playground.connectionRequiredDescription"),
      })
      return false
    }

    return true
  }, [isAuthenticated, requiresConnection, isConnected, t])

  const DEFAULT_MESSAGES: Message[] = useMemo(() => [
    {
      id: createId(),
      role: "system",
      content: "You are a specialized AI assistant focused solely on using and demonstrating available tools. Your core directives:\n\n1. ONLY answer questions that can be addressed using the available tools\n2. If a question is unrelated to the tools or cannot be answered using them, politely explain that you can only help with tool-related queries\n3. For any query, first explain which tools are relevant and how they will help\n4. When using tools, clearly explain what you're doing and why\n5. If the user asks questions unrelated to tools (e.g., general knowledge, casual conversation), redirect them to tool-related topics and explain what tools can do",
    },
    {
      id: createId(),
      role: "assistant",
      content: t("playground.welcomeMessage"),
    },
  ], [t])
  const [toolContext, setToolContext] = useState<string | undefined>()

  const {
    messages,
    input,
    handleInputChange,
    handleSubmit,
    error,
    isLoading,
    stop,
    setMessages,
    setInput,
  } = useChat({
    initialMessages: DEFAULT_MESSAGES,
    initialInput: t("playground.defaultInput"),
    api: `/api/chat?id=${provider}`,
    onFinish: (message) => {
      // Extract tool context from message parts
      const context = message.toolInvocations?.reduce((acc, toolInvocation) => {
        return acc + (JSON.stringify((toolInvocation as any).result) || "")
      }, "")
      setToolContext(context)

      // Auto-focus the input after message completion
      setTimeout(() => {
        if (textareaRef.current) {
          textareaRef.current.focus()
        }
      }, 100)
    },
  })

  // Save chat history to localStorage whenever messages change (only if user is authenticated)
  useEffect(() => {
    if (user && messages.length > DEFAULT_MESSAGES.length) {
      saveChatHistory(provider, messages, user.id)
    }
  }, [messages, provider, user, DEFAULT_MESSAGES.length])

  // Save input history to localStorage whenever it changes (only if user is authenticated)
  useEffect(() => {
    if (user && inputHistory.length > 0) {
      saveInputHistory(provider, inputHistory, user.id)
    }
  }, [inputHistory, provider, user])

  // Load chat and input history on mount and when provider/user changes
  useEffect(() => {
    if (user) {
      // User is logged in, load their chat history
      const loadedChatHistory = loadChatHistory(provider, user.id)
      const loadedInputHistory = loadInputHistory(provider, user.id)

      if (loadedChatHistory && loadedChatHistory.length > DEFAULT_MESSAGES.length) {
        setMessages(loadedChatHistory)
      } else {
        setMessages(DEFAULT_MESSAGES)
      }

      setInputHistory(loadedInputHistory)

      // 只在没有输入历史时才设置默认输入
      if (!loadedInputHistory || loadedInputHistory.length === 0) {
        setInput(t("playground.defaultInput"))
      } else {
        setInput("")
      }
    } else {
      // User is not logged in, reset to default state
      setMessages(DEFAULT_MESSAGES)
      setInputHistory([])
      setInput(t("playground.defaultInput"))
    }

    setHistoryIndex(-1)
    setSelectedOperation("")
    setTriggerUpdate(0)
  }, [provider, user, setMessages, setInput, setSelectedOperation, setTriggerUpdate, t, DEFAULT_MESSAGES])

  // Reset selected operation when provider changes
  useEffect(() => {
    setSelectedOperation("")
    setTriggerUpdate(0)
    // 不再在这里设置默认输入，让上面的 useEffect 处理输入逻辑
  }, [provider, setSelectedOperation, setTriggerUpdate])

  // Update input when an operation is selected
  useEffect(() => {
    if (selectedOperation && triggerUpdate > 0) {
      setInput(selectedOperation)
      // Focus the textarea after setting the input
      setTimeout(() => {
        if (textareaRef.current) {
          textareaRef.current.focus()
          // Move cursor to the end of the text
          const length = selectedOperation.length
          textareaRef.current.setSelectionRange(length, length)
        }
      }, 0)
    }
  }, [selectedOperation, triggerUpdate, setInput])

  const handleClear = () => {
    setMessages(DEFAULT_MESSAGES)
    // Clear localStorage for this provider (only if user is authenticated)
    if (user) {
      try {
        const chatKey = getChatHistoryKey(provider, user.id)
        localStorage.removeItem(chatKey)
        playgroundLogger.info("Cleared chat history from localStorage", { provider, userId: user.id })
      } catch (error) {
        playgroundLogger.error("Failed to clear chat history from localStorage", { error })
      }
    }
  }

  const onSubmit = () => {
    // Check authentication and connection before submitting
    if (!checkAuthAndConnectionBeforeSubmit()) {
      return
    }

    // Add current input to history when submitting
    if (input.trim() && !inputHistory.includes(input.trim())) {
      setInputHistory(prev => [input.trim(), ...prev].slice(0, 50)) // Keep last 50 inputs
    }
    setHistoryIndex(-1) // Reset history index
    handleSubmit()
  }

  const onValueChange = (value: string) => {
    handleInputChange({ target: { value } } as React.ChangeEvent<HTMLTextAreaElement>)
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "ArrowUp" && !e.shiftKey) {
      e.preventDefault()
      if (inputHistory.length > 0) {
        const newIndex = Math.min(historyIndex + 1, inputHistory.length - 1)
        setHistoryIndex(newIndex)
        setInput(inputHistory[newIndex])
      }
    } else if (e.key === "ArrowDown" && !e.shiftKey) {
      e.preventDefault()
      if (historyIndex > 0) {
        const newIndex = historyIndex - 1
        setHistoryIndex(newIndex)
        setInput(inputHistory[newIndex])
      } else if (historyIndex === 0) {
        setHistoryIndex(-1)
        setInput("")
      }
    }
  }

  return (
    <Card className={cn(
      "flex flex-col h-[calc(100vh-15rem)] min-h-[500px] gap-0",
      "bg-white/60 dark:bg-white/[0.02] border-base relative",
    )}
    >
      <div className="pb-0 px-5">
        <PlaygroundHeader onClear={handleClear} />
      </div>

      <CardContent className="flex-1 flex flex-col p-4 gap-4 overflow-hidden relative border-t border-base">
        <ChatContainerRoot className="flex-1 px-1 relative">
          <ChatContainerContent className="space-y-4">
            {messages.map((message, index) => (
              <MessageItem
                key={message.id || `message-${index}`}
                message={message}
              />
            ))}
            {error && (
              <div className="flex items-center gap-2 text-red-500 text-sm p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-100 dark:border-red-900/30">
                {error.message}
              </div>
            )}
            <ChatContainerScrollAnchor />
          </ChatContainerContent>

          <div className="absolute bottom-4 right-4 z-10">
            <ScrollButton />
          </div>
        </ChatContainerRoot>
      </CardContent>

      <div className="flex items-center px-5">
        <PromptInput
          isLoading={isLoading}
          value={input}
          onValueChange={onValueChange}
          onSubmit={onSubmit}
          className="relative rounded-lg border border-base w-full bg-neutral-50/60 dark:bg-white/[0.015]"
        >
          {isLoading && (
            <div className="absolute left-2 top-1/2 transform -translate-y-1/2 z-10">
              <Loader2 className="h-4 w-4 font-bold animate-spin text-primary" />
            </div>
          )}
          <PromptInputTextarea
            ref={textareaRef}
            className={cn(`text-black dark:text-white scrollbar-hidden pr-8 px-2`)}
            disabled={isLoading}
            onKeyDown={handleKeyDown}
          />
          <PromptInputActions className="absolute right-0 top-1/2 transform -translate-y-1/2">
            {isLoading
              ? (
                <PromptInputAction tooltip={t("playground.stop")}>
                  <Button
                    type="button"
                    size="sm"
                    variant="ghost"
                    onClick={stop}
                    className="h-8 w-8 p-0"
                  >
                    <Square className="h-4 w-4 text-primary" />
                  </Button>
                </PromptInputAction>
              )
              : (
                <PromptInputAction tooltip={t("playground.send")}>
                  <Button
                    type="submit"
                    size="sm"
                    disabled={!input.trim()}
                    onClick={onSubmit}
                    variant="ghost"
                    className="h-8 w-8 p-0"
                  >
                    <Triangle className="h-4 w-4 rotate-90 text-primary" />
                  </Button>
                </PromptInputAction>
              )}
          </PromptInputActions>
        </PromptInput>
      </div>
    </Card>
  )
}

// Export the clear function for use in auth provider
export { clearUserChatData }
