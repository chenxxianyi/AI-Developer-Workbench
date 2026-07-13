package dto

import "testing"

func TestNormalizeActionItems_NormalizesEnumsAndDefaults(t *testing.T) {
	items := NormalizeActionItems([]ActionItem{
		{
			ID:              " Fix Upload ",
			Title:           " 修复上传键盘操作 ",
			Priority:        "urgent",
			Effort:          "tiny",
			SuggestedPrompt: " 请补齐 Enter/Space 支持 ",
		},
	})

	if len(items) != 1 {
		t.Fatalf("len(items)=%d, want 1", len(items))
	}
	if items[0].ID != "fix-upload" {
		t.Fatalf("ID=%q, want fix-upload", items[0].ID)
	}
	if items[0].Priority != ActionPriorityMedium {
		t.Fatalf("Priority=%q, want medium", items[0].Priority)
	}
	if items[0].Effort != ActionEffortMedium {
		t.Fatalf("Effort=%q, want medium", items[0].Effort)
	}
	if items[0].Category != "general" {
		t.Fatalf("Category=%q, want general", items[0].Category)
	}
	if items[0].IssueTitle != "修复上传键盘操作" {
		t.Fatalf("IssueTitle=%q, want title fallback", items[0].IssueTitle)
	}
}

func TestNormalizeActionItems_GeneratesStableUniqueIDs(t *testing.T) {
	items := NormalizeActionItems([]ActionItem{
		{ID: "same", Title: "A", Priority: "high", Effort: "small"},
		{ID: "same", Title: "B", Priority: "low", Effort: "large"},
		{Title: "C"},
	})

	wantIDs := []string{"same", "same-2", "action-03"}
	for i, want := range wantIDs {
		if items[i].ID != want {
			t.Fatalf("items[%d].ID=%q, want %q", i, items[i].ID, want)
		}
	}
}

func TestNormalizeActionItems_EmptyInputReturnsEmptySlice(t *testing.T) {
	items := NormalizeActionItems(nil)
	if items == nil {
		t.Fatal("NormalizeActionItems(nil) returned nil, want empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("len(items)=%d, want 0", len(items))
	}
}
