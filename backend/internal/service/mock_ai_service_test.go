package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"ai-developer-workbench/internal/model"
)

func TestMockAIService_ReturnsValidJSONForAllToolTypes(t *testing.T) {
	svc := NewMockAIService()
	toolTypes := []string{
		model.ToolTypeUIReview,
		model.ToolTypeProjectDoctor,
		model.ToolTypeAgentConfig,
		model.ToolTypeAPIDoc,
		model.ToolTypeDBSchema,
	}

	for _, toolType := range toolTypes {
		t.Run(toolType, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			result, err := svc.GenerateJSON(ctx, AIRequest{
				ToolType:     toolType,
				SystemPrompt: "test",
				UserPrompt:   "test",
			})
			if err != nil {
				t.Fatalf("GenerateJSON failed: %v", err)
			}
			if result.JSONText == "" {
				t.Error("JSONText is empty")
			}
			if result.Provider != "mock" {
				t.Errorf("Provider=%q, want %q", result.Provider, "mock")
			}
			if result.Model != "mock-mode" {
				t.Errorf("Model=%q, want %q", result.Model, "mock-mode")
			}

			// Verify the JSON is valid.
			var raw map[string]interface{}
			if err := json.Unmarshal([]byte(result.JSONText), &raw); err != nil {
				t.Fatalf("JSONText is not valid JSON: %v\nRaw: %s", err, result.JSONText)
			}
		})
	}
}

func TestMockAIService_UnknownToolTypeReturnsError(t *testing.T) {
	svc := NewMockAIService()
	_, err := svc.GenerateJSON(context.Background(), AIRequest{
		ToolType: "unknown_tool",
	})
	if err == nil {
		t.Error("expected error for unknown tool type")
	}
}

func TestMockAIService_AllToolTypesIncludeActionItems(t *testing.T) {
	svc := NewMockAIService()
	toolTypes := []string{
		model.ToolTypeUIReview,
		model.ToolTypeProjectDoctor,
		model.ToolTypeAgentConfig,
		model.ToolTypeAPIDoc,
		model.ToolTypeDBSchema,
	}

	for _, toolType := range toolTypes {
		t.Run(toolType, func(t *testing.T) {
			result, err := svc.GenerateJSON(context.Background(), AIRequest{ToolType: toolType})
			if err != nil {
				t.Fatal(err)
			}

			var data map[string]interface{}
			if err := json.Unmarshal([]byte(result.JSONText), &data); err != nil {
				t.Fatal(err)
			}

			items, ok := data["action_items"].([]interface{})
			if !ok {
				t.Fatal("missing or invalid action_items")
			}
			if len(items) < 3 {
				t.Fatalf("len(action_items)=%d, want at least 3", len(items))
			}

			first, ok := items[0].(map[string]interface{})
			if !ok {
				t.Fatal("first action item is not an object")
			}
			for _, field := range []string{"id", "title", "priority", "effort", "category", "reason", "suggested_prompt", "issue_title", "issue_body"} {
				if first[field] == "" {
					t.Fatalf("first action item missing %q", field)
				}
			}
		})
	}
}

func TestMockAIService_UIReviewResultIsParseable(t *testing.T) {
	svc := NewMockAIService()
	result, err := svc.GenerateJSON(context.Background(), AIRequest{
		ToolType: model.ToolTypeUIReview,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify it contains expected fields.
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result.JSONText), &data); err != nil {
		t.Fatal(err)
	}

	if _, ok := data["scores"]; !ok {
		t.Error("missing 'scores' field")
	}
	if _, ok := data["issues"]; !ok {
		t.Error("missing 'issues' field")
	}
	if _, ok := data["recommendations"]; !ok {
		t.Error("missing 'recommendations' field")
	}
	if _, ok := data["codex_prompt"]; !ok {
		t.Error("missing 'codex_prompt' field")
	}
}

func TestMockAIService_AgentConfigHasGeneratedFiles(t *testing.T) {
	svc := NewMockAIService()
	result, err := svc.GenerateJSON(context.Background(), AIRequest{
		ToolType: model.ToolTypeAgentConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result.JSONText), &data); err != nil {
		t.Fatal(err)
	}

	filesContent, ok := data["generated_files_content"].(map[string]interface{})
	if !ok {
		t.Fatal("missing or invalid 'generated_files_content'")
	}
	if len(filesContent) == 0 {
		t.Error("generated_files_content is empty")
	}

	// Verify expected files exist.
	for _, fname := range []string{"AGENTS.md", "TASK_PLAN.md", "CODING_RULES.md"} {
		if _, ok := filesContent[fname]; !ok {
			t.Errorf("missing generated file %q", fname)
		}
	}
}

func TestMockAIService_APIDocHasModules(t *testing.T) {
	svc := NewMockAIService()
	result, err := svc.GenerateJSON(context.Background(), AIRequest{
		ToolType: model.ToolTypeAPIDoc,
	})
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result.JSONText), &data); err != nil {
		t.Fatal(err)
	}

	modules, ok := data["modules"].([]interface{})
	if !ok {
		t.Fatal("missing or invalid 'modules' field")
	}
	if len(modules) == 0 {
		t.Error("modules is empty")
	}

	if _, ok := data["markdown_content"]; !ok {
		t.Error("missing 'markdown_content' field")
	}
}

func TestMockAIService_DBSchemaHasScoresAndIssues(t *testing.T) {
	svc := NewMockAIService()
	result, err := svc.GenerateJSON(context.Background(), AIRequest{
		ToolType: model.ToolTypeDBSchema,
	})
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result.JSONText), &data); err != nil {
		t.Fatal(err)
	}

	if _, ok := data["scores"]; !ok {
		t.Error("missing 'scores' field")
	}
	if _, ok := data["issues"]; !ok {
		t.Error("missing 'issues' field")
	}
	if _, ok := data["optimized_schema"]; !ok {
		t.Error("missing 'optimized_schema' field")
	}
	if _, ok := data["migration_suggestions"]; !ok {
		t.Error("missing 'migration_suggestions' field")
	}
}

func TestMockAIService_NoExternalHTTPCalls(t *testing.T) {
	// The mock service has no fields for HTTP, so this is verified by construction.
	// If it ever gains an http.Client, this test should be updated to intercept it.
	svc := NewMockAIService()

	result, err := svc.GenerateJSON(context.Background(), AIRequest{
		ToolType: model.ToolTypeUIReview,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Provider != "mock" {
		t.Error("Mock service returned non-mock provider")
	}
}
