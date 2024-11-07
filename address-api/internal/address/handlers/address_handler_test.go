package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sefikcan/address/address-api/internal/address/dto/request"
	"github.com/sefikcan/address/address-api/internal/address/dto/response"
	"github.com/sefikcan/address/address-api/internal/address/service/mocks"
	"github.com/sefikcan/address/address-api/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestAddressHandler_GetAll(t *testing.T) {
	mockService := mocks.NewAddressService(t)
	handler := NewAddressHandler(mockService)

	app := fiber.New()
	app.Get("/api/v1/addresses", handler.GetAll)

	mockAddresses := []response.AddressResponse{
		{
			Id:          1,
			Country:     "Street 1",
			City:        "City 1",
			FullAddress: "test full",
			UserId:      "1",
		},
		{
			Id:          2,
			Country:     "Street 12",
			City:        "City 21",
			FullAddress: "test full2",
			UserId:      "2",
		},
	}

	pageableAddresses := &common.Pageable[response.AddressResponse]{
		Items:      mockAddresses,
		TotalItems: int64(len(mockAddresses)),
	}

	mockService.On("GetAll", mock.Anything, 0, 10).Return(pageableAddresses, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/addresses", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAddressHandler_Create(t *testing.T) {
	mockService := mocks.NewAddressService(t)
	handler := NewAddressHandler(mockService)
	app := fiber.New()
	app.Post("/api/v1/addresses", handler.Create)

	addressCreateRequest := request.AddressCreateRequest{
		Country:     "Turkey",
		City:        "Istanbul",
		UserId:      "123",
		FullAddress: "123 Main St",
	}
	addressResponse := response.AddressResponse{
		Id:          1,
		Country:     "Turkey",
		City:        "Istanbul",
		UserId:      "123",
		FullAddress: "123 Main St",
	}

	mockService.On("Create", mock.Anything, addressCreateRequest).Return(&addressResponse, nil)

	body, _ := json.Marshal(addressCreateRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/addresses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestAddressHandler_GetById(t *testing.T) {
	mockService := mocks.NewAddressService(t)
	handler := NewAddressHandler(mockService)
	app := fiber.New()
	app.Get("/api/v1/addresses/:id", handler.GetById)

	id := 1
	addressResponse := response.AddressResponse{
		Id:          1,
		Country:     "Turkey",
		City:        "Istanbul",
		UserId:      "123",
		FullAddress: "123 Main St",
	}

	mockService.On("GetById", mock.Anything, id).Return(&addressResponse, nil)

	addressRequest := httptest.NewRequest(http.MethodGet, "/api/v1/addresses/"+strconv.Itoa(id), nil)
	resp, _ := app.Test(addressRequest)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestAddressHandler_Delete(t *testing.T) {
	mockService := mocks.NewAddressService(t)
	handler := NewAddressHandler(mockService)
	app := fiber.New()
	app.Delete("/api/v1/addresses/:id", handler.Delete)

	id := 1

	mockService.On("Delete", mock.Anything, id).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/addresses/"+strconv.Itoa(id), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestAddressHandler_Update(t *testing.T) {
	mockService := mocks.NewAddressService(t)
	handler := NewAddressHandler(mockService)
	app := fiber.New()
	app.Put("/api/v1/addresses/:id", handler.Update)

	id := 1
	addressUpdateRequest := request.AddressUpdateRequest{
		Id:          id,
		Country:     "Turkey",
		City:        "Istanbul",
		UserId:      "123",
		FullAddress: "123 Main St",
	}
	addressResponse := response.AddressResponse{
		Id:          1,
		Country:     "Turkey",
		City:        "Istanbul",
		UserId:      "123",
		FullAddress: "123 Main St",
	}

	mockService.On("Update", mock.Anything, addressUpdateRequest).Return(&addressResponse, nil)

	body, _ := json.Marshal(addressUpdateRequest)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/addresses/"+strconv.Itoa(id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestAddressHandler_Patch(t *testing.T) {
	mockService := mocks.NewAddressService(t)
	handler := NewAddressHandler(mockService)
	app := fiber.New()
	app.Patch("/api/v1/addresses/:id", handler.Patch)

	id := 1
	addressPatchReq := request.AddressPatchRequest{
		Doc: []request.PatchRequest{
			{
				Op:    "replace",
				Path:  "/city",
				Value: "San Francisco",
			},
			{
				Op:    "replace",
				Path:  "/fullAddress",
				Value: "456 Elm St, San Francisco, CA 94101",
			},
		},
	}
	addressResponse := response.AddressResponse{
		Id:          1,
		Country:     "USA",
		City:        "San Francisco",
		FullAddress: "456 Elm St, San Francisco, CA 94101",
		UserId:      "12345",
	}
	mockService.On("Patch", mock.Anything, id, addressPatchReq).Return(&addressResponse, nil)

	body, _ := json.Marshal(addressPatchReq)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/addresses/"+strconv.Itoa(id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}
