package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// Context keys for provider domain
type preferredLanguageKeyType string

var (
	// PreferredLanguageKey is the context key for preferred language in provider requests
	// This key is used by middleware to store and handlers to retrieve the user's preferred language
	PreferredLanguageKey preferredLanguageKeyType = "provider.preferredLanguage"
)

// SupportedLanguages defines the languages supported by the Provider API
var SupportedLanguages = []language.Tag{
	language.English,            // en
	language.SimplifiedChinese,  // zh-CN
	language.TraditionalChinese, // zh-TW
}

// Language constants for easy reference
var (
	English            = language.English
	SimplifiedChinese  = language.SimplifiedChinese
	TraditionalChinese = language.TraditionalChinese
)

var matcher = language.NewMatcher(SupportedLanguages)

// I18nMiddleware creates a middleware that parses Accept-Language header and sets preferred language
// This middleware is specifically designed for Provider API endpoints
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.GetHeader("Accept-Language")
		var preferredLangTag language.Tag // Store the language.Tag object

		if accept != "" {
			tags, _, _ := language.ParseAcceptLanguage(accept)
			if len(tags) > 0 {
				preferredLangTag, _, _ = matcher.Match(tags...)
			} else {
				preferredLangTag = language.English // fallback to English if no valid tags parsed
			}
		} else {
			preferredLangTag = language.English // default to English if no header
		}

		// Store the language.Tag object in gin context and request context
		c.Set(string(PreferredLanguageKey), preferredLangTag)
		ctx := context.WithValue(c.Request.Context(), PreferredLanguageKey, preferredLangTag)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
