package controllers

import (
	"net/http"
	"strconv"
	"time"

	"traveler_agent/middlewares"
	"traveler_agent/models"
	"traveler_agent/repositories"
	"traveler_agent/services"

	"github.com/gin-gonic/gin"
)

// 全局行程服务
var (
	itineraryService services.ItineraryService
)

// 初始化行程服务
func InitItineraryController() {
	itineraryRepository := repositories.NewItineraryRepository(models.DB)
	itineraryService = services.NewItineraryService(itineraryRepository)
}

// ListPublicItineraries 获取公开行程列表
func ListPublicItineraries(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 获取公开行程
	itineraries, err := itineraryService.GetPublicItineraries(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// 假设我们获取总数为结果的长度（实际应该从服务层返回）
	total := len(itineraries)

	// 计算总页数
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   itineraries,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": totalPages,
		},
	})
}

// CreateItinerary 创建行程
func CreateItinerary(c *gin.Context) {
	// 获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "未授权访问",
		})
		return
	}

	var request models.CreateItineraryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "无效的请求数据",
		})
		return
	}

	// 创建行程
	itinerary, err := itineraryService.CreateItinerary(userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "行程创建成功",
		"data":    itinerary,
	})
}

// ListUserItineraries 获取用户行程列表
func ListUserItineraries(c *gin.Context) {
	// 获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "未授权访问",
		})
		return
	}

	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 获取用户行程
	itineraries, err := itineraryService.GetUserItineraries(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// 假设我们获取总数为结果的长度（实际应该从服务层返回）
	total := len(itineraries)

	// 计算总页数
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   itineraries,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": totalPages,
		},
	})
}

// GetItinerary 获取行程详情
func GetItinerary(c *gin.Context) {
	// 解析路径参数
	itineraryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "无效的行程ID",
		})
		return
	}

	// 获取行程详情
	itinerary, err := itineraryService.GetItineraryByID(itineraryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "行程不存在或无权访问",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   itinerary,
	})
}

// UpdateItinerary 更新行程
func UpdateItinerary(c *gin.Context) {
	// 获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "未授权访问",
		})
		return
	}

	// 解析路径参数
	itineraryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "无效的行程ID",
		})
		return
	}

	var request models.CreateItineraryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "无效的请求数据",
		})
		return
	}

	// 更新行程
	itinerary, err := itineraryService.UpdateItinerary(itineraryID, userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "行程更新成功",
		"data":    itinerary,
	})
}

// DeleteItinerary 删除行程
func DeleteItinerary(c *gin.Context) {
	// 获取用户ID
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "未授权访问",
		})
		return
	}

	// 解析路径参数
	itineraryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "无效的行程ID",
		})
		return
	}

	// 删除行程
	err = itineraryService.DeleteItinerary(itineraryID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "行程删除成功",
	})
}

// AddItineraryItem 添加行程项目
func AddItineraryItem(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "error",
		"error":  "添加行程项目功能尚未实现",
	})
}

// UpdateItineraryItem 更新行程项目
func UpdateItineraryItem(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "error",
		"error":  "更新行程项目功能尚未实现",
	})
}

// DeleteItineraryItem 删除行程项目
func DeleteItineraryItem(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "error",
		"error":  "删除行程项目功能尚未实现",
	})
}

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
