package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/like", func(c *gin.Context) {
		type LikeRequest struct {
			PostID string `json:"postId" binding:"required"`
		}

		var req LikeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "postId is required"})
			return
		}

		if err := kp.SendMessage(req.PostID); err != nil {
			fmt.Println("Failed to send message:", err)
			c.JSON(500, gin.H{"status": "error", "message": "Failed to send message"})
			return
		}

		fmt.Println("Message sent:", req.PostID)
		c.JSON(200, gin.H{"status": "success", "postId": req.PostID})
	})

	r.Run(":8000")
}

// commented for debug purpose
// func main() {
// 	brokers := []string{"localhost:9092"}
// 	topic := "like-events"

// 	err := createTopic(brokers, topic)
// 	if err != nil {
// 		fmt.Println("Failed to create topic:", err)
// 	} else {
// 		fmt.Println("Topic created or already exists")
// 	}
// 	kp := NewKafkaProducer(brokers, topic)
// 	defer kp.Close()

// 	postID := "post-123"

// 	for {
// 		if err := kp.SendMessage(postID); err != nil {
// 			fmt.Println(" Failed to send message:", err)
// 			return
// 		}

// 		fmt.Println(" Message sent!")

// 		time.Sleep(time.Second)
// 	}
// }
