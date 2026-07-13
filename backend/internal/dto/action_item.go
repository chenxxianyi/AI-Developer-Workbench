package dto

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	ActionPriorityHigh   = "high"
	ActionPriorityMedium = "medium"
	ActionPriorityLow    = "low"

	ActionEffortSmall  = "small"
	ActionEffortMedium = "medium"
	ActionEffortLarge  = "large"
)

var actionIDInvalidChars = regexp.MustCompile(`[^a-z0-9_-]+`)

// ActionItem is the shared actionable task contract returned by all tools.
type ActionItem struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Priority        string `json:"priority"`
	Effort          string `json:"effort"`
	Category        string `json:"category"`
	Reason          string `json:"reason"`
	SuggestedPrompt string `json:"suggested_prompt"`
	IssueTitle      string `json:"issue_title"`
	IssueBody       string `json:"issue_body"`
}

// NormalizeActionItems keeps action item enums and IDs predictable for the UI.
func NormalizeActionItems(items []ActionItem) []ActionItem {
	if len(items) == 0 {
		return []ActionItem{}
	}

	seenIDs := map[string]int{}
	normalized := make([]ActionItem, 0, len(items))
	for index, item := range items {
		item = trimActionItem(item)
		item.Priority = normalizeActionPriority(item.Priority)
		item.Effort = normalizeActionEffort(item.Effort)
		if item.Category == "" {
			item.Category = "general"
		}
		if item.Title == "" {
			item.Title = fmt.Sprintf("行动项 %d", index+1)
		}
		if item.IssueTitle == "" {
			item.IssueTitle = item.Title
		}

		baseID := normalizeActionID(item.ID, index)
		count := seenIDs[baseID]
		seenIDs[baseID] = count + 1
		if count > 0 {
			item.ID = fmt.Sprintf("%s-%d", baseID, count+1)
		} else {
			item.ID = baseID
		}

		normalized = append(normalized, item)
	}

	return normalized
}

func trimActionItem(item ActionItem) ActionItem {
	item.ID = strings.TrimSpace(item.ID)
	item.Title = strings.TrimSpace(item.Title)
	item.Priority = strings.TrimSpace(strings.ToLower(item.Priority))
	item.Effort = strings.TrimSpace(strings.ToLower(item.Effort))
	item.Category = strings.TrimSpace(item.Category)
	item.Reason = strings.TrimSpace(item.Reason)
	item.SuggestedPrompt = strings.TrimSpace(item.SuggestedPrompt)
	item.IssueTitle = strings.TrimSpace(item.IssueTitle)
	item.IssueBody = strings.TrimSpace(item.IssueBody)
	return item
}

func normalizeActionPriority(priority string) string {
	switch priority {
	case ActionPriorityHigh, ActionPriorityMedium, ActionPriorityLow:
		return priority
	default:
		return ActionPriorityMedium
	}
}

func normalizeActionEffort(effort string) string {
	switch effort {
	case ActionEffortSmall, ActionEffortMedium, ActionEffortLarge:
		return effort
	default:
		return ActionEffortMedium
	}
}

func normalizeActionID(id string, index int) string {
	id = strings.ToLower(strings.TrimSpace(id))
	id = actionIDInvalidChars.ReplaceAllString(id, "-")
	id = strings.Trim(id, "-_")
	if id == "" {
		return fmt.Sprintf("action-%02d", index+1)
	}
	return id
}
