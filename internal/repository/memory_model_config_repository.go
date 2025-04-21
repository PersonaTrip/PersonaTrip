package repository

import (
	"context"
	"errors"
	"sync"

	"personatrip/internal/models"
)

// MemoryModelConfigRepository 是使用内存实现的模型配置仓库
type MemoryModelConfigRepository struct {
	configs map[uint]*models.ModelConfig
	mutex   sync.RWMutex
	nextID  uint
}

// NewMemoryModelConfigRepository 创建新的内存模型配置仓库
func NewMemoryModelConfigRepository() ModelConfigRepository {
	return &MemoryModelConfigRepository{
		configs: make(map[uint]*models.ModelConfig),
		nextID:  1,
	}
}

// Create 创建新的模型配置
func (r *MemoryModelConfigRepository) Create(ctx context.Context, config *models.ModelConfig) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 设置ID
	config.ID = r.nextID
	r.nextID++

	// 存储配置
	r.configs[config.ID] = config
	return nil
}

// Update 更新模型配置
func (r *MemoryModelConfigRepository) Update(ctx context.Context, config *models.ModelConfig) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查ID是否存在
	if _, exists := r.configs[config.ID]; !exists {
		return errors.New("model config not found")
	}

	// 更新配置
	r.configs[config.ID] = config
	return nil
}

// Delete 删除模型配置
func (r *MemoryModelConfigRepository) Delete(ctx context.Context, id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查ID是否存在
	if _, exists := r.configs[id]; !exists {
		return errors.New("model config not found")
	}

	// 删除配置
	delete(r.configs, id)
	return nil
}

// GetByID 根据ID获取模型配置
func (r *MemoryModelConfigRepository) GetByID(ctx context.Context, id uint) (*models.ModelConfig, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 检查ID是否存在
	config, exists := r.configs[id]
	if !exists {
		return nil, errors.New("model config not found")
	}

	return config, nil
}

// GetAll 获取所有模型配置
func (r *MemoryModelConfigRepository) GetAll(ctx context.Context) ([]models.ModelConfig, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	configs := make([]models.ModelConfig, 0, len(r.configs))
	for _, config := range r.configs {
		configs = append(configs, *config)
	}

	return configs, nil
}

// GetActive 获取当前活跃的模型配置
func (r *MemoryModelConfigRepository) GetActive(ctx context.Context) (*models.ModelConfig, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 查找活跃配置
	for _, config := range r.configs {
		if config.IsActive {
			return config, nil
		}
	}

	return nil, errors.New("no active model config found")
}

// SetActive 设置指定ID的配置为活跃
func (r *MemoryModelConfigRepository) SetActive(ctx context.Context, id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查ID是否存在
	targetConfig, exists := r.configs[id]
	if !exists {
		return errors.New("model config not found")
	}

	// 将所有配置设置为非活跃
	for _, config := range r.configs {
		config.IsActive = false
	}

	// 将目标配置设置为活跃
	targetConfig.IsActive = true
	return nil
}
