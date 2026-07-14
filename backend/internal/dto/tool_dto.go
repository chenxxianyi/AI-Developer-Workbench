package dto

// ToolMetaDTO matches frontend ToolMeta type.
type ToolMetaDTO struct {
	ToolType   string `json:"tool_type"`
	Name       string `json:"name"`
	Desc       string `json:"description"`
	Icon       string `json:"icon"`
	Color      string `json:"color"`
	UsageCount int64  `json:"usage_count"`
}

// Tool metadata definitions exposed by the tools API.
var ToolMetaList = []ToolMetaDTO{
	{
		ToolType:   "ui_review",
		Name:       "UI Review",
		Desc:       "AI-powered UI/UX review with screenshot or code analysis",
		Icon:       "eye",
		Color:      "blue",
		UsageCount: 0,
	},
	{
		ToolType:   "project_doctor",
		Name:       "Project Doctor",
		Desc:       "Comprehensive project health check and engineering quality analysis",
		Icon:       "stethoscope",
		Color:      "green",
		UsageCount: 0,
	},
	{
		ToolType:   "agent_config",
		Name:       "Agent Config Studio",
		Desc:       "Generate AI agent configuration files for your project",
		Icon:       "settings",
		Color:      "purple",
		UsageCount: 0,
	},
	{
		ToolType:   "api_doc",
		Name:       "API Doc Builder",
		Desc:       "Auto-generate API documentation from code or descriptions",
		Icon:       "file-text",
		Color:      "orange",
		UsageCount: 0,
	},
	{
		ToolType:   "db_schema",
		Name:       "DB Schema Review",
		Desc:       "Review and optimize database schema with AI analysis",
		Icon:       "database",
		Color:      "teal",
		UsageCount: 0,
	},
}

// GetToolMetaByType returns tool metadata for a given tool type.
func GetToolMetaByType(toolType string) *ToolMetaDTO {
	for i := range ToolMetaList {
		if ToolMetaList[i].ToolType == toolType {
			return &ToolMetaList[i]
		}
	}
	return nil
}
