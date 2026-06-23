package prompts

import (
	"strings"
	"testing"
)

func TestBuildUIReviewPromptRequiresChineseOutput(t *testing.T) {
	systemPrompt, userPrompt := BuildUIReviewPrompt("code", "paste", "登录页", "简洁现代", "检查可用性", "<button>提交</button>", "")

	if !strings.Contains(systemPrompt, "必须使用简体中文") {
		t.Fatalf("system prompt should require Simplified Chinese output, got: %s", systemPrompt)
	}
	if !strings.Contains(systemPrompt, "视觉层级") {
		t.Fatalf("system prompt should use Chinese scoring dimensions, got: %s", systemPrompt)
	}
	if !strings.Contains(userPrompt, "请分析该 UI") {
		t.Fatalf("user prompt should request Chinese UI analysis, got: %s", userPrompt)
	}
}

func TestBuildUIReviewPromptIncludesProjectZipSummary(t *testing.T) {
	_, userPrompt := BuildUIReviewPrompt("code", "project_zip", "Dashboard", "简洁现代", "", "", `{"important_files":[{"path":"src/App.vue"}]}`)

	if !strings.Contains(userPrompt, "前端项目 ZIP 源码摘要") {
		t.Fatalf("user prompt should include project ZIP summary section, got: %s", userPrompt)
	}
	if !strings.Contains(userPrompt, "src/App.vue") {
		t.Fatalf("user prompt should include project ZIP summary content, got: %s", userPrompt)
	}
}
