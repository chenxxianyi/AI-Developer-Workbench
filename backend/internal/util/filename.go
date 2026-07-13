package util

import (
	"regexp"
	"strings"
)

// Allowed generated filenames for Agent Config tool.
var allowedGeneratedFilenames = map[string]bool{
	"AGENTS.md":                true,
	"TASK_PLAN.md":             true,
	"CODING_RULES.md":          true,
	"FRONTEND_STYLE_GUIDE.md":  true,
	"BACKEND_ARCHITECTURE.md":  true,
	"README_AGENT_CONTEXT.md":  true,
	"UI_REVIEW_REPORT.md":      true,
	"PROJECT_DOCTOR_REPORT.md": true,
	"API_DOCUMENTATION.md":     true,
	"DB_SCHEMA_REVIEW.md":      true,
	"openapi.json":             true,
	"migration.sql":            true,
}

// SanitizeFilename removes dangerous characters from a filename.
func SanitizeFilename(name string) string {
	// Remove path separators and NUL.
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "\x00", "")
	// Remove CR/LF to prevent response-header injection.
	name = strings.ReplaceAll(name, "\r", "")
	name = strings.ReplaceAll(name, "\n", "")
	// Neutralize double quotes (used in Content-Disposition filename="...").
	name = strings.ReplaceAll(name, "\"", "'")
	// Remove remaining ASCII control characters (0x00-0x1F and 0x7F).
	name = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7F {
			return -1
		}
		return r
	}, name)

	// Remove leading/trailing dots and spaces.
	name = strings.Trim(name, ". ")
	name = strings.TrimSpace(name)

	// Replace multiple dots to prevent extension spoofing / traversal.
	re := regexp.MustCompile(`\.{2,}`)
	name = re.ReplaceAllString(name, "_")

	// Limit length.
	if len(name) > 255 {
		name = name[:255]
	}

	// Ensure non-empty.
	if name == "" {
		name = "unnamed"
	}

	return name
}

// IsAllowedGeneratedFilename checks if a filename is in the whitelist.
func IsAllowedGeneratedFilename(name string) bool {
	_, ok := allowedGeneratedFilenames[name]
	return ok
}

// SafeFilename returns a safe filename or generates one if invalid.
func SafeFilename(name string) string {
	safe := SanitizeFilename(name)
	if safe == "" || safe == "unnamed" {
		return "output.md"
	}
	return safe
}
