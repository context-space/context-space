export const siteName = "Context Space"
export const title = "Context Space - Tool-first Context Engineering Infrastructure"
export const description = "One Context, One Space, One Step"

export const baseURL = process.env.NEXT_PUBLIC_BASE_URL || "https://context.space"
export const baseAPIURL = process.env.NEXT_PUBLIC_BASE_API_URL || "https://api.context.space"

export const keywords = [
  "MCP",
  "Model Context Protocol",
  "Context Engineering",
  "AI Integration",
  "Workflow Automation",
  "Enterprise Security",
  "Cloud Native",
  "Open Source",
  "API Integration",
  "Team Collaboration",
]

export const turnstile = {
  siteKey: process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY!,
  enabled: process.env.NEXT_PUBLIC_TURNSTILE_ENABLED === "true",
}
