package services

import (
	"context"

	"personatrip/internal/models"
	"personatrip/internal/repository"
)

// ModelConfigService 定义模型配置服务接口
type ModelConfigService interface {
	CreateModelConfig(ctx context.Context, req *models.ModelConfigCreateRequest) (*models.ModelConfig, error)
	UpdateModelConfig(ctx context.Context, id uint, req *models.ModelConfigUpdateRequest) (*models.ModelConfig, error)
	DeleteModelConfig(ctx context.Context, id uint) error
	GetModelConfigByID(ctx context.Context, id uint) (*models.ModelConfig, error)
	GetAllModelConfigs(ctx context.Context) ([]models.ModelConfig, error)
	GetActiveModelConfig(ctx context.Context) (*models.ModelConfig, error)
	SetActiveModelConfig(ctx context.Context, id uint) error
}

// ModelConfigServiceImpl 是模型配置服务的实现
type ModelConfigServiceImpl struct {
	db repository.Database
}

// NewModelConfigService 创建新的模型配置服务
func NewModelConfigService(db repository.Database) ModelConfigService {
	return &ModelConfigServiceImpl{db: db}
}

// CreateModelConfig 创建新的模型配置
func (s *ModelConfigServiceImpl) CreateModelConfig(ctx context.Context, req *models.ModelConfigCreateRequest) (*models.ModelConfig, error) {
	config := &models.ModelConfig{
		Name:        req.Name,
		ModelType:   req.ModelType,
		ModelName:   req.ModelName,
		ApiKey:      req.ApiKey,
		BaseUrl:     req.BaseUrl,
		IsActive:    req.IsActive,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	}

	// 设置默认值
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 2000
	}

	if err := s.db.ModelConfigRepo().Create(ctx, config); err != nil {
		return nil, err
	}
	return config, nil
}

// UpdateModelConfig 更新模型配置
func (s *ModelConfigServiceImpl) UpdateModelConfig(ctx context.Context, id uint, req *models.ModelConfigUpdateRequest) (*models.ModelConfig, error) {
	config, err := s.db.ModelConfigRepo().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != "" {
		config.Name = req.Name
	}
	if req.ModelType != "" {
		config.ModelType = req.ModelType
	}
	if req.ModelName != "" {
		config.ModelName = req.ModelName
	}
	if req.ApiKey != "" {
		config.ApiKey = req.ApiKey
	}
	if req.BaseUrl != "" {
		config.BaseUrl = req.BaseUrl
	}
	config.IsActive = req.IsActive
	if req.Temperature != 0 {
		config.Temperature = req.Temperature
	}
	if req.MaxTokens != 0 {
		config.MaxTokens = req.MaxTokens
	}

	if err := s.db.ModelConfigRepo().Update(ctx, config); err != nil {
		return nil, err
	}
	return config, nil
}

// DeleteModelConfig 删除模型配置
func (s *ModelConfigServiceImpl) DeleteModelConfig(ctx context.Context, id uint) error {
	return s.db.ModelConfigRepo().Delete(ctx, id)
}

// GetModelConfigByID 根据ID获取模型配置
func (s *ModelConfigServiceImpl) GetModelConfigByID(ctx context.Context, id uint) (*models.ModelConfig, error) {
	return s.db.ModelConfigRepo().GetByID(ctx, id)
}

// GetAllModelConfigs 获取所有模型配置
func (s *ModelConfigServiceImpl) GetAllModelConfigs(ctx context.Context) ([]models.ModelConfig, error) {
	return s.db.ModelConfigRepo().GetAll(ctx)
}

// GetActiveModelConfig 获取当前活跃的模型配置
func (s *ModelConfigServiceImpl) GetActiveModelConfig(ctx context.Context) (*models.ModelConfig, error) {
	return s.db.ModelConfigRepo().GetActive(ctx)
}

// SetActiveModelConfig 设置指定ID的配置为活跃
func (s *ModelConfigServiceImpl) SetActiveModelConfig(ctx context.Context, id uint) error {
	return s.db.ModelConfigRepo().SetActive(ctx, id)
}
