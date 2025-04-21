package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Environment    string
	Port           string
	MongoURI       string
	MySQLDSN       string
	JWTSecret      string
	AutoMigrate    bool   // 是否自动迁移数据库
	CreateSuperAdmin bool   // 是否创建超级管理员
	SuperAdminUsername string // 超级管理员用户名
	SuperAdminPassword string // 超级管理员密码
	SuperAdminEmail    string // 超级管理员邮箱
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	// 加载.env文件，如果存在
	_ = godotenv.Load()

	// 设置默认值
	cfg := &Config{
		Environment: getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017/personatrip"),
		MySQLDSN:    getEnv("MYSQL_DSN", "root:password@tcp(localhost:3306)/personatrip?parseTime=true"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		AutoMigrate: getEnvBool("AUTO_MIGRATE", true),
		CreateSuperAdmin: getEnvBool("CREATE_SUPER_ADMIN", true),
		SuperAdminUsername: getEnv("SUPER_ADMIN_USERNAME", "admin"),
		SuperAdminPassword: getEnv("SUPER_ADMIN_PASSWORD", "admin123"),
		SuperAdminEmail:    getEnv("SUPER_ADMIN_EMAIL", "admin@personatrip.com"),
	}

	return cfg, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvBool 获取布尔类型的环境变量
func getEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}
