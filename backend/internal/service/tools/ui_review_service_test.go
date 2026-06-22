package tools

import (
	"path/filepath"
	"testing"
)

func TestResolveUploadPathIncludesConfiguredUploadDirectory(t *testing.T) {
	uploadDir := filepath.Join("var", "ai-workbench", "uploads")
	relativePath := filepath.ToSlash(filepath.Join(
		"report-id",
		"source",
		"screenshot.png",
	))

	got := resolveUploadPath(uploadDir, relativePath)
	want := filepath.Join(uploadDir, "report-id", "source", "screenshot.png")

	if got != want {
		t.Fatalf("resolveUploadPath() = %q, want %q", got, want)
	}
}
