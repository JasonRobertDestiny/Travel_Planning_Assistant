package models

import (
	"database/sql"
	"time"
)

// Itinerary 行程模型
type Itinerary struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Destination string          `json:"destination"` // 目的地/城市
	StartDate   time.Time       `json:"start_date"`
	EndDate     time.Time       `json:"end_date"`
	IsPublic    bool            `json:"is_public"` // 是否公开分享
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Days        []*ItineraryDay `json:"days,omitempty"`
}

// ItineraryDay 行程天数
type ItineraryDay struct {
	ID          int64            `json:"id"`
	ItineraryID int64            `json:"itinerary_id"`
	DayNumber   int              `json:"day_number"` // 第几天
	Date        time.Time        `json:"date"`
	Note        string           `json:"note,omitempty"`
	Items       []*ItineraryItem `json:"items,omitempty"`
}

// ItineraryItem 行程项目（每天的活动，如游览景点、用餐、交通等）
type ItineraryItem struct {
	ID             int64   `json:"id"`
	ItineraryDayID int64   `json:"itinerary_day_id"`
	Type           string  `json:"type"`   // 类型：attraction, hotel, transport, meal, etc.
	RefID          int64   `json:"ref_id"` // 引用ID，如景点ID
	Title          string  `json:"title"`  // 标题
	Description    string  `json:"description,omitempty"`
	StartTime      string  `json:"start_time"` // 开始时间（如 09:00）
	EndTime        string  `json:"end_time"`   // 结束时间（如 11:00）
	Duration       int     `json:"duration"`   // 持续时间（分钟）
	Order          int     `json:"order"`      // 排序
	Location       string  `json:"location,omitempty"`
	Latitude       float64 `json:"lat,omitempty"`
	Longitude      float64 `json:"lng,omitempty"`
	Notes          string  `json:"notes,omitempty"`
}

// CreateItineraryRequest 创建行程请求
type CreateItineraryRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Destination string    `json:"destination" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	IsPublic    bool      `json:"is_public"`
}

// ItineraryResponse 行程响应
type ItineraryResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Destination string    `json:"destination"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	DaysCount   int       `json:"days_count"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToItineraryResponse 将Itinerary转换为ItineraryResponse
func (i *Itinerary) ToItineraryResponse() ItineraryResponse {
	return ItineraryResponse{
		ID:          i.ID,
		Title:       i.Title,
		Description: i.Description,
		Destination: i.Destination,
		StartDate:   i.StartDate,
		EndDate:     i.EndDate,
		DaysCount:   int(i.EndDate.Sub(i.StartDate).Hours()/24) + 1,
		IsPublic:    i.IsPublic,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

// ScanItinerary 从数据库扫描行程记录
func ScanItinerary(row *sql.Row) (*Itinerary, error) {
	var itinerary Itinerary
	var description sql.NullString

	err := row.Scan(
		&itinerary.ID,
		&itinerary.UserID,
		&itinerary.Title,
		&description,
		&itinerary.Destination,
		&itinerary.StartDate,
		&itinerary.EndDate,
		&itinerary.IsPublic,
		&itinerary.CreatedAt,
		&itinerary.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if description.Valid {
		itinerary.Description = description.String
	}

	return &itinerary, nil
}

// AddDay 向行程添加一天
func (i *Itinerary) AddDay(day *ItineraryDay) {
	if i.Days == nil {
		i.Days = make([]*ItineraryDay, 0)
	}
	i.Days = append(i.Days, day)
}

// AddItemToDay 向行程天添加项目
func (d *ItineraryDay) AddItem(item *ItineraryItem) {
	if d.Items == nil {
		d.Items = make([]*ItineraryItem, 0)
	}
	d.Items = append(d.Items, item)
}
