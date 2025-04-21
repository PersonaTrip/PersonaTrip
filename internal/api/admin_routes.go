package api

import (
	"github.com/gin-gonic/gin"
	"personatrip/internal/handlers"
	"personatrip/internal/middleware"
)

// SetupAdminRoutes 设置管理员相关路由
func SetupAdminRoutes(router *gin.Engine, adminHandler *handlers.AdminHandler, modelConfigHandler *handlers.ModelConfigHandler, jwtSecret string) {
	// 管理员API组
	adminGroup := router.Group("/api/admin")

	// 公开路由
	adminGroup.POST("/login", adminHandler.Login)

	// 需要认证的路由
	authGroup := adminGroup.Group("")
	authGroup.Use(middleware.AdminAuthMiddleware(jwtSecret))

	// 管理员管理（仅超级管理员可访问）
	adminManage := authGroup.Group("/admins")
	adminManage.Use(middleware.RequireSuperAdmin())
	{
		adminManage.POST("", adminHandler.Create)
		adminManage.GET("", adminHandler.GetAll)
		adminManage.GET("/:id", adminHandler.GetByID)
		adminManage.PUT("/:id", adminHandler.Update)
		adminManage.DELETE("/:id", adminHandler.Delete)
	}

	// 模型配置管理
	modelGroup := authGroup.Group("/models")
	{
		modelGroup.POST("", modelConfigHandler.Create)
		modelGroup.GET("", modelConfigHandler.GetAll)
		modelGroup.GET("/active", modelConfigHandler.GetActive)
		modelGroup.GET("/:id", modelConfigHandler.GetByID)
		modelGroup.PUT("/:id", modelConfigHandler.Update)
		modelGroup.DELETE("/:id", modelConfigHandler.Delete)
		modelGroup.POST("/:id/activate", modelConfigHandler.SetActive)
		modelGroup.POST("/:id/test", modelConfigHandler.TestModel)
	}
}
