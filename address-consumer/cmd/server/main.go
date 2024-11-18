package main

import (
	"context"
	"github.com/sefikcan/address-consumer/internal/constants"
	"github.com/sefikcan/address-consumer/internal/consumer"
	"github.com/sefikcan/address-consumer/internal/service"
	"github.com/sefikcan/address-consumer/pkg/config"
	"github.com/sefikcan/address-consumer/pkg/logger"
	"sync"
)

func main() {
	cfg := config.NewConfig()
	log := logger.NewLogger(cfg)
	log.InitLogger()

	services := map[string]consumer.BusinessLogic{
		constants.KafkaTopics.AddressCreated: service.NewAddressCreatedService(log),
		constants.KafkaTopics.AddressDeleted: service.NewAddressDeletedService(log),
		constants.KafkaTopics.AddressUpdated: service.NewAddressUpdatedService(log),
	}

	var wg sync.WaitGroup
	for topic, businessLogic := range services {
		wg.Add(1)
		go func(topic string, logic consumer.BusinessLogic) {
			defer wg.Done()
			kafkaConsumer, err := consumer.NewKafkaConsumer(cfg, log, topic, logic)
			if err != nil {
				log.Errorf("Consumer could not be started (%s): %v", topic, err)
				return
			}

			if err := kafkaConsumer.Start(context.Background()); err != nil {
				log.Errorf("An error occurred while running Consumer (%s): %v", topic, err)
			}
		}(topic, businessLogic)
	}

	wg.Wait()
}
