package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"new_analytics_service/internal/app"
	"new_analytics_service/internal/models"
)

func (a *app.ProducerApp) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(models.HealthCheckResponse{Status: "ok"})
}

// warning: even though this is a GET handler, it is NOT idempotent
// this is actually functionally a POST
// the reason it's handled as a GET is so we can use request mirroring with the /{shortCode}
// / endpoint form the URLs service to record visits
func (a *app.ProducerApp) UrlVisitHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("UrlVisitHandler called for %s", r.PathValue("shortCode"))
	w.WriteHeader(http.StatusOK)
}
