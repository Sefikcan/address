package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type Producer interface {
	SendMessage(ctx context.Context, topic, message string) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func (k *KafkaProducer) SendMessage(ctx context.Context, topic, message string) error {
	fmt.Printf("Sending message to Kafka: %s", message) // Mesajı logla

	err := k.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: []byte(message),
	})
	if err != nil {
		log.Printf("Failed to send message to Kafka: %s", err)
		return err
	}

	log.Printf("Message sent to Kafka: %s", message)
	return nil
}

func (k *KafkaProducer) Close() error {
	if err := k.writer.Close(); err != nil {
		log.Printf("Failed to close Kafka writer: %v", err)
		return err
	}

	log.Println("Kafka producer closed successfully")
	return nil
}

func NewKafkaProducer(brokers []string) Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
	}
}
