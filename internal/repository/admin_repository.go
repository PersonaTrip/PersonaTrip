package repository

import (
	"context"
	"errors"

	"personatrip/internal/models"
	"gorm.io/gorm"
)

// AdminRepository 定义管理员仓库接口
type AdminRepository interface {
	Create(ctx context.Context, admin *models.Admin) error
	Update(ctx context.Context, admin *models.Admin) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*models.Admin, error)
	GetByUsername(ctx context.Context, username string) (*models.Admin, error)
	GetAll(ctx context.Context) ([]models.Admin, error)
}

// GormAdminRepository 是使用GORM实现的管理员仓库
type GormAdminRepository struct {
	db *gorm.DB
}

// NewGormAdminRepository 创建新的GORM管理员仓库
func NewGormAdminRepository(db *gorm.DB) AdminRepository {
	return &GormAdminRepository{db: db}
}

// Create 创建新的管理员
func (r *GormAdminRepository) Create(ctx context.Context, admin *models.Admin) error {
	return r.db.WithContext(ctx).Create(admin).Error
}

// Update 更新管理员信息
func (r *GormAdminRepository) Update(ctx context.Context, admin *models.Admin) error {
	return r.db.WithContext(ctx).Save(admin).Error
}

// Delete 删除管理员
func (r *GormAdminRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Admin{}, id).Error
}

// GetByID 根据ID获取管理员
func (r *GormAdminRepository) GetByID(ctx context.Context, id uint) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.WithContext(ctx).First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetByUsername 根据用户名获取管理员
func (r *GormAdminRepository) GetByUsername(ctx context.Context, username string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}
	return &admin, nil
}

// GetAll 获取所有管理员
func (r *GormAdminRepository) GetAll(ctx context.Context) ([]models.Admin, error) {
	var admins []models.Admin
	if err := r.db.WithContext(ctx).Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}
