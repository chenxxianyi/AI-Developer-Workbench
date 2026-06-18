package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ai-developer-workbench/internal/config"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/repository"

	"github.com/google/uuid"
)

// AllowedImageTypes returns allowed image MIME types for screenshots.
func AllowedImageTypes() map[string]string {
	return map[string]string{
		"image/png":  ".png",
		"image/jpeg": ".jpg",
		"image/webp": ".webp",
	}
}

// AllowedArchiveTypes returns allowed archive MIME types for project ZIPs.
func AllowedArchiveTypes() map[string]string {
	return map[string]string{
		"application/zip":              ".zip",
		"application/x-zip-compressed": ".zip",
	}
}

// FileService handles file upload and validation.
type FileService interface {
	SaveUpload(ctx context.Context, reportID string, fileHeader *multipart.FileHeader, assetType string, allowedTypes map[string]string) (*model.ReportAsset, error)
	DeleteReportDir(uploadDir, tempDir, reportID string) error
	ValidateFile(fileHeader *multipart.FileHeader, allowedTypes map[string]string) error
}

type fileService struct {
	cfg       *config.Config
	assetRepo repository.ReportAssetRepository
}

// NewFileService creates a new file service.
func NewFileService(cfg *config.Config, assetRepo repository.ReportAssetRepository) FileService {
	return &fileService{cfg: cfg, assetRepo: assetRepo}
}

// ValidateFile checks file extension, Content-Type, and magic bytes.
func (s *fileService) ValidateFile(fileHeader *multipart.FileHeader, allowedTypes map[string]string) error {
	// 1. Check extension.
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	validExt := false
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			validExt = true
			break
		}
	}
	if !validExt {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	// 2. Check Content-Type from header.
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "" {
		validMIME := false
		for allowedMIME := range allowedTypes {
			if strings.EqualFold(contentType, allowedMIME) {
				validMIME = true
				break
			}
		}
		if !validMIME {
			return fmt.Errorf("unsupported content type: %s", contentType)
		}
	}

	// 3. Check magic bytes (file header).
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file for validation: %w", err)
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %w", err)
	}
	detectedType := http.DetectContentType(buf[:n])
	validDetected := false
	for allowedMIME := range allowedTypes {
		if strings.EqualFold(detectedType, allowedMIME) {
			validDetected = true
			break
		}
	}
	if !validDetected {
		// Some ZIP files are detected as application/octet-stream.
		if detectedType == "application/octet-stream" {
			for allowedMIME := range allowedTypes {
				if strings.EqualFold(allowedMIME, "application/zip") || strings.EqualFold(allowedMIME, "application/x-zip-compressed") {
					validDetected = true
					break
				}
			}
		}
		if !validDetected {
			return fmt.Errorf("file content does not match allowed types (detected: %s)", detectedType)
		}
	}

	// 4. Check file size.
	maxBytes := int64(s.cfg.Upload.MaxUploadSizeMB) * 1024 * 1024
	if fileHeader.Size > maxBytes {
		return fmt.Errorf("file too large: %d bytes (max %d)", fileHeader.Size, maxBytes)
	}

	return nil
}

// SaveUpload saves an uploaded file and creates a ReportAsset record.
func (s *fileService) SaveUpload(ctx context.Context, reportID string, fileHeader *multipart.FileHeader, assetType string, allowedTypes map[string]string) (*model.ReportAsset, error) {
	// Validate file.
	if err := s.ValidateFile(fileHeader, allowedTypes); err != nil {
		return nil, err
	}

	// Determine file extension.
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		contentType := fileHeader.Header.Get("Content-Type")
		if e, ok := allowedTypes[contentType]; ok {
			ext = e
		}
	}

	// Generate stored filename.
	storedName := uuid.New().String() + ext

	// Create report directory.
	reportDir := filepath.Join(s.cfg.Upload.Dir, reportID, "source")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Verify path is within configured uploads dir.
	destPath := filepath.Join(reportDir, storedName)
	if err := validatePathWithinDir(destPath, s.cfg.Upload.Dir); err != nil {
		return nil, fmt.Errorf("path validation failed: %w", err)
	}

	// Open source file.
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file.
	dst, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy and compute SHA-256.
	hasher := sha256.New()
	size, err := io.Copy(dst, io.TeeReader(src, hasher))
	if err != nil {
		os.Remove(destPath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	sha256Sum := fmt.Sprintf("%x", hasher.Sum(nil))

	// Build relative path.
	relPath, err := filepath.Rel(s.cfg.Upload.Dir, destPath)
	if err != nil {
		relPath = filepath.Join(reportID, "source", storedName)
	}

	// Detect MIME type.
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Create asset record.
	asset := &model.ReportAsset{
		ReportID:     reportID,
		AssetType:    assetType,
		OriginalName: sanitizeOriginalName(fileHeader.Filename),
		StoredName:   storedName,
		RelativePath: filepath.ToSlash(relPath),
		MimeType:     mimeType,
		SizeBytes:    uint64(size),
		SHA256:       sha256Sum,
	}

	if err := s.assetRepo.Create(ctx, asset); err != nil {
		os.Remove(destPath)
		return nil, fmt.Errorf("failed to save asset record: %w", err)
	}

	return asset, nil
}

// DeleteReportDir removes uploaded and temp directories for a report.
func (s *fileService) DeleteReportDir(uploadDir, tempDir, reportID string) error {
	dirs := []string{
		filepath.Join(uploadDir, reportID),
		filepath.Join(tempDir, reportID),
	}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			slog.Warn("Failed to remove report directory", "dir", dir, "error", err)
		}
	}
	return nil
}

// validatePathWithinDir ensures a path is within the allowed directory.
func validatePathWithinDir(target, baseDir string) error {
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return err
	}
	if strings.HasPrefix(rel, "..") {
		return fmt.Errorf("path escapes base directory")
	}
	return nil
}

// sanitizeOriginalName removes path separators and dangerous characters from original filename.
func sanitizeOriginalName(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, "\x00", "")
	if name == "." || name == ".." || name == "" {
		return "upload"
	}
	return name
}
