package api

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"personatrip/internal/config"
	"personatrip/internal/handlers"
	"personatrip/internal/middleware"
	"personatrip/internal/repository"
	"personatrip/internal/services"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(router *gin.Engine) {
	// 加载配置
	cfg, _ := config.Load()

	// 初始化服务和存储
	einoService := services.NewEinoService(cfg.EinoAPIKey)
	
	// 初始化MongoDB存储层
	var repo handlers.TripRepository
	
	// 尝试连接MongoDB
	mongoDB, err := repository.NewMongoDB(cfg.MongoURI)
	if err != nil {
		// 如果MongoDB连接失败，使用内存存储
		log.Printf("Failed to connect to MongoDB: %v, using in-memory storage instead", err)
		repo = repository.NewMemoryStore()
	} else {
		repo = mongoDB
	}
	
	// 初始化MySQL存储层
	mysqlDB, err := repository.NewMySQL(cfg.MySQLDSN)
	if err != nil {
		log.Printf("Failed to connect to MySQL: %v", err)
		// 如果无法连接MySQL，我们将无法提供用户认证
		// 在实际生产环境中，这里应该处理得更优雅
	}
	
	// 初始化认证服务
	authService := services.NewAuthService(mysqlDB, cfg.JWTSecret)
	
	// 初始化处理程序
	tripHandler := handlers.NewTripHandler(einoService, repo)
	authHandler := handlers.NewAuthHandler(authService)
	
	// 中间件
	authMiddleware := middleware.AuthMiddleware(authService)
	
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
	
	// Swagger文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
