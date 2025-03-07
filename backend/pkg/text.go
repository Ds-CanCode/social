package funcs

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

// IsValidInput sanitizes and validates input against common attacks
func IsValidInput(s string, max int) string {
	s = strings.TrimSpace(s)

	// Remove non-printable Unicode control characters
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, s)

	// Escape HTML to prevent XSS
	s = html.EscapeString(s)

	if len(s) > max || s == "" {
		return ""
	}

	// Remove potentially dangerous patterns (basic)
	s = sanitizeDangerousInput(s)

	return s
}

func sanitizeDangerousInput(s string) string {
	// Prevent <script> tags
	scriptRegex := regexp.MustCompile(`(?i)<\s*script.*?>.*?<\s*/\s*script\s*>`)
	s = scriptRegex.ReplaceAllString(s, "")

	// Prevent common SQL Injection patterns
	sqlInjectionRegex := regexp.MustCompile(`(?i)(union\s+select|select\s+\*|insert\s+into|update\s+set|drop\s+table|delete\s+from)`)
	s = sqlInjectionRegex.ReplaceAllString(s, "")

	// Prevent NULL byte attacks
	s = strings.ReplaceAll(s, "\x00", "")

	return s
}
