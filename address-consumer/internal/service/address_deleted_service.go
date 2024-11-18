package service

import (
	"context"
	"github.com/sefikcan/address-consumer/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type AddressDeletedService struct {
	logger logger.Logger
}

func NewAddressDeletedService(logger logger.Logger) *AddressDeletedService {
	return &AddressDeletedService{
		logger: logger,
	}
}

func (s *AddressDeletedService) ProcessMessage(ctx context.Context, msg kafka.Message) error {
	s.logger.Infof("Address Deleted being processed: %s", string(msg.Value))

	return nil
}
