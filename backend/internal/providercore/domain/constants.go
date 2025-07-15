package domain

// Context keys for provider domain
type preferredLanguageKeyType string

var (
	// PreferredLanguageKey is the context key for preferred language in provider requests
	// This key is used by middleware to store and handlers to retrieve the user's preferred language
	PreferredLanguageKey preferredLanguageKeyType = "provider.preferredLanguage"
)

// ProviderAPIDocURLs maps provider identifiers to their API documentation URLs
var ProviderAPIDocURLs = map[string]string{
	"airtable":          "https://airtable.com/developers/web/api/introduction",
	"eodhd":             "https://eodhd.com/financial-apis/quick-start-with-our-financial-data-apis",
	"fetch":             "",
	"github":            "https://docs.github.com/en/rest",
	"notion":            "https://developers.notion.com/reference/intro",
	"search":            "",
	"serper":            "https://serper.dev/api-keys",
	"slack":             "https://api.slack.com/web",
	"spotify":           "https://developer.spotify.com/documentation/web-api",
	"stripe":            "https://docs.stripe.com/keys",
	"zoom":              "https://developers.zoom.us/docs/api",
	"figma":             "https://www.figma.com/developers/api",
	"hubspot":           "https://developers.hubspot.com/docs/api/overview",
	"cfa_knowledgebase": "",
	"openweathermap":    "https://openweathermap.org/api",
}

// GetProviderAPIDocURL returns the API documentation URL for a given provider identifier
func GetProviderAPIDocURL(providerIdentifier string) string {
	if url, exists := ProviderAPIDocURLs[providerIdentifier]; exists {
		return url
	}
	return ""
}
