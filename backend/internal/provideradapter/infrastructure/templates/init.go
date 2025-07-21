package templates

import (
	// Import all provider templates to register them via init()
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/airtable"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/eodhd"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/fetch"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/figma"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/github"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/hubspot"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/knowledgebase"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/mcp"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/notion"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/search"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/serper"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/slack"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/spotify"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/stripe"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/zoom"

	// Add more provider imports here as they are implemented
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/amap"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/openweathermap"
	_ "github.com/context-space/context-space/backend/internal/provideradapter/infrastructure/adapters/tmdb"
)

// Init function explicitly ensures this package is imported
func Init() {
	// This function doesn't need to do anything,
	// its purpose is just to ensure this package is imported
}
