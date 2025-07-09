export type MessageRole = "user" | "assistant" | "system" | "function"

export interface BaseMessage {
  role: MessageRole
  content: string | null
}

export interface UserMessage extends BaseMessage {
  role: "user"
  content: string
}

export interface AssistantMessage extends BaseMessage {
  role: "assistant"
  content: string
  function_call?: {
    name: string
    arguments: string
  }
}

export interface SystemMessage extends BaseMessage {
  role: "system"
  content: string
}

export interface FunctionMessage extends BaseMessage {
  role: "function"
  name: string
  content: string
}

export type Message = UserMessage | AssistantMessage | SystemMessage | FunctionMessage

export interface FunctionCallInfo {
  name: string
  status: "start" | "running" | "end"
  result?: string
}

export interface MCPMessage {
  type: "message" | "function"
  message: string
  done: boolean
}

export interface MCPError {
  type: "error"
  error: string
}

export type MCPResponse = MCPMessage | MCPError
