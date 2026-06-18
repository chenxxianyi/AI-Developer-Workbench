package util

import "strings"

// NormalizeScore clamps a score to the valid 0-100 range.
func NormalizeScore(score int) int {
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

// NormalizeSeverity maps severity to the whitelist.
func NormalizeSeverity(s string) string {
	switch strings.ToLower(s) {
	case "high", "critical", "error":
		return "high"
	case "medium", "warning", "warn":
		return "medium"
	case "low", "info", "information":
		return "low"
	default:
		return "medium"
	}
}

// ComputeGrade returns a letter grade based on score.
func ComputeGrade(score int) string {
	if score >= 90 {
		return "A"
	}
	if score >= 80 {
		return "B"
	}
	if score >= 70 {
		return "C"
	}
	if score >= 60 {
		return "D"
	}
	return "F"
}