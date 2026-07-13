package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WorkspaceService manages project workspace directories with path security.
type WorkspaceService struct {
	rootDir string
}

// NewWorkspaceService creates a workspace service rooted at the given directory.
func NewWorkspaceService(rootDir string) *WorkspaceService {
	return &WorkspaceService{rootDir: rootDir}
}

// RootDir returns the workspace root.
func (s *WorkspaceService) RootDir() string { return s.rootDir }

// ProjectDir returns the project-specific workspace path.
func (s *WorkspaceService) ProjectDir(projectID string) string {
	return filepath.Join(s.rootDir, projectID)
}

// EnsureProjectDir creates the project workspace directory if it doesn't exist.
func (s *WorkspaceService) EnsureProjectDir(projectID string) error {
	dir := s.ProjectDir(projectID)
	return os.MkdirAll(dir, 0o750)
}

// SafeResolve resolves a user-supplied path within the project workspace.
// Returns an error if path traversal or absolute path is detected.
func (s *WorkspaceService) SafeResolve(projectID, userPath string) (string, error) {
	root := s.ProjectDir(projectID)

	// 1. Clean the user path
	clean := filepath.Clean(userPath)

	// 2. Reject absolute paths
	if filepath.IsAbs(clean) {
		return "", fmt.Errorf("绝对路径不允许: %s", userPath)
	}

	// 3. Join and clean again
	full := filepath.Clean(filepath.Join(root, clean))

	// 4. Verify within root
	if !strings.HasPrefix(full, root) {
		return "", fmt.Errorf("路径穿越检测: %s", userPath)
	}

	return full, nil
}

// ListDir lists files and directories within the project workspace.
func (s *WorkspaceService) ListDir(projectID, subPath string) ([]os.DirEntry, error) {
	fullPath, err := s.SafeResolve(projectID, subPath)
	if err != nil {
		return nil, err
	}
	return os.ReadDir(fullPath)
}

// ReadFile reads a file within the project workspace (max 1MB).
func (s *WorkspaceService) ReadFile(projectID, filePath string) ([]byte, error) {
	fullPath, err := s.SafeResolve(projectID, filePath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %s", filePath)
	}

	const maxSize = 1 << 20 // 1 MB
	if info.Size() > maxSize {
		return nil, fmt.Errorf("文件过大 (最大 1MB)")
	}

	return os.ReadFile(fullPath)
}

// WriteFile writes data to a file within the project workspace.
func (s *WorkspaceService) WriteFile(projectID, filePath string, data []byte) error {
	fullPath, err := s.SafeResolve(projectID, filePath)
	if err != nil {
		return err
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}

	return os.WriteFile(fullPath, data, 0o640)
}

// DeleteProject removes the entire project workspace directory.
func (s *WorkspaceService) DeleteProject(projectID string) error {
	dir := s.ProjectDir(projectID)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(dir)
}

// IsTextFile checks if a file is likely a text file based on extension.
func IsTextFile(filename string) bool {
	textExts := map[string]bool{
		".go": true, ".vue": true, ".ts": true, ".js": true, ".jsx": true, ".tsx": true,
		".html": true, ".css": true, ".scss": true, ".json": true, ".yaml": true, ".yml": true,
		".md": true, ".txt": true, ".sql": true, ".xml": true, ".svg": true, ".env": true,
		".gitignore": true, ".dockerignore": true, ".editorconfig": true,
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return textExts[ext]
}
