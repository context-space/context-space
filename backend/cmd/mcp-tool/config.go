package main

import (
	"fmt"
	"strings"
)

// ParseArgs parses shell-style arguments with proper quote handling
func ParseArgs(argsStr string) ([]string, error) {
	var args []string
	var current strings.Builder
	var inSingleQuote, inDoubleQuote bool
	var escaped bool

	runes := []rune(argsStr)
	for _, char := range runes {
		if escaped {
			// Handle escaped character
			switch char {
			case 'n':
				current.WriteRune('\n')
			case 't':
				current.WriteRune('\t')
			case 'r':
				current.WriteRune('\r')
			case '\\':
				current.WriteRune('\\')
			case '"':
				current.WriteRune('"')
			case '\'':
				current.WriteRune('\'')
			default:
				current.WriteRune(char)
			}
			escaped = false
			continue
		}

		switch char {
		case '\\':
			if !inSingleQuote {
				escaped = true
				continue
			}
			current.WriteRune(char)
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			} else {
				current.WriteRune(char)
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			} else {
				current.WriteRune(char)
			}
		case ' ', '\t', '\n', '\r':
			if !inSingleQuote && !inDoubleQuote {
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	// Check for unclosed quotes
	if inSingleQuote {
		return nil, fmt.Errorf("unclosed single quote in arguments")
	}
	if inDoubleQuote {
		return nil, fmt.Errorf("unclosed double quote in arguments")
	}
	if escaped {
		return nil, fmt.Errorf("trailing escape character in arguments")
	}

	// Add the last argument if it exists
	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args, nil
}

// ParseEnvVars parses environment variables from a comma-separated KEY=VALUE string
func ParseEnvVars(envStr string) (map[string]string, error) {
	envVars := make(map[string]string)

	if envStr == "" {
		return envVars, nil
	}

	// Split by comma, but handle escaped commas
	var parts []string
	var current strings.Builder
	var escaped bool

	for _, char := range envStr {
		if escaped {
			current.WriteRune(char)
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if char == ',' {
			if current.Len() > 0 {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
			}
		} else {
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Split KEY=VALUE
		eqIndex := strings.Index(part, "=")
		if eqIndex == -1 {
			return nil, fmt.Errorf("invalid environment variable format: %s (expected KEY=VALUE)", part)
		}

		key := strings.TrimSpace(part[:eqIndex])
		value := part[eqIndex+1:] // Don't trim value to preserve spaces

		if key == "" {
			return nil, fmt.Errorf("empty environment variable key in: %s", part)
		}

		envVars[key] = value
	}

	return envVars, nil
}
