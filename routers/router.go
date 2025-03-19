package routers

import (
	"traveler_agent/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	// 如果是生产环境，设置为生产模式
	// gin.SetMode(gin.ReleaseMode)

	r := gin.New() // 创建不含中间件的路由

	// 添加自定义中间件
	r.Use(middlewares.Logger())
	r.Use(gin.Recovery()) // 内置恢复中间件
	r.Use(middlewares.ErrorHandler())
	r.Use(middlewares.CORS())

	// API 版本
	v1 := r.Group("/api/v1")
	{
		// 健康检查路由
		v1.GET("/health", HealthCheck)

		// 基础测试路由
		v1.GET("/ping", Ping)

		// 用户相关路由组
		_ = v1.Group("/users")
		{
			// 这里可以添加用户相关路由
			// user.POST("/register", controllers.RegisterUser)
			// user.POST("/login", controllers.LoginUser)
			// 需要验证的路由
			// authorized := user.Group("/")
			// authorized.Use(middlewares.JWTAuth())
			// {
			//     authorized.GET("/profile", controllers.GetUserProfile)
			// }
		}

		// 景点相关路由组
		_ = v1.Group("/attractions")
		{
			// 这里可以添加景点相关路由
			// attraction.GET("/", controllers.ListAttractions)
			// attraction.GET("/:id", controllers.GetAttraction)
		}

		// 行程相关路由组
		_ = v1.Group("/itineraries")
		{
			// 这里可以添加行程相关路由
			// itinerary.POST("/", controllers.CreateItinerary)
			// itinerary.GET("/", controllers.ListItineraries)
			// itinerary.GET("/:id", controllers.GetItinerary)
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
