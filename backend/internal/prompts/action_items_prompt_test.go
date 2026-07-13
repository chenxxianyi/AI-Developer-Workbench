package prompts

import (
	"strings"
	"testing"
)

func TestToolPromptsIncludeActionItemsContract(t *testing.T) {
	tests := []struct {
		name   string
		system string
	}{
		{"ui review", first(BuildUIReviewPrompt("code", "paste", "dashboard", "quiet", "", "code", ""))},
		{"project doctor", first(BuildProjectDoctorPrompt("demo", "Vue Go", "", "standard", "summary"))},
		{"agent config", first(BuildAgentConfigPrompt("demo", "web", "Vue", "Go", "MySQL", "quiet", "", ""))},
		{"api doc", first(BuildAPIDocPrompt("manual", "Gin", "", "GET /api/reports", "both", ""))},
		{"db schema", first(BuildDBSchemaPrompt("sql", "mysql", "", "CREATE TABLE reports(id int);", "review"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			required := []string{
				`"action_items"`,
				`"priority": "high|medium|low"`,
				`"effort": "small|medium|large"`,
				`"suggested_prompt"`,
				`"issue_body"`,
				"at least 3 and at most 12 action items",
				"SAFETY BOUNDARY",
			}
			for _, want := range required {
				if !strings.Contains(tt.system, want) {
					t.Fatalf("system prompt missing %q\nprompt:\n%s", want, tt.system)
				}
			}
		})
	}
}

func TestUIReviewPromptKeepsChineseActionItemRequirement(t *testing.T) {
	systemPrompt, _ := BuildUIReviewPrompt("code", "paste", "dashboard", "professional", "", "<button>提交</button>", "")
	required := []string{
		"action_items 中 title、reason、suggested_prompt、issue_title 和 issue_body 必须使用简体中文",
		"必须使用简体中文输出所有面向用户的内容",
	}
	for _, want := range required {
		if !strings.Contains(systemPrompt, want) {
			t.Fatalf("UI review prompt missing %q", want)
		}
	}
}

func first(system, _ string) string {
	return system
}
