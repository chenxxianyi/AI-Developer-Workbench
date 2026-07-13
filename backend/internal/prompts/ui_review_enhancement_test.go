package prompts

import (
	"strings"
	"testing"
)

func TestUIReviewPromptRequiresViewportAnnotations(t *testing.T) {
	system, _ := BuildUIReviewPrompt("screenshot", "paste", "dashboard", "", "", "", "")
	for _, want := range []string{"desktop or mobile", "region coordinates", "component_prompt", "WCAG"} {
		if !strings.Contains(system, want) {
			t.Fatalf("prompt missing %q", want)
		}
	}
}
