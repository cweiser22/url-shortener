package repository

import (
	"database/sql"
	"errors"

	"new_analytics_service/internal/models"
)

type postgresAnalyticsRepository struct {
	db *sql.DB
}

// NewAnalyticsRepository returns a concrete implementation of AnalyticsRepository.
func NewAnalyticsRepository(db *sql.DB) AnalyticsRepository {
	return &postgresAnalyticsRepository{db: db}
}

func (r *postgresAnalyticsRepository) SaveUrlVisit(visit models.UrlVisit) error {
	// Upsert query: insert a new row if it doesn't exist, or update visit_count & last_access if it does
	query := `
        INSERT INTO url_access_log (short_code, visit_count, last_access)
        VALUES ($1, 1, NOW())
        ON CONFLICT (short_code)
        DO UPDATE 
          SET visit_count = url_access_log.visit_count + 1,
              last_access = NOW();
    `

	_, err := r.db.Exec(query, visit.ShortCode)
	return err
}

func (r *postgresAnalyticsRepository) GetUrlStats(shortCode string) (models.UrlStats, error) {
	// Retrieve the total visits and last_access for the given short code
	query := `
        SELECT visit_count, last_access
        FROM url_access_log
        WHERE short_code = $1
    `

	row := r.db.QueryRow(query, shortCode)

	var stats models.UrlStats
	err := row.Scan(&stats.TotalVisits, &stats.LastAccessed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// You can decide whether to return zero-values or a custom "not found" error
			return stats, nil
		}
		return stats, err
	}

	return stats, nil
}
