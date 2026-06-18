package util

import (
	"regexp"
	"strings"
)

// Patterns for secret detection.
var secretPatterns = []struct {
	pattern *regexp.Regexp
	name    string
}{
	{regexp.MustCompile(`(?i)(api[_-]?key|apikey|access[_-]?key|secret[_-]?key)\s*[=:]\s*['"]?[a-zA-Z0-9_\-]{20,}['"]?`), "api_key"},
	{regexp.MustCompile(`(?i)authorization\s*:\s*['"]?(bearer|token)\s+[a-zA-Z0-9_\-\.]{20,}['"]?`), "bearer_token"},
	{regexp.MustCompile(`(?i)sk-[a-zA-Z0-9]{20,}`), "openai_key"},
	{regexp.MustCompile(`(?i)github_pat_[a-zA-Z0-9_]{20,}`), "github_pat"},
	{regexp.MustCompile(`(?i)ghp_[a-zA-Z0-9]{20,}`), "github_token"},
	{regexp.MustCompile(`(?i)gho_[a-zA-Z0-9]{20,}`), "github_oauth"},
	{regexp.MustCompile(`(?i)ghs_[a-zA-Z0-9]{20,}`), "github_server"},
	{regexp.MustCompile(`(?i)ghr_[a-zA-Z0-9]{20,}`), "github_release"},
	{regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[=:]\s*['"]?[^\s'"<>]{8,}['"]?`), "password"},
	{regexp.MustCompile(`(?i)(mysql|postgres|postgresql|mongodb|redis)://[^:]+:[^@]+@`), "dsn"},
	{regexp.MustCompile(`(?i)(aws[_-]?access[_-]?key[_-]?id|AKIA[A-Z0-9]{16})`), "aws_key"},
	{regexp.MustCompile(`(?i)(aws[_-]?secret[_-]?access[_-]?key)[=:]\s*['"]?[a-zA-Z0-9+/]{40}['"]?`), "aws_secret"},
	{regexp.MustCompile(`(?i)private[_-]?key\s*[=:]\s*['"]?-----BEGIN`), "private_key"},
}

// Sensitive field names to redact in maps.
var sensitiveFieldNames = []string{
	"password", "passwd", "pwd", "secret", "token", "api_key", "apikey",
	"access_key", "secret_key", "auth", "authorization", "credential",
	"private_key", "privatekey", "dsn", "connection_string",
	"api_secret", "apisecret", "client_secret", "clientsecret",
}

// RedactText replaces secrets in text with [REDACTED].
func RedactText(text string) string {
	for _, sp := range secretPatterns {
		text = sp.pattern.ReplaceAllStringFunc(text, func(match string) string {
			// Preserve the key name but redact the value.
			return sp.name + "= [REDACTED]"
		})
	}
	return text
}

// RedactMap recursively redacts sensitive values in a map.
func RedactMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		lk := strings.ToLower(k)
		isSensitive := false
		for _, sf := range sensitiveFieldNames {
			if strings.Contains(lk, sf) {
				isSensitive = true
				break
			}
		}
		if isSensitive {
			result[k] = "[REDACTED]"
		} else {
			switch val := v.(type) {
			case map[string]interface{}:
				result[k] = RedactMap(val)
			case []interface{}:
				result[k] = redactSlice(val)
			default:
				result[k] = v
			}
		}
	}
	return result
}

// redactSlice redacts sensitive values in a slice.
func redactSlice(data []interface{}) []interface{} {
	result := make([]interface{}, len(data))
	for i, v := range data {
		switch val := v.(type) {
		case map[string]interface{}:
			result[i] = RedactMap(val)
		case []interface{}:
			result[i] = redactSlice(val)
		default:
			result[i] = v
		}
	}
	return result
}