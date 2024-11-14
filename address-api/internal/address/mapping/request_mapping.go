package mapping

import (
	"github.com/sefikcan/address-api/internal/address/dto/request"
	"github.com/sefikcan/address-api/internal/address/entity"
)

func CreateMapEntity(address *request.AddressCreateRequest) entity.Address {
	return entity.Address{
		City:        address.City,
		Country:     address.Country,
		FullAddress: address.FullAddress,
		UserId:      address.UserId,
	}
}
