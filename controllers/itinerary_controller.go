package controllers

import (
	"net/http"
	"strconv"
	"time"

	"traveler_agent/middlewares"
	"traveler_agent/models"
	"traveler_agent/services"

	"github.com/gin-gonic/gin"
)

// ItineraryController 行程控制器
type ItineraryController struct {
	itineraryService services.ItineraryService
}

// NewItineraryController 创建行程控制器
func NewItineraryController(itineraryService services.ItineraryService) *ItineraryController {
	return &ItineraryController{
		itineraryService: itineraryService,
	}
}

// CreateItinerary 创建行程
func (c *ItineraryController) CreateItinerary(ctx *gin.Context) {
	var request models.CreateItineraryRequest

	// 绑定JSON请求
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 获取用户ID
	userID, exists := middlewares.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 创建行程
	response, err := c.itineraryService.CreateItinerary(userID, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"itinerary": response})
}

// GetItinerary 获取行程详情
func (c *ItineraryController) GetItinerary(ctx *gin.Context) {
	// 解析路径参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的行程ID"})
		return
	}

	// 获取行程
	itinerary, err := c.itineraryService.GetItineraryByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, exists := middlewares.GetUserID(ctx)
	// 检查是否有权限访问
	if !itinerary.IsPublic && (!exists || itinerary.UserID != userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权访问此行程"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"itinerary": itinerary})
}

// GetUserItineraries 获取用户行程列表
func (c *ItineraryController) GetUserItineraries(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := middlewares.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	// 获取行程列表
	itineraries, err := c.itineraryService.GetUserItineraries(userID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"itineraries": itineraries,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateItinerary 更新行程
func (c *ItineraryController) UpdateItinerary(ctx *gin.Context) {
	// 解析路径参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的行程ID"})
		return
	}

	// 获取用户ID
	userID, exists := middlewares.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var request models.CreateItineraryRequest
	// 绑定JSON请求
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 更新行程
	response, err := c.itineraryService.UpdateItinerary(id, userID, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"itinerary": response})
}

// DeleteItinerary 删除行程
func (c *ItineraryController) DeleteItinerary(ctx *gin.Context) {
	// 解析路径参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的行程ID"})
		return
	}

	// 获取用户ID
	userID, exists := middlewares.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 删除行程
	err = c.itineraryService.DeleteItinerary(id, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "行程已删除"})
}

// AddItineraryDay 添加行程天数
func (c *ItineraryController) AddItineraryDay(ctx *gin.Context) {
	// 解析路径参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的行程ID"})
		return
	}

	// 获取用户ID
	userID, exists := middlewares.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 解析请求数据
	type AddDayRequest struct {
		Date string `json:"date" binding:"required"`
		Note string `json:"note"`
	}
	var request AddDayRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 解析日期
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的日期格式"})
		return
	}

	// 添加天数
	day, err := c.itineraryService.AddDay(id, userID, date, request.Note)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"day": day})
}

// GetPublicItineraries 获取公开行程
func (c *ItineraryController) GetPublicItineraries(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	// 获取公开行程
	itineraries, err := c.itineraryService.GetPublicItineraries(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"itineraries": itineraries,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// RegisterRoutes 注册路由
func (c *ItineraryController) RegisterRoutes(router *gin.RouterGroup) {
	itineraries := router.Group("/itineraries")
	{
		// 公开路由
		itineraries.GET("/public", c.GetPublicItineraries)

		// 需要认证的路由
		auth := itineraries.Group("")
		auth.Use(middlewares.JWTAuth())
		{
			auth.POST("", c.CreateItinerary)
			auth.GET("", c.GetUserItineraries)
			auth.GET("/:id", c.GetItinerary)
			auth.PUT("/:id", c.UpdateItinerary)
			auth.DELETE("/:id", c.DeleteItinerary)
			auth.POST("/:id/days", c.AddItineraryDay)
		}
	}
}
