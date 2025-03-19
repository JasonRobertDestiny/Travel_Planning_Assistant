package services

import (
	"context"
	"errors"
	"time"

	"traveler_agent/models"
	"traveler_agent/repositories"
	"traveler_agent/utils"
)

// UserService 用户服务接口
type UserService interface {
	GetUserByID(ctx context.Context, id int64) (*models.UserResponse, error)
	UpdateUserProfile(ctx context.Context, id int64, user *models.User) error
	UpdatePassword(ctx context.Context, id int64, currentPassword, newPassword string) error
	DeleteUser(ctx context.Context, id int64) error
	GetUserPreferences(ctx context.Context, userID int64) (*models.UserPreference, error)
	SaveUserPreferences(ctx context.Context, userID int64, pref *models.UserPreference) error
	GetTravelPreferences(ctx context.Context, userID int64) (*models.UserPreference, error)
	SaveTravelPreferences(ctx context.Context, userID int64, pref *models.UserPreference) error
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImpl{userRepo: userRepo}
}

// GetUserByID 获取用户信息
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id int64) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("用户不存在")
	}

	if user.IsActive == nil || !*user.IsActive {
		return nil, errors.New("账户已禁用")
	}

	response := user.ToUserResponse()
	return &response, nil
}

// UpdateUserProfile 更新用户资料
func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, id int64, updateData *models.User) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("用户不存在")
	}

	// 检查邮箱是否已被其他用户使用
	if updateData.Email != user.Email {
		existingUser, _ := s.userRepo.GetByEmail(ctx, updateData.Email)
		if existingUser != nil && existingUser.ID != id {
			return errors.New("邮箱已被其他账户使用")
		}
	}

	// 检查用户名是否已被其他用户使用
	if updateData.Username != user.Username {
		existingUser, _ := s.userRepo.GetByUsername(ctx, updateData.Username)
		if existingUser != nil && existingUser.ID != id {
			return errors.New("用户名已被其他账户使用")
		}
	}

	// 更新用户信息
	user.Username = updateData.Username
	user.Email = updateData.Email
	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName
	user.Phone = updateData.Phone
	user.Avatar = updateData.Avatar
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

// UpdatePassword 更新用户密码
func (s *UserServiceImpl) UpdatePassword(ctx context.Context, id int64, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("用户不存在")
	}

	// 验证当前密码
	if !utils.CheckPasswordHash(currentPassword, user.Password) {
		return errors.New("当前密码不正确")
	}

	// 哈希新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码处理失败")
	}

	return s.userRepo.UpdatePassword(ctx, id, hashedPassword)
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("用户不存在")
	}

	return s.userRepo.Delete(ctx, id)
}

// GetUserPreferences 获取用户偏好设置
func (s *UserServiceImpl) GetUserPreferences(ctx context.Context, userID int64) (*models.UserPreference, error) {
	preferences, err := s.userRepo.GetUserPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 如果没有偏好设置，返回默认值
	if preferences == nil {
		return &models.UserPreference{
			UserID:              userID,
			Language:            "zh-CN",
			Currency:            "CNY",
			NotificationEnabled: true,
			Theme:               "light",
		}, nil
	}

	return preferences, nil
}

// SaveUserPreferences 保存用户偏好设置
func (s *UserServiceImpl) SaveUserPreferences(ctx context.Context, userID int64, pref *models.UserPreference) error {
	pref.UserID = userID
	return s.userRepo.SaveUserPreferences(ctx, pref)
}

// GetTravelPreferences 获取用户旅行偏好设置
func (s *UserServiceImpl) GetTravelPreferences(ctx context.Context, userID int64) (*models.UserPreference, error) {
	preferences, err := s.userRepo.GetTravelPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 如果没有旅行偏好设置，返回默认值
	if preferences == nil {
		return &models.UserPreference{
			UserID:      userID,
			TravelStyle: "balanced",
			BudgetLevel: "medium",
			PreferredTags: []string{
				"景点", "美食", "文化",
			},
			TransportPrefer: "public",
		}, nil
	}

	return preferences, nil
}

// SaveTravelPreferences 保存用户旅行偏好设置
func (s *UserServiceImpl) SaveTravelPreferences(ctx context.Context, userID int64, pref *models.UserPreference) error {
	pref.UserID = userID
	return s.userRepo.SaveTravelPreferences(ctx, pref)
}
