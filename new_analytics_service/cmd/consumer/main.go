package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"new_analytics_service/internal/repository"
	"new_analytics_service/internal/service"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	streamName    = "visits"       // Redis stream name
	consumerGroup = "visitors_grp" // Redis consumer group
	consumerName  = "consumer_1"   // Unique consumer name (change per instance)
)

type AnalyticsConsumer struct {
	rdb              *redis.Client
	db               *sql.DB
	analyticsService service.AnalyticsService
}

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Println("Skipping local .env...")
	}

	options, err := redis.ParseURL(os.Getenv("REDIS_URI"))

	// Connect to Redis
	rdb := redis.NewClient(options)

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))

	repo := repository.NewAnalyticsRepository(db)
	analyticsService := service.NewAnalyticsService(repo)

	consumer := AnalyticsConsumer{
		db:               db,
		rdb:              rdb,
		analyticsService: analyticsService,
	}

	ctx := context.Background()

	// Ensure consumer group exists (safe to call multiple times)
	consumer.createConsumerGroup(ctx)

	log.Printf("Consumer %s is running...\n", consumerName)

	for {
		// Read messages from Redis Stream
		consumer.readFromStream(ctx)
		time.Sleep(1 * time.Second) // Small delay before next poll
	}
}

// Creates a Redis Stream consumer group if it doesn't exist
func (c *AnalyticsConsumer) createConsumerGroup(ctx context.Context) {
	err := c.rdb.XGroupCreateMkStream(ctx, streamName, consumerGroup, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	log.Printf("Consumer group '%s' is ready.\n", consumerGroup)
}

// Reads messages from the Redis Stream
func (c *AnalyticsConsumer) readFromStream(ctx context.Context) {
	streams, err := c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: consumerName,
		Streams:  []string{streamName, ">"},
		Count:    5,    // Process up to 5 messages at a time
		Block:    5000, // Block for 5 seconds waiting for new messages
	}).Result()

	if err == redis.Nil {
		// No new messages
		return
	} else if err != nil {
		log.Printf("Error reading from stream: %v", err)
		return
	}

	for _, stream := range streams {
		for _, message := range stream.Messages {
			c.processMessage(ctx, message)
		}
	}
}

// Processes and acknowledges each message
func (c *AnalyticsConsumer) processMessage(ctx context.Context, msg redis.XMessage) {

	shortCode := fmt.Sprintf("%v", msg.Values["short_code"])

	log.Printf("Processing message: %v\n", msg.Values)

	err := c.analyticsService.RecordUrlVisit(shortCode)
	if err != nil {
		log.Printf("Failed to record visit for %s.\n", shortCode)
		return
	}
	log.Printf("Recorded visit for %s to DB.", shortCode)

	// Acknowledge message as processed
	_, err = c.rdb.XAck(ctx, streamName, consumerGroup, msg.ID).Result()
	if err != nil {
		log.Printf("Error acknowledging message %s: %v", msg.ID, err)
		return
	}

	log.Printf("Message %s acknowledged\n", msg.ID)
}
