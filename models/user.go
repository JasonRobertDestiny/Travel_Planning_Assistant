package models

import (
	"database/sql"
	"time"
)

// User 用户模型
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // 不在JSON响应中返回密码
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	IsActive  *bool     `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserPreference 用户偏好模型
type UserPreference struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"user_id"`
	Language            string    `json:"language,omitempty"`
	Currency            string    `json:"currency,omitempty"`
	NotificationEnabled bool      `json:"notification_enabled"`
	Theme               string    `json:"theme,omitempty"`
	TravelStyle         string    `json:"travel_style,omitempty"` // 旅行风格: 文化/自然/美食/购物等
	BudgetLevel         string    `json:"budget_level,omitempty"` // 预算级别: 经济/舒适/豪华
	PreferredTags       []string  `json:"preferred_tags,omitempty"`
	ExcludedTags        []string  `json:"excluded_tags,omitempty"`
	TransportPrefer     string    `json:"transport_prefer,omitempty"` // 交通偏好: 公共交通/步行/自驾等
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"omitempty"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应 (过滤敏感信息)
type UserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar,omitempty"`
}

// ToUserResponse 将User转换为UserResponse
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Avatar:    u.Avatar,
	}
}

// ScanUser 从数据库扫描用户记录
func ScanUser(row *sql.Row) (*User, error) {
	var user User
	var phone, avatar sql.NullString
	var isActive sql.NullBool

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&phone,
		&avatar,
		&isActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if phone.Valid {
		user.Phone = phone.String
	}

	if avatar.Valid {
		user.Avatar = avatar.String
	}

	if isActive.Valid {
		user.IsActive = &isActive.Bool
	} else {
		active := true // 默认为活跃状态
		user.IsActive = &active
	}

	return &user, nil
}
