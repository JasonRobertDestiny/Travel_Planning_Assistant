package middlewares

import (
	"log"
	"net/http"
	"time"

	"traveler_agent/configs"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件，记录请求处理时间和状态
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 请求结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 获取状态码和错误信息
		status := c.Writer.Status()
		var errMsg string
		if len(c.Errors) > 0 {
			errMsg = c.Errors.String()
		}

		// 构建日志信息
		log.Printf("[GIN] %s | %3d | %13v | %15s | %-7s %s%s %s",
			end.Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			c.ClientIP(),
			c.Request.Method,
			path,
			func() string {
				if raw != "" {
					return "?" + raw
				}
				return ""
			}(),
			errMsg,
		)
	}
}

// ErrorHandler 错误处理中间件，统一处理错误响应
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果有错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors.Last()

			// 根据环境决定是否返回详细错误信息
			var message string
			if configs.IsDevelopment() {
				message = err.Error()
			} else {
				message = "服务器内部错误"
			}

			// 如果响应已经发送，不再处理
			if c.Writer.Written() {
				return
			}

			// 返回JSON错误响应
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": message,
			})
		}
	}
}

// CORS 跨域资源共享中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
