package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AdminAuthMiddleware 是管理员认证中间件
func AdminAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// 检查令牌格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		// 解析令牌
		tokenString := parts[1]
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// 将管理员信息存储在上下文中
		c.Set("adminID", uint(claims["id"].(float64)))
		c.Set("adminUsername", claims["username"].(string))
		c.Set("adminRole", claims["role"].(string))

		c.Next()
	}
}

// RequireSuperAdmin 是超级管理员权限检查中间件
func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("adminRole")
		if !exists || role.(string) != "super_admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "super admin privileges required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
