package main

import (
	"database/sql"
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	app2 "new_analytics_service/internal/app"
	"new_analytics_service/internal/handlers"
	"new_analytics_service/internal/repository"
	"new_analytics_service/internal/service"
	"os"
	"strings"
)

// NewProducer initializes a new Kafka producer
func NewProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// Create Kafka producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	log.Println("Kafka producer initialized")
	return producer, nil
}

func main() {
	godotenv.Load(".env.local")
	router := http.NewServeMux()

	brokerList := os.Getenv("KAFKA_BROKER_URLS")
	brokers := strings.Split(brokerList, ",")
	producer, err := NewProducer(brokers)

	if err != nil {
		log.Fatal("Failed to connect to Kafka with URL %s, %s", brokerList, err)
	}

	app := app2.ProducerApp{
		KafkaProducer: producer,
	}

	DB, err := sql.Open("postgres", "postgres://admin:admin101@postgres:5432/analytics_db?sslmode=disable")
	if err != nil {
		log.Fatal("Could not connect to DB", err)
	}

	repo := repository.NewAnalyticsRepository(DB)

	ah := handlers.AnalyticsHandler{
		Producer:         app.KafkaProducer,
		AnalyticsService: service.NewAnalyticsService(repo),
	}

	router.HandleFunc("/", ah.HealthCheckHandler)
	router.HandleFunc("/{shortCode}", ah.UrlVisitHandler)
	router.HandleFunc("/api/analytics/{shortCode}/stats", ah.UrlStatsHandler)
	log.Fatal(http.ListenAndServe(":8002", router))
}
