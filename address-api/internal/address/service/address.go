package service

import (
	"context"
	"encoding/json"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/sefikcan/address-api/internal/address/dto/request"
	"github.com/sefikcan/address-api/internal/address/dto/response"
	"github.com/sefikcan/address-api/internal/address/entity"
	"github.com/sefikcan/address-api/internal/address/event"
	"github.com/sefikcan/address-api/internal/address/mapping"
	"github.com/sefikcan/address-api/internal/address/repository"
	"github.com/sefikcan/address-api/internal/common"
	"github.com/sefikcan/address-api/internal/constants"
	"github.com/sefikcan/address-api/pkg/config"
	"github.com/sefikcan/address-api/pkg/kafka"
	"github.com/sefikcan/address-api/pkg/logger"
)

type AddressService interface {
	Create(ctx context.Context, request request.AddressCreateRequest) (*response.AddressResponse, error)
	Update(ctx context.Context, request request.AddressUpdateRequest) (*response.AddressResponse, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, patchRequest request.AddressPatchRequest) (*response.AddressResponse, error)
	GetById(ctx context.Context, id int) (*response.AddressResponse, error)
	GetAll(ctx context.Context, page, pageSize int) (*common.Pageable[response.AddressResponse], error)
}

type addressService struct {
	cfg               *config.Config
	addressRepository repository.AddressRepository
	logger            logger.Logger
	messageBroker     kafka.Producer
}

func (a addressService) Update(ctx context.Context, request request.AddressUpdateRequest) (*response.AddressResponse, error) {
	currentAddress, err := a.addressRepository.GetById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if request.City != "" {
		currentAddress.City = request.City
	}

	if request.Country != "" {
		currentAddress.Country = request.Country
	}

	if request.FullAddress != "" {
		currentAddress.FullAddress = request.FullAddress
	}

	updatedAddress, err := a.addressRepository.Update(ctx, currentAddress)
	if err != nil {
		return nil, err
	}

	addressEvent := event.AddressEvent{
		EventType:   "AddressUpdated",
		AddressId:   updatedAddress.Id,
		City:        updatedAddress.City,
		Country:     updatedAddress.Country,
		FullAddress: updatedAddress.FullAddress,
		UserId:      updatedAddress.UserId,
	}

	eventBytes, err := json.Marshal(addressEvent)
	if err != nil {
		return nil, err
	}

	eventMessage := string(eventBytes)

	err = a.messageBroker.SendMessage(ctx, constants.KafkaTopics.AddressUpdated, eventMessage)

	mappedResponse := mapping.MapDto(updatedAddress)

	return mappedResponse, nil
}

func (a addressService) GetAll(ctx context.Context, page, pageSize int) (*common.Pageable[response.AddressResponse], error) {
	addresses, err := a.addressRepository.GetAll(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Map entities to DTOs
	var addressDTOs []response.AddressResponse
	for _, addr := range addresses.Items {
		addressDTOs = append(addressDTOs, response.AddressResponse{
			Id:          addr.Id,
			City:        addr.City,
			Country:     addr.Country,
			FullAddress: addr.FullAddress,
			UserId:      addr.UserId,
		})
	}

	return &common.Pageable[response.AddressResponse]{
		Items:       addressDTOs,
		TotalItems:  addresses.TotalItems,
		TotalPages:  addresses.TotalPages,
		CurrentPage: addresses.CurrentPage,
		PageSize:    addresses.PageSize,
	}, nil
}

func (a addressService) Create(ctx context.Context, request request.AddressCreateRequest) (*response.AddressResponse, error) {
	address := mapping.CreateMapEntity(&request)
	resp, err := a.addressRepository.Create(ctx, address)
	if err != nil {
		return nil, err
	}

	addressEvent := event.AddressEvent{
		EventType:   "AddressCreated",
		AddressId:   address.Id,
		City:        request.City,
		Country:     request.Country,
		FullAddress: request.FullAddress,
		UserId:      request.UserId,
	}

	eventBytes, err := json.Marshal(addressEvent)
	if err != nil {
		return nil, err
	}

	eventMessage := string(eventBytes)

	err = a.messageBroker.SendMessage(ctx, constants.KafkaTopics.AddressCreated, eventMessage)

	mappedResponse := mapping.MapDto(resp)

	return mappedResponse, err
}

func (a addressService) Delete(ctx context.Context, id int) error {
	_, err := a.addressRepository.GetById(ctx, id)
	if err != nil {
		return err
	}

	if err = a.addressRepository.Delete(ctx, id); err != nil {
		return err
	}

	addressEvent := event.AddressEvent{
		EventType: "AddressDeleted",
		AddressId: id,
	}

	eventBytes, err := json.Marshal(addressEvent)
	if err != nil {
		return err
	}

	eventMessage := string(eventBytes)

	err = a.messageBroker.SendMessage(ctx, constants.KafkaTopics.AddressDeleted, eventMessage)

	return nil
}

func (a addressService) Patch(ctx context.Context, id int, patchRequest request.AddressPatchRequest) (*response.AddressResponse, error) {
	currentAddress, err := a.addressRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	currentAddressBytes, err := json.Marshal(currentAddress)
	if err != nil {
		return nil, err
	}

	patchRequestBytes, err := json.Marshal(patchRequest.Doc)
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.DecodePatch(patchRequestBytes)
	if err != nil {
		return nil, err
	}

	modifiedAddressBytes, err := patch.Apply(currentAddressBytes)
	if err != nil {
		return nil, err
	}

	var updatedAddress entity.Address
	if err := json.Unmarshal(modifiedAddressBytes, &updatedAddress); err != nil {
		return nil, err
	}

	if _, err := a.addressRepository.Update(ctx, updatedAddress); err != nil {
		return nil, err
	}

	addressEvent := event.AddressEvent{
		EventType:   "AddressUpdated",
		AddressId:   updatedAddress.Id,
		City:        updatedAddress.City,
		Country:     updatedAddress.Country,
		FullAddress: updatedAddress.FullAddress,
		UserId:      updatedAddress.UserId,
	}

	eventBytes, err := json.Marshal(addressEvent)
	if err != nil {
		return nil, err
	}

	eventMessage := string(eventBytes)

	err = a.messageBroker.SendMessage(ctx, constants.KafkaTopics.AddressUpdated, eventMessage)

	mappedResponse := mapping.MapDto(updatedAddress)

	return mappedResponse, nil
}

func (a addressService) GetById(ctx context.Context, id int) (*response.AddressResponse, error) {
	address, err := a.addressRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedResponse := mapping.MapDto(address)

	return mappedResponse, nil
}

func NewAddressService(cfg *config.Config, addressRepository repository.AddressRepository, logger logger.Logger, messageBroker kafka.Producer) AddressService {
	return &addressService{
		cfg:               cfg,
		addressRepository: addressRepository,
		logger:            logger,
		messageBroker:     messageBroker,
	}
}
