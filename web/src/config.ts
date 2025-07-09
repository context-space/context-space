export const siteName = process.env.NEXT_PUBLIC_SITE_NAME || "Context Space"
export const title = process.env.NEXT_PUBLIC_TITLE || "Context Space - Ultimate Context Engineering Infrastructure"
export const description = process.env.NEXT_PUBLIC_DESCRIPTION || "One Context, One Space, One Step"

// meta base url
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

// Turnstile configuration
export const turnstile = {
  siteKey: process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY!,
  enabled: process.env.NEXT_PUBLIC_TURNSTILE_ENABLED === "true",
}
