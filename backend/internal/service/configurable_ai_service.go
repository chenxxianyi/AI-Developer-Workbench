package service

import (
	"context"

	"ai-developer-workbench/internal/config"
	"ai-developer-workbench/internal/model"

	"gorm.io/gorm"
)

// ConfigurableAIService resolves the active/default model preset on every
// request. This makes the admin model management page affect real AI calls
// without requiring a backend restart. Secrets are still read from .env; model
// presets intentionally do not expose or persist API keys.
type ConfigurableAIService struct {
	db      *gorm.DB
	baseCfg config.AIConfig
}

func NewConfigurableAIService(db *gorm.DB, cfg *config.AIConfig) *ConfigurableAIService {
	base := config.AIConfig{}
	if cfg != nil {
		base = *cfg
	}
	return &ConfigurableAIService{db: db, baseCfg: base}
}

func (s *ConfigurableAIService) GenerateJSON(ctx context.Context, input AIRequest) (*AIResult, error) {
	cfg := s.resolveConfig()
	return NewOpenAICompatibleService(&cfg).GenerateJSON(ctx, input)
}

func (s *ConfigurableAIService) resolveConfig() config.AIConfig {
	cfg := s.baseCfg
	if s.db == nil {
		return cfg
	}
	var preset model.ModelPreset
	if err := s.db.Where("is_default = ? AND status = ?", true, "active").Order("updated_at desc").First(&preset).Error; err != nil {
		return cfg
	}
	if preset.Provider != "" {
		cfg.Provider = preset.Provider
	}
	if preset.BaseURL != "" {
		cfg.BaseURL = preset.BaseURL
	}
	if preset.Model != "" {
		cfg.Model = preset.Model
	}
	if preset.VisionModel != "" {
		cfg.VisionModel = preset.VisionModel
	}
	if preset.TimeoutSeconds > 0 {
		cfg.TimeoutSeconds = preset.TimeoutSeconds
	}
	if preset.MaxRetries >= 0 {
		cfg.MaxRetries = preset.MaxRetries
	}
	return cfg
}
