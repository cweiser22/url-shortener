package main

import (
	"github.com/IBM/sarama"
	"log"
	"net/http"
	app2 "new_analytics_service/internal/app"
	"new_analytics_service/internal/handlers"
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
	router := http.NewServeMux()

	brokers := []string{"localhost:9092"}
	producer, err := NewProducer(brokers)

	app := app2.ProducerApp{
		KafkaProducer: producer,
	}

	router.HandleFunc("/", app.HealthCheckHandler)
	router.HandleFunc("/{shortCode}", handlers.UrlVisitHandler)
	err := http.ListenAndServe(":8002", router)
	if err != nil {
		log.Fatal("HTTP producer failed.")
		return
	}
}
