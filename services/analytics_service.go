package services

import (
	"blog/repositories"
)

type AnalyticsService interface {
	CountViews() (int64, error)
}

type analyticsService struct {
	analyticsRepo *repositories.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo *repositories.AnalyticsRepository) AnalyticsService {
	return &analyticsService{
		analyticsRepo: analyticsRepo,
	}
}

func (s *analyticsService) CountViews() (int64, error) {
	return s.analyticsRepo.CountViews()
}
