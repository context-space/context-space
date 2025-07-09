package middleware

import (
	"context"

	"github.com/context-space/context-space/backend/internal/providercore/domain"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
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
		c.Set(string(domain.PreferredLanguageKey), preferredLangTag)
		ctx := context.WithValue(c.Request.Context(), domain.PreferredLanguageKey, preferredLangTag)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
