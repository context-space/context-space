import { createProxyHandlers } from "@/lib/api-proxy"

const handlers = createProxyHandlers({
  basePathPrefix: "",
  loggerTag: "proxy-api",
})

export const { GET, POST, PUT, DELETE, PATCH } = handlers