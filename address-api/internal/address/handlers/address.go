package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sefikcan/address-api/internal/address/dto/request"
	"github.com/sefikcan/address-api/internal/address/service"
	"strconv"
)

type AddressHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Patch(c *fiber.Ctx) error
	GetAllV2(c *fiber.Ctx) error
}

type addressHandler struct {
	addressService service.AddressService
}

// GetAll godoc
// @Summary Get all addresses
// @Description Get all addresses with pagination
// @Tags addresses
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {array} response.AddressResponse
// @Router /api/v1/addresses [get]
func (a addressHandler) GetAll(c *fiber.Ctx) error {
	page := c.Query("page", "0")
	size := c.Query("size", "10")

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		pageNumber = 0
	}
	pageSize, err := strconv.Atoi(size)
	if err != nil {
		pageSize = 10
	}

	addresses, err := a.addressService.GetAll(c.Context(), pageNumber, pageSize)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to retrieve addresses")
	}

	return c.JSON(addresses)
}

// GetById godoc
// @Summary Get an address by ID
// @Description Retrieve an address by its ID
// @Tags addresses
// @Param id path int true "Address ID"
// @Success 200 {object} response.AddressResponse
// @Router /api/v1/addresses/{id} [get]
func (a addressHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Id")
	}

	currentAddress, err := a.addressService.GetById(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(currentAddress)
}

// Delete godoc
// @Summary Delete an address
// @Description Delete an address by its ID
// @Tags addresses
// @Param id path int true "Address ID"
// @Success 204
// @Router /api/v1/addresses/{id} [delete]
func (a addressHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Id")
	}

	err = a.addressService.Delete(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Update godoc
// @Summary Update an address
// @Description Update an address by its ID
// @Tags addresses
// @Param id path int true "Address ID"
// @Param address body request.AddressUpdateRequest true "Address update payload"
// @Success 200 {object} response.AddressResponse
// @Router /api/v1/addresses/{id} [put]
func (a addressHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Id")
	}

	address := request.AddressUpdateRequest{}
	if err := c.BodyParser(&address); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	address.Id = id
	updatedAddress, err := a.addressService.Update(c.Context(), address)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(updatedAddress)
}

// Create godoc
// @Summary Create a new address
// @Description Create a new address entry
// @Tags addresses
// @Param address body request.AddressCreateRequest true "Address creation payload"
// @Success 201 {object} response.AddressResponse
// @Router /api/v1/addresses [post]
func (a addressHandler) Create(c *fiber.Ctx) error {
	var address request.AddressCreateRequest
	if err := c.BodyParser(&address); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	response, err := a.addressService.Create(c.Context(), address)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// Patch godoc
// @Summary Patch an address
// @Description Patch (partial update) an address by its ID
// @Tags addresses
// @Param id path int true "Address ID"
// @Param address body request.AddressPatchRequest true "Address patch payload"
// @Success 200 {object} response.AddressResponse
// @Router /api/v1/addresses/{id} [patch]
func (a addressHandler) Patch(c *fiber.Ctx) error {
	id := c.Params("id")
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Id")
	}

	var patchRequest request.AddressPatchRequest
	if err := c.BodyParser(&patchRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	updatedAddress, err := a.addressService.Patch(c.Context(), convertedId, patchRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(updatedAddress)
}

// GetAllV2 godoc
// @Summary Get all addresses
// @Description Get all addresses with pagination
// @Tags addresses
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {array} response.AddressResponse
// @Router /api/v2/addresses [get]
func (a addressHandler) GetAllV2(c *fiber.Ctx) error {
	page := c.Query("page", "0")
	size := c.Query("size", "10")

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		pageNumber = 0
	}
	pageSize, err := strconv.Atoi(size)
	if err != nil {
		pageSize = 10
	}

	addresses, err := a.addressService.GetAll(c.Context(), pageNumber, pageSize)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to retrieve addresses")
	}

	return c.JSON(addresses)
}

func NewAddressHandler(addressService service.AddressService) AddressHandler {
	return &addressHandler{
		addressService: addressService,
	}
}
