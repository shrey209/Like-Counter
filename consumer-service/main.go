package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	broker := "localhost:9092"
	topic := "like-events"
	groupID := "like-consumer-debug-1" // force fresh offset

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer1 := NewLikeConsumer(broker, topic, groupID)
	go consumer1.Start(ctx, 1, 2, 5*time.Second)

	consumer2 := NewLikeConsumer(broker, topic, groupID)
	go consumer2.Start(ctx, 2, 2, 5*time.Second)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("🔻 Shutting down.....")
	cancel()
	consumer1.Reader.Close()
	consumer2.Reader.Close()
}
