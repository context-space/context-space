import { ImageResponse } from "next/og"
import { description, siteName } from "@/config"

// Primary color in RGB (converted from oklch(0.6947 0.0525 283.69))
const primaryColor = "rgb(148, 150, 228)"
const primaryColorAlpha = (alpha: number) => addAlpha(primaryColor, alpha)

// Helper function to add alpha transparency to any color
function addAlpha(color: string, alpha: number): string {
  // Handle rgb() format
  if (color.startsWith("rgb(")) {
    const rgbValues = color.slice(4, -1)
    return `rgba(${rgbValues}, ${alpha})`
  }

  // Handle hex format
  if (color.startsWith("#")) {
    const hex = color.slice(1)
    const r = Number.parseInt(hex.slice(0, 2), 16)
    const g = Number.parseInt(hex.slice(2, 4), 16)
    const b = Number.parseInt(hex.slice(4, 6), 16)
    return `rgba(${r}, ${g}, ${b}, ${alpha})`
  }

  // Fallback: return original color with alpha appended
  return `${color}${Math.round(alpha * 255).toString(16).padStart(2, "0")}`
}

// Feature item component
function FeatureItem({ icon, label, color }: { icon: string, label: string, color: string }) {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gap: "8px",
      }}
    >
      <div
        style={{
          fontSize: "24px",
          color,
        }}
      >
        {icon}
      </div>
      <span
        style={{
          fontSize: "14px",
          color: "rgba(255, 255, 255, 0.6)",
          fontWeight: "500",
          textTransform: "capitalize",
        }}
      >
        {label}
      </span>
    </div>
  )
}

// Default platform OG image
function DefaultOGImage({
  title,
  description: desc,
}: {
  title: string
  description: string
}) {
  return (
    <div
      style={{
        height: "100%",
        width: "100%",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        background: "#1a1a1a",
        backgroundImage: `radial-gradient(ellipse 80% 80% at 50% -30%, ${primaryColorAlpha(0.3)}, rgba(255, 255, 255, 0))`,
        padding: "48px",
      }}
    >
      {/* Card Container */}
      <div
        style={{
          position: "relative",
          display: "flex",
          flexDirection: "column",
          width: "100%",
          maxWidth: "1000px",
          height: "480px",
          padding: "48px",
          borderRadius: "24px",
          backgroundColor: "rgba(255, 255, 255, 0.05)",
          border: `1px solid ${primaryColorAlpha(0.2)}`,
          boxShadow: `0 25px 50px -12px ${primaryColorAlpha(0.15)}`,
        }}
      >
        {/* Header section */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            marginBottom: "40px",
          }}
        >
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: "24px",
            }}
          >
            {/* Logo */}
            <div
              style={{
                width: "64px",
                height: "64px",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <img
                src="https://context.space/logo.svg"
                alt="Context Space logo"
                style={{
                  width: "48px",
                  height: "48px",
                  objectFit: "contain",
                }}
              />
            </div>
            {/* Title */}
            <h1
              style={{
                fontSize: "48px",
                fontWeight: "700",
                color: "#ffffff",
                margin: 0,
                lineHeight: "1.1",
              }}
            >
              {title}
            </h1>
          </div>

          {/* Platform badge */}
          <div
            style={{
              padding: "8px 12px",
              borderRadius: "8px",
              backgroundColor: addAlpha(primaryColor, 0.15),
              border: `1px solid ${addAlpha(primaryColor, 0.4)}`,
              color: primaryColor,
              fontSize: "16px",
              fontWeight: "600",
            }}
          >
            Alpha
          </div>
        </div>

        {/* Main content area */}
        <div
          style={{
            flex: 1,
            display: "flex",
            flexDirection: "column",
            gap: "32px",
          }}
        >
          {/* Main description */}
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              gap: "16px",
            }}
          >
            <p
              style={{
                fontSize: "28px",
                lineHeight: "1.4",
                color: "#ffffff",
                margin: 0,
                fontWeight: "600",
              }}
            >
              Tool-first Context Engineering Infrastructure for AI Agents
            </p>
            <p
              style={{
                fontSize: "20px",
                lineHeight: "1.5",
                color: "rgba(255, 255, 255, 0.7)",
                margin: 0,
              }}
            >
              {desc || "Connect, orchestrate, and scale your AI workflows with enterprise-grade security and seamless integrations."}
            </p>
          </div>

          {/* Platform features */}
          <div
            style={{
              display: "flex",
              gap: "40px",
              alignItems: "center",
            }}
          >
            <FeatureItem
              icon="🔗"
              label="MCP Protocol"
              color="#22c55e"
            />
            <FeatureItem
              icon="⚡"
              label="Fast Integration"
              color={primaryColor}
            />
            <FeatureItem
              icon="🛡️"
              label="Enterprise Security"
              color="#f59e0b"
            />
            <FeatureItem
              icon="🚀"
              label="Cloud Native"
              color={primaryColor}
            />
          </div>
        </div>

        {/* Bottom section */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            paddingTop: "24px",
            borderTop: "1px solid rgba(255, 255, 255, 0.1)",
          }}
        >
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: "12px",
              color: "rgba(255, 255, 255, 0.5)",
              fontSize: "18px",
              fontWeight: "500",
            }}
          >
            <span>🌐</span>
            <span>context.space</span>
          </div>
          <div
            style={{
              display: "flex",
              gap: "24px",
              alignItems: "center",
              color: "rgba(255, 255, 255, 0.4)",
              fontSize: "14px",
            }}
          >
            <span>Open Source</span>
            <span>•</span>
            <span>AI-First</span>
            <span>•</span>
            <span>Developer Friendly</span>
          </div>
        </div>
      </div>
    </div>
  )
}

// Integration-specific OG image
function IntegrationOGImage({
  title,
  description: desc,
  integrationName,
  integrationIcon,
  categories,
  authType,
  connectionStatus,
}: {
  title: string
  description: string
  integrationName: string
  integrationIcon: string | null
  categories: string[]
  authType: "oauth" | "apikey" | "none" | null
  connectionStatus: "connected" | "free" | "unconnected" | null
}) {
  // Get status badge content and color
  const getStatusInfo = () => {
    switch (connectionStatus) {
      case "connected":
        return { text: "Connected", color: "#22c55e" }
      case "free":
        return { text: "Free", color: "#3b82f6" }
      case "unconnected":
        return { text: "Available", color: "#f59e0b" }
      default:
        return { text: "Beta", color: primaryColor }
    }
  }

  // Get auth info
  const getAuthInfo = () => {
    switch (authType) {
      case "oauth":
        return { icon: "🔐", label: "OAuth Auth", color: "#22c55e" }
      case "apikey":
        return { icon: "🔑", label: "API Key Auth", color: "#f59e0b" }
      default:
        return { icon: "🆓", label: "Free Access", color: "#3b82f6" }
    }
  }

  const statusInfo = getStatusInfo()
  const authInfo = getAuthInfo()

  return (
    <div
      style={{
        height: "100%",
        width: "100%",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        background: "#1a1a1a",
        backgroundImage: `radial-gradient(ellipse 80% 80% at 50% -30%, ${primaryColorAlpha(0.3)}, rgba(255, 255, 255, 0))`,
        padding: "48px",
      }}
    >
      {/* Card Container */}
      <div
        style={{
          position: "relative",
          display: "flex",
          flexDirection: "column",
          width: "100%",
          maxWidth: "1000px",
          height: "480px",
          padding: "48px",
          borderRadius: "24px",
          backgroundColor: "rgba(255, 255, 255, 0.05)",
          border: `1px solid ${primaryColorAlpha(0.2)}`,
          boxShadow: `0 25px 50px -12px ${primaryColorAlpha(0.15)}`,
        }}
      >
        {/* Header section */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            marginBottom: "40px",
          }}
        >
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: "24px",
            }}
          >
            {/* Integration Icon */}
            <div
              style={{
                width: "64px",
                height: "64px",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                borderRadius: "12px",
                backgroundColor: "rgba(255, 255, 255, 0.1)",
              }}
            >
              {integrationIcon
                ? (
                    <img
                      src={integrationIcon}
                      alt={`${integrationName} logo`}
                      style={{
                        width: "48px",
                        height: "48px",
                        objectFit: "contain",
                      }}
                    />
                  )
                : (
                    <img
                      src="https://context.space/logo-color.svg"
                      alt="Context Space logo"
                      style={{
                        width: "48px",
                        height: "48px",
                        objectFit: "contain",
                      }}
                    />
                  )}
            </div>
            {/* Title and Categories */}
            <div style={{ display: "flex", flexDirection: "column", gap: "4px" }}>
              <h1
                style={{
                  fontSize: "40px",
                  fontWeight: "700",
                  color: "#ffffff",
                  margin: 0,
                  lineHeight: "1.1",
                }}
              >
                {title}
              </h1>
              {categories.length > 0 && (
                <div
                  style={{
                    display: "flex",
                    gap: "8px",
                    flexWrap: "wrap",
                  }}
                >
                  {categories.slice(0, 2).map((category, index) => (
                    <span
                      key={index}
                      style={{
                        fontSize: "12px",
                        color: "rgba(255, 255, 255, 0.6)",
                        backgroundColor: "rgba(255, 255, 255, 0.1)",
                        padding: "4px 8px",
                        borderRadius: "8px",
                        textTransform: "capitalize",
                      }}
                    >
                      {category}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Status badge */}
          <div
            style={{
              padding: "8px 12px",
              borderRadius: "8px",
              backgroundColor: addAlpha(statusInfo.color, 0.15),
              border: `1px solid ${addAlpha(statusInfo.color, 0.4)}`,
              color: statusInfo.color,
              fontSize: "16px",
              fontWeight: "600",
            }}
          >
            {statusInfo.text}
          </div>
        </div>

        {/* Main content area */}
        <div
          style={{
            flex: 1,
            display: "flex",
            flexDirection: "column",
            gap: "32px",
          }}
        >
          {/* Main description */}
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              gap: "16px",
            }}
          >
            <p
              style={{
                fontSize: "24px",
                lineHeight: "1.4",
                color: "#ffffff",
                margin: 0,
                fontWeight: "600",
              }}
            >
              Connect
              {" "}
              {integrationName}
              {" "}
              with Context Space
            </p>
            <p
              style={{
                fontSize: "20px",
                lineHeight: "1.5",
                color: "rgba(255, 255, 255, 0.7)",
                margin: 0,
              }}
            >
              {desc}
            </p>
          </div>

          {/* Integration features */}
          <div
            style={{
              display: "flex",
              gap: "40px",
              alignItems: "center",
            }}
          >
            <FeatureItem
              icon={authInfo.icon}
              label={authInfo.label}
              color={authInfo.color}
            />
            <FeatureItem
              icon="⚡"
              label="Quick Setup"
              color={primaryColor}
            />
            <FeatureItem
              icon="🔗"
              label="MCP Protocol"
              color="#22c55e"
            />
          </div>
        </div>

        {/* Bottom section */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            paddingTop: "24px",
            borderTop: "1px solid rgba(255, 255, 255, 0.1)",
          }}
        >
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: "12px",
              color: "rgba(255, 255, 255, 0.5)",
              fontSize: "18px",
              fontWeight: "500",
            }}
          >
            <span>🌐</span>
            <span>context.space</span>
          </div>
          <div
            style={{
              display: "flex",
              gap: "24px",
              alignItems: "center",
              color: "rgba(255, 255, 255, 0.4)",
              fontSize: "14px",
            }}
          >
            <span>Integration</span>
            <span>•</span>
            <span>MCP Protocol</span>
            <span>•</span>
            <span>Context Space</span>
          </div>
        </div>
      </div>
    </div>
  )
}

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url)

    // Check if this is for a specific integration via passed parameters
    const integrationName = searchParams.get("integration_name")
    const integrationIcon = searchParams.get("integration_icon")
    const integrationCategories = searchParams.get("integration_categories")
    const authType = searchParams.get("auth_type") as "oauth" | "apikey" | "none" | null
    const connectionStatus = searchParams.get("connection_status") as "connected" | "free" | "unconnected" | null

    // Extract title and description parameters
    const hasTitle = searchParams.has("title")
    const title = hasTitle
      ? searchParams.get("title")?.slice(0, 100)
      : integrationName
        ? `${integrationName} Integration`
        : siteName

    const hasDescription = searchParams.has("description")
    const ogDescription = hasDescription
      ? searchParams.get("description")?.slice(0, 200)
      : description

    // Integration-specific customization
    const isIntegrationPage = Boolean(integrationName)
    const categories = integrationCategories ? integrationCategories.split(",").map(c => c.trim()) : []

    const ogComponent = isIntegrationPage
      ? (
          <IntegrationOGImage
            title={title || ""}
            description={ogDescription || ""}
            integrationName={integrationName || ""}
            integrationIcon={integrationIcon}
            categories={categories}
            authType={authType}
            connectionStatus={connectionStatus}
          />
        )
      : (
          <DefaultOGImage
            title={title || ""}
            description={ogDescription || ""}
          />
        )

    return new ImageResponse(
      ogComponent,
      {
        width: 1200,
        height: 630,
      },
    )
  } catch (error) {
    console.error("OG image generation error:", error)
    return new Response(`Failed to generate the image`, {
      status: 500,
    })
  }
}
