package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"new_analytics_service/internal/handlers"
	"new_analytics_service/internal/repository"
	"new_analytics_service/internal/service"
	"os"
)

func main() {
	godotenv.Load(".env.local")
	router := http.NewServeMux()

	//app := app2.ProducerApp{}

	connString := os.Getenv("POSTGRES_URI")
	log.Println(connString)
	DB, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatal("Could not connect to DB", err)
	}

	repo := repository.NewAnalyticsRepository(DB)

	// set up DI
	ah := handlers.AnalyticsHandler{
		AnalyticsService: service.NewAnalyticsService(repo),
	}

	// TODO: maybe use Chi because this is a bit messy
	apiPrefix := "/analytics/api/v1"
	router.HandleFunc(fmt.Sprintf("GET %s/{shortCode}/stats", apiPrefix), ah.UrlStatsHandler)
	router.HandleFunc(fmt.Sprintf("GET %s/{shortCode}", apiPrefix), ah.UrlVisitHandler)
	router.HandleFunc(fmt.Sprintf("GET %s/", apiPrefix), ah.HealthCheckHandler)

	log.Fatal(http.ListenAndServe(":80", router))
}
