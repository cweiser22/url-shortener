package service

import (
	"new_analytics_service/internal/models"
	"new_analytics_service/internal/repository"
)

type AnalyticsService struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(repo repository.AnalyticsRepository) AnalyticsService {
	return AnalyticsService{
		repo: repo,
	}
}

func (s *AnalyticsService) RecordUrlVisit(shortCode string) error {
	visit := models.UrlVisit{
		ShortCode: shortCode,
	}
	err := s.repo.SaveUrlVisit(visit)

	return err
}

func (s *AnalyticsService) GetUrlStats(shortCode string) (models.UrlStats, error) {
	return s.repo.GetUrlStats(shortCode)
}
