package security

import (
	"fmt"
	"net/url"
	"strings"
)

// RedirectURLValidator provides secure validation for OAuth redirect URLs
type RedirectURLValidator struct {
	allowedDomains []string
	allowedSchemes []string
}

// NewRedirectURLValidator creates a new redirect URL validator with secure defaults
func NewRedirectURLValidator(allowedDomains []string, allowedSchemes []string) *RedirectURLValidator {
	return &RedirectURLValidator{
		allowedDomains: allowedDomains,
		allowedSchemes: allowedSchemes,
	}
}

// ValidateRedirectURL validates that a redirect URL is safe and allowed
func (v *RedirectURLValidator) ValidateRedirectURL(redirectURL string) error {
	if redirectURL == "" {
		return fmt.Errorf("redirect URL cannot be empty")
	}

	// Parse the URL
	parsedURL, err := url.Parse(redirectURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Validate scheme
	if err := v.validateScheme(parsedURL); err != nil {
		return err
	}

	// Validate domain
	if err := v.validateDomain(parsedURL); err != nil {
		return err
	}

	// Check for malicious patterns
	if err := v.checkMaliciousPatterns(parsedURL); err != nil {
		return err
	}

	return nil
}

// checkMaliciousPatterns checks for common malicious URL patterns
func (v *RedirectURLValidator) checkMaliciousPatterns(parsedURL *url.URL) error {
	// Check for javascript: or data: schemes
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme == "javascript" || scheme == "data" || scheme == "vbscript" {
		return fmt.Errorf("dangerous scheme detected: %s", scheme)
	}

	// Check for double slashes that could indicate protocol confusion
	if strings.Contains(parsedURL.Host, "//") || strings.Contains(parsedURL.Path, "//") {
		return fmt.Errorf("double slashes detected in URL")
	}

	// Check for suspicious characters
	suspiciousChars := []string{"<", ">", "\"", "'", "\\r", "\\n"}
	fullURL := parsedURL.String()
	for _, char := range suspiciousChars {
		if strings.Contains(fullURL, char) {
			return fmt.Errorf("suspicious character detected: %s", char)
		}
	}

	return nil
}

// validateScheme ensures the URL scheme is allowed
func (v *RedirectURLValidator) validateScheme(parsedURL *url.URL) error {
	scheme := strings.ToLower(parsedURL.Scheme)

	for _, allowedScheme := range v.allowedSchemes {
		if scheme == allowedScheme {
			return nil
		}
	}

	return fmt.Errorf("scheme '%s' is not allowed", scheme)
}

// validateDomain ensures the URL domain is in the allowlist
func (v *RedirectURLValidator) validateDomain(parsedURL *url.URL) error {
	host := strings.ToLower(parsedURL.Hostname())

	for _, allowedDomain := range v.allowedDomains {
		if host == allowedDomain {
			return nil
		}

		// Check for subdomain matches
		if strings.HasSuffix(host, "."+allowedDomain) {
			return nil
		}
	}

	return fmt.Errorf("domain '%s' is not in the allowed list", host)
}
