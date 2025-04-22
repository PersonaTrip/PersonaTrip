package app

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"personatrip/internal/api"
	"personatrip/internal/config"
	"personatrip/internal/handlers"
	"personatrip/internal/middleware"
	"personatrip/internal/models"
	"personatrip/internal/repository"
	"personatrip/internal/services"
)

// Application 表示应用程序实例
type Application struct {
	Router       *gin.Engine
	Cfg          *config.Config
	Repositories *Repositories
	Services     *Services
	Handlers     *Handlers
}

// Repositories 包含所有仓库实例
type Repositories struct {
	MySQL           *repository.MySQL
	AdminRepo       repository.AdminRepository
	ModelConfigRepo repository.ModelConfigRepository
	TripRepo        handlers.TripRepository
}

// Services 包含所有服务实例
type Services struct {
	AuthService        *services.AuthService
	AdminService       services.AdminService
	ModelConfigService services.ModelConfigService
	EinoService        handlers.EinoServiceInterface
}

// Handlers 包含所有处理程序实例
type Handlers struct {
	AuthHandler        *handlers.AuthHandler
	AdminHandler       *handlers.AdminHandler
	ModelConfigHandler *handlers.ModelConfigHandler
	TripHandler        *handlers.TripHandler
}

// New 创建并初始化一个新的应用实例
func New() *Application {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建应用实例
	app := &Application{
		Router: gin.Default(),
		Cfg:    cfg,
	}

	// 初始化仓库、服务和处理程序
	app.initRepositories()
	app.initServices()
	app.initHandlers()
	app.setupRoutes()

	// 创建超级管理员（如果配置了）
	app.createSuperAdminIfNeeded()

	return app
}

// initRepositories 初始化所有仓库
func (a *Application) initRepositories() {
	a.Repositories = &Repositories{}

	// 初始化MySQL存储层
	mysqlDB, err := repository.NewMySQL(a.Cfg.MySQLDSN)
	if err != nil {
		log.Printf("Failed to connect to MySQL: %v", err)
		// 如果无法连接MySQL，我们将无法提供用户认证和管理员功能
		// 在实际生产环境中，这里应该处理得更优雅
	}
	a.Repositories.MySQL = mysqlDB

	// 初始化仓库
	// 使用GORM数据库存储管理员和模型配置
	a.Repositories.AdminRepo = repository.NewGormAdminRepository(mysqlDB.DB)
	a.Repositories.ModelConfigRepo = repository.NewGormModelConfigRepository(mysqlDB.DB)

	// 初始化MongoDB存储层或内存存储
	log.Printf(a.Cfg.MongoURI)
	mongoDB, err := repository.NewMongoDB(a.Cfg.MongoURI)
	if err != nil {
		// 如果MongoDB连接失败，使用内存存储
		log.Printf("Failed to connect to MongoDB: %v, using in-memory storage instead", err)
		a.Repositories.TripRepo = repository.NewMemoryStore()
	} else {
		a.Repositories.TripRepo = mongoDB
	}
}

// initServices 初始化所有服务
func (a *Application) initServices() {
	a.Services = &Services{
		AuthService:        services.NewAuthService(a.Repositories.MySQL, a.Cfg.JWTSecret),
		AdminService:       services.NewAdminService(a.Repositories.AdminRepo, a.Cfg.JWTSecret),
		ModelConfigService: services.NewModelConfigService(a.Repositories.ModelConfigRepo),
	}

	// 初始化Eino服务
	a.Services.EinoService = services.NewEinoService(a.Services.ModelConfigService)
}

// initHandlers 初始化所有处理程序
func (a *Application) initHandlers() {
	a.Handlers = &Handlers{
		AuthHandler:        handlers.NewAuthHandler(a.Services.AuthService),
		AdminHandler:       handlers.NewAdminHandler(a.Services.AdminService),
		ModelConfigHandler: handlers.NewModelConfigHandler(a.Services.ModelConfigService, a.Services.EinoService),
		TripHandler:        handlers.NewTripHandler(a.Services.EinoService, a.Repositories.TripRepo),
	}
}

// setupRoutes 设置路由
func (a *Application) setupRoutes() {
	// 获取中间件
	authMiddleware := middleware.AuthMiddleware(a.Services.AuthService)

	// 设置路由
	api.SetupRoutes(
		a.Router,
		a.Handlers.AuthHandler,
		a.Handlers.TripHandler,
		a.Handlers.AdminHandler,
		a.Handlers.ModelConfigHandler,
		authMiddleware,
		a.Cfg.JWTSecret,
	)
}

// createSuperAdminIfNeeded 如果配置了，创建超级管理员
func (a *Application) createSuperAdminIfNeeded() {
	if a.Cfg.CreateSuperAdmin {
		log.Println("Creating super admin if not exists...")
		// 检查超级管理员是否存在
		_, err := a.Services.AdminService.GetAdminByUsername(context.Background(), a.Cfg.SuperAdminUsername)
		if err != nil {
			// 创建超级管理员
			_, err := a.Services.AdminService.CreateAdmin(context.Background(), &models.AdminCreateRequest{
				Username: a.Cfg.SuperAdminUsername,
				Password: a.Cfg.SuperAdminPassword,
				Email:    a.Cfg.SuperAdminEmail,
				Role:     "super_admin",
			})
			if err != nil {
				log.Printf("Failed to create super admin: %v", err)
			} else {
				log.Println("Super admin created successfully")
			}
		}
	}
}

// Run 启动应用程序
func (a *Application) Run() error {
	return a.Router.Run(a.Cfg.ServerAddress)
}
