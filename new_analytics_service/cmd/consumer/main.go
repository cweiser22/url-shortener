package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"new_analytics_service/internal/repository"
	"new_analytics_service/internal/service"
	"os"
	"os/signal"
	"strings"
)

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	service service.AnalyticsService
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (h consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Receive messages until session or claim is closed
	for message := range claim.Messages() {
		fmt.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, partition = %d, offset = %d\n",
			string(message.Value),
			message.Timestamp,
			message.Topic,
			message.Partition,
			message.Offset,
		)

		shortCode := string(message.Value)
		err := h.service.RecordUrlVisit(shortCode)
		if err != nil {
			log.Printf("Failed to record visit for %s\n", shortCode)
			log.Println(err)
		}

		// Mark the message as processed
		session.MarkMessage(message, "")
	}
	return nil
}

func main() {
	godotenv.Load(".env.local")
	// Kafka configuration
	config := sarama.NewConfig()
	// Set a matching Kafka version if needed. For example:
	// config.Version = sarama.V2_8_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Broker address and topic
	//brokers := []string{"kafka:9092"}
	brokerList := os.Getenv("KAFKA_BROKER_URLS")
	brokers := strings.Split(brokerList, ",")
	groupID := "visit-consumer-group"
	topics := []string{"visit"}

	// Create new consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Fatalf("Error closing consumer group: %v", err)
		}
	}()

	// Trap SIGINT to trigger a graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	go func() {
		<-sigterm
		cancel()
	}()

	connString := os.Getenv("POSTGRES_URI")

	DB, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Could not connect to DB", err)
	}

	repo := repository.NewAnalyticsRepository(DB)

	// Consume messages in a loop
	handler := consumerGroupHandler{
		service: service.NewAnalyticsService(repo),
	}
	for {
		err := consumerGroup.Consume(ctx, topics, handler)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
		}
		// Check if context was cancelled, signaling a graceful shutdown
		if ctx.Err() != nil {
			return
		}
	}
}
