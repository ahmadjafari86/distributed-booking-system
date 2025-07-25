package kafka

import (
	"context"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type MessageHandler func(msg kafka.Message)

func StartConsumer(brokerURL, groupID string, topic string, handler MessageHandler) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerURL},
		GroupID:     groupID,
		Topic:       topic,
		MaxAttempts: 3,
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		},
	})

	log.Printf("Kafka consumer started for topics %v with GroupID: %s", topic, groupID)

	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v", err)
				// If you want to stop on error, uncomment the next line
				// return
				continue
			}
			log.Printf("Message received from topic %s, partition %d, offset %d, key: %s, value: %s",
				m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			handler(m)
		}
	}()
}
