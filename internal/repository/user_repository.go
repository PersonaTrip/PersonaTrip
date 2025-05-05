package repository

import (
	"personatrip/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	CreateUser(ctx *gin.Context, user *models.UserMySQL) (*models.UserMySQL, error)
	GetUserByID(ctx *gin.Context, userID string) (*models.UserMySQL, error)
	GetUserByUsername(ctx *gin.Context, username string) (*models.UserMySQL, error)
	CheckUserCredentials(ctx *gin.Context, username, password string) (*models.UserMySQL, error)
}

// GormUserRepository 是GORM实现的用户仓库
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository 创建新的GORM用户仓库
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

// CreateUser 实现UserRepository接口
func (r *GormUserRepository) CreateUser(ctx *gin.Context, user *models.UserMySQL) (*models.UserMySQL, error) {
	mysql := &MySQL{DB: r.db}
	return mysql.CreateUser(ctx, user)
}

// GetUserByID 实现UserRepository接口
func (r *GormUserRepository) GetUserByID(ctx *gin.Context, userID string) (*models.UserMySQL, error) {
	mysql := &MySQL{DB: r.db}
	return mysql.GetUserByID(ctx, userID)
}

// GetUserByUsername 实现UserRepository接口
func (r *GormUserRepository) GetUserByUsername(ctx *gin.Context, username string) (*models.UserMySQL, error) {
	mysql := &MySQL{DB: r.db}
	return mysql.GetUserByUsername(ctx, username)
}

// CheckUserCredentials 实现UserRepository接口
func (r *GormUserRepository) CheckUserCredentials(ctx *gin.Context, username, password string) (*models.UserMySQL, error) {
	mysql := &MySQL{DB: r.db}
	return mysql.CheckUserCredentials(ctx, username, password)
}
