package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"personatrip/internal/models"
	"personatrip/internal/repository"
)

// AdminService 定义管理员服务接口
type AdminService interface {
	CreateAdmin(ctx context.Context, req *models.AdminCreateRequest) (*models.Admin, error)
	UpdateAdmin(ctx context.Context, id uint, req *models.AdminUpdateRequest) (*models.Admin, error)
	DeleteAdmin(ctx context.Context, id uint) error
	GetAdminByID(ctx context.Context, id uint) (*models.Admin, error)
	GetAdminByUsername(ctx context.Context, username string) (*models.Admin, error)
	GetAllAdmins(ctx context.Context) ([]models.Admin, error)
	Login(ctx context.Context, req *models.AdminLoginRequest) (string, error)
}

// AdminServiceImpl 是管理员服务的实现
type AdminServiceImpl struct {
	repo      repository.AdminRepository
	jwtSecret string
}

// NewAdminService 创建新的管理员服务
func NewAdminService(repo repository.AdminRepository, jwtSecret string) AdminService {
	return &AdminServiceImpl{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// CreateAdmin 创建新的管理员
func (s *AdminServiceImpl) CreateAdmin(ctx context.Context, req *models.AdminCreateRequest) (*models.Admin, error) {
	// 检查用户名是否已存在
	_, err := s.repo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	admin := &models.Admin{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
	}

	// 如果没有指定角色，默认为admin
	if admin.Role == "" {
		admin.Role = "admin"
	}

	// 设置密码
	if err := admin.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, admin); err != nil {
		return nil, err
	}
	return admin, nil
}

// UpdateAdmin 更新管理员信息
func (s *AdminServiceImpl) UpdateAdmin(ctx context.Context, id uint, req *models.AdminUpdateRequest) (*models.Admin, error) {
	admin, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Email != "" {
		admin.Email = req.Email
	}
	if req.Role != "" {
		admin.Role = req.Role
	}
	if req.Password != "" {
		if err := admin.SetPassword(req.Password); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Update(ctx, admin); err != nil {
		return nil, err
	}
	return admin, nil
}

// DeleteAdmin 删除管理员
func (s *AdminServiceImpl) DeleteAdmin(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// GetAdminByID 根据ID获取管理员
func (s *AdminServiceImpl) GetAdminByID(ctx context.Context, id uint) (*models.Admin, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAdminByUsername 根据用户名获取管理员
func (s *AdminServiceImpl) GetAdminByUsername(ctx context.Context, username string) (*models.Admin, error) {
	return s.repo.GetByUsername(ctx, username)
}

// GetAllAdmins 获取所有管理员
func (s *AdminServiceImpl) GetAllAdmins(ctx context.Context) ([]models.Admin, error) {
	return s.repo.GetAll(ctx)
}

// Login 管理员登录
func (s *AdminServiceImpl) Login(ctx context.Context, req *models.AdminLoginRequest) (string, error) {
	admin, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// 验证密码
	if !admin.CheckPassword(req.Password) {
		return "", errors.New("invalid username or password")
	}

	// 生成JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       admin.ID,
		"username": admin.Username,
		"role":     admin.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
	})

	// 签名令牌
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
