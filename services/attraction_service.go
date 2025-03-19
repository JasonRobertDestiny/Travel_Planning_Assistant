package services

import (
	"context"
	"errors"

	"traveler_agent/models"
	"traveler_agent/repositories"
)

// AttractionService 景点服务接口
type AttractionService interface {
	GetAttractionByID(ctx context.Context, id int64) (*models.Attraction, error)
	SearchAttractions(ctx context.Context, params *models.SearchAttractionsRequest) ([]*models.AttractionListItem, int, error)
	GetPopularAttractions(ctx context.Context, limit int) ([]*models.AttractionListItem, error)
	GetAttractionsByCategory(ctx context.Context, category string, limit int) ([]*models.AttractionListItem, error)
	GetAttractionsByCountry(ctx context.Context, country string, limit int) ([]*models.AttractionListItem, error)
	GetAttractionsByCityOrCountry(ctx context.Context, destination string, limit int) ([]*models.AttractionListItem, error)
}

// AttractionServiceImpl 景点服务实现
type AttractionServiceImpl struct {
	attractionRepo repositories.AttractionRepository
}

// NewAttractionService 创建景点服务
func NewAttractionService(attractionRepo repositories.AttractionRepository) AttractionService {
	return &AttractionServiceImpl{
		attractionRepo: attractionRepo,
	}
}

// GetAttractionByID 根据ID获取景点详情
func (s *AttractionServiceImpl) GetAttractionByID(ctx context.Context, id int64) (*models.Attraction, error) {
	attraction, err := s.attractionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("获取景点失败: " + err.Error())
	}
	if attraction == nil {
		return nil, errors.New("景点不存在")
	}
	return attraction, nil
}

// SearchAttractions 搜索景点
func (s *AttractionServiceImpl) SearchAttractions(ctx context.Context, params *models.SearchAttractionsRequest) ([]*models.AttractionListItem, int, error) {
	attractions, total, err := s.attractionRepo.Search(ctx, params)
	if err != nil {
		return nil, 0, errors.New("搜索景点失败: " + err.Error())
	}
	return attractions, total, nil
}

// GetPopularAttractions 获取热门景点
func (s *AttractionServiceImpl) GetPopularAttractions(ctx context.Context, limit int) ([]*models.AttractionListItem, error) {
	attractions, err := s.attractionRepo.GetPopular(ctx, limit)
	if err != nil {
		return nil, errors.New("获取热门景点失败: " + err.Error())
	}
	return attractions, nil
}

// GetAttractionsByCategory 根据分类获取景点
func (s *AttractionServiceImpl) GetAttractionsByCategory(ctx context.Context, category string, limit int) ([]*models.AttractionListItem, error) {
	attractions, err := s.attractionRepo.GetByCategory(ctx, category, limit)
	if err != nil {
		return nil, errors.New("获取分类景点失败: " + err.Error())
	}
	return attractions, nil
}

// GetAttractionsByCountry 根据国家获取景点
func (s *AttractionServiceImpl) GetAttractionsByCountry(ctx context.Context, country string, limit int) ([]*models.AttractionListItem, error) {
	attractions, err := s.attractionRepo.GetByCountry(ctx, country, limit)
	if err != nil {
		return nil, errors.New("获取国家景点失败: " + err.Error())
	}
	return attractions, nil
}

// GetAttractionsByCityOrCountry 根据城市或国家获取景点
func (s *AttractionServiceImpl) GetAttractionsByCityOrCountry(ctx context.Context, destination string, limit int) ([]*models.AttractionListItem, error) {
	attractions, err := s.attractionRepo.GetByCityOrCountry(ctx, destination, limit)
	if err != nil {
		return nil, errors.New("获取目的地景点失败: " + err.Error())
	}
	return attractions, nil
}
