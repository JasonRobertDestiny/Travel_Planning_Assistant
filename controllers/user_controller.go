package controllers

import (
	"net/http"

	"traveler_agent/middlewares"
	"traveler_agent/models"
	"traveler_agent/repositories"
	"traveler_agent/services"

	"github.com/gin-gonic/gin"
)

// 全局用户服务 - 供路由使用
var (
	userService services.UserService
)

// 初始化UserService
func InitUserController() {
	// 创建用户仓库并初始化用户服务
	userRepo := repositories.NewUserRepository(models.DB)
	userService = services.NewUserService(userRepo)
}

// GetUserProfile 获取用户个人资料
func GetUserProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取用户信息
	user, err := userService.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

// UpdateUserProfile 更新用户个人资料
func UpdateUserProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var request models.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 创建用户模型
	user := &models.User{
		Username:  request.Username,
		Email:     request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Phone:     request.Phone,
		Avatar:    request.Avatar,
	}

	// 更新用户信息
	if err := userService.UpdateUserProfile(c, userID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的用户信息
	updatedUser, _ := userService.GetUserByID(c, userID)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "用户资料更新成功",
		"data":    updatedUser,
	})
}

// UpdatePassword 更新用户密码
func UpdatePassword(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var request models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 更新密码
	err := userService.UpdatePassword(c, userID, request.OldPassword, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "密码更新成功",
	})
}

// GetUserPreferences 获取用户偏好设置
func GetUserPreferences(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取用户偏好
	preferences, err := userService.GetUserPreferences(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   preferences,
	})
}

// SaveUserPreferences 保存用户偏好设置
func SaveUserPreferences(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var preferences models.UserPreference
	if err := c.ShouldBindJSON(&preferences); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 保存用户偏好
	err := userService.SaveUserPreferences(c, userID, &preferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "用户偏好保存成功",
	})
}

// GetTravelPreferences 获取旅行偏好设置
func GetTravelPreferences(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取旅行偏好
	preferences, err := userService.GetTravelPreferences(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   preferences,
	})
}

// SaveTravelPreferences 保存旅行偏好设置
func SaveTravelPreferences(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var preferences models.UserPreference
	if err := c.ShouldBindJSON(&preferences); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 保存旅行偏好
	err := userService.SaveTravelPreferences(c, userID, &preferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "旅行偏好保存成功",
	})
}
