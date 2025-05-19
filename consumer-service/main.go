package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	broker := "localhost:9092"
	topic := "like-events"
	groupID := "like-consumer-group"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := NewLikeConsumer(broker, topic, groupID)

	// Start 2 consumers in the same group
	for i := 1; i <= 2; i++ {
		go consumer.Start(ctx, i)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("🔻 Shutting down...")
	cancel()
	consumer.Reader.Close()
}
