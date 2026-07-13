package tools

import (
	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/service"
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"
	"testing"
)

func TestUIReviewSendsDesktopAndMobileScreenshotsWithViewports(t *testing.T) {
	ai := &fakeUIReviewAIService{result: &service.AIResult{JSONText: `{"scores":[],"issues":[],"recommendations":[],"codex_prompt":""}`}}
	svc := NewUIReviewService(ai, newFakeUIReviewReportService(), &fakeUIReviewFileService{}, &fakeUIReviewZipService{}, filepath.Join("tmp", "uploads"))
	_, err := svc.Run(context.Background(), UIReviewFormInput{Title: "responsive", ReviewMode: "screenshot", DesktopScreenshot: &multipart.FileHeader{Filename: "desktop.png"}, MobileScreenshot: &multipart.FileHeader{Filename: "mobile.png"}, DesktopViewport: "1366x768", MobileViewport: "375x812"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ai.lastRequest.ImagePaths) != 2 {
		t.Fatalf("image paths=%d", len(ai.lastRequest.ImagePaths))
	}
	if !strings.Contains(ai.lastRequest.UserPrompt, "desktop=1366x768; mobile=375x812") {
		t.Fatal("viewport mapping missing")
	}
}
func TestUIReviewClampsAnnotationRegions(t *testing.T) {
	result := dto.UIReviewResult{Issues: []dto.IssueItem{{Region: &dto.IssueRegion{X: -5, Width: 110}}}}
	(&UIReviewService{}).normalizeResult(&result)
	if result.Issues[0].Region.X != 0 || result.Issues[0].Region.Width != 100 {
		t.Fatalf("region not clamped: %+v", result.Issues[0].Region)
	}
}

func TestUIReviewAddsPromptForHighSeverityIssue(t *testing.T) {
	result := dto.UIReviewResult{Issues: []dto.IssueItem{{Title: "navigation overflow", Category: "responsive", Severity: "high"}}}
	(&UIReviewService{}).normalizeResult(&result)
	if result.Issues[0].ComponentPrompt == "" {
		t.Fatal("high severity issue must have component prompt")
	}
}
