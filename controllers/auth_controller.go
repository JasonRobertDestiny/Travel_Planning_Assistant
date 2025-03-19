package controllers

import (
	"net/http"

	"traveler_agent/middlewares"
	"traveler_agent/models"
	"traveler_agent/services"

	"github.com/gin-gonic/gin"
)

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

// Register 注册
// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var request models.CreateUserRequest

	// 绑定JSON请求
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证请求数据
	if request.Username == "" || request.Email == "" || request.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名、邮箱和密码必须提供"})
		return
	}

	// 调用服务层处理注册
	user, token, err := c.authService.Register(ctx, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回用户信息和令牌
	ctx.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var request models.LoginRequest

	// 绑定JSON请求
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证请求数据
	if request.Email == "" || request.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱和密码必须提供"})
		return
	}

	// 调用服务层处理登录
	user, token, err := c.authService.Login(ctx, request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 返回用户信息和令牌
	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
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
		auth.POST("/register", c.Register)
		auth.POST("/login", c.Login)
	}

	profile := router.Group("/profile")
	profile.Use(middlewares.JWTAuth())
	{
		profile.GET("", c.GetProfile)
	}
}
