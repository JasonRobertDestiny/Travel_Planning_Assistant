package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"traveler_agent/models"
)

// ItineraryRepository 行程仓库接口
type ItineraryRepository interface {
	Create(itinerary *models.Itinerary) (int64, error)
	GetByID(id int64) (*models.Itinerary, error)
	GetByUserID(userID int64, limit, offset int) ([]*models.Itinerary, error)
	Update(itinerary *models.Itinerary) error
	Delete(id int64) error
	AddDay(day *models.ItineraryDay) (int64, error)
	UpdateDay(day *models.ItineraryDay) error
	DeleteDay(dayID int64) error
	AddItem(item *models.ItineraryItem) (int64, error)
	UpdateItem(item *models.ItineraryItem) error
	DeleteItem(itemID int64) error
	GetPublicItineraries(limit, offset int) ([]*models.Itinerary, error)
}

// PostgresItineraryRepository PostgreSQL行程仓库实现
type PostgresItineraryRepository struct {
	db *sql.DB
}

// NewItineraryRepository 创建行程仓库
func NewItineraryRepository(db *sql.DB) ItineraryRepository {
	return &PostgresItineraryRepository{db: db}
}

// Create 创建行程
func (r *PostgresItineraryRepository) Create(itinerary *models.Itinerary) (int64, error) {
	// 当前时间作为创建时间和更新时间
	now := time.Now()
	itinerary.CreatedAt = now
	itinerary.UpdatedAt = now

	// 开始事务
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// 插入行程数据
	var id int64
	query := `
		INSERT INTO itineraries (user_id, title, description, destination, start_date, end_date, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	err = tx.QueryRow(
		query,
		itinerary.UserID,
		itinerary.Title,
		itinerary.Description,
		itinerary.Destination,
		itinerary.StartDate,
		itinerary.EndDate,
		itinerary.IsPublic,
		itinerary.CreatedAt,
		itinerary.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	itinerary.ID = id

	// 如果有天行程，则插入天行程
	if len(itinerary.Days) > 0 {
		for i, day := range itinerary.Days {
			day.ItineraryID = id
			day.DayNumber = i + 1 // 从1开始的天数

			dayID, err := r.addDayTx(tx, day)
			if err != nil {
				return 0, err
			}

			day.ID = dayID
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

// addDayTx 在事务中添加行程天数
func (r *PostgresItineraryRepository) addDayTx(tx *sql.Tx, day *models.ItineraryDay) (int64, error) {
	var id int64
	query := `
		INSERT INTO itinerary_days (itinerary_id, day_number, date, note)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := tx.QueryRow(
		query,
		day.ItineraryID,
		day.DayNumber,
		day.Date,
		day.Note,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	// 如果有行程项目，则插入行程项目
	if len(day.Items) > 0 {
		for i, item := range day.Items {
			item.ItineraryDayID = id
			item.Order = i + 1 // 从1开始的排序

			itemID, err := r.addItemTx(tx, item)
			if err != nil {
				return 0, err
			}

			item.ID = itemID
		}
	}

	return id, nil
}

// addItemTx 在事务中添加行程项目
func (r *PostgresItineraryRepository) addItemTx(tx *sql.Tx, item *models.ItineraryItem) (int64, error) {
	var id int64
	query := `
		INSERT INTO itinerary_items (itinerary_day_id, type, ref_id, title, description, start_time, end_time, duration, item_order, location, latitude, longitude, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`
	err := tx.QueryRow(
		query,
		item.ItineraryDayID,
		item.Type,
		item.RefID,
		item.Title,
		item.Description,
		item.StartTime,
		item.EndTime,
		item.Duration,
		item.Order,
		item.Location,
		item.Latitude,
		item.Longitude,
		item.Notes,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID 通过ID查询行程
func (r *PostgresItineraryRepository) GetByID(id int64) (*models.Itinerary, error) {
	// 查询行程信息
	query := `
		SELECT id, user_id, title, description, destination, start_date, end_date, is_public, created_at, updated_at
		FROM itineraries
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	itinerary, err := scanItinerary(row)
	if err != nil {
		return nil, err
	}

	// 查询行程天数
	daysQuery := `
		SELECT id, itinerary_id, day_number, date, note
		FROM itinerary_days
		WHERE itinerary_id = $1
		ORDER BY day_number
	`
	daysRows, err := r.db.Query(daysQuery, id)
	if err != nil {
		return nil, err
	}
	defer daysRows.Close()

	var days []*models.ItineraryDay
	for daysRows.Next() {
		var day models.ItineraryDay
		err := daysRows.Scan(
			&day.ID,
			&day.ItineraryID,
			&day.DayNumber,
			&day.Date,
			&day.Note,
		)
		if err != nil {
			return nil, err
		}

		// 查询行程项目
		itemsQuery := `
			SELECT id, itinerary_day_id, type, ref_id, title, description, start_time, end_time, duration, item_order, location, latitude, longitude, notes
			FROM itinerary_items
			WHERE itinerary_day_id = $1
			ORDER BY item_order
		`
		itemsRows, err := r.db.Query(itemsQuery, day.ID)
		if err != nil {
			return nil, err
		}
		defer itemsRows.Close()

		var items []*models.ItineraryItem
		for itemsRows.Next() {
			var item models.ItineraryItem
			var refID, startTime, endTime, latitude, longitude sql.NullInt64
			var description, location, notes sql.NullString

			err := itemsRows.Scan(
				&item.ID,
				&item.ItineraryDayID,
				&item.Type,
				&refID,
				&item.Title,
				&description,
				&startTime,
				&endTime,
				&item.Duration,
				&item.Order,
				&location,
				&latitude,
				&longitude,
				&notes,
			)
			if err != nil {
				return nil, err
			}

			if refID.Valid {
				item.RefID = refID.Int64
			}
			if description.Valid {
				item.Description = description.String
			}
			if startTime.Valid {
				item.StartTime = formatMinutesToTime(int(startTime.Int64))
			}
			if endTime.Valid {
				item.EndTime = formatMinutesToTime(int(endTime.Int64))
			}
			if location.Valid {
				item.Location = location.String
			}
			if latitude.Valid {
				item.Latitude = float64(latitude.Int64)
			}
			if longitude.Valid {
				item.Longitude = float64(longitude.Int64)
			}
			if notes.Valid {
				item.Notes = notes.String
			}

			items = append(items, &item)
		}

		day.Items = items
		days = append(days, &day)
	}

	itinerary.Days = days
	return itinerary, nil
}

// scanItinerary 扫描行程记录
func scanItinerary(row *sql.Row) (*models.Itinerary, error) {
	var itinerary models.Itinerary
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

// GetByUserID 获取用户的行程列表
func (r *PostgresItineraryRepository) GetByUserID(userID int64, limit, offset int) ([]*models.Itinerary, error) {
	query := `
		SELECT id, user_id, title, description, destination, start_date, end_date, is_public, created_at, updated_at
		FROM itineraries
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itineraries []*models.Itinerary
	for rows.Next() {
		var itinerary models.Itinerary
		var description sql.NullString

		err := rows.Scan(
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

		// 查询天数
		daysCountQuery := `
			SELECT COUNT(*) FROM itinerary_days WHERE itinerary_id = $1
		`
		var daysCount int
		err = r.db.QueryRow(daysCountQuery, itinerary.ID).Scan(&daysCount)
		if err != nil {
			return nil, err
		}

		// 创建空的天数数组以避免nil
		itinerary.Days = make([]*models.ItineraryDay, 0)

		itineraries = append(itineraries, &itinerary)
	}

	return itineraries, nil
}

// Update 更新行程
func (r *PostgresItineraryRepository) Update(itinerary *models.Itinerary) error {
	// 更新时间为当前时间
	itinerary.UpdatedAt = time.Now()

	query := `
		UPDATE itineraries
		SET title = $1, description = $2, destination = $3, start_date = $4, end_date = $5, is_public = $6, updated_at = $7
		WHERE id = $8
	`
	result, err := r.db.Exec(
		query,
		itinerary.Title,
		itinerary.Description,
		itinerary.Destination,
		itinerary.StartDate,
		itinerary.EndDate,
		itinerary.IsPublic,
		itinerary.UpdatedAt,
		itinerary.ID,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("行程不存在")
	}

	return nil
}

// Delete 删除行程
func (r *PostgresItineraryRepository) Delete(id int64) error {
	// 开始事务
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// 删除行程项目
	_, err = tx.Exec(`
		DELETE FROM itinerary_items 
		WHERE itinerary_day_id IN (SELECT id FROM itinerary_days WHERE itinerary_id = $1)
	`, id)
	if err != nil {
		return err
	}

	// 删除行程天数
	_, err = tx.Exec(`DELETE FROM itinerary_days WHERE itinerary_id = $1`, id)
	if err != nil {
		return err
	}

	// 删除行程
	result, err := tx.Exec(`DELETE FROM itineraries WHERE id = $1`, id)
	if err != nil {
		return err
	}

	// 检查是否删除成功
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("行程不存在")
	}

	// 提交事务
	return tx.Commit()
}

// AddDay 添加行程天数
func (r *PostgresItineraryRepository) AddDay(day *models.ItineraryDay) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	dayID, err := r.addDayTx(tx, day)
	if err != nil {
		return 0, err
	}

	return dayID, tx.Commit()
}

// UpdateDay 更新行程天数
func (r *PostgresItineraryRepository) UpdateDay(day *models.ItineraryDay) error {
	query := `
		UPDATE itinerary_days
		SET day_number = $1, date = $2, note = $3
		WHERE id = $4
	`
	result, err := r.db.Exec(
		query,
		day.DayNumber,
		day.Date,
		day.Note,
		day.ID,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("行程天数不存在")
	}

	return nil
}

// DeleteDay 删除行程天数
func (r *PostgresItineraryRepository) DeleteDay(dayID int64) error {
	// 开始事务
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// 删除行程项目
	_, err = tx.Exec(`DELETE FROM itinerary_items WHERE itinerary_day_id = $1`, dayID)
	if err != nil {
		return err
	}

	// 删除行程天数
	result, err := tx.Exec(`DELETE FROM itinerary_days WHERE id = $1`, dayID)
	if err != nil {
		return err
	}

	// 检查是否删除成功
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("行程天数不存在")
	}

	// 提交事务
	return tx.Commit()
}

// AddItem 添加行程项目
func (r *PostgresItineraryRepository) AddItem(item *models.ItineraryItem) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	itemID, err := r.addItemTx(tx, item)
	if err != nil {
		return 0, err
	}

	return itemID, tx.Commit()
}

// UpdateItem 更新行程项目
func (r *PostgresItineraryRepository) UpdateItem(item *models.ItineraryItem) error {
	query := `
		UPDATE itinerary_items
		SET type = $1, ref_id = $2, title = $3, description = $4, start_time = $5, end_time = $6, 
			duration = $7, item_order = $8, location = $9, latitude = $10, longitude = $11, notes = $12
		WHERE id = $13
	`
	result, err := r.db.Exec(
		query,
		item.Type,
		item.RefID,
		item.Title,
		item.Description,
		item.StartTime,
		item.EndTime,
		item.Duration,
		item.Order,
		item.Location,
		item.Latitude,
		item.Longitude,
		item.Notes,
		item.ID,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("行程项目不存在")
	}

	return nil
}

// DeleteItem 删除行程项目
func (r *PostgresItineraryRepository) DeleteItem(itemID int64) error {
	result, err := r.db.Exec(`DELETE FROM itinerary_items WHERE id = $1`, itemID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("行程项目不存在")
	}

	return nil
}

// GetPublicItineraries 获取公开的行程
func (r *PostgresItineraryRepository) GetPublicItineraries(limit, offset int) ([]*models.Itinerary, error) {
	query := `
		SELECT id, user_id, title, description, destination, start_date, end_date, is_public, created_at, updated_at
		FROM itineraries
		WHERE is_public = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itineraries []*models.Itinerary
	for rows.Next() {
		var itinerary models.Itinerary
		var description sql.NullString

		err := rows.Scan(
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

		// 查询天数
		daysCountQuery := `
			SELECT COUNT(*) FROM itinerary_days WHERE itinerary_id = $1
		`
		var daysCount int
		err = r.db.QueryRow(daysCountQuery, itinerary.ID).Scan(&daysCount)
		if err != nil {
			return nil, err
		}

		// 创建空的天数数组以避免nil
		itinerary.Days = make([]*models.ItineraryDay, 0)

		itineraries = append(itineraries, &itinerary)
	}

	return itineraries, nil
}

// formatMinutesToTime 将分钟数转换为时间字符串 (如: 9:30)
func formatMinutesToTime(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60
	return fmt.Sprintf("%d:%02d", hours, mins)
}
