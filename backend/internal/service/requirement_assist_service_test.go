package service

import (
	"context"
	"testing"

	"ai-developer-workbench/internal/model"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRequirementAssistPreservesConfirmedFieldsAndNormalizesSuggestions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Project{}))
	projectID := uuid.NewString()
	require.NoError(t, db.Create(&model.Project{ID: projectID, Name: "象棋练习", ProjectType: "interactive_app"}).Error)

	ai := &captureGenerationAI{result: `{
		"spec":{
			"schema_version":2,
			"app_type":"landing_page",
			"goal":"AI 改写的目标",
			"target_users":["象棋玩家"],
			"primary_scenarios":["打开后开始下棋"],
			"must_have_features":["合法走棋","电脑回应","合法走棋"],
			"should_have_features":["难度选择"],
			"screens":["游戏主界面"],
			"interaction_rules":["红方先行"],
			"data_and_storage":{"persistence":"不用保存","backend_required":false},
			"visual_preferences":{"style":"传统中式","primary_color":"朱红与木色"},
			"responsive_targets":["desktop","mobile"],
			"non_functional_requirements":["手机可用"],
			"acceptance_criteria":["玩家走棋后电脑自动回应"],
			"out_of_scope":["登录"]
		},
		"inferred_fields":["must_have_features"],
		"questions":[],
		"ready":true,
		"summary":{"product":"象棋应用","users":"象棋玩家","features":[],"success":[]}
	}`}
	service := NewRequirementAssistService(db, ai)
	result, err := service.Assist(context.Background(), projectID, RequirementAssistInput{
		Description:     "我想做一个可以和电脑下象棋的网站",
		CurrentStep:     "idea",
		ConfirmedFields: []string{"app_type", "goal"},
		CurrentSpec:     RequirementAssistSpec{SchemaVersion: 2, AppType: "interactive_app", Goal: "用户确认的目标"},
	})

	require.NoError(t, err)
	require.Equal(t, "interactive_app", result.Spec.AppType)
	require.Equal(t, "用户确认的目标", result.Spec.Goal)
	require.Equal(t, []string{"合法走棋", "电脑回应"}, result.Spec.MustHaveFeatures)
	require.Equal(t, result.Spec.MustHaveFeatures, result.Summary.Features)
	require.Equal(t, result.Spec.AcceptanceCriteria, result.Summary.Success)
	require.True(t, result.Ready)
	require.Equal(t, ToolTypeRequirementAssist, ai.request.ToolType)
	require.Contains(t, ai.request.UserPrompt, "不得擅自增加收费")
}

func TestRequirementAssistRequiresDescription(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Project{}))
	projectID := uuid.NewString()
	require.NoError(t, db.Create(&model.Project{ID: projectID, Name: "空需求", ProjectType: "utility_app"}).Error)

	service := NewRequirementAssistService(db, &captureGenerationAI{})
	_, err = service.Assist(context.Background(), projectID, RequirementAssistInput{})
	require.ErrorContains(t, err, "请先简单描述")
}
