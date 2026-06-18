package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// TruncateText truncates text to maxBytes, adding ellipsis marker.
func TruncateText(text string, maxBytes int) string {
	if len(text) <= maxBytes {
		return text
	}

	// Truncate at byte boundary.
	truncated := text[:maxBytes]
	// Remove partial UTF-8 characters.
	for i := len(truncated) - 1; i >= 0 && truncated[i] >= 0x80; i-- {
		truncated = truncated[:i]
	}

	return truncated + "\n[...truncated...]"
}

// EncodeBase64 encodes data as base64 string.
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes a base64 string.
func DecodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// SafeJSONString ensures a string is valid JSON by escaping.
func SafeJSONString(s string) string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return `""`
	}
	return string(bytes)
}

// NormalizeURL ensures a URL has proper scheme.
func NormalizeURL(rawURL string) string {
	if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		return rawURL
	}
	return "https://" + rawURL
}

// ValidateURL checks if a URL is valid.
func ValidateURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if parsed.Scheme == "" {
		return fmt.Errorf("URL missing scheme")
	}
	if parsed.Host == "" {
		return fmt.Errorf("URL missing host")
	}
	return nil
}

// ContainsOnly checks if string contains only allowed characters.
func ContainsOnly(s string, allowed string) bool {
	for _, c := range s {
		if !strings.ContainsRune(allowed, c) {
			return false
		}
	}
	return true
}