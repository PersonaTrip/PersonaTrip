package repository

import (
	"fmt"
	"time"

	"personatrip/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQL 实现用户数据存储，基于GORM
type MySQL struct {
	DB *gorm.DB // 公开的GORM数据库连接，可以被其他仓库使用
}

// NewMySQL 创建新的MySQL存储实例
func NewMySQL(dsn string) (*MySQL, error) {
	// 使用GORM连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 获取并配置底层的SQL DB连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库表结构
	if err := autoMigrate(db); err != nil {
		return nil, err
	}

	return &MySQL{DB: db}, nil
}

// 自动迁移数据库表结构
func autoMigrate(db *gorm.DB) error {
	// 自动迁移用户和管理员表
	err := db.AutoMigrate(
		&models.UserMySQL{},
		&models.Admin{},
		&models.ModelConfig{},
	)
	return err
}

// Close 关闭数据库连接
func (m *MySQL) Close() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CreateUser 创建新用户
func (m *MySQL) CreateUser(ctx *gin.Context, user *models.UserMySQL) (*models.UserMySQL, error) {
	// 检查用户名是否已存在
	var count int64
	if err := m.DB.Model(&models.UserMySQL{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("username already exists")
	}

	// 检查邮箱是否已存在
	if err := m.DB.Model(&models.UserMySQL{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 设置加密后的密码
	user.Password = string(hashedPassword)

	// 使用GORM创建用户
	if err := m.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 返回前清除密码
	userCopy := *user
	userCopy.Password = ""
	return &userCopy, nil
}

// GetUserByID 通过ID获取用户
func (m *MySQL) GetUserByID(ctx *gin.Context, userID string) (*models.UserMySQL, error) {
	var user models.UserMySQL

	result := m.DB.Select("id, user_id, username, email, created_at, updated_at").
		Where("user_id = ?", userID).
		First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

// GetUserByUsername 通过用户名获取用户
func (m *MySQL) GetUserByUsername(ctx *gin.Context, username string) (*models.UserMySQL, error) {
	var user models.UserMySQL

	result := m.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

// CheckUserCredentials 检查用户凭据
func (m *MySQL) CheckUserCredentials(ctx *gin.Context, username, password string) (*models.UserMySQL, error) {
	// 获取用户（包括密码）
	user, err := m.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 清除密码
	user.Password = ""
	return user, nil
}
