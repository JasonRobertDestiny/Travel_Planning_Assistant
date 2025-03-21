package routers

import (
	"traveler_agent/controllers"
	"traveler_agent/middlewares"
	"traveler_agent/utils"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	config := utils.GetConfig()

	// 如果是生产环境，设置为生产模式
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New() // 创建不含中间件的路由

	// 添加自定义中间件
	r.Use(utils.LoggerMiddleware())
	r.Use(gin.Recovery()) // 内置恢复中间件
	r.Use(middlewares.CORSMiddleware())

	// API 版本
	v1 := r.Group("/api/v1")
	{
		// 健康检查路由
		v1.GET("/health", HealthCheck)

		// 基础测试路由
		v1.GET("/ping", Ping)

		// 认证相关路由
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// 用户相关路由组
		user := v1.Group("/users")
		{
			// 需要验证的路由
			authorized := user.Group("/")
			authorized.Use(middlewares.JWTAuth())
			{
				authorized.GET("/profile", controllers.GetUserProfile)
				authorized.PUT("/profile", controllers.UpdateUserProfile)
				authorized.PUT("/password", controllers.UpdatePassword)
				authorized.GET("/preferences", controllers.GetUserPreferences)
				authorized.POST("/preferences", controllers.SaveUserPreferences)
			}
		}

		// 景点相关路由组
		attraction := v1.Group("/attractions")
		{
			attraction.GET("/", controllers.ListAttractions)
			attraction.GET("/:id", controllers.GetAttraction)

			// 需要验证的景点管理路由
			adminAttractions := attraction.Group("/")
			adminAttractions.Use(middlewares.JWTAuth())
			{
				adminAttractions.POST("/", controllers.CreateAttraction)
				adminAttractions.PUT("/:id", controllers.UpdateAttraction)
				adminAttractions.DELETE("/:id", controllers.DeleteAttraction)
			}
		}

		// 行程相关路由组
		itinerary := v1.Group("/itineraries")
		{
			// 公开行程
			itinerary.GET("/public", controllers.ListPublicItineraries)

			// 需要验证的行程管理路由
			userItineraries := itinerary.Group("/")
			userItineraries.Use(middlewares.JWTAuth())
			{
				userItineraries.POST("/", controllers.CreateItinerary)
				userItineraries.GET("/", controllers.ListUserItineraries)
				userItineraries.GET("/:id", controllers.GetItinerary)
				userItineraries.PUT("/:id", controllers.UpdateItinerary)
				userItineraries.DELETE("/:id", controllers.DeleteItinerary)

				// 行程项目管理
				userItineraries.POST("/:id/items", controllers.AddItineraryItem)
				userItineraries.PUT("/:id/items/:itemId", controllers.UpdateItineraryItem)
				userItineraries.DELETE("/:id/items/:itemId", controllers.DeleteItineraryItem)
			}
		}
	}

	return r
}

// HealthCheck 健康检查处理函数
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "服务运行正常",
	})
}

// Ping 测试处理函数
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
