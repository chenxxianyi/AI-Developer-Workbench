package prompts

import (
	"strings"
	"testing"
)

func TestBuildUIReviewPromptRequiresChineseOutput(t *testing.T) {
	systemPrompt, userPrompt := BuildUIReviewPrompt("code", "登录页", "简洁现代", "检查可用性", "<button>提交</button>")

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
