package service

import (
	"context"
	"github.com/sefikcan/address-consumer/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type AddressCreatedService struct {
	logger logger.Logger
}

func NewAddressCreatedService(logger logger.Logger) *AddressCreatedService {
	return &AddressCreatedService{
		logger: logger,
	}
}

func (s *AddressCreatedService) ProcessMessage(ctx context.Context, msg kafka.Message) error {
	s.logger.Infof("Address Created being processed: %s", string(msg.Value))

	return nil
}
