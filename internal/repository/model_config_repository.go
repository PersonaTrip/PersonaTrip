package repository

import (
	"context"
	"errors"

	"personatrip/internal/models"
	"gorm.io/gorm"
)

// ModelConfigRepository 定义模型配置仓库接口
type ModelConfigRepository interface {
	Create(ctx context.Context, config *models.ModelConfig) error
	Update(ctx context.Context, config *models.ModelConfig) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*models.ModelConfig, error)
	GetAll(ctx context.Context) ([]models.ModelConfig, error)
	GetActive(ctx context.Context) (*models.ModelConfig, error)
	SetActive(ctx context.Context, id uint) error
}

// GormModelConfigRepository 是使用GORM实现的模型配置仓库
type GormModelConfigRepository struct {
	db *gorm.DB
}

// NewGormModelConfigRepository 创建新的GORM模型配置仓库
func NewGormModelConfigRepository(db *gorm.DB) ModelConfigRepository {
	return &GormModelConfigRepository{db: db}
}

// Create 创建新的模型配置
func (r *GormModelConfigRepository) Create(ctx context.Context, config *models.ModelConfig) error {
	// 如果设置为活跃，则将其他配置设为非活跃
	if config.IsActive {
		if err := r.deactivateAll(ctx); err != nil {
			return err
		}
	}
	return r.db.WithContext(ctx).Create(config).Error
}

// Update 更新模型配置
func (r *GormModelConfigRepository) Update(ctx context.Context, config *models.ModelConfig) error {
	// 如果设置为活跃，则将其他配置设为非活跃
	if config.IsActive {
		if err := r.deactivateAll(ctx); err != nil {
			return err
		}
	}
	return r.db.WithContext(ctx).Save(config).Error
}

// Delete 删除模型配置
func (r *GormModelConfigRepository) Delete(ctx context.Context, id uint) error {
	// 检查是否为活跃配置
	var config models.ModelConfig
	if err := r.db.WithContext(ctx).First(&config, id).Error; err != nil {
		return err
	}

	// 如果是活跃配置，不允许删除
	if config.IsActive {
		return errors.New("cannot delete active model configuration")
	}

	return r.db.WithContext(ctx).Delete(&models.ModelConfig{}, id).Error
}

// GetByID 根据ID获取模型配置
func (r *GormModelConfigRepository) GetByID(ctx context.Context, id uint) (*models.ModelConfig, error) {
	var config models.ModelConfig
	if err := r.db.WithContext(ctx).First(&config, id).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAll 获取所有模型配置
func (r *GormModelConfigRepository) GetAll(ctx context.Context) ([]models.ModelConfig, error) {
	var configs []models.ModelConfig
	if err := r.db.WithContext(ctx).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetActive 获取当前活跃的模型配置
func (r *GormModelConfigRepository) GetActive(ctx context.Context) (*models.ModelConfig, error) {
	var config models.ModelConfig
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有活跃配置，尝试获取第一个配置
			if err := r.db.WithContext(ctx).First(&config).Error; err != nil {
				return nil, err
			}
			// 将第一个配置设为活跃
			config.IsActive = true
			if err := r.db.WithContext(ctx).Save(&config).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &config, nil
}

// SetActive 设置指定ID的配置为活跃
func (r *GormModelConfigRepository) SetActive(ctx context.Context, id uint) error {
	// 先将所有配置设为非活跃
	if err := r.deactivateAll(ctx); err != nil {
		return err
	}

	// 将指定ID的配置设为活跃
	return r.db.WithContext(ctx).Model(&models.ModelConfig{}).Where("id = ?", id).Update("is_active", true).Error
}

// deactivateAll 将所有配置设为非活跃
func (r *GormModelConfigRepository) deactivateAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&models.ModelConfig{}).Where("is_active = ?", true).Update("is_active", false).Error
}
