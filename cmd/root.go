package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"personatrip/internal/app"
	"personatrip/internal/config"
	"personatrip/internal/utils/logger"

	"github.com/gin-gonic/gin"
)

// Execute 启动服务器并处理优雅关闭
func Execute() error {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// 初始化日志
	logPath := ""
	if cfg.LogConfig != nil {
		logger.SetLogLevel(cfg.LogConfig.Level)
		logPath = cfg.LogConfig.Path
	}
	if logPath != "" {
		logger.SetLogOutput(logPath)
	}

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化应用程序
	application, err := app.New()
	if err != nil {
		return err
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: application.Router,
	}

	// 在goroutine中启动服务器
	go func() {
		logger.Infof("Server starting on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")
	return nil
}
