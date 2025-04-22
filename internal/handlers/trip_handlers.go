package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"personatrip/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EinoServiceInterface 定义Eino服务接口
type EinoServiceInterface interface {
	GenerateTripPlan(ctx context.Context, req *models.PlanRequest) (*models.TripPlan, error)
	GenerateDestinationRecommendations(ctx context.Context, preferences *models.UserPreferences) ([]string, error)
	TestGenerateText(ctx context.Context, prompt string) (string, error)
	RefreshModelConfig(ctx context.Context) error
}

// TripHandler 处理旅行相关的请求
type TripHandler struct {
	einoService EinoServiceInterface
	repository  TripRepository
}

// TripRepository 定义仓库接口
type TripRepository interface {
	CreateTripPlan(ctx *gin.Context, plan *models.TripPlan) (*models.TripPlan, error)
	GetTripPlanByID(ctx *gin.Context, id primitive.ObjectID) (*models.TripPlan, error)
	GetTripPlansByUserID(ctx *gin.Context, userID primitive.ObjectID) ([]*models.TripPlan, error)
	UpdateTripPlan(ctx *gin.Context, plan *models.TripPlan) error
	DeleteTripPlan(ctx *gin.Context, id primitive.ObjectID) error
}

// NewTripHandler 创建新的旅行处理程序
func NewTripHandler(einoService EinoServiceInterface, repository TripRepository) *TripHandler {
	return &TripHandler{
		einoService: einoService,
		repository:  repository,
	}
}

// GenerateTripPlan 生成旅行计划
// @Summary 生成AI旅行计划
// @Description 根据用户输入的偏好生成个性化旅行计划
// @Tags trips
// @Accept json
// @Produce json
// @Param request body models.PlanRequest true "旅行计划请求"
// @Success 200 {object} models.TripPlan
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/trips/generate [post]
func (h *TripHandler) GenerateTripPlan(c *gin.Context) {
	var req models.PlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 验证日期
	if req.StartDate.After(req.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date cannot be after end date"})
		return
	}

	// 调用Eino服务生成旅行计划
	plan, err := h.einoService.GenerateTripPlan(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate trip plan"})
		return
	}

	// 获取用户ID（假设从认证中间件中获取）
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		// 如果没有认证，可以使用一个默认ID或返回错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 设置用户ID和标题
	plan.UserID = userID
	if plan.Title == "" {
		plan.Title = req.Destination + " Trip " + time.Now().Format("2006-01-02")
	}

	// 保存到数据库
	savedPlan, err := h.repository.CreateTripPlan(c, plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save trip plan"})
		return
	}

	c.JSON(http.StatusOK, savedPlan)
}

// GetTripPlan 获取旅行计划
// @Summary 获取旅行计划
// @Description 通过ID获取旅行计划详情
// @Tags trips
// @Accept json
// @Produce json
// @Param id path string true "旅行计划ID"
// @Success 200 {object} models.TripPlan
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/trips/{id} [get]
func (h *TripHandler) GetTripPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	plan, err := h.repository.GetTripPlanByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip plan not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// GetUserTripPlans 获取用户的所有旅行计划
// @Summary 获取用户旅行计划
// @Description 获取当前用户的所有旅行计划
// @Tags trips
// @Accept json
// @Produce json
// @Success 200 {array} models.TripPlan
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/trips/user [get]
func (h *TripHandler) GetUserTripPlans(c *gin.Context) {
	// 获取用户ID（假设从认证中间件中获取）
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	fmt.Println(userID)
	plans, err := h.repository.GetTripPlansByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trip plans" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// UpdateTripPlan 更新旅行计划
// @Summary 更新旅行计划
// @Description 更新现有旅行计划
// @Tags trips
// @Accept json
// @Produce json
// @Param id path string true "旅行计划ID"
// @Param plan body models.TripPlan true "更新的旅行计划"
// @Success 200 {object} models.TripPlan
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/trips/{id} [put]
func (h *TripHandler) UpdateTripPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 检查计划是否存在
	existingPlan, err := h.repository.GetTripPlanByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip plan not found"})
		return
	}

	// 获取用户ID（假设从认证中间件中获取）
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 检查是否是用户自己的计划
	if existingPlan.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this plan"})
		return
	}

	// 解析请求体
	var updatedPlan models.TripPlan
	if err := c.ShouldBindJSON(&updatedPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 保持原始ID和用户ID
	updatedPlan.ID = id
	updatedPlan.UserID = userID

	// 更新计划
	if err := h.repository.UpdateTripPlan(c, &updatedPlan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trip plan"})
		return
	}

	c.JSON(http.StatusOK, updatedPlan)
}

// DeleteTripPlan 删除旅行计划
// @Summary 删除旅行计划
// @Description 删除现有旅行计划
// @Tags trips
// @Accept json
// @Produce json
// @Param id path string true "旅行计划ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/trips/{id} [delete]
func (h *TripHandler) DeleteTripPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 检查计划是否存在
	existingPlan, err := h.repository.GetTripPlanByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip plan not found"})
		return
	}

	// 获取用户ID（假设从认证中间件中获取）
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 检查是否是用户自己的计划
	if existingPlan.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this plan"})
		return
	}

	// 删除计划
	if err := h.repository.DeleteTripPlan(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trip plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trip plan deleted successfully"})
}

// GenerateDestinationRecommendations 生成目的地推荐
// @Summary 生成目的地推荐
// @Description 根据用户偏好生成目的地推荐
// @Tags recommendations
// @Accept json
// @Produce json
// @Param preferences body models.UserPreferences true "用户偏好"
// @Success 200 {array} string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/recommendations/destinations [post]
func (h *TripHandler) GenerateDestinationRecommendations(c *gin.Context) {
	var preferences models.UserPreferences
	if err := c.ShouldBindJSON(&preferences); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 调用Eino服务生成推荐
	recommendations, err := h.einoService.GenerateDestinationRecommendations(c.Request.Context(), &preferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recommendations"})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}
