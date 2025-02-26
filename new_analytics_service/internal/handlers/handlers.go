package handlers

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"net/http"
	"new_analytics_service/internal/models"
	"new_analytics_service/internal/service"
)

type AnalyticsHandler struct {
	Producer         sarama.SyncProducer
	AnalyticsService service.AnalyticsService
}

func (h *AnalyticsHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(models.HealthCheckResponse{Status: "ok"})
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

	// Send Kafka message asynchronously to avoid blocking response
	go func() {
		msg := &sarama.ProducerMessage{Topic: "visit", Value: sarama.StringEncoder(shortCode)}
		_, _, err := h.Producer.SendMessage(msg)
		if err != nil {
			log.Println("Kafka error:", err)
		}
		log.Println("Sent message:", shortCode)
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
	json.NewEncoder(w).Encode(stats)
}
