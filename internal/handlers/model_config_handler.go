package handlers

import (
	"strconv"

	"personatrip/internal/models"
	"personatrip/internal/services"
	"personatrip/internal/utils/httputil"

	"github.com/gin-gonic/gin"
)

// ModelConfigHandler 处理模型配置相关的请求
type ModelConfigHandler struct {
	configService services.ModelConfigService
	einoService   EinoServiceInterface
}

// NewModelConfigHandler 创建新的模型配置处理器
func NewModelConfigHandler(configService services.ModelConfigService, einoService EinoServiceInterface) *ModelConfigHandler {
	return &ModelConfigHandler{
		configService: configService,
		einoService:   einoService,
	}
}

// Create 创建新的模型配置
func (h *ModelConfigHandler) Create(c *gin.Context) {
	var req models.ModelConfigCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnBadRequest(c, err.Error())
		return
	}

	config, err := h.configService.CreateModelConfig(c.Request.Context(), &req)
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	// 如果设置为活跃，刷新Eino服务的模型配置
	if config.IsActive {
		h.einoService.RefreshModelConfig(c.Request.Context())
	}

	httputil.ReturnCreated(c, "模型配置创建成功", config.ToResponse())
}

// Update 更新模型配置
func (h *ModelConfigHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httputil.ReturnBadRequest(c, "无效的ID")
		return
	}

	var req models.ModelConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnBadRequest(c, err.Error())
		return
	}

	config, err := h.configService.UpdateModelConfig(c.Request.Context(), uint(id), &req)
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	// 如果设置为活跃，刷新Eino服务的模型配置
	if config.IsActive {
		h.einoService.RefreshModelConfig(c.Request.Context())
	}

	httputil.ReturnSuccessWithBean(c, "模型配置更新成功", config.ToResponse())
}

// Delete 删除模型配置
func (h *ModelConfigHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httputil.ReturnBadRequest(c, "无效的ID")
		return
	}

	if err := h.configService.DeleteModelConfig(c.Request.Context(), uint(id)); err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	httputil.ReturnSuccess(c, "模型配置删除成功")
}

// GetByID 根据ID获取模型配置
func (h *ModelConfigHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httputil.ReturnBadRequest(c, "无效的ID")
		return
	}

	config, err := h.configService.GetModelConfigByID(c.Request.Context(), uint(id))
	if err != nil {
		httputil.ReturnNotFound(c, err.Error())
		return
	}

	httputil.ReturnSuccessWithBean(c, "获取模型配置成功", config.ToResponse())
}

// GetAll 获取所有模型配置
func (h *ModelConfigHandler) GetAll(c *gin.Context) {
	configs, err := h.configService.GetAllModelConfigs(c.Request.Context())
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	var response []models.ModelConfigResponse
	for _, config := range configs {
		response = append(response, config.ToResponse())
	}

	httputil.ReturnSuccessWithList(c, "获取所有模型配置成功", response)
}

// GetActive 获取当前活跃的模型配置
func (h *ModelConfigHandler) GetActive(c *gin.Context) {
	config, err := h.configService.GetActiveModelConfig(c.Request.Context())
	if err != nil {
		httputil.ReturnNotFound(c, err.Error())
		return
	}

	httputil.ReturnSuccessWithBean(c, "获取当前活跃的模型配置成功", config.ToResponse())
}

// SetActive 设置指定ID的配置为活跃
func (h *ModelConfigHandler) SetActive(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httputil.ReturnBadRequest(c, "无效的ID")
		return
	}

	err = h.configService.SetActiveModelConfig(c.Request.Context(), uint(id))
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	// 刷新Eino服务的模型配置
	h.einoService.RefreshModelConfig(c.Request.Context())

	// 获取更新后的配置
	config, err := h.configService.GetModelConfigByID(c.Request.Context(), uint(id))
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	httputil.ReturnSuccessWithBean(c, "模型配置激活成功", config.ToResponse())
}

// TestModel 测试模型配置
func (h *ModelConfigHandler) TestModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httputil.ReturnBadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Prompt string `json:"prompt" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnBadRequest(c, err.Error())
		return
	}

	// 获取指定ID的模型配置
	config, err := h.configService.GetModelConfigByID(c.Request.Context(), uint(id))
	if err != nil {
		httputil.ReturnNotFound(c, err.Error())
		return
	}

	// 创建临时的Eino客户端
	client := services.NewEinoServiceWithConfig(config)

	// 测试生成文本
	result, err := client.TestGenerateText(c.Request.Context(), req.Prompt)
	if err != nil {
		httputil.ReturnInternalError(c, err.Error())
		return
	}

	httputil.ReturnSuccessWithData(c, "模型测试成功", map[string]string{"result": result})
}
