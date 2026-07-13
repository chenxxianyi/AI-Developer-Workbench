package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"ai-developer-workbench/internal/dto"
)

// issueKey returns a stable matching key for an issue. IssueItem has no id, so
// we fall back to normalized category + title.
func issueKey(category, title string) string {
	c := strings.ToLower(strings.TrimSpace(category))
	t := strings.ToLower(strings.TrimSpace(title))
	return c + "|" + t
}

// extractIssues pulls issues out of a report_data JSON blob regardless of tool.
func extractIssues(raw json.RawMessage) []dto.IssueItem {
	var probe struct {
		Issues []dto.IssueItem `json:"issues"`
	}
	if err := json.Unmarshal(raw, &probe); err != nil {
		return nil
	}
	return probe.Issues
}

// extractActionItems pulls action_items out of a report_data JSON blob.
func extractActionItems(raw json.RawMessage) []dto.ActionItem {
	var probe struct {
		ActionItems []dto.ActionItem `json:"action_items"`
	}
	if err := json.Unmarshal(raw, &probe); err != nil {
		return nil
	}
	return probe.ActionItems
}

// IssueCountDeltaFrom computes severity count deltas (target - baseline).
func IssueCountDeltaFrom(baseline, target []dto.IssueItem) dto.IssueCountDelta {
	countBySev := func(items []dto.IssueItem) (high, med, low int) {
		for _, it := range items {
			switch it.Severity {
			case "high":
				high++
			case "medium":
				med++
			case "low":
				low++
			}
		}
		return
	}
	bh, bm, bl := countBySev(baseline)
	th, tm, tl := countBySev(target)
	return dto.IssueCountDelta{
		High:   th - bh,
		Medium: tm - bm,
		Low:    tl - bl,
		Total:  len(target) - len(baseline),
	}
}

// compareIssues matches issues by category+title (no stable id exists on IssueItem).
func compareIssues(baseline, target []dto.IssueItem) dto.IssueComparison {
	baseMap := make(map[string]dto.IssueItem, len(baseline))
	for _, it := range baseline {
		baseMap[issueKey(it.Category, it.Title)] = it
	}
	tgtMap := make(map[string]dto.IssueItem, len(target))
	for _, it := range target {
		tgtMap[issueKey(it.Category, it.Title)] = it
	}

	var resolved, newItems, persist []dto.IssueMatch
	for k, b := range baseMap {
		if t, ok := tgtMap[k]; ok {
			persist = append(persist, dto.IssueMatch{Title: b.Title, Category: b.Category, Severity: t.Severity, InBaseline: true, InTarget: true})
		} else {
			resolved = append(resolved, dto.IssueMatch{Title: b.Title, Category: b.Category, Severity: b.Severity, InBaseline: true, InTarget: false})
		}
	}
	for k, t := range tgtMap {
		if _, ok := baseMap[k]; !ok {
			newItems = append(newItems, dto.IssueMatch{Title: t.Title, Category: t.Category, Severity: t.Severity, InBaseline: false, InTarget: true})
		}
	}
	return dto.IssueComparison{Resolved: resolved, New: newItems, Persist: persist}
}

// compareActionItems matches action items by stable id, falling back to category+title.
func compareActionItems(baseline, target []dto.ActionItem) dto.ActionItemDelta {
	baseByID := make(map[string]dto.ActionItem, len(baseline))
	for _, a := range baseline {
		baseByID[a.ID] = a
	}
	tgtByID := make(map[string]dto.ActionItem, len(target))
	for _, a := range target {
		tgtByID[a.ID] = a
	}

	var resolved, newItems, persist []dto.ActionItem
	for id, b := range baseByID {
		if t, ok := tgtByID[id]; ok {
			// prefer the target version (most recent state)
			persist = append(persist, t)
		} else {
			// fallback to category+title for cross-version stability
			matched := false
			for _, t := range target {
				if t.Category == b.Category && t.Title == b.Title {
					matched = true
					persist = append(persist, t)
					break
				}
			}
			if !matched {
				resolved = append(resolved, b)
			}
		}
	}
	for id, t := range tgtByID {
		if _, ok := baseByID[id]; ok {
			continue
		}
		matched := false
		for _, b := range baseline {
			if b.Category == t.Category && b.Title == t.Title {
				matched = true
				break
			}
		}
		if !matched {
			newItems = append(newItems, t)
		}
	}
	return dto.ActionItemDelta{Resolved: resolved, New: newItems, Persist: persist}
}

// gradeDelta returns a human-readable grade change (e.g. "B → A", "A → A (unchanged)").
func gradeDelta(baseline, target string) string {
	if baseline == target {
		return fmt.Sprintf("%s → %s (unchanged)", baseline, target)
	}
	return fmt.Sprintf("%s → %s", baseline, target)
}
