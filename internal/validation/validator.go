package validation

import (
	"html"
	"strings"
)

// Validator provides validation functionality
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// SanitizeString sanitizes string input using the validator instance
func (v *Validator) SanitizeString(input string) string {
	return SanitizeString(input)
}

// SanitizeURL sanitizes URL input using the validator instance
func (v *Validator) SanitizeURL(input string) string {
	return SanitizeURL(input)
}

// SanitizeString sanitizes string input
func SanitizeString(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// HTML escape to prevent XSS
	input = html.EscapeString(input)

	return input
}

// SanitizeURL sanitizes URL input
func SanitizeURL(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Basic validation - just check if it starts with http/https
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		return ""
	}

	return input
}
