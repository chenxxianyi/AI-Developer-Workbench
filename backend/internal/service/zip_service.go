package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"ai-developer-workbench/internal/util"
)

// Directories to ignore during project scanning.
var ignoreDirs = map[string]bool{
	"node_modules": true, ".git": true, "dist": true, "build": true,
	"coverage": true, "vendor": true, ".idea": true, ".vscode": true,
	"__pycache__": true, ".next": true, ".nuxt": true, "target": true,
	"bin": true, "obj": true, ".cache": true,
}

// Files to never read (security).
var sensitiveFilePatterns = []string{
	".env", ".env.local", ".env.production", ".env.development",
	".pem", ".key", ".p12", ".pfx", ".jks",
	"id_rsa", "id_ed25519", "id_ecdsa",
	".credentials", ".service-account",
}

// Priority files to read first.
var priorityFiles = []string{
	"README.md", "package.json", "go.mod", "requirements.txt",
	"pyproject.toml", "Dockerfile", "docker-compose.yml",
	".env.example", "AGENTS.md", "CLAUDE.md",
}

// Priority directory patterns.
var priorityDirPatterns = []string{
	"src/router", "src/api", "internal", "cmd", "app",
	"api", "routes", "handlers",
}

// ZipLimits holds limits for ZIP extraction.
type ZipLimits struct {
	MaxFiles            int
	MaxUncompressedMB   int
	MaxFileReadBytes    int
	MaxTotalReadBytes   int
	MaxCompressionRatio float64
}

// DefaultZipLimits returns default limits from config.
func DefaultZipLimits(maxFiles, maxUncompressedMB, maxFileReadBytes, maxTotalReadBytes int) ZipLimits {
	return ZipLimits{
		MaxFiles:            maxFiles,
		MaxUncompressedMB:   maxUncompressedMB,
		MaxFileReadBytes:    maxFileReadBytes,
		MaxTotalReadBytes:   maxTotalReadBytes,
		MaxCompressionRatio: 100.0, // Max 100:1 compression ratio
	}
}

// FileSummary holds content of a scanned file.
type FileSummary struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int    `json:"size"`
	Truncated bool `json:"truncated"`
}

// ProjectSummary holds the result of project scanning.
type ProjectSummary struct {
	Tree              []string               `json:"tree"`
	DetectedStack     []string               `json:"detected_stack"`
	ImportantFiles    []FileSummary          `json:"important_files"`
	EngineeringSignals map[string]interface{} `json:"engineering_signals"`
	Truncation        TruncationInfo         `json:"truncation"`
}

// TruncationInfo tracks how much data was omitted.
type TruncationInfo struct {
	FilesSkipped int   `json:"files_skipped"`
	BytesOmitted int64 `json:"bytes_omitted"`
}

// ZipService handles ZIP extraction and project scanning.
type ZipService interface {
	ExtractAndAnalyze(zipPath string, limits ZipLimits) (*ProjectSummary, error)
	Cleanup(dir string) error
}

type zipService struct {
	tempDir string
}

// NewZipService creates a new ZIP service.
func NewZipService(tempDir string) ZipService {
	return &zipService{tempDir: tempDir}
}

// ExtractAndAnalyze extracts a ZIP file and performs static analysis.
func (s *zipService) ExtractAndAnalyze(zipPath string, limits ZipLimits) (*ProjectSummary, error) {
	// Open the ZIP file.
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ZIP: %w", err)
	}
	defer reader.Close()

	summary := &ProjectSummary{
		Tree:               make([]string, 0),
		DetectedStack:      make([]string, 0),
		ImportantFiles:     make([]FileSummary, 0),
		EngineeringSignals: make(map[string]interface{}),
	}

	var totalUncompressed int64
	var totalReadBytes int64
	var fileCount int
	var filesSkipped int
	var bytesOmitted int64

	// Collect and sort entries for priority reading.
	entries := make([]*zip.File, 0, len(reader.File))
	for _, f := range reader.File {
		entries = append(entries, f)
	}
	sort.Slice(entries, func(i, j int) bool {
		return priorityScore(entries[i].Name) > priorityScore(entries[j].Name)
	})

	for _, f := range entries {
		// Security: check for Zip Slip.
		cleanName := filepath.Clean(f.Name)
		if strings.HasPrefix(cleanName, "..") || strings.HasPrefix(cleanName, string(os.PathSeparator)) {
			slog.Warn("ZIP entry with path traversal rejected", "entry", f.Name)
			continue
		}
		// Windows drive letter check.
		if len(cleanName) >= 2 && cleanName[1] == ':' {
			slog.Warn("ZIP entry with Windows drive letter rejected", "entry", f.Name)
			continue
		}

		// Skip directories.
		if f.FileInfo().IsDir() {
			continue
		}

		// Security: reject symlinks and non-regular files.
		if f.Mode()&os.ModeSymlink != 0 {
			slog.Warn("ZIP symlink entry rejected", "entry", f.Name)
			continue
		}

		// File count limit.
		fileCount++
		if fileCount > limits.MaxFiles {
			filesSkipped++
			continue
		}

		// ZIP bomb: check compression ratio.
		if f.CompressedSize64 > 0 {
			ratio := float64(f.UncompressedSize64) / float64(f.CompressedSize64)
			if ratio > limits.MaxCompressionRatio {
				slog.Warn("ZIP entry with suspicious compression ratio", "entry", f.Name, "ratio", ratio)
				return nil, fmt.Errorf("unsafe archive: suspicious compression ratio for %s", f.Name)
			}
		}

		// Track total uncompressed size.
		totalUncompressed += int64(f.UncompressedSize64)
		maxUncompressed := int64(limits.MaxUncompressedMB) * 1024 * 1024
		if totalUncompressed > maxUncompressed {
			return nil, fmt.Errorf("unsafe archive: uncompressed size exceeds %d MB", limits.MaxUncompressedMB)
		}

		// Check if in ignored directory.
		if isInIgnoredDir(cleanName) {
			continue
		}

		// Check if sensitive file.
		if isSensitiveFile(cleanName) {
			continue
		}

		// Add to tree.
		summary.Tree = append(summary.Tree, filepath.ToSlash(cleanName))

		// Detect stack from filenames.
		detectStack(cleanName, summary)

		// Read file content if text file and within budget.
		if totalReadBytes >= int64(limits.MaxTotalReadBytes) {
			filesSkipped++
			continue
		}

		content, truncated, readSize, err := readFileFromZip(f, limits.MaxFileReadBytes)
		if err != nil {
			slog.Warn("Failed to read ZIP entry", "entry", f.Name, "error", err)
			continue
		}

		// Skip binary files.
		if isBinaryContent(content) {
			continue
		}

		// Redact secrets.
		content = util.RedactText(content)

		totalReadBytes += int64(readSize)
		if truncated {
			bytesOmitted += int64(readSize) - int64(limits.MaxFileReadBytes)
		}

		summary.ImportantFiles = append(summary.ImportantFiles, FileSummary{
			Path:      filepath.ToSlash(cleanName),
			Content:   content,
			Size:      readSize,
			Truncated: truncated,
		})
	}

	summary.Truncation = TruncationInfo{
		FilesSkipped: filesSkipped,
		BytesOmitted: bytesOmitted,
	}

	// Deduplicate detected stack.
	summary.DetectedStack = uniqueStrings(summary.DetectedStack)

	return summary, nil
}

// Cleanup removes a temporary directory.
func (s *zipService) Cleanup(dir string) error {
	return os.RemoveAll(dir)
}

// readFileFromZip reads a file from a ZIP entry with size limits.
func readFileFromZip(f *zip.File, maxBytes int) (string, bool, int, error) {
	rc, err := f.Open()
	if err != nil {
		return "", false, 0, err
	}
	defer rc.Close()

	// Limit read size.
	limitedReader := io.LimitReader(rc, int64(maxBytes)+1)
	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, limitedReader)
	if err != nil {
		return "", false, 0, err
	}

	content := buf.String()
	truncated := n > int64(maxBytes)
	if truncated {
		content = content[:maxBytes]
	}

	return content, truncated, int(n), nil
}

// isBinaryContent checks if content appears to be binary.
func isBinaryContent(content string) bool {
	if len(content) == 0 {
		return false
	}
	// Check first 512 bytes.
	checkLen := 512
	if len(content) < checkLen {
		checkLen = len(content)
	}
	detected := http.DetectContentType([]byte(content[:checkLen]))
	return !strings.HasPrefix(detected, "text/") && detected != "application/json" && detected != "application/xml"
}

// isInIgnoredDir checks if a path is within an ignored directory.
func isInIgnoredDir(path string) bool {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for _, part := range parts {
		if ignoreDirs[part] {
			return true
		}
	}
	return false
}

// isSensitiveFile checks if a file should never be read.
func isSensitiveFile(path string) bool {
	base := filepath.Base(path)
	for _, pattern := range sensitiveFilePatterns {
		if base == pattern || strings.HasSuffix(base, pattern) {
			return true
		}
	}
	// Check for private key patterns.
	if strings.Contains(base, "id_rsa") || strings.Contains(base, "id_ed25519") || strings.Contains(base, "id_ecdsa") {
		return true
	}
	return false
}

// priorityScore determines reading priority for a file.
func priorityScore(path string) int {
	base := filepath.Base(path)
	for i, pf := range priorityFiles {
		if base == pf {
			return 100 - i
		}
	}
	dir := filepath.Dir(path)
	for i, pd := range priorityDirPatterns {
		if strings.HasPrefix(filepath.ToSlash(dir), pd) {
			return 50 - i
		}
	}
	return 0
}

// detectStack adds detected technology stack from filename.
func detectStack(path string, summary *ProjectSummary) {
	base := filepath.Base(path)
	switch base {
	case "package.json":
		summary.DetectedStack = append(summary.DetectedStack, "node.js")
	case "go.mod":
		summary.DetectedStack = append(summary.DetectedStack, "go")
	case "requirements.txt", "pyproject.toml", "Pipfile":
		summary.DetectedStack = append(summary.DetectedStack, "python")
	case "Cargo.toml":
		summary.DetectedStack = append(summary.DetectedStack, "rust")
	case "pom.xml", "build.gradle":
		summary.DetectedStack = append(summary.DetectedStack, "java")
	case "Dockerfile":
		summary.DetectedStack = append(summary.DetectedStack, "docker")
	case "docker-compose.yml", "docker-compose.yaml":
		summary.DetectedStack = append(summary.DetectedStack, "docker-compose")
	}
}

// uniqueStrings deduplicates a string slice.
func uniqueStrings(s []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(s))
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Ensure RedactText is available via the util package.
var _ = regexp.MustCompile // ensure regexp import available if needed
