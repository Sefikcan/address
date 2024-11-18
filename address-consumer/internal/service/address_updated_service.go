package service

import (
	"context"
	"github.com/sefikcan/address-consumer/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type AddressUpdatedService struct {
	logger logger.Logger
}

func NewAddressUpdatedService(logger logger.Logger) *AddressUpdatedService {
	return &AddressUpdatedService{
		logger: logger,
	}
}

func (s *AddressUpdatedService) ProcessMessage(ctx context.Context, msg kafka.Message) error {
	s.logger.Infof("Address Updated being processed: %s", string(msg.Value))

	return nil
}
