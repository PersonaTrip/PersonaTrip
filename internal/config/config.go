package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Environment string
	Port        string
	MongoURI    string
	MySQLDSN    string
	JWTSecret   string
	EinoAPIKey  string
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
		EinoAPIKey:  getEnv("EINO_API_KEY", ""),
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
