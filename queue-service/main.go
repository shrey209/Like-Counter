package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func createTopic(brokers []string, topic string) error {
	conn, err := kafka.DialContext(context.Background(), "tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.DialContext(context.Background(), "tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to dial controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	return nil
}

func main() {
	brokers := []string{"localhost:9092"}
	topic := "like-events"

	err := createTopic(brokers, topic)
	if err != nil {
		fmt.Println("Failed to create topic:", err)
	} else {
		fmt.Println("Topic created or already exists")
	}
	kp := NewKafkaProducer(brokers, topic)
	defer kp.Close()

	postID := "post-123"

	if err := kp.SendMessage(postID); err != nil {
		fmt.Println("❌ Failed to send message:", err)
		return
	}
	fmt.Println("📤 Message sent!")

	time.Sleep(time.Second)
}
