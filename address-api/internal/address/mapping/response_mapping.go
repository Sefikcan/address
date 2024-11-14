package mapping

import (
	"github.com/sefikcan/address-api/internal/address/dto/response"
	"github.com/sefikcan/address-api/internal/address/entity"
)

func MapDto(a entity.Address) *response.AddressResponse {
	return &response.AddressResponse{
		Id:          a.Id,
		City:        a.City,
		Country:     a.Country,
		FullAddress: a.FullAddress,
		UserId:      a.UserId,
	}
}
