package repository

import (
	"database/sql"
	"new_analytics_service/internal/models"
)

type AnalyticsRepository interface {
	SaveUrlVisit(visit models.UrlVisit) error
	GetUrlStats(shortCode string) (models.UrlStats, error)
}

func NewPostgresAnalyticsRepository(db *sql.DB) AnalyticsRepository {
	return &postgresAnalyticsRepository{db: db}
}
