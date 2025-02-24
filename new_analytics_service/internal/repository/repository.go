package repository

type AnalyticsRepository interface {
	SaveUrlVisit(visit models.UrlVisit) error
	GetUrlStats(shortCode string)
}
