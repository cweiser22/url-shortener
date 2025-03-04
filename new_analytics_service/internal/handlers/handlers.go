package handlers

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"new_analytics_service/internal/models"
	"new_analytics_service/internal/service"
)

// used for dependency injection
type AnalyticsHandler struct {
	AnalyticsService service.AnalyticsService
	RedisClient      *redis.Client
}

func (h *AnalyticsHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(models.HealthCheckResponse{Status: "ok"})
}

// Function to publish a visit event
func publishVisit(ctx context.Context, rdb *redis.Client, shortCode string) error {
	msgID, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "visits",
		Values: map[string]interface{}{
			"short_code": shortCode,
		},
	}).Result()

	if err != nil {
		return err
	}

	log.Printf("Visit published for %s and id %s\n", shortCode, msgID)
	return nil
}

// warning: even though this is a GET handler, it is NOT idempotent
// this is actually functionally a POST
// the reason it's handled as a GET is so we can use request mirroring with the /{shortCode}
// / endpoint form the URLs service to record visits
func (h *AnalyticsHandler) UrlVisitHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK")) // Immediately send response
	if f, ok := w.(http.Flusher); ok {
		f.Flush() // Flush response to client
	}

	log.Printf("Visiting %s\n", shortCode)

	// Send Kafka message asynchronously to avoid blocking response
	go func() {
		ctx := context.Background()

		err := publishVisit(ctx, h.RedisClient, shortCode)
		if err != nil {
			log.Printf("Error publishing visit: %v", err)
		}
	}()
}

// make a get url stats endpoint that gets the stats for the url passed as the shortCode
func (h *AnalyticsHandler) UrlStatsHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	stats, err := h.AnalyticsService.GetUrlStats(shortCode)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	log.Printf("Stats for %s: %v\n", shortCode, stats)
	json.NewEncoder(w).Encode(stats)
}
