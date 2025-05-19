package main

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type LikeConsumer struct {
	Reader *kafka.Reader
}

func NewLikeConsumer(broker, topic, groupID string) *LikeConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{broker},
		GroupID:         groupID,
		Topic:           topic,
		MinBytes:        1,    // 1B
		MaxBytes:        10e3, // 10KB
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
	})

	return &LikeConsumer{
		Reader: r,
	}
}

func (c *LikeConsumer) Start(ctx context.Context, id int) {
	log.Printf("🔄 Consumer-%d started\n", id)

	for {
		m, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("❌ Consumer-%d error: %v\n", id, err)
			break
		}

		// TODO: Add your processing logic here
		log.Printf("📥 Consumer-%d received: %s\n", id, string(m.Value))
	}
}
