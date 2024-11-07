package service

import (
	"context"
	"github.com/sefikcan/address/address-api/internal/address/dto/request"
	"github.com/sefikcan/address/address-api/internal/address/entity"
	"github.com/sefikcan/address/address-api/internal/address/repository/mocks"
	mocks2 "github.com/sefikcan/address/address-api/internal/address/service/mocks"
	"github.com/sefikcan/address/address-api/internal/common"
	"github.com/sefikcan/address/address-api/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAddressService_Create(t *testing.T) {
	mockRepo := new(mocks.AddressRepository)
	mockLogger := new(mocks2.Logger)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(entity.Address{
		Id:          1,
		City:        "Test City",
		Country:     "Test Country",
		FullAddress: "123 Test St",
		UserId:      "1",
	}, nil)

	addressService := NewAddressService(&config.Config{}, mockRepo, mockLogger)

	createReq := request.AddressCreateRequest{
		City:        "Test City",
		Country:     "Test Country",
		FullAddress: "123 Test St",
		UserId:      "1",
	}

	// Test Create Address
	resp, err := addressService.Create(context.Background(), createReq)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Id)
	assert.Equal(t, "Test City", resp.City)
	assert.Equal(t, "Test Country", resp.Country)
	assert.Equal(t, "123 Test St", resp.FullAddress)

	// Verify the mock interactions
	mockRepo.AssertExpectations(t)
}

func TestAddressService_GetAll(t *testing.T) {
	mockRepo := new(mocks.AddressRepository)
	mockLogger := new(mocks2.Logger)
	mockRepo.On("GetAll", mock.Anything, mock.Anything, mock.Anything).Return(&common.Pageable[entity.Address]{
		Items: []entity.Address{
			{Id: 1, City: "Test City", Country: "Test Country", FullAddress: "123 Test St", UserId: "1"},
		},
		TotalItems:  1,
		TotalPages:  1,
		CurrentPage: 1,
		PageSize:    10,
	}, nil)

	addressService := NewAddressService(&config.Config{}, mockRepo, mockLogger)

	resp, err := addressService.GetAll(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, 1, resp.Items[0].Id)

	mockRepo.AssertExpectations(t)
}

func TestAddressService_GetById(t *testing.T) {
	mockRepo := new(mocks.AddressRepository)
	mockLogger := new(mocks2.Logger)
	mockRepo.On("GetById", mock.Anything, mock.Anything).Return(entity.Address{
		Id:          1,
		City:        "Test City",
		Country:     "Test Country",
		FullAddress: "123 Test St",
		UserId:      "1",
	}, nil)

	addressService := NewAddressService(&config.Config{}, mockRepo, mockLogger)
	resp, err := addressService.GetById(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Id)
	assert.Equal(t, "Test City", resp.City)
	assert.Equal(t, "Test Country", resp.Country)
	assert.Equal(t, "123 Test St", resp.FullAddress)

	mockRepo.AssertExpectations(t)
}

func TestAddressService_Delete(t *testing.T) {
	mockRepo := new(mocks.AddressRepository)
	mockLogger := new(mocks2.Logger)
	mockRepo.On("GetById", mock.Anything, mock.Anything).Return(entity.Address{
		Id:          1,
		City:        "City",
		Country:     "Country",
		FullAddress: "123 St",
		UserId:      "1",
	}, nil)
	mockRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)

	addressService := NewAddressService(&config.Config{}, mockRepo, mockLogger)

	err := addressService.Delete(context.Background(), 1)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAddressService_Update(t *testing.T) {
	mockRepo := new(mocks.AddressRepository)
	mockLogger := new(mocks2.Logger)
	mockRepo.On("GetById", mock.Anything, mock.Anything).Return(entity.Address{
		Id:          1,
		City:        "Old City",
		Country:     "Old Country",
		FullAddress: "123 Old St",
		UserId:      "1",
	}, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(entity.Address{
		Id:          1,
		City:        "New City",
		Country:     "New Country",
		FullAddress: "123 New St",
		UserId:      "1",
	}, nil)

	addressService := NewAddressService(&config.Config{}, mockRepo, mockLogger)

	updateReq := request.AddressUpdateRequest{
		Id:          1,
		City:        "New City",
		Country:     "New Country",
		FullAddress: "123 New St",
	}

	resp, err := addressService.Update(context.Background(), updateReq)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Id)
	assert.Equal(t, "New City", resp.City)
	assert.Equal(t, "New Country", resp.Country)
	assert.Equal(t, "123 New St", resp.FullAddress)

	mockRepo.AssertExpectations(t)
}

func TestAddressService_Patch(t *testing.T) {
	mockRepo := new(mocks.AddressRepository)
	mockLogger := new(mocks2.Logger)
	mockRepo.On("GetById", mock.Anything, mock.Anything).Return(entity.Address{
		Id:          1,
		City:        "Old City",
		Country:     "Old Country",
		FullAddress: "123 Old St",
		UserId:      "1",
	}, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(entity.Address{
		Id:          1,
		City:        "New City",
		Country:     "New Country",
		FullAddress: "123 New St",
		UserId:      "1",
	}, nil)

	addressService := NewAddressService(&config.Config{}, mockRepo, mockLogger)

	patchReq := request.AddressPatchRequest{
		Doc: []request.PatchRequest{
			{
				Op:    "replace",
				Path:  "/city",
				Value: "New City",
			},
			{
				Op:    "replace",
				Path:  "/country",
				Value: "New Country",
			},
		},
	}

	resp, err := addressService.Patch(context.Background(), 1, patchReq)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Id)
	assert.Equal(t, "New City", resp.City)

	mockRepo.AssertExpectations(t)
}
