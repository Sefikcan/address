package consumer

import (
	"context"
	"github.com/sefikcan/address-consumer/pkg/config"
	"github.com/sefikcan/address-consumer/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type BusinessLogic interface {
	ProcessMessage(ctx context.Context, msg kafka.Message) error
}

type KafkaConsumer struct {
	reader  *kafka.Reader
	logger  logger.Logger
	topic   string
	handler BusinessLogic
}

func NewKafkaConsumer(cfg *config.Config, logger logger.Logger, topic string, handler BusinessLogic) (*KafkaConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
		GroupID: "1",
	})

	return &KafkaConsumer{
		reader:  reader,
		logger:  logger,
		topic:   topic,
		handler: handler,
	}, nil
}

func (kc *KafkaConsumer) Start(ctx context.Context) error {
	kc.logger.Info("Kafka Consumer started. Topic: %s", kc.topic)

	for {
		msg, err := kc.reader.ReadMessage(ctx)
		if err != nil {
			kc.logger.Errorf("Message could not be read (%s): %v", kc.topic, string(msg.Value))
			continue
		}

		kc.logger.Infof("Message received (%s): %s", kc.topic, string(msg.Value))

		if err := kc.handler.ProcessMessage(ctx, msg); err != nil {
			kc.logger.Errorf("Message could not be processed (%s): %v", kc.topic, err)
		}
	}
}
