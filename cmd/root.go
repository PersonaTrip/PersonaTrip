package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"personatrip/internal/app"
	"personatrip/internal/config"
	"personatrip/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Execute 启动服务器并处理优雅关闭
func Execute() error {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化应用程序
	application := app.New()

	// 添加CORS中间件，允许所有跨域请求
	application.Router.Use(middleware.CORS())

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: application.Router,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("Server starting on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
	return nil
}
