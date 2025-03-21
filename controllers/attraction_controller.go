package controllers

import (
	"net/http"
	"strconv"

	"traveler_agent/models"
	"traveler_agent/services"

	"traveler_agent/repositories"

	"github.com/gin-gonic/gin"
)

// 全局景点服务
var (
	attractionService services.AttractionService
)

// 初始化景点服务
func InitAttractionController() {
	// 创建景点仓库并初始化景点服务
	attractionRepo := repositories.NewAttractionRepository(models.DB)
	attractionService = services.NewAttractionService(attractionRepo)
}

// AttractionController 景点控制器
type AttractionController struct {
	attractionService services.AttractionService
}

// NewAttractionController 创建景点控制器
func NewAttractionController(attractionService services.AttractionService) *AttractionController {
	return &AttractionController{
		attractionService: attractionService,
	}
}

// GetAttractionByID 获取景点详情
func (c *AttractionController) GetAttractionByID(ctx *gin.Context) {
	// 解析路径参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的景点ID"})
		return
	}

	// 获取景点信息
	attraction, err := c.attractionService.GetAttractionByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"attraction": attraction})
}

// SearchAttractions 搜索景点
func (c *AttractionController) SearchAttractions(ctx *gin.Context) {
	var request models.SearchAttractionsRequest

	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的查询参数"})
		return
	}

	// 设置默认分页参数
	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.Limit > 100 {
		request.Limit = 100
	}

	// 搜索景点
	attractions, total, err := c.attractionService.SearchAttractions(ctx, &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"attractions": attractions,
		"meta": gin.H{
			"total": total,
			"limit": request.Limit,
			"page":  request.Page,
		},
	})
}

// GetPopularAttractions 获取热门景点
func (c *AttractionController) GetPopularAttractions(ctx *gin.Context) {
	// 获取查询参数
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 获取热门景点
	attractions, err := c.attractionService.GetPopularAttractions(ctx, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"attractions": attractions})
}

// GetAttractionsByCategory 按分类获取景点
func (c *AttractionController) GetAttractionsByCategory(ctx *gin.Context) {
	// 获取路径参数和查询参数
	category := ctx.Param("category")
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 获取分类景点
	attractions, err := c.attractionService.GetAttractionsByCategory(ctx, category, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"attractions": attractions})
}

// GetAttractionsByCountry 按国家获取景点
func (c *AttractionController) GetAttractionsByCountry(ctx *gin.Context) {
	// 获取路径参数和查询参数
	country := ctx.Param("country")
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 获取国家景点
	attractions, err := c.attractionService.GetAttractionsByCountry(ctx, country, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"attractions": attractions})
}

// RegisterRoutes 注册路由
func (c *AttractionController) RegisterRoutes(router *gin.RouterGroup) {
	attractions := router.Group("/attractions")
	{
		attractions.GET("", c.SearchAttractions)
		attractions.GET("/popular", c.GetPopularAttractions)
		attractions.GET("/category/:category", c.GetAttractionsByCategory)
		attractions.GET("/country/:country", c.GetAttractionsByCountry)
		attractions.GET("/:id", c.GetAttractionByID)
	}
}

// ListAttractions 获取景点列表
func ListAttractions(c *gin.Context) {
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

	// 构建搜索请求
	request := &models.SearchAttractionsRequest{
		Name:     c.Query("query"),
		Category: c.Query("category"),
		Country:  c.Query("country"),
		City:     c.Query("city"),
		Page:     page,
		Limit:    limit,
	}

	// 获取景点列表
	attractions, total, err := attractionService.SearchAttractions(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// 计算总页数
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   attractions,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": totalPages,
		},
	})
}

// GetAttraction 获取景点详情
func GetAttraction(c *gin.Context) {
	// 解析路径参数
	attractionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "无效的景点ID",
		})
		return
	}

	// 获取景点详情
	attraction, err := attractionService.GetAttractionByID(c, attractionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   attraction,
	})
}

// CreateAttraction 创建新景点
func CreateAttraction(c *gin.Context) {
	// 简单实现，不实际创建景点
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "error",
		"error":  "景点创建功能尚未实现",
	})
}

// UpdateAttraction 更新景点信息
func UpdateAttraction(c *gin.Context) {
	// 简单实现，不实际更新景点
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "error",
		"error":  "景点更新功能尚未实现",
	})
}

// DeleteAttraction 删除景点
func DeleteAttraction(c *gin.Context) {
	// 简单实现，不实际删除景点
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "error",
		"error":  "景点删除功能尚未实现",
	})
}
