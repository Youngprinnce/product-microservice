package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	validator := NewValidator()

	t.Run("SanitizeString method", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "trim whitespace",
				input:    "  hello world  ",
				expected: "hello world",
			},
			{
				name:     "escape HTML",
				input:    "<script>alert('xss')</script>",
				expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
			},
			{
				name:     "trim and escape",
				input:    "  <div>content</div>  ",
				expected: "&lt;div&gt;content&lt;/div&gt;",
			},
			{
				name:     "normal text unchanged",
				input:    "hello world",
				expected: "hello world",
			},
			{
				name:     "empty string",
				input:    "",
				expected: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := validator.SanitizeString(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("SanitizeURL method", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "valid http URL",
				input:    "http://example.com",
				expected: "http://example.com",
			},
			{
				name:     "valid https URL",
				input:    "https://example.com",
				expected: "https://example.com",
			},
			{
				name:     "trim whitespace from valid URL",
				input:    "  https://example.com  ",
				expected: "https://example.com",
			},
			{
				name:     "invalid URL without protocol",
				input:    "example.com",
				expected: "",
			},
			{
				name:     "invalid URL with wrong protocol",
				input:    "ftp://example.com",
				expected: "",
			},
			{
				name:     "empty string",
				input:    "",
				expected: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := validator.SanitizeURL(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trim whitespace",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "escape HTML",
			input:    "<script>alert('xss')</script>",
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			name:     "normal text",
			input:    "Normal product name",
			expected: "Normal product name",
		},
		{
			name:     "trim and escape combined",
			input:    "  <img src='x' onerror='alert(1)'>  ",
			expected: "&lt;img src=&#39;x&#39; onerror=&#39;alert(1)&#39;&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid HTTP URL",
			input:    "http://example.com/file.pdf",
			expected: "http://example.com/file.pdf",
		},
		{
			name:     "valid HTTPS URL",
			input:    "https://example.com/download/file.zip",
			expected: "https://example.com/download/file.zip",
		},
		{
			name:     "invalid URL",
			input:    "not-a-url",
			expected: "",
		},
		{
			name:     "URL with whitespace",
			input:    "  https://example.com/file.pdf  ",
			expected: "https://example.com/file.pdf",
		},
		{
			name:     "invalid protocol",
			input:    "ftp://example.com",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
