"use client"
import type { Message } from "ai"
import { useMemo } from "react"
import { MessageContent } from "@/components/integration/playground/prompt-kit/message"
import { cn } from "@/lib/utils"
import { FunctionCallStatus } from "./function-call-status"

interface MessageItemProps {
  message: Message
  className?: string
}

function UserMessage({ content, className }: { content: string, className?: string }) {
  if (!content) {
    return null
  }

  return (
    <div className={cn("flex justify-end", className)}>
      <div className={cn(
        "rounded-lg max-w-[80%]",
        "bg-neutral-50/60 dark:bg-white/[0.015]",
        "border border-base",
      )}
      >
        <MessageContent markdown={true} className="prose dark:prose-invert max-w-none text-sm bg-transparent">
          {content}
        </MessageContent>
      </div>
    </div>
  )
}

function AssistantMessage({ content, className }: { content: string, className?: string }) {
  if (!content) {
    return null
  }

  return (
    <div className={cn("flex justify-start", className)}>
      <div className={cn(
        "rounded-lg max-w-[80%]",
        "bg-neutral-50/60 dark:bg-white/[0.015]",
        "border border-base",
      )}
      >
        <MessageContent
          markdown={true}
          className="prose dark:prose-invert max-w-none text-sm bg-transparent"
        >
          {content}
        </MessageContent>
      </div>
    </div>
  )
}

export function MessageItem({ message, className }: MessageItemProps) {
  const isToolCall = useMemo(() => {
    return message.role === "assistant" && message.toolInvocations && message.toolInvocations.length > 0
  }, [message.role, message.toolInvocations])

  const isSystemMessage = message.role === "system"

  if (isSystemMessage) {
    return null // 不显示系统消息
  }

  if (isToolCall && message.toolInvocations) {
    return (
      <>
        <div className={cn("flex flex-col gap-4", className)}>
          {message.toolInvocations.map((toolInvocation, idx) => (
            <FunctionCallStatus
              key={toolInvocation.toolCallId || idx}
              name={toolInvocation.toolName}
              args={toolInvocation.args}
              result={(toolInvocation as any).result}
              state={toolInvocation.state === "result" ? "success" : toolInvocation.state}
            />
          ))}
        </div>
        {message.content && <AssistantMessage content={message.content} />}
      </>
    )
  }

  return (
    <div className={cn("flex flex-col gap-4", className)}>
      {message.role === "user" && <UserMessage content={message.content} />}
      {message.role === "assistant" && <AssistantMessage content={message.content} />}
    </div>
  )
}
