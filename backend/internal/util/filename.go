package util

import (
	"regexp"
	"strings"
)

// Allowed generated filenames for Agent Config tool.
var allowedGeneratedFilenames = map[string]bool{
	"AGENTS.md":              true,
	"TASK_PLAN.md":          true,
	"CODING_RULES.md":       true,
	"FRONTEND_STYLE_GUIDE.md": true,
	"BACKEND_ARCHITECTURE.md": true,
	"README_AGENT_CONTEXT.md": true,
	"UI_REVIEW_REPORT.md":   true,
	"PROJECT_DOCTOR_REPORT.md": true,
	"API_DOCUMENTATION.md":  true,
	"DB_SCHEMA_REVIEW.md":   true,
	"openapi.json":          true,
	"migration.sql":         true,
}

// SanitizeFilename removes dangerous characters from a filename.
func SanitizeFilename(name string) string {
	// Remove path separators.
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "\x00", "")

	// Remove leading/trailing dots and spaces.
	name = strings.Trim(name, ". ")
	name = strings.TrimSpace(name)

	// Replace multiple dots.
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
	if safe == "" {
		return "output.md"
	}
	return safe
}