package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
	Topic  string
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	})
	return &KafkaProducer{
		Writer: writer,
		Topic:  topic,
	}
}

func (kp *KafkaProducer) SendMessage(postID string) error {
	msgValue, err := json.Marshal(map[string]string{"post_id": postID})
	if err != nil {
		return err
	}

	message := kafka.Message{
		Value: msgValue,
		Time:  time.Now(),
	}
	return kp.Writer.WriteMessages(context.Background(), message)
}

func (kp *KafkaProducer) Close() error {
	return kp.Writer.Close()
}
