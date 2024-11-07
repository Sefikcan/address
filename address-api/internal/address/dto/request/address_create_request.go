package request

type AddressCreateRequest struct {
	City        string `json:"city" validate:"required,min=3,max=20"`
	Country     string `json:"country" validate:"required,min=3,max=20"`
	FullAddress string `json:"fullAddress" validate:"required,min=10,max=100"`
	UserId      string `json:"userId" validate:"required"`
}
