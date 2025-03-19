package services

import (
	"context"
	"errors"
	"strings"

	"traveler_agent/models"
	"traveler_agent/repositories"
	"traveler_agent/utils"
)

// AuthService 认证服务接口
type AuthService interface {
	Register(ctx context.Context, request *models.CreateUserRequest) (*models.UserResponse, string, error)
	Login(ctx context.Context, email, password string) (*models.UserResponse, string, error)
	GetUserByID(ctx context.Context, id int64) (*models.UserResponse, error)
}

// AuthServiceImpl 认证服务实现
type AuthServiceImpl struct {
	userRepo repositories.UserRepository
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &AuthServiceImpl{userRepo: userRepo}
}

// Register 用户注册
func (s *AuthServiceImpl) Register(ctx context.Context, request *models.CreateUserRequest) (*models.UserResponse, string, error) {
	// 检查邮箱是否已存在
	existingUser, _ := s.userRepo.GetByEmail(ctx, request.Email)
	if existingUser != nil {
		return nil, "", errors.New("邮箱已被注册")
	}

	// 检查用户名是否已存在
	existingUser, _ = s.userRepo.GetByUsername(ctx, request.Username)
	if existingUser != nil {
		return nil, "", errors.New("用户名已被使用")
	}

	// 密码哈希
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, "", errors.New("密码处理失败")
	}

	// 设置用户为活跃状态
	isActive := true

	// 创建用户对象
	user := &models.User{
		Username:  request.Username,
		Email:     strings.ToLower(request.Email), // 转为小写
		Password:  hashedPassword,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Phone:     request.Phone,
		IsActive:  &isActive,
	}

	// 保存用户
	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, "", errors.New("用户创建失败: " + err.Error())
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(userID, user.Username)
	if err != nil {
		return nil, "", errors.New("令牌生成失败")
	}

	// 返回UserResponse和令牌
	return &models.UserResponse{
		ID:        userID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, token, nil
}

// Login 用户登录
func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (*models.UserResponse, string, error) {
	// 查询用户
	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return nil, "", errors.New("用户不存在")
	}

	// 检查用户是否存在
	if user == nil {
		return nil, "", errors.New("用户不存在")
	}

	// 检查用户状态
	if user.IsActive == nil || !*user.IsActive {
		return nil, "", errors.New("账户已禁用")
	}

	// 验证密码
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("密码错误")
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", errors.New("令牌生成失败")
	}

	// 返回UserResponse和令牌
	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	}, token, nil
}

// GetUserByID 获取用户信息
func (s *AuthServiceImpl) GetUserByID(ctx context.Context, id int64) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user == nil {
		return nil, errors.New("用户不存在")
	}

	if user.IsActive == nil || !*user.IsActive {
		return nil, errors.New("账户已禁用")
	}

	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	}, nil
}
