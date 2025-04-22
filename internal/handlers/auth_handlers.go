package handlers

import (
	"personatrip/internal/models"
	"personatrip/internal/services"
	"personatrip/internal/utils/httputil"

	"github.com/gin-gonic/gin"
)

// AuthHandler 处理认证相关的请求
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler 创建新的认证处理程序
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册信息"
// @Success 201 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnBadRequest(c, "无效的请求格式")
		return
	}

	// 注册用户
	user, err := h.authService.Register(c, &req)
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	httputil.ReturnCreated(c, "用户注册成功", user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录并获取JWT令牌
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录信息"
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 401 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnBadRequest(c, "无效的请求格式")
		return
	}

	// 登录用户
	response, err := h.authService.Login(c, &req)
	if err != nil {
		httputil.ReturnUnauthorized(c, err.Error())
		return
	}

	httputil.ReturnSuccessWithBean(c, "登录成功", response)
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的资料
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.ApiResponse
// @Failure 401 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /api/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID := c.GetString("user_id")
	if userID == "" {
		httputil.ReturnUnauthorized(c, "用户未认证")
		return
	}

	// 获取用户资料
	user, err := h.authService.GetUserFromToken(c, c.GetHeader("Authorization")[7:]) // 移除"Bearer "前缀
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	// 确保获取的用户与令牌中的用户匹配
	if userID != user.UserID {
		httputil.ReturnForbidden(c, "无权访问用户资料")
		return
	}

	httputil.ReturnSuccessWithBean(c, "获取用户资料成功", user)
}
