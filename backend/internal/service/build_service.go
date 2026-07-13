package service

import (
	"context"
	"fmt"
	"os/exec"
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
	return &BuildService{
		ws:        ws,
		timeout:   time.Duration(timeoutSec) * time.Second,
		semaphore: make(chan struct{}, maxConcurrent),
	}
}

// Build runs npm install && npm run build in the project workspace.
func (s *BuildService) Build(projectID string) (string, error) {
	dir := s.ws.ProjectDir(projectID)

	// Acquire concurrency slot
	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	// npm install
	installCmd := exec.CommandContext(ctx, "npm", "install")
	installCmd.Dir = dir
	installOut, err := installCmd.CombinedOutput()
	if err != nil {
		return string(installOut), fmt.Errorf("npm install failed: %w\n%s", err, installOut)
	}

	// npm run build
	buildCmd := exec.CommandContext(ctx, "npm", "run", "build")
	buildCmd.Dir = dir
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		return string(buildOut), fmt.Errorf("npm build failed: %w\n%s", err, buildOut)
	}

	return string(buildOut), nil
}
