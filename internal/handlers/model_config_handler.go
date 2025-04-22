package handlers

import (
	"net/http"
	"strconv"

	"personatrip/internal/models"
	"personatrip/internal/services"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := h.configService.CreateModelConfig(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果设置为活跃，刷新Eino服务的模型配置
	if config.IsActive {
		h.einoService.RefreshModelConfig(c.Request.Context())
	}

	c.JSON(http.StatusCreated, config.ToResponse())
}

// Update 更新模型配置
func (h *ModelConfigHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.ModelConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := h.configService.UpdateModelConfig(c.Request.Context(), uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果设置为活跃，刷新Eino服务的模型配置
	if config.IsActive {
		h.einoService.RefreshModelConfig(c.Request.Context())
	}

	c.JSON(http.StatusOK, config.ToResponse())
}

// Delete 删除模型配置
func (h *ModelConfigHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.configService.DeleteModelConfig(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "model configuration deleted successfully"})
}

// GetByID 根据ID获取模型配置
func (h *ModelConfigHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	config, err := h.configService.GetModelConfigByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config.ToResponse())
}

// GetAll 获取所有模型配置
func (h *ModelConfigHandler) GetAll(c *gin.Context) {
	configs, err := h.configService.GetAllModelConfigs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []models.ModelConfigResponse
	for _, config := range configs {
		response = append(response, config.ToResponse())
	}

	c.JSON(http.StatusOK, response)
}

// GetActive 获取当前活跃的模型配置
func (h *ModelConfigHandler) GetActive(c *gin.Context) {
	config, err := h.configService.GetActiveModelConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config.ToResponse())
}

// SetActive 设置指定ID的配置为活跃
func (h *ModelConfigHandler) SetActive(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.configService.SetActiveModelConfig(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 刷新Eino服务的模型配置
	h.einoService.RefreshModelConfig(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{"message": "model configuration activated successfully"})
}

// TestModel 测试模型配置
func (h *ModelConfigHandler) TestModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Prompt string `json:"prompt" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取指定ID的模型配置
	config, err := h.configService.GetModelConfigByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 创建临时的Eino客户端
	client := services.NewEinoServiceWithConfig(config)

	// 测试生成文本
	result, err := client.TestGenerateText(c.Request.Context(), req.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
