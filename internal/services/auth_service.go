package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"personatrip/internal/models"
)

// AuthService 处理用户认证
type AuthService struct {
	userRepo   UserRepository
	jwtSecret  string
	expiration time.Duration
}

// UserRepository 用户存储接口
type UserRepository interface {
	CreateUser(ctx *gin.Context, user *models.UserMySQL) (*models.UserMySQL, error)
	GetUserByID(ctx *gin.Context, userID string) (*models.UserMySQL, error)
	GetUserByUsername(ctx *gin.Context, username string) (*models.UserMySQL, error)
	CheckUserCredentials(ctx *gin.Context, username, password string) (*models.UserMySQL, error)
}

// NewAuthService 创建新的认证服务
func NewAuthService(userRepo UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		expiration: 24 * time.Hour, // 令牌有效期24小时
	}
}

// Register 用户注册
func (s *AuthService) Register(ctx *gin.Context, req *models.RegisterRequest) (*models.UserMySQL, error) {
	// 创建用户
	userMongoID := primitive.NewObjectID()
	userID := userMongoID.Hex()
	user := &models.UserMySQL{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		UserID:   userID,
	}

	return s.userRepo.CreateUser(ctx, user)
}

// Login 用户登录
func (s *AuthService) Login(ctx *gin.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// 验证用户凭据
	user, err := s.userRepo.CheckUserCredentials(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	// 生成JWT令牌
	token, err := s.GenerateToken(user.UserID)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GenerateToken 生成JWT令牌
func (s *AuthService) GenerateToken(userID string) (string, error) {
	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.expiration).Unix(),
	})

	// 签名令牌
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken 验证JWT令牌
func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// 验证令牌有效性
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// 获取声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// 获取用户ID
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID in token")
	}

	return userID, nil
}

// GetUserFromToken 从令牌获取用户
func (s *AuthService) GetUserFromToken(ctx *gin.Context, tokenString string) (*models.UserMySQL, error) {
	// 验证令牌
	userID, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	// 获取用户
	return s.userRepo.GetUserByID(ctx, userID)
}
