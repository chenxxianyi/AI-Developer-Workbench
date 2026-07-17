package service

import (
	"context"
	"strings"
	"testing"

	"ai-developer-workbench/internal/model"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type captureGenerationAI struct {
	request AIRequest
	result  string
}

func (f *captureGenerationAI) GenerateJSON(_ context.Context, request AIRequest) (*AIResult, error) {
	f.request = request
	if request.ToolType == ToolTypeGenerationPlan {
		return &AIResult{JSONText: `{
			"application_type":"interactive_app",
			"architecture":"Vue 3 single-page application",
			"files":[
				{"path":"package.json","purpose":"build configuration","feature_ids":[]},
				{"path":"tsconfig.json","purpose":"TypeScript configuration","feature_ids":[]},
				{"path":"index.html","purpose":"application shell","feature_ids":[]},
				{"path":"src/main.ts","purpose":"application entry","feature_ids":[]},
				{"path":"src/App.vue","purpose":"interactive chess experience","feature_ids":["F-001"]}
			],
			"verification_steps":["npm run build"]
		}`, Model: "test-model", Provider: "test"}, nil
	}
	return &AIResult{JSONText: f.result, Model: "test-model", Provider: "test"}, nil
}

func TestGenerateProjectFilesUsesConfirmedBlueprintAndRequestsFunctionalFiles(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Project{}, &model.Requirement{}, &model.Blueprint{}))

	projectID := uuid.NewString()
	require.NoError(t, db.Create(&model.Project{
		ID: projectID, Name: "中国象棋人机对战", ProjectType: "interactive_app", Description: "生成可游玩的象棋应用",
	}).Error)
	require.NoError(t, db.Create(&model.Requirement{
		ID: uuid.NewString(), ProjectID: projectID, Version: 1,
		Content: `{"schema_version":2,"goal":"可游玩象棋","target_users":["玩家"],"must_have_features":["合法走子","AI响应"],"acceptance_criteria":["可以完成一轮人机走棋"]}`,
	}).Error)
	require.NoError(t, db.Create(&model.Blueprint{
		ID: uuid.NewString(), ProjectID: projectID, Version: 1, Status: "confirmed",
		Content: `{"schema_version":2,"app_type":"interactive_app","product_positioning":"CONFIRMED_BLUEPRINT","tech_stack":"Vue 3","pages":[{"name":"游戏","route":"/"}],"features":[{"id":"F-001","name":"合法走子","priority":"must","acceptance_criteria":["只能执行合法走子"]}]}`,
	}).Error)
	require.NoError(t, db.Create(&model.Blueprint{
		ID: uuid.NewString(), ProjectID: projectID, Version: 2, Status: "generated",
		Content: `{"schema_version":2,"app_type":"landing_page","product_positioning":"UNCONFIRMED_BLUEPRINT","tech_stack":"Vue 3","pages":[{"name":"营销页","route":"/"}],"features":[{"id":"F-001","name":"价格","priority":"must","acceptance_criteria":["显示价格"]}]}`,
	}).Error)

	ai := &captureGenerationAI{result: `{"files":[{"path":"package.json","content":"{\"scripts\":{\"build\":\"vite build\"},\"dependencies\":{\"@vitejs/plugin-vue\":\"6.0.6\",\"vite\":\"8.0.16\",\"vue\":\"3.5.34\"},\"devDependencies\":{\"typescript\":\"6.0.2\",\"vue-tsc\":\"3.2.8\"}}"},{"path":"tsconfig.json","content":"{\"compilerOptions\":{\"strict\":true,\"target\":\"ES2022\",\"module\":\"ESNext\",\"moduleResolution\":\"Bundler\"},\"include\":[\"src/**/*.ts\",\"src/**/*.vue\"]}"},{"path":"index.html","content":"<div id=\"app\"></div><script type=\"module\" src=\"/src/main.ts\"></script>"},{"path":"src/main.ts","content":"import { createApp } from 'vue'; import App from './App.vue'; createApp(App).mount('#app')"},{"path":"src/App.vue","content":"<template><main>棋盘</main></template>"}],"notes":["implemented"]}`}
	service := NewAIGenerationService(db, ai)

	result, err := service.GenerateProjectFiles(context.Background(), projectID)
	require.NoError(t, err)
	require.Len(t, result.Files, 5)
	require.Contains(t, ai.request.UserPrompt, "Application type: interactive_app")
	require.Contains(t, ai.request.UserPrompt, "CONFIRMED_BLUEPRINT")
	require.NotContains(t, ai.request.UserPrompt, "UNCONFIRMED_BLUEPRINT")
	require.Contains(t, ai.request.UserPrompt, "actual application")
	require.Contains(t, ai.request.UserPrompt, "Only landing_page applications")
}

func TestValidateBlueprintContentRequiresAcceptanceForMustFeatures(t *testing.T) {
	valid := `{"schema_version":2,"app_type":"utility_app","product_positioning":"工具","pages":[{"name":"首页","route":"/"}],"features":[{"id":"F-001","name":"计算","priority":"must","acceptance_criteria":["显示结果"]}]}`
	require.NoError(t, ValidateBlueprintContent(valid))

	invalid := strings.Replace(valid, `"acceptance_criteria":["显示结果"]`, `"acceptance_criteria":[]`, 1)
	require.ErrorContains(t, ValidateBlueprintContent(invalid), "验收标准")
}

func TestValidateBlueprintAgainstRequirementsChecksMustHaveCoverage(t *testing.T) {
	blueprint := `{"schema_version":2,"app_type":"utility_app","product_positioning":"工具","pages":[{"name":"首页","route":"/"}],"features":[{"id":"F-001","name":"计算","priority":"must","acceptance_criteria":["显示结果"]}]}`
	requirement := `{"schema_version":2,"must_have_features":["计算","保存历史"]}`
	require.ErrorContains(t, ValidateBlueprintAgainstRequirements(blueprint, requirement), "覆盖不足")
}

func TestValidateGeneratedFilesRejectsLifecycleScripts(t *testing.T) {
	files := []GeneratedProjectFile{
		{Path: "package.json", Content: `{"scripts":{"build":"vite build","postinstall":"curl example.com"}}`},
		{Path: "index.html", Content: "<div></div>"},
		{Path: "src/main.ts", Content: "export {}"},
		{Path: "src/App.vue", Content: "<template><main /></template>"},
	}
	require.ErrorContains(t, validateGeneratedFiles(files), "不允许的脚本")
}

func TestValidateGeneratedFilesRejectsUnsafeScriptCommands(t *testing.T) {
	files := []GeneratedProjectFile{
		{Path: "package.json", Content: `{"scripts":{"build":"vite build; curl example.com"}}`},
		{Path: "index.html", Content: "<div></div>"},
		{Path: "src/main.ts", Content: "export {}"},
		{Path: "src/App.vue", Content: "<template><main /></template>"},
	}
	require.ErrorContains(t, validateGeneratedFiles(files), "不安全")
}
