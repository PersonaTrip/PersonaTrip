package repository

import (
	"gorm.io/gorm"
)

// AdminRepository 管理员仓库接口已经存在

// ModelConfigRepository 模型配置仓库接口已经存在

// Database 定义数据库操作的抽象接口
type Database interface {
	// 通用操作
	GetDB() interface{}
	Close() error

	// 获取各个子仓库
	UserRepo() UserRepository
	AdminRepo() AdminRepository
	ModelConfigRepo() ModelConfigRepository
}

// GormDatabase 实现了Database接口的MySQL(GORM)版本
type GormDatabase struct {
	DB              *gorm.DB
	userRepo        UserRepository
	adminRepo       AdminRepository
	modelConfigRepo ModelConfigRepository
}

// NewGormDatabase 创建一个新的GORM数据库实例,新加入的模型必须修改的地方
func NewGormDatabase(db *gorm.DB) Database {
	return &GormDatabase{
		DB:              db,
		userRepo:        NewGormUserRepository(db),
		adminRepo:       NewGormAdminRepository(db),
		modelConfigRepo: NewGormModelConfigRepository(db),
	}
}

// GetDB 返回底层数据库连接
func (g *GormDatabase) GetDB() interface{} {
	return g.DB
}

// Close 关闭数据库连接
func (g *GormDatabase) Close() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// UserRepo 返回用户仓库
func (g *GormDatabase) UserRepo() UserRepository {
	return g.userRepo
}

// AdminRepo 返回管理员仓库
func (g *GormDatabase) AdminRepo() AdminRepository {
	return g.adminRepo
}

// ModelConfigRepo 返回模型配置仓库
func (g *GormDatabase) ModelConfigRepo() ModelConfigRepository {
	return g.modelConfigRepo
}
