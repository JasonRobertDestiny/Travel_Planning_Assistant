package services

import (
	"errors"
	"time"

	"traveler_agent/models"
	"traveler_agent/repositories"
)

// ItineraryService 行程服务接口
type ItineraryService interface {
	CreateItinerary(userID int64, request *models.CreateItineraryRequest) (*models.ItineraryResponse, error)
	GetItineraryByID(id int64) (*models.Itinerary, error)
	GetUserItineraries(userID int64, page, pageSize int) ([]*models.ItineraryResponse, error)
	UpdateItinerary(id int64, userID int64, request *models.CreateItineraryRequest) (*models.ItineraryResponse, error)
	DeleteItinerary(id int64, userID int64) error
	AddDay(itineraryID int64, userID int64, date time.Time, note string) (*models.ItineraryDay, error)
	UpdateDay(dayID int64, userID int64, date time.Time, note string) (*models.ItineraryDay, error)
	DeleteDay(dayID int64, userID int64) error
	AddItem(dayID int64, userID int64, item *models.ItineraryItem) (*models.ItineraryItem, error)
	UpdateItem(itemID int64, userID int64, item *models.ItineraryItem) (*models.ItineraryItem, error)
	DeleteItem(itemID int64, userID int64) error
	GetPublicItineraries(page, pageSize int) ([]*models.ItineraryResponse, error)
}

// ItineraryServiceImpl 行程服务实现
type ItineraryServiceImpl struct {
	itineraryRepo repositories.ItineraryRepository
}

// NewItineraryService 创建行程服务
func NewItineraryService(itineraryRepo repositories.ItineraryRepository) ItineraryService {
	return &ItineraryServiceImpl{
		itineraryRepo: itineraryRepo,
	}
}

// CreateItinerary 创建行程
func (s *ItineraryServiceImpl) CreateItinerary(userID int64, request *models.CreateItineraryRequest) (*models.ItineraryResponse, error) {
	// 验证请求
	if request.Title == "" {
		return nil, errors.New("标题不能为空")
	}
	if request.Destination == "" {
		return nil, errors.New("目的地不能为空")
	}
	if request.StartDate.IsZero() || request.EndDate.IsZero() {
		return nil, errors.New("开始和结束日期不能为空")
	}
	if request.EndDate.Before(request.StartDate) {
		return nil, errors.New("结束日期不能早于开始日期")
	}

	// 创建行程对象
	itinerary := &models.Itinerary{
		UserID:      userID,
		Title:       request.Title,
		Description: request.Description,
		Destination: request.Destination,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		IsPublic:    request.IsPublic,
	}

	// 保存行程
	id, err := s.itineraryRepo.Create(itinerary)
	if err != nil {
		return nil, errors.New("创建行程失败: " + err.Error())
	}

	// 获取刚创建的行程
	itinerary, err = s.itineraryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("获取行程失败: " + err.Error())
	}

	// 返回行程响应
	response := itinerary.ToItineraryResponse()
	return &response, nil
}

// GetItineraryByID 获取行程详情
func (s *ItineraryServiceImpl) GetItineraryByID(id int64) (*models.Itinerary, error) {
	itinerary, err := s.itineraryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("行程不存在")
	}
	return itinerary, nil
}

// GetUserItineraries 获取用户的行程列表
func (s *ItineraryServiceImpl) GetUserItineraries(userID int64, page, pageSize int) ([]*models.ItineraryResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	itineraries, err := s.itineraryRepo.GetByUserID(userID, pageSize, offset)
	if err != nil {
		return nil, errors.New("获取行程列表失败: " + err.Error())
	}

	var responses []*models.ItineraryResponse
	for _, itinerary := range itineraries {
		response := itinerary.ToItineraryResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

// UpdateItinerary 更新行程
func (s *ItineraryServiceImpl) UpdateItinerary(id int64, userID int64, request *models.CreateItineraryRequest) (*models.ItineraryResponse, error) {
	// 获取行程
	itinerary, err := s.itineraryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("行程不存在")
	}

	// 验证所有权
	if itinerary.UserID != userID {
		return nil, errors.New("无权修改此行程")
	}

	// 验证请求
	if request.Title == "" {
		return nil, errors.New("标题不能为空")
	}
	if request.Destination == "" {
		return nil, errors.New("目的地不能为空")
	}
	if request.StartDate.IsZero() || request.EndDate.IsZero() {
		return nil, errors.New("开始和结束日期不能为空")
	}
	if request.EndDate.Before(request.StartDate) {
		return nil, errors.New("结束日期不能早于开始日期")
	}

	// 更新行程
	itinerary.Title = request.Title
	itinerary.Description = request.Description
	itinerary.Destination = request.Destination
	itinerary.StartDate = request.StartDate
	itinerary.EndDate = request.EndDate
	itinerary.IsPublic = request.IsPublic

	// 保存更新
	err = s.itineraryRepo.Update(itinerary)
	if err != nil {
		return nil, errors.New("更新行程失败: " + err.Error())
	}

	// 获取更新后的行程
	itinerary, err = s.itineraryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("获取行程失败: " + err.Error())
	}

	// 返回行程响应
	response := itinerary.ToItineraryResponse()
	return &response, nil
}

// DeleteItinerary 删除行程
func (s *ItineraryServiceImpl) DeleteItinerary(id int64, userID int64) error {
	// 获取行程
	itinerary, err := s.itineraryRepo.GetByID(id)
	if err != nil {
		return errors.New("行程不存在")
	}

	// 验证所有权
	if itinerary.UserID != userID {
		return errors.New("无权删除此行程")
	}

	// 删除行程
	err = s.itineraryRepo.Delete(id)
	if err != nil {
		return errors.New("删除行程失败: " + err.Error())
	}

	return nil
}

// AddDay 添加行程天数
func (s *ItineraryServiceImpl) AddDay(itineraryID int64, userID int64, date time.Time, note string) (*models.ItineraryDay, error) {
	// 获取行程
	itinerary, err := s.itineraryRepo.GetByID(itineraryID)
	if err != nil {
		return nil, errors.New("行程不存在")
	}

	// 验证所有权
	if itinerary.UserID != userID {
		return nil, errors.New("无权修改此行程")
	}

	// 验证日期
	if date.IsZero() {
		return nil, errors.New("日期不能为空")
	}

	// 获取当前天数
	dayNumber := len(itinerary.Days) + 1

	// 创建天数对象
	day := &models.ItineraryDay{
		ItineraryID: itineraryID,
		DayNumber:   dayNumber,
		Date:        date,
		Note:        note,
		Items:       make([]*models.ItineraryItem, 0),
	}

	// 保存天数
	dayID, err := s.itineraryRepo.AddDay(day)
	if err != nil {
		return nil, errors.New("添加天数失败: " + err.Error())
	}

	day.ID = dayID
	return day, nil
}

// UpdateDay 更新行程天数
func (s *ItineraryServiceImpl) UpdateDay(dayID int64, userID int64, date time.Time, note string) (*models.ItineraryDay, error) {
	// 获取行程
	itinerary, err := s.itineraryRepo.GetByID(0) // 这里需要根据dayID反查行程ID
	if err != nil {
		return nil, errors.New("无法获取行程信息")
	}

	// 验证所有权
	if itinerary.UserID != userID {
		return nil, errors.New("无权修改此行程")
	}

	// 在这个简化的版本中，我们假设能够找到对应的天数
	var day *models.ItineraryDay
	for _, d := range itinerary.Days {
		if d.ID == dayID {
			day = d
			break
		}
	}

	if day == nil {
		return nil, errors.New("天数不存在")
	}

	// 验证日期
	if date.IsZero() {
		return nil, errors.New("日期不能为空")
	}

	// 更新天数
	day.Date = date
	day.Note = note

	// 保存更新
	err = s.itineraryRepo.UpdateDay(day)
	if err != nil {
		return nil, errors.New("更新天数失败: " + err.Error())
	}

	return day, nil
}

// DeleteDay 删除行程天数
func (s *ItineraryServiceImpl) DeleteDay(dayID int64, userID int64) error {
	// 获取行程
	itinerary, err := s.itineraryRepo.GetByID(0) // 这里需要根据dayID反查行程ID
	if err != nil {
		return errors.New("无法获取行程信息")
	}

	// 验证所有权
	if itinerary.UserID != userID {
		return errors.New("无权修改此行程")
	}

	// 删除天数
	err = s.itineraryRepo.DeleteDay(dayID)
	if err != nil {
		return errors.New("删除天数失败: " + err.Error())
	}

	return nil
}

// AddItem 添加行程项目
func (s *ItineraryServiceImpl) AddItem(dayID int64, userID int64, item *models.ItineraryItem) (*models.ItineraryItem, error) {
	// 获取行程
	itinerary, err := s.itineraryRepo.GetByID(0) // 这里需要根据dayID反查行程ID
	if err != nil {
		return nil, errors.New("无法获取行程信息")
	}

	// 验证所有权
	if itinerary.UserID != userID {
		return nil, errors.New("无权修改此行程")
	}

	// 验证项目
	if item.Title == "" {
		return nil, errors.New("标题不能为空")
	}

	// 设置天ID
	item.ItineraryDayID = dayID

	// 获取当前项目数量
	var day *models.ItineraryDay
	for _, d := range itinerary.Days {
		if d.ID == dayID {
			day = d
			break
		}
	}

	if day == nil {
		return nil, errors.New("天数不存在")
	}

	item.Order = len(day.Items) + 1

	// 保存项目
	itemID, err := s.itineraryRepo.AddItem(item)
	if err != nil {
		return nil, errors.New("添加项目失败: " + err.Error())
	}

	item.ID = itemID
	return item, nil
}

// UpdateItem 更新行程项目
func (s *ItineraryServiceImpl) UpdateItem(itemID int64, userID int64, item *models.ItineraryItem) (*models.ItineraryItem, error) {
	// 获取行程
	itinerary, err := s.itineraryRepo.GetByID(0) // 这里需要根据itemID反查行程ID
	if err != nil {
		return nil, errors.New("无法获取行程信息")
	}

	// 验证所有权
	if itinerary.UserID != userID {
		return nil, errors.New("无权修改此行程")
	}

	// 验证项目
	if item.Title == "" {
		return nil, errors.New("标题不能为空")
	}

	// 设置ID
	item.ID = itemID

	// 保存更新
	err = s.itineraryRepo.UpdateItem(item)
	if err != nil {
		return nil, errors.New("更新项目失败: " + err.Error())
	}

	return item, nil
}

// DeleteItem 删除行程项目
func (s *ItineraryServiceImpl) DeleteItem(itemID int64, userID int64) error {
	// 获取行程
	itinerary, err := s.itineraryRepo.GetByID(0) // 这里需要根据itemID反查行程ID
	if err != nil {
		return errors.New("无法获取行程信息")
	}

	// 验证所有权
	if itinerary.UserID != userID {
		return errors.New("无权修改此行程")
	}

	// 删除项目
	err = s.itineraryRepo.DeleteItem(itemID)
	if err != nil {
		return errors.New("删除项目失败: " + err.Error())
	}

	return nil
}

// GetPublicItineraries 获取公开的行程
func (s *ItineraryServiceImpl) GetPublicItineraries(page, pageSize int) ([]*models.ItineraryResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	itineraries, err := s.itineraryRepo.GetPublicItineraries(pageSize, offset)
	if err != nil {
		return nil, errors.New("获取公开行程失败: " + err.Error())
	}

	var responses []*models.ItineraryResponse
	for _, itinerary := range itineraries {
		response := itinerary.ToItineraryResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}
