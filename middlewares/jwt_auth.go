package middlewares

import (
	"net/http"
	"strings"

	"traveler_agent/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查令牌格式，通常是"Bearer {token}"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式无效"})
			c.Abort()
			return
		}

		// 获取令牌部分
		tokenString := parts[1]

		// 验证令牌
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			if err == utils.ErrExpiredToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌已过期"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌无效"})
			}
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	id, ok := userID.(int64)
	return id, ok
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}

	name, ok := username.(string)
	return name, ok
}
