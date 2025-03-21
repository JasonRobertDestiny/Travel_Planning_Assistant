package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"traveler_agent/utils"
)

// CORSMiddleware 处理跨域请求
func CORSMiddleware() gin.HandlerFunc {
	config := utils.GetConfig()

	return cors.New(cors.Config{
		AllowOrigins:     config.CORS.AllowedOrigins,
		AllowMethods:     config.CORS.AllowedMethods,
		AllowHeaders:     config.CORS.AllowedHeaders,
		ExposeHeaders:    config.CORS.ExposeHeaders,
		AllowCredentials: config.CORS.AllowCredentials,
		MaxAge:           time.Duration(config.CORS.MaxAge) * time.Second,
		AllowWildcard:    true,
	})
}
