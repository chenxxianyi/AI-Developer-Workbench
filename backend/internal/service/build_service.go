package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// BuildService executes project builds (npm install + build) in isolated processes.
type BuildService struct {
	ws        *WorkspaceService
	timeout   time.Duration
	semaphore chan struct{} // concurrency control
}

// NewBuildService creates a build service.
func NewBuildService(ws *WorkspaceService, timeoutSec, maxConcurrent int) *BuildService {
	if timeoutSec <= 0 {
		timeoutSec = 600
	}
	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}
	return &BuildService{
		ws:        ws,
		timeout:   time.Duration(timeoutSec) * time.Second,
		semaphore: make(chan struct{}, maxConcurrent),
	}
}

// Build runs npm install && npm run build in the project workspace.
func (s *BuildService) Build(parent context.Context, projectID string) (string, error) {
	dir := s.ws.ProjectDir(projectID)

	ctx, cancel := context.WithTimeout(parent, s.timeout)
	defer cancel()

	// Waiting for a concurrency slot is part of the build timeout and remains
	// cancellable when the client disconnects.
	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }()
	case <-ctx.Done():
		return "", fmt.Errorf("等待构建资源超时: %w", ctx.Err())
	}

	// npm install
	installCmd := exec.CommandContext(ctx, "npm", "install", "--ignore-scripts", "--no-audit", "--no-fund")
	installCmd.Dir = dir
	installOut, err := installCmd.CombinedOutput()
	if err != nil {
		return string(installOut), fmt.Errorf("npm install failed: %w\n%s", err, installOut)
	}

	var manifest struct {
		Scripts         map[string]string `json:"scripts"`
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	manifestBytes, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return "", fmt.Errorf("read package.json before verification: %w", err)
	}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return "", fmt.Errorf("parse package.json before verification: %w", err)
	}
	if manifest.Scripts["test"] != "" {
		testCmd := exec.CommandContext(ctx, "npm", "run", "test")
		testCmd.Dir = dir
		testOut, err := testCmd.CombinedOutput()
		if err != nil {
			return string(testOut), fmt.Errorf("npm test failed: %w\n%s", err, testOut)
		}
	}
	if manifest.Dependencies["vue-tsc"] != "" || manifest.DevDependencies["vue-tsc"] != "" {
		typecheckCmd := exec.CommandContext(ctx, "npm", "exec", "--no", "--", "vue-tsc", "--noEmit")
		typecheckCmd.Dir = dir
		typecheckOut, err := typecheckCmd.CombinedOutput()
		if err != nil {
			return string(typecheckOut), fmt.Errorf("vue typecheck failed: %w\n%s", err, typecheckOut)
		}
	}

	// npm run build
	previewBase := "/api/projects/" + url.PathEscape(projectID) + "/preview/"
	buildCmd := exec.CommandContext(ctx, "npm", "run", "build", "--", "--base="+previewBase)
	buildCmd.Dir = dir
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		return string(buildOut), fmt.Errorf("npm build failed: %w\n%s", err, buildOut)
	}

	return string(buildOut), nil
}
