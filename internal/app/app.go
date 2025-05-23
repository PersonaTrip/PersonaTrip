package app

import (
	"context"

	"personatrip/internal/api"
	"personatrip/internal/config"
	"personatrip/internal/handlers"
	"personatrip/internal/middleware"
	"personatrip/internal/models"
	"personatrip/internal/repository"
	"personatrip/internal/services"
	"personatrip/internal/utils/logger"

	"github.com/gin-gonic/gin"
)

// Application 表示应用程序实例
type Application struct {
	Router       *gin.Engine
	Cfg          *config.Config
	DB           repository.Database
	Repositories *Repositories
	Services     *Services
	Handlers     *Handlers
}

// Repositories 包含所有仓库实例
type Repositories struct {
	TripRepo handlers.TripRepository
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
func New() (*Application, error) {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		logger.Errorf("Failed to load config: %v", err)
		return nil, err
	}
	// 创建应用实例
	app := &Application{
		Router: gin.Default(),
		Cfg:    cfg,
	}
	// 初始化数据库、仓库、服务和处理程序
	err = app.initDatabase()
	if err != nil {
		return nil, err
	}
	err = app.initRepositories()
	if err != nil {
		return nil, err
	}
	app.initServices()
	app.initHandlers()
	app.setupRoutes()

	// 创建超级管理员（如果配置了）
	err = app.createSuperAdminIfNeeded()
	if err != nil {
		return nil, err
	}
	// 创建默认的大模型配置
	app.createDefaultModelConfigIfNeeded()

	return app, nil
}

// initDatabase 初始化数据库
func (a *Application) initDatabase() error {
	// 初始化MySQL存储层
	mysqlDB, err := repository.NewMySQL(a.Cfg.MySQLDSN)
	if err != nil {
		logger.Errorf("Failed to connect to MySQL: %v", err)
		return err
	}

	// 创建数据库抽象层
	a.DB = repository.NewGormDatabase(mysqlDB.DB)
	return nil
}

// initRepositories 初始化所有仓库
func (a *Application) initRepositories() error {
	a.Repositories = &Repositories{}

	// 初始化MongoDB存储层或内存存储
	mongoDB, err := repository.NewMongoDB(a.Cfg.MongoURI)
	if err != nil {
		// 如果MongoDB连接失败，使用内存存储
		logger.Errorf("Failed to connect to MongoDB: %v, using in-memory storage instead", err)
		return err
	} else {
		a.Repositories.TripRepo = mongoDB
	}
	return nil
}

// initServices 初始化所有服务
func (a *Application) initServices() {
	a.Services = &Services{
		AuthService:        services.NewAuthService(a.DB, a.Cfg.JWTSecret),
		AdminService:       services.NewAdminService(a.DB, a.Cfg.JWTSecret),
		ModelConfigService: services.NewModelConfigService(a.DB),
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
func (a *Application) createSuperAdminIfNeeded() error {
	if a.Cfg.CreateSuperAdmin {
		logger.Info("Creating super admin if not exists...")
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
				logger.Errorf("Failed to create super admin: %v", err)
				return err
			} else {
				logger.Info("Super admin created successfully")
			}
		}
	}
	return nil
}

// createDefaultModelConfigIfNeeded 如果不存在，创建默认的大模型配置
func (a *Application) createDefaultModelConfigIfNeeded() {
	logger.Info("正在检查默认大模型配置...")

	// 尝试获取配置列表
	configs, err := a.Services.ModelConfigService.GetAllModelConfigs(context.Background())
	if err != nil {
		logger.Errorf("获取模型配置失败: %v", err)
		return
	}

	// 如果没有任何配置，创建默认配置
	if len(configs) == 0 {
		logger.Info("创建默认大模型配置...")
		// 创建默认的ARK大模型配置
		defaultConfig := &models.ModelConfigCreateRequest{
			Name:        "默认ARK配置",
			ModelType:   "ark",
			ModelName:   "",
			ApiKey:      "",
			BaseUrl:     "https://ark.cn-beijing.volces.com/api/v3",
			IsActive:    true,
			Temperature: 0.7,
			MaxTokens:   2000,
		}

		config, err := a.Services.ModelConfigService.CreateModelConfig(context.Background(), defaultConfig)
		if err != nil {
			logger.Errorf("创建默认大模型配置失败: %v", err)
		} else {
			logger.Infof("默认大模型配置创建成功，ID: %d", config.ID)
		}
	} else {
		logger.Info("已存在大模型配置，跳过创建默认配置")
	}
}

// Run 启动应用程序
func (a *Application) Run() error {
	return a.Router.Run(a.Cfg.ServerAddress)
}
