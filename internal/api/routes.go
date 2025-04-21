package api

import (
	"context"
	"log"

	"personatrip/internal/config"
	"personatrip/internal/handlers"
	"personatrip/internal/middleware"
	"personatrip/internal/models"
	"personatrip/internal/repository"
	"personatrip/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(router *gin.Engine) {
	// 加载配置
	cfg, _ := config.Load()

	// 初始化MySQL存储层
	mysqlDB, err := repository.NewMySQL(cfg.MySQLDSN)
	if err != nil {
		log.Printf("Failed to connect to MySQL: %v", err)
		// 如果无法连接MySQL，我们将无法提供用户认证和管理员功能
		// 在实际生产环境中，这里应该处理得更优雅
	}

	// 自动迁移数据库
	if cfg.AutoMigrate && mysqlDB != nil {
		log.Println("Auto migrating database...")
		if err := mysqlDB.AutoMigrate(&models.Admin{}, &models.ModelConfig{}); err != nil {
			log.Printf("Failed to auto migrate: %v", err)
		}
	}

	// 初始化仓库
	// 使用MySQL数据库存储管理员和模型配置
	adminRepo := repository.NewSQLAdminRepository(mysqlDB.DB)
	modelConfigRepo := repository.NewSQLModelConfigRepository(mysqlDB.DB)

	// 初始化服务
	adminService := services.NewAdminService(adminRepo, cfg.JWTSecret)
	modelConfigService := services.NewModelConfigService(modelConfigRepo)

	// 创建超级管理员
	if cfg.CreateSuperAdmin {
		log.Println("Creating super admin if not exists...")
		// 检查超级管理员是否存在
		_, err := adminService.GetAdminByUsername(context.Background(), cfg.SuperAdminUsername)
		if err != nil {
			// 创建超级管理员
			_, err := adminService.CreateAdmin(context.Background(), &models.AdminCreateRequest{
				Username: cfg.SuperAdminUsername,
				Password: cfg.SuperAdminPassword,
				Email:    cfg.SuperAdminEmail,
				Role:     "super_admin",
			})
			if err != nil {
				log.Printf("Failed to create super admin: %v", err)
			} else {
				log.Println("Super admin created successfully")
			}
		}
	}

	// 初始化Eino服务
	einoService := services.NewEinoService(modelConfigService)

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

	// 初始化认证服务
	authService := services.NewAuthService(mysqlDB, cfg.JWTSecret)

	// 初始化处理程序
	tripHandler := handlers.NewTripHandler(einoService, repo)
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)
	modelConfigHandler := handlers.NewModelConfigHandler(modelConfigService, einoService)

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

	// 设置管理员路由
	SetupAdminRoutes(router, adminHandler, modelConfigHandler, cfg.JWTSecret)

	// Swagger文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
