package util

import (
	"encoding/json"
	"regexp"
	"strings"
)

// ParseAIResponse attempts to extract valid JSON from AI response text.
func ParseAIResponse[T any](rawText string) (*T, error) {
	var result T

	// 1. Try direct unmarshal.
	if err := json.Unmarshal([]byte(rawText), &result); err == nil {
		return &result, nil
	}

	// 2. Strip Markdown code fences.
	cleaned := stripCodeFences(rawText)
	if err := json.Unmarshal([]byte(cleaned), &result); err == nil {
		return &result, nil
	}

	// 3. Extract first balanced JSON object.
	jsonStr := extractBalancedJSON(cleaned)
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &result); err == nil {
			return &result, nil
		}
	}

	// 4. Try parsing as map then converting.
	var rawMap map[string]interface{}
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &rawMap); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal([]byte(cleaned), &rawMap); err != nil {
			return nil, err
		}
	}

	NormalizeMap(rawMap)
	convertMapToStruct(rawMap, &result)

	return &result, nil
}

// ParseAIResponseInto parses AI response into an existing struct.
func ParseAIResponseInto(rawText string, target interface{}) error {
	// 1. Try direct unmarshal.
	if err := json.Unmarshal([]byte(rawText), target); err == nil {
		return nil
	}

	// 2. Strip code fences.
	cleaned := stripCodeFences(rawText)
	if err := json.Unmarshal([]byte(cleaned), target); err == nil {
		return nil
	}

	// 3. Extract balanced JSON.
	jsonStr := extractBalancedJSON(cleaned)
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), target); err == nil {
			return nil
		}
	}

	return json.Unmarshal([]byte(jsonStr), target)
}

// stripCodeFences removes Markdown code fences from text.
func stripCodeFences(text string) string {
	re := regexp.MustCompile("(?s)```(?:json)?\\s*([\\s\\S]*?)\\s*```")
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return strings.TrimSpace(text)
}

// extractBalancedJSON extracts the first balanced JSON object or array.
func extractBalancedJSON(text string) string {
	start := strings.Index(text, "{")
	arrStart := strings.Index(text, "[")

	if arrStart >= 0 && (start < 0 || arrStart < start) {
		start = arrStart
	}

	if start < 0 {
		return ""
	}

	openChar := text[start]
	closeChar := byte('}')
	if openChar == '[' {
		closeChar = byte(']')
	}

	depth := 0
	for i := start; i < len(text); i++ {
		c := text[i]
		if c == byte(openChar) {
			depth++
		} else if c == closeChar {
			depth--
			if depth == 0 {
				return text[start : i+1]
			}
		}
		// Skip escaped quotes inside strings.
		if c == '"' {
			j := i + 1
			for j < len(text) {
				if text[j] == '\\' {
					j += 2
				} else if text[j] == '"' {
					break
				} else {
					j++
				}
			}
			i = j
		}
	}

	return ""
}

// NormalizeMap normalizes values in a map for consistency.
func NormalizeMap(data map[string]interface{}) {
	for k, v := range data {
		switch val := v.(type) {
		case float64:
			if strings.Contains(strings.ToLower(k), "score") {
				data[k] = NormalizeScore(int(val))
			}
		case string:
			if strings.Contains(strings.ToLower(k), "severity") {
				data[k] = NormalizeSeverity(val)
			}
		case map[string]interface{}:
			NormalizeMap(val)
		case []interface{}:
			NormalizeSlice(val)
		case nil:
			if strings.Contains(strings.ToLower(k), "s") ||
				strings.Contains(strings.ToLower(k), "items") ||
				strings.Contains(strings.ToLower(k), "list") {
				data[k] = []interface{}{}
			}
		}
	}
}

// NormalizeSlice normalizes values in a slice.
func NormalizeSlice(data []interface{}) {
	for i, v := range data {
		switch val := v.(type) {
		case map[string]interface{}:
			NormalizeMap(val)
			data[i] = val
		case []interface{}:
			NormalizeSlice(val)
			data[i] = val
		}
	}
}

// convertMapToStruct converts a map to a struct via JSON round-trip.
func convertMapToStruct[T any](m map[string]interface{}, target *T) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}