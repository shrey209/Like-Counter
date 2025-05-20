package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	log.Printf(" Consumer-%d processing batch of size %d", id, len(batch))

	postCount := make(map[string]int)
	for _, msg := range batch {
		var parsedMsg struct {
			PostID string `json:"post_id"`
		}

		if err := json.Unmarshal(msg.Value, &parsedMsg); err != nil {
			log.Printf(" Consumer-%d: Failed to parse message: %v\n", id, err)
			continue
		}

		postCount[parsedMsg.PostID]++
		log.Printf(" Consumer-%d received post_id: %s\n", id, parsedMsg.PostID)
	}

	type LikePayload struct {
		PostID    string `json:"PostID"`
		LikeCount int    `json:"LikeCount"`
	}

	var payload []LikePayload
	for postID, count := range postCount {
		payload = append(payload, LikePayload{
			PostID:    postID,
			LikeCount: count,
		})
	}
	for _, it := range payload {
		fmt.Println("postid is ", it.PostID)
		fmt.Println("count is ", it.LikeCount)
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Consumer-%d: Failed to marshal payload: %v", id, err)
		return
	}

	resp, err := http.Post("http://localhost:8001/likes/batch", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("❌ Consumer-%d: HTTP request failed: %v", id, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("📬 Consumer-%d: Sent batch to /likes/batch, got status: %s", id, resp.Status)
}
