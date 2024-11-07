package request

type AddressUpdateRequest struct {
	Id          int    `json:"id"`
	City        string `json:"city" validate:"min=3,max=20"`
	Country     string `json:"country" validate:"min=3,max=20"`
	FullAddress string `json:"fullAddress" validate:"min=10,max=100"`
	UserId      string `json:"userId"`
}
