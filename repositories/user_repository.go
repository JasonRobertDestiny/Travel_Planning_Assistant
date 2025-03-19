package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"traveler_agent/models"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (int64, error)
	Update(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error
	Delete(ctx context.Context, id int64) error
	GetUserPreferences(ctx context.Context, userID int64) (*models.UserPreference, error)
	SaveUserPreferences(ctx context.Context, pref *models.UserPreference) error
	GetTravelPreferences(ctx context.Context, userID int64) (*models.UserPreference, error)
	SaveTravelPreferences(ctx context.Context, pref *models.UserPreference) error
}

// UserRepositoryImpl 用户仓库实现
type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository 创建一个新的用户仓库实例
func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// GetByID 根据ID获取用户
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, 
		       phone, avatar, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	var isActive sql.NullBool
	var phone, avatar sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.FirstName, &user.LastName, &phone, &avatar,
		&isActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 用户不存在
		}
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

	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, 
		       phone, avatar, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	var isActive sql.NullBool
	var phone, avatar sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.FirstName, &user.LastName, &phone, &avatar,
		&isActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 用户不存在
		}
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

	return user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, 
		       phone, avatar, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &models.User{}
	var isActive sql.NullBool
	var phone, avatar sql.NullString

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.FirstName, &user.LastName, &phone, &avatar,
		&isActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 用户不存在
		}
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

	return user, nil
}

// Create 创建新用户
func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.User) (int64, error) {
	query := `
		INSERT INTO users (
			username, email, password, first_name, last_name,
			phone, avatar, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	if user.IsActive == nil {
		isActive := true
		user.IsActive = &isActive
	}

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.Password,
		user.FirstName, user.LastName, user.Phone,
		user.Avatar, user.IsActive, user.CreatedAt, user.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update 更新用户信息
func (r *UserRepositoryImpl) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, first_name = $3, last_name = $4,
			phone = $5, avatar = $6, is_active = $7, updated_at = $8
		WHERE id = $9
	`

	user.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		user.Username, user.Email, user.FirstName, user.LastName,
		user.Phone, user.Avatar, user.IsActive, user.UpdatedAt, user.ID,
	)

	return err
}

// UpdatePassword 更新用户密码
func (r *UserRepositoryImpl) UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, hashedPassword, time.Now(), userID)
	return err
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetUserPreferences 获取用户偏好设置
func (r *UserRepositoryImpl) GetUserPreferences(ctx context.Context, userID int64) (*models.UserPreference, error) {
	query := `
		SELECT id, user_id, language, currency, 
               notification_enabled, theme, created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`

	pref := &models.UserPreference{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&pref.ID, &pref.UserID, &pref.Language, &pref.Currency,
		&pref.NotificationEnabled, &pref.Theme, &pref.CreatedAt, &pref.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 用户偏好不存在
		}
		return nil, err
	}

	return pref, nil
}

// SaveUserPreferences 保存用户偏好设置
func (r *UserRepositoryImpl) SaveUserPreferences(ctx context.Context, pref *models.UserPreference) error {
	// 检查是否已存在
	existing, err := r.GetUserPreferences(ctx, pref.UserID)
	if err != nil {
		return err
	}

	now := time.Now()
	pref.UpdatedAt = now

	if existing == nil {
		// 创建新记录
		query := `
			INSERT INTO user_preferences (
				user_id, language, currency, notification_enabled, 
				theme, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`

		pref.CreatedAt = now
		err := r.db.QueryRowContext(ctx, query,
			pref.UserID, pref.Language, pref.Currency,
			pref.NotificationEnabled, pref.Theme,
			pref.CreatedAt, pref.UpdatedAt,
		).Scan(&pref.ID)

		return err
	} else {
		// 更新现有记录
		query := `
			UPDATE user_preferences
			SET language = $1, currency = $2, notification_enabled = $3,
				theme = $4, updated_at = $5
			WHERE user_id = $6
		`

		_, err := r.db.ExecContext(ctx, query,
			pref.Language, pref.Currency, pref.NotificationEnabled,
			pref.Theme, pref.UpdatedAt, pref.UserID,
		)

		return err
	}
}

// GetTravelPreferences 获取用户旅行偏好设置
func (r *UserRepositoryImpl) GetTravelPreferences(ctx context.Context, userID int64) (*models.UserPreference, error) {
	query := `
		SELECT id, user_id, travel_style, budget_level, 
               transport_prefer, preferred_tags, excluded_tags
		FROM user_travel_preferences
		WHERE user_id = $1
	`

	var pref models.UserPreference
	var preferredTags, excludedTags []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&pref.ID, &pref.UserID, &pref.TravelStyle, &pref.BudgetLevel,
		&pref.TransportPrefer, &preferredTags, &excludedTags,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 用户旅行偏好不存在
		}
		return nil, err
	}

	// 解析JSON数组
	if len(preferredTags) > 0 {
		if err := json.Unmarshal(preferredTags, &pref.PreferredTags); err != nil {
			return nil, err
		}
	}

	if len(excludedTags) > 0 {
		if err := json.Unmarshal(excludedTags, &pref.ExcludedTags); err != nil {
			return nil, err
		}
	}

	return &pref, nil
}

// SaveTravelPreferences 保存用户旅行偏好设置
func (r *UserRepositoryImpl) SaveTravelPreferences(ctx context.Context, pref *models.UserPreference) error {
	// 检查是否已存在
	existing, err := r.GetTravelPreferences(ctx, pref.UserID)
	if err != nil {
		return err
	}

	// 将字符串数组转为JSON
	preferredTagsJSON, err := json.Marshal(pref.PreferredTags)
	if err != nil {
		return err
	}

	excludedTagsJSON, err := json.Marshal(pref.ExcludedTags)
	if err != nil {
		return err
	}

	now := time.Now()
	pref.UpdatedAt = now

	if existing == nil {
		// 创建新记录
		query := `
			INSERT INTO user_travel_preferences (
				user_id, travel_style, budget_level, transport_prefer,
				preferred_tags, excluded_tags, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`

		pref.CreatedAt = now
		err := r.db.QueryRowContext(ctx, query,
			pref.UserID, pref.TravelStyle, pref.BudgetLevel, pref.TransportPrefer,
			preferredTagsJSON, excludedTagsJSON, pref.CreatedAt, pref.UpdatedAt,
		).Scan(&pref.ID)

		return err
	} else {
		// 更新现有记录
		query := `
			UPDATE user_travel_preferences
			SET travel_style = $1, budget_level = $2, transport_prefer = $3,
				preferred_tags = $4, excluded_tags = $5, updated_at = $6
			WHERE user_id = $7
		`

		_, err := r.db.ExecContext(ctx, query,
			pref.TravelStyle, pref.BudgetLevel, pref.TransportPrefer,
			preferredTagsJSON, excludedTagsJSON, pref.UpdatedAt, pref.UserID,
		)

		return err
	}
}
