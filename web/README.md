# Context Space Web Frontend

A modern, production-ready Next.js application that serves as the frontend interface for Context Space - the ultimate context engineering infrastructure platform.

## Tech Stack

### Core Framework
- **Next.js 15** - React framework with App Router and Turbopack
- **React 19** - Latest React with concurrent features
- **TypeScript 5.8** - Type-safe development experience

### Styling & UI
- **Tailwind CSS 4.1** - Utility-first CSS framework
- **Radix UI** - Headless, accessible component primitives
- **Shadcn/ui** - Beautiful, customizable component library
- **Lucide React** - Consistent iconography

### Authentication & Data
- **Supabase** - Authentication and database integration
- **TanStack Query** - Server state management
- **React Hook Form** - Form state management with validation
- **Zod** - Runtime type validation

### Content & Internationalization
- **Contentlayer** - Type-safe Markdown/MDX processing
- **Next-intl** - Internationalization (en, zh, zh-TW)
- **MDX** - Enhanced Markdown with React components

### Developer Experience
- **ESLint** - Code linting with @antfu/eslint-config
- **Commitlint** - Conventional commit message enforcement
- **Vitest** - Fast unit testing framework
- **Simple Git Hooks** - Git workflow automation

## Project Structure

```
src/
├── app/                    # Next.js App Router
│   ├── [locale]/          # Internationalized routes
│   │   ├── integration/   # OAuth integration flows
│   │   ├── integrations/  # Provider catalog
│   │   └── account/       # User dashboard
│   └── api/               # API routes (chat, MCP, search)
├── components/            # Reusable UI components
│   ├── ui/               # Shadcn/ui base components
│   ├── auth/             # Authentication components
│   ├── integration/      # Provider integration UI
│   └── landing/          # Marketing page components
├── lib/                  # Utility libraries
│   ├── supabase/        # Database client setup
│   ├── utils/           # Helper functions
│   └── mcp/             # MCP protocol integration
├── hooks/               # Custom React hooks
├── services/           # API service layer
└── typings/           # TypeScript type definitions
```

## Key Features

### 🔐 OAuth Integration Management
- Streamlined OAuth flows for 14+ providers (GitHub, Slack, Notion, etc.)
- Secure credential management with HashiCorp Vault integration
- Visual connection status and health monitoring

### 🚀 MCP Provider Playground
- Interactive testing interface for MCP operations
- Real-time operation execution and response inspection
- Parameter validation and documentation

### 🌐 Internationalization
- Full i18n support (English, Chinese Simplified, Chinese Traditional)
- Localized content and UI strings
- Dynamic locale switching

### 📱 Responsive Design
- Mobile-first responsive layout
- Dark/light theme support with `next-themes`
- Accessible components following WCAG guidelines

### 🔍 Content Management
- Blog system with MDX support
- Technical documentation with syntax highlighting
- SEO-optimized with Open Graph and Twitter cards

## Quick Start

### Prerequisites
- Node.js 20+ (defined in `.nvmrc`)
- pnpm 10.12.1+ (specified in `packageManager` field)

### Installation

```bash
# Clone the repository
git clone https://github.com/context-space/context-space.git
cd context-space/web

# Install dependencies
pnpm install

# install git hooks
pnpm update:hooks

# Start development server
pnpm dev
```

The application will be available at `http://localhost:4321`.

### Environment Variables

Create a `.env` file with required configuration:

```bash
# Site Configuration
NEXT_PUBLIC_SITE_NAME="Context Space"
NEXT_PUBLIC_BASE_URL="http://localhost:4321"
NEXT_PUBLIC_BASE_API_URL="http://localhost:8080"

# Supabase Configuration
NEXT_PUBLIC_SUPABASE_URL="your-supabase-url"
NEXT_PUBLIC_SUPABASE_ANON_KEY="your-anon-key"

# Optional Analytics
NEXT_PUBLIC_GOOGLE_ANALYTICS_ID="GA_MEASUREMENT_ID"
NEXT_PUBLIC_MICROSOFT_CLARITY="clarity-project-id"
GOOGLE_SITE_VERIFICATION="google-verification-code"

# Security
NEXT_PUBLIC_TURNSTILE_SITE_KEY="cloudflare-turnstile-key"
NEXT_PUBLIC_TURNSTILE_ENABLED="true"
```

## License

This project is licensed under the AGPL v3 License - see the [LICENSE](../LICENSE) file for details.