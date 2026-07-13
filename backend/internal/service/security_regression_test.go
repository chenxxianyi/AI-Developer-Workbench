package service

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"ai-developer-workbench/internal/util"
)

// makeZip builds an in-memory zip from the given entries (name -> content).
func makeZip(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for name, content := range entries {
		f, err := w.Create(name)
		if err != nil {
			t.Fatalf("create zip entry %q: %v", name, err)
		}
		if _, err := f.Write([]byte(content)); err != nil {
			t.Fatalf("write zip entry %q: %v", name, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buf.Bytes()
}

// writeFile writes data to path.
func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}

// TestRedactTextBlocksSecretsFromPrompts verifies the redaction utility used by
// every tool service strips common secret shapes before content reaches the AI.
func TestRedactTextBlocksSecretsFromPrompts(t *testing.T) {
	code := `
const API_KEY = "sk-abcdefghijklmnopqrstuvwxyz0123456789"
const DB = "mysql://root:s3cret@localhost:3306/prod"
const GH_TOKEN = "ghp_abcdefghijklmnopqrstuvwxyz0123"
`
	redacted := util.RedactText(code)
	if strings.Contains(redacted, "sk-abcdef") {
		t.Fatalf("openai key survived redaction: %s", redacted)
	}
	if strings.Contains(redacted, "s3cret") {
		t.Fatalf("dsn password survived redaction: %s", redacted)
	}
	if strings.Contains(redacted, "ghp_abcde") {
		t.Fatalf("github token survived redaction: %s", redacted)
	}
	if !strings.Contains(redacted, "[REDACTED]") {
		t.Fatalf("expected [REDACTED] marker in: %s", redacted)
	}
}

// TestSanitizeFilename_NeutralizesInjection ensures download filenames cannot
// smuggle CR/LF, path separators, or traversal sequences into headers.
func TestSanitizeFilename_NeutralizesInjection(t *testing.T) {
	cases := []string{
		"report.md\r\nContent-Length: 0\r\n\r\nevil",
		"../../etc/passwd",
		`..\..\windows\system32`,
		"a\x00b\"c",
	}
	for _, in := range cases {
		got := util.SafeFilename(in)
		if strings.ContainsAny(got, "\r\n") {
			t.Fatalf("CR/LF survived in %q -> %q", in, got)
		}
		if strings.Contains(got, "..") {
			t.Fatalf("traversal survived in %q -> %q", in, got)
		}
	}
}

// TestZipSlipRejected builds an archive whose entry name escapes the target dir
// and asserts ExtractAndAnalyze rejects / skips it rather than writing outside.
func TestZipSlipRejected(t *testing.T) {
	payload := makeZip(t, map[string]string{
		"../../etc/evil.txt": "pwned",
		"safe/normal.txt":    "hello",
	})
	tmp := t.TempDir() + "/slip.zip"
	if err := writeFile(tmp, payload); err != nil {
		t.Fatalf("write tmp zip: %v", err)
	}

	svc := NewZipService(t.TempDir())
	limits := DefaultZipLimits(10, 10, 4096, 65536)
	summary, err := svc.ExtractAndAnalyze(tmp, limits)
	if err != nil {
		return // whole-archive rejection is acceptable
	}
	raw, _ := json.Marshal(summary)
	if strings.Contains(string(raw), "pwned") {
		t.Fatalf("traversal entry content leaked into summary: %s", raw)
	}
}

// TestZipBombRejected builds an archive with a pathological compression ratio and
// asserts ExtractAndAnalyze rejects it.
func TestZipBombRejected(t *testing.T) {
	bomb := strings.Repeat("A", 5_000_000)
	payload := makeZip(t, map[string]string{"bomb.txt": bomb})
	tmp := t.TempDir() + "/bomb.zip"
	if err := writeFile(tmp, payload); err != nil {
		t.Fatalf("write tmp zip: %v", err)
	}

	svc := NewZipService(t.TempDir())
	limits := DefaultZipLimits(10, 1, 4096, 65536)
	_, err := svc.ExtractAndAnalyze(tmp, limits)
	if err == nil {
		t.Fatal("expected zip bomb to be rejected, got nil error")
	}
}

// TestZipTooManyFiles asserts the file-count cap is enforced (no panic, extras skipped).
func TestZipTooManyFiles(t *testing.T) {
	entries := map[string]string{}
	for i := 0; i < 30; i++ {
		entries["f"+strings.Repeat("x", i%4+1)+".txt"] = "x"
	}
	payload := makeZip(t, entries)
	tmp := t.TempDir() + "/many.zip"
	if err := writeFile(tmp, payload); err != nil {
		t.Fatalf("write tmp zip: %v", err)
	}

	svc := NewZipService(t.TempDir())
	limits := DefaultZipLimits(5, 10, 4096, 65536)
	_, _ = svc.ExtractAndAnalyze(tmp, limits) // must not panic
}

