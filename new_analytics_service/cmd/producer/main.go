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

	brokers := []string{"kafka:9092"}
	producer, err := NewProducer(brokers)

	if err != nil {
		log.Fatal("Failed to connect to Kafka.")
		return
	}

	app := app2.ProducerApp{
		KafkaProducer: producer,
	}

	ah := handlers.AnalyticsHandler{
		Producer: app.KafkaProducer,
	}

	router.HandleFunc("/", ah.HealthCheckHandler)
	router.HandleFunc("/{shortCode}", ah.UrlVisitHandler)
	log.Fatal(http.ListenAndServe(":8002", router))
}
