package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"traveler_agent/models"
)

// AttractionRepository 景点仓库接口
type AttractionRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Attraction, error)
	Search(ctx context.Context, params *models.SearchAttractionsRequest) ([]*models.AttractionListItem, int, error)
	GetPopular(ctx context.Context, limit int) ([]*models.AttractionListItem, error)
	GetByCategory(ctx context.Context, category string, limit int) ([]*models.AttractionListItem, error)
	GetByCountry(ctx context.Context, country string, limit int) ([]*models.AttractionListItem, error)
	GetByCityOrCountry(ctx context.Context, destination string, limit int) ([]*models.AttractionListItem, error)
}

// PostgresAttractionRepository PostgreSQL景点仓库实现
type PostgresAttractionRepository struct {
	db *sql.DB
}

// NewAttractionRepository 创建景点仓库
func NewAttractionRepository(db *sql.DB) AttractionRepository {
	return &PostgresAttractionRepository{db: db}
}

// GetByID 通过ID查询景点
func (r *PostgresAttractionRepository) GetByID(ctx context.Context, id int64) (*models.Attraction, error) {
	var attraction models.Attraction
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, city, country, address, latitude, longitude, 
		image_url, category, tags, open_hours, ticket_price, duration, rating, 
		created_at, updated_at
		FROM attractions WHERE id = $1
	`, id).Scan(
		&attraction.ID, &attraction.Name, &attraction.Description, &attraction.City,
		&attraction.Country, &attraction.Address, &attraction.Latitude, &attraction.Longitude,
		&attraction.ImageURL, &attraction.Category, &attraction.Tags, &attraction.OpenHours,
		&attraction.TicketPrice, &attraction.Duration, &attraction.Rating,
		&attraction.CreatedAt, &attraction.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 返回nil表示未找到
		}
		return nil, err
	}
	return &attraction, nil
}

// buildSearchQuery 构建搜索查询
func (r *PostgresAttractionRepository) buildSearchQuery(params *models.SearchAttractionsRequest) (string, []interface{}, error) {
	// 构建查询条件
	conditions := []string{}
	args := []interface{}{}
	paramIndex := 1

	// 添加搜索条件
	if params.Name != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", paramIndex, paramIndex))
		args = append(args, "%"+params.Name+"%")
		paramIndex++
	}

	if params.City != "" {
		conditions = append(conditions, fmt.Sprintf("city ILIKE $%d", paramIndex))
		args = append(args, "%"+params.City+"%")
		paramIndex++
	}

	if params.Country != "" {
		conditions = append(conditions, fmt.Sprintf("country ILIKE $%d", paramIndex))
		args = append(args, "%"+params.Country+"%")
		paramIndex++
	}

	if params.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", paramIndex))
		args = append(args, params.Category)
		paramIndex++
	}

	if params.MinRating > 0 {
		conditions = append(conditions, fmt.Sprintf("rating >= $%d", paramIndex))
		args = append(args, params.MinRating)
		paramIndex++
	}

	// 基础查询
	query := `
		SELECT id, name, description, city, country, address, image_url, category, rating, latitude, longitude, duration
		FROM attractions
	`

	// 添加条件到查询
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// 添加排序
	orderBy := "rating DESC"
	if params.SortBy != "" {
		switch params.SortBy {
		case "name":
			orderBy = "name ASC"
		case "rating":
			orderBy = "rating DESC"
		case "popularity":
			orderBy = "popularity DESC"
		}
	}
	query += " ORDER BY " + orderBy

	// 添加分页
	limit := 10
	if params.Limit > 0 && params.Limit <= 50 {
		limit = params.Limit
	}

	offset := 0
	if params.Page > 1 {
		offset = (params.Page - 1) * limit
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	args = append(args, limit, offset)

	return query, args, nil
}

// Search 搜索景点
func (r *PostgresAttractionRepository) Search(ctx context.Context, params *models.SearchAttractionsRequest) ([]*models.AttractionListItem, int, error) {
	if params == nil {
		params = &models.SearchAttractionsRequest{
			Page:  1,
			Limit: 10,
		}
	}

	// 计算总记录数
	countQuery := "SELECT COUNT(*) FROM attractions"
	countConditions := []string{}
	countArgs := []interface{}{}
	paramIndex := 1

	if params.Name != "" {
		countConditions = append(countConditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", paramIndex, paramIndex))
		countArgs = append(countArgs, "%"+params.Name+"%")
		paramIndex++
	}

	if params.City != "" {
		countConditions = append(countConditions, fmt.Sprintf("city ILIKE $%d", paramIndex))
		countArgs = append(countArgs, "%"+params.City+"%")
		paramIndex++
	}

	if params.Country != "" {
		countConditions = append(countConditions, fmt.Sprintf("country ILIKE $%d", paramIndex))
		countArgs = append(countArgs, "%"+params.Country+"%")
		paramIndex++
	}

	if params.Category != "" {
		countConditions = append(countConditions, fmt.Sprintf("category = $%d", paramIndex))
		countArgs = append(countArgs, params.Category)
		paramIndex++
	}

	if params.MinRating > 0 {
		countConditions = append(countConditions, fmt.Sprintf("rating >= $%d", paramIndex))
		countArgs = append(countArgs, params.MinRating)
		paramIndex++
	}

	if len(countConditions) > 0 {
		countQuery += " WHERE " + strings.Join(countConditions, " AND ")
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 构建查询
	query, args, err := r.buildSearchQuery(params)
	if err != nil {
		return nil, 0, err
	}

	// 执行查询
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 处理结果
	var results []*models.AttractionListItem
	for rows.Next() {
		var attraction models.AttractionListItem

		err := rows.Scan(
			&attraction.ID,
			&attraction.Name,
			&attraction.Description,
			&attraction.City,
			&attraction.Country,
			&attraction.Address,
			&attraction.ImageURL,
			&attraction.Category,
			&attraction.Rating,
			&attraction.Latitude,
			&attraction.Longitude,
			&attraction.Duration,
		)
		if err != nil {
			return nil, 0, err
		}

		results = append(results, &attraction)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetPopular 获取热门景点
func (r *PostgresAttractionRepository) GetPopular(ctx context.Context, limit int) ([]*models.AttractionListItem, error) {
	if limit <= 0 {
		limit = 10 // 默认限制
	}

	query := `
		SELECT id, name, description, city, country, address, image_url, category, rating, latitude, longitude, duration
		FROM attractions
		ORDER BY rating DESC, popularity DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.AttractionListItem
	for rows.Next() {
		var attraction models.AttractionListItem

		err := rows.Scan(
			&attraction.ID,
			&attraction.Name,
			&attraction.Description,
			&attraction.City,
			&attraction.Country,
			&attraction.Address,
			&attraction.ImageURL,
			&attraction.Category,
			&attraction.Rating,
			&attraction.Latitude,
			&attraction.Longitude,
			&attraction.Duration,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &attraction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetByCategory 按分类获取景点
func (r *PostgresAttractionRepository) GetByCategory(ctx context.Context, category string, limit int) ([]*models.AttractionListItem, error) {
	if limit <= 0 {
		limit = 10 // 默认限制
	}

	query := `
		SELECT id, name, description, city, country, address, image_url, category, rating, latitude, longitude, duration
		FROM attractions
		WHERE category = $1
		ORDER BY rating DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, category, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.AttractionListItem
	for rows.Next() {
		var attraction models.AttractionListItem

		err := rows.Scan(
			&attraction.ID,
			&attraction.Name,
			&attraction.Description,
			&attraction.City,
			&attraction.Country,
			&attraction.Address,
			&attraction.ImageURL,
			&attraction.Category,
			&attraction.Rating,
			&attraction.Latitude,
			&attraction.Longitude,
			&attraction.Duration,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &attraction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetByCountry 按国家获取景点
func (r *PostgresAttractionRepository) GetByCountry(ctx context.Context, country string, limit int) ([]*models.AttractionListItem, error) {
	if limit <= 0 {
		limit = 10 // 默认限制
	}

	query := `
		SELECT id, name, description, city, country, address, image_url, category, rating, latitude, longitude, duration
		FROM attractions
		WHERE country ILIKE $1
		ORDER BY rating DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, "%"+country+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.AttractionListItem
	for rows.Next() {
		var attraction models.AttractionListItem

		err := rows.Scan(
			&attraction.ID,
			&attraction.Name,
			&attraction.Description,
			&attraction.City,
			&attraction.Country,
			&attraction.Address,
			&attraction.ImageURL,
			&attraction.Category,
			&attraction.Rating,
			&attraction.Latitude,
			&attraction.Longitude,
			&attraction.Duration,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &attraction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetByCityOrCountry 按城市或国家获取景点
func (r *PostgresAttractionRepository) GetByCityOrCountry(ctx context.Context, destination string, limit int) ([]*models.AttractionListItem, error) {
	if limit <= 0 {
		limit = 10 // 默认限制
	}

	query := `
		SELECT id, name, description, city, country, address, image_url, category, rating, latitude, longitude, duration
		FROM attractions
		WHERE city ILIKE $1 OR country ILIKE $1
		ORDER BY rating DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, "%"+destination+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.AttractionListItem
	for rows.Next() {
		var attraction models.AttractionListItem

		err := rows.Scan(
			&attraction.ID,
			&attraction.Name,
			&attraction.Description,
			&attraction.City,
			&attraction.Country,
			&attraction.Address,
			&attraction.ImageURL,
			&attraction.Category,
			&attraction.Rating,
			&attraction.Latitude,
			&attraction.Longitude,
			&attraction.Duration,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &attraction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
