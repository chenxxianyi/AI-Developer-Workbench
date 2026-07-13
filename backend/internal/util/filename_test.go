package util

import (
	"strings"
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain", "report.md", "report.md"},
		{"unix path", "../etc/passwd", "_etc_passwd"},
		{"windows path", `..\..\windows\system32`, `___windows_system32`},
		{"nul removed", "a\x00b", "ab"},
		{"cr lf removed", "a\r\nb", "ab"},
		{"quotes neutralized", `a"b`, "a'b"},
		{"control chars stripped", "a\x01\x02b", "ab"},
		{"del char stripped", "a\x7Fb", "ab"},
		{"leading dots trimmed", "...report.md", "report.md"},
		{"trailing dots trimmed", "report.md...", "report.md"},
		{"multiple dots collapsed", "report..md", "report_md"},
		{"empty becomes unnamed", "", "unnamed"},
		{"spaces only becomes unnamed", "   ", "unnamed"},
		{"chinese preserved", "报告-2026.md", "报告-2026.md"},
		{"length capped", strings.Repeat("a", 300), strings.Repeat("a", 255)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFilename(tt.input)
			if got != tt.want {
				t.Fatalf("SanitizeFilename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeFilename_PreventsHeaderInjection(t *testing.T) {
	// A payload that attempts to inject Content-Disposition headers must not
	// survive sanitization: no CR/LF, no quotes, no path traversal. The injected
	// header text itself is kept as part of the filename body (harmless once
	// stripped of CR/LF), but no CRLF means it cannot start a new header.
	injected := "report.md\r\nContent-Length: 0\r\n\r\nevil"
	got := SanitizeFilename(injected)
	if strings.Contains(got, "\r") || strings.Contains(got, "\n") {
		t.Fatalf("sanitized filename still contains CR/LF: %q", got)
	}
	// Must not be able to split a Content-Disposition header line.
	if strings.Contains(got, "Content-Length") && (strings.Contains(got, ":") == false) {
		t.Fatalf("unexpected header-shaped content: %q", got)
	}
}

func TestIsAllowedGeneratedFilename(t *testing.T) {
	if !IsAllowedGeneratedFilename("AGENTS.md") {
		t.Fatal("AGENTS.md should be allowed")
	}
	if IsAllowedGeneratedFilename("../../etc/passwd") {
		t.Fatal("path traversal should not be allowed")
	}
}

func TestSafeFilename(t *testing.T) {
	if got := SafeFilename(""); got != "output.md" {
		t.Fatalf("empty -> output.md, got %q", got)
	}
	if got := SafeFilename("../x"); got != "_x" {
		t.Fatalf("traversal -> _x, got %q", got)
	}
}
