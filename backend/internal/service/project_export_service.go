package service

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ProjectExportService handles ZIP export of project workspaces.
type ProjectExportService struct {
	ws *WorkspaceService
}

// NewProjectExportService creates a project export service.
func NewProjectExportService(ws *WorkspaceService) *ProjectExportService {
	return &ProjectExportService{ws: ws}
}

// ExportExcludes lists patterns excluded from ZIP export (security-sensitive).
var ExportExcludes = []string{
	".env", ".env.local", ".env.production",
	"node_modules", ".git", ".npm-cache",
	"secrets", "*.pem", "*.key",
}

// ExportProject creates a ZIP archive of the project workspace.
func (s *ProjectExportService) ExportProject(projectID string, w io.Writer) error {
	dir := s.ws.ProjectDir(projectID)
	zw := zip.NewWriter(w)
	defer zw.Close()

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		for _, pattern := range ExportExcludes {
			if matched, _ := filepath.Match(pattern, filepath.Base(relPath)); matched {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.Contains(relPath, pattern) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open %s: %w", relPath, err)
		}
		defer f.Close()
		w, err := zw.Create(relPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		return err
	})
}
