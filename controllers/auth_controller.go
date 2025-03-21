package controllers

import (
	"net/http"

	"traveler_agent/middlewares"
	"traveler_agent/models"
	"traveler_agent/repositories"
	"traveler_agent/services"

	"github.com/gin-gonic/gin"
)

// 全局认证服务 - 供路由使用
var (
	authService services.AuthService
)

// InitAuthController 初始化认证控制器
func InitAuthController() {
	userRepo := repositories.NewUserRepository(models.DB)
	authService = services.NewAuthService(userRepo)
}

// AuthController 认证控制器
type AuthController struct {
	authService services.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register 用户注册
func Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "无效的请求数据",
		})
		return
	}

	// 注册新用户
	user, token, err := authService.Register(c, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "无效的登录信息",
		})
		return
	}

	// 处理登录请求
	user, token, err := authService.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})
}

// GetProfile 获取用户个人资料
func (c *AuthController) GetProfile(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := middlewares.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取用户信息
	user, err := c.authService.GetUserByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// RegisterRoutes 注册路由
func (c *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", Register)
		auth.POST("/login", Login)
	}

	profile := router.Group("/profile")
	profile.Use(middlewares.JWTAuth())
	{
		profile.GET("", c.GetProfile)
	}
}
