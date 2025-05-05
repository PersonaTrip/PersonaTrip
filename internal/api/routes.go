package api

import (
	"personatrip/internal/handlers"
	"personatrip/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes 注册所有API路由
func SetupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	tripHandler *handlers.TripHandler,
	adminHandler *handlers.AdminHandler,
	modelConfigHandler *handlers.ModelConfigHandler,
	authMiddleware gin.HandlerFunc,
	jwtSecret string,
) {
	router.Use(middleware.CORS())
	// API路由组
	api := router.Group("/api")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", authMiddleware, authHandler.GetProfile)
		}

		// 旅行计划相关路由
		trips := api.Group("/trips")
		{
			trips.POST("/generate", authMiddleware, tripHandler.GenerateTripPlan)
			trips.GET("/:id", tripHandler.GetTripPlan)
			trips.GET("/user", authMiddleware, tripHandler.GetUserTripPlans)
			trips.PUT("/:id", authMiddleware, tripHandler.UpdateTripPlan)
			trips.DELETE("/:id", authMiddleware, tripHandler.DeleteTripPlan)
		}

		// 推荐相关路由
		recommendations := api.Group("/recommendations")
		{
			recommendations.POST("/destinations", tripHandler.GenerateDestinationRecommendations)
		}
	}

	// 设置管理员路由
	SetupAdminRoutes(router, adminHandler, modelConfigHandler, jwtSecret)

	// Swagger文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
