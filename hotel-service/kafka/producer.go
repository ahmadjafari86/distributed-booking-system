package kafka

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func InitProducer(brokerURL string) {
	writer = &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("Kafka producer initialized.")
}

func ProduceMessage(topic string, key, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to write message to topic %s: %v", topic, err)
		return err
	}
	log.Printf("Message sent to topic %s, key: %s, value: %s", topic, string(key), string(value))
	return nil
}
