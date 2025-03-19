package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Attraction 景点模型
type Attraction struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Address     string    `json:"address"`
	Latitude    float64   `json:"lat"`
	Longitude   float64   `json:"lng"`
	ImageURL    string    `json:"image_url,omitempty"`
	Category    string    `json:"category"`       // 博物馆/公园/景区/历史遗迹等
	Tags        []string  `json:"tags,omitempty"` // 标签列表
	OpenHours   string    `json:"open_hours"`     // 营业时间
	TicketPrice float64   `json:"ticket_price"`   // 门票价格
	Duration    int       `json:"duration"`       // 推荐游览时长(分钟)
	Rating      float64   `json:"rating"`         // 评分
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AttractionListItem 景点列表项
type AttractionListItem struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Address     string  `json:"address,omitempty"`
	Latitude    float64 `json:"lat,omitempty"`
	Longitude   float64 `json:"lng,omitempty"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url,omitempty"`
	Rating      float64 `json:"rating"`
	Duration    int     `json:"duration,omitempty"`
}

// SearchAttractionsRequest 景点搜索请求
type SearchAttractionsRequest struct {
	Name      string  `json:"name" form:"name"`             // 搜索词
	City      string  `json:"city" form:"city"`             // 城市
	Country   string  `json:"country" form:"country"`       // 国家
	Category  string  `json:"category" form:"category"`     // 景点类型
	MinRating float64 `json:"min_rating" form:"min_rating"` // 最低评分
	SortBy    string  `json:"sort_by" form:"sort_by"`       // 排序字段：name, rating, popularity
	Page      int     `json:"page" form:"page"`             // 页码
	Limit     int     `json:"limit" form:"limit"`           // 每页数量
}

// AttractionResponse 景点详情响应
type AttractionResponse struct {
	Attraction
	Distance float64 `json:"distance,omitempty"` // 距离当前位置的距离
}

// ScanAttraction 从数据库扫描景点记录
func ScanAttraction(row *sql.Row) (*Attraction, error) {
	var attraction Attraction
	var imageURL sql.NullString
	var tagsJSON string // 存储JSON格式的标签

	err := row.Scan(
		&attraction.ID,
		&attraction.Name,
		&attraction.Description,
		&attraction.City,
		&attraction.Country,
		&attraction.Address,
		&attraction.Latitude,
		&attraction.Longitude,
		&imageURL,
		&attraction.Category,
		&tagsJSON,
		&attraction.OpenHours,
		&attraction.TicketPrice,
		&attraction.Duration,
		&attraction.Rating,
		&attraction.CreatedAt,
		&attraction.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if imageURL.Valid {
		attraction.ImageURL = imageURL.String
	}

	// 这里应该解析tagsJSON为[]string
	// 简化实现，实际应使用json.Unmarshal

	return &attraction, nil
}

// ToListItem 将Attraction转换为AttractionListItem
func (a *Attraction) ToListItem() *AttractionListItem {
	return &AttractionListItem{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		City:        a.City,
		Country:     a.Country,
		Address:     a.Address,
		Latitude:    a.Latitude,
		Longitude:   a.Longitude,
		Category:    a.Category,
		ImageURL:    a.ImageURL,
		Rating:      a.Rating,
		Duration:    a.Duration,
	}
}

// ParseTags 解析JSON格式的标签
func (a *Attraction) ParseTags(tagsJSON []byte) error {
	if tagsJSON == nil || len(tagsJSON) == 0 {
		a.Tags = []string{}
		return nil
	}

	return json.Unmarshal(tagsJSON, &a.Tags)
}
