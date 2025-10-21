package config

import (
	"testing"
)

func TestNormalizeContentType(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		expectedResult string
	}{
		{
			name:           "simple content type",
			contentType:    "text/plain",
			expectedResult: "text/plain",
		},
		{
			name:           "content type with charset parameter",
			contentType:    "text/yaml; charset=utf-8",
			expectedResult: "text/yaml",
		},
		{
			name:           "content type with multiple parameters",
			contentType:    "application/yaml; charset=utf-8; boundary=something",
			expectedResult: "application/yaml",
		},
		{
			name:           "content type with quoted parameters",
			contentType:    `text/x-yaml; charset="utf-8"`,
			expectedResult: "text/x-yaml",
		},
		{
			name:           "invalid content type",
			contentType:    "invalid/content/type; malformed",
			expectedResult: "invalid/content/type; malformed",
		},
		{
			name:           "empty content type",
			contentType:    "",
			expectedResult: "",
		},
		{
			name:           "content type with spaces",
			contentType:    "text/plain ; charset=utf-8",
			expectedResult: "text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeContentType(tt.contentType)
			if result != tt.expectedResult {
				t.Errorf("normalizeContentType(%q) = %q, want %q", tt.contentType, result, tt.expectedResult)
			}
		})
	}
}

func TestAllowedContentTypes(t *testing.T) {
	// Test that our normalization works with the allowed content types
	testCases := []string{
		"text/plain",
		"text/yaml",
		"text/x-yaml",
		"application/yaml",
		"application/x-yaml",
	}

	for _, contentType := range testCases {
		// Test without parameters
		if !allowedContentTypes.Contains(contentType) {
			t.Errorf("allowedContentTypes should contain %q", contentType)
		}

		// Test with charset parameter
		withCharset := contentType + "; charset=utf-8"
		normalized := normalizeContentType(withCharset)
		if !allowedContentTypes.Contains(normalized) {
			t.Errorf("normalized content type %q should be allowed (original: %q)", normalized, withCharset)
		}
	}

	// Test that disallowed content types are rejected
	disallowedTypes := []string{
		"application/json",
		"text/html",
		"application/xml",
	}

	for _, contentType := range disallowedTypes {
		if allowedContentTypes.Contains(contentType) {
			t.Errorf("allowedContentTypes should not contain %q", contentType)
		}

		// Test with parameters
		withCharset := contentType + "; charset=utf-8"
		normalized := normalizeContentType(withCharset)
		if allowedContentTypes.Contains(normalized) {
			t.Errorf("normalized content type %q should not be allowed (original: %q)", normalized, withCharset)
		}
	}
}
