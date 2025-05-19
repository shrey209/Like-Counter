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
		MinBytes:        1,
		MaxBytes:        10e3,
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
		StartOffset:     kafka.FirstOffset,
	})

	return &LikeConsumer{
		Reader: r,
	}
}
func (c *LikeConsumer) Start(ctx context.Context, id int, batchSize int, maxWait time.Duration) {
	log.Printf("🔄 Consumer-%d started\n", id)

	batch := make([]kafka.Message, 0, batchSize)
	timer := time.NewTimer(maxWait)

	for {
		select {
		case <-ctx.Done():
			log.Printf("🛑 Consumer-%d shutting down\n", id)
			return

		default:
			m, err := c.Reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("❌ Consumer-%d error: %v\n", id, err)
				continue
			}

			batch = append(batch, m)

			if len(batch) >= batchSize {
				c.processBatch(batch, id)
				batch = batch[:0]
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(maxWait)
			}
		}

		select {
		case <-timer.C:
			if len(batch) > 0 {
				c.processBatch(batch, id)
				batch = batch[:0]
			}
			timer.Reset(maxWait)
		default:

		}
	}
}

func (c *LikeConsumer) processBatch(batch []kafka.Message, id int) {
	log.Printf("📦 Consumer-%d processing batch of size %d", id, len(batch))

	postCount := make(map[string]int)

	for _, msg := range batch {
		postID := string(msg.Value)
		postCount[postID]++
		log.Printf("📥 Consumer-%d received: %s\n", id, postID)
	}

	log.Printf("🧾 Consumer-%d batch result: %+v\n", id, postCount)
}
