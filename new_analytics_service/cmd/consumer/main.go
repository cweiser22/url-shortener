package main

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Println("Skipping local .env...")
	}

	for {
		time.Sleep(10 * time.Second)
		fmt.Println("Consumer placeholder...")
	}
}
