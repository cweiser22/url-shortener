package models

import "time"

type UrlVisit struct {
	ShortCode string
}

type UrlStats struct {
	TotalVisits  int `json:"totalVisits"`
	LastAccessed time.Time
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}
