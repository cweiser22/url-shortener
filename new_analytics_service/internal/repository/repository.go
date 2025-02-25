package repository

import (
	"new_analytics_service/internal/models"
)

type AnalyticsRepository interface {
	SaveUrlVisit(visit models.UrlVisit) error
	GetUrlStats(shortCode string) (models.UrlStats, error)
}
