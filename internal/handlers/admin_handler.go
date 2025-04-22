package handlers

import (
	"strconv"

	"personatrip/internal/models"
	"personatrip/internal/services"
	"personatrip/internal/utils/httputil"

	"github.com/gin-gonic/gin"
)

// AdminHandler 处理管理员相关的请求
type AdminHandler struct {
	adminService services.AdminService
}

// NewAdminHandler 创建新的管理员处理器
func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// Login 管理员登录
func (h *AdminHandler) Login(c *gin.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnBadRequest(c, err.Error())
		return
	}

	token, err := h.adminService.Login(c.Request.Context(), &req)
	if err != nil {
		httputil.ReturnUnauthorized(c, err.Error())
		return
	}

	httputil.ReturnSuccessWithData(c, "管理员登录成功", map[string]string{"token": token})
}

// Create 创建新管理员
func (h *AdminHandler) Create(c *gin.Context) {
	var req models.AdminCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnBadRequest(c, err.Error())
		return
	}

	admin, err := h.adminService.CreateAdmin(c.Request.Context(), &req)
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	httputil.ReturnCreated(c, "管理员创建成功", admin.ToResponse())
}

// Update 更新管理员信息
func (h *AdminHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httputil.ReturnBadRequest(c, "无效的ID")
		return
	}

	var req models.AdminUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnBadRequest(c, err.Error())
		return
	}

	admin, err := h.adminService.UpdateAdmin(c.Request.Context(), uint(id), &req)
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	httputil.ReturnSuccessWithBean(c, "管理员更新成功", admin.ToResponse())
}

// Delete 删除管理员
func (h *AdminHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httputil.ReturnBadRequest(c, "无效的ID")
		return
	}

	if err := h.adminService.DeleteAdmin(c.Request.Context(), uint(id)); err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	httputil.ReturnSuccess(c, "管理员删除成功")
}

// GetByID 根据ID获取管理员
func (h *AdminHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httputil.ReturnBadRequest(c, "无效的ID")
		return
	}

	admin, err := h.adminService.GetAdminByID(c.Request.Context(), uint(id))
	if err != nil {
		httputil.ReturnNotFound(c, err.Error())
		return
	}

	httputil.ReturnSuccessWithBean(c, "获取管理员成功", admin.ToResponse())
}

// GetAll 获取所有管理员
func (h *AdminHandler) GetAll(c *gin.Context) {
	admins, err := h.adminService.GetAllAdmins(c.Request.Context())
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	var response []models.AdminResponse
	for _, admin := range admins {
		response = append(response, admin.ToResponse())
	}

	httputil.ReturnSuccessWithList(c, "获取所有管理员成功", response)
}
