package repository

import (
	"context"
	"errors"
	"sync"

	"personatrip/internal/models"
)

// MemoryAdminRepository 是使用内存实现的管理员仓库
type MemoryAdminRepository struct {
	admins map[uint]*models.Admin
	mutex  sync.RWMutex
	nextID uint
}

// NewMemoryAdminRepository 创建新的内存管理员仓库
func NewMemoryAdminRepository() AdminRepository {
	return &MemoryAdminRepository{
		admins: make(map[uint]*models.Admin),
		nextID: 1,
	}
}

// Create 创建新的管理员
func (r *MemoryAdminRepository) Create(ctx context.Context, admin *models.Admin) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查用户名是否已存在
	for _, a := range r.admins {
		if a.Username == admin.Username {
			return errors.New("username already exists")
		}
	}

	// 设置ID
	admin.ID = r.nextID
	r.nextID++

	// 存储管理员
	r.admins[admin.ID] = admin
	return nil
}

// Update 更新管理员信息
func (r *MemoryAdminRepository) Update(ctx context.Context, admin *models.Admin) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查ID是否存在
	if _, exists := r.admins[admin.ID]; !exists {
		return errors.New("admin not found")
	}

	// 更新管理员
	r.admins[admin.ID] = admin
	return nil
}

// Delete 删除管理员
func (r *MemoryAdminRepository) Delete(ctx context.Context, id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查ID是否存在
	if _, exists := r.admins[id]; !exists {
		return errors.New("admin not found")
	}

	// 删除管理员
	delete(r.admins, id)
	return nil
}

// GetByID 根据ID获取管理员
func (r *MemoryAdminRepository) GetByID(ctx context.Context, id uint) (*models.Admin, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 检查ID是否存在
	admin, exists := r.admins[id]
	if !exists {
		return nil, errors.New("admin not found")
	}

	return admin, nil
}

// GetByUsername 根据用户名获取管理员
func (r *MemoryAdminRepository) GetByUsername(ctx context.Context, username string) (*models.Admin, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 查找用户名
	for _, admin := range r.admins {
		if admin.Username == username {
			return admin, nil
		}
	}

	return nil, errors.New("admin not found")
}

// GetAll 获取所有管理员
func (r *MemoryAdminRepository) GetAll(ctx context.Context) ([]models.Admin, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	admins := make([]models.Admin, 0, len(r.admins))
	for _, admin := range r.admins {
		admins = append(admins, *admin)
	}

	return admins, nil
}
